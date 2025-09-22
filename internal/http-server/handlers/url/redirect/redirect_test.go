package redirect_test

import (
	"urlShortner/internal/http-server/handlers/url/redirect"
	"urlShortner/internal/http-server/handlers/url/redirect/mocks"
	"urlShortner/internal/lib/api"
	"urlShortner/internal/lib/logger/handlers/slogdiscard"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)



func TestRedirectHandler(t *testing.T){
	cases := []struct{
		name      string
		alias     string
		url       string
		respError string
		mockError error
	}{
		{
			name:	"Success",
			alias:	"test_alias",
			url:	"https://google.com",
		},
		{
			name:	"Success",
			alias:	"test_alias2",
			url:	"https://lol.com",
		},
		
	}
	for _, tc := range cases{
		t.Run(tc.name, func(t *testing.T){
			urlGetterMock := mocks.NewURLGetter(t)

			if tc.respError == "" || tc.mockError != nil{
				urlGetterMock.On("GetURL", tc.alias).
					Return(tc.url, tc.mockError).Once()
			}
		

			r := chi.NewRouter()
			r.Get("/{alias}", redirect.New(slogdiscard.NewDiscardLogger(), urlGetterMock))

			ts:= httptest.NewServer(r)
			defer ts.Close()

			redirectedToURL, err := api.GetRedirect(ts.URL + "/" + tc.alias)
			require.NoError(t, err)

			//Check the final url redirection
			assert.Equal(t, tc.url, redirectedToURL)

			

		})	
	}
}