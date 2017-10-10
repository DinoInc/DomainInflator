package main

import "fmt"
import "os"
import "reflect"
import "flag"
import "regexp"
import "strings"

import (
	"github.com/DinoInc/DomainInflator/Converter"
)

var _ = reflect.TypeOf

var baseDir string
var namespace string
var currentFile string

func isValidIdentifier(s string) bool {
	isMatch, _ := regexp.MatchString(`^[A-z_][A-z0-9._]*$`, s)
	return isMatch
}

func main() {

	pBaseDir := flag.String("schema-dir", "./schemas/", "JSON schema Directory")
	pSchemas := flag.String("schema", "", "schemas to resolve")
	pNamespace := flag.String("namespace", "", "thrift namespace")
	flag.Parse()

	if *pSchemas == "" {
		fmt.Fprintf(os.Stderr, "missing required --schema argument/flag\n")
		os.Exit(1)
	}

	if *pNamespace == "" {
		fmt.Fprintf(os.Stderr, "missing required --namespace argument/flag\n")
		os.Exit(1)
	}

	schemaList := strings.Split(*pSchemas, ",")

	engine := Converter.NewConverter(*pBaseDir)
	engine.NewIDL()
	for _, schema := range schemaList {
		engine.ResolveDefinitionOf(schema)
		engine.Convert()
	}

	fmt.Println(engine.Thrift())

}
