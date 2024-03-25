package handlers

import (
	"database/sql"
	"github.com/Lesnoi3283/url_shortener/internal/app/handlers/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_pingDBHandler_ServeHTTP(t *testing.T) {

	tests := []struct {
		name          string
		mockSetupFunc func(dbInterface *mocks.MockDBInterface)
		statusWant    int
	}{
		{
			name: "db works",
			mockSetupFunc: func(dbInterface *mocks.MockDBInterface) {
				dbInterface.EXPECT().Ping().Return(nil)
			},
			statusWant: http.StatusOK,
		},
		{
			name: "db doesnt work",
			mockSetupFunc: func(dbInterface *mocks.MockDBInterface) {
				dbInterface.EXPECT().Ping().Return(sql.ErrConnDone)
			},
			statusWant: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockController := gomock.NewController(t)
			db := mocks.NewMockDBInterface(mockController)

			tt.mockSetupFunc(db)

			p := &pingDBHandler{
				db: db,
			}

			req, err := http.NewRequest("GET", "/ping", nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()

			p.ServeHTTP(rr, req)
			assert.Equal(t, tt.statusWant, rr.Code, "Status code is not equal.")
		})
	}
}
