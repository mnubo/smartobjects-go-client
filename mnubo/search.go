package mnubo

import (
	"encoding/json"
	"fmt"
)

const (
	searchPath = "/api/v3/search"
)

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

type QueryValidation struct {
	IsValid          bool     `json:"isValid"`
	ValidationErrors []string `json:"validationErrors"`
}

type SearchResultsColumn struct {
	Label string `json:"label"`
	Type  string `json:"type"`
}

type SearchResults struct {
	Columns []SearchResultsColumn `json:"columns"`
	Rows    [][]interface{}       `json:"rows"`
}

func NewSearch(m Mnubo) *Search {
	return &Search{
		Mnubo: m,
	}
}

func (s *Search) CreateBasicQuery(mql interface{}, results interface{}) error {
	payload, err := json.Marshal(mql)

	if err != nil {
		return fmt.Errorf("unable to marshal the mql: %s (%s)", mql, err)
	}

	return s.CreateBasicQueryWithBytes(payload, results)
}

func (s *Search) CreateBasicQueryWithString(mql string, results interface{}) error {
	return s.CreateBasicQueryWithBytes([]byte(mql), results)
}

func (s *Search) CreateBasicQueryWithBytes(mql []byte, results interface{}) error {
	cr := ClientRequest{
		method:      "POST",
		contentType: "application/json",
		path:        fmt.Sprintf("%s/basic", searchPath),
		payload:     mql,
	}

	return s.Mnubo.doRequestWithAuthentication(cr, results)
}

func (s *Search) ValidateQuery(mql interface{}, results *QueryValidation) error {
	payload, err := json.Marshal(mql)

	if err != nil {
		return fmt.Errorf("unable to marshal the mql: %s (%s)", mql, err)
	}

	return s.ValidateQueryWithBytes(payload, results)
}

func (s *Search) ValidateQueryWithString(mql string, results *QueryValidation) error {
	return s.ValidateQueryWithBytes([]byte(mql), results)
}

func (s *Search) ValidateQueryWithBytes(mql []byte, results *QueryValidation) error {
	cr := ClientRequest{
		method:      "POST",
		contentType: "application/json",
		path:        fmt.Sprintf("%s/validateQuery", searchPath),
		payload:     mql,
	}

	return s.Mnubo.doRequestWithAuthentication(cr, results)
}

func (s *Search) GetDatasets(results *[]Dataset) error {
	cr := ClientRequest{
		method:      "GET",
		contentType: "application/json",
		path:        fmt.Sprintf("%s/datasets", searchPath),
	}

	return s.Mnubo.doRequestWithAuthentication(cr, results)
}
