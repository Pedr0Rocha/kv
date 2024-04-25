# Key Value store in memory

Simple key value store in memory using a linked list to implement
transactions and rollback.

The idea is to implement a simple version of a kv store to understand the
concepts and how it can become more complex over time.

## Usage

This example sets "test" to three different values using nested transactions.

Put signature is `key`, `value`, `ttl`.

```go
		kvStore.Put("test", 1, nil)

		kvStore.Begin() // first T1

		kvStore.Put("test", 10, nil)

		kvStore.Begin() // nested T2

		kvStore.Put("test", 20, nil)

		kvStore.Commit() // end nested T2

		kvStore.Commit() // end first T1

		kvStore.Get("test") // value=10
```

In this example, "test" will be assigned to 1 -> 20 -> 10. 10 being the final
value commited from the first transaction.

- [x] KV store
- [x] Transactions/Rollback
- [x] TTL to entries. Expires on get
- [ ] Auto delete expired entries with a go routine
