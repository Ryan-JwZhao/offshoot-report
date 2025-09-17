package main

import (
	"embed"
	"flag"
	"fmt"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// 添加命令行参数解析
	projectName := flag.String("projectName", "", "Project name")
	backupCount := flag.String("backupCount", "", "Backup count")
	mhlPath := flag.String("mhlPath", "", "MHL file path")
	flag.Parse()

	// 如果所有命令行参数都有值，则直接执行生成报告
	if *projectName != "" && *backupCount != "" && *mhlPath != "" {
		app := NewApp()
		request := ReportRequest{
			ProjectTitle: *projectName,
			Backups:      *backupCount,
			FilePaths:    []string{*mhlPath},
		}
		err := app.processFile(*mhlPath, request)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		return
	}

	app := NewApp()

	err := wails.Run(&options.App{
		Title:         "Offshoot Plus",
		Width:         450,
		Height:        400,
		DisableResize: true,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup:  app.startup,
		OnDomReady: app.domReady,
		OnShutdown: app.shutdown,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
