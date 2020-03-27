package mnubo

import (
	"encoding/json"
	"fmt"
)

const (
	modelPath = "/api/v3/model"
)

// Model is a helper to Mnubo client which contains Model related functions.
type Model struct {
	Mnubo *Mnubo
}

type AttributeType struct {
	HighLevelType string `json:"highLevelType"`
	ContainerType string `json:"containerType"`
}

type EventType struct {
	Key            string       `json:"key"`
	DisplayName    string       `json:"displayName"`
	Description    string       `json:"description"`
	Origin         string       `json:"origin"`
	TimeseriesKeys []string     `json:"timeseriesKeys,omitempty"`
	Timeseries     []Timeseries `json:"timeseries,omitempty"`
}

type ObjectAttribute struct {
	Key            string        `json:"key"`
	DisplayName    string        `json:"displayName"`
	Description    string        `json:"description"`
	Type           AttributeType `json:"type,omitempty"`
	ObjectTypeKeys []string      `json:"objectTypeKeys,omitempty"`
}

type ObjectType struct {
	Key                  string            `json:"key"`
	DisplayName          string            `json:"displayName"`
	Description          string            `json:"description"`
	ObjectAttributesKeys []string          `json:"objectAttributesKeys,omitempty"`
	ObjectAttributes     []ObjectAttribute `json:"objectAttributes,omitempty"`
}

type OwnerAttribute struct {
	Key         string        `json:"key"`
	DisplayName string        `json:"displayName"`
	Description string        `json:"description"`
	Type        AttributeType `json:"type"`
}

type ChallengeCode struct {
	Code string `json:"code"`
}

type TimeseriesType struct {
	HighLevelType string `json:"highLevelType"`
}

type Timeseries struct {
	Key           string         `json:"key,omitempty"`
	DisplayName   string         `json:"displayName"`
	Description   string         `json:"description"`
	EventTypeKeys []string       `json:"eventTypeKeys,omitempty"`
	Type          TimeseriesType `json:"type,omitempty"`
}

type Sessionizer struct {
	Key               string `json:"key"`
	DisplayName       string `json:"displayName"`
	Description       string `json:"description"`
	StartEventTypeKey string `json:"startEventTypeKey"`
	EndEventTypeKey   string `json:"endEventTypeKey"`
}

type DataModel struct {
	ObjectTypes     []ObjectType     `json:"objectTypes"`
	EventTypes      []EventType      `json:"eventTypes"`
	OwnerAttributes []OwnerAttribute `json:"ownerAttributes"`
	Sessionizers    []Sessionizer    `json:"sessionizers"`
	Orphans         struct {
		Timeseries []Timeseries `json:"timeseries"`
	} `json:"orphans,omitempty"`
	Enrichers               map[string][]string `json:"enrichers,omitempty"`
	ReservedEnrichersFields []string            `json:"reservedEnrichersFields,omitempty"`
}

func NewModel(m *Mnubo) *Model {
	return &Model{
		Mnubo: m,
	}
}

// Export dumps a JSON object representing the current data model.
// See: https://smartobjects.mnubo.com/documentation/api_modeler.html#exporting-your-data-model
func (m *Model) Export(results *DataModel) error {
	cr := ClientRequest{
		method:      "GET",
		contentType: "application/json",
		path:        fmt.Sprintf("%s/export", modelPath),
	}

	return m.Mnubo.doRequestWithAuthentication(cr, results)
}

// GetTimeseries retrieves the timeseries of the data model.
// See: https://smartobjects.mnubo.com/documentation/api_modeler.html#getting-all-timeseries
func (m *Model) GetTimeseries(results *[]Timeseries) error {
	cr := ClientRequest{
		method:      "GET",
		contentType: "application/json",
		path:        fmt.Sprintf("%s/timeseries", modelPath),
	}

	return m.Mnubo.doRequestWithAuthentication(cr, results)
}

// CreateObjectAttributes creates new object attribute.
// See: https://smartobjects.mnubo.com/documentation/api_modeler.html#creating-object-attributes
func (m *Model) CreateObjectAttributes(oa []ObjectAttribute) error {
	bytes, err := json.Marshal(oa)

	if err != nil {
		return err
	}
	cr := ClientRequest{
		method:      "POST",
		contentType: "application/json",
		path:        fmt.Sprintf("%s/objectAttributes", modelPath),
		payload:     bytes,
	}

	var results interface{}
	return m.Mnubo.doRequestWithAuthentication(cr, &results)
}

