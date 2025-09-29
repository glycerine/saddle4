package main

import (
	"bytes"
	"fmt"
	"os"
	// "strings"
)

func main() {
	a := os.Args
	if len(a) < 2 {
		fmt.Fprintf(os.Stderr, "must supply name of .pony file to format\n")
		os.Exit(1)
	}
	path := os.Args[1]
	if !fileExists(path) {
		fmt.Fprintf(os.Stderr, "error: file not found for path '%v'\n", path)
		os.Exit(1)
	}
	by, err := os.ReadFile(path)
	panicOn(err)
	newline := []byte("\n")

	byby := bytes.Split(by, newline)
	//vv("len byby = %v", len(byby))
	//vv("byby[0] = '%v'", string(byby[0]))

	indented := make([][]byte, len(byby))
	for i, b := range byby {
		lead := firstNonSpace(b)
		//vv("i= %v; lead = %v, b = '%v'", i, lead, string(b))

		if lead >= 0 { // skip no leading space lines; or all spaces.
			delta := closestStopDelta(lead, 4)
			//vv("i = %v; delta = %v", i, delta)
			switch {
			case delta < 0:
				indented[i] = append([]byte{}, b[-delta:]...)
			case delta > 0:
				amt := bytes.Repeat([]byte(" "), delta)
				indented[i] = append(amt, b...)
			case delta == 0:
				indented[i] = append([]byte{}, b...)
			}
		} else {
			indented[i] = append([]byte{}, b...)
		}
	}
	path1 := path + ".ponyfmt"
	path2 := path + ".prev"
	fd, err := os.Create(path1)
	panicOn(err)
	defer fd.Close()
	for _, b := range indented {
		_, err = fd.Write(b)
		panicOn(err)
		_, err = fd.Write(newline)
		panicOn(err)
	}
	// back up original to .prev
	panicOn(os.Rename(path, path2))
	// put .ponyfmt in place of original path.
	panicOn(os.Rename(path1, path))
}

func closestStopDelta(n, tabstop int) int {
	r := n % tabstop
	if r == 0 {
		return 0
	}
	if n > 0 && n < tabstop {
		// move anything indented at all to
		// at least the first tabstop
		return tabstop - n
	}
	// INVAR: after first tabstop: n > tabstop
	k := n / tabstop
	prev := k * tabstop
	next := (k + 1) * tabstop
	dist0 := n - prev
	dist1 := next - n
	if dist1 <= dist0 {
		// bump forward to next, delta > 0
		return next - n
	}
	// pull back to prev, negative delta
	return prev - n
}

func firstNonSpace(by []byte) int {
	for i, c := range by {
		if c != ' ' {
			return i
		}
	}
	return -1
}

func fileExists(name string) bool {
	fi, err := os.Stat(name)
	if err != nil {
		return false
	}
	if fi.IsDir() {
		return false
	}
	return true
}

func dirExists(name string) bool {
	fi, err := os.Stat(name)
	if err != nil {
		return false
	}
	if fi.IsDir() {
		return true
	}
	return false
}

func fileSize(name string) (int64, error) {
	fi, err := os.Stat(name)
	if err != nil {
		return -1, err
	}
	return fi.Size(), nil
}

// func panicOn(err error) {
// 	if err != nil {
// 		panic(err)
// 	}
// }
