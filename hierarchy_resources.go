package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type Fork struct { // criando uma coleção de garfos
    sync.Mutex  // que vai possuir os métodos do mutex 
}

type Philosopher struct { // criando coleção de filósofos
    id              int // com respectivo id
    leftFork, rightFork *Fork // e que tem os garfos da esquerda e da direita (e são do tipo Fork)
}

func (p Philosopher) eat(wg *sync.WaitGroup, hierarchy *sync.Mutex) { // WaitGroup: esperar que várias goroutines terminem
                                                                      // hierarchy tbm é do tipo Mutex
    defer wg.Done() // atrasa a execução até que a próxima função retorne

    // seção crítica
    hierarchy.Lock()
    p.leftFork.Lock()
    p.rightFork.Lock()
    hierarchy.Unlock()
    
    fmt.Printf("Philosopher %d is eating\n", p.id)

    p.rightFork.Unlock()
    p.leftFork.Unlock()
    // fim da seção crítica

    fmt.Printf("Philosopher %d is done eating\n", p.id)
}

func main() {

	start := time.Now() // para medição de tempo

    forks := make([]*Fork, 5) //criando um array zerado do tipo Fork com 5 indices
    for i := 0; i < 5; i++ {// para cada indice de 1 até 5
        forks[i] = &Fork{} // array no indice i é do tipo Fork (&Type{} é uma espécie de ponteiro)
    }

    philosophers := make([]*Philosopher, 5) //criando um array zerado do tipo Philosophers com 5 indices
    for i := 0; i < 5; i++ { // para cada indice do array
        philosophers[i] = &Philosopher{ // cada indice no array será do tipo Philosopher
            id:        i + 1, // seu respectivo id
            leftFork:  forks[i], //garfos da esquerda
            rightFork: forks[(i+1)%5], // e da direita no array de Forks
        }
    }

    var hierarchy sync.Mutex // parâmetros esperados pela função
    var wg sync.WaitGroup  // eat dos Philosophers

    // esta parte é para testar um caso em que eles comam n vezes
	for j := 0; j<60; j++{
		for i := 0; i < 5; i++ {
			wg.Add(1)
			go philosophers[i].eat(&wg, &hierarchy)
		}
	}
    // for i := 0; i < 5; i++ { // para cada filósofo
    //     wg.Add(1) // contabndo cada go routine
    //     philosophers[i].eat(&wg, &hierarchy) // chamando de forma concorrente 
    //                                             // o array de filósfos com o método eat
    // }
		
    wg.Wait() // esperando todas as go routines

	log.Printf("Time took: %s", time.Since(start)) // tempo final
}