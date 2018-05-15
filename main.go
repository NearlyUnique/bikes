package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

type (
	// TflPlace where bikes live
	TflPlace struct {
		Type                 string            `json:"$type"`
		ID                   string            `json:"id"`
		URL                  string            `json:"url"`
		CommonName           string            `json:"commonName"`
		PlaceType            string            `json:"placeType"`
		AdditionalProperties []AdditionalProp  `json:"additionalProperties"`
		Children             []json.RawMessage `json:"children"`
		ChildrenUrls         []json.RawMessage `json:"childrenUrls"`
		Lat                  float64           `json:"lat"`
		Lon                  float64           `json:"lon"`
	}
	// AdditionalProp data
	AdditionalProp struct {
		Type            string    `json:"$type"`
		Category        string    `json:"category"`
		Key             string    `json:"key"`
		SourceSystemKey string    `json:"sourceSystemKey"`
		Value           string    `json:"value"`
		Modified        time.Time `json:"modified"`
	}
	// TflPlaceIndex to lookup based on name
	TflPlaceIndex struct {
		CommonName string  `json:"commonName"`
		ID         string  `json:"id"`
		URL        string  `json:"url"`
		Lat        float64 `json:"lat"`
		Lon        float64 `json:"lon"`
	}
	// Snapshot of current state
	Snapshot struct {
		TerminalName string    `json:"TerminalName"`
		Installed    bool      `json:"Installed"`
		Locked       bool      `json:"Locked"`
		InstallDate  time.Time `json:"InstallDate"`
		RemovalDate  time.Time `json:"RemovalDate"`
		Temporary    bool      `json:"Temporary"`
		NbBikes      int       `json:"NbBikes"`
		NbEmptyDocks int       `json:"NbEmptyDocks"`
		NbDocks      int       `json:"NbDocks"`
	}
)

func main() {
	act := getAction(os.Args)
	switch act {
	case "init":
		createIndex()
	case "find":
		if len(os.Args) < 3 {
			fmt.Println("missing find param")
			os.Exit(1)
		}
		find(os.Args[2:])
	}
}
func createIndex() {
	f, err := os.Create("./index.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable create index.json %v\n", err)
		os.Exit(1)
	}

	resp, err := http.Get("https://api.tfl.gov.uk/bikepoint")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to download: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)

	var results []TflPlace
	dec.Decode(&results)

	index := make([]TflPlaceIndex, len(results), len(results))
	for i, place := range results {
		index[i] = TflPlaceIndex{
			CommonName: place.CommonName,
			ID:         place.ID,
			URL:        place.URL,
			Lat:        place.Lat,
			Lon:        place.Lon,
		}
	}
	enc := json.NewEncoder(f)
	err = enc.Encode(&index)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to save index: %v\n", err)
		os.Exit(1)
	}
}

func getAction(args []string) string {
	if len(args) < 2 {
		return ""
	}
	return args[1]
}

func find(args []string) {
	f, err := os.Open("./index.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot open index:%v\n", err)
		os.Exit(1)
	}
	var index []TflPlaceIndex
	j := json.NewDecoder(f)
	err = j.Decode(&index)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot decode index:%v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Search %s...\n", strings.Join(args, " "))
	found := []TflPlaceIndex{}
	for _, i := range index {
		c := 0
		for _, a := range args {
			if strings.Contains(strings.ToLower(i.CommonName), strings.ToLower(a)) {
				c++
			}
		}
		if c == len(args) {
			found = append(found, i)
		}
	}

	for _, f := range found {
		fmt.Printf("%s %s\n", f.CommonName, f.URL)
	}
}