// UpdateObjectAttribute updates the DisplayName and Description (only those) from an object attribute.
// See: https://smartobjects.mnubo.com/documentation/api_modeler.html#updating-an-object-attribute
func (m *Model) UpdateObjectAttribute(key string, oa ObjectAttribute) error {
	bytes, err := json.Marshal(oa)

	if err != nil {
		return err
	}
	cr := ClientRequest{
		method:      "PUT",
		contentType: "application/json",
		path:        fmt.Sprintf("%s/objectAttributes/%s", modelPath, key),
		payload:     bytes,
	}

	var results interface{}
	return m.Mnubo.doRequestWithAuthentication(cr, &results)
}

// GenerateObjectAttributeDeployCode generates a new challenge code before deploying an object attribute to production.
// See: https://smartobjects.mnubo.com/documentation/api_modeler.html#id10
func (m *Model) GenerateObjectAttributeDeployCode(key string, results *ChallengeCode) error {
	cr := ClientRequest{
		method:      "POST",
		contentType: "application/json",
		path:        fmt.Sprintf("%s/objectAttributes/%s/deploy", modelPath, key),
	}

	return m.Mnubo.doRequestWithAuthentication(cr, results)
}

// ApplyObjectAttributeDeployCode applies the challenge code to deploy the attribute to production.
// See: https://smartobjects.mnubo.com/documentation/api_modeler.html#id11
func (m *Model) ApplyObjectAttributeDeployCode(key string, cc ChallengeCode) error {
	cr := ClientRequest{
		method:      "POST",
		contentType: "application/json",
		path:        fmt.Sprintf("%s/objectAttributes/%s/deploy/%s", modelPath, key, cc.Code),
	}

	var results interface{}
	return m.Mnubo.doRequestWithAuthentication(cr, &results)
}

// DeployObjectAttributeToProduction deploys an object attribute created in sandbox to production.
// Making calls to GenerateObjectAttributeDeployCode and ApplyObjectAttributeDeployCode.
// See: https://smartobjects.mnubo.com/documentation/api_modeler.html#deploying-an-object-attribute-into-production
func (m *Model) DeployObjectAttributeToProduction(key string) error {
	var cc ChallengeCode
	err := m.GenerateObjectAttributeDeployCode(key, &cc)
	if err != nil {
		return err
	}
	return m.ApplyObjectAttributeDeployCode(key, cc)
}

// GetObjectAttributes retrieves the object attributes of the data model.
// See: https://smartobjects.mnubo.com/documentation/api_modeler.html#getting-all-object-attributes
func (m *Model) GetObjectAttributes(results *[]ObjectAttribute) error {
	cr := ClientRequest{
		method:      "GET",
		contentType: "application/json",
		path:        fmt.Sprintf("%s/objectAttributes", modelPath),
	}

	return m.Mnubo.doRequestWithAuthentication(cr, results)
}

// CreateTimeseries creates new timeseries to the data model.
// See: https://smartobjects.mnubo.com/documentation/api_modeler.html#creating-timeseries
func (m *Model) CreateTimeseries(ts []Timeseries) error {
	bytes, err := json.Marshal(ts)

	if err != nil {
		return err
	}
	cr := ClientRequest{
		method:      "POST",
		contentType: "application/json",
		path:        fmt.Sprintf("%s/timeseries", modelPath),
		payload:     bytes,
	}

	var results interface{}
	return m.Mnubo.doRequestWithAuthentication(cr, &results)
}

// UpdateTimeseries updates the DisplayName and Description (only those) from a Timeseries.
// See: https://smartobjects.mnubo.com/documentation/api_modeler.html#updating-timeseries
func (m *Model) UpdateTimeseries(key string, ts Timeseries) error {
	bytes, err := json.Marshal(ts)

	if err != nil {
		return err
	}
	cr := ClientRequest{
		method:      "PUT",
		contentType: "application/json",
		path:        fmt.Sprintf("%s/timeseries/%s", modelPath, key),
		payload:     bytes,
	}

	var results interface{}
	return m.Mnubo.doRequestWithAuthentication(cr, &results)
}

