// You can edit this code!
// Click here and start typing.
package main

import "fmt"

type S struct {
	x int
}

func one() *S {
	return &S{-1}
}

func add() S {
	q := S{1}
	q.x += 3
	return q
}

func another( s *S ) S {
	s.x += 3
	return *s

}

func main() {
	one().x = 5
	add().x = 5
	
}