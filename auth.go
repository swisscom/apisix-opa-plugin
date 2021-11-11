package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	pkgHTTP "github.com/apache/apisix-go-plugin-runner/pkg/http"
	"github.com/apache/apisix-go-plugin-runner/pkg/log"
	"github.com/apache/apisix-go-plugin-runner/pkg/plugin"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type OpaAuth struct {
	client *http.Client
}

type OpaAuthConf struct {
	RulePath string `json:"rule_path"`         // e.g: example.authz/allow
	OpaUrl   string `json:"opa_url,omitempty"` // overwrites OPA_URL
}

func (p *OpaAuth) Name() string {
	return "opa"
}

func (p *OpaAuth) ParseConf(in []byte) (interface{}, error) {
	conf := OpaAuthConf{}
	err := json.Unmarshal(in, &conf)
	if err != nil {
		return nil, err
	}

	if conf.OpaUrl != "" {
		// Opa URL is not empty
		_, err := url.Parse(conf.OpaUrl)
		if err != nil {
			return nil, fmt.Errorf("invalid OPA URL provided: %v", err)
		}
	} else {
		envOpaUrl := os.Getenv("OPA_URL")
		if envOpaUrl != "" {
			_, err := url.Parse(envOpaUrl)
			if err != nil {
				return nil, fmt.Errorf("invalid OPA URL provided: %v", err)
			}
			conf.OpaUrl = envOpaUrl
		} else {
			log.Warnf("no opa_url or OPA_URL defined, falling back to default")
			conf.OpaUrl = "http://127.0.0.1:8181"
		}
	}

	return conf, err
}

/*
	Filter sends the request to OPA on /v1/data/{{ p.RulePath | replace('.', '/') }}
	and proceeds with the request on success or aborts it and returns 403 on failure.
*/
func (p *OpaAuth) Filter(conf interface{}, w http.ResponseWriter, r pkgHTTP.Request) {
	opaConf := conf.(OpaAuthConf)

	if p.client == nil {
		p.initClient()
	}

	opaInputReq := OpaInputRequest{
		Host:    valueOrDefault(r, "server_name", ""),
		SrcIp:   r.SrcIP(),
		Path:    string(r.Path()),
		Headers: r.Header().View(),
	}
	input, err := json.Marshal(&opaInputReq)
	if err != nil {
		log.Errorf("unable to marshal JSON: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequest(http.MethodPost,
		fmt.Sprintf(
			"%s/v1/data/%s",
			opaConf.OpaUrl,
			FormatRulePathUrl(opaConf.RulePath),
		),
		bytes.NewReader(input),
	)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := p.client.Do(req)
	if err != nil {
		log.Errorf("unable to evaluate OPA policy: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if res.StatusCode != http.StatusOK {
		log.Errorf("invalid status code received from OPA: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Parse response from OPA
	var result Result
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&result)

	if err != nil {
		log.Errorf("unable to decode JSON response from OPA: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if result.StatusCode != 0 {
		w.WriteHeader(int(result.StatusCode))
		if result.Reason != "" {
			_, err = w.Write([]byte(result.Reason))
			if err != nil {
				log.Errorf("unable to write reason: %v", err)
			}
		}
	}
}

func valueOrDefault(r pkgHTTP.Request, s string, def string) string {
	v, err := r.Var(s)
	if err != nil {
		return def
	}
	return string(v)
}

/*
	FormatRulePathUrl given a rulePath returns a path to be used
	in an OPA /v1/data/path-goes-here request.
*/

func FormatRulePathUrl(path string) interface{} {
	return strings.Replace(path, ".", "/", -1)
}

func (p *OpaAuth) initClient() {
	c := http.DefaultClient
	c.Timeout = time.Second * 10
	p.client = c
}

func init() {
	err := plugin.RegisterPlugin(&OpaAuth{})
	if err != nil {
		log.Fatalf("failed to register OPA plugin: %v", err)
	}
}
