package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/MohamedBassem/gormgen"
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

	// Remove file if it's already there to avoid parsing it
	if _, err := os.Stat(cnf.output); err == nil {
		if err := os.Remove(cnf.output); err != nil {
			log.Fatalf("Couldn't remove file %v: %v", cnf.output, err.Error())
		}
	}

	parser := gormgen.NewParser()
	parser.ParseDir(wd)

	gen := gormgen.NewGenerator(cnf.output)
	if err := gen.Init(parser, cnf.structs); err != nil {
		log.Fatalf("Error Initializing Generator: %v", err.Error())
	}
	if err := gen.Generate(); err != nil {
		log.Fatalf("Error Generating file: %v", err.Error())
	}
	if err := gen.Format(); err != nil {
		log.Fatalf("Error Formating output file: %v", err.Error())
	}
	if err := gen.Flush(); err != nil {
		log.Fatalf("Error writing output file: %v", err.Error())
	}
}
