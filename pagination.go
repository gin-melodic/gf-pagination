package pagination

import (
	"context"
	"math"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

// Paginator 分页器接口
type Paginator interface {
	Paginate(ctx context.Context, req *Request, query *Query) (*Response, error)
	BuildQuery(ctx context.Context, model *gdb.Model, conditions map[string]interface{}) *gdb.Model
}

// paginator 分页器实现
type paginator struct {
	config *Config
}

// New 创建新的分页器实例
func New(opts ...Option) Paginator {
	return &paginator{
		config: NewConfig(opts...),
	}
}

// Paginate 执行分页查询
func (p *paginator) Paginate(ctx context.Context, req *Request, query *Query) (*Response, error) {
	// 参数验证和默认值设置
	req = req.SetDefaults(p.config)

	var (
		db    = g.DB()
		model *gdb.Model
		err   error
	)

	// 构建基础查询模型
	switch {
	case query.Model != nil:
		model = query.Model
	case query.Table != "":
		model = db.Model(query.Table)
	default:
		return nil, gerror.New("table name or model cannot be empty")
	}

	// 设置查询字段
	if query.Fields != "" {
		model = model.Fields(query.Fields)
	}

	// 构建查询条件
	model = p.BuildQuery(ctx, model, query.Conditions)

	// 处理关联查询
	for _, join := range query.Joins {
		model = model.LeftJoin(join)
	}

	// 处理分组和Having
	switch {
	case query.GroupBy != "":
		model = model.Group(query.GroupBy)
	case query.Having != "":
		model = model.Having(query.Having)
	}

	// 获取总记录数
	total, err := model.Clone().Count()
	if err != nil {
		return nil, gerror.Wrap(err, "failed to get total records")
	}

	// 构建排序
	if req.OrderBy != "" {
		orderStr := req.OrderBy
		if req.Sort != "" {
			orderStr += " " + strings.ToUpper(req.Sort)
		}
		model = model.Order(orderStr)
	}

	// 执行分页查询
	var list []gdb.Record
	err = model.Limit(req.PageSize).Offset(req.GetOffset()).Scan(&list)
	if err != nil {
		return nil, gerror.Wrap(err, "pagination query failed")
	}

	// 计算分页信息
	totalPages := int64(math.Ceil(float64(total) / float64(req.PageSize)))

	return &Response{
		List:       list,
		Page:       req.Page,
		PageSize:   req.PageSize,
		Total:      int64(total),
		TotalPages: totalPages,
		HasNext:    int64(req.Page) < totalPages,
		HasPrev:    req.Page > 1,
	}, nil
}

// queryType represents the type of query condition
type queryType int

const (
	queryEqual queryType = iota
	queryLike
	queryIn
	queryNotIn
	queryGt
	queryGte
	queryLt
	queryLte
	queryBetween
	queryNull
)

// getQueryType determines the query type based on the field suffix
func getQueryType(field string) (string, queryType) {
	switch {
	case strings.HasSuffix(field, "_like"):
		return strings.TrimSuffix(field, "_like"), queryLike
	case strings.HasSuffix(field, "_in"):
		return strings.TrimSuffix(field, "_in"), queryIn
	case strings.HasSuffix(field, "_not_in"):
		return strings.TrimSuffix(field, "_not_in"), queryNotIn
	case strings.HasSuffix(field, "_gt"):
		return strings.TrimSuffix(field, "_gt"), queryGt
	case strings.HasSuffix(field, "_gte"):
		return strings.TrimSuffix(field, "_gte"), queryGte
	case strings.HasSuffix(field, "_lt"):
		return strings.TrimSuffix(field, "_lt"), queryLt
	case strings.HasSuffix(field, "_lte"):
		return strings.TrimSuffix(field, "_lte"), queryLte
	case strings.HasSuffix(field, "_between"):
		return strings.TrimSuffix(field, "_between"), queryBetween
	case strings.HasSuffix(field, "_null"):
		return strings.TrimSuffix(field, "_null"), queryNull
	default:
		return field, queryEqual
	}
}

// BuildQuery 构建查询条件
func (p *paginator) BuildQuery(ctx context.Context, model *gdb.Model, conditions map[string]interface{}) *gdb.Model {
	if conditions == nil {
		return model
	}

	for field, value := range conditions {
		if value == nil || value == "" {
			continue
		}

		realField, qType := getQueryType(field)
		model = p.applyQueryType(model, realField, value, qType)
	}

	return model
}

// applyQueryType applies the appropriate query based on the query type
func (p *paginator) applyQueryType(model *gdb.Model, field string, value interface{}, qType queryType) *gdb.Model {
	switch qType {
	case queryLike:
		// 模糊查询
		return model.Where(field+" LIKE ?", "%"+gconv.String(value)+"%")

	case queryIn:
		// IN查询
		return model.WhereIn(field, value)

	case queryNotIn:
		// NOT IN查询
		return model.WhereNotIn(field, value)

	case queryGt:
		// 大于
		return model.Where(field+" > ?", value)

	case queryGte:
		// 大于等于
		return model.Where(field+" >= ?", value)

	case queryLt:
		// 小于
		return model.Where(field+" < ?", value)

	case queryLte:
		// 小于等于
		return model.Where(field+" <= ?", value)

	case queryBetween:
		// BETWEEN查询
		if arr, ok := value.([]interface{}); ok && len(arr) == 2 {
			return model.Where(field+" BETWEEN ? AND ?", arr[0], arr[1])
		}
		return model

	case queryNull:
		// NULL查询
		if gconv.Bool(value) {
			return model.Where(field + " IS NULL")
		}
		return model.Where(field + " IS NOT NULL")

	default:
		// 等值查询
		return model.Where(field, value)
	}
}

// PaginateWithDB 使用指定数据库连接进行分页查询
func PaginateWithDB(ctx context.Context, db gdb.DB, req *Request, query *Query, opts ...Option) (*Response, error) {
	p := &paginator{config: NewConfig(opts...)}

	// 重写模型构建逻辑使用指定数据库
	if query.Model == nil && query.Table != "" {
		query.Model = db.Model(query.Table)
	}

	return p.Paginate(ctx, req, query)
}
