package main

import (
	"fmt"
	"golang.org/x/crypto/md4"
	"runtime"
	"time"
)

type chanproc struct {
	mdp       chan []byte
	nomorejob chan bool
	temps     chan time.Duration
}

const (
	NBHASH  = 500000000
	LENFILE = 100
)

var (
	NBPROC   int
	tempsavg time.Duration
)

func main() {
	NBPROC = runtime.NumCPU()
	proc := make([]chanproc, NBPROC)

	for i := 0; i < NBPROC; i++ {

		proc[i].mdp = make(chan []byte, LENFILE)
		proc[i].nomorejob = make(chan bool)
		proc[i].temps = make(chan time.Duration)
		go createWork(proc[i], byte(i), NBHASH/NBPROC)
	}

	fmt.Println("Début du compteur.")

	// Launch work
	for i := 0; i < NBPROC; i++ {
		go doWork(proc[i], i)
	}

	for i := 0; i < NBPROC; i++ {
		tempsavg = <-proc[i].temps
		fmt.Printf("\n%f hash/s pour coeur %d\n", float64(NBHASH)/tempsavg.Seconds(), i)
	}

}

func createWork(proc chanproc, i byte, nbtodo int) {
	i = i * 12
	for {
		i = (i + 7) % 200
		proc.mdp <- []byte{i, 0, i / 2, 0, 48, 0, i + 5, 0, i + 2, 0, i + 8, 0, i / 4, 0, i, 0}
		nbtodo--
		if nbtodo == 0 {
			proc.nomorejob <- true
			return
		}
	}
}

func doWork(proc chanproc, id int) {
	start := time.Now()
	for {
		select {
		case pass := <-proc.mdp:
			{
				tmp := md4.New()
				tmp.Write(pass)
				tmp.Sum(nil)
			}
		case <-proc.nomorejob:
			{
				if len(proc.mdp) != 0 {
					for i := 0; i < len(proc.mdp); i++ {
						tmp := md4.New()
						tmp.Write(<-proc.mdp)
						tmp.Sum(nil)
					}
				}
				proc.temps <- time.Since(start)
				return
			}
		}
	}
}
