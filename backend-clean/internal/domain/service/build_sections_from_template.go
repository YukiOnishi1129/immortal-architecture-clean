// Package service contains domain services.
// Package service contains domain services shared across aggregates.
package service

import (
	"immortal-architecture-clean/backend/internal/domain/note"
	"immortal-architecture-clean/backend/internal/domain/template"
)

// BuildSectionsFromTemplate creates empty note sections based on template fields.
func BuildSectionsFromTemplate(fields []template.Field) []note.Section {
	sections := make([]note.Section, 0, len(fields))
	for _, f := range fields {
		sections = append(sections, note.Section{
			FieldID: f.ID,
			NoteID:  "",
			Content: "",
		})
	}
	return sections
}
