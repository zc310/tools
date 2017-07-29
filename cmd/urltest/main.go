package main

import (
	"fmt"

	"github.com/zc310/log"

	"io/ioutil"
	"strings"
	"sync"

	"net/http"
	"time"
)

func main() {
	wg1 := &sync.WaitGroup{}

	f := func(c rune) bool {
		return (c == '\n') || (c == '\r')
	}
	b, err := ioutil.ReadFile("url.txt")
	if err != nil {
		log.Fatal(err)
	}
	list := strings.FieldsFunc(string(b), f)

	tm1 := time.Now()
	wg1.Add(len(list))
	for i := 0; i < len(list); i++ {
		go func(s string) {
			resp, err := http.Get(s)
			if resp != nil {
				defer resp.Body.Close()
			}
			if err != nil {
				log.Print(s, err)
			} else {
				log.Print(s, resp.StatusCode)
			}
			wg1.Done()
		}(list[i])
	}
	wg1.Wait()
	fmt.Println(time.Now().Sub(tm1))
}
