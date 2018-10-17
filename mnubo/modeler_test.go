package mnubo

import (
	"github.com/google/uuid"
	"os"
	"testing"
)

var m = NewClient(os.Getenv("MNUBO_CLIENT_ID"), os.Getenv("MNUBO_CLIENT_SECRET"), os.Getenv("MNUBO_HOST"))

func TestModel_Export(t *testing.T) {
	var dm DataModel
	err := m.Model.Export(&dm)

	if err != nil {
		t.Errorf("client call failed: %+v", err)
	}
}

func TestModel_GetTimeseries(t *testing.T) {
	var ts []Timeseries
	err := m.Model.GetTimeseries(&ts)

	if err != nil {
		t.Errorf("client call failed: %+v", err)
	}
}

func TestModel_CreateTimeseries(t *testing.T) {
	ts := []Timeseries{
		{
			Key:           uuid.New().String(),
			EventTypeKeys: []string{"event_type1"},
			Type: TimeseriesType{
				HighLevelType: "DOUBLE",
			},
		},
	}
	err := m.Model.CreateTimeseries(ts)

	if err != nil {
		t.Errorf("client call failed: %+v", err)
	}
}

func TestModel_UpdateTimeseries(t *testing.T) {
	key := uuid.New().String()
	ts1 := []Timeseries{
		{
			Key:           key,
			Description:   "Description 1",
			DisplayName:   "Display Name 1",
			EventTypeKeys: []string{"event_type1"},
			Type: TimeseriesType{
				HighLevelType: "DOUBLE",
			},
		},
	}

	ts2 := Timeseries{
		DisplayName: "Display Name 2",
		Description: "Description 2",
	}
	e1, e2 := m.Model.CreateTimeseries(ts1), m.Model.UpdateTimeseries(key, ts2)

	if e1 != nil {
		t.Errorf("client call failed: %+v", e1)
	}

	if e2 != nil {
		t.Errorf("client call failed: %+v", e2)
	}
}

func TestModel_DeployTimeseriesToProduction(t *testing.T) {
	key := uuid.New().String()
	ts1 := []Timeseries{
		{
			Key:           key,
			Description:   "Description 1",
			DisplayName:   "Display Name 1",
			EventTypeKeys: []string{"event_type1"},
			Type: TimeseriesType{
				HighLevelType: "DOUBLE",
			},
		},
	}

	e1, e2 := m.Model.CreateTimeseries(ts1), m.Model.DeployTimeseriesToProduction(key)

	if e1 != nil {
		t.Errorf("client call failed: %+v", e1)
	}

	if e2 != nil {
		t.Errorf("client call failed: %+v", e2)
	}
}

func TestModel_ObjectAttributes(t *testing.T) {
	key := uuid.New().String()
	var roa []ObjectAttribute
	oa := []ObjectAttribute{
		{
			Key:         key,
			DisplayName: "test",
			Description: "test",
			Type: AttributeType{
				HighLevelType: "DOUBLE",
				ContainerType: "none",
			},
			ObjectTypeKeys: []string{"cat_detector"},
		},
	}

	cases := []struct {
		Error error
	}{
		{
			Error: m.Model.CreateObjectAttributes(oa),
		},
		{
			Error: m.Model.UpdateObjectAttribute(key, ObjectAttribute{
				DisplayName: "test 2",
				Description: "test 2",
			}),
		},
		{
			Error: m.Model.DeployObjectAttributeToProduction(key),
		},
		{
			Error: m.Model.GetObjectAttributes(&roa),
		},
	}

	for i, c := range cases {
		if c.Error != nil {
			t.Errorf("%d, client call failed: %+v", i, c.Error)
		}
	}
}

func TestModel_OwnerAttributes(t *testing.T) {
	key := uuid.New().String()
	var roa []OwnerAttribute
	oa := []OwnerAttribute{
		{
			Key:         key,
			DisplayName: "test",
			Description: "test",
			Type: AttributeType{
				HighLevelType: "DOUBLE",
				ContainerType: "none",
			},
		},
	}

	cases := []struct {
		Error error
	}{
		{
			Error: m.Model.CreateOwnerAttributes(oa),
		},
		{
			Error: m.Model.UpdateOwnerAttribute(key, OwnerAttribute{
				DisplayName: "test 2",
				Description: "test 2",
			}),
		},
		{
			Error: m.Model.DeployOwnerAttributeToProduction(key),
		},
		{
			Error: m.Model.GetOwnerAttributes(&roa),
		},
	}

	for i, c := range cases {
		if c.Error != nil {
			t.Errorf("%d, client call failed: %+v", i, c.Error)
		}
	}

}

func TestModel_EventTypes(t *testing.T) {
	key := uuid.New().String()
	var ret []EventType
	et := []EventType{
		{
			Key:         key,
			DisplayName: "test",
			Description: "test",
			Origin:      "rule",
		},
	}

	cases := []struct {
		Error error
	}{
		{
			Error: m.Model.CreateEventTypes(et),
		},
		{
			Error: m.Model.UpdateEventType(key, EventType{
				Key:         key,
				DisplayName: "test 2",
				Description: "test 2",
				Origin:      "scheduled",
			}),
		},
		{
			Error: m.Model.DeleteEventType(key),
		},
		{
			Error: m.Model.GetEventTypes(&ret),
		},
	}

	for i, c := range cases {
		if c.Error != nil {
			t.Errorf("%d, client call failed: %+v", i, c.Error)
		}
	}
}

func TestModel_ObjectTypes(t *testing.T) {
	key := uuid.New().String()
	ot := []ObjectType{
		{
			Key:         key,
			DisplayName: "Object Type",
			Description: "Description",
		},
	}
	var rot []ObjectType

	cases := []struct {
		Error error
	}{
		{
			Error: m.Model.CreateObjectTypes(ot),
		},
		{
			Error: m.Model.UpdateObjectType(key, ObjectType{
				Key:         key,
				DisplayName: "Object Type 2",
				Description: "Description 2",
			}),
		},
		{
			Error: m.Model.DeleteObjectType(key),
		},
		{
			Error: m.Model.GetObjectTypes(&rot),
		},
	}

	for i, c := range cases {
		if c.Error != nil {
			t.Errorf("%d, client call failed: %+v", i, c.Error)
		}
	}
}

func TestModel_GenerateResetCode(t *testing.T) {
	var cc ChallengeCode
	err := m.Model.GenerateResetCode(&cc)

	if err != nil {
		t.Errorf("client call failed: %+v", err)
	}
}

//func TestModel_ResetDataModel(t *testing.T) {
//	err := m.Model.ResetDataModel()
//
//	if err != nil {
//		t.Errorf("client call failed: %+v", err)
//	}
//}
