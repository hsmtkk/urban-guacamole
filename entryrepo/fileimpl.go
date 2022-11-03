package entryrepo

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/hsmtkk/urban-guacamole/entry"
	"go.uber.org/zap"
)

type fileImpl struct {
	sugar       *zap.SugaredLogger
	fileName    string
	internalMap map[int64]entry.Entry
}

func NewFileImpl(sugar *zap.SugaredLogger, fileName string) (EntryRepo, error) {
	internalMap, err := load(fileName)
	if err != nil {
		return nil, err
	}
	return &fileImpl{sugar, fileName, internalMap}, nil
}

func load(fileName string) (map[int64]entry.Entry, error) {
	decoded := map[int64]entry.Entry{}
	bs, err := os.ReadFile(fileName)
	if err != nil {
		return decoded, nil
	}
	if err := json.Unmarshal(bs, &decoded); err != nil {
		return nil, fmt.Errorf("json.Unmarshal failed; %w", err)
	}
	return decoded, nil
}

func (impl *fileImpl) Scan() ([]entry.Entry, error) {
	impl.sugar.Info("Scan")
	entries := []entry.Entry{}
	for _, e := range impl.internalMap {
		entries = append(entries, e)
	}
	return entries, nil
}

func (impl *fileImpl) Save(e entry.Entry) error {
	impl.sugar.Info("Save")
	impl.internalMap[e.ID] = e
	return impl.save()
}

func (impl *fileImpl) Get(id int64) (entry.Entry, error) {
	impl.sugar.Info("Get")
	e, ok := impl.internalMap[id]
	if ok {
		return e, nil
	}
	return entry.Entry{}, fmt.Errorf("ID %d was not found", id)
}

func (impl *fileImpl) Delete(id int64) error {
	impl.sugar.Info("Delete")
	delete(impl.internalMap, id)
	return impl.save()
}

func (impl *fileImpl) save() error {
	bs, err := json.Marshal(impl.internalMap)
	if err != nil {
		return fmt.Errorf("json.Marshal failed; %w", err)
	}
	if err := os.WriteFile(impl.fileName, bs, 0644); err != nil {
		return fmt.Errorf("os.WriteFile failed; %w", err)
	}
	return nil
}
