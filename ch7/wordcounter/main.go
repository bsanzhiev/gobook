package main

import (
	"bufio"
	"bytes"
	"fmt"
	"unicode"
	"unicode/utf8"
)

type WordCounter int

func ScanWords(data []byte, atEOF bool) (advance int, token []byte, err error) {
	start := 0
	for width := 0; start < len(data); start += width {
		var r rune
		r, width = utf8.DecodeRune(data[start:])
		if !unicode.IsSpace(r) && !unicode.IsPunct(r) {
			break
		}
	}

	// Scan until space, marking end of word
	for width, i := 0, start; i < len(data); i += width {
		var r rune
		r, width = utf8.DecodeRune(data[i:])
		if unicode.IsSpace(r) || unicode.IsPunct(r) {
			return i + width, data[start:i], nil
		}
	}

	// If we're at EOF, we have a final, non-empty, non-terminated ward.
	// Return it
	if atEOF && len(data) > start {
		return len(data), data[start:], nil
	}
	// Request more data
	return start, nil, nil
}

func (c *WordCounter) Write(p []byte) (int, error) {
	s := bufio.NewScanner(bytes.NewBuffer(p))

	// word split
	s.Split(ScanWords)
	for s.Scan() {
		// debug only
		// fmt.Primtln(s.Text())
		*c++
	}
	return len(p), s.Err()
}

func (c *WordCounter) String() string {
	return fmt.Sprintf("%d word(s)", *c)
}

// Line counter

type LineCounter int

func (l *LineCounter) Write(p []byte) (int, error) {
	s := bufio.NewScanner(bytes.NewBuffer(p))

	for s.Scan() {
		*l++
	}
	return len(p), s.Err()
}

func (l *LineCounter) String() string {
	return fmt.Sprintf("%d line(s)", *l)
}

func main() {
	var c WordCounter

	c.Write([]byte("hello, 世界! 哈哈,哈"))
	fmt.Println(&c)

	c = 0
	var name = "Dolly"
	fmt.Fprintf(&c, "hello, %s", name)
	fmt.Println(&c)

	var l LineCounter
	l.Write([]byte("你\n好，世\n界\n\n你好\n吗，世界\n世界\n世界1111\n"))
	fmt.Println(&l)
}
