package notionapi_test

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/rclod/notion-go"
)

func TestPageClient(t *testing.T) {
	timestamp, err := time.Parse(time.RFC3339, "2021-05-24T05:06:34.827Z")
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Get", func(t *testing.T) {
		tests := []struct {
			name       string
			filePath   string
			statusCode int
			id         notionapi.PageID
			want       *notionapi.Page
			wantErr    bool
			err        error
		}{
			{
				name:       "returns page by id",
				id:         "some_id",
				filePath:   "testdata/page_get.json",
				statusCode: http.StatusOK,
				want: &notionapi.Page{
					Object:         notionapi.ObjectTypePage,
					ID:             "some_id",
					CreatedTime:    timestamp,
					LastEditedTime: timestamp,
					CreatedBy: notionapi.User{
						Object: "user",
						ID:     "some_id",
					},
					LastEditedBy: notionapi.User{
						Object: "user",
						ID:     "some_id",
					},
					Parent: notionapi.Parent{
						Type:       notionapi.ParentTypeDatabaseID,
						DatabaseID: "some_id",
					},
					Archived: false,
					URL:      "some_url",
					Properties: notionapi.Properties{
						"Tags": &notionapi.MultiSelectProperty{
							ID:   ";s|V",
							Type: "multi_select",
							MultiSelect: []notionapi.Option{
								{
									ID:    "some_id",
									Name:  "tag",
									Color: "blue",
								},
							},
						},
						"Some another column": &notionapi.PeopleProperty{
							ID:   "rJt\\",
							Type: "people",
							People: []notionapi.User{
								{
									Object:    "user",
									ID:        "some_id",
									Name:      "some name",
									AvatarURL: "some.url",
									Type:      "person",
									Person: &notionapi.Person{
										Email: "some@email.com",
									},
								},
							},
						},
						"SomeColumn": &notionapi.RichTextProperty{
							ID:   "~j_@",
							Type: "rich_text",
							RichText: []notionapi.RichText{
								{
									Type: "text",
									Text: &notionapi.Text{
										Content: "some text",
									},
									Annotations: &notionapi.Annotations{
										Color: "default",
									},
									PlainText: "some text",
								},
							},
						},
						"Some Files": &notionapi.FilesProperty{
							ID:   "files",
							Type: "files",
							Files: []notionapi.File{
								{
									Name: "https://google.com",
									Type: "external",
									External: &notionapi.FileObject{
										URL: "https://google.com",
									},
								},
							},
						},
						"Name": &notionapi.TitleProperty{
							ID:   "title",
							Type: "title",
							Title: []notionapi.RichText{
								{
									Type: "text",
									Text: &notionapi.Text{
										Content: "Hello",
									},
									Annotations: &notionapi.Annotations{
										Color: "default",
									},
									PlainText: "Hello",
								},
							},
						},
						"RollupArray": &notionapi.RollupProperty{
							ID:   "abcd",
							Type: "rollup",
							Rollup: notionapi.Rollup{
								Type: "array",
								Array: notionapi.PropertyArray{
									&notionapi.NumberProperty{
										Type:   "number",
										Number: 42.2,
									},
									&notionapi.NumberProperty{
										Type:   "number",
										Number: 56,
									},
								},
							},
						},
					},
				},
			},
			{
				name:       "returns validation error for invalid request",
				id:         "some_id",
				filePath:   "testdata/validation_error.json",
				statusCode: http.StatusBadRequest,
				wantErr:    true,
				err: &notionapi.Error{
					Object:  notionapi.ObjectTypeError,
					Status:  http.StatusBadRequest,
					Code:    "validation_error",
					Message: "The provided page ID is not a valid Notion UUID: bla bla.",
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				c := newMockedClient(t, tt.filePath, tt.statusCode)
				client := notionapi.NewClient("some_token", notionapi.WithHTTPClient(c))

				got, err := client.Page.Get(context.Background(), tt.id)
				if err != nil {
					if tt.wantErr {
						if !reflect.DeepEqual(err, tt.err) {
							t.Errorf("Get error() got = %v, want %v", err, tt.err)
						}
					} else {
						t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)

					}
					return
				}

				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Get() got = %v, want %v", got, tt.want)
				}
			})
		}
	})

	t.Run("Create", func(t *testing.T) {
		tests := []struct {
			name       string
			filePath   string
			statusCode int
			id         notionapi.PageID
			request    *notionapi.PageCreateRequest
			want       *notionapi.Page
			wantErr    bool
			err        error
		}{
			{
				name:       "returns a new page",
				filePath:   "testdata/page_create.json",
				statusCode: http.StatusOK,
				request: &notionapi.PageCreateRequest{
					Parent: notionapi.Parent{
						Type:       notionapi.ParentTypeDatabaseID,
						DatabaseID: "f830be5eff534859932e5b81542b3c7b",
					},
					Properties: notionapi.Properties{
						"Name": notionapi.TitleProperty{
							Title: []notionapi.RichText{
								{Text: &notionapi.Text{Content: "hello"}},
							},
						},
					},
				},
				want: &notionapi.Page{
					Object:         notionapi.ObjectTypePage,
					ID:             "some_id",
					CreatedTime:    timestamp,
					LastEditedTime: timestamp,
					CreatedBy: notionapi.User{
						Object: "user",
						ID:     "some_id",
					},
					LastEditedBy: notionapi.User{
						Object: "user",
						ID:     "some_id",
					},
					Parent: notionapi.Parent{
						Type:       notionapi.ParentTypeDatabaseID,
						DatabaseID: "some_id",
					},
					Archived: false,
					URL:      "some_url",
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				c := newMockedClient(t, tt.filePath, tt.statusCode)
				client := notionapi.NewClient("some_token", notionapi.WithHTTPClient(c))
				got, err := client.Page.Create(context.Background(), tt.request)
				if (err != nil) != tt.wantErr {
					t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				// TODO: remove properties from comparing for a while. Have to compare with interface somehow
				got.Properties = nil
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Create() got = %v, want %v", got, tt.want)
				}
			})
		}
	})

	t.Run("Move", func(t *testing.T) {
		tests := []struct {
			name       string
			filePath   string
			statusCode int
			id         notionapi.PageID
			request    *notionapi.PageMoveRequest
			wantErr    bool
		}{
			{
				name:       "moves a page to a new parent",
				id:         "page_moved_id",
				filePath:   "testdata/page_move.json",
				statusCode: http.StatusOK,
				request: &notionapi.PageMoveRequest{
					Parent: notionapi.Parent{
						Type:   notionapi.ParentTypePageID,
						PageID: "new_parent_page_id",
					},
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				c := newMockedClient(t, tt.filePath, tt.statusCode)
				client := notionapi.NewClient("some_token", notionapi.WithHTTPClient(c))

				got, err := client.Page.Move(context.Background(), tt.id, tt.request)
				if (err != nil) != tt.wantErr {
					t.Errorf("Move() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				if got.Parent.Type != notionapi.ParentTypePageID {
					t.Errorf("Move() parent.type = %v, want %v", got.Parent.Type, notionapi.ParentTypePageID)
				}
				if got.Parent.PageID != "new_parent_page_id" {
					t.Errorf("Move() parent.page_id = %v, want new_parent_page_id", got.Parent.PageID)
				}
			})
		}
	})

	t.Run("GetMarkdown", func(t *testing.T) {
		tests := []struct {
			name       string
			filePath   string
			statusCode int
			id         notionapi.PageID
			wantErr    bool
		}{
			{
				name:       "returns page markdown",
				id:         "page_md_id",
				filePath:   "testdata/page_markdown_get.json",
				statusCode: http.StatusOK,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				c := newMockedClient(t, tt.filePath, tt.statusCode)
				client := notionapi.NewClient("some_token", notionapi.WithHTTPClient(c))

				got, err := client.Page.GetMarkdown(context.Background(), tt.id)
				if (err != nil) != tt.wantErr {
					t.Errorf("GetMarkdown() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				if got.Object != notionapi.ObjectTypePageMarkdown {
					t.Errorf("GetMarkdown() object = %v, want %v", got.Object, notionapi.ObjectTypePageMarkdown)
				}
				if string(got.ID) != "page_md_id" {
					t.Errorf("GetMarkdown() id = %v, want page_md_id", got.ID)
				}
				if got.Markdown == "" {
					t.Error("GetMarkdown() markdown is empty")
				}
				if got.Truncated {
					t.Error("GetMarkdown() truncated = true, want false")
				}
				if len(got.UnknownBlockIDs) != 0 {
					t.Errorf("GetMarkdown() unknown_block_ids = %v, want empty", got.UnknownBlockIDs)
				}
			})
		}
	})

	t.Run("UpdateMarkdown", func(t *testing.T) {
		tests := []struct {
			name       string
			filePath   string
			statusCode int
			id         notionapi.PageID
			request    *notionapi.MarkdownUpdateRequest
			wantErr    bool
		}{
			{
				name:       "updates page markdown",
				id:         "page_md_id",
				filePath:   "testdata/page_markdown_update.json",
				statusCode: http.StatusOK,
				request: &notionapi.MarkdownUpdateRequest{
					Type: "insert_content",
					InsertContent: &notionapi.InsertContent{
						Content: "- Item 3\n",
					},
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				c := newMockedClient(t, tt.filePath, tt.statusCode)
				client := notionapi.NewClient("some_token", notionapi.WithHTTPClient(c))

				got, err := client.Page.UpdateMarkdown(context.Background(), tt.id, tt.request)
				if (err != nil) != tt.wantErr {
					t.Errorf("UpdateMarkdown() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				if got.Object != notionapi.ObjectTypePageMarkdown {
					t.Errorf("UpdateMarkdown() object = %v, want %v", got.Object, notionapi.ObjectTypePageMarkdown)
				}
				if len(got.UnknownBlockIDs) != 1 {
					t.Errorf("UpdateMarkdown() unknown_block_ids = %v, want 1 element", got.UnknownBlockIDs)
				}
			})
		}
	})

	t.Run("Update", func(t *testing.T) {
		tests := []struct {
			name       string
			filePath   string
			statusCode int
			id         notionapi.PageID
			request    *notionapi.PageUpdateRequest
			want       *notionapi.Page
			wantErr    bool
			err        error
		}{
			{
				name:       "change requested properties and return the result",
				id:         "some_id",
				filePath:   "testdata/page_update.json",
				statusCode: http.StatusOK,
				request: &notionapi.PageUpdateRequest{
					Properties: notionapi.Properties{
						"SomeColumn": notionapi.RichTextProperty{
							Type: notionapi.PropertyTypeRichText,
							RichText: []notionapi.RichText{
								{
									Type: notionapi.RichTextTypeText,
									Text: &notionapi.Text{Content: "patch"},
								},
							},
						},
						"Important Files": notionapi.FilesProperty{
							Type: "files",
							Files: []notionapi.File{
								{
									Type: "external",
									Name: "https://google.com",
									External: &notionapi.FileObject{
										URL: "https://google.com",
									},
								},
								{
									Type: "external",
									Name: "https://123.com",
									External: &notionapi.FileObject{
										URL: "https://123.com",
									},
								},
							},
						},
					},
				},
				want: &notionapi.Page{
					Object:         notionapi.ObjectTypePage,
					ID:             "some_id",
					CreatedTime:    timestamp,
					LastEditedTime: timestamp,
					CreatedBy: notionapi.User{
						Object: "user",
						ID:     "some_id",
					},
					LastEditedBy: notionapi.User{
						Object: "user",
						ID:     "some_id",
					},
					Parent: notionapi.Parent{
						Type:       notionapi.ParentTypeDatabaseID,
						DatabaseID: "some_id",
					},
					Archived: false,
					URL:      "some_url",
				},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				c := newMockedClient(t, tt.filePath, tt.statusCode)
				client := notionapi.NewClient("some_token", notionapi.WithHTTPClient(c))
				got, err := client.Page.Update(context.Background(), tt.id, tt.request)
				if (err != nil) != tt.wantErr {
					t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				// TODO: remove properties from comparing for a while. Have to compare with interface somehow
				got.Properties = nil
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Update() got = %v, want %v", got, tt.want)
				}
			})
		}
	})
}

