package tests

import (
	"context"
	"fmt"
	"testing"

	pagination "github.com/gin-melodic/gf-pagination"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/test/gtest"

	_ "github.com/gogf/gf/contrib/drivers/sqlite/v2"
)

func TestMain(m *testing.M) {
	setupTestDB()
	m.Run()
}

func setupTestDB() {
	ctx := gctx.New()

	// Solution 1: Use the correct SQLite in-memory database configuration format
	configContent := `
database:
  default:
    link: "sqlite::@file(:memory:)?cache=shared"
    debug: false
`
	adapter, err := gcfg.NewAdapterContent(configContent)
	if err != nil {
		panic(err)
	}
	g.Cfg().SetAdapter(adapter)

	db := g.DB()

	// Verify database connection
	if err := db.PingMaster(); err != nil {
		panic("database connection failed: " + err.Error())
	}

	// Create test table
	_, err = db.Exec(ctx, `
        CREATE TABLE test_users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name VARCHAR(100) NOT NULL,
            email VARCHAR(100),
            age INTEGER,
            status INTEGER DEFAULT 1,
            city VARCHAR(50),
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP
        )
    `)
	if err != nil {
		panic("failed to create test table: " + err.Error())
	}

	// Insert test data
	cities := []string{"Beijing", "Shanghai", "Guangzhou", "Shenzhen"}
	for i := 1; i <= 100; i++ {
		_, err = db.Exec(ctx,
			"INSERT INTO test_users (name, email, age, status, city) VALUES (?, ?, ?, ?, ?)",
			fmt.Sprintf("user%d", i),
			fmt.Sprintf("user%d@test.com", i),
			20+i%50,
			i%3+1,
			cities[i%4],
		)
		if err != nil {
			panic("failed to insert test data: " + err.Error())
		}
	}
}

func TestPagination_Basic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		paginator := pagination.New()

		req := &pagination.Request{
			Page:     1,
			PageSize: 10,
		}

		query := &pagination.Query{
			Table: "test_users",
		}

		result, err := paginator.Paginate(ctx, req, query)
		t.AssertNil(err)
		t.AssertNE(result, nil)
		t.Assert(result.Page, 1)
		t.Assert(result.PageSize, 10)
		t.Assert(result.Total, 100)
		t.Assert(result.TotalPages, 10)
		t.Assert(result.HasNext, true)
		t.Assert(result.HasPrev, false)
	})
}

func TestPagination_WithConditions(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		paginator := pagination.New()

		req := &pagination.Request{
			Page:     1,
			PageSize: 5,
			OrderBy:  "age",
			Sort:     "desc",
		}

		query := &pagination.Query{
			Table: "test_users",
			Conditions: map[string]interface{}{
				"status": 1,
				"age_gt": 30,
			},
		}

		result, err := paginator.Paginate(ctx, req, query)
		t.AssertNil(err)
		t.AssertNE(result, nil)
		t.Assert(result.Page, 1)
		t.Assert(result.PageSize, 5)
	})
}

func TestPagination_LikeQuery(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		paginator := pagination.New()

		req := &pagination.Request{
			Page:     1,
			PageSize: 20,
		}

		query := &pagination.Query{
			Table: "test_users",
			Conditions: map[string]interface{}{
				"name_like": "user1",
			},
		}

		result, err := paginator.Paginate(ctx, req, query)
		t.AssertNil(err)
		t.AssertNE(result, nil)
		// user1, user10, user11, ..., user19 = 11 records
		t.Assert(result.Total, 11)
	})
}

func TestPagination_InQuery(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		paginator := pagination.New()

		req := &pagination.Request{
			Page:     1,
			PageSize: 50,
		}

		query := &pagination.Query{
			Table: "test_users",
			Conditions: map[string]interface{}{
				"status_in": []int{1, 2},
			},
		}

		result, err := paginator.Paginate(ctx, req, query)
		t.AssertNil(err)
		t.AssertNE(result, nil)
		// Records with status 1 or 2
		t.AssertGT(result.Total, 0)
	})
}

func TestPagination_BetweenQuery(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		paginator := pagination.New()

		req := &pagination.Request{
			Page:     1,
			PageSize: 100,
		}

		query := &pagination.Query{
			Table: "test_users",
			Conditions: map[string]interface{}{
				"age_between": []interface{}{25, 35},
			},
		}

		result, err := paginator.Paginate(ctx, req, query)
		t.AssertNil(err)
		t.AssertNE(result, nil)
		t.AssertGT(result.Total, 0)
	})
}

func TestPagination_EmptyResult(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		paginator := pagination.New()

		req := &pagination.Request{
			Page:     1,
			PageSize: 10,
		}

		query := &pagination.Query{
			Table: "test_users",
			Conditions: map[string]interface{}{
				"name": "nonexistent_user",
			},
		}

		result, err := paginator.Paginate(ctx, req, query)
		t.AssertNil(err)
		t.AssertNE(result, nil)
		t.Assert(result.Total, 0)
		t.Assert(result.TotalPages, 0)
		t.Assert(result.HasNext, false)
		t.Assert(result.HasPrev, false)
	})
}

func TestPagination_DefaultValues(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		paginator := pagination.New()

		// Test empty pagination request
		req := &pagination.Request{}

		query := &pagination.Query{
			Table: "test_users",
		}

		result, err := paginator.Paginate(ctx, req, query)
		t.AssertNil(err)
		t.AssertNE(result, nil)
		t.Assert(result.Page, 1)      // Default page 1
		t.Assert(result.PageSize, 20) // Default 20 items per page
	})
}

func TestPagination_CustomConfig(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()

		// Use custom configuration
		paginator := pagination.New(
			pagination.WithDefaultPageSize(15),
			pagination.WithMaxPageSize(30),
			pagination.WithDefaultSort("asc"),
		)

		req := &pagination.Request{
			PageSize: 50, // Exceeds maximum limit
		}

		query := &pagination.Query{
			Table: "test_users",
		}

		result, err := paginator.Paginate(ctx, req, query)
		t.AssertNil(err)
		t.AssertNE(result, nil)
		t.Assert(result.PageSize, 30) // Limited to maximum value
	})
}

func TestPagination_InvalidTable(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		paginator := pagination.New()

		req := &pagination.Request{
			Page:     1,
			PageSize: 10,
		}

		query := &pagination.Query{
			// Empty table name
		}

		_, err := paginator.Paginate(ctx, req, query)
		t.AssertNE(err, nil)
		t.AssertIN(err.Error(), "table name or model cannot be empty")
	})
}
