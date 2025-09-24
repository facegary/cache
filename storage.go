package cache

import (
	"sync"
	"time"
)

type Storage interface {
	// Добавить объект в хранилище
	/* При ttl <= 0 кеш не истекает */
	Set(key any, value any, ttl time.Duration)
	// Получить значение из хранилища
	/* Вернёт false, если время кеша истекло */
	Get(key any) (value any, ok bool)
	// Удалить значение в хранилище
	Delete(key any)
	// Запуск очистки хранилища
	Cleanup()
	// Проверка на наличие данных в хранилище
	Exists(key any) bool
	// Получить хранимые ключи
	Keys() []any
	// Установить задержку между очисткой
	SetCleanupDelay(d time.Duration)
}

type storage struct {
	bucket     sync.Map
	cleanerJob ijob
}

// Инициализировать локальный Storage
/* cleanupDelay == 0 отключает периодическую чистку */
func NewStorage(cleanupDelay time.Duration) Storage {
	s := &storage{}
	s.runCleaner(cleanupDelay)
	return s
}

func (s *storage) Set(key, value any, ttl time.Duration) {
	obj := newObject(value, ttl)
	s.bucket.Store(key, obj)
}

func (s *storage) Get(key any) (any, bool) {
	unit, ok := s.bucket.Load(key)
	if !ok {
		return nil, false
	}
	obj := unit.(iObject)
	if obj.IsExpired() {
		return nil, false
	}
	return obj.Value(), true
}

func (s *storage) Delete(key any) {
	s.bucket.Delete(key)
}

func (s *storage) Cleanup() {
	s.bucket.Range(func(key, value any) bool {
		go func(key, value any) {
			if obj := value.(iObject); obj.IsExpired() {
				s.Delete(key)
			}
		}(key, value)
		return true
	})
}

func (s *storage) Exists(key any) bool {
	unit, ok := s.bucket.Load(key)
	if ok {
		obj := unit.(iObject)
		if obj.IsExpired() {
			return false
		}
	}
	return ok
}

func (s *storage) Keys() []any {
	var result []any
	swg := new(sync.WaitGroup)
	s.bucket.Range(func(key, value any) bool {
		swg.Add(1)
		go func(key, value any) {
			defer swg.Done()
			if obj := value.(iObject); !obj.IsExpired() {
				result = append(result, key)
			}
		}(key, value)
		return true
	})
	swg.Wait()
	return result
}

func (s *storage) runCleaner(d time.Duration) {
	if s.cleanerJob != nil {
		s.cleanerJob.Stop()
	}
	if d > 0 {
		s.cleanerJob = newJob(func(t time.Time) { s.Cleanup() }, d)
	}
}

func (s *storage) SetCleanupDelay(d time.Duration) {
	s.runCleaner(d)
}
