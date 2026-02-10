package main

import (
	"fmt"
	"sync"
	"time"
)

/*
ЗАДАНИЕ:
Реализуйте in-memory key-value хранилище с TTL, аналогичное Redis.

ТРЕБОВАНИЯ:
1. Реализуйте структуру, удовлетворяющую интерфейсу KVCache
2. Хранилище должно быть безопасным для конкурентного доступа (concurrent safe)
3. Реализуйте автоматическую очистку просроченных ключей в фоне (background cleanup)
4. Учитывайте различные edge cases:
   - TTL = 0 (ключ должен быть вечным)
   - Отрицательный TTL (можно считать как 0 или вернуть ошибку)
   - Повторная установка значения для существующего ключа (должен обновиться TTL)
   - Конкурентные операции чтения/записи

ПОДСКАЗКИ:
1. Для хранения данных можно использовать sync.Map или map + sync.RWMutex
2. Для фоновой очистки используйте горутину + ticker
3. Для хранения TTL удобно хранить время истечения (time.Time), а не duration
4. Не забудьте корректно останавливать фоновые горутины при вызове Stop()

БОНУСНЫЕ ЗАДАЧИ (по желанию):
1. Добавьте метод Flush() для очистки всех ключей - сделал
2. Добавьте метод Keys() для получения списка всех не просроченных ключей - сделал
3. Реализуйте возможность обновления TTL для существующего ключа
4. Добавьте метрики: количество ключей, количество операций get/set

ТЕСТИРОВАНИЕ:
1. Запустите main() после реализации
2. Убедитесь, что через 1 секунду часть ключей еще доступна
3. Убедитесь, что через 3 секунды все ключи с TTL=1s удалены
4. Проверьте работу с вечными ключами (TTL=0)
5. Проверьте конкурентный доступ (тест в конце main())
*/

// KVCache - интерфейс для key-value хранилища с TTL
type KVCache interface {
	// Run запускает фоновые процессы кэша (например, очистку просроченных ключей)
	Run()

	// Set устанавливает значение по ключу с указанным TTL (временем жизни)
	// Если TTL <= 0, ключ должен быть вечным (или иметь очень большое TTL)
	Set(key string, value any, ttl time.Duration)

	// Get возвращает значение по ключу и флаг, указывающий, найден ли ключ
	// Если ключ просрочен, должен вернуться false, и ключ должен быть удален
	Get(key string) (any, bool)

	// Delete удаляет ключ из кэша
	Delete(key string)

	// Stop останавливает фоновые процессы кэша
	Stop()

	// Flush очищает все данные
	Flush()

	// Возвращает список ключей с TTL>0
	Keys() []string
}

type dataValue struct {
	Value any
	TTL   time.Time
	inf   bool
}

type cache struct {
	data   sync.Map
	ticker time.Ticker
}

const (
	FirstReadTimeout  = 1 * time.Second
	SecondReadTimeout = 2 * time.Second
)

var (
	CacheValues = [7]int{500, 200, 2000, 5000, 200, 1500, 5000}
)

func (c *cache) Run() {
	go func() {
		for {
			_ = <-c.ticker.C // единственный канал во всем задании, и тот пишет в никуда :(
			c.data.Range(func(key, value any) bool {
				v := value.(dataValue)
				// Если время TTL раньше, чем сейчас, то TTL истек, и если ключ не бесконечный, то удаляем
				if v.TTL.Before(time.Now()) && !v.inf {
					c.data.Delete(key.(string))
				}
				return true
			})
		}
	}()
}

func (c *cache) Set(key string, value any, ttl time.Duration) {
	infiniteKey := false
	// Если TTL 0 или меньше, ключ бесконечный
	if ttl <= 0 {
		infiniteKey = true
	}
	// TTL считается как кол-во секунд с момента добавления ключа
	t := time.Now().Add(ttl)
	data := dataValue{value, t, infiniteKey}
	c.data.Store(key, data)
}

func (c *cache) Get(key string) (any, bool) {
	val, ok := c.data.Load(key)
	v, _ := val.(dataValue)
	returnVal := v.Value
	// Та же проверка на удаление
	if v.TTL.Before(time.Now()) && !v.inf {
		ok = false
		c.Delete(key)
	}
	return returnVal, ok
}

