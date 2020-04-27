package main

import (
	"io/ioutil"
	"log"
)

func main() {
	log.SetOutput(ioutil.Discard)
	log.SetFlags(0)
}
