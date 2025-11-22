package template

// Filters for listing templates.
type Filters struct {
	Query   *string
	OwnerID *string
}

type WithUsage struct {
	Template Template
	Fields   []Field
	IsUsed   bool
}
