// Package template holds template domain models.
package template

import domainerr "immortal-architecture-clean/backend/internal/domain/errors"

// ReplaceFields enforces updates through the aggregate root.
func (t *Template) ReplaceFields(fields []Field) error {
	fields, err := NormalizeAndValidate(fields)
	if err != nil {
		return err
	}
	t.Fields = fields
	return nil
}

// EnsureOwner set owner if empty.
func (t *Template) EnsureOwner(ownerID string) error {
	if ownerID == "" {
		return domainerr.ErrTemplateOwnerRequired
	}
	t.OwnerID = ownerID
	return nil
}