func (c *cache) Delete(key string) {
	c.data.Delete(key)
}

func (c *cache) Stop() {
	c.ticker.Stop()
}

func (c *cache) Flush() {
	c.data.Clear()
}

func (c *cache) Keys() []string {
	nonExpiredKeys := make([]string, 0)
	c.data.Range(func(key, value any) bool {
		v := value.(dataValue)
		// Если TTL не истек или ключ бесконечный, добавить ключ в список
		if v.TTL.After(time.Now()) || v.inf {
			nonExpiredKeys = append(nonExpiredKeys, key.(string))
		}
		return true
	})
	return nonExpiredKeys
}

func NewCache() *cache {
	cacheStruct := &cache{sync.Map{}, *time.NewTicker(FirstReadTimeout)}
	return cacheStruct
}

func main() {
	cache := NewCache()
	cache.Run()
	defer cache.Stop()

	wg := sync.WaitGroup{}

	wg.Go(func() {
		for index, value := range CacheValues {
			cache.Set(fmt.Sprintf("key_%d", index), value, time.Duration(value)*time.Millisecond)
		}
	})

	wg.Go(func() {
		time.Sleep(FirstReadTimeout)

		for index := range len(CacheValues) {
			key := fmt.Sprintf("key_%d", index)
			value, ok := cache.Get(key)
			if !ok {
				fmt.Printf("%s deleted\n", key)
				continue
			}
			fmt.Printf("%s:%s\n", key, value)
		}
	})

	wg.Wait()

	fmt.Println()

	// Тест для Keys
	nonExpiredKeys := cache.Keys()
	fmt.Printf("Ключи, которые не истекли спустя %v: %v\n", FirstReadTimeout, nonExpiredKeys)

	fmt.Println()
	time.Sleep(SecondReadTimeout)

	for index := range len(CacheValues) {
		key := fmt.Sprintf("key_%d", index)
		value, ok := cache.Get(key)
		if !ok {
			fmt.Printf("%s deleted\n", key)
			continue
		}
		fmt.Printf("%s:%s\n", key, value)
	}

	// Тест 1: Вечный ключ (TTL = 0)
	cache.Set("forever", "I live forever", 0)
	time.Sleep(3 * time.Second)
	if val, ok := cache.Get("forever"); ok {
		fmt.Printf("Вечный ключ: %v\n", val)
	}

	// Тест 2: Удаление
	cache.Set("todelete", "delete me", 10*time.Second)
	cache.Delete("todelete")
	if _, ok := cache.Get("todelete"); !ok {
		fmt.Println("Ключ 'todelete' успешно удален")
	}

	// Тест 3: Конкурентный доступ
	var wg2 sync.WaitGroup
	for i := range 100 {
		wg2.Go(func() {
			//defer wg2.Done() // Этой строчки тут быть не должно, потому что Done используется вместе со старым Add,
			// Go сам убирает таску из WaitGroup, когда функция внутри возвращает что-то
			// Если её не закомментить, то падает с паникой panic: sync: negative WaitGroup counter
			cache.Set(fmt.Sprintf("concurrent_%d", i), i, 5*time.Second)
			cache.Get(fmt.Sprintf("concurrent_%d", i))
		})
	}
	wg2.Wait()
	fmt.Println("Конкурентные операции завершены")

	// Вот как тест 3 выглядел бы с использованием Add и Done
	var wg3 sync.WaitGroup
	for i := range 100 {
		wg3.Add(1)
		go func() {
			defer wg3.Done()
			cache.Set(fmt.Sprintf("concurrent_%d", i), i, 5*time.Second)
			cache.Get(fmt.Sprintf("concurrent_%d", i))
		}()
	}
	wg3.Wait()
	fmt.Println("Конкурентные операции завершены ещё раз")

	// Тест Flush
	_, isDeleted := cache.Get("forever")
	fmt.Printf("Ключ forever удален: %v\n", !isDeleted)

	cache.Flush()
	_, isDeleted = cache.Get("forever")
	fmt.Printf("Ключ forever удален: %v\n", !isDeleted)

}
