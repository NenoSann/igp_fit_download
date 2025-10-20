package main

import (
	"encoding/json"
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

	// 获取所有活动记录
	activities, err := client.GetAllActivities()
	if err != nil {
		log.Fatal("Failed to fetch activities:", err)
	}
	log.Printf("Found %d activities\n", len(activities))

	// 持久化所有活动记录到本地 JSON 文件
	err = persistAllActivities(activities)
	if err != nil {
		log.Fatal("Failed to persist activities:", err)
	}

	// 创建下载器
	downloader, err := igpsport.NewDownloader(client, igpsport.DownloadOptions{
		DownloadDir: "downloaded_fit_files",
	})
	if err != nil {
		log.Fatal("Failed to create downloader:", err)
	}

	fmt.Println("Starting download...")

	// 下载所有 FIT 文件
	err = downloader.DownloadAll(activities, igpsport.DownloadOptions{
		RequestDelay: 500 * time.Millisecond,
		OnProgress: func(current, total int, activity igpsport.Activity) {
			log.Printf("[%d/%d] Processing activity: %s (RideID: %d)\n",
				current, total, activity.Title, activity.RideID)
		},
		OnError: func(activity igpsport.Activity, err error) {
			log.Printf("Failed to download activity %d: %v\n", activity.RideID, err)
		},
	})

	if err != nil {
		log.Fatal("Download failed:", err)
	}

	log.Println("Download completed!")
}

// persistAllActivities 持久化所有活动记录到本地 JSON 文件
func persistAllActivities(activities []igpsport.Activity) error {
	filePath := time.Now().Format("2006-01-02 15:04:05") + "-activities.json"
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(activities); err != nil {
		return fmt.Errorf("failed to write JSON to file: %w", err)
	}

	log.Printf("Activities saved to %s\n", filePath)
	return nil
}
