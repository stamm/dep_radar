package http

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	stdhttp "net/http"
	"strings"
	"time"

	i "github.com/stamm/dep_radar/src/interfaces"
)

var (
	_                 i.IWebClient = &Client{}
	defaultHTTPClient *stdhttp.Client
)

func init() {
	defaultHTTPClient = &stdhttp.Client{
		Timeout: 30 * time.Second,
		Transport: &stdhttp.Transport{
			MaxIdleConns:          100,
			MaxIdleConnsPerHost:   100,
			IdleConnTimeout:       30 * time.Second,
			ResponseHeaderTimeout: 30 * time.Second,
		},
	}
}

// Options for http client
type Options struct {
	URL      string
	User     string
	Password string
}

// Client gets html pages
type Client struct {
	Options    Options
	Limit      int
	httpClient *stdhttp.Client
	limitCh    chan struct{}
}

// NewClient returns our http client
func NewClient(op Options, limit int) *Client {
	return &Client{
		Options: op,
		Limit:   limit,
		limitCh: make(chan struct{}, limit),
	}
}

// Get the html
func (c *Client) Get(ctx context.Context, uri string) ([]byte, error) {
	log.Printf("Start getting url %s\n", uri)
	c.limitCh <- struct{}{}
	defer func() {
		<-c.limitCh
	}()
	url := uri
	if c.Options.URL != "" {
		url = strings.Trim(c.Options.URL, "/") + "/" + uri
	}
	// fmt.Printf("url = %+v\n", url)
	req, err := stdhttp.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/json")
	if c.Options.User != "" && c.Options.Password != "" {
		// fmt.Println("USER!")
		req.SetBasicAuth(c.Options.User, c.Options.Password)
	}

	start := time.Now()
	resp, err := c.getHTTPClient().Do(req)
	log.Printf("time %s for %s", time.Since(start), url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("StatusCode is not 200: %d", resp.StatusCode)
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return buf, err
	}
	// fmt.Printf("result %s", string(buf))
	return buf, err
}

func (c Client) getHTTPClient() *stdhttp.Client {
	if c.httpClient != nil {
		return c.httpClient
	}
	c.httpClient = defaultHTTPClient
	return c.httpClient
}