func TestPageCreateRequest_MarshallJSON(t *testing.T) {
	timeObj, err := time.Parse(time.RFC3339, "2020-12-08T12:00:00Z")
	if err != nil {
		t.Error(err)
		return
	}

	dateObj := notionapi.Date(timeObj)
	tests := []struct {
		name    string
		req     *notionapi.PageCreateRequest
		want    []byte
		wantErr bool
	}{
		{
			name: "create a page",
			req: &notionapi.PageCreateRequest{
				Parent: notionapi.Parent{
					DatabaseID: "some_id",
				},
				Properties: notionapi.Properties{
					"Type": notionapi.SelectProperty{
						Select: notionapi.Option{
							ID:    "some_id",
							Name:  "Article",
							Color: notionapi.ColorDefault,
						},
					},
					"Name": notionapi.TitleProperty{
						Title: []notionapi.RichText{
							{Text: &notionapi.Text{Content: "New Media Article"}},
						},
					},
					"Publishing/Release Date": notionapi.DateProperty{
						Date: &notionapi.DateObject{
							Start: &dateObj,
						},
					},
					"Link": notionapi.URLProperty{
						URL: "some_url",
					},
					"Summary": notionapi.TextProperty{
						Text: []notionapi.RichText{
							{
								Type: notionapi.RichTextTypeText,
								Text: &notionapi.Text{
									Content: "Some content",
								},
								Annotations: &notionapi.Annotations{
									Bold:  true,
									Color: notionapi.ColorBlue,
								},
								PlainText: "Some content",
							},
						},
					},
					"Read": notionapi.CheckboxProperty{
						Checkbox: false,
					},
				},
			},
			want: []byte(`{"parent":{"database_id":"some_id"},"properties":{"Link":{"url":"some_url"},"Name":{"title":[{"text":{"content":"New Media Article"}}]},"Publishing/Release Date":{"date":{"start":"2020-12-08T12:00:00Z","end":null}},"Read":{"checkbox":false},"Summary":{"text":[{"type":"text","text":{"content":"Some content"},"annotations":{"bold":true,"italic":false,"strikethrough":false,"underline":false,"code":false,"color":"blue"},"plain_text":"Some content"}]},"Type":{"select":{"id":"some_id","name":"Article","color":"default"}}}}`),
		},
		{
			name: "create a page with content",
			req: &notionapi.PageCreateRequest{
				Parent: notionapi.Parent{
					DatabaseID: "some_id",
				},
				Properties: notionapi.Properties{
					"Type": notionapi.SelectProperty{
						Select: notionapi.Option{
							ID:    "some_id",
							Name:  "Article",
							Color: notionapi.ColorDefault,
						},
					},
					"Name": notionapi.TitleProperty{
						Title: []notionapi.RichText{
							{Text: &notionapi.Text{Content: "New Media Article"}},
						},
					},
					"Publishing/Release Date": notionapi.DateProperty{
						Date: &notionapi.DateObject{
							Start: &dateObj,
						},
					},
					"Link": notionapi.URLProperty{
						URL: "some_url",
					},
					"Summary": notionapi.TextProperty{
						Text: []notionapi.RichText{
							{
								Type: notionapi.RichTextTypeText,
								Text: &notionapi.Text{
									Content: "Some content",
								},
								Annotations: &notionapi.Annotations{
									Bold:  true,
									Color: notionapi.ColorBlue,
								},
								PlainText: "Some content",
							},
						},
					},
					"Read": notionapi.CheckboxProperty{
						Checkbox: false,
					},
				},
				Children: []notionapi.Block{
					notionapi.Heading2Block{
						BasicBlock: notionapi.BasicBlock{
							Object: notionapi.ObjectTypeBlock,
							Type:   notionapi.BlockTypeHeading2,
						},
						Heading2: notionapi.Heading{
							RichText: []notionapi.RichText{
								{
									Type: notionapi.RichTextTypeText,
									Text: &notionapi.Text{Content: "Lacinato"},
								},
							},
						},
					},
					notionapi.ParagraphBlock{
						BasicBlock: notionapi.BasicBlock{
							Object: notionapi.ObjectTypeBlock,
							Type:   notionapi.BlockTypeParagraph,
						},
						Paragraph: notionapi.Paragraph{
							RichText: []notionapi.RichText{
								{
									Text: &notionapi.Text{
										Content: "Lacinato",
										Link: &notionapi.Link{
											Url: "some_url",
										},
									},
								},
							},
							Children: nil,
						},
					},
				},
			},
			want: []byte(`{"parent":{"database_id":"some_id"},"properties":{"Link":{"url":"some_url"},"Name":{"title":[{"text":{"content":"New Media Article"}}]},"Publishing/Release Date":{"date":{"start":"2020-12-08T12:00:00Z","end":null}},"Read":{"checkbox":false},"Summary":{"text":[{"type":"text","text":{"content":"Some content"},"annotations":{"bold":true,"italic":false,"strikethrough":false,"underline":false,"code":false,"color":"blue"},"plain_text":"Some content"}]},"Type":{"select":{"id":"some_id","name":"Article","color":"default"}}},"children":[{"object":"block","type":"heading_2","heading_2":{"rich_text":[{"type":"text","text":{"content":"Lacinato"}}]}},{"object":"block","type":"paragraph","paragraph":{"rich_text":[{"text":{"content":"Lacinato","link":{"url":"some_url"}}}]}}]}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalJSON() got = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestPageUpdateRequest_MarshallJSON(t *testing.T) {
	tests := []struct {
		name    string
		req     *notionapi.PageUpdateRequest
		want    []byte
		wantErr bool
	}{
		{
			name: "update checkbox",
			req: &notionapi.PageUpdateRequest{
				Properties: map[string]notionapi.Property{
					"Checked": notionapi.CheckboxProperty{
						Checkbox: false,
					},
				},
			},
			want: []byte(`{"properties":{"Checked":{"checkbox":false}},"archived":false}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalJSON() got = %s, want %s", got, tt.want)
			}
		})
	}
}
