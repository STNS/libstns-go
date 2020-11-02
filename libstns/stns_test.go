package libstns

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/STNS/STNS/v2/model"
)

func TestSTNS_ListUser(t *testing.T) {
	tests := []struct {
		name         string
		want         []*model.User
		wantErr      bool
		responseCode int
		responseUser []model.User
	}{
		{
			name:         "ok",
			responseCode: http.StatusOK,
			responseUser: []model.User{
				model.User{
					Base: model.Base{
						ID:   1,
						Name: "example1",
					},
				},
				model.User{
					Base: model.Base{
						ID:   2,
						Name: "example2",
					},
				},
			},
			want: []*model.User{
				&model.User{
					Base: model.Base{
						ID:   1,
						Name: "example1",
					},
				},
				&model.User{
					Base: model.Base{
						ID:   2,
						Name: "example2",
					},
				},
			},
			wantErr: false,
		},
		{
			name:         "notfound",
			responseCode: http.StatusNotFound,
			responseUser: []model.User{},
			want:         nil,
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.String() == "/users" {
					w.WriteHeader(tt.responseCode)
					if tt.responseUser != nil {
						rp, err := json.Marshal(tt.responseUser)
						if err != nil {
							t.Error(err)
						}
						fmt.Fprintf(w, string(rp))
					}
				}
			}))
			defer ts.Close()

			s := &STNS{
				client: &Client{
					ApiEndpoint: ts.URL,
					opt:         &ClientOptions{},
				},
			}
			got, err := s.ListUser()
			if (err != nil) != tt.wantErr {
				t.Errorf("STNS.ListUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("STNS.ListUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSTNS_GetUserByName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    *model.User
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				name: "example1",
			},
			want: &model.User{
				Base: model.Base{
					ID:   1,
					Name: "example1",
				},
			},
			wantErr: false,
		},
		{
			name: "notfound",
			args: args{
				name: "example2",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.FormValue("name") == "example1" && r.URL.String() == "/users?name=example1" {
					rp, err := json.Marshal([]*model.User{
						&model.User{
							Base: model.Base{
								ID:   1,
								Name: "example1",
							},
						}})

					if err != nil {
						t.Error(err)
					}
					fmt.Fprintf(w, string(rp))

					w.WriteHeader(http.StatusOK)
				} else {
					w.WriteHeader(http.StatusNotFound)
				}
			}))
			defer ts.Close()

			s := &STNS{
				client: &Client{
					ApiEndpoint: ts.URL,
					opt:         &ClientOptions{},
				},
			}

			got, err := s.GetUserByName(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("STNS.GetUserByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("STNS.GetUserByName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSTNS_GetUserByID(t *testing.T) {
	type args struct {
		id int
	}
	tests := []struct {
		name    string
		args    args
		want    *model.User
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				id: 1,
			},
			want: &model.User{
				Base: model.Base{
					ID:   1,
					Name: "example1",
				},
			},
			wantErr: false,
		},
		{
			name: "notfound",
			args: args{
				id: 2,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.FormValue("id") == "1" && r.URL.String() == "/users?id=1" {
					rp, err := json.Marshal([]*model.User{
						&model.User{
							Base: model.Base{
								ID:   1,
								Name: "example1",
							},
						}})

					if err != nil {
						t.Error(err)
					}
					fmt.Fprintf(w, string(rp))

					w.WriteHeader(http.StatusOK)
				} else {
					w.WriteHeader(http.StatusNotFound)
				}
			}))
			defer ts.Close()

			s := &STNS{
				client: &Client{
					ApiEndpoint: ts.URL,
					opt:         &ClientOptions{},
				},
			}

			got, err := s.GetUserByID(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("STNS.GetUserByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("STNS.GetUserByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSTNS_ListGroup(t *testing.T) {
	tests := []struct {
		name          string
		want          []*model.Group
		wantErr       bool
		responseCode  int
		responseGroup []model.Group
	}{
		{
			name:         "ok",
			responseCode: http.StatusOK,
			responseGroup: []model.Group{
				model.Group{
					Base: model.Base{
						ID:   1,
						Name: "example1",
					},
				},
				model.Group{
					Base: model.Base{
						ID:   2,
						Name: "example2",
					},
				},
			},
			want: []*model.Group{
				&model.Group{
					Base: model.Base{
						ID:   1,
						Name: "example1",
					},
				},
				&model.Group{
					Base: model.Base{
						ID:   2,
						Name: "example2",
					},
				},
			},
			wantErr: false,
		},
		{
			name:          "notfound",
			responseCode:  http.StatusNotFound,
			responseGroup: []model.Group{},
			want:          nil,
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.String() == "/groups" {
					w.WriteHeader(tt.responseCode)
					if tt.responseGroup != nil {
						rp, err := json.Marshal(tt.responseGroup)
						if err != nil {
							t.Error(err)
						}
						fmt.Fprintf(w, string(rp))
					}
				}
			}))
			defer ts.Close()

			s := &STNS{
				client: &Client{
					ApiEndpoint: ts.URL,
					opt:         &ClientOptions{},
				},
			}
			got, err := s.ListGroup()
			if (err != nil) != tt.wantErr {
				t.Errorf("STNS.ListGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("STNS.ListGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSTNS_GetGroupByName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    *model.Group
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				name: "example1",
			},
			want: &model.Group{
				Base: model.Base{
					ID:   1,
					Name: "example1",
				},
			},
			wantErr: false,
		},
		{
			name: "notfound",
			args: args{
				name: "example2",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.FormValue("name") == "example1" && r.URL.String() == "/groups?name=example1" {
					rp, err := json.Marshal([]*model.Group{
						&model.Group{
							Base: model.Base{
								ID:   1,
								Name: "example1",
							},
						}})

					if err != nil {
						t.Error(err)
					}
					fmt.Fprintf(w, string(rp))

					w.WriteHeader(http.StatusOK)
				} else {
					w.WriteHeader(http.StatusNotFound)
				}
			}))
			defer ts.Close()

			s := &STNS{
				client: &Client{
					ApiEndpoint: ts.URL,
					opt:         &ClientOptions{},
				},
			}

			got, err := s.GetGroupByName(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("STNS.GetGroupByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("STNS.GetGroupByName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSTNS_GetGroupByID(t *testing.T) {
	type args struct {
		id int
	}
	tests := []struct {
		name    string
		args    args
		want    *model.Group
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				id: 1,
			},
			want: &model.Group{
				Base: model.Base{
					ID:   1,
					Name: "example1",
				},
			},
			wantErr: false,
		},
		{
			name: "notfound",
			args: args{
				id: 2,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.FormValue("id") == "1" && r.URL.String() == "/groups?id=1" {
					rp, err := json.Marshal([]*model.Group{
						&model.Group{
							Base: model.Base{
								ID:   1,
								Name: "example1",
							},
						}})

					if err != nil {
						t.Error(err)
					}
					fmt.Fprintf(w, string(rp))

					w.WriteHeader(http.StatusOK)
				} else {
					w.WriteHeader(http.StatusNotFound)
				}
			}))
			defer ts.Close()

			s := &STNS{
				client: &Client{
					ApiEndpoint: ts.URL,
					opt:         &ClientOptions{},
				},
			}

			got, err := s.GetGroupByID(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("STNS.GetGroupByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("STNS.GetGroupByID() = %v, want %v", got, tt.want)
			}
		})
	}
}
