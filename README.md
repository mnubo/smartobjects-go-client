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
)

type SelectOperation struct {
	Count string `json:"count"`
}

type SimpleQuery struct {
	From   string            `json:"from"`
	Select []SelectOperation `json:"select"`
}

type SimpleOwner struct {
	Username string `json:"username"`
	Age      int    `json:"age"`
}

type SimpleObject struct {
	XDeviceID   string `json:"x_device_id"`
	XObjectType string `json:"x_object_type"`
	Color       string `json:"color"`
}

func main() {
	var m *mnubo.Mnubo
	var res interface{}
	var exist mnubo.EntitiesExist

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
		Request:  true, // will send "Content-Encoding: gzip"
		Response: true, // will send "Accept-Encoding: gzip"
	}
	m.Compression = comp

	// Create, Update, Delete Owners
	ow := "user@example.com"
	so := SimpleOwner{
		Username: ow, // Username is mandatory, the rest depends on the data model
		Age:      18,
	}
	m.Owners.Create(so, &res)
	m.Owners.Update([]SimpleOwner{
		{
			Username: ow,
			Age:      19,
		},
	}, &res)
	m.Owners.Delete(ow)

	// Update Owner Password
	m.Owners.UpdateOwnerPassword(ow, "new-password")

	// Check if owners already exist
	m.Owners.Exist([]string{ow, "does-not-exist@example.com"}, &exist)

	// Create, Update, Delete Objects and check if they exist
	ob := "16EAD3C6-48FB-4F34-BC7E-5C45519E2F40"
	sob := SimpleObject{
		XDeviceID: ob, // x_device_id is mandatory, the rest depends on the data model
		Color:     "red",
	}
	m.Objects.Create(sob, &res)
	m.Objects.Update([]SimpleObject{
		{
			XDeviceID: ob,
			Color:     "blue",
		},
	}, &res)
	m.Objects.Delete(ob)
	m.Objects.Exist([]string{ob}, &exist)

	// Claim / Unclaim objects
	var cr []mnubo.ClaimResult
	oop := []mnubo.ObjectOwnerPair{
		{
			XDeviceID: ob,
			Username:  ow,
		},
	}
	m.Owners.Claim(oop, &cr)
	m.Owners.Unclaim(oop, &cr)

	// Getting Datasets for querying
	var ds []mnubo.Dataset
	m.Search.GetDatasets(&ds)

	// Creating a MQL
	var sr mnubo.SearchResults
	qs := `{ "from": "event", "select": [ { "count": "*" } ] }`
	m.Search.CreateBasicQueryWithString(qs, &sr)

	q := SimpleQuery{
		From: "event",
		Select: []SelectOperation{
			{
				Count: "*",
			},
		},
	}

	// Or if you prefer a typed structure
	m.Search.CreateBasicQuery(q, &sr)

	// Validate Query
	var qv mnubo.QueryValidation
	m.Search.ValidateQuery(q, &qv)
	m.Search.ValidateQueryWithString(qs, &qv)
}
```

## References

[mnubo documentation](https://smartobjects.mnubo.com/documentation/)
