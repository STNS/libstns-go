package libstns

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/user"
	"strings"

	"github.com/STNS/STNS/v2/model"
	"github.com/caarlos0/env"
	"golang.org/x/crypto/ssh"
)

type STNS struct {
	client             *client
	PrivatekeyPath     string `env:"STNS_PRIVATE_KEY" envDefault:"~/.ssh/id_rsa"`
	PrivatekeyPassword string `env:"STNS_PRIVATE_KEY_PASSWORD"`
}

func NewSTNS(endpoint string, opt *ClientOptions) (*STNS, error) {
	s := &STNS{}
	if err := env.Parse(s); err != nil {
		return nil, err
	}
	c, err := newClient(endpoint, opt)
	if err != nil {
		return nil, err
	}
	s.client = c
	return s, nil
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

func (c *STNS) Sign(msg []byte) (*ssh.Signature, error) {
	privateKey, err := c.loadPrivateKey()
	if err != nil {
		return nil, err
	}
	return privateKey.Sign(rand.Reader, msg)
}

func (c *STNS) VerifyWithUser(name string, msg, signature []byte) error {
	user, err := c.GetUserByName(name)
	if err != nil {
		return err
	}
	return c.Verify(msg, []byte(strings.Join(user.Keys, "\n")), signature)
}

func (c *STNS) Verify(msg, publicKeyBytes, signature []byte) error {
	publicKey, _, _, _, err := ssh.ParseAuthorizedKey(publicKeyBytes)
	if err != nil {
		return fmt.Errorf("can't read public key %s", err.Error())
	}

	var sig ssh.Signature
	if err := json.Unmarshal(signature, &sig); err != nil {
		return err
	}
	return publicKey.Verify(msg, &sig)
}

func (c *STNS) loadPrivateKey() (ssh.Signer, error) {
	usr, _ := user.Current()
	priv, err := ioutil.ReadFile(strings.Replace(c.PrivatekeyPath, "~", usr.HomeDir, 1))
	if err != nil {
		return nil, fmt.Errorf("error:%s path:%s", err.Error(), c.PrivatekeyPath)
	}

	if c.PrivatekeyPassword != "" {
		return ssh.ParsePrivateKeyWithPassphrase(priv, []byte(c.PrivatekeyPassword))
	}
	return ssh.ParsePrivateKey(priv)

}
