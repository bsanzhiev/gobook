// Давайте проверим это, используя новый тип. 
// Метод Write типа *ByteCounter ниже просто подсчитывает 
// записанные в него байты перед их отбрасыванием. 
// (Преобразование необходимо, чтобы типы len(p) 
// и *c совпадали в операторе присваивания +=.)

// Bytecounter demonstrates an implementation 
// of io.Writer that counts bytes.
package main

import (
	"fmt"
)

type ByteCounter int

func (c *ByteCounter) Write(p []byte) (int, error) {
	*c += ByteCounter(len(p))
	return len(p), nil
}

func main() {
	var c ByteCounter
	c.Write([]byte("hello"))
	fmt.Println(c)

	c = 0
	var name = "Dolly"
	fmt.Fprintf(&c, "hello, %s", name)
	fmt.Println(c)
}
