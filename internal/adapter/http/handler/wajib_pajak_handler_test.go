package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bpka/tps-pkb/internal/adapter/http/handler"
	mw "github.com/bpka/tps-pkb/internal/adapter/http/middleware"
	"github.com/bpka/tps-pkb/internal/domain/dto"
	"github.com/bpka/tps-pkb/internal/domain/entity"
	"github.com/bpka/tps-pkb/internal/domain/port"
	"github.com/bpka/tps-pkb/internal/mock"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func withAdminClaims(ctx context.Context) context.Context {
	return context.WithValue(ctx, mw.CtxKeyClaims, &port.Claims{
		UserID:   uuid.New(),
		Username: "admin",
		Role:     entity.RoleAdmin,
	})
}

func newValidate() *validator.Validate {
	v := validator.New()
	_ = v.RegisterValidation("nik", func(fl validator.FieldLevel) bool {
		s := fl.Field().String()
		if len(s) != 16 {
			return false
		}
		for _, c := range s {
			if c < '0' || c > '9' {
				return false
			}
		}
		return true
	})
	_ = v.RegisterValidation("npwp", func(fl validator.FieldLevel) bool {
		s := fl.Field().String()
		return len(s) == 15
	})
	return v
}

func TestWajibPajakHandler_Create_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock.NewMockWajibPajakService(ctrl)
	expectedWP := &entity.WajibPajak{
		ID:    uuid.New(),
		Nama:  "Budi Santoso",
		Jenis: entity.JenisIndividu,
	}
	mockSvc.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(expectedWP, nil)

	h := handler.NewWajibPajakHandler(mockSvc, newValidate())

	body := `{
		"nama":"Budi Santoso",
		"jenis":"INDIVIDU",
		"nik":"1234567890123456",
		"alamat":"Jl. Merdeka No. 1"
	}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/wajib-pajak", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(withAdminClaims(req.Context()))

	rr := httptest.NewRecorder()
	h.Create(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var resp map[string]any
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
	assert.Equal(t, true, resp["success"])
}

func TestWajibPajakHandler_Create_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock.NewMockWajibPajakService(ctrl)
	// Create must NOT be called when validation fails
	mockSvc.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0)

	h := handler.NewWajibPajakHandler(mockSvc, newValidate())

	body := `{"nama":"","jenis":"INDIVIDU","alamat":"Jl. A"}` // nama too short
	req := httptest.NewRequest(http.MethodPost, "/api/v1/wajib-pajak", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(withAdminClaims(req.Context()))

	rr := httptest.NewRecorder()
	h.Create(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var resp map[string]any
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
	assert.Equal(t, false, resp["success"])
}

func TestWajibPajakHandler_Create_Conflict(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock.NewMockWajibPajakService(ctrl)
	mockSvc.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, dto.ErrConflict)

	h := handler.NewWajibPajakHandler(mockSvc, newValidate())

	body := `{"nama":"Budi","jenis":"INDIVIDU","nik":"1234567890123456","alamat":"Jl. A"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/wajib-pajak", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(withAdminClaims(req.Context()))

	rr := httptest.NewRecorder()
	h.Create(rr, req)

	assert.Equal(t, http.StatusConflict, rr.Code)
}

func TestWajibPajakHandler_List_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock.NewMockWajibPajakService(ctrl)
	mockSvc.EXPECT().
		List(gomock.Any(), gomock.Any()).
		Return([]entity.WajibPajak{{ID: uuid.New(), Nama: "Test"}}, 1, nil)

	h := handler.NewWajibPajakHandler(mockSvc, newValidate())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/wajib-pajak?page=1&limit=10", nil)
	req = req.WithContext(withAdminClaims(req.Context()))

	rr := httptest.NewRecorder()
	h.List(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp map[string]any
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
	assert.Equal(t, true, resp["success"])
	assert.NotNil(t, resp["meta"])
}
