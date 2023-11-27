package main

import (
	"fmt"
	"sync"
)

/*
Описание структуры
*/
type CopyOnWriteArrayList struct {
	mu     sync.RWMutex   // Мьютекс
	values []*interface{} // Срез указателей на интерфейсы, хранящие элементы списка
}

/*
*
Создания нового экземпляра CopyOnWriteArrayList
*/
func NewCopyOnWriteArrayList() *CopyOnWriteArrayList {
	return &CopyOnWriteArrayList{
		values: make([]*interface{}, 0), // Инициализация списка с пустым срезом указателей
	}
}

/*
*
Добавления элемента в список
*/
func (cow *CopyOnWriteArrayList) Add(value interface{}) {
	cow.mu.Lock()
	defer cow.mu.Unlock()

	copiedValues := make([]*interface{}, len(cow.values))
	copy(copiedValues, cow.values)

	// Добавление указателя на новый элемент в конец списка - использование указателей помогает нам избежать полного копирования
	copiedValues = append(copiedValues, &value)
	cow.values = copiedValues // Обновление списка
}

/*
*
Получения итератора по списку
*/
func (cow *CopyOnWriteArrayList) Iterator() <-chan interface{} {
	cow.mu.RLock()
	defer cow.mu.RUnlock()

	ch := make(chan interface{})

	go func() {
		cow.mu.RLock()
		defer cow.mu.RUnlock()

		for _, value := range cow.values {
			ch <- *value
		}
		close(ch)
	}()

	return ch
}

/*
*
Пример работы
*/
func main() {
	list := NewCopyOnWriteArrayList()
	list.Add(1)
	list.Add(2)
	list.Add(3)

	iterator := list.Iterator()
	for value := range iterator {
		fmt.Println(value)
	}

	list.Add(4)

	fmt.Println("После добавления значения:")

	iterator = list.Iterator()
	for value := range iterator {
		fmt.Println(value)
	}
}
