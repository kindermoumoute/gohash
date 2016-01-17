package main

import (
	"fmt"
	"golang.org/x/crypto/md4"
	"time"
)

const (
	NBHASH       = 100000000
	NBSTRUCT     = 10000
	NBCREATEWORK = 1
	NBDOWORK     = 8
)

func main() {
	mdp := make(chan []byte, NBSTRUCT)
	done := make(chan bool, NBSTRUCT)
	fmt.Println("Début du compteur.")
	start := time.Now()

	// Create work
	for i := 0; i < NBCREATEWORK; i++ {
		go createWork(mdp)

	}

	// Do work
	for i := 0; i < NBDOWORK; i++ {
		go doWork(mdp, done)

	}

	for i := 0; i < NBHASH; i++ {
		<-done
	}

	elapsed := time.Since(start)
	fmt.Printf("\n%f hash/s\n", (NBHASH / elapsed.Seconds()))
}

func createWork(mdp chan []byte) {
	var i byte = 12
	for {
		i = (i + 7) % 200
		mdp <- []byte{i, 0, i / 2, 0, 48, 0, i + 5, 0, i + 2, 0, i + 8, 0, i / 4, 0, i, 0}
	}
}

func doWork(mdp chan []byte, done chan bool) {
	for {
		tmp := md4.New()
		tmp.Write(<-mdp)
		tmp.Sum(nil)
		done <- true
	}
}
