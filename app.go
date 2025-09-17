package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// 添加读取配置文件的方法
func readConfig() (string, string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", "", err
	}

	configPath := filepath.Join(homeDir, "Library", "Application Support", "Offshoot Plus", "config.txt")

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", "", nil
		}
		return "", "", err
	}

	lines := strings.Split(string(data), "\n")
	if len(lines) < 2 {
		return "", "", nil
	}

	return lines[0], lines[1], nil
}

func (a *App) SaveSettings(projectTitle, backups string) error {
	if projectTitle == "" || backups == "" {
		return fmt.Errorf("项目名称和备份数量不能为空")
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configDir := filepath.Join(homeDir, "Library", "Application Support", "Offshoot Plus")
	configPath := filepath.Join(configDir, "config.txt")

	// 创建目录
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	// 写入配置文件
	return os.WriteFile(configPath, []byte(fmt.Sprintf("%s\n%s", projectTitle, backups)), 0644)
}

type App struct {
	ctx context.Context
}

type ReportRequest struct {
	ProjectTitle string   `json:"projectTitle"`
	Backups      string   `json:"backups"`
	FilePaths    []string `json:"filePaths"`
}

type FileInfo struct {
	Name      string `json:"name"`
	FileType  string `json:"fileType"`
	FileSize  string `json:"fileSize"`
	HashValue string `json:"hashValue"`
	Status    string `json:"status"`
}

type ClipsOverview struct {
	VideoFiles int    `json:"videoFiles"`
	AudioFiles int    `json:"audioFiles"`
	OtherFiles int    `json:"otherFiles"`
	TotalFiles int    `json:"totalFiles"`
	VideoSize  string `json:"videoSize"`
	AudioSize  string `json:"audioSize"`
	OtherSize  string `json:"otherSize"`
	TotalSize  string `json:"totalSize"`
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) domReady(ctx context.Context) {
}

func (a *App) shutdown(ctx context.Context) {
}

func (a *App) SelectFiles() ([]string, error) {
	files, err := runtime.OpenMultipleFilesDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select MHL Files",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "MHL Files (*.mhl)",
				Pattern:     "*.mhl",
			},
			{
				DisplayName: "Text Files (*.txt)",
				Pattern:     "*.txt",
			},
			{
				DisplayName: "All Files (*.*)",
				Pattern:     "*.*",
			},
		},
	})
	return files, err
}

// 修改GenerateReport方法
func (a *App) GenerateReport(request ReportRequest) error {
	projectTitle := request.ProjectTitle
	backups := request.Backups

	// 尝试从配置文件读取
	cfgProj, cfgBkup, err := readConfig()
	if err != nil {
		return err
	}

	// 判断使用哪个来源的数据
	useConfig := true
	if cfgProj == "" || cfgBkup == "" {
		useConfig = false
		if projectTitle == "" || backups == "" {
			return fmt.Errorf("请填写完整信息")
		}
	} else {
		projectTitle = cfgProj
		backups = cfgBkup
	}

	// 如果GUI数据有效且配置文件数据不完整，则保存到配置文件
	if !useConfig && projectTitle != "" && backups != "" {
		err := a.SaveSettings(projectTitle, backups)
		if err != nil {
			return err
		}
	}

	// 修改输出路径
	for _, filePath := range request.FilePaths {
		// 创建输出目录
		homeDir, _ := os.UserHomeDir()
		outputDir := filepath.Join(homeDir, "Documents", "Offshoot Reports")
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return err
		}

		// 修改outputPath生成逻辑
		logData, err := parseLogFile(filePath)
		if err != nil {
			return err
		}
		outputFileName := fmt.Sprintf("%s_DMT Report.pdf", logData.ReelName)
		outputPath := filepath.Join(outputDir, outputFileName)

		err = generatePDF(outputPath, logData, ReportRequest{
			ProjectTitle: projectTitle,
			Backups:      backups,
			FilePaths:    []string{filePath},
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) processFile(filePath string, request ReportRequest) error {
	// Parse the log file
	logData, err := parseLogFile(filePath)
	if err != nil {
		return err
	}

	// Get output directory (same as input file)
	outputDir := filepath.Dir(filePath)
	outputFileName := fmt.Sprintf("%s_DMT Report.pdf", logData.ReelName)
	outputPath := filepath.Join(outputDir, outputFileName)

	// --- 在这里添加调试信息 ---
	//fmt.Printf("--- DEBUG: About to call generatePDF ---\n")
	//fmt.Printf("Passing ProjectTitle to generatePDF: '%s'\n", request.ProjectTitle)
	//fmt.Printf("---------------------------------------\n")
	// --- 调试信息结束 ---

	// Generate PDF
	err = generatePDF(outputPath, logData, request)
	if err != nil {
		return err
	}

	return nil
}

func getSuffix(filename string) string {
	ext := filepath.Ext(filename)
	if len(ext) > 1 {
		return strings.ToUpper(ext[1:])
	}
	return ""
}

func humanReadableSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
