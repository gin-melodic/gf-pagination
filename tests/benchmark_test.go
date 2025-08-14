package tests

import (
	"context"
	"testing"

	pagination "github.com/gin-melodic/gf-pagination"
)

func BenchmarkPagination_BasicQuery(b *testing.B) {
	ctx := context.Background()
	paginator := pagination.New()

	req := &pagination.Request{
		Page:     1,
		PageSize: 20,
	}

	query := &pagination.Query{
		Table: "test_users",
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := paginator.Paginate(ctx, req, query)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkPagination_WithConditions(b *testing.B) {
	ctx := context.Background()
	paginator := pagination.New()

	req := &pagination.Request{
		Page:     1,
		PageSize: 20,
		OrderBy:  "age",
		Sort:     "desc",
	}

	query := &pagination.Query{
		Table: "test_users",
		Conditions: map[string]interface{}{
			"status":    1,
			"age_gt":    25,
			"name_like": "user",
		},
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := paginator.Paginate(ctx, req, query)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkPagination_ComplexQuery(b *testing.B) {
	ctx := context.Background()
	paginator := pagination.New()

	req := &pagination.Request{
		Page:     1,
		PageSize: 10,
		OrderBy:  "created_at",
		Sort:     "desc",
	}

	query := &pagination.Query{
		Table:  "test_users",
		Fields: "id,name,email,age,status,city",
		Conditions: map[string]interface{}{
			"status_in":   []int{1, 2},
			"age_between": []interface{}{20, 50},
			"city_in":     []string{"北京", "上海"},
			"name_like":   "user",
		},
		GroupBy: "city",
		Having:  "COUNT(*) > 0",
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := paginator.Paginate(ctx, req, query)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkPagination_LargePageSize(b *testing.B) {
	ctx := context.Background()
	paginator := pagination.New(
		pagination.WithMaxPageSize(1000),
	)

	req := &pagination.Request{
		Page:     1,
		PageSize: 100,
	}

	query := &pagination.Query{
		Table: "test_users",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := paginator.Paginate(ctx, req, query)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkPagination_DeepPaging(b *testing.B) {
	ctx := context.Background()
	paginator := pagination.New()

	req := &pagination.Request{
		Page:     10, // 深度分页
		PageSize: 10,
		OrderBy:  "id",
		Sort:     "asc",
	}

	query := &pagination.Query{
		Table: "test_users",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := paginator.Paginate(ctx, req, query)
		if err != nil {
			b.Fatal(err)
		}
	}
}
