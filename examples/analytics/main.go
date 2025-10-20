package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

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

	fmt.Println("Fetching activities...")

	// 获取所有活动
	activities, err := client.GetAllActivities()
	if err != nil {
		log.Fatal("Failed to fetch activities:", err)
	}

	fmt.Printf("Found %d activities\n", len(activities))

	// 分析统计数据
	stats := analyzeActivities(activities)

	// 打印统计信息
	fmt.Println("\n=== Activity Statistics ===")
	fmt.Printf("Total activities: %d\n", stats.TotalCount)
	fmt.Printf("Total distance: %.2f km\n", stats.TotalDistance/1000)
	fmt.Printf("Total time: %.2f hours\n", stats.TotalTime/3600)
	fmt.Printf("Average distance: %.2f km\n", stats.AvgDistance/1000)
	fmt.Printf("Average speed: %.2f km/h\n", stats.AvgSpeed)
	fmt.Printf("Longest ride: %.2f km\n", stats.LongestRide/1000)
	fmt.Printf("Total ascent: %d m\n", stats.TotalAscent)

	// 保存到 JSON 文件
	saveToJSON(activities, "activities_analysis.json")
	fmt.Println("\n✅ Full data saved to activities_analysis.json")
}

type ActivityStats struct {
	TotalCount    int
	TotalDistance float64
	TotalTime     float64
	TotalAscent   int
	AvgDistance   float64
	AvgSpeed      float64
	LongestRide   float64
}

func analyzeActivities(activities []igpsport.Activity) ActivityStats {
	stats := ActivityStats{
		TotalCount: len(activities),
	}

	for _, activity := range activities {
		stats.TotalDistance += activity.RideDistance
		stats.TotalTime += activity.TotalMovingTime
		stats.TotalAscent += activity.TotalAscent

		if activity.RideDistance > stats.LongestRide {
			stats.LongestRide = activity.RideDistance
		}
	}

	if stats.TotalCount > 0 {
		stats.AvgDistance = stats.TotalDistance / float64(stats.TotalCount)
		if stats.TotalTime > 0 {
			stats.AvgSpeed = (stats.TotalDistance / stats.TotalTime) * 3.6 // m/s to km/h
		}
	}

	return stats
}

func saveToJSON(activities []igpsport.Activity, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		log.Printf("Failed to create file: %v", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(activities); err != nil {
		log.Printf("Failed to write JSON: %v", err)
	}
}
