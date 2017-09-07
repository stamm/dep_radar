package http

import (
	"fmt"
	"io/ioutil"
	stdhttp "net/http"
	"strings"
	"time"

	i "github.com/stamm/dep_radar/interfaces"
)

var (
	_                 i.IWebClient = &Client{}
	defaultHttpClient *stdhttp.Client
)

func init() {
	tr := &stdhttp.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:     30 * time.Second,
	}
	defaultHttpClient = &stdhttp.Client{
		Timeout:   1 * time.Second,
		Transport: tr,
	}
}

type Options struct {
	URL      string
	User     string
	Password string
}

type Client struct {
	Options    Options
	Limit      int
	httpClient *stdhttp.Client
	limitCh    chan struct{}
}

func NewClient(op Options, limit int) *Client {
	return &Client{
		Options: op,
		Limit:   limit,
		limitCh: make(chan struct{}, limit),
	}
}

func (r *Client) Get2(url string) ([]byte, error) {
	client := &stdhttp.Client{}
	req, _ := stdhttp.NewRequest("GET", url, nil)

	resp, err := client.Do(req)
	if nil != err {
		return nil, err
	}
	if resp.StatusCode > 300 {
		return nil, fmt.Errorf("Response code is %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (c *Client) Get(uri string) ([]byte, error) {
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

	req.Header.Set("Content-Type", "application/json")
	if c.Options.User != "" && c.Options.Password != "" {
		// fmt.Println("USER!")
		req.SetBasicAuth(c.Options.User, c.Options.Password)
	}

	client := c.getHTTPClient()
	resp, err := client.Do(req)
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
	c.httpClient = defaultHttpClient
	return c.httpClient
}
