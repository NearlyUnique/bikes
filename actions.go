package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
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
		CommonName   string     `json:"commonName,omitempty"`
		UpdatedAt    time.Time  `json:"updatedAt,omitempty"`
		TerminalName string     `json:"terminalName,omitempty"`
		Installed    bool       `json:"installed,omitempty"`
		Locked       bool       `json:"locked,omitempty"`
		InstallDate  *time.Time `json:"installDate,omitempty"`
		RemovalDate  *time.Time `json:"removalDate,omitempty"`
		Temporary    bool       `json:"temporary,omitempty"`
		Bikes        int        `json:"bikes,omitempty"`
		EmptyDocks   int        `json:"emptyDocks,omitempty"`
		Docks        int        `json:"docks,omitempty"`
	}
)

func createIndex() error {
	f, err := os.Create("./index.json")
	if err != nil {
		return errors.Wrap(err, "Unable create index.json")
	}

	c := newBikeClient()
	results, err := c.ViewIndex()
	if err != nil {
		return err
	}
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
	return errors.Wrap(err, "Unable to save index2")
}

func loadIndex() ([]TflPlaceIndex, error) {
	f, err := os.Open("./index.json")
	if err != nil {
		return nil, errors.Wrap(err, "cannot open index")
	}
	var index []TflPlaceIndex
	j := json.NewDecoder(f)
	err = j.Decode(&index)
	return index, errors.Wrap(err, "cannot decode index")
}

func find(args []string) error {
	index, err := loadIndex()
	if err != nil {
		return errors.Wrap(err, "Failed to load index")
	}
	fmt.Printf("Search %q...\n", strings.Join(args, " "))
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
	return nil
}

func show(args []string) error {
	client := newBikeClient()
	place, err := client.ViewDocking(args[0])
	if err != nil {
		return errors.Wrap(err, "unable to view docking info")
	}
	snap := Snapshot{
		CommonName: place.CommonName,
		UpdatedAt:  time.Now().UTC(),
	}

	for _, prop := range place.AdditionalProperties {
		switch prop.Key {
		case "NbBikes":
			v, _ := strconv.Atoi(prop.Value)
			snap.Bikes = v
		case "NbDocks":
			v, _ := strconv.Atoi(prop.Value)
			snap.Docks = v
		case "NbEmptyDocks":
			v, _ := strconv.Atoi(prop.Value)
			snap.EmptyDocks = v
		case "TerminalName":
			snap.TerminalName = prop.Value
		}
	}

	enc := json.NewEncoder(os.Stdout)
	enc.Encode(snap)

	return nil
}
