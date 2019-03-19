package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/novikov1981/experiments"
	"github.com/novikov1981/experiments/repository/sqllite"
	"github.com/novikov1981/experiments/validation"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
	"log"
	"os"
)

const (
	sqlLiteFile  = "./synthesis.db"
	validateMode = "validate"
	saveMode     = "save"
	searchMode   = "search"
)

func main() {
	var synthesisName = flag.String("name", "generic", "the name of the sysnthesis under interest")
	var synthesisScale = flag.Int64("scale", 1, "the scale of the sysnthesis under interest")
	var filePath = flag.String("path", "", "the file path with the oligs to be parsed and saved")
	var mode = flag.String("mode", validateMode, "the mode to run the application: validate, save, search")
	var oligPattern = flag.String("oligPattern", "", "the pattern of the olig's name to search for in the database")

	flag.Parse()

	fmt.Printf("synthesis '%s', scale %d from file '%s' running in mode '%s', search pattern '%s' (if search mode)\n", *synthesisName, *synthesisScale, *filePath, *mode, *oligPattern)

	repo, err := sqllite.NewRepository(sqlLiteFile)
	if err != nil {
		log.Fatal(err)
	}

	validator, err := validation.NewValidator()
	if err != nil {
		log.Fatal(err)
	}

	switch *mode {
	case "validate":
		oligs, err := readOligsFromFile(*filePath)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Validate initial set of oligs:")
		printOligs(oligs)
		if err = validator.Validate(oligs); err != nil {
			log.Printf("Validation finished with error: %s\n", err.Error())
		}
	case "save":
		oligs, err := readOligsFromFile(*filePath)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Validate initial set of oligs before save:")
		printOligs(oligs)
		if err = validator.Validate(oligs); err != nil {
			log.Printf("Validation finished with error: %s. ATTENTION, cannot save the oligs!\n", err.Error())
			return
		}
		err = repo.InsertSynthesis(*synthesisName, *synthesisScale, oligs)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Synthesis %s saved with success.", *synthesisName)
	case "search":
		foundSynthesis, err := repo.FindSynthesis(*oligPattern)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Found %d synthesis containing requested olig pattern %s.\n", len(foundSynthesis), *oligPattern)
		printSynthesis(foundSynthesis)
	}
}

func readOligsFromFile(filePath string) ([]string, error) {
	oligs := make([]string, 0)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	utf8Reader := transform.NewReader(file, charmap.Windows1251.NewDecoder())
	scanner := bufio.NewScanner(utf8Reader)
	for scanner.Scan() {
		oligs = append(oligs, scanner.Text())
	}
	return oligs, nil
}

func printOligs(oo []string) {
	for i, o := range oo {
		fmt.Printf("Olig %d: %s\n", i+1, o)
	}
}

func printSynthesis(ss []experiments.Synthesis) {
	for _, s := range ss {
		fmt.Printf("Synthes %+v\n", s)
	}
}
