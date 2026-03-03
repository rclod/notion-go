package notionapi_test

import (
	"context"
	"net/http"
	"reflect"
	"testing"
	"time"

	notionapi "github.com/rclod/notion-go"
)

func TestDataSourceClient(t *testing.T) {
	timestamp, err := time.Parse(time.RFC3339, "2021-05-24T05:06:34.827Z")
	if err != nil {
		t.Fatal(err)
	}

	var user = notionapi.User{
		Object: "user",
		ID:     "some_id",
	}

	t.Run("Get", func(t *testing.T) {
		tests := []struct {
			name       string
			filePath   string
			statusCode int
			id         notionapi.DataSourceID
			wantErr    bool
		}{
			{
				name:       "returns data source by id",
				id:         "ds_some_id",
				filePath:   "testdata/data_source_get.json",
				statusCode: http.StatusOK,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				c := newMockedClient(t, tt.filePath, tt.statusCode)
				client := notionapi.NewClient("some_token", notionapi.WithHTTPClient(c))

				got, err := client.DataSource.Get(context.Background(), tt.id)
				if (err != nil) != tt.wantErr {
					t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				if got.Object != notionapi.ObjectTypeDataSource {
					t.Errorf("Get() object = %v, want %v", got.Object, notionapi.ObjectTypeDataSource)
				}
				if string(got.ID) != "ds_some_id" {
					t.Errorf("Get() id = %v, want ds_some_id", got.ID)
				}
				if got.Parent.DatabaseID != "parent_db_id" {
					t.Errorf("Get() parent.database_id = %v, want parent_db_id", got.Parent.DatabaseID)
				}
				if got.Properties == nil {
					t.Error("Get() properties is nil")
				}
				if _, ok := got.Properties["Name"]; !ok {
					t.Error("Get() missing 'Name' property")
				}
				if _, ok := got.Properties["Status"]; !ok {
					t.Error("Get() missing 'Status' property")
				}
			})
		}
	})

	t.Run("Update", func(t *testing.T) {
		tests := []struct {
			name       string
			filePath   string
			statusCode int
			id         notionapi.DataSourceID
			request    *notionapi.DataSourceUpdateRequest
			wantErr    bool
		}{
			{
				name:       "returns updated data source",
				filePath:   "testdata/data_source_update.json",
				statusCode: http.StatusOK,
				id:         "ds_some_id",
				request: &notionapi.DataSourceUpdateRequest{
					Properties: notionapi.PropertyConfigs{
						"Priority": notionapi.SelectPropertyConfig{
							Type: notionapi.PropertyConfigTypeSelect,
							Select: notionapi.Select{
								Options: []notionapi.Option{
									{Name: "High", Color: "red"},
								},
							},
						},
					},
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				c := newMockedClient(t, tt.filePath, tt.statusCode)
				client := notionapi.NewClient("some_token", notionapi.WithHTTPClient(c))

				got, err := client.DataSource.Update(context.Background(), tt.id, tt.request)
				if (err != nil) != tt.wantErr {
					t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				if got.Object != notionapi.ObjectTypeDataSource {
					t.Errorf("Update() object = %v, want %v", got.Object, notionapi.ObjectTypeDataSource)
				}
				if len(got.Title) == 0 || got.Title[0].PlainText != "Updated Source" {
					t.Errorf("Update() title = %v, want 'Updated Source'", got.Title)
				}
				if _, ok := got.Properties["Priority"]; !ok {
					t.Error("Update() missing 'Priority' property")
				}
			})
		}
	})

	t.Run("Query", func(t *testing.T) {
		tests := []struct {
			name       string
			filePath   string
			statusCode int
			id         notionapi.DataSourceID
			request    *notionapi.DatabaseQueryRequest
			want       *notionapi.DatabaseQueryResponse
			wantErr    bool
		}{
			{
				name:       "returns query results",
				id:         "ds_some_id",
				filePath:   "testdata/data_source_query.json",
				statusCode: http.StatusOK,
				request: &notionapi.DatabaseQueryRequest{
					Filter: &notionapi.PropertyFilter{
						Property: "Name",
						RichText: &notionapi.TextFilterCondition{
							Contains: "Test",
						},
					},
				},
				want: &notionapi.DatabaseQueryResponse{
					Object: notionapi.ObjectTypeList,
					Results: []notionapi.Page{
						{
							Object:         notionapi.ObjectTypePage,
							ID:             "page_id_1",
							CreatedTime:    timestamp,
							LastEditedTime: timestamp,
							CreatedBy:      user,
							LastEditedBy:   user,
							Parent: notionapi.Parent{
								Type:         notionapi.ParentTypeDataSourceID,
								DataSourceID: "ds_some_id",
							},
							Archived: false,
							URL:      "some_url",
						},
					},
					HasMore:    false,
					NextCursor: "",
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				c := newMockedClient(t, tt.filePath, tt.statusCode)
				client := notionapi.NewClient("some_token", notionapi.WithHTTPClient(c))

				got, err := client.DataSource.Query(context.Background(), tt.id, tt.request)
				if (err != nil) != tt.wantErr {
					t.Errorf("Query() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				// Clear properties for comparison (same pattern as database tests)
				got.Results[0].Properties = nil
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Query() got = %v, want %v", got, tt.want)
				}
			})
		}
	})
}
