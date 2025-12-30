package documentstore

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
)

type StoreDump struct {
	Collections map[string]CollectionDump `json:"collections"`
}

type CollectionDump struct {
	Config    CollectionConfig `json:"config"`
	Documents []Document       `json:"documents"`
}

func (s *Store) Dump() ([]byte, error) {
	s.logger.Info("starting store dump")

	dump := StoreDump{
		Collections: make(map[string]CollectionDump),
	}

	for name, coll := range s.collections {
		dump.Collections[name] = CollectionDump{
			Config:    coll.cfg,
			Documents: coll.List(),
		}
	}

	data, err := json.Marshal(dump)
	if err != nil {
		s.logger.Error("failed to marshal store dump", "error", err)
		return nil, err
	}

	s.logger.Info("store dump completed", "collections_count", len(s.collections))
	return data, nil
}

func NewStoreFromDump(dump []byte) (*Store, error) {
	var storeDump StoreDump
	if err := json.Unmarshal(dump, &storeDump); err != nil {
		return nil, err
	}

	store := NewStore()

	for name, collDump := range storeDump.Collections {
		coll, err := store.CreateCollection(name, &collDump.Config)
		if err != nil {
			return nil, err
		}

		for _, doc := range collDump.Documents {
			if err := coll.Put(doc); err != nil {
				return nil, err
			}
		}
	}

	store.logger.Info("store loaded from dump", "collections_count", len(storeDump.Collections))
	return store, nil
}

func (s *Store) DumpToFile(filename string) error {
	s.logger.Info("dumping store to file", "filename", filename)

	data, err := s.Dump()
	if err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		s.logger.Error("failed to create dump file", "filename", filename, "error", err)
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	if _, err := writer.Write(data); err != nil {
		s.logger.Error("failed to write dump to file", "filename", filename, "error", err)
		return err
	}

	if err := writer.Flush(); err != nil {
		s.logger.Error("failed to flush writer", "filename", filename, "error", err)
		return err
	}

	if err := file.Sync(); err != nil {
		s.logger.Error("failed to sync file", "filename", filename, "error", err)
		return err
	}

	s.logger.Info("store dumped to file successfully", "filename", filename)
	return nil
}

func NewStoreFromFile(filename string) (*Store, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	store, err := NewStoreFromDump(data)
	if err != nil {
		return nil, err
	}

	store.logger.Info("store loaded from file", "filename", filename)
	return store, nil
}
