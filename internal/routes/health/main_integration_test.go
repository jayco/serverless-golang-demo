package main_test

import (
	"net/http"
	"reflect"
	"testing"

	api "github.com/jayco/serverless-golang-demo/internal/routes/health"
	integration "github.com/jayco/serverless-golang-demo/pkg"
)

// Integration Tests
func TestHandlerIntegration(t *testing.T) {
	template := "../../../sam-local-template.yml"
	handler := "Health"
	tests := []struct {
		name       string
		args       string
		want       api.Response
		errStrRegx string
	}{
		{
			"lambda responds with the correct 204 and custom headers",
			"",
			api.Response{StatusCode: http.StatusNoContent, Headers: expectedHeaders},
			"errorMessage:",
		},
	}
	for _, tt := range tests {
		got, err := integration.Run(template, handler, tt.args, tt.errStrRegx)
		if err != nil {
			t.Error("unexpected error:", err)
		}

		if !reflect.DeepEqual(got.StatusCode, tt.want.StatusCode) {
			t.Errorf("Handler() = StatusCode %v, want %v", got.StatusCode, tt.want.StatusCode)
		}
		if !reflect.DeepEqual(got.Headers, tt.want.Headers) {
			t.Errorf("Handler() = Headers %v, want %v", got.Headers, tt.want.Headers)
		}
	}
}
