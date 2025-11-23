package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"

	domainerr "immortal-architecture-clean/backend/internal/domain/errors"
	"immortal-architecture-clean/backend/internal/domain/template"
	"immortal-architecture-clean/backend/internal/port"
	uc "immortal-architecture-clean/backend/internal/usecase"
	mockusecase "immortal-architecture-clean/backend/internal/usecase/mock"
)

func TestTemplateInteractor_Create(t *testing.T) {
	tests := []struct {
		name       string
		input      port.TemplateCreateInput
		created    *template.Template
		withFields *template.WithUsage
		createErr  error
		wantError  error
	}{
		{
			name: "[Success] create with fields",
			input: port.TemplateCreateInput{
				Name:    "Template",
				OwnerID: "owner-1",
				Fields: []template.Field{
					{ID: "f1", Label: "Title", Order: 1, IsRequired: true},
				},
			},
			created: &template.Template{ID: "tpl-1", Name: "Template", OwnerID: "owner-1"},
			withFields: &template.WithUsage{
				Template: template.Template{
					ID:      "tpl-1",
					Name:    "Template",
					OwnerID: "owner-1",
					Fields:  []template.Field{{ID: "f1", Label: "Title", Order: 1, IsRequired: true}},
				},
				Fields: []template.Field{{ID: "f1", Label: "Title", Order: 1, IsRequired: true}},
			},
		},
		{
			name: "[Fail] validation error",
			input: port.TemplateCreateInput{
				Name:    "",
				OwnerID: "owner-1",
				Fields:  []template.Field{{ID: "f1", Label: "Title", Order: 1}},
			},
			wantError: domainerr.ErrTemplateNameRequired,
		},
		{
			name: "[Fail] repo create error",
			input: port.TemplateCreateInput{
				Name:    "Template",
				OwnerID: "owner-1",
				Fields:  []template.Field{{ID: "f1", Label: "Title", Order: 1}},
			},
			createErr: errors.New("repo error"),
			wantError: errors.New("repo error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mockusecase.NewMockTemplateRepository(ctrl)
			tx := mockusecase.NewMockTxManager(ctrl)
			out := mockusecase.NewMockTemplateOutputPort(ctrl)

			// set expectations based on test case data
			if tt.created != nil || tt.createErr != nil {
				tx.EXPECT().WithinTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(_ context.Context, fn func(context.Context) error) error {
						return fn(context.Background())
					},
				)
			}

			if tt.created != nil || tt.createErr != nil {
				repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(tt.created, tt.createErr)
			}
			if tt.created != nil && tt.createErr == nil {
				repo.EXPECT().ReplaceFields(gomock.Any(), tt.created.ID, gomock.Any()).Return(nil)
				repo.EXPECT().Get(gomock.Any(), tt.created.ID).Return(tt.withFields, nil)
				out.EXPECT().PresentTemplate(gomock.Any(), tt.withFields).Return(nil)
			}

			interactor := uc.NewTemplateInteractor(repo, tx, out)
			err := interactor.Create(context.Background(), tt.input)

			if tt.wantError == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantError != nil {
				if tt.wantError.Error() != err.Error() {
					t.Fatalf("want %v, got %v", tt.wantError, err)
				}
			}
		})
	}
}
