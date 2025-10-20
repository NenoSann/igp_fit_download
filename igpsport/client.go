package igpsport

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/dsnet/compress/brotli"
)

// Client iGPSport API 客户端
type Client struct {
	authToken  string
	httpClient *http.Client
	baseURL    string
}

// Config 客户端配置
type Config struct {
	AuthToken string        // Authorization token（必需）
	Timeout   time.Duration // HTTP 请求超时时间（可选，默认 30 秒）
}

// NewClient 创建新的 iGPSport 客户端
func NewClient(config Config) *Client {
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	return &Client{
		authToken: config.AuthToken,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
		baseURL: "https://prod.zh.igpsport.com/service/web-gateway/web-analyze/activity",
	}
}

// GetAllActivities 获取所有活动记录
func (c *Client) GetAllActivities() ([]Activity, error) {
	return c.GetActivitiesWithFilter(ActivityFilter{})
}

// GetActivitiesWithFilter 根据过滤条件获取活动记录
func (c *Client) GetActivitiesWithFilter(filter ActivityFilter) ([]Activity, error) {
	var allActivities []Activity
	pageNo := 1
	pageSize := filter.PageSize
	if pageSize == 0 {
		pageSize = 20
	}

	for {
		url := c.buildListURL(pageNo, pageSize, filter)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		c.setHeaders(req)
		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to send request: %w", err)
		}

		// 处理响应
		activityResp, err := c.parseActivityResponse(resp)
		resp.Body.Close()
		if err != nil {
			return nil, err
		}

		if activityResp.Code != 0 {
			return nil, fmt.Errorf("API error: %s", activityResp.Msg)
		}

		allActivities = append(allActivities, activityResp.Data.Rows...)

		// 检查是否还有下一页
		if pageNo >= activityResp.Data.TotalPage {
			break
		}

		pageNo++

		// 添加延迟，避免请求过快
		if filter.RequestDelay > 0 {
			time.Sleep(filter.RequestDelay)
		} else {
			time.Sleep(500 * time.Millisecond)
		}
	}

	return allActivities, nil
}

// GetDownloadURL 获取指定活动的 FIT 文件下载链接
func (c *Client) GetDownloadURL(rideID int) (string, error) {
	url := fmt.Sprintf("%s/getDownloadUrl/%d", c.baseURL, rideID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	c.setHeaders(req)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var downloadResp DownloadResponse
	if err := json.Unmarshal(body, &downloadResp); err != nil {
		return "", fmt.Errorf("failed to parse JSON: %w", err)
	}

	if downloadResp.Code != 0 {
		return "", fmt.Errorf("API error: %s", downloadResp.Msg)
	}

	return downloadResp.Data, nil
}

// DownloadFitFile 下载 FIT 文件内容
func (c *Client) DownloadFitFile(downloadURL string) ([]byte, error) {
	req, err := http.NewRequest("GET", downloadURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return data, nil
}

// buildListURL 构建活动列表请求 URL
func (c *Client) buildListURL(pageNo, pageSize int, filter ActivityFilter) string {
	url := fmt.Sprintf("%s/queryMyActivity?pageNo=%d&pageSize=%d&reqType=%d&sort=%d",
		c.baseURL, pageNo, pageSize, filter.ReqType, filter.Sort)

	if filter.BeginTime != "" {
		url += "&beginTime=" + filter.BeginTime
	}
	if filter.EndTime != "" {
		url += "&endTime=" + filter.EndTime
	}

	return url
}

// setHeaders 设置 HTTP 请求头
func (c *Client) setHeaders(req *http.Request) {
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Set("Accept-Language", "zh-Hans")
	req.Header.Set("Authorization", c.authToken)
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("DNT", "1")
	req.Header.Set("hasloading", "true")
	req.Header.Set("Origin", "https://app.zh.igpsport.com")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("qiwu-app-version", "1.0.0")
	req.Header.Set("Referer", "https://app.zh.igpsport.com/")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Timezone", "Asia/Shanghai")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0.0.0 Safari/537.36")
}

// parseActivityResponse 解析活动列表响应（处理压缩）
func (c *Client) parseActivityResponse(resp *http.Response) (*ActivityResponse, error) {
	var reader io.ReadCloser
	var err error

	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to create gzip reader: %w", err)
		}
		defer reader.Close()
	case "br":
		reader, err = brotli.NewReader(resp.Body, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create brotli reader: %w", err)
		}
		defer reader.Close()
	default:
		reader = resp.Body
	}

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var activityResp ActivityResponse
	if err := json.Unmarshal(data, &activityResp); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return &activityResp, nil
}
