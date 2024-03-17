package persistence

import (
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"
)

type CrudStorage[E any] interface {
	Add(c context.Context, id uuid.UUID, entity E) error
	GetById(c context.Context, id uuid.UUID) (E, error)
	Delete(c context.Context, id uuid.UUID) error
	GetAll(c context.Context) []E
}

type Storage[E any] struct {
	entities map[uuid.UUID]E
	lock     sync.RWMutex
}

func New[E any]() *Storage[E] {
	return &Storage[E]{
		entities: make(map[uuid.UUID]E),
	}
}

func (s *Storage[E]) Add(c context.Context, id uuid.UUID, entity E) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	_, ok := s.entities[id]
	if ok {
		return errors.New("already exists")
	}

	s.entities[id] = entity
	return nil
}

func (s *Storage[E]) GetById(c context.Context, id uuid.UUID) (E, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	var zero E

	entity, ok := s.entities[id]
	if !ok {
		return zero, errors.New("not found")
	}
	return entity, nil
}

func (s *Storage[E]) Delete(c context.Context, id uuid.UUID) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	delete(s.entities, id)
	return nil
}

func (s *Storage[E]) GetAll(c context.Context) []E {
	s.lock.RLock()
	defer s.lock.RUnlock()

	ents := make([]E, 0, len(s.entities))
	for _, e := range s.entities {
		ents = append(ents, e)
	}

	return ents
}
