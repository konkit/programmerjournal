package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"

	"github.com/go-chi/chi/v5"
)

func TestStaticFilesHandler(t *testing.T) {
	router := chi.NewMux()
	staticFilesHandler(router, fstest.MapFS{
		"index.html":     {Data: []byte("<app-root>")},
		"main-ABC123.js": {Data: []byte("console.log('app')")},
	})

	testCases := []struct {
		name     string
		path     string
		wantCode int
		wantBody string
	}{
		{name: "root", path: "/", wantCode: http.StatusOK, wantBody: "<app-root>"},
		{name: "existing_asset", path: "/main-ABC123.js", wantCode: http.StatusOK, wantBody: "console.log('app')"},
		{name: "app_route", path: "/notes", wantCode: http.StatusOK, wantBody: "<app-root>"},
		{name: "nested_app_route", path: "/day/2024-05-01", wantCode: http.StatusOK, wantBody: "<app-root>"},
		{name: "app_route_with_param", path: "/tags/team", wantCode: http.StatusOK, wantBody: "<app-root>"},
		{name: "missing_asset", path: "/main-OLD.js", wantCode: http.StatusNotFound},
		{name: "unknown_api_route", path: "/api/nope", wantCode: http.StatusNotFound},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := httptest.NewRecorder()
			router.ServeHTTP(res, httptest.NewRequest(http.MethodGet, tc.path, nil))

			if res.Code != tc.wantCode {
				t.Fatalf("GET %s: status %d, want %d", tc.path, res.Code, tc.wantCode)
			}
			if tc.wantBody != "" && res.Body.String() != tc.wantBody {
				t.Errorf("GET %s: body %q, want %q", tc.path, res.Body.String(), tc.wantBody)
			}
		})
	}
}
