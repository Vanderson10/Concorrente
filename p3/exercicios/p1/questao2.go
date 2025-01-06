package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func main() {
	n, _ := strconv.Atoi(os.Args[1])

	chJoin := make(chan int, n)
	channels := make([]chan int, n)

	for i := range channels {
		channels[i] = make(chan int, 1)
	}

	for i := 0; i < n; i++ {
		fmt.Printf("A thread %d colocou em: %d \n", i, (i+1)%n)
		go sleeper(chJoin, channels[i], channels[(i+1)%n])
	}

	for i := 0; i < n; i++ {
		<-chJoin
	}
}

func sleeper(chJoin chan<- int, chIn <-chan int, chOut chan<- int) {
	randomNumber := rand.Intn(5)
	time.Sleep(time.Duration(randomNumber) * time.Second)

	newRandomnumber := rand.Intn(10)
	chOut <- newRandomnumber
	fmt.Printf("A thread colocou : %d \n", newRandomnumber)

	//IMPLEMENTAÇAÕ DA BARREIRA (QUE NÃO PRECISA COLOCAR)

	timeToSleep := <-chIn
	fmt.Printf("A thread retirou : %d \n", timeToSleep)
	time.Sleep(time.Duration(timeToSleep) * time.Second)

	chJoin <- 1
}
