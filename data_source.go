package notionapi

import (
	"context"
	"fmt"
	"net/http"
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
	Update(ctx context.Context, id DataSourceID, request *DataSourceUpdateRequest) (*DataSource, error)
	Query(ctx context.Context, id DataSourceID, request *DatabaseQueryRequest) (*DatabaseQueryResponse, error)
}

// DataSource represents a data source within a database.
type DataSource struct {
	Object         ObjectType      `json:"object"`
	ID             ObjectID        `json:"id"`
	Title          []RichText      `json:"title,omitempty"`
	Properties     PropertyConfigs `json:"properties,omitempty"`
	Parent         Parent          `json:"parent"`
	Description    []RichText      `json:"description,omitempty"`
	Icon           *Icon           `json:"icon,omitempty"`
	Cover          *Image          `json:"cover,omitempty"`
}

func (ds *DataSource) GetObject() ObjectType {
	return ds.Object
}

// DataSourceUpdateRequest represents the request body for DataSourceClient.Update.
type DataSourceUpdateRequest struct {
	Title       []RichText      `json:"title,omitempty"`
	Properties  PropertyConfigs `json:"properties,omitempty"`
	Description []RichText      `json:"description,omitempty"`
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

// Query returns pages from a data source, with optional filtering and sorting.
//
// See https://developers.notion.com/reference/post-data-source-query
func (dc *DataSourceClient) Query(ctx context.Context, id DataSourceID, requestBody *DatabaseQueryRequest) (*DatabaseQueryResponse, error) {
	return doRequest[DatabaseQueryResponse](dc.apiClient, ctx, http.MethodPost, fmt.Sprintf("data_sources/%s/query", id.String()), nil, requestBody)
}
