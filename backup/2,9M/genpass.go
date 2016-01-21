package main

import (
	"time"
)

const (
	MINLENPASS = 1
	MAXLENPASS = 12
)

type chanproc struct {
	currentPass []byte
	temps       chan time.Duration
}

func (proc *chanproc) GenPass(nbhashstart int) {
	proc.NextPass()
	for i := 0; i < nbhashstart; i++ {
		proc.NextPass()
	}
}

func (proc *chanproc) NextPass() {
	ln := len(proc.currentPass) / 2

	//fmt.Println("proc pass : ", proc.currentPass, " len ", ln, " addr ")
	if ln < MINLENPASS || ln == MAXLENPASS {
		proc.currentPass = make([]byte, MINLENPASS*2, MAXLENPASS*2)
		return
	}
	i := 0
	for {
		proc.currentPass[i*2] = proc.currentPass[i*2] + 1
		if proc.currentPass[i*2] != 0 {
			return
		}
		i++
		if i == ln {
			proc.currentPass = append(proc.currentPass, 1, 0)
			return
		}
	}
}
