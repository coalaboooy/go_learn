package main

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

// go test -v homework_test.go

type OrderedMap struct {
	Key          int
	Value        int
	ContainsItem bool
	Left, Right  *OrderedMap
}

func NewOrderedMap() OrderedMap {
	return OrderedMap{}
}

func (m *OrderedMap) Insert(key, value int) {
	if !m.ContainsItem {
		// Если дерево пустое, вставляем узел
		m.Key = key
		m.Value = value
		m.ContainsItem = true
		return
	}

	if key == m.Key {
		// Если такой ключ нашелся, меняем значение
		m.Value = value
		return
	}

	// Ищем по дереву
	if key > m.Key {
		if m.Right != nil {
			// Если есть следующий узел, повторяем там
			m.Right.Insert(key, value)
		} else {
			// Если его нету, создаем
			newNode := NewOrderedMap()
			m.Right = &newNode
			m.Right.Key = key
			m.Right.Value = value
			m.Right.ContainsItem = true
		}
	} else {
		if m.Left != nil {
			m.Left.Insert(key, value)
		} else {
			// При создании узла почему то кстати работает только в таком порядке,
			// если создать, новому все назначить а потом его присвоить Left (или Right),
			// то не работает
			newNode := NewOrderedMap()
			m.Left = &newNode
			m.Left.Key = key
			m.Left.Value = value
			m.Left.ContainsItem = true
		}
	}
}

func (m *OrderedMap) Erase(key int) {
	if !m.Contains(key) {
		// Если такого ключа нету, ничего не делаем
		return
	}

	// Ищем по дереву,
	if key > m.Key {
		m.Right.Erase(key)
	} else if key < m.Key {
		m.Left.Erase(key)
	} else {
		//  и если ключ в текущем узле совпал, удаляем
		// Тут с указателями я просто перебором всех возможных комбинаций * и & подобрал рабочий вариант уже
		if m.Right == nil {
			// Если правого узла нет, то можно просто заменить текущий на левый. так на Хабре сказали
			if m.Left != nil {
				*m = *m.Left
			} else {
				// Но почему-то когда левый nil, то все падает, поэтому создаю новый и он почему-то считается пустым?
				// Ну, главное, что работает
				*m = NewOrderedMap()
			}
		} else {
			// Если правый есть, то надо найти наименьший узел в правом поддереве, убрать его и заменить текущий на него - тоже на Хабре написали
			// Хорошо вообще что в домашке ссылка была на статью эту, потому что в лекциях ничего про деревья эти не рассказали
			leftmostNodePointer := &m.Right
			for (*leftmostNodePointer).Left != nil {
				leftmostNodePointer = &(*leftmostNodePointer).Left
			}
			m.Value = (*leftmostNodePointer).Value
			m.Key = (*leftmostNodePointer).Key
			*leftmostNodePointer = nil
		}
	}
}

func (m *OrderedMap) Contains(key int) bool {
	// Просто поиск ключа по дереву и проброс флага наверх
	// Наверное, можно было бы и через ForEach сделать
	contains := false
	if key == m.Key {
		contains = true
		return contains
	}

	if key > m.Key {
		if m.Right != nil {
			return m.Right.Contains(key)
		} else {
			return contains
		}
	} else {
		if m.Left != nil {
			return m.Left.Contains(key)
		} else {
			return contains
		}
	}
}

func (m *OrderedMap) ContainsFE(key int) bool {
	// Ну я ж говорил. Интересно, какая из них быстрее работает но мне чет не хочется с бенчмарками разбираться
	contains := false
	m.ForEach(func(mkey, _ int) {
		if mkey == key {
			contains = true
		}
	})
	return contains
}

func (m *OrderedMap) Size() int {
	// Тоже проход по каждому элементу и увеличение счетчика
	// Вот это прям 100% можно через ForEach сделать
	size := 0
	if m.ContainsItem {
		size += 1
	}
	if m.Left != nil {
		size += m.Left.Size()
	}
	if m.Right != nil {
		size += m.Right.Size()
	}
	return size
}

