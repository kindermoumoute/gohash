package main

import (
	"fmt"
	"golang.org/x/crypto/md4"
	//"log"
	"runtime"
	"time"
)

type chanproc struct {
	id   int
	pass chan []byte
}

const (
	NBHASH  = 10000000
	LENFILE = 100

	MINLENPASS = 1
	MAXLENPASS = 12

	CHARSTART = 33
	CHAREND   = 126
)

var (
	NBPROC int
)

func main() {
	NBPROC = runtime.NumCPU()
	runtime.GOMAXPROCS(NBPROC)
	proc := make([]chanproc, NBPROC)

	//work := make(chan []byte, LENFILE)
	fmt.Println("Le calcul s'effectuera en parralère sur ", NBPROC, " processeur(s).")
	for i := 0; i < NBPROC; i++ {
		go createWork(proc[i].pass, nil)
		proc[i].id = i
		proc[i].pass = make(chan []byte, LENFILE)
	}

	fmt.Println("Lancement du calcul des hash.")

	// Launch work
	start := time.Now()
	for i := 0; i < NBPROC; i++ {
		go doWork(proc[i])
	}

	for i := 0; i < NBPROC; i++ {
		for len(proc[i].pass) > 0 {

		}
	}

	fmt.Printf("\n%f hash/s\n", float64(NBHASH)/time.Since(start).Seconds())

}

func doWork(proc chanproc) {
	//i := NBHASH / NBPROC
	for {
		pass := <-proc.pass
		tmp := md4.New()
		tmp.Write(pass)
		//fmt.Println(proc.id, " pass : ", pass, " len pass ", len(proc.pass))
		tmp.Sum(nil)
		/*
				f, err := os.OpenFile("testlogfile", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
				if err != nil {
				    t.Fatalf("error opening file: %v", err)
				}
				defer f.Close()

				log.SetOutput(f)
				log.Println("This is a test log entry")


			i--
			if i == 0 {
				proc.temps <- time.Since(start)
				return
			}*/
	}
}

// Generate password from the password "start"
func createWork(pass chan []byte, start []byte) {
	var b []byte
	if start == nil {
		b = make([]byte, MINLENPASS*2, MAXLENPASS*2)
		for i := range b {
			if i%2 == 0 {
				b[i] = CHARSTART
			}

		}
	} else {
		b = start
	}

	hashMax := NBHASH
	ln := len(b) / 2

	// lancement de la génération des passwords
	for {
		pass <- b
		hashMax--
		if hashMax == 0 || ln == MAXLENPASS {
			return
		}
		i := 0
		for {
			b[i*2] = b[i*2] + 1
			if b[i*2] != CHAREND+1 {
				break
			}
			b[i*2] = CHARSTART
			i++
			if i == ln {
				b = append(b, CHARSTART, 0)
				ln++
				break
			}
		}
	}
}
