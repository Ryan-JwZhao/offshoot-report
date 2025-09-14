package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jung-kurt/gofpdf"
)

func generatePDF(outputPath string, logData *LogData, request ReportRequest) error {
	//fmt.Printf("--- DEBUG: Received request in Go ---\n")
	//fmt.Printf("ProjectTitle: '%s'\n", request.ProjectTitle)
	//fmt.Printf("Backups: '%s'\n", request.Backups)
	//fmt.Printf("FilePaths: %v\n", request.FilePaths)
	//fmt.Printf("------------------------------------\n")

	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.SetMargins(10, 10, 10)
	pdf.SetAutoPageBreak(true, 10)
	pdf.AddPage()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err) // 或者更优雅的错误处理
	}

	// 构造完整的字体路径
	fontPath_Normal := filepath.Join(homeDir, "Library", "Fonts", "MapleMonoNormalNL-CN-Light.ttf")
	fontPath_Bold := filepath.Join(homeDir, "Library", "Fonts", "MapleMonoNormalNL-CN-Medium.ttf")

	// Header

	pdf.AddUTF8Font("MapleMono", "", fontPath_Normal)
	pdf.AddUTF8Font("MapleMono", "B", fontPath_Bold)

	pdf.SetFont("MapleMono", "", 12)
	pdf.SetXY(10, 10)
	pdf.Cell(200, 10, "Clips Report")
	pdf.SetXY(250, 10)
	pdf.Cell(40, 10, logData.GenerationTime)
	pdf.Ln(15)

	// Title
	pdf.SetFont("MapleMono", "B", 20)
	pdf.Cell(0, 15, logData.ReelName)
	pdf.Ln(20)

	// 在 generatePDF 中，就在 pdf.Cell(...) 之前
	//fmt.Printf("--- DEBUG (inside generatePDF): Drawing ProjectTitle Cell ---\n")
	//fmt.Printf("ProjectTitle string passed to pdf.Cell: '%s'\n", request.ProjectTitle)
	//fmt.Printf("Current PDF X, Y before Cell: %f, %f\n", pdf.GetX(), pdf.GetY())
	//fmt.Printf("----------------------------------------------------------------\n")

	// Project info
	pdf.SetFont("MapleMono", "", 10)
	pdf.SetTextColor(100, 100, 100)
	pdf.Cell(0, 8, request.ProjectTitle)
	pdf.Ln(8)

	if logData.StartTime != "" && logData.FinishTime != "" {
		pdf.Cell(0, 8, fmt.Sprintf("Offloaded between %s", logData.StartTime))
		pdf.Ln(8)
		pdf.Cell(0, 8, fmt.Sprintf("and %s", logData.FinishTime))
		pdf.Ln(15)
	}

	// Clips Overview
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("MapleMono", "B", 11)
	pdf.Cell(0, 10, "Clips Overview")
	pdf.Ln(12)

	// Overview table
	pdf.SetFont("MapleMono", "B", 10)
	pdf.SetFillColor(240, 240, 240)
	pdf.CellFormat(60, 8, "", "", 0, "", true, 0, "")
	pdf.CellFormat(30, 8, "Clips", "", 0, "C", true, 0, "")
	pdf.CellFormat(30, 8, "Files", "", 0, "C", true, 0, "")
	pdf.CellFormat(40, 8, "Size", "", 1, "C", true, 0, "")

	pdf.SetFont("MapleMono", "", 10)
	pdf.SetFillColor(255, 255, 255)

	// Video clips row
	pdf.CellFormat(60, 7, "Video Clips", "", 0, "", true, 0, "")
	pdf.CellFormat(30, 7, fmt.Sprintf("%d", logData.ClipsOverview.VideoFiles), "", 0, "C", true, 0, "")
	pdf.CellFormat(30, 7, fmt.Sprintf("%d", logData.ClipsOverview.VideoFiles), "", 0, "C", true, 0, "")
	pdf.CellFormat(40, 7, logData.ClipsOverview.VideoSize, "", 1, "C", true, 0, "")

	// Audio clips row
	pdf.SetFillColor(240, 240, 240)
	pdf.CellFormat(60, 7, "Audio Clips", "", 0, "", true, 0, "")
	pdf.CellFormat(30, 7, fmt.Sprintf("%d", logData.ClipsOverview.AudioFiles), "", 0, "C", true, 0, "")
	pdf.CellFormat(30, 7, fmt.Sprintf("%d", logData.ClipsOverview.AudioFiles), "", 0, "C", true, 0, "")
	pdf.CellFormat(40, 7, logData.ClipsOverview.AudioSize, "", 1, "C", true, 0, "")

	// Other clips row
	pdf.SetFillColor(255, 255, 255)
	pdf.CellFormat(60, 7, "Other Files", "", 0, "", true, 0, "")
	pdf.CellFormat(30, 7, "0", "", 0, "C", true, 0, "")
	pdf.CellFormat(30, 7, fmt.Sprintf("%d", logData.ClipsOverview.OtherFiles), "", 0, "C", true, 0, "")
	pdf.CellFormat(40, 7, logData.ClipsOverview.OtherSize, "", 1, "C", true, 0, "")

	// Total row
	pdf.SetFillColor(240, 240, 240)
	pdf.CellFormat(60, 7, "Total", "", 0, "", true, 0, "")
	totalClips := logData.ClipsOverview.VideoFiles + logData.ClipsOverview.AudioFiles
	pdf.CellFormat(30, 7, fmt.Sprintf("%d", totalClips), "", 0, "C", true, 0, "")
	pdf.CellFormat(30, 7, fmt.Sprintf("%d", logData.ClipsOverview.TotalFiles), "", 0, "C", true, 0, "")
	pdf.CellFormat(40, 7, logData.ClipsOverview.TotalSize, "", 1, "C", true, 0, "")

	pdf.Ln(15)

	// File details table header
	pdf.SetFont("MapleMono", "B", 9)
	pdf.SetFillColor(240, 240, 240)
	pdf.CellFormat(50, 8, "Name", "", 0, "", true, 0, "")
	pdf.CellFormat(20, 8, "Type", "", 0, "C", true, 0, "")
	pdf.CellFormat(25, 8, "Size", "", 0, "C", true, 0, "")
	pdf.CellFormat(70, 8, "Hash Values", "", 0, "", true, 0, "")
	pdf.CellFormat(25, 8, "Backups", "", 0, "C", true, 0, "")
	pdf.CellFormat(20, 8, "Status", "", 1, "C", true, 0, "")

	// File details rows
	pdf.SetFont("MapleMono", "", 8)
	for i, file := range logData.Files {
		if i%2 == 0 {
			pdf.SetFillColor(255, 255, 255)
		} else {
			pdf.SetFillColor(247, 247, 247)
		}

		pdf.CellFormat(50, 6, file.Name, "", 0, "", true, 0, "")
		pdf.CellFormat(20, 6, file.FileType, "", 0, "C", true, 0, "")
		pdf.CellFormat(25, 6, file.FileSize, "", 0, "C", true, 0, "")
		pdf.CellFormat(70, 6, file.HashValue, "", 0, "", true, 0, "")
		pdf.CellFormat(25, 6, request.Backups, "", 0, "C", true, 0, "")
		pdf.CellFormat(20, 6, file.Status, "", 1, "C", true, 0, "")

		// Check if we need a new page
		if pdf.GetY() > 200 {
			pdf.AddPage()
		}
	}

	return pdf.OutputFileAndClose(outputPath)
}
