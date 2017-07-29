package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("One argument, the json file to pretty-print is required")
		os.Exit(-1)
	}

	fileName := os.Args[1]
	byt, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	var obj map[string]interface{}
	var arr []map[string]interface{}
	if err := json.Unmarshal(byt, &obj); err != nil {
		if err := json.Unmarshal(byt, &arr); err != nil {
			panic(err)
		}
		pp(arr)
		return
	}
	pp(obj)

}
func pp(v interface{}) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic(err)
	}
	b2 := append(b, '\n')
	os.Stdout.Write(b2)
}
