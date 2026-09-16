//go:build e2e

// Package client 提供黑盒 E2E 测试使用的最小 HTTP 客户端：统一注入鉴权头、编码请求 JSON，
// 并把网关统一错误体 {"error":{...}} 解码为 gRPC 状态码，供用例按错误语义断言。
//
// 该包只依赖标准库。E2E 测的是对外 HTTP 契约，禁止 import 本服务的业务包
// （api/、handler/、internal/、modeler/），详见 test/e2e/README.md。
package client

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Client 是面向单个调用主体的 HTTP 客户端；凭据变更时通过 WithToken/WithBasicAuth 派生新实例。
type Client struct {
	baseURL string
	http    *http.Client
	// authorization 为完整的 Authorization 头值，如 "Basic dXNlcjE6..." 或 "Bearer xxx"；空则匿名。
	authorization string
}

// Option 是 Client 的构造选项。
type Option func(*Client)

// WithBaseURL 设置服务基础地址，如 http://127.0.0.1:8080。
func WithBaseURL(baseURL string) Option {
	return func(c *Client) { c.baseURL = strings.TrimRight(baseURL, "/") }
}

// WithBasicAuth 设置 HTTP Basic 凭据，对应配置中的 security.authentication.http_users。
func WithBasicAuth(username, password string) Option {
	return func(c *Client) { c.authorization = basicHeader(username, password) }
}

// WithToken 设置 Bearer 访问令牌，对应配置中的 security.authentication.oidc_provider。
func WithToken(token string) Option {
	return func(c *Client) { c.authorization = "Bearer " + token }
}

// WithTimeout 设置单次请求超时，默认 30 秒。
func WithTimeout(d time.Duration) Option {
	return func(c *Client) { c.http.Timeout = d }
}

// New 构造 Client，默认 30 秒超时。
func New(opts ...Option) *Client {
	c := &Client{http: &http.Client{Timeout: 30 * time.Second}}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// WithBasicAuth 返回携带新 Basic 凭据的 Client 副本，用于多主体用例。
func (c *Client) WithBasicAuth(username, password string) *Client {
	next := *c
	next.authorization = basicHeader(username, password)
	return &next
}

// WithToken 返回携带新令牌的 Client 副本，用于多主体用例。
func (c *Client) WithToken(token string) *Client {
	next := *c
	next.authorization = "Bearer " + token
	return &next
}

// Anonymous 返回去掉全部凭据的 Client 副本，用于鉴权负例。
func (c *Client) Anonymous() *Client {
	next := *c
	next.authorization = ""
	return &next
}

// BaseURL 返回该客户端的服务基础地址。
func (c *Client) BaseURL() string { return c.baseURL }

func basicHeader(username, password string) string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+password))
}

// Response 保留原始状态码、响应头与响应体；用例按需解码，不强制结构，
// 以便覆盖「零值字段可能省略」这类契约宽容性。
type Response struct {
	StatusCode int
	Header     http.Header
	Body       []byte
}

// errorBody 对应网关统一错误结构 {"error":{"code","status","message","details"}}。
type errorBody struct {
	Error struct {
		Code    int    `json:"code"`
		Status  string `json:"status"`
		Message string `json:"message"`
	} `json:"error"`
}

// Decode 把响应体解码到 v。
func (r *Response) Decode(v any) error {
	return json.Unmarshal(r.Body, v)
}

// IsOK 报告响应是否为 2xx。
func (r *Response) IsOK() bool { return r.StatusCode >= 200 && r.StatusCode < 300 }

// RPCCode 返回错误体中的 gRPC 状态字符串（如 "NotFound"）；成功响应返回空串。
func (r *Response) RPCCode() string {
	if r.IsOK() {
		return ""
	}
	var eb errorBody
	if err := json.Unmarshal(r.Body, &eb); err != nil {
		return ""
	}
	return eb.Error.Status
}

// RPCMessage 返回错误体中的 message；解码失败时返回原始响应体，便于失败诊断。
func (r *Response) RPCMessage() string {
	if r.IsOK() {
		return ""
	}
	var eb errorBody
	if err := json.Unmarshal(r.Body, &eb); err != nil {
		return string(r.Body)
	}
	return eb.Error.Message
}

// HasRPCCode 按 gRPC 状态码断言比较（大小写与分隔符不敏感，
// 兼容 "NotFound" 与 "NOT_FOUND" 两种表示）。
func (r *Response) HasRPCCode(code string) bool {
	return normalizeCode(r.RPCCode()) == normalizeCode(code)
}

func normalizeCode(code string) string {
	replacer := strings.NewReplacer("_", "", "-", "", " ", "")
	return strings.ToLower(replacer.Replace(code))
}

// Do 发送 JSON 请求并返回原始响应。body 为 nil 时不携带请求体。
func (c *Client) Do(ctx context.Context, method, path string, body any) (*Response, error) {
	var reader io.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request body: %w", err)
		}
		reader = bytes.NewReader(raw)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reader)
	if err != nil {
		return nil, fmt.Errorf("build request %s %s: %w", method, path, err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if c.authorization != "" {
		req.Header.Set("Authorization", c.authorization)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do %s %s: %w", method, path, err)
	}
	defer func() { _ = resp.Body.Close() }()

	raw, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, fmt.Errorf("read response of %s %s: %w", method, path, err)
	}
	return &Response{StatusCode: resp.StatusCode, Header: resp.Header, Body: raw}, nil
}
