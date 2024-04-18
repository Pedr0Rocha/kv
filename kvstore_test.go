package main

import (
	"testing"
	"time"
)

// clear store and kvStore.for every test
func resetAllWrapper(t *testing.T, name string, testFunc func(t *testing.T)) {
	t.Run(name, func(t *testing.T) {
		kvStore = NewKvStore()
		testFunc(t)
	})
}

func TestKvStore(t *testing.T) {
	resetAllWrapper(t, "should insert into global store", func(t *testing.T) {
		kvStore.Put("test", 1, nil)

		if kvStore.Get("test") != 1 {
			t.Errorf("not inserted. expected=%d, got=%d", 1, kvStore.Get("test"))
		}
	})

	resetAllWrapper(t, "should insert and update into global store", func(t *testing.T) {
		kvStore.Put("test", 1, nil)
		kvStore.Put("test2", 1, nil)
		kvStore.Put("test2", 5, nil)
		kvStore.Put("test", 2, nil)

		if kvStore.Get("test") != 2 {
			t.Errorf("not inserted/updated. expected=%d, got=%d", 2, kvStore.Get("test"))
		}
		if kvStore.Get("test2") != 5 {
			t.Errorf("not inserted/updated. expected=%d, got=%d", 5, kvStore.Get("test2"))
		}
	})

	resetAllWrapper(t, "should delete from global store", func(t *testing.T) {
		kvStore.Put("test", 2, nil)
		kvStore.Put("test2", 5, nil)

		kvStore.Delete("test2")

		if kvStore.Get("test") != 2 {
			t.Errorf("was deleted. expected=%d, got=%d", 2, kvStore.Get("test"))
		}
		if value := kvStore.Get("test2"); value != -1 {
			t.Errorf("not deleted. expected=%d, got=%d", -1, kvStore.Get("test2"))
		}
	})
}

func TestKvStoreWithTtl(t *testing.T) {
	resetAllWrapper(t, "should insert into global store with ttl", func(t *testing.T) {
		ttl := time.Duration(50 * time.Millisecond)
		kvStore.Put("test", 1, &ttl)

		time.Sleep(100 * time.Millisecond)

		if value := kvStore.Get("test"); value != -1 {
			t.Errorf("key did not expire. expected=%d, got=%d", -1, value)
		}
	})
}

func TestKvStoreTransaction(t *testing.T) {
	resetAllWrapper(t, "should create and commit transaction", func(t *testing.T) {
		kvStore.Put("test", 1, nil)

		kvStore.Begin()
		kvStore.Put("test", 10, nil)
		if kvStore.Get("test") != 1 {
			t.Errorf("modified global store. expected=%d, got=%d", 1, kvStore.Get("test"))
		}
		kvStore.Commit()

		if kvStore.Get("test") != 10 {
			t.Errorf("not modified global store. expected=%d, got=%d", 10, kvStore.Get("test"))
		}
	})

	resetAllWrapper(t, "should rollback transaction", func(t *testing.T) {
		kvStore.Put("test", 1, nil)

		kvStore.Begin()
		kvStore.Put("test", 10, nil)
		if kvStore.Get("test") != 1 {
			t.Errorf("modified global store. expected=%d, got=%d", 1, kvStore.Get("test"))
		}
		kvStore.Rollback()

		if kvStore.Get("test") != 1 {
			t.Errorf("modified global store after rollback. expected=%d, got=%d", 1, kvStore.Get("test"))
		}
	})

	resetAllWrapper(t, "should handle nested transactions", func(t *testing.T) {
		kvStore.Put("test", 1, nil)

		// first T
		kvStore.Begin()
		kvStore.Put("test", 10, nil)
		if kvStore.Get("test") != 1 {
			t.Errorf("modified global store. expected=%d, got=%d", 1, kvStore.Get("test"))
		}
		if kvStore.Peek().local["test"].value != 10 {
			t.Errorf("not modified kvStore. expected=%d, got=%d", 10, kvStore.Peek().local["test"].value)
		}
		// nested T
		kvStore.Begin()
		if kvStore.Get("test") != 1 {
			t.Errorf("modified global store. expected=%d, got=%d", 1, kvStore.Get("test"))
		}
		kvStore.Put("test", 20, nil)
		if kvStore.stack.size != 2 {
			t.Errorf("wrong kvStore.size. expected=%d, got=%d", 2, kvStore.stack.size)
		}
		// end nested T
		kvStore.Commit()
		if kvStore.Get("test") != 20 {
			t.Errorf("did not modified global store. expected=%d, got=%d", 20, kvStore.Get("test"))
		}

		// end first T
		kvStore.Commit()

		if kvStore.Get("test") != 10 {
			t.Errorf("did not modified global store. expected=%d, got=%d", 10, kvStore.Get("test"))
		}

		if kvStore.stack.size != 0 {
			t.Errorf("kvStore.should be empty after all commits. expected=%d, got=%d", 0, kvStore.stack.size)
		}
	})

	resetAllWrapper(t, "should error if there are not transactions when commiting or rolling back", func(t *testing.T) {
		err := kvStore.Commit()
		if err != ErrNoActiveTransaction {
			t.Errorf("commit should error on empty transactions. expected=%v, got=%v", ErrNoActiveTransaction, err)
		}

		err = kvStore.Rollback()
		if err != ErrNoActiveTransaction {
			t.Errorf("rollback should error on empty transactions. expected=%v, got=%v", ErrNoActiveTransaction, err)
		}

		if kvStore.stack.size != 0 {
			t.Errorf("kvStore.should be empty. expected=%d, got=%d", 0, kvStore.stack.size)
		}
	})
}
