package main

import (
	"sync"
)

type Lock struct {
	users     map[string]int
	exclusive bool
}

func (l *Lock) Add(user string) int {
	u, ok := l.users[user]
	if ok {
		l.users[user] = u + 1
	} else {
		l.users[user] = 1
	}
	return len(l.users)
}

func (l *Lock) Del(user string) int {
	delete(l.users, user)
	return len(l.users)
}

type LockManager struct {
	locks map[string]*Lock
	mutex sync.Mutex
}

func NewLockManager() *LockManager {
	m := &LockManager{}
	m.locks = make(map[string]*Lock)
	return m
}

// TODO rimeout
func (m *LockManager) Require(user, key string, exclusive bool) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	l, exist := m.locks[key]

	if !exist {
		l = &Lock{}
		l.users = make(map[string]int)
		l.exclusive = exclusive
		l.Add(user)
		return true
	}

	if l.exclusive {
		return false
	}

	if exclusive {
		return false
	}

	l.Add(user)

	return true
}

func (m *LockManager) Relase(user, key string) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	l, exist := m.locks[key]
	if !exist {
		return true
	}

	n := l.Del(user)

	if n <= 0 {
		delete(m.locks, key)
	}
	return true
}

func (m *LockManager) RelaseUser(user string) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for key, l := range m.locks {
		n := l.Del(user)
		if n <= 0 {
			delete(m.locks, key)
		}
	}

	return true
}
