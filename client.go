package hubspot

import (
	"fmt"
	"bytes"
	"io/ioutil"
	"net/http"
	"log"
	"net/url"
)

type Client struct {
	HAPIKey string
	Token   string
}

const ContentType = "application/json"
const baseUrl = "https://api.hubapi.com"

func NewHAPIClient(api_key string) *Client {
	return &Client{
		HAPIKey: api_key,
	}
}

func NewTokenClient(token string) *Client {
	return &Client{
		Token: token,
	}
}

func (c *Client) doRequest(endpoint string, http_type string, data []byte, params ...map[string]interface{}) ([]byte, error) {
	url := c.collectURIForRequest(endpoint, params)
	req, err := http.NewRequest(http_type, url, bytes.NewBuffer(data))
	c.setDefaultHeaders(req)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("HUBSPOT ERROR: ", err)
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if 200 != resp.StatusCode {
		return nil, fmt.Errorf("%s", body)
	}
	return body, nil
}

func (c *Client) setDefaultHeaders(req *http.Request) {
	req.Header.Set("Accept", ContentType)
	req.Header.Set("Content-Type", ContentType)
	if len(c.Token) > 0 {
		req.Header.Set("Authorization", "Bearer "+c.HAPIKey)
	}
}

func (c *Client) collectURIForRequest(endpoint string, params []map[string]interface{}) string {
	endpoint_url := fmt.Sprintf(baseUrl+"/%s", endpoint)
	u, err := url.Parse(endpoint_url)
	if err != nil {
		log.Fatal(err)
	}
	q := u.Query()
	if len(c.HAPIKey) > 0 {
		q.Set("hapikey", c.HAPIKey)
	}
	if len(params) > 0 {
		c.paramsToQuery(params[0], q)
	}
	u.RawQuery = q.Encode()
	return u.String()
}

func (c *Client) paramsToQuery(params map[string]interface{}, q url.Values) error {
	for key, value := range params {
		switch v := value.(type) {
		case int:
			q.Set(key, fmt.Sprintf("%v", v))
			break
		case string:
			q.Set(key, v)
			break
		case []string:
			for _, item := range v {
				q.Add(key, item)
			}
			break
		}
	}
	return nil
}
