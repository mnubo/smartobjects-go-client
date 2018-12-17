package mnubo

import (
	"encoding/json"
	"fmt"
	"net/url"
)

const (
	eventsPath  = "/api/v3/events"
	objectsPath = "/api/v3/objects"
	ownersPath  = "/api/v3/owners"
)

// Events is a helper to Mnubo client which contains Events related functions.
type Events struct {
	Mnubo Mnubo
}

// Objects is a helper to Mnubo client which contains Objects related functions.
type Objects struct {
	Mnubo Mnubo
}

// Owners is a helper to Mnubo client which contains Owners related functions.
type Owners struct {
	Mnubo Mnubo
}

// Search is a helper to Mnubo client which contains Search related functions.
type Search struct {
	Mnubo Mnubo
}

// SendEventsOptions helps configure the Send events function.
type SendEventsOptions struct {
	ReportResults    bool
	ObjectsMustExist bool
}

// SendEventsReport contains information when sending events with ReportResults set to true.
type SendEventsReport struct {
	ID           string `json:"id"`
	Result       string `json:"result"`
	ObjectExists bool   `json:"objectExists"`
}

// EntitiesExist is an array of map useful when checking if events, objects or owners exist.
type EntitiesExist map[string]bool

// ObjectOwnerPair can be used to claim and unclaim devices.
type ObjectOwnerPair struct {
	XDeviceID string `json:"x_device_id"`
	Username  string `json:"username"`
}

// ClaimResult contains useful information after claiming obects.
type ClaimResult struct {
	ID      string `json:"id"`
	Result  string `json:"result"`
	Message string `json:"message"`
}

// PasswordUpdatePayload helps building the update password function.
type PasswordUpdatePayload struct {
	XPassword string `json:"x_password"`
}

// NewEvents creates an Events wrapper for Mnubo client.
func NewEvents(m Mnubo) *Events {
	return &Events{
		Mnubo: m,
	}
}

// NewObjects creates an Objects wrapper for Mnubo client.
func NewObjects(m Mnubo) *Objects {
	return &Objects{
		Mnubo: m,
	}
}

// NewOwners creates an Owners wrapper for Mnubo client.
func NewOwners(m Mnubo) *Owners {
	return &Owners{
		Mnubo: m,
	}
}

// buildEventsClientRequest is an internal function to help send events to SmartObjects.
func buildEventsClientRequest(events interface{}, options SendEventsOptions, path string) (ClientRequest, error) {
	bytes, err := json.Marshal(events)

	if err != nil {
		return ClientRequest{}, err
	}

	q := url.Values{}

	if options.ObjectsMustExist {
		q.Add("objects_must_exist", "true")
	}

	if options.ReportResults {
		q.Add("report_results", "true")
	}

	return ClientRequest{
		method:      "POST",
		contentType: "application/json",
		path:        path,
		urlQuery:    q,
		payload:     bytes,
	}, nil
}

// Send allows to post events to SmartObjects.
// The events payload depends on the data model.
// See: https://smartobjects.mnubo.com/documentation/api_ingestion.html#post-api-v3-events
func (e *Events) Send(events interface{}, options SendEventsOptions, results interface{}) error {
	cr, err := buildEventsClientRequest(events, options, eventsPath)

	if err != nil {
		return err
	}

	return e.Mnubo.doRequestWithAuthentication(cr, results)
}

// SendFromDevice allows to post events to SmartObjects from one device.
// See: https://smartobjects.mnubo.com/documentation/api_ingestion.html#post-api-v3-objects-x-device-id-events
func (e *Events) SendFromDevice(deviceId string, events interface{}, options SendEventsOptions, results interface{}) error {
	cr, err := buildEventsClientRequest(events, options, fmt.Sprintf("%s/%s/events", objectsPath, deviceId))

	if err != nil {
		return err
	}

	return e.Mnubo.doRequestWithAuthentication(cr, results)
}

