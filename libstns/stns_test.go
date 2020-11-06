package libstns

import (
	"encoding/json"
	"errors"
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
				client: &client{
					ApiEndpoint: ts.URL,
					opt:         &Options{},
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
				client: &client{
					ApiEndpoint: ts.URL,
					opt:         &Options{},
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
				client: &client{
					ApiEndpoint: ts.URL,
					opt:         &Options{},
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
				client: &client{
					ApiEndpoint: ts.URL,
					opt:         &Options{},
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
				client: &client{
					ApiEndpoint: ts.URL,
					opt:         &Options{},
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
				client: &client{
					ApiEndpoint: ts.URL,
					opt:         &Options{},
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

func TestSTNS_Signature(t *testing.T) {
	type fields struct {
		client             *client
		PrivatekeyPath     string
		PrivatekeyPassword string
	}
	tests := []struct {
		name    string
		fields  fields
		msg     []byte
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				PrivatekeyPath:     "./testdata/id_rsa",
				PrivatekeyPassword: "test",
			},
			msg:     []byte("test"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &STNS{
				client: tt.fields.client,
				opt: &Options{
					PrivatekeyPath:     tt.fields.PrivatekeyPath,
					PrivatekeyPassword: tt.fields.PrivatekeyPassword,
				},
				makeChallengeCode: func() ([]byte, error) {
					return []byte("dummy"), nil
				},
				storeChallengeCode: func(name string, code []byte) error {
					if name == tt.name && reflect.DeepEqual(code, []byte("dummy")) {
						return nil
					}
					return errors.New("unmatch store code")
				},
			}
			_, err := c.Signature(tt.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("STNS.Signature() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestSTNS_loadPrivateKey(t *testing.T) {
	type fields struct {
		client             *client
		PrivatekeyPath     string
		PrivatekeyPassword string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				PrivatekeyPath:     "./testdata/id_rsa",
				PrivatekeyPassword: "test",
			},
			wantErr: false,
		},
		{
			name: "missmatch key",
			fields: fields{
				PrivatekeyPath:     "./testdata/id_rsa",
				PrivatekeyPassword: "error",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &STNS{
				client: tt.fields.client,
				opt: &Options{
					PrivatekeyPath:     tt.fields.PrivatekeyPath,
					PrivatekeyPassword: tt.fields.PrivatekeyPassword,
				},
			}
			_, err := c.loadPrivateKey()
			if (err != nil) != tt.wantErr {
				t.Errorf("STNS.loadPrivateKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestSTNS_Verify(t *testing.T) {
	type fields struct {
		client             *client
		PrivatekeyPath     string
		PrivatekeyPassword string
	}
	type args struct {
		msg            []byte
		publicKeyBytes []byte
		signature      []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				PrivatekeyPath:     "./testdata/id_rsa",
				PrivatekeyPassword: "test",
			},
			args: args{
				msg:            []byte("test"),
				publicKeyBytes: []byte("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCrCVhMfKrZqx/f0QCoyWwNX1VJOtSH5uY+rV1qkyt335Q4Yj5Gln/jQdRvwoS9gClXgsbXQKFg5+eSgClONE6H9iBV5QSOcHUAQTjJClxdPT9FEfq9sLSn0tlBMP0nwRaKMHpR3PPC7AB8SmPvsLoJnnE0muBSkColWOJCoTuYQKsAdD63ieqozs2LuDbaNiTGbZMyjwrSn6SjOLOAGFwpkekvlxfOOTzO11vBs+DaUnZQ1U9ZFgNAmqOp4ELhCXBC8yorXKY8T9CyG4dTkD2Zz5tMXkqo+3NpdBXqEqxmr23V3YRtZ8++QdePjknSixnpL+TmdW2K2yB+7DighuMQlkwhITj261jTQEo0AUFi6OpWAwMY0cEfMbdK0Erkw7EZg4dJTqHgp4f67wezmlv3kdPz9UruwMtVfY1uSZ4hZDkQyp4X75FKxh3dTi0K8aEt9gW5HdwOKm5MWmeJbz3r8wQpawsPotFnK8ssfHaUyqlN1Qs3UwKDMpO281R3ksc="),
				signature:      []byte(`{"Format":"ssh-rsa","Blob":"YL0Elcvpcs2RfBMesMPmDiBI0ppwiPpdmwWDnOAEAzKvshWoBTPWhy7qE/VDwnwcTvbk9SyJovrAWPdOmPcnXuKtgQam4NorFWQWMRFI6/tL1C3JWjY/uuALD+bH0WUbmCvUCnCQ3s7tG6UNx0JDyP//bh2IV/B1tds24c2hd36MRpufknUbsD303welyxVcdGRKm0bi3hp+X/NFWHFo5TWe1qw+5mJpqR32flkflGXcHZkzRk+5hs3YbrM/Je5hEur55lCVz2pLv/zzF72nAyMxM3rpHJi8aNe5uI1ZuJW+GsK5SW3V7Wn1uJgfI5et7P7H2/i03EHgV8RGUR5bIMF+PC8dG4T9v5dJr/rZy8QfcQrVTcMioDcjB0BW/wE/JmAy3NvmOhkJb68Ec5FF7nisbpdGv+5UcH+3stSZR3g6mp1cYxNtTx+01AV1+ZIiC6dhh0lz+O81NrY+A98ldBdnXi7pDD/0B+UldDzkPZ/gqidSQLHKEFqJTEnPFud1","Rest":null}`),
			},

			wantErr: false,
		},
		{
			name: "multikey",
			fields: fields{
				PrivatekeyPath:     "./testdata/id_rsa",
				PrivatekeyPassword: "test",
			},
			args: args{
				msg:            []byte("test"),
				publicKeyBytes: []byte("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCrCVhMfKrZqx/f0QCoyWwNX1VJOtSH5uY+rV1qkyt335Q4Yj5Gln/jQdRvwoS9gClXgsbXQKFg5+eSgClONE6H9iBV5QSOcHUAQTjJClxdPT9FEfq9sLSn0tlBMP0nwRaKMHpR3PPC7AB8SmPvsLoJnnE0muBSkColWOJCoTuYQKsAdD63ieqozs2LuDbaNiTGbZMyjwrSn6SjOLOAGFwpkekvlxfOOTzO11vBs+DaUnZQ1U9ZFgNAmqOp4ELhCXBC8yorXKY8T9CyG4dTkD2Zz5tMXkqo+3NpdBXqEqxmr23V3YRtZ8++QdePjknSixnpL+TmdW2K2yB+7DighuMQlkwhITj261jTQEo0AUFi6OpWAwMY0cEfMbdK0Erkw7EZg4dJTqHgp4f67wezmlv3kdPz9UruwMtVfY1uSZ4hZDkQyp4X75FKxh3dTi0K8aEt9gW5HdwOKm5MWmeJbz3r8wQpawsPotFnK8ssfHaUyqlN1Qs3UwKDMpO281R3ksc=\nssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCwayeIUERUKcKX4PYJpBRN6Sd7QpecD026HFJiiOi6UlEEEKcgPykB5+UVYXaU+jCJK/b5+pPqWXm848furoVL0qMxR/k+tBH9jMZgkeHumoM6YOQYOi6SvxC7Bqo4846DD63aHvaDLwixVGtJYRQBXlWD2AGJDSZVxeiJ8b72LnUdMhEhHs+GHAcXumxxlEl1XPBkVE8ncB10utcAxiQC9+DRKwrtwGwHnBQ2Zu6Ms9s2BkI6RxEDqqjGq2sqMiulvG68hLAhHPSwBnyBPfzQJCnP+xPqw1j+2Pl4hdseW4Lf0Kdet2tkf6fz93XAfdkr3nAUNOY8fJ3GZQ+xvV/Y2DkEPAocKKi4A3w0MonMLSO/aowArrJWNOCyaUOgpgvcb4d4rRWKF/fHq0SYkVGg7eKnTBPByqiB6KZfLSrE9flptzAfY5hokLx2tIEV5jsG0arzTks5j8uS+U/Om9UiFrymZNALoapiKH+SwbqQfi87oInHMSVsLxBtFyamhD0="),
				signature:      []byte(`{"Format":"ssh-rsa","Blob":"YL0Elcvpcs2RfBMesMPmDiBI0ppwiPpdmwWDnOAEAzKvshWoBTPWhy7qE/VDwnwcTvbk9SyJovrAWPdOmPcnXuKtgQam4NorFWQWMRFI6/tL1C3JWjY/uuALD+bH0WUbmCvUCnCQ3s7tG6UNx0JDyP//bh2IV/B1tds24c2hd36MRpufknUbsD303welyxVcdGRKm0bi3hp+X/NFWHFo5TWe1qw+5mJpqR32flkflGXcHZkzRk+5hs3YbrM/Je5hEur55lCVz2pLv/zzF72nAyMxM3rpHJi8aNe5uI1ZuJW+GsK5SW3V7Wn1uJgfI5et7P7H2/i03EHgV8RGUR5bIMF+PC8dG4T9v5dJr/rZy8QfcQrVTcMioDcjB0BW/wE/JmAy3NvmOhkJb68Ec5FF7nisbpdGv+5UcH+3stSZR3g6mp1cYxNtTx+01AV1+ZIiC6dhh0lz+O81NrY+A98ldBdnXi7pDD/0B+UldDzkPZ/gqidSQLHKEFqJTEnPFud1","Rest":null}`),
			},

			wantErr: false,
		},
		{
			name: "msg unmatch",
			fields: fields{
				PrivatekeyPath:     "./testdata/id_rsa",
				PrivatekeyPassword: "test",
			},
			args: args{
				msg:            []byte("unmatch"),
				publicKeyBytes: []byte("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCrCVhMfKrZqx/f0QCoyWwNX1VJOtSH5uY+rV1qkyt335Q4Yj5Gln/jQdRvwoS9gClXgsbXQKFg5+eSgClONE6H9iBV5QSOcHUAQTjJClxdPT9FEfq9sLSn0tlBMP0nwRaKMHpR3PPC7AB8SmPvsLoJnnE0muBSkColWOJCoTuYQKsAdD63ieqozs2LuDbaNiTGbZMyjwrSn6SjOLOAGFwpkekvlxfOOTzO11vBs+DaUnZQ1U9ZFgNAmqOp4ELhCXBC8yorXKY8T9CyG4dTkD2Zz5tMXkqo+3NpdBXqEqxmr23V3YRtZ8++QdePjknSixnpL+TmdW2K2yB+7DighuMQlkwhITj261jTQEo0AUFi6OpWAwMY0cEfMbdK0Erkw7EZg4dJTqHgp4f67wezmlv3kdPz9UruwMtVfY1uSZ4hZDkQyp4X75FKxh3dTi0K8aEt9gW5HdwOKm5MWmeJbz3r8wQpawsPotFnK8ssfHaUyqlN1Qs3UwKDMpO281R3ksc="),
				signature:      []byte(`{"Format":"ssh-rsa","Blob":"YL0Elcvpcs2RfBMesMPmDiBI0ppwiPpdmwWDnOAEAzKvshWoBTPWhy7qE/VDwnwcTvbk9SyJovrAWPdOmPcnXuKtgQam4NorFWQWMRFI6/tL1C3JWjY/uuALD+bH0WUbmCvUCnCQ3s7tG6UNx0JDyP//bh2IV/B1tds24c2hd36MRpufknUbsD303welyxVcdGRKm0bi3hp+X/NFWHFo5TWe1qw+5mJpqR32flkflGXcHZkzRk+5hs3YbrM/Je5hEur55lCVz2pLv/zzF72nAyMxM3rpHJi8aNe5uI1ZuJW+GsK5SW3V7Wn1uJgfI5et7P7H2/i03EHgV8RGUR5bIMF+PC8dG4T9v5dJr/rZy8QfcQrVTcMioDcjB0BW/wE/JmAy3NvmOhkJb68Ec5FF7nisbpdGv+5UcH+3stSZR3g6mp1cYxNtTx+01AV1+ZIiC6dhh0lz+O81NrY+A98ldBdnXi7pDD/0B+UldDzkPZ/gqidSQLHKEFqJTEnPFud1","Rest":null}`),
			},

			wantErr: true,
		},
		{
			name: "pubkey unmatch",
			fields: fields{
				PrivatekeyPath:     "./testdata/id_rsa",
				PrivatekeyPassword: "test",
			},
			args: args{
				msg:            []byte("test"),
				publicKeyBytes: []byte("ssh-rsa BBBB3NzaC1yc2EAAAADAQABAAABgQCrCVhMfKrZqx/f0QCoyWwNX1VJOtSH5uY+rV1qkyt335Q4Yj5Gln/jQdRvwoS9gClXgsbXQKFg5+eSgClONE6H9iBV5QSOcHUAQTjJClxdPT9FEfq9sLSn0tlBMP0nwRaKMHpR3PPC7AB8SmPvsLoJnnE0muBSkColWOJCoTuYQKsAdD63ieqozs2LuDbaNiTGbZMyjwrSn6SjOLOAGFwpkekvlxfOOTzO11vBs+DaUnZQ1U9ZFgNAmqOp4ELhCXBC8yorXKY8T9CyG4dTkD2Zz5tMXkqo+3NpdBXqEqxmr23V3YRtZ8++QdePjknSixnpL+TmdW2K2yB+7DighuMQlkwhITj261jTQEo0AUFi6OpWAwMY0cEfMbdK0Erkw7EZg4dJTqHgp4f67wezmlv3kdPz9UruwMtVfY1uSZ4hZDkQyp4X75FKxh3dTi0K8aEt9gW5HdwOKm5MWmeJbz3r8wQpawsPotFnK8ssfHaUyqlN1Qs3UwKDMpO281R3ksc="),
				signature:      []byte(`{"Format":"ssh-rsa","Blob":"YL0Elcvpcs2RfBMesMPmDiBI0ppwiPpdmwWDnOAEAzKvshWoBTPWhy7qE/VDwnwcTvbk9SyJovrAWPdOmPcnXuKtgQam4NorFWQWMRFI6/tL1C3JWjY/uuALD+bH0WUbmCvUCnCQ3s7tG6UNx0JDyP//bh2IV/B1tds24c2hd36MRpufknUbsD303welyxVcdGRKm0bi3hp+X/NFWHFo5TWe1qw+5mJpqR32flkflGXcHZkzRk+5hs3YbrM/Je5hEur55lCVz2pLv/zzF72nAyMxM3rpHJi8aNe5uI1ZuJW+GsK5SW3V7Wn1uJgfI5et7P7H2/i03EHgV8RGUR5bIMF+PC8dG4T9v5dJr/rZy8QfcQrVTcMioDcjB0BW/wE/JmAy3NvmOhkJb68Ec5FF7nisbpdGv+5UcH+3stSZR3g6mp1cYxNtTx+01AV1+ZIiC6dhh0lz+O81NrY+A98ldBdnXi7pDD/0B+UldDzkPZ/gqidSQLHKEFqJTEnPFud1","Rest":null}`),
			},

			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &STNS{
				client: tt.fields.client,
				opt: &Options{
					PrivatekeyPath:     tt.fields.PrivatekeyPath,
					PrivatekeyPassword: tt.fields.PrivatekeyPassword,
				},
			}
			if err := c.verify(tt.args.msg, tt.args.publicKeyBytes, tt.args.signature); (err != nil) != tt.wantErr {
				t.Errorf("STNS.Verify() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
