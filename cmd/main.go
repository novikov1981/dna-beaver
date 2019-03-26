package main

import (
	"bufio"
	"flag"
	"fmt"
	dna_beaver "github.com/novikov1981/dna-beaver"
	"github.com/novikov1981/dna-beaver/repository/sqllite"
	validation "github.com/novikov1981/dna-beaver/validations"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
	"log"
	"os"
)

func main() {
	var synthesisName = flag.String("name", "generic", "the name of the sysnthesis under interest")
	var synthesisScale = flag.Int64("scale", 1, "the scale of the sysnthesis under interest")
	var filePath = flag.String("path", "", "the file path with the oligs to be parsed and saved")
	var mode = flag.String("mode", dna_beaver.ValidateMode, "the mode to run the application: validate, save, search")
	var searchPattern = flag.String("searchPattern", "", "the pattern of the olig's name to search for in the database")

	flag.Parse()

	// print initial parameters
	log.Printf("start DNA BEAVER APPLICATION")
	defer log.Printf("finish DNA BEAVER APPLICATION")
	log.Printf("synthesis '%s', scale %d, synthesis file '%s' running in mode '%s'", *synthesisName, *synthesisScale, *filePath, *mode)
	if *mode == dna_beaver.SearchMode {
		log.Printf("search pattern '%s'", *searchPattern)
	}

	// flags validation
	if *synthesisScale <= 0 {
		log.Fatalf("error: wrong scale provided %d, should be above zero", *synthesisScale)
	}
	if *filePath == "" {
		log.Fatalf("error: empty file path with synthesis provided %s", *filePath)
	}

	repo, err := sqllite.NewRepository(dna_beaver.SqlLiteFile)
	if err != nil {
		log.Fatal(err)
	}

	validator, err := validation.NewValidator()
	if err != nil {
		log.Fatal(err)
	}

	var oligs []string
	if *mode == dna_beaver.ValidateMode || *mode == dna_beaver.SaveMode {
		oligs, err = readOligsFromFile(*filePath)
		if err != nil {
			log.Fatalf("error: oligs file reading problem: %s", err.Error())
		}
		log.Println("validate initial set of oligs:")
		printOligs(oligs)
		validator.Validate(oligs)
		}
	}
	if *mode == dna_beaver.SaveMode {
		err = repo.InsertSynthesis(*synthesisName, *synthesisScale, oligs)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("synthesis %s saved with success.", *synthesisName)
	}
	if *mode == dna_beaver.SearchMode {
		foundSynthesis, err := repo.FindSynthesis(*searchPattern)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("found %d synthesis containing requested olig pattern %s.\n", len(foundSynthesis), *searchPattern)
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

func printSynthesis(ss []dna_beaver.Synthesis) {
	for _, s := range ss {
		fmt.Printf("Synthes %+v\n", s)
	}
}
