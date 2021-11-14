package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	sem := make(chan int, 5)
	for x := 0; x <= 30; x++ {
		sem <- 1 // will block if there is MAX ints in sem
		fmt.Println(len(sem))
		go worker(x, sem)
	}
}

func worker(i int, out <-chan int) {

	min, max := 1, 10
	sleep := time.Duration(rand.Intn(max-min) + min)

	fmt.Printf("doing work on %d, sleeping for %d\n", i, sleep)
	time.Sleep(sleep * time.Second)
	<-out
}
