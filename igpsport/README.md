# iGPSport Go Package

这是一个用于访问 iGPSport API 并下载 FIT 文件的 Go 包。

## 安装

```bash
go get github.com/NenoSann/igp_fit_download/igpsport
```

## 快速开始

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/NenoSann/igp_fit_download/igpsport"
)

func main() {
    // 创建客户端
    client := igpsport.NewClient(igpsport.Config{
        AuthToken: "your_authorization_token",
    })
    
    // 获取所有活动
    activities, err := client.GetAllActivities()
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Found %d activities\n", len(activities))
    
    // 创建下载器
    downloader, err := igpsport.NewDownloader(client, igpsport.DownloadOptions{
        DownloadDir: "fit_files",
    })
    if err != nil {
        log.Fatal(err)
    }
    
    // 下载所有文件
    err = downloader.DownloadAll(activities, igpsport.DownloadOptions{
        OnProgress: func(current, total int, activity igpsport.Activity) {
            fmt.Printf("[%d/%d] Downloading: %s\n", current, total, activity.Title)
        },
    })
    if err != nil {
        log.Fatal(err)
    }
}
```

## 主要功能

### 1. 创建客户端

```go
client := igpsport.NewClient(igpsport.Config{
    AuthToken: "your_token",
    Timeout:   30 * time.Second, // 可选，默认 30 秒
})
```

### 2. 获取活动列表

**获取所有活动：**

```go
activities, err := client.GetAllActivities()
```

**根据条件过滤：**

```go
activities, err := client.GetActivitiesWithFilter(igpsport.ActivityFilter{
    BeginTime: "2024-01-01",
    EndTime:   "2024-12-31",
    PageSize:  50,
})
```

### 3. 下载 FIT 文件

**批量下载：**

```go
downloader, _ := igpsport.NewDownloader(client, igpsport.DownloadOptions{
    DownloadDir: "fit_files",
})

downloader.DownloadAll(activities, igpsport.DownloadOptions{
    RequestDelay: 500 * time.Millisecond,
    OnProgress: func(current, total int, activity igpsport.Activity) {
        fmt.Printf("[%d/%d] %s\n", current, total, activity.Title)
    },
    OnError: func(activity igpsport.Activity, err error) {
        fmt.Printf("Error downloading %s: %v\n", activity.Title, err)
    },
})
```

**下载单个文件：**

```go
err := downloader.Download(activity)
```

**手动下载：**

```go
// 获取下载链接
downloadURL, err := client.GetDownloadURL(activity.RideID)

// 下载文件内容
data, err := client.DownloadFitFile(downloadURL)

// 保存到文件
os.WriteFile("ride.fit", data, 0644)
```

## API 文档

### Client

#### `NewClient(config Config) *Client`

创建新的 iGPSport 客户端。

**参数：**
- `config.AuthToken` (string, 必需): Authorization token
- `config.Timeout` (time.Duration, 可选): HTTP 超时时间，默认 30 秒

#### `GetAllActivities() ([]Activity, error)`

获取所有活动记录。

#### `GetActivitiesWithFilter(filter ActivityFilter) ([]Activity, error)`

根据过滤条件获取活动记录。

**过滤选项：**
- `BeginTime` (string): 开始时间，格式 "2023-10-01"
- `EndTime` (string): 结束时间，格式 "2025-10-20"
- `PageSize` (int): 每页数量，默认 20
- `RequestDelay` (time.Duration): 请求间延迟，默认 500ms

#### `GetDownloadURL(rideID int) (string, error)`

获取指定活动的 FIT 文件下载链接。

#### `DownloadFitFile(downloadURL string) ([]byte, error)`

下载 FIT 文件内容。

### Downloader

#### `NewDownloader(client *Client, options DownloadOptions) (*Downloader, error)`

创建新的下载器。

**选项：**
- `DownloadDir` (string): 下载目录，默认 "downloaded_fit_files"

#### `DownloadAll(activities []Activity, options DownloadOptions) error`

批量下载活动的 FIT 文件。

**选项：**
- `RequestDelay` (time.Duration): 请求间延迟
- `OnProgress` (func): 进度回调函数
- `OnError` (func): 错误回调函数

#### `Download(activity Activity) error`

下载单个活动的 FIT 文件。

## 数据结构

### Activity

```go
type Activity struct {
    ID              string
    RideID          int
    ExerciseType    int
    Title           string
    StartTime       string
    RideDistance    float64
    TotalMovingTime float64
    AvgSpeed        float64
    DataStatus      int
    ErrorType       int
    AnalysisStatus  int
    Label           int
    IsOpen          int
    UnRead          bool
    Icon            string
    TotalAscent     int
}
```

## 完整示例

查看 `examples/` 目录中的完整示例。

## License

MIT
