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

	fmt.Printf("start DNA BEAVER APPLICATION\n\n")
	defer fmt.Printf("finish DNA BEAVER APPLICATION")

	if len(os.Args) < 2 {
		fmt.Println("validate/save/search/print command is required")
		return
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
		return
	}

	repo, err := sqllite.NewRepository(dna_beaver.SqlLiteFile)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	/////// VALIDATE
	if validateMode.Parsed() {
		fmt.Printf("validate synthesis file %s\n", *validateFilePath)
		if *validateFilePath == "" {
			fmt.Println("error: empty file path with synthesis provided\n")
			return
		}
		var oligs []string
		oligs, err = readOligsFromFile(*validateFilePath)
		if err != nil {
			fmt.Printf("error: oligs file reading problem - %s\n\n", err.Error())
			return
		}
		fmt.Println("validate initial set of oligs:")
		printOligs(oligs)
		fmt.Println()
		err = validation.Validate(oligs)
		if err != nil {
			fmt.Print(err.Error())
		} else {
			fmt.Println("validated successfully - synthesis does not contain errors\n")
		}
		statistics := measurements.Measure(oligs)
		fmt.Printf("statistics for the synthesis: \noligs count   %d\nwrong symbols %d\nlinks number  %d\n\n",
			statistics.Oligs, statistics.WrongSymbols, statistics.Links)
		fmt.Printf("count by every link symbol:\n")
		for _, r := range dna_beaver.ValidNotations {
			fmt.Printf("%s:%d; ", string(r), statistics.LinksCount[string(r)])
		}
		fmt.Println()
		fmt.Printf("\ncount by amedit:\n")
		for _, amedit := range dna_beaver.Amedits {
			fmt.Printf("%s:%.2f; ", amedit, statistics.AmeditCount[amedit])
		}
		fmt.Println("\n")
	}
	/////// SAVE
	if saveMode.Parsed() {
		fmt.Printf("save synthesis '%s' from file %s, scale %d\n", *saveSynthesisName, *saveFilePath, *saveSynthesisScale)
		var oligs []string
		oligs, err = readOligsFromFile(*saveFilePath)
		if err != nil {
			fmt.Printf("error: oligs file reading problem - %s\n\n", err.Error())
			return
		}
		err = repo.InsertSynthesis(*saveSynthesisName, *saveSynthesisScale, oligs)
		if err != nil {
			fmt.Printf("cannot save synthesis because of the error: %s\n\n", err.Error())
			return
		}
		fmt.Printf("synthesis '%s' saved with success\n\n", *saveSynthesisName)
	}
	/////// SEACH
	if searchMode.Parsed() {
		fmt.Printf("search by pattern '%s' in synthesis name or oligs sequences\n", *searchPattern)
		foundSynthesis, err := repo.FindSynthesis(*searchPattern)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Printf("found %d synthesis containing requested olig pattern '%s'\n", len(foundSynthesis), *searchPattern)
		printSearchResults(foundSynthesis)
		fmt.Println()
	}
	/////// PRINT
	if printMode.Parsed() {
		fmt.Printf("print synthesis '%s'\n", *printSynthesisName)
		if *printSynthesisName == "" {
			fmt.Println("error: empty sequence name to print provided")
		}
		synt, err := repo.GetSynthesis(*printSynthesisName)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		if synt == nil {
			fmt.Printf("no sequence found for name '%s'\n", *printSynthesisName)
			return
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
		fmt.Printf("olig %d: %s\n", i+1, o)
	}
}

func printSearchResults(ss []dna_beaver.Synthesis) {
	for _, s := range ss {
		fmt.Printf("synthesis uuid=%s, name '%s', saved %s, scale %d\n", s.Uuid, s.Name, s.CreatedAt, s.Scale)
		for _, o := range s.Oligs {
			fmt.Printf("%d %s\n", o.Position, o.Content)
		}
	}
}

func printSynthesis(s dna_beaver.Synthesis) {
	fmt.Printf("synthesis uuid=%s, name '%s', saved %s, scale %d\n", s.Uuid, s.Name, s.CreatedAt, s.Scale)
	for _, o := range s.Oligs {
		fmt.Printf("%s\n", o.Content)
	}
	fmt.Println()
}
