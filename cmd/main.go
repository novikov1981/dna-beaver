package main

import (
	"bufio"
	"flag"
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
	var synthesisName = flag.String("name", "generic", "the name of the sysnthesis under interest")
	var synthesisScale = flag.Int64("scale", 1, "the scale of the sysnthesis under interest")
	var filePath = flag.String("path", "", "the file path with the oligs to be parsed and saved")
	var mode = flag.String("mode", dna_beaver.ValidateMode, "the mode to run the application: validate, save, search")
	var searchPattern = flag.String("searchPattern", "", "the pattern of the olig's name to search for in the database")
	var force = flag.Bool("forceSave", false, "force saving of non-validated synthesis")

	flag.Parse()

	// print initial parameters
	log.Printf("start DNA BEAVER APPLICATION")
	defer log.Printf("finish DNA BEAVER APPLICATION")
	log.Printf("synthesis '%s', scale %d, synthesis file '%s' running in mode '%s'", *synthesisName, *synthesisScale, *filePath, *mode)

	// flags validation
	if *synthesisScale <= 0 {
		log.Fatalf("error: wrong scale provided %d, should be above zero", *synthesisScale)
	}

	repo, err := sqllite.NewRepository(dna_beaver.SqlLiteFile)
	if err != nil {
		log.Fatal(err)
	}

	validated := true
	if *mode == dna_beaver.ValidateMode || *mode == dna_beaver.SaveMode {
		if *filePath == "" {
			log.Fatalf("error: empty file path with synthesis provided")
		}
		var oligs []string
		oligs, err = readOligsFromFile(*filePath)
		if err != nil {
			log.Fatalf("error: oligs file reading problem: %s", err.Error())
		}
		log.Println("validate initial set of oligs:")
		printOligs(oligs)
		err = validation.Validate(oligs)
		if err != nil {
			log.Print(err.Error())
			validated = false
		} else {
			log.Print("validated successfully - synthesis does not contain errors")
		}
		statistics := measurements.Measure(oligs)
		log.Printf("statistics for the synthesis:")
		log.Printf("oligs number %d, wrong symbols %d, links number %d", statistics.Oligs, statistics.WrongSymbols, statistics.Links)
		log.Printf("count by every link symbol %v", statistics.LinksCount)
	}
	if *mode == dna_beaver.SaveMode {
		if !validated {
			if !*force {
				log.Fatalf("synthesis will not be saved because of validation errors")
			} else {
				log.Print("force saving non validated synthesis")
			}
		}
		var oligs []string
		err = repo.InsertSynthesis(*synthesisName, *synthesisScale, oligs)
		if err != nil {
			log.Fatal("cannot save synthesis because of the error: " + err.Error())
		}
		log.Printf("synthesis '%s' saved with success", *synthesisName)
	} else if *mode == dna_beaver.SearchMode {
		log.Printf("search pattern '%s'", *searchPattern)
		foundSynthesis, err := repo.FindSynthesis(*searchPattern)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("found %d synthesis containing requested olig pattern '%s'\n", len(foundSynthesis), *searchPattern)
		printSearchResults(foundSynthesis)
	} else if *mode == dna_beaver.PrintMode {
		if *synthesisName == "" {
			log.Fatalf("error: empty sequence name to print provided")
		}
		synt, err := repo.GetSynthesis(*synthesisName)
		if err != nil {
			log.Fatal(err)
		}
		if synt == nil {
			log.Fatalf("no sequence found for name '%s'", *synthesisName)
		}
		log.Printf("print sequence '%s'", *synthesisName)
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
		log.Printf("synthesis uuid=%s, name '%s', saved %s, scale %d oligs", s.Uuid, s.Name, s.CreatedAt, s.Scale)
		for _, o := range s.Oligs {
			log.Printf("%d %s", o.Position, o.Content)
		}
	}
}

func printSynthesis(s dna_beaver.Synthesis) {
	log.Printf("synthesis uuid=%s, name '%s', saved %s, scale %d oligs", s.Uuid, s.Name, s.CreatedAt, s.Scale)
	for _, o := range s.Oligs {
		log.Printf("%s", o.Content)
	}
}
