package source

import (
	"FeedCraft/internal/constant"
	"fmt"
	"sync"
)

var (
	registry = make(map[constant.SourceType]SourceFactory)
	lock     = new(sync.RWMutex)
)

// Register a new source factory. Panics if type is already registered.
func Register(sourceType constant.SourceType, factory SourceFactory) {
	lock.Lock()
	defer lock.Unlock()
	if _, exists := registry[sourceType]; exists {
		panic(fmt.Sprintf("source factory for type '%s' already registered", sourceType))
	}
	registry[sourceType] = factory
}

// Get a source factory by type.
func Get(sourceType constant.SourceType) (SourceFactory, error) {
	lock.RLock()
	defer lock.RUnlock()
	factory, ok := registry[sourceType]
	if !ok {
		return nil, fmt.Errorf("no source factory registered for type '%s'", sourceType)
	}
	return factory, nil
}
