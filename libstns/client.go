package libstns

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/caarlos0/env"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/sirupsen/logrus"
	"github.com/thoas/go-funk"
)

var version = "0.0.1"

type TLS struct {
	CA   string
	Cert string
	Key  string
}

var DefaultTimeout = 15
var DefaultRetry = 3

type client struct {
	ApiEndpoint string
	opt         *Options
}

type Response struct {
	StatusCode int
	Headers    map[string]string
	Body       []byte
}

func newClient(endpoint string, opt *Options) (*client, error) {

	if err := env.Parse(opt); err != nil {
		return nil, err
	}
	if opt.UserAgent == "" {
		opt.UserAgent = fmt.Sprintf("%s/%s", "libstns-go", version)
	}

	if opt.RequestTimeout == 0 {
		opt.RequestTimeout = DefaultTimeout
	}

	if opt.RequestRetry == 0 {
		opt.RequestRetry = DefaultRetry
	}

	return &client{
		ApiEndpoint: endpoint,
		opt:         opt,
	}, nil
}
func (h *client) RequestURL(requestPath, query string) (*url.URL, error) {
	u, err := url.Parse(h.ApiEndpoint)
	if err != nil {
		return nil, err
	}

	u.Path = path.Join(u.Path, requestPath)
	u.RawQuery = query
	return u, nil

}

func (h *client) Request(path, query string) (*Response, error) {
	supportHeaders := []string{
		"user-highest-id",
		"user-lowest-id",
		"group-highest-id",
		"group-lowest-id",
	}

	u, err := h.RequestURL(path, query)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		logrus.Errorf("make http request error:%s", err.Error())
		return nil, err
	}

	h.setHeaders(req)
	h.setBasicAuth(req)

	tr := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: time.Duration(h.opt.RequestTimeout) * time.Second,
		}).Dial,
	}
	if strings.Index(h.ApiEndpoint, "https") == 0 {
		tc, err := h.tlsConfig()
		if err != nil {
			logrus.Errorf("make tls config error:%s", err.Error())
			return nil, err
		}

		tr.TLSClientConfig = tc
	}

	tr.Proxy = http.ProxyFromEnvironment
	if h.opt.HttpProxy != "" {
		proxyUrl, err := url.Parse(h.opt.HttpProxy)
		if err == nil {
			tr.Proxy = http.ProxyURL(proxyUrl)
		}
	}
	retryclient := retryablehttp.NewClient()
	retryclient.RetryMax = h.opt.RequestRetry

	client := retryclient.StandardClient()
	client.Transport = tr
	resp, err := client.Do(req)
	if err != nil {
		logrus.Errorf("http request error:%s", err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	headers := map[string]string{}
	for k, v := range resp.Header {
		if funk.ContainsString(supportHeaders, strings.ToLower(k)) {
			headers[k] = v[0]
		}
	}

	switch resp.StatusCode {
	case http.StatusOK:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		r := Response{
			StatusCode: resp.StatusCode,
			Body:       body,
			Headers:    headers,
		}

		return &r, nil
	default:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("status code=%d, body=%s", resp.StatusCode, string(body))
	}
}

func (h *client) setHeaders(req *http.Request) {
	if len(h.opt.HttpHeaders) > 0 {
		for k, v := range h.opt.HttpHeaders {
			req.Header.Add(k, v)
		}
	}

	req.Header.Set("User-Agent", h.opt.UserAgent)

	if h.opt.AuthToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("token %s", h.opt.AuthToken))
	}
}

func (h *client) setBasicAuth(req *http.Request) {
	if h.opt.User != "" && h.opt.Password != "" {
		req.SetBasicAuth(h.opt.User, h.opt.Password)
	}
}

func (h *client) tlsConfig() (*tls.Config, error) {
	tlsConfig := &tls.Config{InsecureSkipVerify: h.opt.SkipSSLVerify}
	if h.opt.TLS.CA != "" {
		CA_Pool := x509.NewCertPool()

		severCert, err := ioutil.ReadFile(h.opt.TLS.CA)
		if err != nil {
			return nil, err
		}
		CA_Pool.AppendCertsFromPEM(severCert)

		tlsConfig.RootCAs = CA_Pool
	}

	if h.opt.TLS.Cert != "" && h.opt.TLS.Key != "" {
		x509Cert, err := tls.LoadX509KeyPair(h.opt.TLS.Cert, h.opt.TLS.Key)
		if err != nil {
			return nil, err
		}
		tlsConfig.Certificates = make([]tls.Certificate, 1)
		tlsConfig.Certificates[0] = x509Cert
	}

	if len(tlsConfig.Certificates) == 0 && tlsConfig.RootCAs == nil {
		tlsConfig = nil
	}

	return tlsConfig, nil
}
