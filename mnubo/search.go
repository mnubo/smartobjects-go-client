package mnubo

import (
	"encoding/json"
	"fmt"
)

const (
	searchPath = "/api/v3/search"
)

// Dataset is the main structure used to make queries to SmartObjects.
// Its key represent the "From" to use when making queries.
type Dataset struct {
	Key         string      `json:"key"`
	Description interface{} `json:"description"`
	DisplayName string      `json:"displayName"`
	Fields      []struct {
		Key           string `json:"key"`
		HighLevelType string `json:"highLevelType"`
		DisplayName   string `json:"displayName"`
		Description   string `json:"description"`
		ContainerType string `json:"containerType"`
		PrimaryKey    bool   `json:"primaryKey"`
	} `json:"fields"`
}

// QueryValidation is a helper structure that contains useful information about MQL queries validity.
type QueryValidation struct {
	IsValid          bool     `json:"isValid"`
	ValidationErrors []string `json:"validationErrors"`
}

// SearchResultsColumn is the definition of the a SearchResults column.
type SearchResultsColumn struct {
	Label string `json:"label"`
	Type  string `json:"type"`
}

// SearchResults is the main structure that contain results after a valid query has been sent.
type SearchResults struct {
	Columns []SearchResultsColumn `json:"columns"`
	Rows    [][]interface{}       `json:"rows"`
}

// NewSearch creates a Search wrapper for Mnubo client.
func NewSearch(m Mnubo) *Search {
	return &Search{
		Mnubo: m,
	}
}

// CreateBasicQuery is the main function to make analytics queries to SmartObjects.
// Its mql data structure varies depending on the data model.
// This method trust the user will use the proper structure based on its model.
// See: https://smartobjects.mnubo.com/documentation/api_search.html#basic
func (s *Search) CreateBasicQuery(mql interface{}, results interface{}) error {
	payload, err := json.Marshal(mql)

	if err != nil {
		return fmt.Errorf("unable to marshal the mql: %s (%s)", mql, err)
	}

	return s.CreateBasicQueryWithBytes(payload, results)
}

// CreateBasicQueryWithString is a helper to use a string instead of creating a dedicated structure.
func (s *Search) CreateBasicQueryWithString(mql string, results interface{}) error {
	return s.CreateBasicQueryWithBytes([]byte(mql), results)
}

// CreateBasicQueryWithString is a helper to use an array of bytes.
func (s *Search) CreateBasicQueryWithBytes(mql []byte, results interface{}) error {
	cr := ClientRequest{
		method:      "POST",
		contentType: "application/json",
		path:        fmt.Sprintf("%s/basic", searchPath),
		payload:     mql,
	}

	return s.Mnubo.doRequestWithAuthentication(cr, results)
}

// ValidateQuery is the main function to use to understand why an MQL query is not valid.
// See: https://smartobjects.mnubo.com/documentation/api_search.html#validate
func (s *Search) ValidateQuery(mql interface{}, results *QueryValidation) error {
	payload, err := json.Marshal(mql)

	if err != nil {
		return fmt.Errorf("unable to marshal the mql: %s (%s)", mql, err)
	}

	return s.ValidateQueryWithBytes(payload, results)
}

// ValidateQueryWithString is a helper that allows to send a string instead of creating
// a dedicated structure.
func (s *Search) ValidateQueryWithString(mql string, results *QueryValidation) error {
	return s.ValidateQueryWithBytes([]byte(mql), results)
}

// ValidateQueryWithBytes is a helper that allows to send bytes.
func (s *Search) ValidateQueryWithBytes(mql []byte, results *QueryValidation) error {
	cr := ClientRequest{
		method:      "POST",
		contentType: "application/json",
		path:        fmt.Sprintf("%s/validateQuery", searchPath),
		payload:     mql,
	}

	return s.Mnubo.doRequestWithAuthentication(cr, results)
}

// GetDatasets returns an array of SmartObjects datasets to perform queries.
func (s *Search) GetDatasets(results *[]Dataset) error {
	cr := ClientRequest{
		method:      "GET",
		contentType: "application/json",
		path:        fmt.Sprintf("%s/datasets", searchPath),
	}

	return s.Mnubo.doRequestWithAuthentication(cr, results)
}
