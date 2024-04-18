package main

import (
	"errors"
	"log/slog"
)

var ErrNoActiveTransaction = errors.New("ErrNoActiveTransaction")

type Transaction struct {
	local KvStore
	next  *Transaction
}

func NewTransaction() *Transaction {
	return &Transaction{
		local: NewKvStore(),
		next:  nil,
	}
}

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

func (s *Stack) Peek() *Transaction {
	return s.top
}

func (s *Stack) Pop() error {
	if s.size <= 0 {
		return ErrNoActiveTransaction
	}

	currTop := s.Peek()

	s.top = currTop.next
	s.size--

	return nil
}

func (s *Stack) Begin() *Transaction {
	currTop := s.Peek()

	newTransaction := NewTransaction()
	s.size++
	s.top = newTransaction
	s.top.next = currTop

	return newTransaction
}

func (s *Stack) Commit() error {
	currTrans := s.Peek()

	if currTrans == nil {
		slog.Error("Commit failed. No active transaction.")
		return ErrNoActiveTransaction
	}

	for k, v := range currTrans.local {
		store[k] = v
	}

	s.Pop()

	return nil
}

func (s *Stack) Rollback() error {
	if err := s.Pop(); err != nil {
		if err == ErrNoActiveTransaction {
			slog.Error("Rollback failed. No active transaction.")
			return err
		}
		slog.Error("Rollback failed.", "error", err)
		return err
	}
	return nil
}
