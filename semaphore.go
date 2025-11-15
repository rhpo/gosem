package gosem

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var (
	semaphores = make(map[string]*semaphore)
	globalMu   sync.Mutex
	wg         sync.WaitGroup
)

type semaphore struct {
	count int
	mu    sync.Mutex
	cond  *sync.Cond
}

// I creates a new semaphore with the given name and initial count.
func I(name string, initial int) {
	globalMu.Lock()
	defer globalMu.Unlock()
	s := &semaphore{
		count: initial,
	}
	s.cond = sync.NewCond(&s.mu)
	semaphores[name] = s
}

// P decrements the semaphore count and waits if the count is less than zero.
func P(name string) {
	s := getSemaphore(name)
	s.mu.Lock()
	defer s.mu.Unlock()

	s.count--
	for s.count < 0 {
		s.cond.Wait()
	}
}

// V increments the semaphore count and signals any waiting goroutines.
func V(name string) {
	s := getSemaphore(name)
	s.mu.Lock()
	defer s.mu.Unlock()

	s.count++
	s.cond.Signal()
}

// getSemaphore retrieves a semaphore by its name.
// It panics if the semaphore is not found.
func getSemaphore(name string) *semaphore {
	globalMu.Lock()
	defer globalMu.Unlock()
	s, ok := semaphores[name]
	if !ok {
		panic("semaphore not found: " + name)
	}
	return s
}

// -------------------------
// Gestion simple des goroutines
// -------------------------
func RandomDelay() {
	time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
}

// Repeat calls the provided function fn a specified number of times or indefinitely if times is less than or equal to zero.
func Repeat(times int, fn func()) {
	if times <= 0 {
		for {
			fn()
			RandomDelay()
		}
	} else {
		for range times {
			fn()
		}
	}
}

func Loop(fn func()) {
	Repeat(0, fn)
}

// Process starts a goroutine and adds it to the internal WaitGroup.
func Process(f func()) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(time.Duration(rand.Intn(5)) * time.Millisecond) // Small startup delay
		f()
	}()
}

// Wait blocks until all goroutines launched via Process() have completed.
func Wait() {
	wg.Wait()
}

var (
	semaphoreNameIndex = 0
	semNameMu          sync.Mutex
)

// NomSemaphore generates a unique semaphore name.
func NomSemaphore() string {
	semNameMu.Lock()
	defer semNameMu.Unlock()
	semaphoreNameIndex++
	return fmt.Sprint("s", semaphoreNameIndex)
}
