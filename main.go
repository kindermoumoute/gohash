package main

import (
	"fmt"
	"github.com/kindermoumoute/gohash/md4"
	"github.com/kindermoumoute/gohash/passtype"
	"log"
	"os"
	"runtime"
	"time"
)

const (
	NBHASH  = 1000
	LENFILE = 100
)

var (
	NBPROC     int
	tempsavg   time.Duration
	hashToFind []byte
)

func main() {
	/*
		from := []byte{'a', 0, 'a', 0, 'a', 0, 'a', 0, 'a', 0, 'a', 0}
		to := []byte{'b', 0, 'a', 0, 'a', 0, 'a', 0, 'a', 0, 'a', 0}
		toFind := []byte{40, 172, 43, 162, 224, 67, 89, 237, 203, 42, 92, 181, 80, 57, 177, 198}
		charset := "abcdefghijklmnopqrstuvwxyz"
		CalculHash(from, to, toFind, charset)
	*/
	Charset = []byte("0123456789ABCDEF")
	main := []byte{10, 4, 4, 4}
	tmp := []byte{8, 5, 4, 4}
}

func CalculHash(from []byte, to []byte, toFind []byte, charset []byte) {
	hashToFind = toFind
	Charset = charset
	NBPROC = runtime.NumCPU()
	runtime.GOMAXPROCS(NBPROC)
	fmt.Println("Le calcul s'effectuera en parralelle sur ", NBPROC, " processeur(s).")

	proc := make([]chanproc, NBPROC)
	fmt.Println("Initialisation des mots de passe...")
	for i := 0; i < NBPROC; i++ {
		proc[i].id = i
		proc[i].temps = make(chan time.Duration)
		if i > 0 {
			proc[i].currentPass = proc[i-1].currentPass
			for j := 0; j < NBHASH/NBPROC; j++ {
				proc[i].NextPass()
			}
		}
		proc[i].NextPass()
		//fmt.Printf("\n\nMot de passe de départ processeur %d : %s", i, string(toUTF8(proc[i].currentPass)))
	}

	fmt.Println("Lancement du calcul des hash.")

	// Launch work
	for i := 0; i < NBPROC; i++ {
		go doWork(proc[i])
	}

	for i := 0; i < NBPROC; i++ {
		tempsavg = <-proc[i].temps + tempsavg

	}
	fmt.Printf("\n%f hash/s\n", float64(NBHASH)/(tempsavg.Seconds()/float64(NBPROC)))
}

func doWork(proc chanproc) {
	i := NBHASH / NBPROC
	start := time.Now()
	f, err := os.OpenFile("loghash"+string(proc.id+48)+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("error opening file: %v", err)
		return
	}
	defer f.Close()

	log.SetOutput(f)

	for {
		tmp := md4.New()
		tmp.Write(proc.currentPass)
		pass := tmp.Sum(nil)
		fmt.Printf("%s : %x", proc.currentPass, pass)
		fmt.Println(proc.currentPass, pass)

		i--
		if i == 0 {
			proc.temps <- time.Since(start)
			return
		}
		proc.NextPass()
	}
}
