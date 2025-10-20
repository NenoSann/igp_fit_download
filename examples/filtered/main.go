package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"fit-viewer/igpsport"

	"github.com/joho/godotenv"
)

func main() {
	// 加载 .env 文件
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	authToken := os.Getenv("AUTHORIZATION")
	if authToken == "" {
		log.Fatal("AUTHORIZATION token not found in .env file")
	}

	// 创建客户端
	client := igpsport.NewClient(igpsport.Config{
		AuthToken: authToken,
	})

	fmt.Println("Fetching activities from 2024...")

	// 使用过滤条件获取 2024 年的活动
	activities, err := client.GetActivitiesWithFilter(igpsport.ActivityFilter{
		BeginTime:    "2024-01-01",
		EndTime:      "2024-12-31",
		PageSize:     50,
		RequestDelay: 300 * time.Millisecond,
	})

	if err != nil {
		log.Fatal("Failed to fetch activities:", err)
	}

	fmt.Printf("Found %d activities in 2024\n\n", len(activities))

	// 显示活动信息
	for i, activity := range activities {
		fmt.Printf("%d. %s - %s (%.2f km)\n",
			i+1,
			activity.StartTime,
			activity.Title,
			activity.RideDistance/1000)
	}

	// 询问是否下载
	fmt.Print("\nDo you want to download all FIT files? (y/n): ")
	var answer string
	fmt.Scanln(&answer)

	if answer != "y" && answer != "Y" {
		fmt.Println("Download cancelled.")
		return
	}

	// 创建下载器
	downloader, err := igpsport.NewDownloader(client, igpsport.DownloadOptions{
		DownloadDir: "fit_files_2024",
	})
	if err != nil {
		log.Fatal("Failed to create downloader:", err)
	}

	// 下载文件
	successCount := 0
	failCount := 0

	err = downloader.DownloadAll(activities, igpsport.DownloadOptions{
		RequestDelay: 500 * time.Millisecond,
		OnProgress: func(current, total int, activity igpsport.Activity) {
			fmt.Printf("[%d/%d] Downloading: %s\n", current, total, activity.Title)
		},
		OnError: func(activity igpsport.Activity, err error) {
			fmt.Printf("❌ Error: %s - %v\n", activity.Title, err)
			failCount++
		},
	})

	if err != nil {
		log.Fatal("Download failed:", err)
	}

	successCount = len(activities) - failCount

	fmt.Printf("\n✅ Download completed!\n")
	fmt.Printf("Success: %d, Failed: %d\n", successCount, failCount)
	fmt.Printf("Files saved to: %s\n", downloader.GetDownloadDir())
}
