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

type Arbitrator struct { // árbitro que tem um array de garfos com propriedades mutex
    forks []*Fork        // além de ser mutex também
    mutex sync.Mutex
}

func (a *Arbitrator) requestForks(left, right int) { // funç. para pedir os garfos para o árbitro
    a.mutex.Lock() // seção crítica
    for { // for looping baseado em condição (tipo um for true {...})
        if a.forks[left] != nil && a.forks[right] != nil {  // se os garfos que o árbitro tem
            break                                       // não tem valor inicializado, para o laço
        }
    }
    a.forks[left].Lock() // seção crítica
    a.forks[right].Lock() // seção crítica
    a.mutex.Unlock() // fim da seção crítica
}

func (a *Arbitrator) releaseForks(left, right int) { // funç. para entregar os garfos de volta
    a.mutex.Lock()
    a.forks[left].Unlock()
    a.forks[right].Unlock()
    a.mutex.Unlock()
}

func (p Philosopher) eat(wg *sync.WaitGroup) { // para o filósofo comer ele precisa 
    defer wg.Done()

    p.arbitrator.requestForks(p.id, (p.id + 1) % 5) // fazer a requisição dos garfos

    fmt.Printf("Philosopher %d is eating\n", p.id) 

    p.arbitrator.releaseForks(p.id, (p.id + 1) % 5) // entregar os garfos de volta

    fmt.Printf("Philosopher %d is done eating\n", p.id)
}

func main() {

	start := time.Now()
    forks := make([]*Fork, 5)
    for i := 0; i < 5; i++ { // inicializando os garfos enumeradamente
        forks[i] = &Fork{}
    }

    arbitrator := &Arbitrator{ // dando os garfos ao árbitro
        forks: forks,
    }

    philosophers := make([]*Philosopher, 5) // inicializando os filósofos
    for i := 0; i < 5; i++ {
        philosophers[i] = &Philosopher{
            id:        i + 1,
            leftFork:  forks[i],
            rightFork: forks[(i+1)%5],
            arbitrator: arbitrator,
        }
    }

    var wg sync.WaitGroup

    // esta parte é para testar um caso em que eles comam n vezes
	// for j := 0; j<60; j++{
	// 	for i := 0; i < 5; i++ {
	// 		wg.Add(1)
	// 		go philosophers[i].eat(&wg)
	// 	}
	// }

    for i := 0; i < 5; i++ {
        wg.Add(1)
        go philosophers[i].eat(&wg)
    }
    wg.Wait() // esperando todas as go routines

	log.Printf("Time took: %s", time.Since(start))
}
