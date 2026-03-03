package notionapi

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type DataSourceID string

func (id DataSourceID) String() string {
	return string(id)
}

// DataSourceService provides methods for interacting with Notion data sources.
// In API v2025-09-03, databases are containers with one or more data sources,
// each with independent schemas. Most operations that previously targeted
// /v1/databases/:id now target /v1/data_sources/:id.
type DataSourceService interface {
	Get(ctx context.Context, id DataSourceID) (*DataSource, error)
	Create(ctx context.Context, request *DataSourceCreateRequest) (*DataSource, error)
	Update(ctx context.Context, id DataSourceID, request *DataSourceUpdateRequest) (*DataSource, error)
	Query(ctx context.Context, id DataSourceID, request *DatabaseQueryRequest) (*DatabaseQueryResponse, error)
	ListTemplates(ctx context.Context, id DataSourceID, pagination *Pagination) (*TemplateListResponse, error)
}

// DataSource represents a data source within a database.
type DataSource struct {
	Object         ObjectType      `json:"object"`
	ID             ObjectID        `json:"id"`
	Title          []RichText      `json:"title,omitempty"`
	Properties     PropertyConfigs `json:"properties,omitempty"`
	Parent         Parent          `json:"parent"`
	DatabaseParent Parent          `json:"database_parent,omitempty"`
	Description    []RichText      `json:"description,omitempty"`
	Icon           *Icon           `json:"icon,omitempty"`
	Cover          *Image          `json:"cover,omitempty"`
	IsInline       bool            `json:"is_inline,omitempty"`
	Archived       bool            `json:"archived,omitempty"`
	InTrash        bool            `json:"in_trash,omitempty"`
	CreatedTime    time.Time       `json:"created_time,omitempty"`
	LastEditedTime time.Time       `json:"last_edited_time,omitempty"`
	CreatedBy      User            `json:"created_by,omitempty"`
	LastEditedBy   User            `json:"last_edited_by,omitempty"`
	URL            string          `json:"url,omitempty"`
	PublicURL      string          `json:"public_url,omitempty"`
}

func (ds *DataSource) GetObject() ObjectType {
	return ds.Object
}

// DataSourceCreateRequest represents the request body for DataSourceClient.Create.
type DataSourceCreateRequest struct {
	Parent     Parent          `json:"parent"`
	Properties PropertyConfigs `json:"properties"`
	Title      []RichText      `json:"title,omitempty"`
	Icon       *Icon           `json:"icon,omitempty"`
}

// DataSourceUpdateRequest represents the request body for DataSourceClient.Update.
type DataSourceUpdateRequest struct {
	Title       []RichText      `json:"title,omitempty"`
	Properties  PropertyConfigs `json:"properties,omitempty"`
	Description []RichText      `json:"description,omitempty"`
	Icon        *Icon           `json:"icon,omitempty"`
	InTrash     *bool           `json:"in_trash,omitempty"`
	Archived    *bool           `json:"archived,omitempty"`
}

// TemplateRef is a reference to a template within a data source.
type TemplateRef struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	IsDefault bool   `json:"is_default"`
}

// TemplateListResponse is the response from DataSourceClient.ListTemplates.
type TemplateListResponse struct {
	Templates  []TemplateRef `json:"templates"`
	HasMore    bool          `json:"has_more"`
	NextCursor Cursor        `json:"next_cursor"`
}

// DataSourceClient implements DataSourceService.
type DataSourceClient struct {
	apiClient *Client
}

// Get retrieves a data source by ID.
//
// See https://developers.notion.com/reference/retrieve-a-data-source
func (dc *DataSourceClient) Get(ctx context.Context, id DataSourceID) (*DataSource, error) {
	return doRequest[DataSource](dc.apiClient, ctx, http.MethodGet, fmt.Sprintf("data_sources/%s", id.String()), nil, nil)
}

// Update modifies a data source's schema or metadata.
//
// See https://developers.notion.com/reference/update-a-data-source
func (dc *DataSourceClient) Update(ctx context.Context, id DataSourceID, requestBody *DataSourceUpdateRequest) (*DataSource, error) {
	return doRequest[DataSource](dc.apiClient, ctx, http.MethodPatch, fmt.Sprintf("data_sources/%s", id.String()), nil, requestBody)
}

// Create creates a new data source within a database.
//
// See https://developers.notion.com/reference/create-a-data-source
func (dc *DataSourceClient) Create(ctx context.Context, requestBody *DataSourceCreateRequest) (*DataSource, error) {
	return doRequest[DataSource](dc.apiClient, ctx, http.MethodPost, "data_sources", nil, requestBody)
}

// Query returns pages from a data source, with optional filtering and sorting.
//
// See https://developers.notion.com/reference/post-data-source-query
func (dc *DataSourceClient) Query(ctx context.Context, id DataSourceID, requestBody *DatabaseQueryRequest) (*DatabaseQueryResponse, error) {
	return doRequest[DatabaseQueryResponse](dc.apiClient, ctx, http.MethodPost, fmt.Sprintf("data_sources/%s/query", id.String()), nil, requestBody)
}

// ListTemplates returns templates available for a data source.
//
// See https://developers.notion.com/reference/list-data-source-templates
func (dc *DataSourceClient) ListTemplates(ctx context.Context, id DataSourceID, pagination *Pagination) (*TemplateListResponse, error) {
	return doRequest[TemplateListResponse](dc.apiClient, ctx, http.MethodGet, fmt.Sprintf("data_sources/%s/templates", id.String()), pagination.ToQuery(), nil)
}
