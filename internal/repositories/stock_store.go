package repositories

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log/slog"
	"os"
	"sort"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/parquet-go/parquet-go"
)

const (
	r2AccessKeyIDEnv        = "R2_ACCESS_KEY_ID"
	r2SecretAccessKeyEnv    = "R2_SECRET_ACCESS_KEY"
	r2EndpointEnv           = "R2_ENDPOINT"
	r2BucketEnv             = "R2_BUCKET_NAME"
	goEnv                   = "GO_ENV"
	metaDir                 = "meta"
	stockDir                = "stock"
	featureDir              = "feature"
	stockListFileNamePrefix = "stock_list"
)

const maxLRUCapacity = 100

type StockStore struct {
	r2             *s3.Client
	logger         *slog.Logger
	bucket         string
	runtimeEnv     string
	inMemStockList map[string]StockInfo
	stockCache     *LRUCache[string, []StockPrice]
	featureCache   *LRUCache[string, []StockFeature]
}

func NewStockStore(ctx context.Context, runtimeEnv string, logger *slog.Logger) (*StockStore, error) {
	// read env variables
	var r2AccessKeyID, r2SecretAccessKey, r2Endpoint, r2BucketName string
	if r2AccessKeyID = os.Getenv(r2AccessKeyIDEnv); len(r2AccessKeyID) == 0 {
		return nil, fmt.Errorf("%s shouldn't be an empty string", r2AccessKeyIDEnv)
	}
	if r2SecretAccessKey = os.Getenv(r2SecretAccessKeyEnv); len(r2SecretAccessKey) == 0 {
		return nil, fmt.Errorf("%s shouldn't be an empty string", r2SecretAccessKeyEnv)
	}
	if r2Endpoint = os.Getenv(r2EndpointEnv); len(r2Endpoint) == 0 {
		return nil, fmt.Errorf("%s shouldn't be an empty string", r2EndpointEnv)
	}
	if r2BucketName = os.Getenv(r2BucketEnv); len(r2BucketName) == 0 {
		return nil, fmt.Errorf("%s shouldn't be an empty string", r2BucketEnv)
	}
	// init connection to r2 storage
	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion("auto"),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				r2AccessKeyID, r2SecretAccessKey, "",
			),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load R2 storage config: %w", err)
	}
	r2 := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(r2Endpoint)
	})
	if r2 == nil {
		return nil, fmt.Errorf("failed to create r2 client")
	}
	// load in-memory stock list
	stockList, err := loadStockList(ctx, r2, r2BucketName, runtimeEnv, logger)
	if err != nil {
		return nil, err
	}
	logger.Info("stock list loaded", "count", len(stockList))
	stockCache := NewLRUCache[string, []StockPrice](maxLRUCapacity)
	featureCache := NewLRUCache[string, []StockFeature](maxLRUCapacity)
	return &StockStore{
		r2:             r2,
		runtimeEnv:     runtimeEnv,
		bucket:         r2BucketName,
		inMemStockList: stockList,
		logger:         logger,
		stockCache:     stockCache,
		featureCache:   featureCache,
	}, nil
}

func (s *StockStore) ValidateTicker(ticker string) bool {
	if _, ok := s.inMemStockList[ticker]; !ok {
		return false
	}
	return true
}

func (s *StockStore) GetStockList() map[string]StockInfo {
	return s.inMemStockList
}

func (s *StockStore) GetStockInfo(ticker string) (*StockInfo, error) {
	stockInfo, ok := s.inMemStockList[ticker]
	if !ok {
		return nil, fmt.Errorf("invalid ticker name: %q", ticker)
	}
	return &stockInfo, nil
}

func (s *StockStore) ReadStock(ctx context.Context, ticker string) ([]StockPrice, error) {
	// ensure the ticker is in the in-mem list
	if _, ok := s.inMemStockList[ticker]; !ok {
		return nil, fmt.Errorf("invalid ticker name: %q", ticker)
	}
	cachedPrices, found := s.stockCache.Get(ticker)
	if found {
		return cachedPrices, nil
	}
	stockPrices, err := readStock(ctx, s.r2, s.bucket, ticker)
	if err != nil {
		return nil, err
	}
	s.stockCache.Put(ticker, stockPrices)
	return stockPrices, nil
}

func (s *StockStore) ReadFeatures(ctx context.Context, ticker string) ([]StockFeature, error) {
	if _, ok := s.inMemStockList[ticker]; !ok {
		return nil, fmt.Errorf("invalid ticker name: %q", ticker)
	}
	cachedFeatures, found := s.featureCache.Get(ticker)
	if found {
		return cachedFeatures, nil
	}
	stockFeatures, err := readFeatures(ctx, s.r2, s.bucket, ticker)
	if err != nil {
		return nil, err
	}
	s.featureCache.Put(ticker, stockFeatures)
	return stockFeatures, nil
}

