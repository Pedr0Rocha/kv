package main

import (
	"errors"
)

var ErrEmptyStack = errors.New("ErrEmptyStack")

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

func (kv *KvStore) Pop() error {
	if kv.stack.size <= 0 {
		return ErrEmptyStack
	}

	currTop := kv.Peek()

	kv.stack.top = currTop.next
	kv.stack.size--

	return nil
}

func (kv *KvStore) Push(t *Transaction) {
	previousTop := kv.Peek()

	kv.stack.size++
	kv.stack.top = t
	kv.stack.top.next = previousTop
}