// Exists checks if an event has already been submitted.
// See: https://smartobjects.mnubo.com/documentation/api_ingestion.html#post-api-v3-events-exists
func (e *Events) Exists(eventIds []string, results *EntitiesExist) error {
	if *results == nil {
		res := make(EntitiesExist)
		results = &res
	}

	bytes, err := json.Marshal(eventIds)
	if err != nil {
		return err
	}

	cr := ClientRequest{
		method:      "POST",
		contentType: "application/json",
		path:        fmt.Sprintf("%s/exists", eventsPath),
		payload:     bytes,
	}

	rawResults := []map[string]bool{}
	// this endpoint returns an array of objects
	err = e.Mnubo.doRequestWithAuthentication(cr, &rawResults)
	if err != nil {
		return err
	}

	// Flatten the objects so it can be easily used to check for existence
	for _, rr := range rawResults {
		for k, v := range rr {
			(*results)[k] = v
		}
	}

	return nil
}

// Create creates an object to SmartObjects.
// The objects payload is based on the data model.
// See: https://smartobjects.mnubo.com/documentation/api_ingestion.html#post-api-v3-objects
func (o *Objects) Create(objects interface{}, results interface{}) error {
	bytes, err := json.Marshal(objects)

	if err != nil {
		return err
	}

	cr := ClientRequest{
		method:      "POST",
		contentType: "application/json",
		path:        fmt.Sprintf("%s", objectsPath),
		payload:     bytes,
	}

	return o.Mnubo.doRequestWithAuthentication(cr, results)
}

// Update creates and / or updates a batch of objects at once.
// See: https://smartobjects.mnubo.com/documentation/api_ingestion.html#put-api-v3-objects-batch
func (o *Objects) Update(objects interface{}, results interface{}) error {
	bytes, err := json.Marshal(objects)

	if err != nil {
		return err
	}

	cr := ClientRequest{
		method:      "PUT",
		contentType: "application/json",
		path:        fmt.Sprintf("%s", objectsPath),
		payload:     bytes,
	}

	return o.Mnubo.doRequestWithAuthentication(cr, results)
}

// Delete deletes an object
// See: https://smartobjects.mnubo.com/documentation/api_ingestion.html#delete-api-v3-objects-x-device-id
func (o *Objects) Delete(deviceId string) error {
	cr := ClientRequest{
		method: "DELETE",
		path:   fmt.Sprintf("%s/%s", objectsPath, deviceId),
	}

	var results interface{}
	return o.Mnubo.doRequestWithAuthentication(cr, &results)
}

// Exist checks if an array of objects have been created.
// See: https://smartobjects.mnubo.com/documentation/api_ingestion.html#post-api-v3-objects-exists
func (o *Objects) Exist(deviceIds []string, results *EntitiesExist) error {
	if *results == nil {
		res := make(EntitiesExist)
		results = &res
	}

	bytes, err := json.Marshal(deviceIds)
	if err != nil {
		return err
	}

	cr := ClientRequest{
		method:      "POST",
		contentType: "application/json",
		path:        fmt.Sprintf("%s/exists", objectsPath),
		payload:     bytes,
	}

	rawResults := []map[string]bool{}
	// this endpoint returns an array of objects
	err = o.Mnubo.doRequestWithAuthentication(cr, &rawResults)
	if err != nil {
		return err
	}

	// Flatten the objects so it can be easily used to check for existence
	for _, rr := range rawResults {
		for k, v := range rr {
			(*results)[k] = v
		}
	}

	return nil
}

// Create creates a new owner to SmartObjects.
// The owner payload is based on the data model.
// See: https://smartobjects.mnubo.com/documentation/api_ingestion.html#post-api-v3-owners
func (o *Owners) Create(owners interface{}, results interface{}) error {
	bytes, err := json.Marshal(owners)

	if err != nil {
		return err
	}

	cr := ClientRequest{
		method:      "POST",
		contentType: "application/json",
		path:        fmt.Sprintf("%s", ownersPath),
		payload:     bytes,
	}

	return o.Mnubo.doRequestWithAuthentication(cr, results)
}

