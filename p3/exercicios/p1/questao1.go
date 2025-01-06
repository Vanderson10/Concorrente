package main

import (
	"math/rand"
	"os"
	"strconv"
	"time"
)

func main() {
	n, _ := strconv.Atoi(os.Args[1])

	ch := make(chan int, n)

	for i := 0; i < n; i++ {
		go sleeper(ch)
	}

	for i := 0; i < n; i++ {
		<-ch
	}
}

func sleeper(ch chan<- int) {
	randomNumber := rand.Intn(5)
	time.Sleep(time.Duration(randomNumber) * time.Second)
	ch <- 1
}
