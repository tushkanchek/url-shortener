package delete_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"urlShortner/internal/http-server/handlers/url/delete"
	"urlShortner/internal/http-server/handlers/url/delete/mocks"
	"urlShortner/internal/lib/api/response"
	"urlShortner/internal/lib/logger/handlers/slogdiscard"
	"urlShortner/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteHandler(t *testing.T) {
	cases := []struct {
		name           string
		alias          string
		mockError      error
		expectedStatus int
		expectedResp   response.Response
	}{
		
		{
			name:           "url not found",
			alias:          "alias2",
			mockError:      storage.ErrURLNotFound,
			expectedStatus: http.StatusOK,
			expectedResp:   response.Error("not found"),
		},
		{
			name:           "internal error",
			alias:          "alias3",
			mockError:      errors.New("db error"),
			expectedStatus: http.StatusOK,
			expectedResp:   response.Error("internal error"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			urlDeleterMock := mocks.NewURLDeleter(t)

			// только если alias не пустой
			if tc.alias != "" {
				urlDeleterMock.
					On("DeleteURL", tc.alias).
					Return(tc.mockError).
					Once()
			}

			r := chi.NewRouter()
			r.Delete("/{alias}", delete.New(slogdiscard.NewDiscardLogger(), urlDeleterMock))

			req := httptest.NewRequest(http.MethodDelete, "/"+tc.alias, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			require.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t,
				response.ToJSON(tc.expectedResp),
				w.Body.String(),
			)
		})
	}
}
