package main

import (
	"fmt"
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
	chBid := handle(n)

	for bid := range chBid {
		fmt.Println(bid.item.nome)
	}
}

func handle(nServers int) chan Bid {
	chItem := itemsStream()
	chBid := make(chan Bid)
	chJoin := make(chan int, nServers)

	for i := 0; i < nServers; i++ {
		go server(chItem, chBid, chJoin)
	}

	go func() {
		for i := 0; i < nServers; i++ {
			<-chJoin
		}
		close(chBid)
	}()

	return chBid
}

func server(chItem <-chan Item, chBid chan<- Bid, chJoin chan int) {
	for item := range chItem {
		chBid <- bid(item)
	}
	chJoin <- 1
}

func bid(item Item) Bid {
	time.Sleep(3 * time.Second)
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
