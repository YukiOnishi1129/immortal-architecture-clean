// Code generated manually for gomock-based tests (account/note).
package mockusecase

import (
	"context"
	"reflect"

	"github.com/golang/mock/gomock"

	"immortal-architecture-clean/backend/internal/domain/account"
	"immortal-architecture-clean/backend/internal/domain/note"
)

// ---- Account ----

// MockAccountRepository is a mock of port.AccountRepository.
type MockAccountRepository struct {
	ctrl     *gomock.Controller
	recorder *MockAccountRepositoryMockRecorder
}

// MockAccountRepositoryMockRecorder records invocations.
type MockAccountRepositoryMockRecorder struct {
	mock *MockAccountRepository
}

// NewMockAccountRepository creates a new mock.
func NewMockAccountRepository(ctrl *gomock.Controller) *MockAccountRepository {
	mock := &MockAccountRepository{ctrl: ctrl}
	mock.recorder = &MockAccountRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns recorder.
func (m *MockAccountRepository) EXPECT() *MockAccountRepositoryMockRecorder {
	return m.recorder
}

func (m *MockAccountRepository) UpsertOAuthAccount(ctx context.Context, input account.OAuthAccountInput) (*account.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpsertOAuthAccount", ctx, input)
	res0, _ := ret[0].(*account.Account)
	res1, _ := ret[1].(error)
	return res0, res1
}

func (mr *MockAccountRepositoryMockRecorder) UpsertOAuthAccount(ctx, input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertOAuthAccount", reflect.TypeOf((*MockAccountRepository)(nil).UpsertOAuthAccount), ctx, input)
}

func (m *MockAccountRepository) GetByID(ctx context.Context, id string) (*account.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	res0, _ := ret[0].(*account.Account)
	res1, _ := ret[1].(error)
	return res0, res1
}

func (mr *MockAccountRepositoryMockRecorder) GetByID(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockAccountRepository)(nil).GetByID), ctx, id)
}

// MockAccountOutputPort is a mock of port.AccountOutputPort.
type MockAccountOutputPort struct {
	ctrl     *gomock.Controller
	recorder *MockAccountOutputPortMockRecorder
}

// MockAccountOutputPortMockRecorder records invocations.
type MockAccountOutputPortMockRecorder struct {
	mock *MockAccountOutputPort
}

// NewMockAccountOutputPort creates a new mock.
func NewMockAccountOutputPort(ctrl *gomock.Controller) *MockAccountOutputPort {
	mock := &MockAccountOutputPort{ctrl: ctrl}
	mock.recorder = &MockAccountOutputPortMockRecorder{mock}
	return mock
}

// EXPECT returns recorder.
func (m *MockAccountOutputPort) EXPECT() *MockAccountOutputPortMockRecorder {
	return m.recorder
}

func (m *MockAccountOutputPort) PresentAccount(ctx context.Context, acc *account.Account) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PresentAccount", ctx, acc)
	res0, _ := ret[0].(error)
	return res0
}

func (mr *MockAccountOutputPortMockRecorder) PresentAccount(ctx, acc any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PresentAccount", reflect.TypeOf((*MockAccountOutputPort)(nil).PresentAccount), ctx, acc)
}

// ---- Note ----

// MockNoteRepository is a mock of port.NoteRepository.
type MockNoteRepository struct {
	ctrl     *gomock.Controller
	recorder *MockNoteRepositoryMockRecorder
}

// MockNoteRepositoryMockRecorder records invocations.
type MockNoteRepositoryMockRecorder struct {
	mock *MockNoteRepository
}

// NewMockNoteRepository creates a new mock.
func NewMockNoteRepository(ctrl *gomock.Controller) *MockNoteRepository {
	mock := &MockNoteRepository{ctrl: ctrl}
	mock.recorder = &MockNoteRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns recorder.
func (m *MockNoteRepository) EXPECT() *MockNoteRepositoryMockRecorder {
	return m.recorder
}

func (m *MockNoteRepository) List(ctx context.Context, filters note.Filters) ([]note.WithMeta, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, filters)
	res0, _ := ret[0].([]note.WithMeta)
	res1, _ := ret[1].(error)
	return res0, res1
}

func (mr *MockNoteRepositoryMockRecorder) List(ctx, filters any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockNoteRepository)(nil).List), ctx, filters)
}

