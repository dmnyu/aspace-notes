package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"

	aspace "github.com/nyudlts/go-aspace"
)

var (
	env        string
	config     string
	repoID     int
	helpmsg    bool
	outputFile string
	version    = "v0.1.0"
)

func init() {
	flag.StringVar(&env, "env", "", "")
	flag.StringVar(&config, "config", "", "")
	flag.IntVar(&repoID, "repo-id", 2, "")
	flag.BoolVar(&helpmsg, "help", false, "")
	flag.StringVar(&outputFile, "output-file", "output.csv", "")
}

func help() {
	fmt.Println("Usage: aspace-notes [options]")
	fmt.Println("Options:")
	fmt.Println("  --config path/to/ go-aspace config file, Mandatory")
	fmt.Println("  --env the environment to run program against, Mandatory")
	fmt.Println("  --repo-id the id of the repository to run program against, default: `2`")
	fmt.Println("  --output-file the path/to the output file, default: output.csv")
	fmt.Println("  --help print this help message")
	os.Exit(0)
}

func main() {
	fmt.Printf("aspace-notes, %v", version)
	//parse flags
	flag.Parse()

	//check for help
	if helpmsg {
		help()
	}

	//get client
	client, err := aspace.NewClient(config, env, 20)
	if err != nil {
		panic(err)
	}

	//create output file
	outfile, err := os.Create(outputFile)
	if err != nil {
		panic(err)
	}
	defer outfile.Close()

	//create a csv writer
	writer := csv.NewWriter(outfile)
	writer.Write([]string{"ID", "TITLE", "NOTE"})

	//request resource IDs for repository
	resourceIds, err := client.GetResourceIDs(repoID)
	if err != nil {
		panic(err)
	}

	//loop through resource ids
	for _, resourceID := range resourceIds {
		//Get a resource
		resource, err := client.GetResource(repoID, resourceID)
		if err != nil {
			panic(err)
		}
		fmt.Printf("checking %s %s\n", resource.MergeIDs("."), resource.Title)

		//get notes
		notes := resource.Notes
		selectedNotes := []aspace.Note{}
		for _, note := range notes {
			if strings.Contains(note.Label, "Conditions Governing Access") {
				selectedNotes = append(selectedNotes, note)
			}
		}

		//if matching notes were found
		if len(selectedNotes) > 0 {
			var collectedNotes = []string{}

			//gather the notes
			for _, note := range selectedNotes {
				for _, subnote := range note.Subnotes {
					collectedNotes = append(collectedNotes, strings.ReplaceAll(subnote.Content, "\n", ""))
				}
			}

			//write to output file
			line := []string{resource.MergeIDs("."), resource.Title}
			line = append(line, collectedNotes...)
			writer.Write(line)
			writer.Flush()
		}

	}
	writer.Flush() //just in case
}