func (m *OrderedMap) SizeFE() int {
	// Ну да, это совсем просто было. Я ForEach писал после Size и Contains, так что ради интереса это сделал просто
	size := 0
	m.ForEach(func(_, _ int) {
		size += 1
	})
	return size
}

func (m *OrderedMap) ForEach(action func(int, int)) {
	// Вот это я сначала написал в порядке центр-лево-право и сидел думал, как мне отфильтровать дерево,
	// а потом меня озарило. Очень горд собой

	// А работает оно потому что всегда сначала идем в самый левый (меньший) элемент,
	// и потом выполняем функцию, а потом в следующий правый (больший), и так по всем проходим
	if m.Left != nil {
		m.Left.ForEach(action)
	}
	if m.ContainsItem {
		action(m.Key, m.Value)
	}
	if m.Right != nil {
		m.Right.ForEach(action)
	}
}

func (m *OrderedMap) GetValue(key int) (int, bool) {
	// Функция получения значения, написал, чтобы протестировать замену значения в Insert
	// Если ok не true, то значит, что такого ключа нет
	resultValue := -1
	ok := false
	m.ForEach(func(mkey, value int) {
		if mkey == key {
			resultValue = value
			ok = true
		}
	})
	return resultValue, ok
}

func TestOrderedMap(t *testing.T) {
	data := NewOrderedMap()
	assert.Zero(t, data.Size())
	assert.Zero(t, data.SizeFE())

	// Проверка на то, что ForEach ничего не делает с пустым OrderedMap
	var emptyKeys []int
	data.ForEach(func(key, _ int) {
		emptyKeys = append(emptyKeys, key)
	})
	assert.Zero(t, len(emptyKeys))

	data.Insert(10, 10)
	data.Insert(5, 5)
	data.Insert(15, 15)
	data.Insert(2, 2)
	data.Insert(4, 4)
	data.Insert(12, 12)
	data.Insert(14, 14)

	assert.Equal(t, 7, data.Size())
	assert.Equal(t, 7, data.SizeFE())
	assert.True(t, data.Contains(4))
	assert.True(t, data.ContainsFE(4))
	assert.True(t, data.Contains(12))
	assert.True(t, data.ContainsFE(12))
	assert.False(t, data.Contains(3))
	assert.False(t, data.ContainsFE(3))
	assert.False(t, data.Contains(13))
	assert.False(t, data.ContainsFE(13))

	var keys []int
	expectedKeys := []int{2, 4, 5, 10, 12, 14, 15}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))

	data.Erase(15)
	data.Erase(14)
	data.Erase(2)

	// Проверка на удаление несуществующего ключа (как минимум на то, что не падает)
	data.Erase(99)

	assert.Equal(t, 4, data.Size())
	assert.Equal(t, 4, data.SizeFE())
	assert.True(t, data.Contains(4))
	assert.True(t, data.ContainsFE(4))
	assert.True(t, data.Contains(12))
	assert.True(t, data.ContainsFE(12))
	assert.False(t, data.Contains(2))
	assert.False(t, data.ContainsFE(2))
	assert.False(t, data.Contains(14))
	assert.False(t, data.ContainsFE(14))

	keys = nil
	expectedKeys = []int{4, 5, 10, 12}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))

	// Проверка GetValue
	NEValue, ok := data.GetValue(14)
	assert.False(t, ok)
	assert.Equal(t, -1, NEValue)

	value, ok := data.GetValue(4)
	assert.True(t, ok)
	assert.Equal(t, 4, value)

	// Проверка, что Insert меняет значения
	newData := NewOrderedMap()
	newData.Insert(1, 10)

	assert.True(t, newData.Contains(1))
	assert.True(t, newData.ContainsFE(1))

	newValue, ok := newData.GetValue(1)
	assert.True(t, ok)
	assert.Equal(t, 10, newValue)

	newData.Insert(1, 45)

	newValue, ok = newData.GetValue(1)
	assert.True(t, ok)
	assert.Equal(t, 45, newValue)
}