// ================== private methods ==================

// load stock list from r2 storage
func loadStockList(ctx context.Context, r2 *s3.Client, bucket string, runtimeEnv string, logger *slog.Logger) (map[string]StockInfo, error) {
	fileKey := fmt.Sprintf("%s/%s_%s.csv", metaDir, stockListFileNamePrefix, runtimeEnv)
	object, err := r2.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileKey),
	})
	if err != nil {
		return nil, fmt.Errorf("get stock-list object %q: %w", fileKey, err)
	}
	defer object.Body.Close()
	return readStockListCSV(object.Body, logger)
}

// parse stock list csv file to GO struct
func readStockListCSV(source io.Reader, logger *slog.Logger) (map[string]StockInfo, error) {
	funcLogger := logger.With("path", "internal/repositories/stock_store/readStockListCSV")
	r := csv.NewReader(source)
	r.TrimLeadingSpace = true
	r.FieldsPerRecord = 8 // ref: StockInfo

	header, err := r.Read()
	if err != nil {
		return nil, fmt.Errorf("read stock-list header: %w", err)
	}
	expected := []string{
		"ts_code", "company_name", "sector", "sub_industry",
		"headquarters", "date_added", "cik", "founded"}
	for i, name := range expected {
		if header[i] != name {
			return nil, fmt.Errorf(
				"unexpected CSV column %d: got %q, expect %q",
				i, header[i], name)
		}
	}

	stockList := make(map[string]StockInfo)
	for line := 2; ; line++ {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read stock-list row %d: %w", line, err)
		}
		stockInfo := StockInfo{
			TSCode:       record[0],
			CompanyName:  record[1],
			Sector:       record[2],
			SubIndustry:  record[3],
			Headquarters: record[4],
			DateAdded:    record[5],
			CIK:          record[6],
			Founded:      record[7],
		}
		if stockInfo.TSCode == "" {
			return nil, fmt.Errorf("stock-list row %d has empty ts_code", line)
		}
		if _, exists := stockList[stockInfo.TSCode]; exists {
			funcLogger.Warn("duplicate ts_code", "ts_code", stockInfo.TSCode)
			continue
		}
		stockList[stockInfo.TSCode] = stockInfo
	}
	if len(stockList) == 0 {
		return nil, fmt.Errorf("Failed to load stock list, current list is empty")
	}
	return stockList, nil
}

// parse stock prices parquet file to GO struct
func readStock(ctx context.Context, r2 *s3.Client, bucket, ticker string) ([]StockPrice, error) {
	fileKey := fmt.Sprintf("%s/%s.parquet", stockDir, ticker)
	object, err := r2.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileKey),
	})
	if err != nil {
		return nil, fmt.Errorf("get stock object %q: %w", fileKey, err)
	}
	defer object.Body.Close()
	payload, err := io.ReadAll(object.Body)
	if err != nil {
		return nil, fmt.Errorf("read stock object %q: %w", fileKey, err)
	}
	reader := parquet.NewGenericReader[StockPrice](bytes.NewReader(payload))
	defer reader.Close()

	var prices []StockPrice
	buffer := make([]StockPrice, 1024)
	for {
		count, err := reader.Read(buffer)
		prices = append(prices, buffer[:count]...)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("decode parquet %q: %w", fileKey, err)
		}
	}
	sort.Slice(prices, func(i, j int) bool {
		return prices[i].TradeDate < prices[j].TradeDate
	})
	return prices, nil
}

// parse stock feature parquet file to GO struct
func readFeatures(ctx context.Context, r2 *s3.Client, bucket, ticker string) ([]StockFeature, error) {
	fileKey := fmt.Sprintf("%s/%s.parquet", featureDir, ticker)
	object, err := r2.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileKey),
	})
	if err != nil {
		return nil, fmt.Errorf("get feature object %q: %w", fileKey, err)
	}
	defer object.Body.Close()
	payload, err := io.ReadAll(object.Body)
	if err != nil {
		return nil, fmt.Errorf("read feature object %q: %w", fileKey, err)
	}
	reader := parquet.NewGenericReader[StockFeature](bytes.NewReader(payload))
	defer reader.Close()

	var features []StockFeature
	buffer := make([]StockFeature, 1024)
	for {
		count, err := reader.Read(buffer)
		features = append(features, buffer[:count]...)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("decode stock feature parquet %q: %w", fileKey, err)
		}
	}
	sort.Slice(features, func(i, j int) bool {
		return features[i].TradeDate < features[j].TradeDate
	})
	return features, nil
}
