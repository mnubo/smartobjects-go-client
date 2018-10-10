package mnubo

import (
	"encoding/json"
	"fmt"
	"net/url"
)

const (
	eventsPath  = "/api/v3/events"
	objectsPath = "/api/v3/objects"
)

type Events struct {
	Mnubo Mnubo
}

type SendEventsOptions struct {
	ReportResults    bool
	ObjectsMustExist bool
}

type SendEventsReport struct {
	ID           string `json:"id"`
	Result       string `json:"result"`
	ObjectExists bool   `json:"objectExists"`
}

type EventsExist []map[string]bool

func NewEvents(m Mnubo) *Events {
	return &Events{
		Mnubo: m,
	}
}

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

func (e *Events) Send(events interface{}, options SendEventsOptions, results interface{}) error {
	cr, err := buildEventsClientRequest(events, options, eventsPath)

	if err != nil {
		return err
	}

	return e.Mnubo.doRequestWithAuthentication(cr, results)
}

func (e *Events) SendFromDevice(deviceId string, events interface{}, options SendEventsOptions, results interface{}) error {
	cr, err := buildEventsClientRequest(events, options, fmt.Sprintf("%s/%s/events", objectsPath, deviceId))

	if err != nil {
		return err
	}

	return e.Mnubo.doRequestWithAuthentication(cr, results)
}

func (e *Events) Exists(eventIds []string, results interface{}) error {
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

	return e.Mnubo.doRequestWithAuthentication(cr, results)
}