// GenerateTimeseriesDeployCode generates a new challenge code before deploying a timeseries to production.
// See: https://smartobjects.mnubo.com/documentation/api_modeler.html#part-1-getting-code
func (m *Model) GenerateTimeseriesDeployCode(key string, results *ChallengeCode) error {
	cr := ClientRequest{
		method:      "POST",
		contentType: "application/json",
		path:        fmt.Sprintf("%s/timeseries/%s/deploy", modelPath, key),
	}

	return m.Mnubo.doRequestWithAuthentication(cr, results)
}

// ApplyTimeseriesDeployCode applies the challenge code to deploy the timeseries to production.
// See: https://smartobjects.mnubo.com/documentation/api_modeler.html#part-2-challenging-code
func (m *Model) ApplyTimeseriesDeployCode(key string, cc ChallengeCode) error {
	cr := ClientRequest{
		method:      "POST",
		contentType: "application/json",
		path:        fmt.Sprintf("%s/timeseries/%s/deploy/%s", modelPath, key, cc.Code),
	}

	var results interface{}
	return m.Mnubo.doRequestWithAuthentication(cr, &results)
}

// DeployTimeseriesToProduction deploys a timeseries created in sandbox to production.
// Making calls to GenerateTimeseriesDeployCode and ApplyTimeseriesDeployCode.
// See: https://smartobjects.mnubo.com/documentation/api_modeler.html#deploying-a-timeseries-into-production
func (m *Model) DeployTimeseriesToProduction(key string) error {
	var cc ChallengeCode
	err := m.GenerateTimeseriesDeployCode(key, &cc)
	if err != nil {
		return err
	}
	return m.ApplyTimeseriesDeployCode(key, cc)
}

// CreateOwnerAttributes creates new owner attribute.
// See: https://smartobjects.mnubo.com/documentation/api_modeler.html#creating-object-attributes
func (m *Model) CreateOwnerAttributes(oa []OwnerAttribute) error {
	bytes, err := json.Marshal(oa)

	if err != nil {
		return err
	}
	cr := ClientRequest{
		method:      "POST",
		contentType: "application/json",
		path:        fmt.Sprintf("%s/ownerAttributes", modelPath),
		payload:     bytes,
	}

	var results interface{}
	return m.Mnubo.doRequestWithAuthentication(cr, &results)
}

// UpdateOwnerAttribute updates the DisplayName and Description (only those) from an owner attribute.
// See: https://smartobjects.mnubo.com/documentation/api_modeler.html#updating-an-object-attribute
func (m *Model) UpdateOwnerAttribute(key string, oa OwnerAttribute) error {
	bytes, err := json.Marshal(oa)

	if err != nil {
		return err
	}
	cr := ClientRequest{
		method:      "PUT",
		contentType: "application/json",
		path:        fmt.Sprintf("%s/ownerAttributes/%s", modelPath, key),
		payload:     bytes,
	}

	var results interface{}
	return m.Mnubo.doRequestWithAuthentication(cr, &results)
}

// GenerateOwnerAttributeDeployCode generates a new challenge code before deploying an owner attribute to production.
// See: https://smartobjects.mnubo.com/documentation/api_modeler.html#id10
func (m *Model) GenerateOwnerAttributeDeployCode(key string, results *ChallengeCode) error {
	cr := ClientRequest{
		method:      "POST",
		contentType: "application/json",
		path:        fmt.Sprintf("%s/ownerAttributes/%s/deploy", modelPath, key),
	}

	return m.Mnubo.doRequestWithAuthentication(cr, results)
}

// ApplyOwnerAttributeDeployCode applies the challenge code to deploy the attribute to production.
// See: https://smartobjects.mnubo.com/documentation/api_modeler.html#id11
func (m *Model) ApplyOwnerAttributeDeployCode(key string, cc ChallengeCode) error {
	cr := ClientRequest{
		method:      "POST",
		contentType: "application/json",
		path:        fmt.Sprintf("%s/ownerAttributes/%s/deploy/%s", modelPath, key, cc.Code),
	}

	var results interface{}
	return m.Mnubo.doRequestWithAuthentication(cr, &results)
}

