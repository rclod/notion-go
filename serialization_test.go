package notionapi

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func TestStatusPropertyConfig_OmitsEmptyFields(t *testing.T) {
	t.Run("empty ID is omitted", func(t *testing.T) {
		p := StatusPropertyConfig{
			Type: PropertyConfigStatus,
			Status: StatusConfig{
				Options: []Option{{Name: "Done", Color: "green"}},
			},
		}
		b, err := json.Marshal(p)
		if err != nil {
			t.Fatal(err)
		}
		s := string(b)
		if strings.Contains(s, `"id"`) {
			t.Errorf("expected empty ID to be omitted, got %s", s)
		}
	})

	t.Run("non-empty ID is included", func(t *testing.T) {
		p := StatusPropertyConfig{
			ID:   "abc",
			Type: PropertyConfigStatus,
		}
		b, err := json.Marshal(p)
		if err != nil {
			t.Fatal(err)
		}
		s := string(b)
		if !strings.Contains(s, `"id":"abc"`) {
			t.Errorf("expected ID to be present, got %s", s)
		}
	})
}

func TestStatusConfig_OmitsEmptyOptionsAndGroups(t *testing.T) {
	t.Run("nil options and groups are omitted", func(t *testing.T) {
		sc := StatusConfig{}
		b, err := json.Marshal(sc)
		if err != nil {
			t.Fatal(err)
		}
		s := string(b)
		if strings.Contains(s, `"options"`) {
			t.Errorf("expected empty options to be omitted, got %s", s)
		}
		if strings.Contains(s, `"groups"`) {
			t.Errorf("expected empty groups to be omitted, got %s", s)
		}
		want := `{}`
		if s != want {
			t.Errorf("got %s, want %s", s, want)
		}
	})

	t.Run("populated options are included", func(t *testing.T) {
		sc := StatusConfig{
			Options: []Option{{Name: "Active", Color: "blue"}},
		}
		b, err := json.Marshal(sc)
		if err != nil {
			t.Fatal(err)
		}
		s := string(b)
		if !strings.Contains(s, `"options"`) {
			t.Errorf("expected options to be present, got %s", s)
		}
		if strings.Contains(s, `"groups"`) {
			t.Errorf("expected empty groups to be omitted, got %s", s)
		}
	})
}

func TestNumberFormat_OmitsEmptyFormat(t *testing.T) {
	t.Run("empty format is omitted", func(t *testing.T) {
		nf := NumberFormat{}
		b, err := json.Marshal(nf)
		if err != nil {
			t.Fatal(err)
		}
		want := `{}`
		if string(b) != want {
			t.Errorf("got %s, want %s", string(b), want)
		}
	})

	t.Run("populated format is included", func(t *testing.T) {
		nf := NumberFormat{Format: FormatDollar}
		b, err := json.Marshal(nf)
		if err != nil {
			t.Fatal(err)
		}
		want := `{"format":"dollar"}`
		if string(b) != want {
			t.Errorf("got %s, want %s", string(b), want)
		}
	})
}

func TestSearchFilter_OmitsEmptyFields(t *testing.T) {
	t.Run("empty filter serializes to empty object", func(t *testing.T) {
		sf := SearchFilter{}
		b, err := json.Marshal(sf)
		if err != nil {
			t.Fatal(err)
		}
		want := `{}`
		if string(b) != want {
			t.Errorf("got %s, want %s", string(b), want)
		}
	})

	t.Run("populated filter includes fields", func(t *testing.T) {
		sf := SearchFilter{Value: "page", Property: "object"}
		b, err := json.Marshal(sf)
		if err != nil {
			t.Fatal(err)
		}
		want := `{"value":"page","property":"object"}`
		if string(b) != want {
			t.Errorf("got %s, want %s", string(b), want)
		}
	})

	t.Run("partial filter omits empty field", func(t *testing.T) {
		sf := SearchFilter{Value: "database"}
		b, err := json.Marshal(sf)
		if err != nil {
			t.Fatal(err)
		}
		s := string(b)
		if !strings.Contains(s, `"value":"database"`) {
			t.Errorf("expected value to be present, got %s", s)
		}
		if strings.Contains(s, `"property"`) {
			t.Errorf("expected empty property to be omitted, got %s", s)
		}
	})
}

