package main

import (
	"context"

	pagination "github.com/gin-melodic/gf-pagination"
	"github.com/gogf/gf/v2/frame/g"
)

func main() {
	// Complex query example
	paginator := pagination.New()

	req := &pagination.Request{
		Page:     1,
		PageSize: 15,
		OrderBy:  "score",
		Sort:     "desc",
	}

	query := &pagination.Query{
		Table:  "articles a",
		Fields: "a.id,a.title,a.score,u.username,c.name as category_name",
		Conditions: map[string]interface{}{
			"a.status":             1,
			"a.score_gte":          80,
			"a.created_at_between": []interface{}{"2024-01-01", "2024-12-31"},
			"c.id_in":              []int{1, 2, 3},
		},
		Joins: []string{
			"users u ON a.user_id = u.id",
			"categories c ON a.category_id = c.id",
		},
		GroupBy: "a.id",
		Having:  "a.score > 85",
	}

	result, err := paginator.Paginate(context.Background(), req, query)
	if err != nil {
		g.Log().Fatal(context.Background(), err)
	}

	g.Log().Info(context.Background(), "Pagination query result:", result)
}