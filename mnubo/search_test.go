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

type SearchResultsColumn struct {
	Label string `json:"label"`
	Type  string `json:"type"`
}

type SearchResults struct {
	Columns []SearchResultsColumn `json:"columns"`
	Rows    [][]interface{}       `json:"rows"`
}

func TestCreateQuery(t *testing.T) {
	m := NewClient(os.Getenv("MNUBO_CLIENT_ID"), os.Getenv("MNUBO_CLIENT_SECRET"), os.Getenv("MNUBO_HOST"))
	m.getAccessToken()

	var results = [2]SearchResults{}
	cases := []struct {
		Error error
	}{
		{
			Error: m.createBasicQuery(SimpleQuery{
				From: "event",
				Select: []SelectOperation{
					{
						Count: "*",
					},
				},
			}, &results[0]),
		},
		{
			Error: m.createBasicQueryWithString(`
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

func TestValidateQuery(t *testing.T) {
	m := NewClient(os.Getenv("MNUBO_CLIENT_ID"), os.Getenv("MNUBO_CLIENT_SECRET"), os.Getenv("MNUBO_HOST"))
	m.getAccessToken()

	var results = [2]QueryValidation{}
	cases := []struct {
		Error error
	}{
		{
			Error: m.validateQuery(SimpleQuery{
				From: "event",
				Select: []SelectOperation{
					{
						Count: "*",
					},
				},
			}, &results[0]),
		},
		{
			Error: m.validateQueryWithString(`
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

func TestDatasets(t *testing.T) {
	m := NewClient(os.Getenv("MNUBO_CLIENT_ID"), os.Getenv("MNUBO_CLIENT_SECRET"), os.Getenv("MNUBO_HOST"))
	m.getAccessToken()

	var results = [1][]Dataset{}
	cases := []struct {
		Error error
	}{
		{
			Error: m.getDatasets(&results[0]),
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
