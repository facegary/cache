package cache

import (
	"time"
)

var generalStorage Storage

func init() {
	generalStorage = NewStorage(time.Second * 10)
}

// Добавить объект
/* При отрицательном ttl кеш не истекает */
func Set(key string, value any, ttl time.Duration) {
	generalStorage.Set(key, value, ttl)
}

// Получить значение
/* Вернёт false, если ttl истек */
func Get(key string) (value any, isExist bool) {
	return generalStorage.Get(key)
}

// Удалить значение
func Delete(key string) {
	generalStorage.Delete(key)
}

// Запустить очистку
func Cleanup() {
	generalStorage.Cleanup()
}

// Проверить существует ли кеш по ключу
func Exists(key string) bool {
	return generalStorage.Exists(key)
}

// Получить хранимые ключи
func Keys() []any {
	return generalStorage.Keys()
}

// Установить задержку между очисткой
func SetCleanupDelay(d time.Duration) {
	generalStorage.SetCleanupDelay(d)
}
