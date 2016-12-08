package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	var src string
	if len(os.Args) > 1 {
		src = os.Args[1]
	} else {
		src = "."
	}
	fs, err := ioutil.ReadDir(src)
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
