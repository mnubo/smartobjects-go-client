# SmartObjects Go Client

[![Build status](https://travis-ci.org/mnubo/smartobjects-go-client.svg?branch=master)](https://travis-ci.org/mnubo/smartobjects-go-client)

## Installation

```bash
go get github.com/mnubo/smartobjects-go-client
```

## Usage

```go
package main

import (
	"github.com/mnubo/smartobjects-go-client/mnubo"
	"log"
)

func main() {
	var m *mnubo.Mnubo
	
	// Creating new client with client id and secret.
	// Get them by going to the Security app: https://smartobjects.mnubo.com/apps/security
	m = mnubo.NewClient("YOUR_CLIENT_ID", "YOUR_CLIENT_SECRET", "YOUR_HOST_URL")
	// Creating new client with a static token that you manage yourself
	// Create one by going to the Security app: https://smartobjects.mnubo.com/apps/security
	m = mnubo.NewClientWithToken("YOUR_STATIC_TOKEN", "YOUR_HOST_URL")
	
	// Activate compression (optional).
	// See: https://smartobjects.mnubo.com/documentation/api_basics.html#compression-support
	// Two modes are available:
	// - request Useful when ingesting a lot of data (will compress using `gzip.BestSpeed`)
	// - response Useful when retrieving a lot of data (will ask for response to be gzipped)
	comp := mnubo.CompressionConfig{
		Request: true, // will send "Content-Encoding: gzip"
		Response: true, // will send "Accept-Encoding: gzip"
	}
	m.Compression = comp
	
	// Getting Datasets for querying
	var ds []mnubo.Dataset
	m.Search.GetDatasets(&ds)
	
	// Creating a MQL
    var results mnubo.SearchResults
	qs := `{ "from": "event", "select": [ { "count": "*" } ] }`
	m.Search.CreateBasicQueryWithString(qs, &results)
	
	// Or if you prefer a typed structure
	type SelectOperation struct {
    	Count string `json:"count"`
    }
    type SimpleQuery struct {
    	From   string            `json:"from"`
    	Select []SelectOperation `json:"select"`
    }
	
	q := SimpleQuery {
		From: "event",
		Select: []SelectOperation {
			{
				Count: "*",
			},
		},
	}

	m.Search.CreateBasicQuery(q, &results)
	
	// Validate Query
	var qv mnubo.QueryValidation
	m.Search.ValidateQuery(q, &qv)
	m.Search.ValidateQueryWithString(qs, &qv)
}
```

## References

[mnubo documentation](https://smartobjects.mnubo.com/documentation/)
