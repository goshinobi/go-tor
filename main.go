package main

import (
	"fmt"
	"log"
	"time"

	"github.com/goshinobi/go-tor/tor"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	t := tor.New("hoge")
	fmt.Println(t)

	fmt.Println("start")
	if err := t.Start(); err != nil {
		log.Fatalln(err)
	}
	defer t.Kill()
	time.Sleep(60 * time.Second)
	fmt.Println("reload")
	t.Reload()
	time.Sleep(60 * time.Second)
}
