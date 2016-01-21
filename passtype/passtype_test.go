// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package passtype

import (
	"bytes"
	"fmt"
	"testing"
)

type passCalcul struct {
	slice1, slice2, result []byte
	myfunc                 func([]byte, []byte) []byte
}

type passConv struct {
	slice1, slice2 []byte
	myfunc         func([]byte) []byte
}

var (
	tmpCharset = []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'A', 'B', 'C', 'D', 'E', 'F'}
	calculs    = []passCalcul{
		{[]byte{10, 8, 8, 8}, []byte{0, 2, 2, 2}, []byte{10, 10, 10, 10}, addPass},
		{[]byte{10, 8, 8, 8}, []byte{0, 7, 7, 8}, []byte{11, 0, 0, 0}, addPass},
		{[]byte{10, 8, 8, 8}, []byte{7, 7, 8}, []byte{11, 0, 0, 0}, addPass},
	}
	conversions = []passConv{
		{[]byte{5, 12, 3, 7}, []byte{'5', 0, 'C', 0, '3', 0, '7', 0}, passToTruePass},  // pass to true pass
		{[]byte{'A', 0, 'B', 0, '3', 0, '7', 0}, []byte{10, 11, 3, 7}, truePassToPass}, // true pass to pass
	}
)

func TestConversion(t *testing.T) {
	Charset = tmpCharset
	for i, jeu := range conversions {
		tmp := jeu.myfunc(jeu.slice1)
		if bytes.Compare(tmp, jeu.slice2) != 0 {
			t.Fatalf("jeu d'essai N°%d des tests de conversion incorrecte.\n%v n'est pas converti en %v mais en %v", i, jeu.slice1, jeu.slice2, tmp)
		}
	}
}

func TestCalcul(t *testing.T) {
	Charset = tmpCharset
	for i, jeu := range calculs {
		tmpslice1 := make([]byte, len(jeu.slice1))
		copy(tmpslice1, jeu.slice1)
		tmp := jeu.myfunc(jeu.slice1, jeu.slice2)
		if bytes.Compare(tmp, jeu.result) != 0 {
			t.Fatalf("jeu d'essai N°%d des tests de calcul incorrecte.\n%v\net %v\nne donne pas %v\nmais %v", i, tmpslice1, jeu.slice2, jeu.result, tmp)
		}
	}
}
