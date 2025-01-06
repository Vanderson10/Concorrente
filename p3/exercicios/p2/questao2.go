package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type Item struct {
	nome string
}

type Bid struct {
	item      Item
	bidValue  float32
	bidFailed bool
}

func main() {
	n, _ := strconv.Atoi(os.Args[1])
	timeoutSecs := 5
	chBid := handle(n, timeoutSecs)

	for bid := range chBid {
		fmt.Println(bid)
	}
}

func handle(nServers int, timeoutSecs int) chan Bid {
	chItem := itemsStream()
	chBid := make(chan Bid)
	chJoin := make(chan int, nServers)

	for i := 0; i < nServers; i++ {
		go server(chItem, chBid, chJoin, timeoutSecs)
	}

	go func() {
		for i := 0; i < nServers; i++ {
			<-chJoin
		}
		close(chBid)
	}()

	return chBid
}

func server(chItem <-chan Item, chBid chan<- Bid, chJoin chan int, timeoutSecs int) {

	for item := range chItem {
		chNewBid := make(chan Bid)
		go func() {
			chNewBid <- bid(item)
		}()
		func() {
			for {
				tick := time.Tick(time.Duration(timeoutSecs) * time.Second)
				select {
				case bid := <-chNewBid:
					chBid <- bid
					return
				case <-tick:
					chBid <- Bid{item, -1, true}
				}
			}
		}()
	}
	chJoin <- 1
}

func bid(item Item) Bid {
	time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
	return Bid{item: item, bidValue: 100, bidFailed: false}
}

func itemsStream() chan Item {
	chItem := make(chan Item, 30)
	for i := 0; i < 30; i++ {
		name := "item - " + strconv.Itoa(i)
		chItem <- Item{nome: name}
	}
	close(chItem)
	return chItem
}
