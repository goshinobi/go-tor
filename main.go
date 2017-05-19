package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"time"

	"github.com/goshinobi/go-tor/tor"
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	runtime.GOMAXPROCS(runtime.NumCPU())
	rand.Seed(time.Now().UnixNano())
	if t, ok := http.DefaultClient.Transport.(*http.Transport); ok {
		t.MaxIdleConnsPerHost = 4
	}
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
