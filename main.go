package main

import "fmt"
import "os"
import "reflect"
import "flag"
import "regexp"
import "strings"

import (
	"github.com/DinoInc/DomainInflator/Converter"
	"github.com/DinoInc/DomainInflator/Utils"
)

var _ = reflect.TypeOf

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var baseDir string
var namespaces arrayFlags
var currentFile string

func isValidIdentifier(s string) bool {
	isMatch, _ := regexp.MatchString(`^[A-z_][A-z0-9._]*$`, s)
	return isMatch
}

func main() {

	flag.Var(&namespaces, "namespace", "thrift namespace")

	pBaseDir := flag.String("schema-dir", "./schemas/", "JSON schema Directory")
	pSchemas := flag.String("schema", "", "schemas to resolve")
	pThriftFile := flag.String("thrift-file", "", "thrift file")
	pSedFile := flag.String("sed-file", "", "sed file")
	flag.Parse()

	if *pThriftFile == "" {
		fmt.Fprintf(os.Stderr, "missing required --thrift-file argument/flag\n")
		os.Exit(1)
	}

	if *pSedFile == "" {
		fmt.Fprintf(os.Stderr, "missing required --sed-file argument/flag\n")
		os.Exit(1)
	}

	if *pSchemas == "" {
		fmt.Fprintf(os.Stderr, "missing required --schema argument/flag\n")
		os.Exit(1)
	}

	if len(namespaces) == 0 {
		fmt.Fprintf(os.Stderr, "missing required --namespace argument/flag\n")
		os.Exit(1)
	}

	schemaList := strings.Split(*pSchemas, ",")

	engine := Converter.NewConverter(*pBaseDir)

	if Utils.FileExists(*pThriftFile) {
		engine.ReadIDL(*pThriftFile)
	} else {
		engine.NewIDL()
	}

	for _, schema := range schemaList {
		engine.ResolveDefinitionOf(schema)
		engine.Convert()
	}

	thriftFile, _ := os.Create(*pThriftFile)
	defer thriftFile.Close()

	sedFile, _ := os.Create(*pSedFile)
	defer sedFile.Close()

	for _, ns := range namespaces {
		s := strings.Split(ns, ":")

		language := s[0]
		namespace := s[1]

		fmt.Fprintf(thriftFile, "namespace %s %s\n", language, namespace)
	}
	fmt.Fprintf(thriftFile, "\n")

	thriftFile.WriteString(engine.Thrift())
	sedFile.WriteString(engine.Deviations())
}