func TestDatabaseCreateRequest_InitialDataSource(t *testing.T) {
	t.Run("omits initial_data_source when nil", func(t *testing.T) {
		req := DatabaseCreateRequest{
			Parent: Parent{Type: ParentTypePageID, PageID: "pid"},
			Title:  []RichText{{Type: RichTextTypeText, Text: &Text{Content: "Test"}}},
		}
		b, err := json.Marshal(req)
		if err != nil {
			t.Fatal(err)
		}
		s := string(b)
		if strings.Contains(s, `"initial_data_source"`) {
			t.Errorf("expected initial_data_source to be omitted, got %s", s)
		}
	})

	t.Run("includes initial_data_source when set", func(t *testing.T) {
		req := DatabaseCreateRequest{
			Parent: Parent{Type: ParentTypePageID, PageID: "pid"},
			Title:  []RichText{{Type: RichTextTypeText, Text: &Text{Content: "Test"}}},
			InitialDataSource: &InitialDataSource{
				Properties: PropertyConfigs{
					"Name": TitlePropertyConfig{Type: PropertyConfigTypeTitle},
				},
			},
		}
		b, err := json.Marshal(req)
		if err != nil {
			t.Fatal(err)
		}
		s := string(b)
		if !strings.Contains(s, `"initial_data_source"`) {
			t.Errorf("expected initial_data_source to be present, got %s", s)
		}
	})
}

func TestDatabaseWithDataSources_Unmarshal(t *testing.T) {
	data := `{
		"object": "database",
		"id": "db_id",
		"created_time": "2021-05-24T05:06:34.827Z",
		"last_edited_time": "2021-05-24T05:06:34.827Z",
		"title": [],
		"parent": {"type": "page_id", "page_id": "pid"},
		"url": "",
		"public_url": "",
		"description": [],
		"is_inline": false,
		"archived": false,
		"data_sources": [
			{"id": "ds1", "name": "Main"},
			{"id": "ds2", "name": "External"}
		]
	}`
	var db Database
	if err := json.Unmarshal([]byte(data), &db); err != nil {
		t.Fatal(err)
	}
	if len(db.DataSources) != 2 {
		t.Fatalf("expected 2 data sources, got %d", len(db.DataSources))
	}
	if db.DataSources[0].ID != "ds1" || db.DataSources[0].Name != "Main" {
		t.Errorf("data_sources[0] = %+v, want {ds1, Main}", db.DataSources[0])
	}
	if db.DataSources[1].ID != "ds2" || db.DataSources[1].Name != "External" {
		t.Errorf("data_sources[1] = %+v, want {ds2, External}", db.DataSources[1])
	}
}

func TestParentWithDataSourceID(t *testing.T) {
	t.Run("unmarshal", func(t *testing.T) {
		data := `{"type":"data_source_id","data_source_id":"ds_abc"}`
		var p Parent
		if err := json.Unmarshal([]byte(data), &p); err != nil {
			t.Fatal(err)
		}
		if p.Type != ParentTypeDataSourceID {
			t.Errorf("type = %v, want %v", p.Type, ParentTypeDataSourceID)
		}
		if p.DataSourceID != "ds_abc" {
			t.Errorf("data_source_id = %v, want ds_abc", p.DataSourceID)
		}
	})

	t.Run("marshal", func(t *testing.T) {
		p := Parent{
			Type:         ParentTypeDataSourceID,
			DataSourceID: "ds_abc",
		}
		b, err := json.Marshal(p)
		if err != nil {
			t.Fatal(err)
		}
		want := `{"type":"data_source_id","data_source_id":"ds_abc"}`
		if string(b) != want {
			t.Errorf("got %s, want %s", string(b), want)
		}
	})
}

func TestSearchResponse_DataSourceObject(t *testing.T) {
	data := `{
		"object": "list",
		"results": [
			{
				"object": "data_source",
				"id": "ds_id",
				"title": [{"type":"text","text":{"content":"Test"},"plain_text":"Test"}],
				"properties": {
					"Name": {"id":"title","type":"title","title":{}}
				},
				"parent": {"type":"database_id","database_id":"db_id"}
			}
		],
		"next_cursor": null,
		"has_more": false
	}`

	var sr SearchResponse
	if err := json.Unmarshal([]byte(data), &sr); err != nil {
		t.Fatal(err)
	}
	if len(sr.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(sr.Results))
	}

	ds, ok := sr.Results[0].(*DataSource)
	if !ok {
		t.Fatalf("expected *DataSource, got %T", sr.Results[0])
	}
	if ds.Object != ObjectTypeDataSource {
		t.Errorf("object = %v, want %v", ds.Object, ObjectTypeDataSource)
	}
	if string(ds.ID) != "ds_id" {
		t.Errorf("id = %v, want ds_id", ds.ID)
	}
}

