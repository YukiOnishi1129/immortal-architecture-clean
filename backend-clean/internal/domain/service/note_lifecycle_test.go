package service

import (
	"errors"
	"testing"

	domainerr "immortal-architecture-clean/backend/internal/domain/errors"
	"immortal-architecture-clean/backend/internal/domain/note"
	"immortal-architecture-clean/backend/internal/domain/template"
)

func TestBuildNote(t *testing.T) {
	tpl := template.Template{
		ID:      "tpl-1",
		Name:    "Template",
		OwnerID: "owner-1",
		Fields: []template.Field{
			{ID: "f1", Label: "Title", Order: 1, IsRequired: true},
		},
	}
	sections := []note.Section{{FieldID: "f1", Content: "value"}}

	tests := []struct {
		name      string
		title     string
		tpl       template.Template
		ownerID   string
		sections  []note.Section
		wantError error
	}{
		{
			name:     "[Success] builds draft note",
			title:    "Note title",
			tpl:      tpl,
			ownerID:  "owner-1",
			sections: sections,
		},
		{
			name:      "[Fail] missing owner",
			title:     "Note title",
			tpl:       tpl,
			ownerID:   "",
			sections:  sections,
			wantError: domainerr.ErrOwnerRequired,
		},
		{
			name:      "[Fail] missing required section",
			title:     "Note title",
			tpl:       tpl,
			ownerID:   "owner-1",
			sections:  []note.Section{{FieldID: "f1", Content: ""}},
			wantError: domainerr.ErrRequiredFieldEmpty,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n, err := BuildNote(tt.title, tt.tpl, tt.ownerID, tt.sections)
			if tt.wantError == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantError != nil && !errors.Is(err, tt.wantError) {
				t.Fatalf("want %v, got %v", tt.wantError, err)
			}
			if err == nil {
				if n.Status != note.StatusDraft || n.OwnerID != tt.ownerID || n.TemplateID != tt.tpl.ID {
					t.Fatalf("unexpected note: %+v", n)
				}
			}
		})
	}
}
