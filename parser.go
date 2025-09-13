package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type LogData struct {
	ReelName       string
	StartTime      string
	FinishTime     string
	ProjectTitle   string
	Files          []FileInfo
	ClipsOverview  ClipsOverview
	TotalSize      int64
	HashType       string
	GenerationTime string
}

func parseLogFile(filePath string) (*LogData, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	contentStr := string(content)
	logData := &LogData{
		HashType:       "xxhash64be",
		GenerationTime: time.Now().Format("2006/01/02 15:04"),
	}

	// Check if it's a Hedge log or MHL file
	if strings.Contains(contentStr, "Hedge") {
		return parseHedgeLog(contentStr, logData)
	} else {
		return parseMHLFile(contentStr, logData)
	}
}

func parseHedgeLog(content string, logData *LogData) (*LogData, error) {
	scanner := bufio.NewScanner(strings.NewReader(content))

	var currentFile *FileInfo

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "Source:") {
			// Extract reel name from source path
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				sourcePath := parts[1]
				logData.ReelName = filepath.Base(sourcePath)
			}
		} else if strings.HasPrefix(line, "Started:") {
			logData.StartTime = strings.TrimPrefix(line, "Started: ")
		} else if strings.HasPrefix(line, "Finished:") {
			logData.FinishTime = strings.TrimPrefix(line, "Finished: ")
		} else if strings.HasPrefix(line, "#") && strings.Contains(line, ":") {
			// New file entry
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				fileName := strings.TrimSpace(parts[1])
				currentFile = &FileInfo{
					Name:     strings.TrimSuffix(fileName, filepath.Ext(fileName)),
					FileType: getSuffix(fileName),
					Status:   "verified",
				}
				logData.Files = append(logData.Files, *currentFile)
			}
		} else if strings.HasPrefix(line, "Size:") && len(logData.Files) > 0 {
			sizeStr := strings.TrimPrefix(line, "Size: ")
			sizeStr = strings.TrimSuffix(sizeStr, " bytes")
			if size, err := strconv.ParseInt(sizeStr, 10, 64); err == nil {
				logData.TotalSize += size
				logData.Files[len(logData.Files)-1].FileSize = humanReadableSize(size)
			}
		} else if strings.HasPrefix(line, "Source hash:") && len(logData.Files) > 0 {
			hashValue := strings.TrimPrefix(line, "Source hash: ")
			logData.Files[len(logData.Files)-1].HashValue = fmt.Sprintf("%s: %s", logData.HashType, hashValue)
		}
	}

	logData.ClipsOverview = calculateOverview(logData.Files, logData.TotalSize)
	return logData, nil
}

func parseMHLFile(content string, logData *LogData) (*LogData, error) {
	// 定义局部正则
	reelNameRe := regexp.MustCompile(`<sourceInfoField name="Source Name">(.*?)<`)
	startTimeRe := regexp.MustCompile(`<startdate>(.*?)<`)
	finishTimeRe := regexp.MustCompile(`<finishdate>(.*?)<`)
	fileRe := regexp.MustCompile(`<file>(.*?)<`)
	sizeRe := regexp.MustCompile(`<size>(.*?)<`)
	hashRe := regexp.MustCompile(`<xxhash64be>(.*?)<`)

	// Extract reel name
	if match := reelNameRe.FindStringSubmatch(content); len(match) > 1 {
		logData.ReelName = match[1]
	}

	// Extract times
	if match := startTimeRe.FindStringSubmatch(content); len(match) > 1 {
		timeStr := match[1]
		timeStr = strings.Replace(timeStr, "T", " ", 1)
		timeStr = strings.TrimSuffix(timeStr, "Z")
		timeStr = strings.Replace(timeStr, "-", "/", -1)
		logData.StartTime = timeStr
	}

	if match := finishTimeRe.FindStringSubmatch(content); len(match) > 1 {
		timeStr := match[1]
		timeStr = strings.Replace(timeStr, "T", " ", 1)
		timeStr = strings.TrimSuffix(timeStr, "Z")
		timeStr = strings.Replace(timeStr, "-", "/", -1)
		logData.FinishTime = timeStr
	}

	// Extract files
	fileMatches := fileRe.FindAllStringSubmatch(content, -1)
	sizeMatches := sizeRe.FindAllStringSubmatch(content, -1)
	hashMatches := hashRe.FindAllStringSubmatch(content, -1)

	for i, fileMatch := range fileMatches {
		if len(fileMatch) > 1 {
			fileName := filepath.Base(fileMatch[1])
			fileInfo := FileInfo{
				Name:     strings.TrimSuffix(fileName, filepath.Ext(fileName)),
				FileType: getSuffix(fileName),
				Status:   "verified",
			}

			if i < len(sizeMatches) && len(sizeMatches[i]) > 1 {
				if size, err := strconv.ParseInt(sizeMatches[i][1], 10, 64); err == nil {
					logData.TotalSize += size
					fileInfo.FileSize = humanReadableSize(size)
				}
			}

			if i < len(hashMatches) && len(hashMatches[i]) > 1 {
				fileInfo.HashValue = fmt.Sprintf("%s: %s", logData.HashType, hashMatches[i][1])
			}

			logData.Files = append(logData.Files, fileInfo)
		}
	}

	logData.ClipsOverview = calculateOverview(logData.Files, logData.TotalSize)
	return logData, nil
}

func calculateOverview(files []FileInfo, totalSize int64) ClipsOverview {
	videoExts := map[string]bool{
		"MOV": true, "MXF": true, "MP4": true, "ARI": true, "ARX": true,
		"MPG": true, "MPEG": true, "DNG": true, "BRAW": true, "CRM": true,
		"MTS": true, "VRW": true, "R3D": true, "CINE": true, "AVI": true,
		"MKV": true, "RMF": true, "KRW": true, "RED": true, "KWV": true,
	}

	audioExts := map[string]bool{
		"WAV": true, "MP3": true, "AAC": true, "M4A": true,
		"APE": true, "FLAC": true, "WMA": true,
	}

	overview := ClipsOverview{}

	for _, file := range files {
		ext := strings.ToUpper(file.FileType)

		if videoExts[ext] {
			overview.VideoFiles++
		} else if audioExts[ext] {
			overview.AudioFiles++
		} else {
			overview.OtherFiles++
		}
	}

	overview.TotalFiles = len(files)
	overview.TotalSize = humanReadableSize(totalSize)

	// 简化处理
	if overview.VideoFiles > 0 {
		overview.VideoSize = overview.TotalSize
	} else {
		overview.VideoSize = "0 B"
	}

	if overview.AudioFiles > 0 {
		overview.AudioSize = "0 B"
	} else {
		overview.AudioSize = "0 B"
	}

	if overview.OtherFiles > 0 {
		overview.OtherSize = "0 B"
	} else {
		overview.OtherSize = "0 B"
	}

	return overview
}
