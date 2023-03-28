package v2ex

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type TokenScope string
type TokenExpiration int

type Client struct {
	token string
}

// V2EX API 接口返回值公共参数 struct
type v2exResponse struct {
	Message string           `json:"message"`
	Result  *json.RawMessage `json:"result"`
	Success bool             `json:"success"`
}

func (r *v2exResponse) Error() string {
	return fmt.Sprintf("API error, message: %s", r.Message)
}

type v2exNode struct {
	Avatar       string `json:"avatar"`
	Created      int    `json:"created"`
	Footer       string `json:"footer"`
	Header       string `json:"header"`
	Id           int    `json:"id"`
	LastModified int    `json:"last_modified"`
	Name         string `json:"name"`
	Title        string `json:"title"`
	Topics       int    `json:"topics"`
	Url          string `json:"url"`
}

type v2exTopic struct {
	Content         string `json:"content"`
	ContentRendered string `json:"content_rendered"`
	Created         int    `json:"created"`
	Id              int    `json:"id"`
	LastModified    int    `json:"last_modified"`
	LastReplyBy     string `json:"last_reply_by"`
	LastTouched     int    `json:"last_touched"`
	Replies         int    `json:"replies"`
	Syntax          int    `json:"syntax"`
	Title           string `json:"title"`
	Url             string `json:"url"`
}

type v2exMember struct {
	Avatar   string `json:"avatar"`
	Bio      string `json:"bio"`
	Created  int    `json:"created"`
	Github   string `json:"github"`
	Id       int    `json:"id"`
	Url      string `json:"url"`
	Username string `json:"username"`
	Website  string `json:"website"`
}

type v2exSupplement struct {
	Content          string `json:"content"`
	Content_Rendered string `json:"content_rendered"`
	Created          int    `json:"created"`
	Id               int    `json:"id"`
	Syntax           int    `json:"syntax"`
}

type v2exReply struct {
	Content         string     `json:"content"`
	ContentRendered string     `json:"content_rendered"`
	Created         int        `json:"created"`
	Id              int        `json:"id"`
	Member          v2exMember `json:"member"`
}

type v2exToken struct {
	Created     int    `json:"created"`
	Expiration  int    `json:"expiration"`
	GoodForDays int    `json:"good_for_days"`
	LastUsed    int    `json:"last_used"`
	Scope       string `json:"scope"`
	Token       string `json:"token"`
	TotalUsed   int    `json:"total_used"`
}

type v2exSelfProfile struct {
	AvatarLarge  string `json:"avatar_large"`
	AvatarMini   string `json:"avatar_mini"`
	AvatarNormal string `json:"avatar_normal"`
	Bio          string `json:"bio"`
	Btc          string `json:"btc"`
	Created      int    `json:"created"`
	Github       string `json:"github"`
	Id           int    `json:"id"`
	LastModified int    `json:"last_modified"`
	Location     string `json:"location"`
	Psn          string `json:"psn"`
	Tagline      string `json:"tagline"`
	Twitter      string `json:"twitter"`
	Url          string `json:"url"`
	Username     string `json:"username"`
	Website      string `json:"website"`
}

type v2exNotification struct {
	Created     int `json:"created"`
	ForMemberId int `json:"for_member_id"`
	Id          int `json:"id"`
	Member      struct {
		Username string `json:"username"`
	} `json:"member"`
	MemberId        int    `json:"member_id"`
	Payload         string `json:"payload"`
	PayloadRendered string `json:"payload_rendered"`
	Text            string `json:"text"`
}

type getTopicResult struct {
	v2exTopic
	Member      v2exMember       `json:"member"`
	Node        v2exNode         `json:"node"`
	Supplements []v2exSupplement `json:"supplements"`
}

type createTokenResult struct {
	Token string `json:"token"`
}

func (c Client) request(method string, path string, params map[string]string, data map[string]interface{}, result interface{}) error {
	jsonBody, err := json.Marshal(data)
	if err != nil {
		return err
	}

	url, err := url.JoinPath(baseURL, path)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
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
		return err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var response v2exResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return err
	}

	if !response.Success {
		return &response
	}

	if result != nil {
		if err := json.Unmarshal(*response.Result, result); err != nil {
			return err
		}
	}

	return nil
}

// 获取指定节点
func (c Client) GetNode(nodeName string) (result v2exNode, err error) {
	err = c.request("GET", fmt.Sprintf("/nodes/%s", nodeName), nil, nil, &result)
	return
}

// 获取指定节点下的主题
func (c Client) GetNodeTopics(nodeName string, page int) (result []v2exTopic, err error) {
	err = c.request("GET", fmt.Sprintf("/nodes/%s/topics", nodeName), map[string]string{"p": strconv.Itoa(page)}, nil, &result)
	return
}

// 获取指定主题
func (c Client) GetTopic(topicID int) (result getTopicResult, err error) {
	err = c.request("GET", fmt.Sprintf("/topics/%d", topicID), nil, nil, &result)
	return
}

// 获取指定主题下的回复
func (c Client) GetTopicReplies(topicID int, page int) (result []v2exReply, err error) {
	err = c.request("GET", fmt.Sprintf("/topics/%d/replies", topicID), map[string]string{"p": strconv.Itoa(page)}, nil, &result)
	return
}

// 查看当前使用的令牌
func (c Client) GetToken() (result v2exToken, err error) {
	err = c.request("GET", "/token", nil, nil, &result)
	return
}

// 获取自己的 Profile
func (c Client) GetSelfProfile() (result v2exSelfProfile, err error) {
	err = c.request("GET", "/member", nil, nil, &result)
	return
}

// 获取最新的提醒
func (c Client) GetNotifications(page int) (result []v2exNotification, err error) {
	err = c.request("GET", "/notifications", map[string]string{"p": strconv.Itoa(page)}, nil, &result)
	return
}

// 创建新的令牌
func (c Client) CreateToken(scope TokenScope, expiration TokenExpiration) (result createTokenResult, err error) {
	err = c.request("POST", "/tokens", nil, map[string]interface{}{"scope": scope, "expiration": expiration}, &result)
	return
}

// 删除指定的提醒
func (c Client) DeleteNotification(notificationId int) (err error) {
	err = c.request("DELETE", fmt.Sprintf("/notifications/%d", notificationId), nil, nil, nil)
	return
}

// Create V2EX API Client
func NewClient(token string) Client {
	return Client{token}
}
