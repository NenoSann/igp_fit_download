#!/bin/bash

# iGPSport FIT 文件下载器启动脚本

echo "=================================="
echo "iGPSport FIT 文件下载器"
echo "=================================="
echo ""

# 检查 .env 文件是否存在
if [ ! -f .env ]; then
    echo "❌ 错误: .env 文件不存在"
    echo ""
    echo "请按照以下步骤配置："
    echo "1. 复制示例配置: cp .env.example .env"
    echo "2. 从浏览器获取 Authorization token"
    echo "3. 编辑 .env 文件，填入你的 token"
    echo ""
    echo "详细步骤请查看 QUICKSTART.md"
    exit 1
fi

# 检查 AUTHORIZATION 是否配置
if ! grep -q "AUTHORIZATION=." .env; then
    echo "❌ 错误: .env 文件中未配置 AUTHORIZATION"
    echo ""
    echo "请编辑 .env 文件，填入你的 Authorization token"
    echo "详细步骤请查看 QUICKSTART.md"
    exit 1
fi

# 检查 Go 是否安装
if ! command -v go &> /dev/null; then
    echo "❌ 错误: Go 未安装"
    echo "请访问 https://golang.org/dl/ 下载并安装 Go"
    exit 1
fi

echo "✅ 配置检查通过"
echo ""
echo "开始下载 FIT 文件..."
echo ""

# 运行爬虫程序
go run iGPSport_Crawler.go

# 检查运行结果
if [ $? -eq 0 ]; then
    echo ""
    echo "=================================="
    echo "✅ 下载完成！"
    echo "=================================="
    echo ""
    echo "FIT 文件已保存到: downloaded_fit_files/"
    echo ""
    echo "查看下载的文件:"
    echo "  ls -l downloaded_fit_files/"
else
    echo ""
    echo "=================================="
    echo "❌ 下载过程中出现错误"
    echo "=================================="
    echo ""
    echo "常见问题:"
    echo "1. Token 过期 - 需要重新获取"
    echo "2. 网络问题 - 检查网络连接"
    echo "3. 查看上方错误信息了解详情"
fi
