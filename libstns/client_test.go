package libstns

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestClient_Request(t *testing.T) {
	tests := []struct {
		name         string
		opt          *Options
		path         string
		query        string
		responseCode int
		responseBody string
		want         *Response
		wantErr      bool
	}{
		{
			name:         "ok",
			opt:          &Options{},
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
			opt:          &Options{},
			path:         "test",
			responseCode: http.StatusNotFound,
			responseBody: "notfound",
			want: &Response{
				StatusCode: http.StatusNotFound,
				Body:       []byte("notfound"),
				Headers:    map[string]string{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.responseCode)
				fmt.Fprint(w, tt.responseBody)
			}))
			defer ts.Close()
			h, err := newClient(
				ts.URL,
				tt.opt,
			)
			if err != nil {
				t.Error(err)
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

func Test_newClient(t *testing.T) {
	type args struct {
		endpoint string
		opt      *Options
	}
	tests := []struct {
		name string
		args args
		envs map[string]string
		want *client
	}{
		{
			name: "default value ok",
			args: args{
				endpoint: "http://localhost",
			},
			want: &client{
				ApiEndpoint: "http://localhost",
				opt: &Options{
					PrivatekeyPath: "~/.ssh/id_rsa",
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
				opt: &Options{
					UserAgent:      "libstns-go/update",
					RequestTimeout: 30,
					RequestRetry:   6,
				},
			},
			want: &client{
				ApiEndpoint: "http://localhost",
				opt: &Options{
					PrivatekeyPath: "~/.ssh/id_rsa",
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
			want: &client{
				ApiEndpoint: "http://localhost",
				opt: &Options{
					PrivatekeyPath: "~/.ssh/id_rsa",
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
			for k, v := range tt.envs {
				t.Setenv(k, v)
			}
			got, err := newClient(tt.args.endpoint, tt.args.opt)
			if err != nil {
				t.Errorf("newClient() error = %v", err)
				return
			}
			if got.ApiEndpoint != tt.want.ApiEndpoint {
				t.Errorf("newClient() ApiEndpoint = %v, want %v", got.ApiEndpoint, tt.want.ApiEndpoint)
			}
			if !reflect.DeepEqual(got.opt, tt.want.opt) {
				t.Errorf("newClient() opt = %v, want %v", got.opt, tt.want.opt)
			}
		})
	}
}
