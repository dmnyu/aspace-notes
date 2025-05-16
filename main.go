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
	env    string
	config string
	repoID int
)

func init() {
	flag.StringVar(&env, "env", "", "")
	flag.StringVar(&config, "config", "", "")
	flag.IntVar(&repoID, "repo-id", 0, "")
}

func main() {

	//parse flags
	flag.Parse()

	//get client
	client, err := aspace.NewClient(config, env, 20)
	if err != nil {
		panic(err)
	}

	//create output file
	outfile, err := os.Create("out.csv")
	if err != nil {
		panic(err)
	}
	defer outfile.Close()
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
				for _, v := range note.Subnotes {
					collectedNotes = append(collectedNotes, strings.ReplaceAll(v.Content, "\n", ""))
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
