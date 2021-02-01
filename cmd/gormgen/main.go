package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/MohamedBassem/gormgen"
)

type config struct {
	output   string
	structs  []string
	template string
}

var cnf config

func parseFlags() {
	var output, structs, template string
	flag.StringVar(&structs, "structs", "", "[Required] The name of schema structs to generate structs for, comma seperated")
	flag.StringVar(&output, "output", "", "[Required] The name of the output file")
	flag.StringVar(&template, "template", "", "[Optional] Specify the go-template file; default template content (see: https://github.com/MohamedBassem/gormgen/blob/master/template.go)")
	flag.Parse()

	if output == "" || structs == "" {
		flag.Usage()
		os.Exit(1)
	}

	cnf = config{
		output:   output,
		structs:  strings.Split(structs, ","),
		template: template,
	}
}

func main() {
	parseFlags()
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	parser := gormgen.NewParser()
	parser.ParseDir(wd)

	gen := gormgen.NewGenerator(cnf.output, cnf.template)
	if err := gen.Init(parser, cnf.structs); err != nil {
		log.Fatalf("Error Initializing Generator: %v", err.Error())
	}
	if err := gen.Generate(); err != nil {
		log.Fatalf("Error Generating file: %v", err.Error())
	}
	if err := gen.Imports(); err != nil {
		log.Fatalf("Error adding imports to output file: %v", err.Error())
	}
	if err := gen.Format(); err != nil {
		log.Fatalf("Error Formating output file: %v", err.Error())
	}
	if err := gen.Flush(); err != nil {
		log.Fatalf("Error writing output file: %v", err.Error())
	}
}
