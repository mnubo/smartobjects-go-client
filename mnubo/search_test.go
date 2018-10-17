package mnubo

import (
	"os"
	"testing"
)

type SelectOperation struct {
	Count string `json:"count"`
}

type SimpleQuery struct {
	From   string            `json:"from"`
	Select []SelectOperation `json:"select"`
}

func TestSearch_CreateBasicQuery(t *testing.T) {
	m := NewClient(os.Getenv("MNUBO_CLIENT_ID"), os.Getenv("MNUBO_CLIENT_SECRET"), os.Getenv("MNUBO_HOST"))

	var results = [2]SearchResults{}
	cases := []struct {
		Error error
	}{
		{
			Error: m.Search.CreateBasicQuery(SimpleQuery{
				From: "event",
				Select: []SelectOperation{
					{
						Count: "*",
					},
				},
			}, &results[0]),
		},
		{
			Error: m.Search.CreateBasicQueryWithString(`
				{
				    "from": "event",
				    "select": [
				        { "count": "*" }
				    ]
				}
			`, &results[1]),
		},
	}

	for i, c := range cases {
		if c.Error != nil {
			t.Errorf("%d, could not create basic query: %t", i, c.Error)
		}

		if len(results[i].Rows) != 1 || len(results[i].Rows[0]) != 1 {
			t.Errorf("%d, expecting results to have a count in firt row and cell", i)
		}
	}
}

func TestSearch_ValidateQuery(t *testing.T) {
	m := NewClient(os.Getenv("MNUBO_CLIENT_ID"), os.Getenv("MNUBO_CLIENT_SECRET"), os.Getenv("MNUBO_HOST"))

	var results = [2]QueryValidation{}
	cases := []struct {
		Error error
	}{
		{
			Error: m.Search.ValidateQuery(SimpleQuery{
				From: "event",
				Select: []SelectOperation{
					{
						Count: "*",
					},
				},
			}, &results[0]),
		},
		{
			Error: m.Search.ValidateQueryWithString(`
				{
				    "from": "event",
				    "select": [
				        { "count": "*" }
				    ]
				}
			`, &results[1]),
		},
	}

	for i, c := range cases {
		if c.Error != nil {
			t.Errorf("%d, could not validate basic query: %t", i, c.Error)
		}

		if results[i].IsValid != true {
			t.Errorf("%d, expecting the query to be valid", i)
		}
	}
}

func TestSearch_GetDatasets(t *testing.T) {
	m := NewClient(os.Getenv("MNUBO_CLIENT_ID"), os.Getenv("MNUBO_CLIENT_SECRET"), os.Getenv("MNUBO_HOST"))

	var results = [1][]Dataset{}
	cases := []struct {
		Error error
	}{
		{
			Error: m.Search.GetDatasets(&results[0]),
		},
	}

	for i, c := range cases {
		if c.Error != nil {
			t.Errorf("%d, could not create basic query: %t", i, c.Error)
		}

		if len(results[i]) <= 3 {
			t.Errorf("%d, expecting datasets to be equal or greater than 3", i)
		}
	}
}
