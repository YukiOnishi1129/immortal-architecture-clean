package factory

import httppresenter "immortal-architecture-clean/backend/internal/adapter/http/presenter"

func NewAccountOutputFactory() func() *httppresenter.AccountPresenter {
	return func() *httppresenter.AccountPresenter {
		return httppresenter.NewAccountPresenter()
	}
}

func NewTemplateOutputFactory() func() *httppresenter.TemplatePresenter {
	return func() *httppresenter.TemplatePresenter {
		return httppresenter.NewTemplatePresenter()
	}
}

func NewNoteOutputFactory() func() *httppresenter.NotePresenter {
	return func() *httppresenter.NotePresenter {
		return httppresenter.NewNotePresenter()
	}
}
