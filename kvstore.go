package main

import (
	"errors"
	"log/slog"
)

var ErrValueNotFound = errors.New("ErrValueNotFound")
var ErrNoActiveTransaction = errors.New("ErrNoActiveTransaction")

type KvStore struct {
	store map[string]int
	stack Stack
}

func NewKvStore() KvStore {
	return KvStore{
		store: make(map[string]int),
		stack: NewStack(),
	}
}

type Transaction struct {
	local map[string]int
	next  *Transaction
}

func NewTransaction() *Transaction {
	return &Transaction{
		local: make(map[string]int),
		next:  nil,
	}
}

func (kv KvStore) Put(key string, value int) {
	s := kv.Peek()

	// no transaction running, update global store
	if s == nil {
		kv.store[key] = value
		return
	}
	s.local[key] = value
}

func (kv KvStore) Delete(key string) {
	s := kv.Peek()

	if s == nil {
		delete(kv.store, key)
		return
	}
	delete(s.local, key)
}

func (kv KvStore) Get(key string) int {
	if value, ok := kv.store[key]; ok {
		return value
	}
	slog.Error("error", ErrValueNotFound, "key", key)
	return -1
}

func (kv *KvStore) Begin() {
	kv.Push(NewTransaction())
}

func (kv *KvStore) Commit() error {
	currTrans := kv.Peek()

	if currTrans == nil {
		slog.Error("Commit failed. No active transaction.")
		return ErrNoActiveTransaction
	}

	for k, v := range currTrans.local {
		kv.store[k] = v
	}

	kv.Pop()

	return nil
}

func (kv *KvStore) Rollback() error {
	if err := kv.Pop(); err != nil {
		if err == ErrEmptyStack {
			slog.Error("Rollback failed. No active transaction.")
			return ErrNoActiveTransaction
		}
		slog.Error("Rollback failed.", "error", err)
		return err
	}
	return nil
}
