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

// -------------------------
// Sémaphore procédural
// -------------------------
func I(name string, initial int) {
	globalMu.Lock()
	defer globalMu.Unlock()
	s := &semaphore{
		count: initial,
	}
	s.cond = sync.NewCond(&s.mu)
	semaphores[name] = s
}

func P(name string) {
	s := getSemaphore(name)
	s.mu.Lock()
	defer s.mu.Unlock()

	s.count--
	for s.count < 0 {
		s.cond.Wait()
	}
}

func V(name string) {
	s := getSemaphore(name)
	s.mu.Lock()
	defer s.mu.Unlock()

	s.count++
	s.cond.Signal()
}

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

// Process lance une goroutine et l'ajoute au WaitGroup interne
func Process(f func()) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(time.Duration(rand.Intn(5)) * time.Millisecond) // Small startup delay
		f()
	}()
}

// Wait attend que toutes les goroutines lancées via Process() soient terminées
func Wait() {
	wg.Wait()
}

var (
	semaphoreNameIndex = 0
	semNameMu          sync.Mutex
)

func NomSemaphore() string {
	semNameMu.Lock()
	defer semNameMu.Unlock()
	semaphoreNameIndex++
	return fmt.Sprint("s", semaphoreNameIndex)
}
