package note

// Filters for listing notes.
type Filters struct {
	Status     *NoteStatus
	TemplateID *string
	OwnerID    *string
	Query      *string
}

type SectionWithField struct {
	Section    Section
	FieldLabel string
	FieldOrder int
	IsRequired bool
}

type WithMeta struct {
	Note          Note
	TemplateName  string
	OwnerFullName string
	Sections      []SectionWithField
}
