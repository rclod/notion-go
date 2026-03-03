# notion-go

A Go client library for the [Notion API](https://developers.notion.com/).

Fork of [jomei/notionapi](https://github.com/jomei/notionapi) with serialization bug fixes and support for Notion API version `2025-09-03` (multi-source databases).

## Changes from upstream

- **Serialization fixes**: Added `omitempty` to `StatusConfig`, `StatusPropertyConfig.ID`, `NumberFormat.Format`, and `SearchFilter` fields that caused invalid JSON when empty
- **API v2025-09-03**: Data source support (`DataSourceService` with Get/Update/Query), `DataSourceRef` on `Database`, `InitialDataSource` on `DatabaseCreateRequest`, `DataSourceID` on `Parent`, new object/parent type constants
- **Module path**: `github.com/rclod/notion-go`

## Installation

```bash
go get github.com/rclod/notion-go
```

## Usage

First, follow the [Getting Started Guide](https://developers.notion.com/docs/getting-started) to obtain an integration token.

### Initialization

```go
import notionapi "github.com/rclod/notion-go"

client := notionapi.NewClient("your_integration_token")
```

### Pages and Databases

```go
page, err := client.Page.Get(context.Background(), "your_page_id")

db, err := client.Database.Get(context.Background(), "your_database_id")
```

### Data Sources (API v2025-09-03)

In the latest API version, databases are containers with one or more data sources. Schema operations (properties) target data sources rather than databases directly.

```go
// Get a data source (schema + properties)
ds, err := client.DataSource.Get(context.Background(), "your_data_source_id")

// Update data source schema
_, err = client.DataSource.Update(context.Background(), "your_data_source_id", &notionapi.DataSourceUpdateRequest{
    Properties: notionapi.PropertyConfigs{
        "Status": notionapi.SelectPropertyConfig{
            Type: notionapi.PropertyConfigTypeSelect,
            Select: notionapi.Select{
                Options: []notionapi.Option{
                    {Name: "To Do", Color: "red"},
                    {Name: "Done", Color: "green"},
                },
            },
        },
    },
})

// Query pages in a data source
results, err := client.DataSource.Query(context.Background(), "your_data_source_id", &notionapi.DatabaseQueryRequest{
    Filter: &notionapi.PropertyFilter{
        Property: "Status",
        Select:   &notionapi.SelectFilterCondition{Equals: "Done"},
    },
})
```

### Creating a Database with Initial Data Source

```go
db, err := client.Database.Create(context.Background(), &notionapi.DatabaseCreateRequest{
    Parent: notionapi.Parent{
        Type:   notionapi.ParentTypePageID,
        PageID: "parent_page_id",
    },
    Title: []notionapi.RichText{
        {Type: notionapi.RichTextTypeText, Text: &notionapi.Text{Content: "My Database"}},
    },
    InitialDataSource: &notionapi.InitialDataSource{
        Properties: notionapi.PropertyConfigs{
            "Name": notionapi.TitlePropertyConfig{Type: notionapi.PropertyConfigTypeTitle},
        },
    },
})
```
