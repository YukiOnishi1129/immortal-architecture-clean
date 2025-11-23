package factory

import httppresenter "immortal-architecture-clean/backend/internal/adapter/http/presenter"

// NewAccountOutputFactory returns a factory for AccountPresenter.
func NewAccountOutputFactory() func() *httppresenter.AccountPresenter {
	return func() *httppresenter.AccountPresenter {
		return httppresenter.NewAccountPresenter()
	}
}

// NewTemplateOutputFactory returns a factory for TemplatePresenter.
func NewTemplateOutputFactory() func() *httppresenter.TemplatePresenter {
	return func() *httppresenter.TemplatePresenter {
		return httppresenter.NewTemplatePresenter()
	}
}

// NewNoteOutputFactory returns a factory for NotePresenter.
func NewNoteOutputFactory() func() *httppresenter.NotePresenter {
	return func() *httppresenter.NotePresenter {
		return httppresenter.NewNotePresenter()
	}
}