// DeployOwnerAttributeToProduction deploys an owner attribute created in sandbox to production.
// Making calls to GenerateOwnerAttributeDeployCode and ApplyOwnerAttributeDeployCode.
// See: https://smartobjects.mnubo.com/documentation/api_modeler.html#deploying-an-object-attribute-into-production
func (m *Model) DeployOwnerAttributeToProduction(key string) error {
	var cc ChallengeCode
	err := m.GenerateOwnerAttributeDeployCode(key, &cc)
	if err != nil {
		return err
	}
	return m.ApplyOwnerAttributeDeployCode(key, cc)
}

// GetOwnerAttributes retrieves the owner attributes of the data model.
// See: https://smartobjects.mnubo.com/documentation/api_modeler.html#getting-all-owner-attributes
func (m *Model) GetOwnerAttributes(results *[]OwnerAttribute) error {
	cr := ClientRequest{
		method:      "GET",
		contentType: "application/json",
		path:        fmt.Sprintf("%s/ownerAttributes", modelPath),
	}

	return m.Mnubo.doRequestWithAuthentication(cr, results)
}

// GetEventTypes retrieves the event types of the data model.
// See: https://smartobjects.mnubo.com/documentation/api_modeler.html#getting-all-event-types
func (m *Model) GetEventTypes(results *[]EventType) error {
	cr := ClientRequest{
		method:      "GET",
		contentType: "application/json",
		path:        fmt.Sprintf("%s/eventTypes", modelPath),
	}

	return m.Mnubo.doRequestWithAuthentication(cr, results)
}

// CreateEventTypes creates an array of event types in the data model.
// See: https://smartobjects.mnubo.com/documentation/api_modeler.html#creating-event-types
func (m *Model) CreateEventTypes(et []EventType) error {
	bytes, err := json.Marshal(et)

	if err != nil {
		return err
	}
	cr := ClientRequest{
		method:      "POST",
		contentType: "application/json",
		path:        fmt.Sprintf("%s/eventTypes", modelPath),
		payload:     bytes,
	}

	var results interface{}
	return m.Mnubo.doRequestWithAuthentication(cr, &results)
}

// UpdateEventType updates an event type in the data model.
// See: https://smartobjects.mnubo.com/documentation/api_modeler.html#updating-an-event-type
func (m *Model) UpdateEventType(key string, et EventType) error {
	bytes, err := json.Marshal(et)

	if err != nil {
		return err
	}
	cr := ClientRequest{
		method:      "PUT",
		contentType: "application/json",
		path:        fmt.Sprintf("%s/eventTypes/%s", modelPath, key),
		payload:     bytes,
	}

	var results interface{}
	return m.Mnubo.doRequestWithAuthentication(cr, &results)
}

// DeleteEventType deletes an event type from the data model.
// See: https://smartobjects.mnubo.com/documentation/api_modeler.html#deleting-an-event-type
func (m *Model) DeleteEventType(key string) error {
	cr := ClientRequest{
		method:      "DELETE",
		contentType: "application/json",
		path:        fmt.Sprintf("%s/eventTypes/%s", modelPath, key),
	}

	var results interface{}
	return m.Mnubo.doRequestWithAuthentication(cr, &results)
}

// Add a relation to a timeseries.
// See: https://smartobjects.mnubo.com/documentation/api_modeler.html#linking-a-timeseries-to-an-event-type
func (m *Model) AddEventTypeRelation(typeKey string, entityKey string) error {
	cr := ClientRequest{
		method: "POST",
		path:   fmt.Sprintf("%s/eventTypes/%s/timeseries/%s", modelPath, typeKey, entityKey),
	}

	var results interface{}
	return m.Mnubo.doRequestWithAuthentication(cr, &results)
}

// Remove a relation to a timeseries.
// See: https://smartobjects.mnubo.com/documentation/api_modeler.html#deleting-the-link-between-a-timeseries-and-an-event-type
func (m *Model) RemoveEventTypeRelation(typeKey string, entityKey string) error {
	cr := ClientRequest{
		method: "DELETE",
		path:   fmt.Sprintf("%s/eventTypes/%s/timeseries/%s", modelPath, typeKey, entityKey),
	}

	var results interface{}
	return m.Mnubo.doRequestWithAuthentication(cr, &results)
}

// GetObjectTypes retrieves the object types of the data model.
// See: https://smartobjects.mnubo.com/documentation/api_modeler.html#getting-all-event-types
func (m *Model) GetObjectTypes(results *[]ObjectType) error {
	cr := ClientRequest{
		method:      "GET",
		contentType: "application/json",
		path:        fmt.Sprintf("%s/objectTypes", modelPath),
	}

	return m.Mnubo.doRequestWithAuthentication(cr, results)
}

