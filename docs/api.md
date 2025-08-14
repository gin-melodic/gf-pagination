# GF-Pagination API Documentation

Complete API reference for the GF-Pagination library - a high-performance pagination component for GoFrame framework.

## Table of Contents

- [Overview](#overview)
- [Quick Start](#quick-start)
- [API Reference](#api-reference)
  - [Core Types](#core-types)
  - [Paginator Interface](#paginator-interface)
  - [Configuration](#configuration)
  - [Query Operators](#query-operators)
- [Examples](#examples)
- [Best Practices](#best-practices)
- [Error Handling](#error-handling)

## Overview

GF-Pagination provides a standardized way to handle pagination in GoFrame applications with support for complex queries, joins, and customizable configurations.

### Key Features

- **Type-safe pagination**: Complete type definitions with validation
- **Flexible query conditions**: Support for various operators (LIKE, IN, BETWEEN, etc.)
- **Join support**: Complex table relationships
- **Configurable defaults**: Customizable page size, sort order, etc.
- **Performance optimized**: Efficient query execution with minimal overhead

## Quick Start

### Installation

```bash
go get github.com/gin-melodic/gf-pagination
```

### Basic Usage

```go
import "github.com/gin-melodic/gf-pagination"

// Create paginator
paginator := pagination.New()

// Execute pagination query
result, err := paginator.Paginate(ctx, req, query)
```

## API Reference

### Core Types

#### Request

The `Request` struct defines pagination parameters.

```go
type Request struct {
    Page     int    `json:"page" v:"min:1#Page number must be at least 1"`
    PageSize int    `json:"page_size" v:"min:1|max:500#Page size must be between 1 and 500"`
    OrderBy  string `json:"order_by" dc:"Sort field"`
    Sort     string `json:"sort" v:"in:asc,desc#Sort order must be either asc or desc"`
}
```

**Fields:**
- `Page` (int): Page number, starts from 1
- `PageSize` (int): Items per page, range 1-500
- `OrderBy` (string): Field to sort by
- `Sort` (string): Sort direction, "asc" or "desc"

**Methods:**
- `SetDefaults(config *Config) *Request`: Apply default values based on configuration
- `GetOffset() int`: Calculate offset for database query

#### Response

The `Response` struct contains pagination results and metadata.

```go
type Response struct {
    List       interface{} `json:"list"`
    Page       int         `json:"page"`
    PageSize   int         `json:"page_size"`
    Total      int64       `json:"total"`
    TotalPages int64       `json:"total_pages"`
    HasNext    bool        `json:"has_next"`
    HasPrev    bool        `json:"has_prev"`
}
```

**Fields:**
- `List` (interface{}): Array of data records
- `Page` (int): Current page number
- `PageSize` (int): Items per page
- `Total` (int64): Total number of records
- `TotalPages` (int64): Total number of pages
- `HasNext` (bool): Whether there is a next page
- `HasPrev` (bool): Whether there is a previous page

#### Query

The `Query` struct defines query configuration.

```go
type Query struct {
    Model      *gdb.Model                 `json:"-"`
    Table      string                     `json:"table"`
    Fields     string                     `json:"fields"`
    Conditions map[string]interface{}     `json:"conditions"`
    Joins      []string                   `json:"joins"`
    GroupBy    string                     `json:"group_by"`
    Having     string                     `json:"having"`
}
```

**Fields:**
- `Model` (*gdb.Model): GoFrame database model instance
- `Table` (string): Table name for query
- `Fields` (string): SELECT fields specification
- `Conditions` (map[string]interface{}): WHERE conditions
- `Joins` ([]string): JOIN clauses
- `GroupBy` (string): GROUP BY clause
- `Having` (string): HAVING clause

### Paginator Interface

The main interface for pagination operations.

```go
type Paginator interface {
    Paginate(ctx context.Context, req *Request, query *Query) (*Response, error)
    BuildQuery(ctx context.Context, model *gdb.Model, conditions map[string]interface{}) *gdb.Model
}
```

#### Methods

##### Paginate

Execute pagination query with the given parameters.

```go
func Paginate(ctx context.Context, req *Request, query *Query) (*Response, error)
```

**Parameters:**
- `ctx` (context.Context): Request context
- `req` (*Request): Pagination request parameters
- `query` (*Query): Query configuration

**Returns:**
- `*Response`: Pagination results with metadata
- `error`: Error if query fails

**Example:**
```go
result, err := paginator.Paginate(ctx, &pagination.Request{
    Page:     1,
    PageSize: 20,
    OrderBy:  "created_at",
    Sort:     "desc",
}, &pagination.Query{
    Table: "users",
    Fields: "id,username,email",
    Conditions: map[string]interface{}{
        "status": 1,
        "age_gte": 18,
    },
})
```

##### BuildQuery

Build database query with conditions (used internally).

```go
func BuildQuery(ctx context.Context, model *gdb.Model, conditions map[string]interface{}) *gdb.Model
```

**Parameters:**
- `ctx` (context.Context): Request context
- `model` (*gdb.Model): GoFrame database model
- `conditions` (map[string]interface{}): Query conditions

**Returns:**
- `*gdb.Model`: Modified database model with applied conditions

### Configuration

#### Config

Configuration options for pagination behavior.

```go
type Config struct {
    DefaultPage     int    `json:"default_page"`
    DefaultPageSize int    `json:"default_page_size"`
    MaxPageSize     int    `json:"max_page_size"`
    DefaultSort     string `json:"default_sort"`
    DefaultOrderBy  string `json:"default_order_by"`
}
```

#### Option Functions

Configuration can be customized using option functions:

```go
// WithDefaultPage sets the default page number
func WithDefaultPage(page int) Option

// WithDefaultPageSize sets the default page size
func WithDefaultPageSize(size int) Option

// WithMaxPageSize sets the maximum allowed page size
func WithMaxPageSize(size int) Option

// WithDefaultSort sets the default sort direction
func WithDefaultSort(sort string) Option

// WithDefaultOrderBy sets the default sort field
func WithDefaultOrderBy(orderBy string) Option
```

#### Creating Paginator with Custom Config

```go
paginator := pagination.New(
    pagination.WithDefaultPageSize(10),
    pagination.WithMaxPageSize(100),
    pagination.WithDefaultSort("desc"),
    pagination.WithDefaultOrderBy("created_at"),
)
```

#### Default Values

| Option | Default Value | Description |
|--------|---------------|-------------|
| DefaultPage | 1 | Starting page number |
| DefaultPageSize | 20 | Items per page |
| MaxPageSize | 500 | Maximum items per page |
| DefaultSort | "desc" | Sort direction |
| DefaultOrderBy | "created_at" | Sort field |

### Query Operators

GF-Pagination supports various query operators by appending suffixes to field names:

#### Comparison Operators

| Operator | Suffix | SQL Equivalent | Example |
|----------|--------|----------------|---------|
| Equal | (none) | `=` | `"status": 1` → `status = 1` |
| Like | `_like` | `LIKE` | `"name_like": "john"` → `name LIKE '%john%'` |
| Greater Than | `_gt` | `>` | `"age_gt": 18` → `age > 18` |
| Greater Than or Equal | `_gte` | `>=` | `"score_gte": 80` → `score >= 80` |
| Less Than | `_lt` | `<` | `"price_lt": 100` → `price < 100` |
| Less Than or Equal | `_lte` | `<=` | `"quantity_lte": 5` → `quantity <= 5` |
| Not Equal | `_ne` | `!=` | `"status_ne": 0` → `status != 0` |
| In | `_in` | `IN` | `"id_in": []int{1,2,3}` → `id IN (1,2,3)` |
| Not In | `_nin` | `NOT IN` | `"id_nin": []int{1,2,3}` → `id NOT IN (1,2,3)` |
| Between | `_between` | `BETWEEN` | `"age_between": []int{18,65}` → `age BETWEEN 18 AND 65` |
| Not Between | `_nbetween` | `NOT BETWEEN` | `"age_nbetween": []int{18,65}` → `age NOT BETWEEN 18 AND 65` |

**Example:**
```go
conditions := map[string]interface{}{
    "username_like": "john",
    "age_gte": 18,
    "status_in": []int{1, 2},
    "created_at_between": []string{"2023-01-01", "2023-12-31"},
}

result, err := paginator.Paginate(ctx, req, query)
```

### Using Custom Database Connection

```go
// Using specific database connection
result, err := pagination.PaginateWithDB(ctx, customDB, req, query)
```

### GoFrame Integration Example

```go
// In service layer
func (s *sUser) GetUserList(ctx context.Context, req *UserListReq) (*PaginationRes, error) {
    conditions := map[string]interface{}{
        "username_like": req.Username,
        "status": req.Status,
        "created_at_between": []interface{}{req.StartDate, req.EndDate},
    }
    
    return s.paginator.Paginate(ctx, &req.Request, &pagination.Query{
        Table: "users",
        Fields: "id,username,email,status,created_at",
        Conditions: conditions,
    })
}

// In controller layer
func (c *cUser) GetList(ctx context.Context, req *UserListReq) (res *UserListRes, err error) {
    result, err := service.User().GetUserList(ctx, req)
    if err != nil {
        return &UserListRes{
            Code: 1,
            Message: "Failed to get user list",
        }, nil
    }
    
    return &UserListRes{
        Code: 0,
        Message: "success",
        Data: result,
    }, nil
}
```

## Best Practices

### 1. Configuration Management

```go
// Create different paginators for different use cases
userPaginator := pagination.New(
    pagination.WithDefaultPageSize(20),
    pagination.WithMaxPageSize(100),
)

adminPaginator := pagination.New(
    pagination.WithDefaultPageSize(50),
    pagination.WithMaxPageSize(500),
)
```

### 2. Condition Building

```go
func buildUserConditions(req *UserSearchReq) map[string]interface{} {
    conditions := make(map[string]interface{})
    
    if req.Username != "" {
        conditions["username_like"] = req.Username
    }
    
    if req.Status > 0 {
        conditions["status"] = req.Status
    }
    
    if req.AgeMin > 0 && req.AgeMax > 0 {
        conditions["age_between"] = []interface{}{req.AgeMin, req.AgeMax}
    }
    
    return conditions
}
```

### 3. Error Handling

```go
result, err := paginator.Paginate(ctx, req, query)
if err != nil {
    g.Log().Error(ctx, "Pagination failed:", err)
    return nil, gerror.Wrap(err, "Failed to execute pagination query")
}

if result.Total == 0 {
    return &EmptyResponse{
        Code: 0,
        Message: "No data found",
        Data: &pagination.Response{
            List: []interface{}{},
            Page: req.Page,
            PageSize: req.PageSize,
            Total: 0,
            TotalPages: 0,
            HasNext: false,
            HasPrev: false,
        },
    }, nil
}
```

### 4. Performance Optimization

```go
// Use specific fields instead of SELECT *
query := &pagination.Query{
    Table: "users",
    Fields: "id,username,email", // Only select needed fields
    Conditions: conditions,
}

// Add appropriate database indexes for ORDER BY fields
// For example: CREATE INDEX idx_users_created_at ON users(created_at);

// Use LIMIT for large datasets
paginator := pagination.New(
    pagination.WithMaxPageSize(100), // Prevent oversized requests
)
```

### 5. Complex Query Handling

```go
// For complex OR conditions, use raw SQL
conditions := map[string]interface{}{
    "(username LIKE ? OR email LIKE ?)": []interface{}{
        "%" + searchTerm + "%",
        "%" + searchTerm + "%",
    },
    "status": 1,
}
```

## Error Handling

### Common Error Types

1. **Empty Table/Model Error**
   ```
   Error: "Table name or model cannot be empty"
   Solution: Ensure either Query.Table or Query.Model is provided
   ```

2. **Invalid Page Parameters**
   ```
   Error: Validation failed for page parameters
   Solution: Check Page >= 1 and PageSize within allowed range
   ```

3. **Database Connection Error**
   ```
   Error: Database connection failed
   Solution: Verify database configuration and connection
   ```

4. **Query Execution Error**
   ```
   Error: SQL syntax error
   Solution: Check table names, field names, and JOIN syntax
   ```

### Error Response Format

```go
if err != nil {
    return &Response{
        Code: 1,
        Message: "Pagination query failed: " + err.Error(),
        Data: nil,
    }, nil
}
```

***

For more examples and advanced usage, visit the [GitHub repository](https://github.com/gin-melodic/gf-pagination).

**Version**: 1.0.0  
**License**: MIT