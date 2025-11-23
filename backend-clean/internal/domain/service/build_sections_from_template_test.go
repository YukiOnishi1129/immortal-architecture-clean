package service

import (
	"testing"

	"immortal-architecture-clean/backend/internal/domain/template"
)

func TestBuildSectionsFromTemplate(t *testing.T) {
	tests := []struct {
		name   string
		fields []template.Field
	}{
		{
			name: "[Success] preserves field order and IDs",
			fields: []template.Field{
				{ID: "f1", Label: "Title", Order: 1},
				{ID: "f2", Label: "Body", Order: 2},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sections := BuildSectionsFromTemplate(tt.fields)
			if len(sections) != len(tt.fields) {
				t.Fatalf("sections length mismatch: %d", len(sections))
			}
			for i, f := range tt.fields {
				if sections[i].FieldID != f.ID {
					t.Fatalf("field ID mismatch at %d: want %s, got %s", i, f.ID, sections[i].FieldID)
				}
			}
		})
	}
}