// Update creates and / or updates a batch of owners at once.
// See: https://smartobjects.mnubo.com/documentation/api_ingestion.html#put-api-v3-owners-batch
func (o *Owners) Update(owners interface{}, results interface{}) error {
	bytes, err := json.Marshal(owners)

	if err != nil {
		return err
	}

	cr := ClientRequest{
		method:      "PUT",
		contentType: "application/json",
		path:        fmt.Sprintf("%s", ownersPath),
		payload:     bytes,
	}

	return o.Mnubo.doRequestWithAuthentication(cr, results)
}

// UpdateOwnerPassword updates an owner password.
// See: https://smartobjects.mnubo.com/documentation/api_ingestion.html#put-api-v3-owners-username-password
func (o *Owners) UpdateOwnerPassword(username string, password string) error {
	bytes, err := json.Marshal(PasswordUpdatePayload{
		XPassword: password,
	})

	if err != nil {
		return err
	}
	cr := ClientRequest{
		method:      "PUT",
		contentType: "application/json",
		path:        fmt.Sprintf("%s/%s/password", ownersPath, username),
		payload:     bytes,
	}

	var results interface{}
	return o.Mnubo.doRequestWithAuthentication(cr, &results)
}

// Delete deletes an owner from SmartObjects.
// See: https://smartobjects.mnubo.com/documentation/api_ingestion.html#delete-api-v3-owners-username
func (o *Owners) Delete(username string) error {
	cr := ClientRequest{
		method: "DELETE",
		path:   fmt.Sprintf("%s/%s", ownersPath, username),
	}

	var results interface{}
	return o.Mnubo.doRequestWithAuthentication(cr, &results)
}

// Exist checks if an array of owners exist in SmartObjects.
// See: https://smartobjects.mnubo.com/documentation/api_ingestion.html#get-api-v3-owners-exists-username
func (o *Owners) Exist(usernames []string, results *EntitiesExist) error {
	// Check if results was nil
	// Covers cases where the user create the results object with something like
	// `var results EntitiesExist`

	if *results == nil {
		res := make(EntitiesExist)
		results = &res
	}

	bytes, err := json.Marshal(usernames)
	if err != nil {
		return err
	}

	cr := ClientRequest{
		method:      "POST",
		contentType: "application/json",
		path:        fmt.Sprintf("%s/exists", ownersPath),
		payload:     bytes,
	}

	rawResults := []map[string]bool{}
	// this endpoint returns an array of objects
	err = o.Mnubo.doRequestWithAuthentication(cr, &rawResults)
	if err != nil {
		return err
	}

	// Flatten the objects so it can be easily used to check for existence
	for _, rr := range rawResults {
		for k, v := range rr {
			(*results)[k] = v
		}
	}

	return nil
}

// Claim claims an array of object / owner pair.
// See: https://smartobjects.mnubo.com/documentation/api_ingestion.html#post-api-v3-owners-claim-batch
func (o *Owners) Claim(pairs []ObjectOwnerPair, results *[]ClaimResult) error {
	bytes, err := json.Marshal(pairs)

	if err != nil {
		return err
	}

	cr := ClientRequest{
		method:      "POST",
		contentType: "application/json",
		path:        fmt.Sprintf("%s/claim", ownersPath),
		payload:     bytes,
	}

	return o.Mnubo.doRequestWithAuthentication(cr, results)
}

// Unclaim unclaims an array of object / owner pair.
// See: https://smartobjects.mnubo.com/documentation/api_ingestion.html#post-api-v3-owners-unclaim-batch
func (o *Owners) Unclaim(pairs []ObjectOwnerPair, results *[]ClaimResult) error {
	bytes, err := json.Marshal(pairs)

	if err != nil {
		return err
	}

	cr := ClientRequest{
		method:      "POST",
		contentType: "application/json",
		path:        fmt.Sprintf("%s/unclaim", ownersPath),
		payload:     bytes,
	}

	return o.Mnubo.doRequestWithAuthentication(cr, results)
}
