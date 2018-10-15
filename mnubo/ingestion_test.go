package mnubo

import (
	"github.com/google/uuid"
	"os"
	"testing"
)

type XObject struct {
	XDeviceID string `json:"x_device_id"`
}

type EventWithObject struct {
	XObject    XObject `json:"x_object"`
	XEventType string  `json:"x_event_type"`
}

type SimpleEvent struct {
	XEventType string `json:"x_event_type"`
}

type SimpleObject struct {
	XDeviceID   string `json:"x_device_id"`
	XObjectType string `json:"x_object_type"`
}

type SimpleOwner struct {
	Username string `json:"username"`
}

func TestEvents_Send(t *testing.T) {
	m := NewClient(os.Getenv("MNUBO_CLIENT_ID"), os.Getenv("MNUBO_CLIENT_SECRET"), os.Getenv("MNUBO_HOST"))

	var results [3][]SendEventsReport

	cases := []struct {
		Error    error
		Expected []SendEventsReport
	}{
		{
			Error: m.Events.Send([]EventWithObject{
				{
					XObject: XObject{
						XDeviceID: uuid.New().String(),
					},
					XEventType: "event_type1",
				},
			}, SendEventsOptions{
				ReportResults:    false,
				ObjectsMustExist: false,
			}, &results[0]),
			Expected: nil,
		},
		{
			Error: m.Events.Send([]EventWithObject{
				{
					XObject: XObject{
						XDeviceID: uuid.New().String(),
					},
					XEventType: "event_type1",
				},
			}, SendEventsOptions{
				ReportResults:    true,
				ObjectsMustExist: false,
			}, &results[1]),
			Expected: []SendEventsReport{
				{
					Result:       "success",
					ObjectExists: false,
				},
			},
		},
		{
			Error: m.Events.Send([]EventWithObject{
				{
					XObject: XObject{
						XDeviceID: uuid.New().String(),
					},
					XEventType: "event_type1",
				},
			}, SendEventsOptions{
				ReportResults:    true,
				ObjectsMustExist: true,
			}, &results[2]),
			Expected: []SendEventsReport{
				{
					Result:       "error",
					ObjectExists: false,
				},
			},
		},
	}

	for i, c := range cases {
		if c.Error != nil {
			t.Errorf("%d, client call failed: %+v", i, c.Error)
		}

		if len(results[i]) != len(c.Expected) {
			t.Errorf("%d, expecting: %d, got: %d", i, len(c.Expected), len(results[i]))
		}

		for j := range results[i] {
			ra := results[i][j].Result
			re := c.Expected[j].Result
			oa := results[i][j].ObjectExists
			oe := c.Expected[j].ObjectExists
			if ra != re {
				t.Errorf("%d, expecting: %+v, got: %+v", i, ra, re)
			}
			if oa != oe {
				t.Errorf("%d, expecting: %+v, got: %+v", i, oa, oe)
			}
		}
	}
}

func TestEvents_SendFromDevice(t *testing.T) {
	m := NewClient(os.Getenv("MNUBO_CLIENT_ID"), os.Getenv("MNUBO_CLIENT_SECRET"), os.Getenv("MNUBO_HOST"))

	var results [1][]SendEventsReport

	cases := []struct {
		Error          error
		ExpectedLength int
	}{
		{
			Error: m.Events.SendFromDevice(uuid.New().String(),
				[]SimpleEvent{
					{
						XEventType: "event_type1",
					},
				},
				SendEventsOptions{
					ReportResults: true,
				},
				&results[0]),
			ExpectedLength: 1,
		},
	}

	for i, c := range cases {
		if c.Error != nil {
			t.Errorf("%d, client called failed: %+v", i, c.Error)
		}

		if len(results[i]) != c.ExpectedLength {
			t.Errorf("%d, expecting length: %d, got %d", i, c.ExpectedLength, len(results[i]))
		}
	}
}

func TestEvents_Exists(t *testing.T) {
	m := NewClient(os.Getenv("MNUBO_CLIENT_ID"), os.Getenv("MNUBO_CLIENT_SECRET"), os.Getenv("MNUBO_HOST"))

	var results [1]EntitiesExist

	cases := []struct {
		Error          error
		ExpectedLength int
	}{
		{
			Error:          m.Events.Exists([]string{uuid.New().String()}, &results[0]),
			ExpectedLength: 1,
		},
	}

	for i, c := range cases {
		if c.Error != nil {
			t.Errorf("%d, client call failed: %+v", i, c.Error)
		}

		if len(results[i]) != c.ExpectedLength {
			t.Errorf("%d, expecting length: %d, got %d", i, c.ExpectedLength, len(results[i]))
		}
	}
}

func TestObjects_Create(t *testing.T) {
	m := NewClient(os.Getenv("MNUBO_CLIENT_ID"), os.Getenv("MNUBO_CLIENT_SECRET"), os.Getenv("MNUBO_HOST"))
	id := uuid.New().String()

	o := SimpleObject{
		XDeviceID:   id,
		XObjectType: "rand",
	}
	var results interface{}
	err := m.Objects.Create(o, &results)

	if err != nil {
		t.Errorf("client call failed: %+v", err)
	}
}

