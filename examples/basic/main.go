package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/NenoSann/igp_fit_download/igpsport"

	"github.com/joho/godotenv"
)

func main() {
	// 加载 .env 文件
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// 获取 Authorization token
	authToken := os.Getenv("AUTHORIZATION")
	if authToken == "" {
		log.Fatal("AUTHORIZATION token not found in .env file")
	}

	// 创建客户端
	client := igpsport.NewClient(igpsport.Config{
		AuthToken: authToken,
		Timeout:   30 * time.Second,
	})

	fmt.Println("Fetching activities...")

	// 获取所有活动
	activities, err := client.GetAllActivities()
	if err != nil {
		log.Fatal("Failed to fetch activities:", err)
	}

	fmt.Printf("Found %d activities\n\n", len(activities))

	// 创建下载器
	downloader, err := igpsport.NewDownloader(client, igpsport.DownloadOptions{
		DownloadDir: "downloaded_fit_files",
	})
	if err != nil {
		log.Fatal("Failed to create downloader:", err)
	}

	fmt.Println("Starting download...")

	// 下载所有文件，带进度显示
	err = downloader.DownloadAll(activities, igpsport.DownloadOptions{
		RequestDelay: 500 * time.Millisecond,
		OnProgress: func(current, total int, activity igpsport.Activity) {
			fmt.Printf("[%d/%d] Downloading: %s (RideID: %d)\n",
				current, total, activity.Title, activity.RideID)
		},
		OnError: func(activity igpsport.Activity, err error) {
			fmt.Printf("❌ Failed to download %s: %v\n", activity.Title, err)
		},
	})

	if err != nil {
		log.Fatal("Download failed:", err)
	}

	fmt.Println("\n✅ Download completed!")
	fmt.Printf("Files saved to: %s\n", downloader.GetDownloadDir())
}
