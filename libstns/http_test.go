package libstns

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/k0kubun/pp"
)

func TestHttp_Request(t *testing.T) {
	tests := []struct {
		name         string
		opt          *HttpOptions
		path         string
		query        string
		responseCode int
		responseBody string
		want         *Response
		wantErr      bool
	}{
		{
			name:         "ok",
			opt:          &HttpOptions{},
			path:         "test",
			responseCode: http.StatusOK,
			responseBody: "it is ok",
			want: &Response{
				StatusCode: http.StatusOK,
				Body:       []byte("it is ok"),
				Headers:    map[string]string{},
			},
			wantErr: false,
		},
		{
			name:         "notfound",
			opt:          &HttpOptions{},
			path:         "test",
			responseCode: http.StatusNotFound,
			responseBody: "",
			want: &Response{
				StatusCode: http.StatusNotFound,
				Body:       nil,
				Headers:    map[string]string{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.responseCode)
				fmt.Fprintf(w, tt.responseBody)
			}))
			defer ts.Close()
			h := &Http{
				ApiEndpoint: ts.URL,
				opt:         tt.opt,
			}
			got, err := h.Request(tt.path, tt.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("Http.Request() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Http.Request() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewHttp(t *testing.T) {
	type args struct {
		endpoint string
		opt      *HttpOptions
	}
	tests := []struct {
		name string
		args args
		want *Http
	}{
		{
			name: "default value ok",
			args: args{
				endpoint: "http://localhost",
			},
			want: &Http{
				ApiEndpoint: "http://localhost",
				opt: &HttpOptions{
					UserAgent:      "libstns-go/0.0.1",
					RequestTimeout: 15,
					RequestRetry:   3,
				},
			},
		},
		{
			name: "set value ok",
			args: args{
				endpoint: "http://localhost",
				opt: &HttpOptions{
					UserAgent:      "libstns-go/update",
					RequestTimeout: 30,
					RequestRetry:   6,
				},
			},
			want: &Http{
				ApiEndpoint: "http://localhost",
				opt: &HttpOptions{
					UserAgent:      "libstns-go/update",
					RequestTimeout: 30,
					RequestRetry:   6,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHttp(tt.args.endpoint, tt.args.opt); !reflect.DeepEqual(got, tt.want) {
				pp.Println(got)
				t.Errorf("NewHttp() = %v, want %v", got, tt.want)
			}
		})
	}
}
