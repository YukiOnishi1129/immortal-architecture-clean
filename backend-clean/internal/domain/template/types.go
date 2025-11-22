// Package template holds template domain models.
package template

// Filters for listing templates.
type Filters struct {
	Query   *string
	OwnerID *string
}

// WithUsage is a template with usage metadata.
type WithUsage struct {
	Template Template
	Fields   []Field
	IsUsed   bool
}
