# 示例程序

这个目录包含了使用 `igpsport` 包的示例程序。

## 示例列表

### 1. basic - 基础使用

最简单的示例，演示如何获取所有活动并下载 FIT 文件。

```bash
cd examples/basic
go run main.go
```

**功能：**
- 获取所有活动记录
- 批量下载所有 FIT 文件
- 显示下载进度

### 2. filtered - 条件过滤

演示如何使用过滤条件获取特定时间段的活动。

```bash
cd examples/filtered
go run main.go
```

**功能：**
- 获取 2024 年的活动记录
- 显示活动列表
- 交互式确认是否下载
- 统计成功/失败数量

### 3. analytics - 数据分析

演示如何分析活动数据并保存为 JSON。

```bash
cd examples/analytics
go run main.go
```

**功能：**
- 获取所有活动
- 计算统计信息（总距离、平均速度等）
- 保存数据到 JSON 文件

## 运行前准备

1. **配置环境变量**

在项目根目录创建 `.env` 文件：

```bash
cd /Users/liangheng/Code/Web/fit-viewer
cp .env.example .env
# 编辑 .env 文件，填入你的 Authorization token
```

2. **安装依赖**

```bash
cd /Users/liangheng/Code/Web/fit-viewer
go mod tidy
```

## 自定义示例

你可以基于这些示例创建自己的程序：

```go
package main

import (
    "fmt"
    "fit-viewer/igpsport"
)

func main() {
    client := igpsport.NewClient(igpsport.Config{
        AuthToken: "your_token",
    })
    
    // 你的代码...
}
```

## 常见用例

### 只下载最近 10 个活动

```go
activities, _ := client.GetAllActivities()
recent := activities[:10]  // 取前 10 个
downloader.DownloadAll(recent, igpsport.DownloadOptions{})
```

### 只下载特定类型的活动

```go
activities, _ := client.GetAllActivities()
var filtered []igpsport.Activity
for _, a := range activities {
    if a.ExerciseType == 0 {  // 只要骑行活动
        filtered = append(filtered, a)
    }
}
downloader.DownloadAll(filtered, igpsport.DownloadOptions{})
```

### 下载大于 50km 的活动

```go
activities, _ := client.GetAllActivities()
var longRides []igpsport.Activity
for _, a := range activities {
    if a.RideDistance > 50000 {  // 大于 50km (单位是米)
        longRides = append(longRides, a)
    }
}
downloader.DownloadAll(longRides, igpsport.DownloadOptions{})
```

## 在其他项目中使用

如果你想在其他 Go 项目中使用这个包：

```bash
# 在你的项目目录
go mod init your-project
go mod edit -replace fit-viewer=/Users/liangheng/Code/Web/fit-viewer
go mod tidy
```

然后在代码中：

```go
import "fit-viewer/igpsport"
```
