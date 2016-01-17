package main

import (
	"fmt"
	"golang.org/x/crypto/md4"
	"runtime"
	"time"
)

const (
	NBHASH  = 1000000000
	LENFILE = 100
)

var (
	NBPROC   int
	tempsavg time.Duration
)

func main() {
	NBPROC = runtime.NumCPU()
	proc := make([]chanproc, NBPROC)

	fmt.Println("Initialisation des mots de passe...")
	for i := 0; i < NBPROC; i++ {
		proc[i].temps = make(chan time.Duration)
		proc[i].GenPass(i * NBHASH / NBPROC)
	}

	fmt.Println("Lancement du calcul des hash.")

	// Launch work
	for i := 0; i < NBPROC; i++ {
		go doWork(proc[i])
	}

	for i := 0; i < NBPROC; i++ {
		tempsavg = <-proc[i].temps
		fmt.Printf("\n%f hash/s pour coeur %d\n", float64(NBHASH)/tempsavg.Seconds(), i)
	}

}

func doWork(proc chanproc) {
	i := NBHASH / NBPROC
	start := time.Now()
	for {
		tmp := md4.New()
		tmp.Write(proc.currentPass)
		tmp.Sum(nil)

		i--
		if i == 0 {
			proc.temps <- time.Since(start)
			return
		}
		proc.NextPass()
	}
}
