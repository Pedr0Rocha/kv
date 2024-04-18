package main

import (
	"errors"
	"log/slog"
	"time"
)

var ErrValueNotFound = errors.New("ErrValueNotFound")
var ErrEntryExpired = errors.New("ErrEntryExpired")
var ErrNoActiveTransaction = errors.New("ErrNoActiveTransaction")

type KvEntry struct {
	value     int
	expiresAt time.Time
}

func NewKvEntry(value int, ttl *time.Duration) KvEntry {
	var expiresAt time.Time
	if ttl != nil {
		expiresAt = time.Now().Add(*ttl)
	}

	return KvEntry{
		value:     value,
		expiresAt: expiresAt,
	}
}

// if expiresAt is zeroed, the entry will never expire
func (e KvEntry) IsExpired() bool {
	return !e.expiresAt.IsZero() && time.Now().After(e.expiresAt)
}

type KvMap map[string]KvEntry

type KvStore struct {
	store KvMap
	stack Stack
}

func NewKvStore() KvStore {
	return KvStore{
		store: make(KvMap),
		stack: NewStack(),
	}
}

func (kv KvStore) Put(key string, value int, ttl *time.Duration) {
	s := kv.Peek()

	// no transaction running, update global store
	if s == nil {
		kv.store[key] = NewKvEntry(value, ttl)
		return
	}
	s.local[key] = NewKvEntry(value, ttl)
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
	entry, ok := kv.store[key]
	if !ok {
		slog.Error("Get failed. Value not found", "error", ErrValueNotFound, "key", key)
		return -1
	}

	if entry.IsExpired() {
		// clean expired keys
		delete(kv.store, key)
		slog.Error("Get failed. Entry expired.", "error", ErrEntryExpired, "key", key)
		return -1
	}

	return entry.value
}

func (kv *KvStore) Begin() {
	kv.Push(NewTransaction())
}

func (kv *KvStore) Commit() error {
	currTrans, err := kv.Pop()

	if err != nil {
		if err == ErrEmptyStack {
			slog.Error("Commit failed. No active transaction.")
			return ErrNoActiveTransaction
		}
		slog.Error("Commit failed.")
		return err
	}

	for k, v := range currTrans.local {
		kv.store[k] = v
	}

	return nil
}

func (kv *KvStore) Rollback() error {
	if _, err := kv.Pop(); err != nil {
		if err == ErrEmptyStack {
			slog.Error("Rollback failed. No active transaction.")
			return ErrNoActiveTransaction
		}
		slog.Error("Rollback failed.", "error", err)
		return err
	}
	return nil
}
