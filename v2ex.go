package v2ex

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const APIBaseURL string = "https://www.v2ex.com/api/v2/"

type Client struct {
	token string
}

type GetNodeResult struct {
	Avatar       string // `json:avatar`
	Created      int    // `json:created`
	Footer       string // `json:footer`
	Header       string // `json:header`
	Id           int    // `json:id`
	LastModified int    // `json:last_modified`
	Name         string // `json:name`
	Title        string // `json:title`
	Topics       int    // `json:topics`
	Url          string // `json:url`
}

type GetNodeResponse struct {
	Message string        // `json:message`
	Result  GetNodeResult // `json:result`
	Success bool          // `json:success`
}

func (c Client) request(method string, path string, params map[string]string, data map[string]any) (*[]byte, error) {
	jsonBody, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	url, err := url.JoinPath(APIBaseURL, path)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	queryParams := req.URL.Query()
	for k, v := range params {
		queryParams.Add(k, v)
	}
	req.URL.RawQuery = queryParams.Encode()

	req.Header.Set("Authorization", "Bearer "+c.token)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &body, nil
}

func (c Client) GetNode(nodeName string) (GetNodeResponse, error) {
	var resp GetNodeResponse
	resp_body, err := c.request("GET", fmt.Sprintf("/nodes/%s", nodeName), nil, nil)
	if err != nil {
		return resp, err
	}
	if err := json.Unmarshal(*resp_body, &resp); err != nil {
		return resp, err
	}
	return resp, nil
}

func (c Client) GetNodeTopics(nodeName string, page int) {}
func (c Client) GetTopic(topicID int)                    {}
func (c Client) GetTopicReplies(topicID int, page int)   {}
