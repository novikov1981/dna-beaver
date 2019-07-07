package main

import (
	"bufio"
	"flag"
	"fmt"
	dna_beaver "github.com/novikov1981/dna-beaver"
	"github.com/novikov1981/dna-beaver/measurements"
	"github.com/novikov1981/dna-beaver/repository/sqllite"
	validation "github.com/novikov1981/dna-beaver/validations"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
	"log"
	"os"
)

func main() {
	validateMode := flag.NewFlagSet("validate", flag.ExitOnError)
	var validateFilePath = validateMode.String("path", "", "the file path with the oligs to be parsed and saved")

	saveMode := flag.NewFlagSet("save", flag.ExitOnError)
	var saveSynthesisName = saveMode.String("name", "generic", "the name of the sysnthesis under interest")
	var saveSynthesisScale = saveMode.Int64("scale", 1, "the scale of the sysnthesis under interest")
	var saveFilePath = saveMode.String("path", "", "the file path with the oligs to be parsed and saved")

	searchMode := flag.NewFlagSet("search", flag.ExitOnError)
	var searchPattern = searchMode.String("pattern", "", "the pattern of the olig's name to search for in the database")

	printMode := flag.NewFlagSet("print", flag.ExitOnError)
	var printSynthesisName = printMode.String("name", "", "the name of the sysnthesis under interest")

	log.Printf("start DNA BEAVER APPLICATION")
	defer log.Printf("finish DNA BEAVER APPLICATION")

	if len(os.Args) < 2 {
		fmt.Println("validate/save/search/print subcommand is required")
		os.Exit(1)
	}

	switch os.Args[1] {
	case dna_beaver.ValidateMode:
		validateMode.Parse(os.Args[2:])
	case dna_beaver.SaveMode:
		saveMode.Parse(os.Args[2:])
	case dna_beaver.SearchMode:
		searchMode.Parse(os.Args[2:])
	case dna_beaver.PrintMode:
		printMode.Parse(os.Args[2:])
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}

	repo, err := sqllite.NewRepository(dna_beaver.SqlLiteFile)
	if err != nil {
		log.Fatal(err)
	}

	/////// VALIDATE
	if validateMode.Parsed() {
		log.Printf("validate synthesis file %s", *validateFilePath)
		if *validateFilePath == "" {
			log.Fatalf("error: empty file path with synthesis provided")
		}
		var oligs []string
		oligs, err = readOligsFromFile(*validateFilePath)
		if err != nil {
			log.Fatalf("error: oligs file reading problem: %s", err.Error())
		}
		log.Println("validate initial set of oligs:")
		printOligs(oligs)
		err = validation.Validate(oligs)
		if err != nil {
			log.Print(err.Error())
		} else {
			log.Print("validated successfully - synthesis does not contain errors")
		}
		statistics := measurements.Measure(oligs)
		log.Printf(`statistics for the synthesis:
                    oligs count   %d
                    wrong symbols %d
                    links number  %d`, statistics.Oligs, statistics.WrongSymbols, statistics.Links)
		log.Printf(`count by every link symbol:
                    %v`, statistics.LinksCount)
	}
	/////// SAVE
	if saveMode.Parsed() {
		log.Printf("save synthesis '%s' from file %s, scale %d", *saveSynthesisName, *saveFilePath, *saveSynthesisScale)
		var oligs []string
		err = repo.InsertSynthesis(*saveSynthesisName, *saveSynthesisScale, oligs)
		if err != nil {
			log.Fatal("cannot save synthesis because of the error: " + err.Error())
		}
		log.Printf("synthesis '%s' saved with success", *saveSynthesisName)
	}
	/////// SEACH
	if searchMode.Parsed() {
		log.Printf("search by pattern '%s' in synthesis name or oligs sequences", *searchPattern)
		foundSynthesis, err := repo.FindSynthesis(*searchPattern)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("found %d synthesis containing requested olig pattern '%s'\n", len(foundSynthesis), *searchPattern)
		printSearchResults(foundSynthesis)
	}
	/////// PRINT
	if printMode.Parsed() {
		log.Printf("print synthesis '%s'", *printSynthesisName)
		if *printSynthesisName == "" {
			log.Fatalf("error: empty sequence name to print provided")
		}
		synt, err := repo.GetSynthesis(*printSynthesisName)
		if err != nil {
			log.Fatal(err)
		}
		if synt == nil {
			log.Fatalf("no sequence found for name '%s'", *printSynthesisName)
		}
		printSynthesis(*synt)
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
		log.Printf("olig %d: %s\n", i+1, o)
	}
}

func printSearchResults(ss []dna_beaver.Synthesis) {
	for _, s := range ss {
		log.Printf("synthesis uuid=%s, name '%s', saved %s, scale %d", s.Uuid, s.Name, s.CreatedAt, s.Scale)
		for _, o := range s.Oligs {
			log.Printf("%d %s", o.Position, o.Content)
		}
	}
}

func printSynthesis(s dna_beaver.Synthesis) {
	log.Printf("synthesis uuid=%s, name '%s', saved %s, scale %d", s.Uuid, s.Name, s.CreatedAt, s.Scale)
	for _, o := range s.Oligs {
		log.Printf("%s", o.Content)
	}
}
