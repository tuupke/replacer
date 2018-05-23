package main

import (
	"fmt"
	"flag"
	"os"
	"io/ioutil"
)

func main() {
	flag.Parse()
	if fileBytes, err := ioutil.ReadFile(flag.Arg(0)); err == nil {
		fmt.Println(os.ExpandEnv(string(fileBytes)))
	}
}

