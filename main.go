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

func main() {
	parseFlags()
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}
	parser := NewParser()
	parser.ParseDir(wd)

	gen := NewGenerator(cnf.output)
	if err := gen.init(parser, cnf.structs); err != nil {
		log.Fatalf(err.Error())
	}
	if err := gen.Generate(); err != nil {
		log.Fatalf(err.Error())
	}
	if err := gen.Format(); err != nil {
		log.Fatalf(err.Error())
	}
	gen.Flush()

}
