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

func TestSearchRequest_OmitsEmptyFilter(t *testing.T) {
	t.Run("no filter omits filter field", func(t *testing.T) {
		req := SearchRequest{Query: "test", PageSize: 10}
		b, err := json.Marshal(req)
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(string(b), `"filter"`) {
			t.Errorf("expected filter to be omitted, got %s", string(b))
		}
	})

	t.Run("with filter includes filter field", func(t *testing.T) {
		req := SearchRequest{
			Query:  "test",
			Filter: &SearchFilter{Value: "page", Property: "object"},
		}
		b, err := json.Marshal(req)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(b), `"filter"`) {
			t.Errorf("expected filter to be present, got %s", string(b))
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

func TestPageCreateRequest_TemplateFields(t *testing.T) {
	t.Run("with template", func(t *testing.T) {
		req := PageCreateRequest{
			Parent:     Parent{Type: ParentTypeDataSourceID, DataSourceID: "ds_id"},
			Properties: Properties{},
			Template: &PageTemplate{
				Type:       TemplateTypeTemplateID,
				TemplateID: "tmpl_123",
			},
		}
		b, err := json.Marshal(req)
		if err != nil {
			t.Fatal(err)
		}
		s := string(b)
		if !strings.Contains(s, `"template"`) {
			t.Errorf("expected template to be present, got %s", s)
		}
		if !strings.Contains(s, `"template_id":"tmpl_123"`) {
			t.Errorf("expected template_id to be present, got %s", s)
		}
	})

	t.Run("with markdown", func(t *testing.T) {
		req := PageCreateRequest{
			Parent:     Parent{Type: ParentTypeDataSourceID, DataSourceID: "ds_id"},
			Properties: Properties{},
			Markdown:   "# Hello\n\nWorld",
		}
		b, err := json.Marshal(req)
		if err != nil {
			t.Fatal(err)
		}
		s := string(b)
		if !strings.Contains(s, `"markdown":"# Hello\n\nWorld"`) {
			t.Errorf("expected markdown to be present, got %s", s)
		}
	})

	t.Run("with position", func(t *testing.T) {
		req := PageCreateRequest{
			Parent:     Parent{Type: ParentTypeDataSourceID, DataSourceID: "ds_id"},
			Properties: Properties{},
			Position: &PagePosition{
				Type:       PositionTypeAfterBlock,
				AfterBlock: &AfterBlockRef{ID: "block_abc"},
			},
		}
		b, err := json.Marshal(req)
		if err != nil {
			t.Fatal(err)
		}
		s := string(b)
		if !strings.Contains(s, `"position"`) {
			t.Errorf("expected position to be present, got %s", s)
		}
		if !strings.Contains(s, `"after_block"`) {
			t.Errorf("expected after_block to be present, got %s", s)
		}
	})

	t.Run("omits optional fields when nil", func(t *testing.T) {
		req := PageCreateRequest{
			Parent:     Parent{Type: ParentTypeDataSourceID, DataSourceID: "ds_id"},
			Properties: Properties{},
		}
		b, err := json.Marshal(req)
		if err != nil {
			t.Fatal(err)
		}
		s := string(b)
		if strings.Contains(s, `"template"`) {
			t.Errorf("expected template to be omitted, got %s", s)
		}
		if strings.Contains(s, `"markdown"`) {
			t.Errorf("expected markdown to be omitted, got %s", s)
		}
		if strings.Contains(s, `"position"`) {
			t.Errorf("expected position to be omitted, got %s", s)
		}
	})
}

func TestPageUpdateRequest_NewFields(t *testing.T) {
	t.Run("with template and erase_content", func(t *testing.T) {
		eraseContent := true
		req := PageUpdateRequest{
			Template: &PageTemplate{
				Type: TemplateTypeDefault,
			},
			EraseContent: &eraseContent,
		}
		b, err := json.Marshal(req)
		if err != nil {
			t.Fatal(err)
		}
		s := string(b)
		if !strings.Contains(s, `"template"`) {
			t.Errorf("expected template to be present, got %s", s)
		}
		if !strings.Contains(s, `"erase_content":true`) {
			t.Errorf("expected erase_content to be present, got %s", s)
		}
	})

	t.Run("with is_locked and in_trash", func(t *testing.T) {
		isLocked := true
		inTrash := false
		req := PageUpdateRequest{
			IsLocked: &isLocked,
			InTrash:  &inTrash,
		}
		b, err := json.Marshal(req)
		if err != nil {
			t.Fatal(err)
		}
		s := string(b)
		if !strings.Contains(s, `"is_locked":true`) {
			t.Errorf("expected is_locked to be present, got %s", s)
		}
		if !strings.Contains(s, `"in_trash":false`) {
			t.Errorf("expected in_trash to be present, got %s", s)
		}
	})

	t.Run("omits optional fields when nil", func(t *testing.T) {
		req := PageUpdateRequest{}
		b, err := json.Marshal(req)
		if err != nil {
			t.Fatal(err)
		}
		s := string(b)
		if strings.Contains(s, `"template"`) {
			t.Errorf("expected template to be omitted, got %s", s)
		}
		if strings.Contains(s, `"erase_content"`) {
			t.Errorf("expected erase_content to be omitted, got %s", s)
		}
		if strings.Contains(s, `"is_locked"`) {
			t.Errorf("expected is_locked to be omitted, got %s", s)
		}
		if strings.Contains(s, `"in_trash"`) {
			t.Errorf("expected in_trash to be omitted, got %s", s)
		}
	})
}

func TestDataSourceCreateRequest_Marshal(t *testing.T) {
	t.Run("required and optional fields", func(t *testing.T) {
		req := DataSourceCreateRequest{
			Parent: Parent{Type: ParentTypeDatabaseID, DatabaseID: "db_id"},
			Properties: PropertyConfigs{
				"Name": TitlePropertyConfig{Type: PropertyConfigTypeTitle},
			},
			Title: []RichText{{Type: RichTextTypeText, Text: &Text{Content: "Source"}}},
		}
		b, err := json.Marshal(req)
		if err != nil {
			t.Fatal(err)
		}
		s := string(b)
		if !strings.Contains(s, `"parent"`) {
			t.Errorf("expected parent to be present, got %s", s)
		}
		if !strings.Contains(s, `"properties"`) {
			t.Errorf("expected properties to be present, got %s", s)
		}
		if !strings.Contains(s, `"title"`) {
			t.Errorf("expected title to be present, got %s", s)
		}
	})

	t.Run("omits optional fields when empty", func(t *testing.T) {
		req := DataSourceCreateRequest{
			Parent:     Parent{Type: ParentTypeDatabaseID, DatabaseID: "db_id"},
			Properties: PropertyConfigs{},
		}
		b, err := json.Marshal(req)
		if err != nil {
			t.Fatal(err)
		}
		s := string(b)
		if strings.Contains(s, `"icon"`) {
			t.Errorf("expected icon to be omitted, got %s", s)
		}
	})
}

func TestTemplateListResponse_Unmarshal(t *testing.T) {
	data := `{
		"templates": [
			{"id": "t1", "name": "Template A", "is_default": true},
			{"id": "t2", "name": "Template B", "is_default": false}
		],
		"has_more": true,
		"next_cursor": "cursor_abc"
	}`
	var resp TemplateListResponse
	if err := json.Unmarshal([]byte(data), &resp); err != nil {
		t.Fatal(err)
	}
	if len(resp.Templates) != 2 {
		t.Fatalf("expected 2 templates, got %d", len(resp.Templates))
	}
	if resp.Templates[0].ID != "t1" || resp.Templates[0].Name != "Template A" || !resp.Templates[0].IsDefault {
		t.Errorf("templates[0] = %+v", resp.Templates[0])
	}
	if resp.Templates[1].ID != "t2" || resp.Templates[1].IsDefault {
		t.Errorf("templates[1] = %+v", resp.Templates[1])
	}
	if !resp.HasMore {
		t.Error("expected has_more = true")
	}
	if resp.NextCursor != "cursor_abc" {
		t.Errorf("next_cursor = %v, want cursor_abc", resp.NextCursor)
	}
}

func TestPageMarkdown_Unmarshal(t *testing.T) {
	data := `{
		"object": "page_markdown",
		"id": "page_id",
		"markdown": "# Title\n\nContent here.",
		"truncated": true,
		"unknown_block_ids": ["blk_1", "blk_2"]
	}`
	var pm PageMarkdown
	if err := json.Unmarshal([]byte(data), &pm); err != nil {
		t.Fatal(err)
	}
	if pm.Object != ObjectTypePageMarkdown {
		t.Errorf("object = %v, want page_markdown", pm.Object)
	}
	if string(pm.ID) != "page_id" {
		t.Errorf("id = %v, want page_id", pm.ID)
	}
	if pm.Markdown != "# Title\n\nContent here." {
		t.Errorf("markdown = %v", pm.Markdown)
	}
	if !pm.Truncated {
		t.Error("expected truncated = true")
	}
	if len(pm.UnknownBlockIDs) != 2 {
		t.Fatalf("expected 2 unknown_block_ids, got %d", len(pm.UnknownBlockIDs))
	}
}

func TestMarkdownUpdateRequest_Marshal(t *testing.T) {
	t.Run("insert_content", func(t *testing.T) {
		req := MarkdownUpdateRequest{
			Type: "insert_content",
			InsertContent: &InsertContent{
				Content: "## New section\n",
				After:   "block_abc",
			},
		}
		b, err := json.Marshal(req)
		if err != nil {
			t.Fatal(err)
		}
		s := string(b)
		if !strings.Contains(s, `"type":"insert_content"`) {
			t.Errorf("expected type to be insert_content, got %s", s)
		}
		if !strings.Contains(s, `"insert_content"`) {
			t.Errorf("expected insert_content object, got %s", s)
		}
		if strings.Contains(s, `"replace_content_range"`) {
			t.Errorf("expected replace_content_range to be omitted, got %s", s)
		}
	})

	t.Run("replace_content_range", func(t *testing.T) {
		req := MarkdownUpdateRequest{
			Type: "replace_content_range",
			ReplaceContentRange: &ReplaceContentRange{
				Content:              "Updated content",
				ContentRange:         "## Old heading\nOld text",
				AllowDeletingContent: true,
			},
		}
		b, err := json.Marshal(req)
		if err != nil {
			t.Fatal(err)
		}
		s := string(b)
		if !strings.Contains(s, `"type":"replace_content_range"`) {
			t.Errorf("expected type to be replace_content_range, got %s", s)
		}
		if !strings.Contains(s, `"allow_deleting_content":true`) {
			t.Errorf("expected allow_deleting_content, got %s", s)
		}
		if strings.Contains(s, `"insert_content"`) {
			t.Errorf("expected insert_content to be omitted, got %s", s)
		}
	})
}

func TestPageMoveRequest_Marshal(t *testing.T) {
	req := PageMoveRequest{
		Parent: Parent{
			Type:   ParentTypePageID,
			PageID: "new_parent",
		},
	}
	b, err := json.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}
	s := string(b)
	if !strings.Contains(s, `"page_id":"new_parent"`) {
		t.Errorf("expected page_id to be present, got %s", s)
	}
}

func TestDataSource_FullUnmarshal(t *testing.T) {
	data := `{
		"object": "data_source",
		"id": "ds_full",
		"title": [{"type":"text","text":{"content":"Full"},"plain_text":"Full"}],
		"properties": {"Name":{"id":"title","type":"title","title":{}}},
		"parent": {"type":"database_id","database_id":"db_parent"},
		"database_parent": {"type":"database_id","database_id":"db_container"},
		"is_inline": true,
		"archived": false,
		"in_trash": true,
		"created_time": "2024-01-15T10:00:00.000Z",
		"last_edited_time": "2024-06-20T15:30:00.000Z",
		"created_by": {"object":"user","id":"user_1"},
		"last_edited_by": {"object":"user","id":"user_2"},
		"url": "https://notion.so/ds_full",
		"public_url": "https://notion.so/public/ds_full"
	}`

	var ds DataSource
	if err := json.Unmarshal([]byte(data), &ds); err != nil {
		t.Fatal(err)
	}
	if ds.Object != ObjectTypeDataSource {
		t.Errorf("object = %v, want data_source", ds.Object)
	}
	if string(ds.ID) != "ds_full" {
		t.Errorf("id = %v, want ds_full", ds.ID)
	}
	if ds.DatabaseParent.DatabaseID != "db_container" {
		t.Errorf("database_parent.database_id = %v, want db_container", ds.DatabaseParent.DatabaseID)
	}
	if !ds.IsInline {
		t.Error("expected is_inline = true")
	}
	if ds.Archived {
		t.Error("expected archived = false")
	}
	if !ds.InTrash {
		t.Error("expected in_trash = true")
	}
	if ds.CreatedTime.IsZero() {
		t.Error("expected created_time to be set")
	}
	if ds.LastEditedTime.IsZero() {
		t.Error("expected last_edited_time to be set")
	}
	if string(ds.CreatedBy.ID) != "user_1" {
		t.Errorf("created_by.id = %v, want user_1", ds.CreatedBy.ID)
	}
	if string(ds.LastEditedBy.ID) != "user_2" {
		t.Errorf("last_edited_by.id = %v, want user_2", ds.LastEditedBy.ID)
	}
	if ds.URL != "https://notion.so/ds_full" {
		t.Errorf("url = %v", ds.URL)
	}
	if ds.PublicURL != "https://notion.so/public/ds_full" {
		t.Errorf("public_url = %v", ds.PublicURL)
	}
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
