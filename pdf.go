package main

import (
	"fmt"

	"github.com/jung-kurt/gofpdf"
)

func generatePDF(outputPath string, logData *LogData, request ReportRequest) error {
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.AddPage()

	// Header
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(200, 10, "Clips Report")
	pdf.SetXY(250, 10)
	pdf.Cell(40, 10, logData.GenerationTime)
	pdf.Ln(20)

	// Title
	pdf.SetFont("Arial", "B", 20)
	pdf.Cell(0, 15, logData.ReelName)
	pdf.Ln(20)

	// Project info
	pdf.SetFont("Arial", "", 10)
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
	pdf.SetFont("Arial", "B", 11)
	pdf.Cell(0, 10, "Clips Overview")
	pdf.Ln(12)

	// Overview table
	pdf.SetFont("Arial", "B", 10)
	pdf.SetFillColor(240, 240, 240)
	pdf.CellFormat(60, 8, "", "", 0, "", true, 0, "")
	pdf.CellFormat(30, 8, "Clips", "", 0, "C", true, 0, "")
	pdf.CellFormat(30, 8, "Files", "", 0, "C", true, 0, "")
	pdf.CellFormat(40, 8, "Size", "", 1, "C", true, 0, "")

	pdf.SetFont("Arial", "", 10)
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
	pdf.SetFont("Arial", "B", 9)
	pdf.SetFillColor(240, 240, 240)
	pdf.CellFormat(50, 8, "Name", "", 0, "", true, 0, "")
	pdf.CellFormat(20, 8, "Type", "", 0, "C", true, 0, "")
	pdf.CellFormat(25, 8, "Size", "", 0, "C", true, 0, "")
	pdf.CellFormat(70, 8, "Hash Values", "", 0, "", true, 0, "")
	pdf.CellFormat(25, 8, "Backups", "", 0, "C", true, 0, "")
	pdf.CellFormat(20, 8, "Status", "", 1, "C", true, 0, "")

	// File details rows
	pdf.SetFont("Arial", "", 8)
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
		if pdf.GetY() > 180 {
			pdf.AddPage()
		}
	}

	return pdf.OutputFileAndClose(outputPath)
}
