package pagination

// Config Pagination configuration
type Config struct {
	DefaultPage     int    `json:"default_page"`      // Default page number
	DefaultPageSize int    `json:"default_page_size"` // Default page size
	MaxPageSize     int    `json:"max_page_size"`     // Maximum page size
	DefaultSort     string `json:"default_sort"`      // Default sort method
	DefaultOrderBy  string `json:"default_order_by"`  // Default order by field
}

// DefaultConfig Return default configuration
func DefaultConfig() *Config {
	return &Config{
		DefaultPage:     1,
		DefaultPageSize: 20,
		MaxPageSize:     500,
		DefaultSort:     "desc",
		DefaultOrderBy:  "created_at",
	}
}

// NewConfig Create new configuration
func NewConfig(opts ...Option) *Config {
	config := DefaultConfig()
	for _, opt := range opts {
		opt(config)
	}
	return config
}

// Option Configuration option function
type Option func(*Config)

// WithDefaultPage Set default page number
func WithDefaultPage(page int) Option {
	return func(c *Config) {
		c.DefaultPage = page
	}
}

// WithDefaultPageSize Set default page size
func WithDefaultPageSize(size int) Option {
	return func(c *Config) {
		c.DefaultPageSize = size
	}
}

// WithMaxPageSize Set maximum page size
func WithMaxPageSize(size int) Option {
	return func(c *Config) {
		c.MaxPageSize = size
	}
}

// WithDefaultSort Set default sort method
func WithDefaultSort(sort string) Option {
	return func(c *Config) {
		c.DefaultSort = sort
	}
}

// WithDefaultOrderBy Set default order by field
func WithDefaultOrderBy(orderBy string) Option {
	return func(c *Config) {
		c.DefaultOrderBy = orderBy
	}
}
