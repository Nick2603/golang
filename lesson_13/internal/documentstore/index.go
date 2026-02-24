package documentstore

import (
	"sync"

	"github.com/google/btree"
)

type IndexEntry struct {
	Value      string
	PrimaryKey string
}

func (a IndexEntry) Less(b btree.Item) bool {
	return a.Value < b.(IndexEntry).Value
}

type Index struct {
	fieldName string
	tree      *btree.BTree
	mu        sync.RWMutex
}

func newIndex(fieldName string) *Index {
	return &Index{
		fieldName: fieldName,
		tree:      btree.New(32), // degree 32 is a good default
	}
}

func (idx *Index) insert(value string, primaryKey string) {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	entry := IndexEntry{
		Value:      value,
		PrimaryKey: primaryKey,
	}
	idx.tree.ReplaceOrInsert(entry)
}

func (idx *Index) delete(value string, primaryKey string) {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	entry := IndexEntry{
		Value:      value,
		PrimaryKey: primaryKey,
	}
	idx.tree.Delete(entry)
}

type QueryParams struct {
	Desc     bool
	MinValue *string
	MaxValue *string
}

func (idx *Index) query(params QueryParams) []IndexEntry {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	var results []IndexEntry

	iterFunc := func(item btree.Item) bool {
		entry := item.(IndexEntry)

		if params.MinValue != nil && entry.Value < *params.MinValue {
			return true // continue iteration
		}

		if params.MaxValue != nil && entry.Value > *params.MaxValue {
			return false
		}

		results = append(results, entry)
		return true
	}

	if params.Desc {
		idx.tree.Descend(iterFunc)
	} else {
		idx.tree.Ascend(iterFunc)
	}

	return results
}

func (idx *Index) len() int {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	return idx.tree.Len()
}
