package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

var bucketLimit = 20
var tokenRate = 0.5

type Request struct {
	name string
	size int
}

func main() {
	chBucket := make(chan int, bucketLimit)
	go incrementBucket(chBucket)

	chRequests := make(chan Request, 10)
	go generateRequests(chRequests)

	chMutex := make(chan int, 1)
	handleRequests(chMutex, chRequests, chBucket)

}

func generateRequests(chRequests chan Request) {
	for i := 0; i < 10; i++ {
		size := rand.Intn(10)
		name := "Requisicao numero: " + strconv.Itoa(i)
		chRequests <- Request{name: name, size: size}
	}
	close(chRequests)
}

func incrementBucket(chBucket chan int) {
	for {
		time.Sleep(time.Duration(1/tokenRate) * time.Second)
		chBucket <- 1
	}
}

func handleRequests(chMutex chan int, chRequests chan Request, chBucket chan int) {
	chJoin := make(chan int)
	for Request := range chRequests {
		go run(chMutex, Request, chBucket)
	}
	<-chJoin
}

func run(chMutex chan int, req Request, chBucket chan int) {
	limitCap_wait(chMutex, req, chBucket)
	handleReq(req)
}

func limitCap_wait(chMutex chan int, req Request, chBucket chan int) {
	fmt.Println("Requisição no limit cap wait: " + req.name)
	chMutex <- 1
	for i := 0; i < req.size; i++ {
		<-chBucket
	}
	<-chMutex
	fmt.Println(len(chMutex))
}

func handleReq(req Request) {
	fmt.Println("Requisição atendida: " + req.name)
}
