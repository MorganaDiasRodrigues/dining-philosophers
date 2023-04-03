package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type Fork struct {
    sync.Mutex
}

type Philosopher struct {
    id              int
    leftFork, rightFork *Fork
    arbitrator      *Arbitrator
}

type Arbitrator struct {
    forks []*Fork
    mutex sync.Mutex
}

func (a *Arbitrator) requestForks(left, right int) {
    a.mutex.Lock()
    for {
        if a.forks[left] != nil && a.forks[right] != nil {
            break
        }
    }
    a.forks[left].Lock()
    a.forks[right].Lock()
    a.mutex.Unlock()
}

func (a *Arbitrator) releaseForks(left, right int) {
    a.mutex.Lock()
    a.forks[left].Unlock()
    a.forks[right].Unlock()
    a.mutex.Unlock()
}

func (p Philosopher) eat(wg *sync.WaitGroup) {
    defer wg.Done()

    p.arbitrator.requestForks(p.id, (p.id + 1) % 5)

    fmt.Printf("Philosopher %d is eating\n", p.id)

    p.arbitrator.releaseForks(p.id, (p.id + 1) % 5)

    fmt.Printf("Philosopher %d is done eating\n", p.id)
}

func main() {

	start := time.Now()
    forks := make([]*Fork, 5)
    for i := 0; i < 5; i++ {
        forks[i] = &Fork{}
    }

    arbitrator := &Arbitrator{
        forks: forks,
    }

    philosophers := make([]*Philosopher, 5)
    for i := 0; i < 5; i++ {
        philosophers[i] = &Philosopher{
            id:        i + 1,
            leftFork:  forks[i],
            rightFork: forks[(i+1)%5],
            arbitrator: arbitrator,
        }
    }

    var wg sync.WaitGroup
    for i := 0; i < 5; i++ {
        wg.Add(1)
        go philosophers[i].eat(&wg)
    }
    wg.Wait()

	log.Printf("Time took: %s", time.Since(start))
}
