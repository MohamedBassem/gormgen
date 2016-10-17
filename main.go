package main

import (
	"flag"
	"log"
	"os"
	"strings"
)

type config struct {
	output  string
	structs []string
}

var cnf config

func parseFlags() {
	var output, structs string
	flag.StringVar(&structs, "structs", "", "[Required] The name of schema structs to generate structs for, comma seperated")
	flag.StringVar(&output, "output", "", "[Required] The name of the output file")
	flag.Parse()

	if output == "" || structs == "" {
		flag.Usage()
		os.Exit(1)
	}

	cnf = config{
		output:  output,
		structs: strings.Split(structs, ","),
	}
}

func checkStructSuffix() {
	for _, st := range cnf.structs {
		if !strings.HasSuffix(st, "Schema") {
			log.Fatalf("Struct %s must have 'Schema' as a suffix", st)
		}
	}
}

func main() {
	parseFlags()
	checkStructSuffix()
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}
	parser := NewParser()
	parser.ParseDir(wd)

	generator := NewGenerator()
	generator.init()
	for _, s := range cnf.structs {
		ast := parser.GetTypeByName(s)
		if ast == nil {
			log.Printf("Struct %s is not found\n", s)
			continue
		}
		generator.Generate(parser, s)
	}
	generator.Flush()
}
