package libstns

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/k0kubun/pp"
)

func TestClient_Request(t *testing.T) {
	tests := []struct {
		name         string
		opt          *ClientOptions
		path         string
		query        string
		responseCode int
		responseBody string
		want         *Response
		wantErr      bool
	}{
		{
			name:         "ok",
			opt:          &ClientOptions{},
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
			opt:          &ClientOptions{},
			path:         "test",
			responseCode: http.StatusNotFound,
			responseBody: "notfound",
			want:         nil,
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.responseCode)
				fmt.Fprintf(w, tt.responseBody)
			}))
			defer ts.Close()
			h := &Client{
				ApiEndpoint: ts.URL,
				opt:         tt.opt,
			}
			got, err := h.Request(tt.path, tt.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Request() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.Request() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewClient(t *testing.T) {
	type args struct {
		endpoint string
		opt      *ClientOptions
	}
	tests := []struct {
		name string
		args args
		envs map[string]string
		want *Client
	}{
		{
			name: "default value ok",
			args: args{
				endpoint: "http://localhost",
			},
			want: &Client{
				ApiEndpoint: "http://localhost",
				opt: &ClientOptions{
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
				opt: &ClientOptions{
					UserAgent:      "libstns-go/update",
					RequestTimeout: 30,
					RequestRetry:   6,
				},
			},
			want: &Client{
				ApiEndpoint: "http://localhost",
				opt: &ClientOptions{
					UserAgent:      "libstns-go/update",
					RequestTimeout: 30,
					RequestRetry:   6,
				},
			},
		},
		{
			name: "set envs ok",
			args: args{
				endpoint: "http://localhost",
			},
			want: &Client{
				ApiEndpoint: "http://localhost",
				opt: &ClientOptions{
					UserAgent:      "libstns-go/0.0.1",
					RequestTimeout: 15,
					RequestRetry:   3,
					User:           "example user",
					Password:       "example password",
				},
			},
			envs: map[string]string{
				"STNS_USER":     "example user",
				"STNS_PASSWORD": "example password",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.envs) > 0 {
				for k, v := range tt.envs {
					os.Setenv(k, v)
				}
			}
			if got, _ := NewClient(tt.args.endpoint, tt.args.opt); !reflect.DeepEqual(got, tt.want) {
				pp.Println(got)
				pp.Println(tt.want)
				t.Errorf("NewClient() = %v, want %v", got, tt.want)
			}
		})
	}
}
