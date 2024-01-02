package kanthor

import (
	_ "embed"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/scrapnode/kanthor/clients/sdk-go/internal/openapi"
)

//go:embed project.json
var project []byte

//go:embed project.host
var host string

type Project struct {
	Version string `json:"version"`
	Hosts   map[string]string
}

type Options struct {
	Debug  bool
	Scheme string
	Host   string
}

func New(credentials string, opts *Options) (*Kanthor, error) {
	return NewWithOptions(credentials, nil)
}

func NewWithOptions(credentials string, opts *Options) (*Kanthor, error) {
	var proj Project
	if err := json.Unmarshal(project, &proj); err != nil {
		return nil, err
	}

	conf := openapi.NewConfiguration()
	conf.Scheme = "https"
	conf.Middleware = func(r *http.Request) {
		r.Header.Set("Idempotency-Key", uuid.NewString())
	}

	h, err := parse(credentials, &proj)
	if err != nil {
		return nil, err
	}
	conf.Host = h

	// forllowing https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/User-Agent#syntax
	conf.UserAgent = fmt.Sprintf("kanthor/%s OpenAPI/go", proj.Version)
	conf.AddDefaultHeader("Authorization", fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(credentials))))

	// override configuration with custom options'
	if opts != nil {
		conf.Debug = opts.Debug
		if opts.Host != "" {
			conf.Host = opts.Host
			// if the host is localhost, should change the scheme back to HTTP
			conf.Scheme = scheme(conf.Host)
		}
		if opts.Scheme != "" {
			conf.Scheme = opts.Scheme
		}
	}

	api := openapi.NewAPIClient(conf)
	sdk := &Kanthor{
		Message: &Message{api: api},
	}
	return sdk, nil
}

func parse(credentials string, proj *Project) (string, error) {
	segments := strings.Split(credentials, ":")
	if len(segments) != 2 || len(segments[1]) == 0 {
		return "", errors.New("malform credentials")
	}

	parts := strings.Split(segments[1], ".")
	if len(parts) == 2 {
		if h, exist := proj.Hosts[parts[0]]; exist {
			return h, nil
		}
	}

	return host, nil
}

func scheme(host string) string {
	if strings.HasPrefix(host, "localhost") {
		return "http"
	}
	if strings.HasPrefix(host, "127.0.0.1") {
		return "http"
	}
	return "https"
}

type Kanthor struct {
	Account      *Account
	Application  *Application
	Endpoint     *Endpoint
	EndpointRule *EndpointRule
	Message      *Message
}
