package storage

import (
	"archive/zip"
	"bytes"
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/teryble09/17.07.2025/internal/archiver/model"
	"github.com/teryble09/17.07.2025/internal/archiver/repository"
)

type InMemoryStorage struct {
	storage map[model.TaskID]Task
	maxUrl  int
	sync.RWMutex
}

type Task struct {
	urls          []model.Url
	archive       *bytes.Buffer
	archiveWriter *zip.Writer
	mutex         *sync.RWMutex
}

func NewInMemoryStorage(maxUrl int) *InMemoryStorage {
	return &InMemoryStorage{
		storage: make(map[model.TaskID]Task, 1000),
		maxUrl:  maxUrl,
	}
}

func (s *InMemoryStorage) CreateTask() model.TaskID {
	s.Lock()
	defer s.Unlock()

	id := model.TaskID{Id: uuid.NewString()}
	for {
		if _, ok := s.storage[id]; ok {
			id = model.TaskID{Id: uuid.NewString()}
		} else {
			buf := bytes.NewBuffer([]byte{})
			s.storage[id] = Task{urls: nil, archive: buf, archiveWriter: zip.NewWriter(buf), mutex: &sync.RWMutex{}}
			return id
		}
	}
}

func (s *InMemoryStorage) AddURL(id model.TaskID, url string) error {
	s.Lock()
	defer s.Unlock()

	task, ok := s.storage[id]
	if !ok {
		return repository.ErrTaskNotFound
	}

	task.mutex.Lock()
	defer task.mutex.Unlock()

	if len(task.urls) == s.maxUrl {
		return repository.ErrMaximumTaskNumberReached
	}

	task.urls = append(task.urls, model.Url{Address: url, Status: model.Waiting})
	s.storage[id] = task
	return nil
}

func (s *InMemoryStorage) Status(id model.TaskID) ([]model.Url, error) {
	s.RLock()

	task, ok := s.storage[id]
	if !ok {
		return nil, repository.ErrTaskNotFound
	}

	s.RUnlock()

	task.mutex.RLock()
	defer task.mutex.RUnlock()

	return task.urls, nil
}

func (s *InMemoryStorage) LoadArchive(id model.TaskID) ([]byte, error) {
	s.RLock()

	task, ok := s.storage[id]
	if !ok {
		return nil, repository.ErrTaskNotFound
	}

	s.RUnlock()

	task.mutex.RLock()
	defer task.mutex.RUnlock()

	for _, url := range task.urls {
		if url.Status == model.Waiting || url.Status == model.Loaded {
			return nil, repository.ErrArchiveNotReady
		}
	}

	return task.archive.Bytes(), nil
}

func (s *InMemoryStorage) WriteToArchive(id model.TaskID, filename []byte, file []byte) error {
	s.RLock()

	task, ok := s.storage[id]
	if !ok {
		return repository.ErrTaskNotFound
	}

	s.RUnlock()

	task.mutex.Lock()
	defer task.mutex.Unlock()

	f, err := task.archiveWriter.Create(string(filename))
	if err != nil {
		return errors.Join(repository.ErrFailedWrite, err)
	}

	_, err = f.Write(file)
	if err != nil {
		return errors.Join(repository.ErrFailedWrite, err)
	}

	return nil
}
