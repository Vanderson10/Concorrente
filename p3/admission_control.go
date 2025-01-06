package main

import (
	"fmt"
	"time"
)

const maxCapacity = 10

type Request struct {
	ID int
}

func createReq() Request {
	return Request{}
}

func execReq(req Request, workerID int) {
	// Simulando o processamento da requisição
	fmt.Printf("Worker %d processando a requisição %d\n", workerID, req.ID)
	time.Sleep(1 * time.Second) // Simulando o tempo de processamento
}

func main() {
	requestChannel := make(chan Request)
	workerChannel := make(chan int)

	// Crie um pool de workers
	for i := 1; i <= maxCapacity; i++ {
		go func(workerID int) {
			for {
				request := <-requestChannel
				execReq(request, workerID)
				workerChannel <- workerID
			}
		}(i)
	}

	// Loop de criação de requisições
	reqID := 1
	for {
		<-workerChannel // Aguarda um worker estar disponível
		req := createReq()
		req.ID = reqID
		reqID++
		requestChannel <- req
	}
}
