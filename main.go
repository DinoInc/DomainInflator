package main

import "encoding/json"
import "io/ioutil"
import "fmt"
import "os"
import "reflect"
import "flag"
import "regexp"
import "strings"

import (
	"github.com/DinoInc/DomainInflator/Engine"
	"github.com/DinoInc/DomainInflator/Schema"
)

var _ = reflect.TypeOf

// [v] allOf
// [ ] anyOf
// [v] primitive type
// [v] array -> ref
// [v] array -> primitive
// [ ]

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

	engine := Engine.NewEngine(*pBaseDir)

	for _, schema := range schemaList {

		baseDir = *pBaseDir
		currentFile = schema
		namespace = *pNamespace

		content, e := ioutil.ReadFile(baseDir + currentFile + ".schema.json")
		if e != nil {
			fmt.Printf("File error: %v\n", e)
			os.Exit(1)
		}

		jsonContent := json.RawMessage(content)
		ref, isRef := Schema.ReadRef(&jsonContent)

		if !isRef {
			panic("not implemented")
		}

		engine.SetCurrentFile(currentFile)
		engine.Resolve(ref)
		engine.Print()

	}

}
