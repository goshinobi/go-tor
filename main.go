package main

import (
	"fmt"
	"log"

	"github.com/goshinobi/go-tor/tor"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	fmt.Println(tor.New("hoge"))
}
