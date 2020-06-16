package main

import (
	"fmt"

	"github.com/mattn/go-runewidth"
)

func main()  {
	fmt.Printf("tewst= %d\n", runewidth.RuneWidth('\t'))
}