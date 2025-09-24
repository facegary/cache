# Golang __cache__ package
Runtime cache
## Features
 - Базируется на sync.Map
## Usage/Examples
~~~golang
// Управление общедоступным кешем
cache.Set("someKey", "anyValue", time.Minute)
cache.Get("spmeKey")
cache.Delete("someKey")
cache.Cleanup()
cache.Exists("someKey")
cache.Keys()
cache.SetCleanupDelay(time.Second*5)

// Объявление локального хранилища с идентичным интерфейсом
cstorage := cache.NewStorage(time.Second)
~~~
По дефолту очистка проводится каждые 10 секунд. Это значение можно изменить методом _SetCleanupDelay_