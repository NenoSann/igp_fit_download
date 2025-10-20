package igpsport

import "time"

// ActivityResponse 活动列表响应结构
type ActivityResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		PageNo    int        `json:"pageNo"`
		PageSize  int        `json:"pageSize"`
		TotalPage int        `json:"totalPage"`
		TotalRows int        `json:"totalRows"`
		Rows      []Activity `json:"rows"`
	} `json:"data"`
}

// Activity 单个活动记录
type Activity struct {
	ID              string  `json:"id"`
	RideID          int     `json:"rideId"`
	ExerciseType    int     `json:"exerciseType"`
	Title           string  `json:"title"`
	StartTime       string  `json:"startTime"`
	RideDistance    float64 `json:"rideDistance"`
	TotalMovingTime float64 `json:"totalMovingTime"`
	AvgSpeed        float64 `json:"avgSpeed"`
	DataStatus      int     `json:"dataStatus"`
	ErrorType       int     `json:"errorType"`
	AnalysisStatus  int     `json:"analysisStatus"`
	Label           int     `json:"label"`
	IsOpen          int     `json:"isOpen"`
	UnRead          bool    `json:"unRead"`
	Icon            string  `json:"icon"`
	TotalAscent     int     `json:"totalAscent"`
}

// DownloadResponse 下载链接响应结构
type DownloadResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data string `json:"data"` // 下载链接
}

// ActivityFilter 活动查询过滤条件
type ActivityFilter struct {
	BeginTime    string        // 开始时间，格式：2023-10-01
	EndTime      string        // 结束时间，格式：2025-10-20
	ReqType      int           // 请求类型，默认 0
	Sort         int           // 排序方式，默认 1
	PageSize     int           // 每页数量，默认 20
	RequestDelay time.Duration // 请求间延迟，默认 500ms
}
