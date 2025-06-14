package main

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
)

// 版本資訊變數 (編譯時注入)
var (
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
	GitBranch = "unknown"
)

// 版本資訊結構
type VersionInfo struct {
	Version   string `json:"version"`
	BuildTime string `json:"build_time"`
	GitCommit string `json:"git_commit"`
	GitBranch string `json:"git_branch"`
	GoVersion string `json:"go_version"`
	OS        string `json:"os"`
	Arch      string `json:"arch"`
}

// 取得版本資訊
func getVersionInfo() VersionInfo {
	return VersionInfo{
		Version:   Version,
		BuildTime: BuildTime,
		GitCommit: GitCommit,
		GitBranch: GitBranch,
		GoVersion: runtime.Version(),
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
	}
}

// 版本資訊 HTTP 處理器
func versionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(getVersionInfo())
}

// 顯示版本資訊
func showVersion() {
	info := getVersionInfo()
	fmt.Printf("🚀 %s 版本資訊\n", os.Args[0])
	fmt.Printf("版本: %s\n", info.Version)
	fmt.Printf("建置時間: %s\n", info.BuildTime)
	fmt.Printf("Git 提交: %s\n", info.GitCommit)
	fmt.Printf("Git 分支: %s\n", info.GitBranch)
	fmt.Printf("Go 版本: %s\n", info.GoVersion)
	fmt.Printf("作業系統: %s\n", info.OS)
	fmt.Printf("架構: %s\n", info.Arch)
}

func main() {
	logger := &lumberjack.Logger{
		Filename:   "./logs/service.log", // log 檔案路徑
		MaxSize:    10,                   // 每個 log 檔最大 MB
		MaxBackups: 200,                  // 最多保留幾個舊 log
		MaxAge:     60,                   // 最多保留幾天
		Compress:   false,                // 是否壓縮舊 log
		LocalTime:  true,                 // 使用本地時間分割
	}
	multiWriter := io.MultiWriter(os.Stdout, logger)
	log.SetOutput(multiWriter)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Println("service started")

	myApp := app.NewWithID("com.yee2.qc")
	myWindow := myApp.NewWindow("進階 Fyne 範例")
	myWindow.Resize(fyne.NewSize(600, 400))

	myWindow.SetMainMenu(CreateMen(myWindow))

	myWindow.SetContent(CreateMainTab())
	myWindow.ShowAndRun()

	// 檢查命令列參數
	if len(os.Args) > 1 && os.Args[1] == "version" {
		showVersion()
		return
	}

	log.Printf("🚀 啟動 %s v%s\n", os.Args[0], Version)
	log.Printf("建置時間: %s\n", BuildTime)

	// 設定路由
	http.HandleFunc("/version", versionHandler)

	// 啟動服務
	port := "8080"
	log.Printf("🌐 HTTP 服務啟動於端口 %s\n", port)
	http.ListenAndServe(":"+port, nil)
}
