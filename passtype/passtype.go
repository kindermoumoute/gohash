package passtype

import (
	"time"
)

const (
	MINLENPASS = 8
	MAXLENPASS = 12
)

var (
	Charset []byte
)

type chanproc struct {
	currentPass []byte // password in utf8 with charset's indexes
	truePass    []byte // password in utf16
	id          int
	temps       chan time.Duration
}

/*
func initPass(from []byte, to []byte, nbmaxProc int) []byte {
	b := make()
}
*/
func passToTruePass(from []byte) []byte {
	ln := len(from) * 2
	b := make([]byte, ln)
	for i := 0; i < ln/2; i++ {
		b[i*2] = Charset[from[i]]
		b[i*2+1] = 0
	}
	return b
}

func truePassToPass(from []byte) []byte {
	ln := len(from) / 2
	b := make([]byte, ln)
	for i := 0; i < ln; i++ {
		b[i] = pos(from[i*2], Charset)
	}
	return b
}

func pos(value byte, b []byte) byte {
	for p, v := range b {
		if v == value {
			return byte(p)
		}
	}
	return byte(0)
}

//mainPass is higher
func addPass(mainPass, add []byte) []byte {
	lnp := len(mainPass)
	lna := len(add)
	charmax := len(Charset)

	retenu := 0
	for i := lna - 1; i >= 0; i-- {

		intadd := int(mainPass[i]) + int(add[i]) + retenu
		retenu = intadd / charmax
		mainPass[i] = byte(intadd % charmax)
		if i == 0 && retenu > 0 {
			mainPass = append(mainPass, 0)
			j := lnp
			for {
				if j == 0 {
					mainPass[0] = byte(retenu)
					break
				}
				mainPass[j] = mainPass[j-1]
				j--
			}
		}
	}
	return mainPass
}

// 'to' is higher than 'from'
/*
func substractPass(mainPass, sub []byte) []byte {
	lnp := len(mainPass)
	lna := len(add)
	charmax := len(Charset)

	retenu := 0
	for i := lna - 1; i >= 0; i-- {

		intadd := int(mainPass[i]) + int(add[i]) + retenu
		retenu = intadd / charmax
		mainPass[i] = byte(intadd % charmax)
		if i == 0 && retenu > 0 {
			mainPass = append(mainPass, 0)
			j := lnp
			for {
				if j == 0 {
					mainPass[0] = byte(retenu)
					break
				}
				mainPass[j] = mainPass[j-1]
				j--
			}
		}
	}
	return mainPass
}

/*
func (proc *chanproc) NextPass() {

	ln := len(proc.currentPass) / 2

	//fmt.Println("proc pass : ", proc.currentPass, " len ", ln, " addr ")
	if ln < MINLENPASS || ln == MAXLENPASS {
		proc.currentPass = make([]byte, MINLENPASS*2, MAXLENPASS*2)
		for i := range proc.currentPass {
			if i%2 == 0 {
				proc.currentPass[i] = CHARSTART
			}

		}
		return
	}
	i := 0
	for {
		proc.currentPass[i*2] = proc.currentPass[i*2] + 1
		if proc.currentPass[i*2] != CHAREND+1 {
			return
		}
		proc.currentPass[i*2] = CHARSTART
		i++
		if i == ln {
			proc.currentPass = append(proc.currentPass, CHARSTART, 0)
			return
		}
	}
}
*/
