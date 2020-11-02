package libstns

import (
	"encoding/json"
	"fmt"

	"github.com/STNS/STNS/v2/model"
)

type STNS struct {
	client *Client
}

func NewSTNS(client *Client) *STNS {
	return &STNS{
		client: client,
	}
}

const usersEndpoint = "/users"
const groupsEndpoint = "/groups"

func (s *STNS) ListUser() ([]*model.User, error) {
	r, err := s.client.Request(usersEndpoint, "")
	if err != nil {
		return nil, err
	}
	v := []*model.User{}
	if err := json.Unmarshal(r.Body, &v); err != nil {
		return nil, err
	}

	return v, nil
}

func (s *STNS) GetUserByName(name string) (*model.User, error) {
	r, err := s.client.Request(usersEndpoint, fmt.Sprintf("name=%s", name))
	if err != nil {
		return nil, err
	}
	v := []*model.User{}
	if err := json.Unmarshal(r.Body, &v); err != nil {
		return nil, err
	}

	return v[0], nil
}

func (s *STNS) GetUserByID(id int) (*model.User, error) {
	r, err := s.client.Request(usersEndpoint, fmt.Sprintf("id=%d", id))
	if err != nil {
		return nil, err
	}
	v := []*model.User{}
	if err := json.Unmarshal(r.Body, &v); err != nil {
		return nil, err
	}

	return v[0], nil
}

func (s *STNS) ListGroup() ([]*model.Group, error) {
	r, err := s.client.Request(groupsEndpoint, "")
	if err != nil {
		return nil, err
	}
	v := []*model.Group{}
	if err := json.Unmarshal(r.Body, &v); err != nil {
		return nil, err
	}

	return v, nil
}

func (s *STNS) GetGroupByName(name string) (*model.Group, error) {
	r, err := s.client.Request(groupsEndpoint, fmt.Sprintf("name=%s", name))
	if err != nil {
		return nil, err
	}
	v := []*model.Group{}
	if err := json.Unmarshal(r.Body, &v); err != nil {
		return nil, err
	}

	return v[0], nil
}

func (s *STNS) GetGroupByID(id int) (*model.Group, error) {
	r, err := s.client.Request(groupsEndpoint, fmt.Sprintf("id=%d", id))
	if err != nil {
		return nil, err
	}
	v := []*model.Group{}
	if err := json.Unmarshal(r.Body, &v); err != nil {
		return nil, err
	}

	return v[0], nil
}
