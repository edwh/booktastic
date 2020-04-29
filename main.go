package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"log"
)

var sugar *zap.SugaredLogger

func init() {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar = logger.Sugar()
}

func main() {
	verbosePtr := flag.Bool("v", false, "Debug logging")
	inputPtr := flag.String("i", "", "Input file")
	outputPtr := flag.String("o", "", "Output file")

	flag.Parse()

	if !*verbosePtr {
		// Turn off logging.
		log.SetOutput(ioutil.Discard)
		log.SetFlags(0)
	}

	fmt.Println("Input ", *inputPtr, " output ", *outputPtr)

	if len(*inputPtr) > 0 && len(*outputPtr) > 0 {
		data, _ := ioutil.ReadFile(*inputPtr)

		lines, fragments := GetLinesAndFragments(string(data))
		spines, fragments := ExtractSpines(lines, fragments)
		spines, fragments = IdentifyBooks(spines, fragments)

		type output struct {
			spines    []Spine
			fragments []OCRFragment
		}

		outputVal := output{
			spines:    spines,
			fragments: fragments,
		}

		encoded, _ := json.MarshalIndent(outputVal, "", " ")
		_ = ioutil.WriteFile(*outputPtr, encoded, 0644)
	}
}
