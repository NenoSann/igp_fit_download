# iGPSport FIT 文件爬虫

这是一个用 Go 语言编写的工具，用于从 iGPSport 网站获取运动记录并下载 FIT 文件。

## ✨ 新版本特性

**现在这个项目已经重构为一个可重用的 Go 包！** 你可以：

- 🔧 在当前项目中使用（作为爬虫工具）
- 📦 在其他 Go 项目中导入使用（作为库）
- 🎯 使用提供的示例程序快速开始

## 📁 项目结构

```
fit-viewer/
├── igpsport/              # 核心包（可在其他项目中使用）
│   ├── client.go         # API 客户端
│   ├── downloader.go     # 文件下载器
│   └── types.go          # 数据类型
├── examples/             # 示例程序
│   ├── basic/           # 基础使用
│   ├── filtered/        # 条件过滤
│   └── analytics/       # 数据分析
└── iGPSport_Crawler.go  # 主程序
```

## 功能特性

- ✅ 自动获取所有活动记录（支持分页）
- ✅ 批量下载 FIT 文件
- ✅ 条件过滤（时间范围、活动类型等）
- ✅ 自定义下载目录
- ✅ 进度回调和错误处理
- ✅ 请求间自动延迟，避免请求过快
- ✅ 可作为库在其他项目中使用

## 安装步骤

### 方式 1: 作为爬虫工具使用（推荐新手）

这是最简单的方式，适合只想下载 FIT 文件的用户。

#### 1. 安装依赖

```bash
cd /Users/liangheng/Code/Web/fit-viewer
go mod tidy
```

#### 2. 配置 Authorization Token

1. 复制示例配置文件：
```bash
cp .env.example .env
```

2. 获取你的 Authorization token：
   - 打开浏览器，访问 https://app.zh.igpsport.com/
   - 登录你的账号
   - 按 F12 打开开发者工具
   - 切换到 "Network" (网络) 标签
   - 刷新页面或浏览你的活动记录
   - 找到任意一个请求到 `prod.zh.igpsport.com` 的请求
   - 在请求头中找到 `Authorization` 字段
   - 复制完整的 Authorization 值（包括 Bearer 等前缀，如果有的话）

3. 编辑 `.env` 文件，将 `your_authorization_token_here` 替换为你复制的 token：
```
AUTHORIZATION=你的token值
```

#### 3. 运行程序

```bash
# 运行主程序
go run iGPSport_Crawler.go

# 或运行编译后的程序
go build -o igpsport_crawler iGPSport_Crawler.go
./igpsport_crawler
```

### 方式 2: 作为 Go 包在其他项目中使用

如果你想在自己的 Go 项目中使用这个功能，查看 **[包使用指南](PACKAGE_USAGE.md)**。

快速示例：

```go
package main

import (
    "fit-viewer/igpsport"
    "fmt"
)

func main() {
    client := igpsport.NewClient(igpsport.Config{
        AuthToken: "your_token",
    })
    
    activities, _ := client.GetAllActivities()
    fmt.Printf("Found %d activities\n", len(activities))
}
```

### 方式 3: 使用示例程序

查看 `examples/` 目录中的示例：

```bash
cd examples/basic
go run main.go
```

## 运行输出

程序会：
1. 读取 `.env` 文件中的 token
2. 获取所有活动记录
3. 逐个下载 FIT 文件到 `downloaded_fit_files` 目录

## 输出示例

```
2025/10/20 15:30:00 Fetched page 1/13 (20 activities)
2025/10/20 15:30:01 Fetched page 2/13 (20 activities)
...
2025/10/20 15:30:10 Found 247 activities
2025/10/20 15:30:10 [1/247] Processing activity: 户外骑行 (RideID: 40979454)
2025/10/20 15:30:11 Downloaded: ride-40979454-2025.10.18.fit
2025/10/20 15:30:12 [2/247] Processing activity: 户外骑行 (RideID: 40978123)
2025/10/20 15:30:13 Downloaded: ride-40978123-2025.10.17.fit
...
2025/10/20 15:35:00 Download completed!
```

## 文件结构

运行完成后，目录结构如下：

```
fit-viewer/
├── igpsport/                     # 核心包
│   ├── client.go
│   ├── downloader.go
│   └── types.go
├── examples/                     # 示例程序
│   ├── basic/
│   ├── filtered/
│   └── analytics/
├── iGPSport_Crawler.go          # 主程序
├── .env                         # 配置文件（包含 token）
├── .env.example                 # 配置示例
├── go.mod                       # Go 模块文件
├── go.sum                       # Go 依赖校验文件
├── downloaded_fit_files/        # 下载的 FIT 文件目录
│   ├── ride-40979454-2025.10.18.fit
│   └── ...
└── 2025-10-20 15:04:05-activities.json  # 活动记录 JSON
```

## 📚 文档

- **[PACKAGE_USAGE.md](PACKAGE_USAGE.md)** - 如何在其他项目中使用这个包
- **[REFACTORING_SUMMARY.md](REFACTORING_SUMMARY.md)** - 重构总结和新功能说明
- **[QUICKSTART.md](QUICKSTART.md)** - 快速开始指南
- **[igpsport/README.md](igpsport/README.md)** - 包 API 文档
- **[examples/README.md](examples/README.md)** - 示例程序说明

## 注意事项

1. **保护你的 token**：`.env` 文件包含敏感信息，不要上传到 git 仓库
2. **请求频率**：程序已经添加了 500ms 的延迟，避免请求过快被封禁
3. **token 过期**：如果遇到认证错误，可能是 token 过期了，需要重新获取
4. **网络问题**：如果下载失败，程序会继续下载下一个文件

## 故障排除

### 错误：AUTHORIZATION token not found in .env file
- 确保你已经创建了 `.env` 文件
- 确保文件中有 `AUTHORIZATION=...` 这一行

### 错误：API error: Unauthorized
- 你的 token 可能已经过期，需要重新从浏览器获取

### 错误：Failed to download fit file
- 可能是网络问题，稍后重试
- 也可能是某些活动没有对应的 FIT 文件

## 自定义选项

### 使用主程序

你可以修改 `iGPSport_Crawler.go` 中的参数。

### 使用包 API

```go
// 按时间范围过滤
activities, _ := client.GetActivitiesWithFilter(igpsport.ActivityFilter{
    BeginTime: "2024-01-01",
    EndTime:   "2024-12-31",
})

// 自定义下载选项
downloader.DownloadAll(activities, igpsport.DownloadOptions{
    RequestDelay: 1 * time.Second,
    OnProgress: func(current, total int, activity igpsport.Activity) {
        // 自定义进度显示
    },
})
```

查看 **[PACKAGE_USAGE.md](PACKAGE_USAGE.md)** 了解更多用法。

## 🎓 学习和使用

| 你想做什么 | 查看文档 |
|-----------|---------|
| 快速下载 FIT 文件 | 本文档 + `QUICKSTART.md` |
| 在其他项目中使用 | `PACKAGE_USAGE.md` |
| 查看 API 文档 | `igpsport/README.md` |
| 学习使用示例 | `examples/README.md` |
| 了解重构内容 | `REFACTORING_SUMMARY.md` |

## License

MIT