func TestSearchResponse_MixedResults(t *testing.T) {
	data := `{
		"object": "list",
		"results": [
			{
				"object": "page",
				"id": "page_id",
				"created_time": "2021-05-24T05:06:34.827Z",
				"last_edited_time": "2021-05-24T05:06:34.827Z",
				"parent": {"type":"database_id","database_id":"db_id"},
				"archived": false,
				"properties": {"Name":{"id":"title","type":"title","title":[]}}
			},
			{
				"object": "data_source",
				"id": "ds_id",
				"title": [],
				"properties": {"Name":{"id":"title","type":"title","title":{}}},
				"parent": {"type":"database_id","database_id":"db_id"}
			},
			{
				"object": "database",
				"id": "db_id",
				"created_time": "2021-05-24T05:06:34.827Z",
				"last_edited_time": "2021-05-24T05:06:34.827Z",
				"title": [],
				"properties": {"Name":{"id":"title","type":"title","title":{}}},
				"parent": {"type":"page_id","page_id":"pid"}
			}
		],
		"next_cursor": null,
		"has_more": false
	}`

	var sr SearchResponse
	if err := json.Unmarshal([]byte(data), &sr); err != nil {
		t.Fatal(err)
	}
	if len(sr.Results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(sr.Results))
	}

	types := []ObjectType{ObjectTypePage, ObjectTypeDataSource, ObjectTypeDatabase}
	for i, want := range types {
		got := sr.Results[i].GetObject()
		if got != want {
			t.Errorf("result[%d] object = %v, want %v", i, got, want)
		}
	}
}

func TestDataSource_Unmarshal(t *testing.T) {
	data := `{
		"object": "data_source",
		"id": "ds_id",
		"title": [{"type":"text","text":{"content":"Source"},"plain_text":"Source"}],
		"properties": {
			"Name": {"id":"title","type":"title","title":{}},
			"Count": {"id":"cnt","type":"number","number":{"format":"number"}}
		},
		"parent": {"type":"database_id","database_id":"db_parent"},
		"description": [{"type":"text","text":{"content":"A data source"},"plain_text":"A data source"}]
	}`

	var ds DataSource
	if err := json.Unmarshal([]byte(data), &ds); err != nil {
		t.Fatal(err)
	}

	if ds.Object != ObjectTypeDataSource {
		t.Errorf("object = %v, want data_source", ds.Object)
	}
	if string(ds.ID) != "ds_id" {
		t.Errorf("id = %v, want ds_id", ds.ID)
	}
	if len(ds.Title) != 1 || ds.Title[0].PlainText != "Source" {
		t.Errorf("title = %+v, want [Source]", ds.Title)
	}
	if len(ds.Properties) != 2 {
		t.Errorf("expected 2 properties, got %d", len(ds.Properties))
	}
	if ds.Parent.DatabaseID != "db_parent" {
		t.Errorf("parent.database_id = %v, want db_parent", ds.Parent.DatabaseID)
	}
	if len(ds.Description) != 1 {
		t.Errorf("expected 1 description element, got %d", len(ds.Description))
	}
}

func TestDataSourceUpdateRequest_Marshal(t *testing.T) {
	t.Run("omits empty fields", func(t *testing.T) {
		req := DataSourceUpdateRequest{}
		b, err := json.Marshal(req)
		if err != nil {
			t.Fatal(err)
		}
		want := `{}`
		if string(b) != want {
			t.Errorf("got %s, want %s", string(b), want)
		}
	})

	t.Run("includes set fields", func(t *testing.T) {
		req := DataSourceUpdateRequest{
			Title: []RichText{{Type: RichTextTypeText, Text: &Text{Content: "New"}}},
		}
		b, err := json.Marshal(req)
		if err != nil {
			t.Fatal(err)
		}
		s := string(b)
		if !strings.Contains(s, `"title"`) {
			t.Errorf("expected title to be present, got %s", s)
		}
		if strings.Contains(s, `"properties"`) {
			t.Errorf("expected empty properties to be omitted, got %s", s)
		}
	})
}

func TestNumberPropertyConfig_OmitsEmptyFormat(t *testing.T) {
	// Creating a number property without specifying format should not send "format":""
	p := NumberPropertyConfig{
		Type:   PropertyConfigTypeNumber,
		Number: NumberFormat{},
	}
	b, err := json.Marshal(p)
	if err != nil {
		t.Fatal(err)
	}
	want := []byte(`{"type":"number","number":{}}`)
	if !reflect.DeepEqual(b, want) {
		t.Errorf("got %s, want %s", string(b), string(want))
	}
}
