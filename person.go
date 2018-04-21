package solveuranimateur

import (
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type PersonIndex int
type personValue string

type personsMap map[PersonIndex]personValue

func packagePath() string {
	return os.Getenv("GOPATH") + "/src/github.com/jmichiels/SolveurAnimateur/"
}

func GeneratePersonsDefinitions() error {
	transformer := transform.Chain(norm.NFD,
		runes.Remove(runes.In(unicode.Pd)),
		runes.Remove(runes.In(unicode.Mn)),
		runes.Remove(runes.In(unicode.Zs)))

	var def string
	for index, person := range Persons {
		constName, _, _ := transform.String(transformer, string(person))
		def += fmt.Sprintf("%s personIndex = %d\n", constName, index)
	}
	file, err := format.Source([]byte("package solveuranimateur \n const (\n" + def + ")"))
	if err != nil {
		return err
	}
	return ioutil.WriteFile(packagePath()+"person_data_gen.go", file, 0666)
}
