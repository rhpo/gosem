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

func V(name string) {
	s := getSemaphore(name)
	s.mu.Lock()
	s.count++
	if s.count <= 0 {
		s.wq <- struct{}{}
	}
	s.mu.Unlock()
}

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

func RandomDelay() {
	time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)
}

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

// Wait attend que toutes les goroutines lancées via Go() soient terminées
func Wait() {
	wg.Wait()
}

var semaphoreNameIndex = 0

func NomSemaphore() string {
	semaphoreNameIndex++
	return fmt.Sprint("s", semaphoreNameIndex)
}
