package main

import (
	"errors"
)

var ErrEmptyStack = errors.New("ErrEmptyStack")

type Transaction struct {
	local KvMap
	next  *Transaction
}

func NewTransaction() *Transaction {
	return &Transaction{
		local: make(KvMap),
		next:  nil,
	}
}

// This is a linked list implementation of a stack.
// It can also be done with a []int
type Stack struct {
	top  *Transaction
	size int
}

func NewStack() Stack {
	return Stack{
		top:  nil,
		size: 0,
	}
}

func (kv *KvStore) Peek() *Transaction {
	return kv.stack.top
}

func (kv *KvStore) Pop() (*Transaction, error) {
	if kv.stack.size <= 0 {
		return nil, ErrEmptyStack
	}

	currTop := kv.Peek()

	kv.stack.top = currTop.next
	kv.stack.size--

	return currTop, nil
}

func (kv *KvStore) Push(t *Transaction) {
	previousTop := kv.Peek()

	kv.stack.size++
	kv.stack.top = t
	kv.stack.top.next = previousTop
}