// CreateObjectTypes creates an array of object types in the data model.
// See: https://smartobjects.mnubo.com/documentation/api_modeler.html#creating-object-types
func (m *Model) CreateObjectTypes(ot []ObjectType) error {
	bytes, err := json.Marshal(ot)

	if err != nil {
		return err
	}
	cr := ClientRequest{
		method:      "POST",
		contentType: "application/json",
		path:        fmt.Sprintf("%s/objectTypes", modelPath),
		payload:     bytes,
	}

	var results interface{}
	return m.Mnubo.doRequestWithAuthentication(cr, &results)
}

// UpdateObjectType updates an object type in the data model.
// See: https://smartobjects.mnubo.com/documentation/api_modeler.html#updating-an-object-type
func (m *Model) UpdateObjectType(key string, ot ObjectType) error {
	bytes, err := json.Marshal(ot)

	if err != nil {
		return err
	}
	cr := ClientRequest{
		method:      "PUT",
		contentType: "application/json",
		path:        fmt.Sprintf("%s/objectTypes/%s", modelPath, key),
		payload:     bytes,
	}

	var results interface{}
	return m.Mnubo.doRequestWithAuthentication(cr, &results)
}

// DeleteObjectType deletes an object type from the data model.
// See: https://smartobjects.mnubo.com/documentation/api_modeler.html#deleting-an-object-type
func (m *Model) DeleteObjectType(key string) error {
	cr := ClientRequest{
		method:      "DELETE",
		contentType: "application/json",
		path:        fmt.Sprintf("%s/objectTypes/%s", modelPath, key),
	}

	var results interface{}
	return m.Mnubo.doRequestWithAuthentication(cr, &results)
}

// Add a relation to an object attribute.
// See: https://smartobjects.mnubo.com/documentation/api_modeler.html#linking-an-attribute-to-an-object-type
func (m *Model) AddObjectTypeRelation(typeKey string, entityKey string) error {
	cr := ClientRequest{
		method: "POST",
		path:   fmt.Sprintf("%s/objectTypes/%s/objectAttributes/%s", modelPath, typeKey, entityKey),
	}

	var results interface{}
	return m.Mnubo.doRequestWithAuthentication(cr, &results)
}

// Remove a relation to an object attribute.
// See: https://smartobjects.mnubo.com/documentation/api_modeler.html#deleting-the-link-between-an-attribute-and-an-object-type
func (m *Model) RemoveObjectTypeRelation(typeKey string, entityKey string) error {
	cr := ClientRequest{
		method: "DELETE",
		path:   fmt.Sprintf("%s/objectTypes/%s/objectAttributes/%s", modelPath, typeKey, entityKey),
	}

	var results interface{}
	return m.Mnubo.doRequestWithAuthentication(cr, &results)
}

// GenerateResetCode generates a new code that must be used in order to reset a data model
// in sandbox.
// See: https://smartobjects.mnubo.com/documentation/api_modeler.html#resetting-your-sandbox-data-model
func (m *Model) GenerateResetCode(results *ChallengeCode) error {
	cr := ClientRequest{
		method:      "POST",
		contentType: "application/json",
		path:        fmt.Sprintf("%s/reset", modelPath),
	}

	return m.Mnubo.doRequestWithAuthentication(cr, results)
}

// ApplyResetCode sends the reset code for sandbox reset (if the code is valid).
// See: https://smartobjects.mnubo.com/documentation/api_modeler.html#part-2-using-the-code
func (m *Model) ApplyResetCode(cc ChallengeCode) error {
	cr := ClientRequest{
		method:      "POST",
		contentType: "application/json",
		path:        fmt.Sprintf("%s/reset/%s", modelPath, cc.Code),
	}

	var results interface{}
	return m.Mnubo.doRequestWithAuthentication(cr, &results)
}

// ResetDataModel performs both GenerateResetCode and ApplyResetCode at the same time
// for convenience.
func (m *Model) ResetDataModel() error {
	var cc ChallengeCode
	err := m.GenerateResetCode(&cc)
	if err != nil {
		return err
	}
	return m.ApplyResetCode(cc)
}
