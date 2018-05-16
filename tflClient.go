package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type (
	client struct {
		baseURL string
		client  http.Client
	}
)

func newBikeClient() client {
	return client{
		baseURL: "https://api.tfl.gov.uk",
		client: http.Client{
			Timeout: 3 * time.Second,
		},
	}
}

func (c client) ViewIndex() ([]TflPlace, error) {
	var results []TflPlace

	r, err := http.NewRequest("GET", fmt.Sprintf("%s/bikepoint", c.baseURL), nil)
	if err != nil {
		return results, errors.Wrap(err, "cannot create request")
	}
	resp, err := c.client.Do(r)
	if err != nil {
		return results, errors.Wrap(err, "Unable to download")
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)

	dec.Decode(&results)
	return results, nil
}

func (c client) ViewDocking(id string) (TflPlace, error) {
	var results TflPlace

	r, err := http.NewRequest("GET", fmt.Sprintf("%s/Place/BikePoints_%s", c.baseURL, id), nil)
	if err != nil {
		return results, errors.Wrap(err, "cannot create request")
	}
	resp, err := c.client.Do(r)
	if err != nil {
		return results, errors.Wrap(err, "Unable to download")
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return results, errors.Errorf("Id:%q not found", id)
	}

	dec := json.NewDecoder(resp.Body)

	dec.Decode(&results)
	return results, nil
}
