package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	// "strings"
)

func main() {
	a := os.Args

	stdin := false
	stdinbuf := bytes.NewBuffer(nil)

	var path string
	var err error
	var by []byte
	if len(a) < 2 {
		stdin = true
		//fmt.Fprintf(os.Stderr, "must supply name of .pony file to format\n")
		//os.Exit(1)
		io.Copy(stdinbuf, os.Stdin) // dest, source
		by = stdinbuf.Bytes()
	} else {
		path = os.Args[1]
		if !fileExists(path) {
			fmt.Fprintf(os.Stderr, "error: file not found for path '%v'\n", path)
			os.Exit(1)
		}
		by, err = os.ReadFile(path)
		panicOn(err)
	}
	newline := []byte("\n")
	byby := bytes.Split(by, newline)
	//vv("len byby = %v", len(byby))
	//vv("byby[0] = '%v'", string(byby[0]))

	indented := make([][]byte, len(byby))
	var spacelead []int
	var prevDefn bool
	var prevLead int
	var accum int

	for i, b := range byby {
		lead := firstNonSpace(b)
		//vv("i= %v; lead = %v, b = '%v'", i, lead, string(b))
		if lead < prevLead {
			accum = 0
		}

		b = chompTrailingWhitespace(b)
		if lead >= 0 { // skip no leading space lines; or all spaces.
			delta := closestStopDelta(lead, 4)
			if i > 0 {
				if spacelead[i-1] == lead+delta && prevDefn {
					accum += 4
				}
			}
			delta += accum
			spacelead = append(spacelead, lead+delta)

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
			spacelead = append(spacelead, 0)
		}
		prevDefn = defn(indented[i])
		prevLead = lead
	}

	// output =======================
	var w io.WriteCloser
	var path1 string
	if stdin {
		w = os.Stdout
	} else {
		path1 = path + ".ponyfmt"
		fd, err := os.Create(path1)
		panicOn(err)
		defer fd.Close()
		w = fd
	}

	path2 := path + ".prev"
	for _, b := range indented {
		_, err = w.Write(b)
		panicOn(err)
		_, err = w.Write(newline)
		panicOn(err)
	}
	if stdin {
		return
	}
	// back up original to .prev
	panicOn(os.Rename(path, path2))
	// put .ponyfmt in place of original path.
	panicOn(os.Rename(path1, path))

}

// ends in =>
func defn(b []byte) bool {
	n := len(b)
	if n < 3 {
		return false
	}
	if b[n-1] == '>' && b[n-2] == '=' {
		return true
	}
	// can we piggy back then too?
	return then(b)
}

// ends in then
func then(b []byte) bool {
	n := len(b)
	if n < 5 {
		return false
	}
	return string(b[n-4:]) == "then"
}

func chompTrailingWhitespace(b []byte) (r []byte) {
	n := len(b)
	if n == 0 {
		return b
	}
	for i := n - 1; i >= 0; i-- {
		switch b[i] {
		case ' ':
		case '\t':
		default:
			return b[:i+1]
		}
	}
	return b
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
