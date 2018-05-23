package main

import (
	"flag"
	"os"
	"io/ioutil"
)

func main() {
	flag.Parse()
	if fileBytes, err := ioutil.ReadFile(flag.Arg(0)); err == nil {
		ioutil.WriteFile(flag.Arg(1), []byte(os.ExpandEnv(string(fileBytes))), os.ModePerm)
	}
}

