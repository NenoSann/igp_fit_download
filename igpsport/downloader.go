package igpsport

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Downloader FIT 文件下载器
type Downloader struct {
	client      *Client
	downloadDir string
}

// DownloadOptions 下载选项
type DownloadOptions struct {
	DownloadDir  string                                      // 下载目录，默认为 "downloaded_fit_files"
	RequestDelay time.Duration                               // 请求间延迟，默认 500ms
	OnProgress   func(current, total int, activity Activity) // 进度回调
	OnError      func(activity Activity, err error)          // 错误回调
}

// NewDownloader 创建新的下载器
func NewDownloader(client *Client, options DownloadOptions) (*Downloader, error) {
	if options.DownloadDir == "" {
		options.DownloadDir = "downloaded_fit_files"
	}

	// 创建下载目录
	if err := os.MkdirAll(options.DownloadDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create download directory: %w", err)
	}

	return &Downloader{
		client:      client,
		downloadDir: options.DownloadDir,
	}, nil
}

// DownloadAll 下载所有活动的 FIT 文件
func (d *Downloader) DownloadAll(activities []Activity, options DownloadOptions) error {
	if options.RequestDelay == 0 {
		options.RequestDelay = 500 * time.Millisecond
	}

	for i, activity := range activities {
		// 进度回调
		if options.OnProgress != nil {
			options.OnProgress(i+1, len(activities), activity)
		}

		err := d.Download(activity)
		if err != nil {
			// 错误回调
			if options.OnError != nil {
				options.OnError(activity, err)
			} else {
				// 默认错误处理：记录但继续
				fmt.Printf("Failed to download activity %d: %v\n", activity.RideID, err)
			}
			continue
		}

		// 添加延迟
		if i < len(activities)-1 {
			time.Sleep(options.RequestDelay)
		}
	}

	return nil
}

// Download 下载单个活动的 FIT 文件
func (d *Downloader) Download(activity Activity) error {
	// 获取下载链接
	downloadURL, err := d.client.GetDownloadURL(activity.RideID)
	if err != nil {
		return fmt.Errorf("failed to get download URL: %w", err)
	}

	// 下载文件内容
	data, err := d.client.DownloadFitFile(downloadURL)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}

	// 生成文件名
	filename := d.generateFilename(activity)
	filePath := filepath.Join(d.downloadDir, filename)

	// 保存文件
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	return nil
}

// generateFilename 生成文件名
func (d *Downloader) generateFilename(activity Activity) string {
	return fmt.Sprintf("%s-%d-%s.fit",
		activity.Title,
		activity.RideID,
		activity.StartTime)
}

// SetDownloadDir 设置下载目录
func (d *Downloader) SetDownloadDir(dir string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create download directory: %w", err)
	}
	d.downloadDir = dir
	return nil
}

// GetDownloadDir 获取下载目录
func (d *Downloader) GetDownloadDir() string {
	return d.downloadDir
}
