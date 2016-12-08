package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	fs, err := ioutil.ReadDir(os.Args[1])
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, f := range fs {
		str := f.Name()
		if f.IsDir() {
			str += "/"
		}
		fmt.Println(str)
	}
}
