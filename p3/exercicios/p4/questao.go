package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path"
	"time"
)

var chCountDirFscked = make(chan int)
var chCountFileFscked = make(chan int)
var chCountDirDamaged = make(chan int)
var chCountFileDamaged = make(chan int)

var chFilesToFsck = make(chan string, 30)
var chDirsToFsck = make(chan string, 30)

var chJoin = make(chan int)

func main() {
	root := os.Args[1]
	chJoinMain := make(chan int, 16)
	for i := 0; i < 8; i++ {
		go func() {
			for file := range chFilesToFsck {
				chCountFileFscked <- 1
				if fsckFile(file) {
					chCountFileDamaged <- 1
					chDirsToFsck <- path.Dir(file)
				}
			}
			chJoinMain <- 1
		}()

		go func() {
			for dir := range chDirsToFsck {
				chCountDirFscked <- 1
				if fsckDir(dir) {
					chCountDirDamaged <- 1
					chDirsToFsck <- path.Dir(dir)
				} else {
					_, ok := <-chDirsToFsck
					if !ok && len(chDirsToFsck) == 0 {
						break
					}
				}
			}
			chJoinMain <- 1
		}()
	}

	go printProgress(chJoin)
	go analyseDir(root)

	for i := 0; i < 16; i++ {
		<-chJoinMain
	}
	chJoin <- 1
	close(chDirsToFsck)
	close(chJoinMain)
	close(chCountDirDamaged)
	close(chCountDirFscked)
	close(chCountFileDamaged)
	close(chCountFileFscked)

}

func printProgress(chJoin chan int) {
	countDirDamaged := 0
	countDirFscked := 0
	countFileDamaged := 0
	countFIleFscked := 0

	go func() {
		for out := range chCountDirDamaged {
			countDirDamaged += out
		}
	}()

	go func() {
		for out := range chCountDirFscked {
			countDirFscked += out
		}
	}()

	go func() {
		for out := range chCountFileDamaged {
			countFileDamaged += out
		}
	}()

	go func() {
		for out := range chCountFileFscked {
			countFIleFscked += out
		}
	}()

	for {
		select {
		case <-chJoin:
			close(chJoin)
			return
		default:
			time.Sleep(1 * time.Second)
			fmt.Printf("fscked_files %d damaged_files %d fscked_dirs %d damaged_dirs %d", countFIleFscked, countFileDamaged, countDirFscked, countDirDamaged)
		}
	}

}

func analyseDir(root string) {
	files, err := ioutil.ReadDir(root)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() {
			go analyseDir(root + "/" + file.Name())
		} else {
			chFilesToFsck <- root + "/" + file.Name()
		}
	}
	close(chFilesToFsck)
}

func fsckFile(root string) bool {
	//dorme por um tempo aleatório entre 0 e 3 segs. aumente se //preferir
	rSleep := rand.Intn(4)
	time.Sleep(time.Duration(rSleep) * time.Second)
	//retorna true ou false com igual probabilidade
	rn := rand.Intn(2)
	return (rn%2 == 0)
}

func fsckDir(root string) bool {
	//dorme por um tempo aleatório entre 0 e 3 segs. aumente se //preferir
	rSleep := rand.Intn(4)
	time.Sleep(time.Duration(rSleep) * time.Second)
	//retorna true ou false com igual probabilidade
	rn := rand.Intn(2)
	return (rn%2 == 0)
}
