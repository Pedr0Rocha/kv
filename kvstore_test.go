package main

import "testing"

// clear store and stack for every test
func resetAllWrapper(t *testing.T, name string, testFunc func(t *testing.T)) {
	t.Run(name, func(t *testing.T) {
		store = NewKvStore()
		stack = NewStack()
		testFunc(t)
	})
}

func TestKvStore(t *testing.T) {
	resetAllWrapper(t, "should insert into global store", func(t *testing.T) {
		Put("test", 1)

		if store["test"] != 1 {
			t.Errorf("not inserted. expected=%d, got=%d", 1, store["test"])
		}
	})

	resetAllWrapper(t, "should insert and update into global store", func(t *testing.T) {
		Put("test", 1)
		Put("test2", 1)
		Put("test2", 5)
		Put("test", 2)

		if store["test"] != 2 {
			t.Errorf("not inserted/updated. expected=%d, got=%d", 2, store["test"])
		}
		if store["test2"] != 5 {
			t.Errorf("not inserted/updated. expected=%d, got=%d", 5, store["test2"])
		}
	})

	resetAllWrapper(t, "should delete from global store", func(t *testing.T) {
		Put("test", 2)
		Put("test2", 5)

		Delete("test2")

		if store["test"] != 2 {
			t.Errorf("was deleted. expected=%d, got=%d", 2, store["test"])
		}
		if _, ok := store["test2"]; ok {
			t.Errorf("not deleted. expected=nil, got=%d", store["test2"])
		}
	})
}

func TestKvStoreTransaction(t *testing.T) {
	resetAllWrapper(t, "should create and commit transaction", func(t *testing.T) {
		Put("test", 1)

		stack.Begin()

		Put("test", 10)

		if store["test"] != 1 {
			t.Errorf("modified global store. expected=%d, got=%d", 1, store["test"])
		}

		stack.Commit()

		if store["test"] != 10 {
			t.Errorf("not modified global store. expected=%d, got=%d", 10, store["test"])
		}
	})

	resetAllWrapper(t, "should rollback transaction", func(t *testing.T) {
		Put("test", 1)

		stack.Begin()

		Put("test", 10)

		if store["test"] != 1 {
			t.Errorf("modified global store. expected=%d, got=%d", 1, store["test"])
		}

		stack.Rollback()

		if store["test"] != 1 {
			t.Errorf("modified global store after rollback. expected=%d, got=%d", 1, store["test"])
		}
	})

	resetAllWrapper(t, "should handle nested transactions", func(t *testing.T) {
		Put("test", 1)

		// first T
		stack.Begin()

		Put("test", 10)

		if store["test"] != 1 {
			t.Errorf("modified global store. expected=%d, got=%d", 1, store["test"])
		}

		if stack.Peek().local["test"] != 10 {
			t.Errorf("not modified stack. expected=%d, got=%d", 10, stack.Peek().local["test"])
		}

		// nested T
		stack.Begin()

		if store["test"] != 1 {
			t.Errorf("modified global store. expected=%d, got=%d", 1, store["test"])
		}

		Put("test", 20)

		if stack.size != 2 {
			t.Errorf("wrong stack size. expected=%d, got=%d", 2, stack.size)
		}

		// end nested T
		stack.Commit()

		if store["test"] != 20 {
			t.Errorf("did not modified global store. expected=%d, got=%d", 20, store["test"])
		}

		// end first T
		stack.Commit()

		if store["test"] != 10 {
			t.Errorf("did not modified global store. expected=%d, got=%d", 10, store["test"])
		}

		if stack.size != 0 {
			t.Errorf("stack should be empty after all commits. expected=%d, got=%d", 0, stack.size)
		}
	})

	resetAllWrapper(t, "should error if there are not transactions when commiting or rolling back", func(t *testing.T) {
		err := stack.Commit()
		if err != ErrNoActiveTransaction {
			t.Errorf("commit should error on empty transactions. expected=%v, got=%v", ErrNoActiveTransaction, err)
		}

		err = stack.Rollback()
		if err != ErrNoActiveTransaction {
			t.Errorf("rollback should error on empty transactions. expected=%v, got=%v", ErrNoActiveTransaction, err)
		}

		if stack.size != 0 {
			t.Errorf("stack should be empty. expected=%d, got=%d", 0, stack.size)
		}
	})
}
