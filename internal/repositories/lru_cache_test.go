package repositories

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLRUCache(t *testing.T) {

	t.Parallel()
	type operation struct {
		name      string
		action    string // "put" or "get"
		key       string
		value     int
		wantValue int
		wantFound bool
	}

	tests := []struct {
		name     string
		capacity int
		ops      []operation
	}{
		{
			name:     "returns cache miss for unknown key",
			capacity: 2,
			ops: []operation{
				{name: "get missing key", action: "get", key: "missing", wantFound: false},
			},
		},
		{
			name:     "returns stored value",
			capacity: 2,
			ops: []operation{
				{name: "put value", action: "put", key: "a", value: 1},
				{name: "get stored value", action: "get", key: "a", wantValue: 1, wantFound: true},
			},
		},
		{
			name:     "evicts least recently used item",
			capacity: 2,
			ops: []operation{
				{name: "put a", action: "put", key: "a", value: 1},
				{name: "put b", action: "put", key: "b", value: 2},
				// Make "a" most recently used, leaving "b" as LRU.
				{name: "access a", action: "get", key: "a", wantValue: 1, wantFound: true},
				{name: "put c", action: "put", key: "c", value: 3},
				{name: "b was evicted", action: "get", key: "b", wantFound: false},
				{name: "a remains", action: "get", key: "a", wantValue: 1, wantFound: true},
				{name: "c remains", action: "get", key: "c", wantValue: 3, wantFound: true},
			},
		},
		{
			name:     "updates an existing key",
			capacity: 1,
			ops: []operation{
				{name: "put initial value", action: "put", key: "a", value: 1},
				{name: "replace value", action: "put", key: "a", value: 2},
				{name: "read replacement", action: "get", key: "a", wantValue: 2, wantFound: true},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := NewLRUCache[string, int](tt.capacity)
			for _, op := range tt.ops {
				t.Run(op.name, func(t *testing.T) {
					switch op.action {
					case "put":
						cache.Put(op.key, op.value)
					case "get":
						got, found := cache.Get(op.key)
						assert.Equal(t, op.wantFound, found)
						if found {
							assert.Equal(t, op.wantValue, got)
						}
					default:
						t.Fatalf("unknown action: %q", op.action)
					}
				})
			}
		})
	}
}