func (m *MockNoteRepository) Get(ctx context.Context, id string) (*note.WithMeta, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, id)
	res0, _ := ret[0].(*note.WithMeta)
	res1, _ := ret[1].(error)
	return res0, res1
}

func (mr *MockNoteRepositoryMockRecorder) Get(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockNoteRepository)(nil).Get), ctx, id)
}

func (m *MockNoteRepository) Create(ctx context.Context, n note.Note) (*note.Note, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, n)
	res0, _ := ret[0].(*note.Note)
	res1, _ := ret[1].(error)
	return res0, res1
}

func (mr *MockNoteRepositoryMockRecorder) Create(ctx, n any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockNoteRepository)(nil).Create), ctx, n)
}

func (m *MockNoteRepository) Update(ctx context.Context, n note.Note) (*note.Note, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, n)
	res0, _ := ret[0].(*note.Note)
	res1, _ := ret[1].(error)
	return res0, res1
}

func (mr *MockNoteRepositoryMockRecorder) Update(ctx, n any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockNoteRepository)(nil).Update), ctx, n)
}

func (m *MockNoteRepository) UpdateStatus(ctx context.Context, id string, status note.NoteStatus) (*note.Note, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateStatus", ctx, id, status)
	res0, _ := ret[0].(*note.Note)
	res1, _ := ret[1].(error)
	return res0, res1
}

func (mr *MockNoteRepositoryMockRecorder) UpdateStatus(ctx, id, status any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateStatus", reflect.TypeOf((*MockNoteRepository)(nil).UpdateStatus), ctx, id, status)
}

func (m *MockNoteRepository) Delete(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	res0, _ := ret[0].(error)
	return res0
}

func (mr *MockNoteRepositoryMockRecorder) Delete(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockNoteRepository)(nil).Delete), ctx, id)
}

func (m *MockNoteRepository) ReplaceSections(ctx context.Context, noteID string, sections []note.Section) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReplaceSections", ctx, noteID, sections)
	res0, _ := ret[0].(error)
	return res0
}

func (mr *MockNoteRepositoryMockRecorder) ReplaceSections(ctx, noteID, sections any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReplaceSections", reflect.TypeOf((*MockNoteRepository)(nil).ReplaceSections), ctx, noteID, sections)
}

// MockNoteOutputPort is a mock of port.NoteOutputPort.
type MockNoteOutputPort struct {
	ctrl     *gomock.Controller
	recorder *MockNoteOutputPortMockRecorder
}

// MockNoteOutputPortMockRecorder records invocations.
type MockNoteOutputPortMockRecorder struct {
	mock *MockNoteOutputPort
}

// NewMockNoteOutputPort creates a new mock.
func NewMockNoteOutputPort(ctrl *gomock.Controller) *MockNoteOutputPort {
	mock := &MockNoteOutputPort{ctrl: ctrl}
	mock.recorder = &MockNoteOutputPortMockRecorder{mock}
	return mock
}

// EXPECT returns recorder.
func (m *MockNoteOutputPort) EXPECT() *MockNoteOutputPortMockRecorder {
	return m.recorder
}

func (m *MockNoteOutputPort) PresentNoteList(ctx context.Context, notes []note.WithMeta) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PresentNoteList", ctx, notes)
	res0, _ := ret[0].(error)
	return res0
}

func (mr *MockNoteOutputPortMockRecorder) PresentNoteList(ctx, notes any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PresentNoteList", reflect.TypeOf((*MockNoteOutputPort)(nil).PresentNoteList), ctx, notes)
}

func (m *MockNoteOutputPort) PresentNote(ctx context.Context, n *note.WithMeta) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PresentNote", ctx, n)
	res0, _ := ret[0].(error)
	return res0
}

func (mr *MockNoteOutputPortMockRecorder) PresentNote(ctx, n any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PresentNote", reflect.TypeOf((*MockNoteOutputPort)(nil).PresentNote), ctx, n)
}

func (m *MockNoteOutputPort) PresentNoteDeleted(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PresentNoteDeleted", ctx)
	res0, _ := ret[0].(error)
	return res0
}

func (mr *MockNoteOutputPortMockRecorder) PresentNoteDeleted(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PresentNoteDeleted", reflect.TypeOf((*MockNoteOutputPort)(nil).PresentNoteDeleted), ctx)
}
