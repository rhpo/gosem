package sema

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
	wq    chan struct{}
}

// -------------------------
// Sémaphore procédural
// -------------------------

func I(name string, initial int) {
	globalMu.Lock()
	defer globalMu.Unlock()
	semaphores[name] = &semaphore{
		count: initial,
		wq:    make(chan struct{}),
	}
}

// P decrements the semaphore count for the given name and waits if the count is negative.
func P(name string) {
	s := getSemaphore(name)
	s.mu.Lock()
	s.count--
	if s.count < 0 {
		s.mu.Unlock()
		<-s.wq
		return
	}
	s.mu.Unlock()
}

// V increments the semaphore count and signals if it was previously zero.
func V(name string) {
	s := getSemaphore(name)
	s.mu.Lock()
	s.count++
	if s.count <= 0 {
		s.wq <- struct{}{}
	}
	s.mu.Unlock()
}

// getSemaphore retrieves a semaphore by its name.
// It panics if the semaphore is not found.
func getSemaphore(name string) *semaphore {
	globalMu.Lock()
	s, ok := semaphores[name]
	globalMu.Unlock()
	if !ok {
		panic("semaphore not found: " + name)
	}
	return s
}

// -------------------------
// Gestion simple des goroutines
// -------------------------

// RandomDelay pauses execution for a random duration up to 50 milliseconds.
func RandomDelay() {
	time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)
}

// Repeat calls the provided function fn a specified number of times, with a random delay between calls.
func Repeat(times int, fn func()) {
	if times <= 0 {
		for {
			RandomDelay()
			fn()
		}
	} else {
		for range times {
			RandomDelay()
			fn()
		}
	}
}

// Loop calls the provided function fn once.
func Loop(fn func()) {
	Repeat(0, fn)
}

// Go lance une goroutine et l'ajoute au WaitGroup interne
func Process(f func()) {
	wg.Go(func() {
		// Add random delay (to simulate asyncronous behaviour)
		RandomDelay()

		f()
	})
}

// Wait blocks until all goroutines launched via Go() have completed.
func Wait() {
	wg.Wait()
}

var semaphoreNameIndex = 0

// NomSemaphore generates a unique semaphore name.
func NomSemaphore() string {
	semaphoreNameIndex++
	return fmt.Sprint("s", semaphoreNameIndex)
}
