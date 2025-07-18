package storage

import (
	"sync"

	"github.com/google/uuid"
	"github.com/teryble09/17.07.2025/internal/archiver/model"
)

type InMemoryStorage struct {
	storage map[model.TaskID]model.Task
	maxUrl  int
	sync.RWMutex
}

func NewInMemoryStorage(maxUrl int) *InMemoryStorage {
	return &InMemoryStorage{storage: make(map[model.TaskID]model.Task, 1000), maxUrl: maxUrl}
}

func (s *InMemoryStorage) CreateTask() model.TaskID {
	s.Lock()
	defer s.Unlock()
	id := model.TaskID{Id: uuid.NewString()}
	for {
		if _, ok := s.storage[id]; ok {
			id = model.TaskID{Id: uuid.NewString()}
		} else {
			s.storage[id] = model.Task{}
			return id
		}
	}
}

func (s *InMemoryStorage) AddURL(id model.TaskID, url string) error {
	s.Lock()
	defer s.Unlock()
	task, ok := s.storage[id]
	if !ok {
		return model.ErrTaskNotFound
	}
	task.Urls = append(task.Urls, model.Url{Address: url, Status: model.Waiting})
	s.storage[id] = task
	return nil
}

func (s *InMemoryStorage) Status(id model.TaskID) ([]model.Url, error) {
	s.RLock()
	defer s.RUnlock()
	task, ok := s.storage[id]
	if !ok {
		return nil, model.ErrTaskNotFound
	}
	return task.Urls, nil
}
