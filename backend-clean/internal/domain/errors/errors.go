package errors

import "errors"

var (
	ErrNotFound                = errors.New("not found")
	ErrTemplateInUse           = errors.New("template is used by notes")
	ErrUnauthorized            = errors.New("unauthorized")
	ErrInvalidStatus           = errors.New("invalid status")
	ErrInvalidStatusChange     = errors.New("invalid status change")
	ErrInvalidTemplateField    = errors.New("invalid template field")
	ErrTemplateNameRequired    = errors.New("template name is required")
	ErrTemplateOwnerRequired   = errors.New("template owner is required")
	ErrFieldRequired           = errors.New("template requires at least one field")
	ErrFieldOrderInvalid       = errors.New("field order must be greater than zero and unique")
	ErrFieldLabelRequired      = errors.New("field label is required")
	ErrSectionsMissing         = errors.New("sections do not match template fields")
	ErrRequiredFieldEmpty      = errors.New("required field content is empty")
	ErrProviderRequired        = errors.New("provider is required")
	ErrProviderAccountRequired = errors.New("provider account id is required")
	ErrTitleRequired           = errors.New("title is required")
	ErrOwnerRequired           = errors.New("owner is required")
)
