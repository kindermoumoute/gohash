package main

import (
	"bufio"
	"fmt"
	"golang.org/x/crypto/md4"
	"hash"
	"math/rand"
	"os"
	//"runtime"
	"time"
)

const (
	NBHASH       = 3000000
	NBSTRUCT     = 1000
	NBCREATEWORK = 1
	NBDOWORK     = 1
	DATAFILE     = "tmp.dat"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func main() {
	//runtime.GOMAXPROCS(runtime.NumCPU())
	myFile, _ := os.Create(DATAFILE)
	//defer myFile.Close()
	// Generate NBHASH passwords
	for i := 0; i < NBHASH; i++ {
		myFile.Write(genPass(rand.Intn(4) + 6))
	}
	myFile.Close()
	inFile, _ := os.Open(DATAFILE)
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)
	work := make(chan hash.Hash, NBSTRUCT)
	done := make(chan bool)
	fmt.Println("Début du compteur.")
	start := time.Now()

	// Create work
	for i := 0; i < NBCREATEWORK; i++ {
		go createWork(work, scanner, done)

	}

	// Do work
	for i := 0; i < NBDOWORK; i++ {
		go doWork(work)

	}

	// export
	//for i := 0; i < NBJOBDONE; i++ {
	//	go exportWork(jobsDone)

	//}
	<-done
	for {
		if len(work) == 0 {
			break
		}
	}
	elapsed := time.Since(start)
	fmt.Printf("\n%f hash/s\n", (NBHASH / elapsed.Seconds()))
}

func createWork(work chan hash.Hash, scanner *bufio.Scanner, done chan bool) {
	for scanner.Scan() {

		tmp2 := []byte(scanner.Text())
		//fmt.Println(tmp2, "plop")

		tmp := md4.New()
		tmp.Write(tmp2)
		work <- tmp
	}
	fmt.Println("CREATION FINI ! ")
	done <- true
} /*
func createWork(work chan hash.Hash, scanner *bufio.Scanner, done chan bool) {
	i := 0
	for scanner.Scan() {

		tmp := md4.New()
		tmp.Write(genPass(rand.Intn(4) + 6))
		work <- tmp
		i++
		fmt.Println("preworker :", i)
	}
	fmt.Println("CREATION FINI ! ")
	done <- true
}*/

func doWork(work chan hash.Hash) {
	for {
		tmp := <-work
		tmp.Sum(nil)
	}
}

/*
func exportWork(jobsDone chan int) {

	clock := time.NewTicker(time.Second)
	jobsDone <- 0
	for {
		<-clock.C
		i := <-jobsDone
		fmt.Printf("\n%d hash/s\n", i)
		jobsDone <- 0
	}

}*/

func genPass(n int) []byte {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return utfy(string(b) + "\n")
}

func utfy(s string) []byte {
	b := []byte(s)
	ret := make([]byte, len(b)*2)
	for i := range b {
		ret[i*2] = b[i]
	}
	return ret
}
