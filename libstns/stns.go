package libstns

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/STNS/STNS/v2/model"
	"github.com/caarlos0/env"
	"golang.org/x/crypto/ssh"
)

type STNS struct {
	client             *client
	opt                *Options
	makeChallengeCode  func() ([]byte, error)
	storeChallengeCode func(string, []byte) error
	popChallengeCode   func(string) ([]byte, error)
}

// user must be a single path element so the code stays directly under os.TempDir().
func challengeCodePath(user string) (string, error) {
	if user == "" || user == "." || user == ".." || strings.ContainsAny(user, "/\\\x00") {
		return "", errors.New("invalid user name")
	}
	return filepath.Join(os.TempDir(), user), nil
}

func DefaultStoreChallengeCode(user string, code []byte) error {
	p, err := challengeCodePath(user)
	if err != nil {
		return err
	}
	if err := os.Remove(p); err != nil && !os.IsNotExist(err) {
		return err
	}
	f, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return err
	}
	if _, err := f.Write(code); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}

func DefaultPopChallengeCode(user string) ([]byte, error) {
	p, err := challengeCodePath(user)
	if err != nil {
		return nil, err
	}
	fi, err := os.Lstat(p)
	if err != nil {
		return nil, err
	}
	if !fi.Mode().IsRegular() || !isOwnedByCurrentUser(fi) {
		return nil, errors.New("invalid challenge code file")
	}
	defer os.Remove(p)
	return os.ReadFile(p)
}

func DefaultMakeChallengeCode() ([]byte, error) {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return nil, errors.New("rand read error")
	}

	var bs []byte
	for _, v := range b {
		bs = append(bs, letters[int(v)%len(letters)])
	}
	return bs, nil

}

type Options struct {
	AuthToken          string `env:"STNS_AUTH_TOKEN"`
	User               string `env:"STNS_USER"`
	Password           string `env:"STNS_PASSWORD"`
	UserAgent          string
	SkipSSLVerify      bool `env:"STNS_SKIP_VERIFY"`
	HttpProxy          string
	HttpKeepalive      bool `env:"STNS_HTTP_KEEPALIVE"`
	RequestTimeout     int  `env:"STNS_REQUEST_TIMEOUT"`
	RequestRetry       int  `env:"STNS_REQUEST_RETRY"`
	HttpHeaders        map[string]string
	TLS                TLS
	PrivatekeyPath     string `env:"STNS_PRIVATE_KEY"`
	PrivatekeyPassword string `env:"STNS_PRIVATE_KEY_PASSWORD"`
}

func NewSTNS(endpoint string, opt *Options) (*STNS, error) {
	if opt == nil {
		opt = &Options{}
	}
	s := &STNS{
		popChallengeCode:   DefaultPopChallengeCode,
		storeChallengeCode: DefaultStoreChallengeCode,
		makeChallengeCode:  DefaultMakeChallengeCode,
	}
	if err := env.Parse(s); err != nil {
		return nil, err
	}
	c, err := newClient(endpoint, opt)
	if err != nil {
		return nil, err
	}
	s.client = c
	s.opt = opt
	return s, nil
}

const usersEndpoint = "/users"
const groupsEndpoint = "/groups"

func (s *STNS) SetStoreChallengeCode(f func(string, []byte) error) {
	s.storeChallengeCode = f
}

func (s *STNS) SetPopChallengeCode(f func(string) ([]byte, error)) {
	s.popChallengeCode = f
}

func (s *STNS) Request(path, query string) (*Response, error) {
	return s.client.Request(path, query)
}

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
	r, err := s.client.Request(usersEndpoint, url.Values{"name": {name}}.Encode())
	if err != nil {
		return nil, err
	}
	v := []*model.User{}
	if err := json.Unmarshal(r.Body, &v); err != nil {
		return nil, err
	}

	if len(v) == 0 {
		return nil, errors.New("user not found")
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

	if len(v) == 0 {
		return nil, errors.New("user not found")
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
	r, err := s.client.Request(groupsEndpoint, url.Values{"name": {name}}.Encode())
	if err != nil {
		return nil, err
	}
	v := []*model.Group{}
	if err := json.Unmarshal(r.Body, &v); err != nil {
		return nil, err
	}

	if len(v) == 0 {
		return nil, errors.New("group not found")
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

	if len(v) == 0 {
		return nil, errors.New("group not found")
	}

	return v[0], nil
}
func (c *STNS) CreateUserChallengeCode(name string) ([]byte, error) {
	code, err := c.makeChallengeCode()
	if err != nil {
		return nil, err
	}
	err = c.storeChallengeCode(name, code)
	if err != nil {
		return nil, err
	}
	return code, nil
}

func (c *STNS) PopUserChallengeCode(name string) ([]byte, error) {
	return c.popChallengeCode(name)
}

func (c *STNS) Sign(code []byte) ([]byte, error) {
	privateKey, err := c.loadPrivateKey()
	if err != nil {
		return nil, err
	}
	sig, err := privateKey.Sign(rand.Reader, code)
	if err != nil {
		return nil, err
	}

	jsonSig, err := json.Marshal(sig)
	if err != nil {
		return nil, err
	}

	return jsonSig, nil
}

func (c *STNS) VerifyWithUser(name string, msg, signature []byte) error {
	user, err := c.GetUserByName(name)
	if err != nil {
		return err
	}

	return c.Verify(msg, []byte(strings.Join(user.Keys, "\n")), signature)
}

func (c *STNS) Verify(msg, publicKeyBytes, signature []byte) error {
	for len(publicKeyBytes) > 0 {
		publicKey, _, _, rest, err := ssh.ParseAuthorizedKey(publicKeyBytes)
		if err != nil {
			return fmt.Errorf("can't read public key %s", err.Error())
		}

		var sig ssh.Signature
		if err := json.Unmarshal(signature, &sig); err != nil {
			return err
		}

		if err := publicKey.Verify(msg, &sig); err == nil {
			return nil
		}
		publicKeyBytes = rest
	}
	return errors.New("verify failed")

}

func (c *STNS) loadPrivateKey() (ssh.Signer, error) {
	usr, _ := user.Current()
	priv, err := ioutil.ReadFile(strings.Replace(c.opt.PrivatekeyPath, "~", usr.HomeDir, 1))
	if err != nil {
		return nil, fmt.Errorf("error:%s path:%s", err.Error(), c.opt.PrivatekeyPath)
	}

	if c.opt.PrivatekeyPassword != "" {
		return ssh.ParsePrivateKeyWithPassphrase(priv, []byte(c.opt.PrivatekeyPassword))
	}
	return ssh.ParsePrivateKey(priv)

}
