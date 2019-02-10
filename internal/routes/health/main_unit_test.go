package main_test

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	"github.com/aws/aws-lambda-go/lambdacontext"
	api "github.com/jayco/serverless-golang-demo/internal/routes/health"
)

func generateLambdaCtx() context.Context {
	ctx := context.Background()
	lc := new(lambdacontext.LambdaContext)
	ctx = lambdacontext.NewContext(ctx, lc)
	return ctx
}

var expectedHeaders = map[string]string{"Content-Type": "application/json", "X-BaseAPI-V1-Reply": "health-handler"}

// Unit Tests
func TestHandlerUnit(t *testing.T) {
	type args struct {
		ctx     context.Context
		request api.Request
	}
	tests := []struct {
		name    string
		args    args
		want    api.Response
		wantErr bool
	}{
		{
			"the handler responds with the correct 204 and custom headers",
			args{generateLambdaCtx(), api.Request{}},
			api.Response{StatusCode: http.StatusNoContent, Headers: expectedHeaders},
			false,
		},
		{
			"The handler responds with the correct 404 and custom headers when no lambda context",
			args{context.Background(), api.Request{}},
			api.Response{StatusCode: http.StatusNotFound, Headers: expectedHeaders},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := api.Handler(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Handler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.StatusCode, tt.want.StatusCode) {
				t.Errorf("Handler() = StatusCode %v, want %v", got.StatusCode, tt.want.StatusCode)
			}
			if !reflect.DeepEqual(got.Headers, tt.want.Headers) {
				t.Errorf("Handler() = Headers %v, want %v", got.Headers, tt.want.Headers)
			}
		})
	}
}
