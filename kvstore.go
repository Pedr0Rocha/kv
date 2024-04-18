package main

// global store
var store KvStore = NewKvStore()

// stack of transactions
var stack Stack = NewStack()

type KvStore map[string]int

func NewKvStore() KvStore {
	return make(KvStore)
}

func Put(key string, value int) {
	s := stack.Peek()

	// no transaction running, update global store
	if s == nil {
		store[key] = value
		return
	}
	s.local[key] = value
}

func Delete(key string) {
	s := stack.Peek()

	if s == nil {
		delete(store, key)
		return
	}
	delete(s.local, key)
}
