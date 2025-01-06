package main

import (
        "fmt"
        "math/rand"
        "time"
)

func produtorPar(ch chan int) int {
        for {
         rand.Seed(time.Now().UnixNano())
         par := rand.Int()
         if (par % 2 == 0) {
                ch <- par
         }
        }
}

func produtorImpar(ch chan int) int {
        for {
         rand.Seed(time.Now().UnixNano())
         impar := rand.Int()
         if (impar % 2 == 0) {
                ch <- impar
         }
        }
       

}

func consumidor(chPar chan int, chImpar chan int) int {
        for {
                impar := <- chImpar
                par := <- chPar
                fmt.Printf("Impar %d\n", impar)
                fmt.Printf("Par %d\n", par)
        }
       

}

func main() {
        chPar := make(chan int)
        chImpar := make(chan int)
        esperar := make(chan int)
        go produtorPar(chPar)
        go produtorImpar(chImpar)
        go consumidor(chPar, chImpar)
        <- esperar
}