func TestObjects_Update(t *testing.T) {
	m := NewClient(os.Getenv("MNUBO_CLIENT_ID"), os.Getenv("MNUBO_CLIENT_SECRET"), os.Getenv("MNUBO_HOST"))
	id := uuid.New().String()

	o := []SimpleObject{
		{
			XDeviceID:   id,
			XObjectType: "rand",
		},
	}
	var results interface{}
	err := m.Objects.Update(o, &results)

	if err != nil {
		t.Errorf("client call failed: %+v", err)
	}
}

func TestObjects_Delete(t *testing.T) {
	m := NewClient(os.Getenv("MNUBO_CLIENT_ID"), os.Getenv("MNUBO_CLIENT_SECRET"), os.Getenv("MNUBO_HOST"))
	id := uuid.New().String()

	o := SimpleObject{
		XDeviceID:   id,
		XObjectType: "rand",
	}
	var results interface{}
	e1, e2 := m.Objects.Create(o, &results), m.Objects.Delete(id)

	if e1 != nil {
		t.Errorf("client call failed: %+v", e1)
	}

	if e2 != nil {
		t.Errorf("client call failed: %+v", e2)
	}
}

func TestObjects_Exist(t *testing.T) {
	m := NewClient(os.Getenv("MNUBO_CLIENT_ID"), os.Getenv("MNUBO_CLIENT_SECRET"), os.Getenv("MNUBO_HOST"))
	id := uuid.New().String()

	var results EntitiesExist
	err := m.Objects.Exist([]string{id}, &results)

	if err != nil {
		t.Errorf("client call failed: %+v", err)
	}
}

func TestOwners_Create(t *testing.T) {
	m := NewClient(os.Getenv("MNUBO_CLIENT_ID"), os.Getenv("MNUBO_CLIENT_SECRET"), os.Getenv("MNUBO_HOST"))
	id := uuid.New().String()

	o := SimpleOwner{
		Username: id,
	}
	var results interface{}
	err := m.Owners.Create(o, &results)

	if err != nil {
		t.Errorf("client call failed: %+v", err)
	}
}

func TestOwners_Update(t *testing.T) {
	m := NewClient(os.Getenv("MNUBO_CLIENT_ID"), os.Getenv("MNUBO_CLIENT_SECRET"), os.Getenv("MNUBO_HOST"))
	id := uuid.New().String()

	o := []SimpleOwner{
		{
			Username: id,
		},
	}
	var results interface{}
	err := m.Owners.Update(o, &results)

	if err != nil {
		t.Errorf("client call failed: %+v", err)
	}
}

func TestOwners_UpdateOwnerPassword(t *testing.T) {
	m := NewClient(os.Getenv("MNUBO_CLIENT_ID"), os.Getenv("MNUBO_CLIENT_SECRET"), os.Getenv("MNUBO_HOST"))
	id := uuid.New().String()

	o := SimpleOwner{
		Username: id,
	}
	var results interface{}
	e1, e2 := m.Owners.Create(o, &results), m.Owners.UpdateOwnerPassword(id, "test")

	if e1 != nil {
		t.Errorf("client call failed: %+v", e1)
	}

	if e2 != nil {
		t.Errorf("client call failed: %+v", e2)
	}

}

func TestOwners_Delete(t *testing.T) {
	m := NewClient(os.Getenv("MNUBO_CLIENT_ID"), os.Getenv("MNUBO_CLIENT_SECRET"), os.Getenv("MNUBO_HOST"))
	id := uuid.New().String()

	o := SimpleOwner{
		Username: id,
	}
	var results interface{}
	e1, e2 := m.Owners.Create(o, &results), m.Owners.Delete(id)

	if e1 != nil {
		t.Errorf("client call failed: %+v", e1)
	}

	if e2 != nil {
		t.Errorf("client call failed: %+v", e2)
	}
}

func TestOwners_Exist(t *testing.T) {
	m := NewClient(os.Getenv("MNUBO_CLIENT_ID"), os.Getenv("MNUBO_CLIENT_SECRET"), os.Getenv("MNUBO_HOST"))
	id := uuid.New().String()

	var results EntitiesExist
	err := m.Owners.Exist([]string{id}, &results)

	if err != nil {
		t.Errorf("client call failed: %+v", err)
	}
}

func TestOwners_Claim(t *testing.T) {
	m := NewClient(os.Getenv("MNUBO_CLIENT_ID"), os.Getenv("MNUBO_CLIENT_SECRET"), os.Getenv("MNUBO_HOST"))
	ob, ow := uuid.New().String(), uuid.New().String()

	var results []ClaimResult
	err := m.Owners.Claim([]ObjectOwnerPair{
		{
			XDeviceID: ob,
			Username:  ow,
		},
	}, &results)

	if err != nil {
		t.Errorf("client call failed: %+v", err)
	}
}

func TestOwners_Unclaim(t *testing.T) {
	m := NewClient(os.Getenv("MNUBO_CLIENT_ID"), os.Getenv("MNUBO_CLIENT_SECRET"), os.Getenv("MNUBO_HOST"))
	ob, ow := uuid.New().String(), uuid.New().String()

	var results []ClaimResult
	err := m.Owners.Unclaim([]ObjectOwnerPair{
		{
			XDeviceID: ob,
			Username:  ow,
		},
	}, &results)

	if err != nil {
		t.Errorf("client call failed: %+v", err)
	}
}
