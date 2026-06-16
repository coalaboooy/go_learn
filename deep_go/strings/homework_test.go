package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"slices"
	"testing"
	"unsafe"
)

type COWBuffer struct {
	data []byte
	refs *int
}

func NewCOWBuffer(data []byte) COWBuffer {
	refCounter := 1
	return COWBuffer{
		data,
		&refCounter,
	}
}

func (b *COWBuffer) Clone() COWBuffer {
	var refCounter = b.refs
	newBuf := COWBuffer{
		slices.Clone(b.data), // вот это вот вообще ну неочевидно да окей b это ссылка но я то беру b.data это обращение к другой переменной че она по ссылке то тоже идет???
		refCounter,
	}
	*newBuf.refs += 1
	return newBuf
}

func (b *COWBuffer) Close() {
	*b.refs -= 1
}

func (b *COWBuffer) Update(index int, value byte) bool {
	if index < 0 || index >= len(b.data) {
		return false
	}

	if *b.refs == 1 {
		b.data[index] = value
		*b.refs += 1
		return true
	}

	newBuf := b.Clone()        // склонировал
	*newBuf.refs += 1          // счетчик поменял
	newBuf.data[index] = value // изменение внес

	b.data = newBuf.data // сохранил
	b.refs = newBuf.refs

	// все равно не проходит
	return true
}

func (b *COWBuffer) String() string {
	// Я искренне не понимаю, почему когда я возвращаю это,
	// то тесты в строках  не проходят. Если я сравниваю это с твоим методом,
	// то строки равны, но с моим методом они находятся по другому адресу???
	//return string(b.data)
	data := unsafe.SliceData(b.data)
	return unsafe.String(data, len(b.data))
}

func TestCOWBuffer(t *testing.T) {
	data := []byte{'a', 'b', 'c', 'd'}
	buffer := NewCOWBuffer(data)
	defer buffer.Close()

	copy1 := buffer.Clone()
	copy2 := buffer.Clone()

	assert.Equal(t, unsafe.SliceData(data), unsafe.SliceData(buffer.data))
	assert.Equal(t, unsafe.SliceData(buffer.data), unsafe.SliceData(copy1.data))
	assert.Equal(t, unsafe.SliceData(copy1.data), unsafe.SliceData(copy2.data))

	fmt.Println("### data")
	fmt.Printf("%#v\n", (*byte)(unsafe.SliceData(data)))
	fmt.Printf("%#v\n", unsafe.StringData(buffer.String()))

	fmt.Println("### copy1")
	fmt.Printf("%#v\n", (*byte)(unsafe.StringData(buffer.String())))
	fmt.Printf("%#v\n", unsafe.StringData(copy1.String()))

	fmt.Println("### copy2")
	fmt.Printf("%#v\n", (*byte)(unsafe.StringData(copy1.String())))
	fmt.Printf("%#v\n", unsafe.StringData(copy2.String()))

	assert.True(t, (*byte)(unsafe.SliceData(data)) == unsafe.StringData(buffer.String()))
	// я хз че тут честно, почему то в первом оно проходит, а в этих двух тут разные адреса
	// КОНЕЧНО ОНИ РАЗНЫЕ Я ЖЕ КОПИЮ СОЗДАВАЛ
	assert.True(t, (*byte)(unsafe.StringData(buffer.String())) == unsafe.StringData(copy1.String()))
	assert.True(t, (*byte)(unsafe.StringData(copy1.String())) == unsafe.StringData(copy2.String()))

	assert.True(t, buffer.Update(0, 'g'))
	assert.False(t, buffer.Update(-1, 'g'))
	assert.False(t, buffer.Update(4, 'g'))

	fmt.Printf("buffer: %#v\n", buffer.data)
	fmt.Printf("copy1: %#v\n", copy1.data)
	fmt.Printf("copy2: %#v\n", copy2.data)

	// а почему она меняется в копиях вообще если copy берет
	// копию данных из оригинального буфера? Потому что он по ссылке?
	// но оно же обращается к составляющей, т.е. копия? или нет?
	// как тогда клонировать вообще если update меняет по ссылке и все копии
	// так же указывают на эту ссылку?
	assert.True(t, reflect.DeepEqual([]byte{'g', 'b', 'c', 'd'}, buffer.data))
	assert.True(t, reflect.DeepEqual([]byte{'a', 'b', 'c', 'd'}, copy1.data))
	assert.True(t, reflect.DeepEqual([]byte{'a', 'b', 'c', 'd'}, copy2.data))

	assert.NotEqual(t, unsafe.SliceData(buffer.data), unsafe.SliceData(copy1.data))
	assert.Equal(t, unsafe.SliceData(copy1.data), unsafe.SliceData(copy2.data))

	copy1.Close()

	previous := copy2.data
	copy2.Update(0, 'f')
	current := copy2.data

	// 1 reference - don't need to copy buffer during update
	assert.Equal(t, unsafe.SliceData(previous), unsafe.SliceData(current)) //ладно

	copy2.Close()
}
