package libstns

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/sirupsen/logrus"
	"github.com/thoas/go-funk"
)

var version string

type TLS struct {
	CA   string
	Cert string
	Key  string
}

var DefaultTimeout = 15
var DefaultRetry = 3

type HttpOptions struct {
	AuthToken      string
	User           string
	Password       string
	UserAgent      string
	SkipSSLVerify  bool
	HttpProxy      string
	RequestTimeout int
	RequestRetry   int
	HttpHeaders    map[string]string
	TLS            TLS
}
type Http struct {
	ApiEndpoint string
	opt         *HttpOptions
}

type Response struct {
	StatusCode int
	Headers    map[string]string
	Body       []byte
}

func NewHttp(endpoint string, opt *HttpOptions) *Http {
	if opt.UserAgent == "" {
		opt.UserAgent = "libstns-go"
	}

	if opt.RequestTimeout == 0 {
		opt.RequestTimeout = DefaultTimeout
	}

	if opt.RequestRetry == 0 {
		opt.RequestRetry = DefaultRetry
	}

	return &Http{
		ApiEndpoint: endpoint,
		opt:         opt,
	}
}

func (h *Http) Request(path string) (*Response, error) {
	supportHeaders := []string{
		"user-highest-id",
		"user-lowest-id",
		"group-highest-id",
		"group-lowest-id",
	}

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		logrus.Errorf("make http request error:%s", err.Error())
		return nil, err
	}

	h.setHeaders(req)
	h.setBasicAuth(req)

	tc, err := h.tlsConfig()
	if err != nil {
		logrus.Errorf("make tls config error:%s", err.Error())
		return nil, err
	}

	tr := &http.Transport{
		TLSClientConfig: tc,
		Dial: (&net.Dialer{
			Timeout: time.Duration(h.opt.RequestTimeout) * time.Second,
		}).Dial,
	}

	tr.Proxy = http.ProxyFromEnvironment
	if h.opt.HttpProxy != "" {
		proxyUrl, err := url.Parse(h.opt.HttpProxy)
		if err == nil {
			tr.Proxy = http.ProxyURL(proxyUrl)
		}
	}
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = h.opt.RequestRetry

	client := retryClient.StandardClient()
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
		r := Response{
			StatusCode: resp.StatusCode,
			Headers:    headers,
		}
		return &r, nil
	}
}

func (h *Http) setHeaders(req *http.Request) {
	for k, v := range h.opt.HttpHeaders {
		req.Header.Add(k, v)
	}
	req.Header.Set("User-Agent", fmt.Sprintf("%s/%s", h.opt.UserAgent, version))

	if h.opt.AuthToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("token %s", h.opt.AuthToken))
	}
}

func (h *Http) setBasicAuth(req *http.Request) {
	if h.opt.User != "" && h.opt.Password != "" {
		req.SetBasicAuth(h.opt.User, h.opt.Password)
	}
}

func (h *Http) tlsConfig() (*tls.Config, error) {
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
