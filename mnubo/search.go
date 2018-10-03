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

func (m *Mnubo) createBasicQuery(mql interface{}, results interface{}) error {
	payload, err := json.Marshal(mql)

	if err != nil {
		return fmt.Errorf("unable to marshal the mql: %s (%s)", mql, err)
	}

	return m.createBasicQueryWithBytes(payload, results)
}

func (m *Mnubo) createBasicQueryWithString(mql string, results interface{}) error {
	return m.createBasicQueryWithBytes([]byte(mql), results)
}

func (m *Mnubo) createBasicQueryWithBytes(mql []byte, results interface{}) error {
	cr := ClientRequest{
		method:      "POST",
		contentType: "application/json",
		path:        fmt.Sprintf("%s/basic", searchPath),
		payload:     mql,
	}

	return m.doRequestWithAuthentication(cr, results)
}

func (m *Mnubo) validateQuery(mql interface{}, results *QueryValidation) error {
	payload, err := json.Marshal(mql)

	if err != nil {
		return fmt.Errorf("unable to marshal the mql: %s (%s)", mql, err)
	}

	return m.validateQueryWithBytes(payload, results)
}

func (m *Mnubo) validateQueryWithString(mql string, results *QueryValidation) error {
	return m.validateQueryWithBytes([]byte(mql), results)
}

func (m *Mnubo) validateQueryWithBytes(mql []byte, results *QueryValidation) error {
	cr := ClientRequest{
		method:      "POST",
		contentType: "application/json",
		path:        fmt.Sprintf("%s/validateQuery", searchPath),
		payload:     mql,
	}

	return m.doRequestWithAuthentication(cr, results)
}

func (m *Mnubo) getDatasets(results *[]Dataset) error {
	cr := ClientRequest{
		method:      "GET",
		contentType: "application/json",
		path:        fmt.Sprintf("%s/datasets", searchPath),
	}

	return m.doRequestWithAuthentication(cr, results)
}
