package main

import (
	"fmt"
	"sync"
)

type Lock struct {
	users     map[string]int // map[user]count
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
	u, ok := l.users[user]
	if !ok {
		return len(l.users)
	}
	if u <= 1 {
		delete(l.users, user)
	} else {
		l.users[user] = u - 1
	}
	return len(l.users)
}

func (l *Lock) String() string {
	s := fmt.Sprintf("%v:%v[", len(l.users), l.exclusive)
	for user, u := range l.users {
		s += fmt.Sprintf("%v:%v,", user, u)
	}
	s += "]"
	return s
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

func (m *LockManager) Acquire(user, key string, exclusive bool) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	// TODO timeout

	l, exist := m.locks[key]

	// key is not exist, allow it
	if !exist {
		l = &Lock{}
		l.users = make(map[string]int)
		l.exclusive = exclusive
		l.Add(user)

		m.locks[key] = l
		return true
	}

	// check if is the only user
	only := true
	if len(l.users) > 1 {
		only = false
	}
	_, ok := l.users[user]
	if !ok {
		only = false
	}

	// lock is exclusive
	if l.exclusive {
		if !only {
			return false
		}
	}

	// request is exclusive
	if exclusive {
		if !only {
			return false
		}
	}

	// allow multiple users to lock
	l.Add(user)

	l.exclusive = l.exclusive || exclusive

	return true
}

func (m *LockManager) Relase(user, key string) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	l, exist := m.locks[key]
	if !exist {
		return false
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

func (m *LockManager) String() string {
	s := fmt.Sprintf("%v{\n", len(m.locks))
	for key, l := range m.locks {
		s += fmt.Sprintf("{%v:%v}\n", key, l.String())
	}
	s += "}"
	return s
}
