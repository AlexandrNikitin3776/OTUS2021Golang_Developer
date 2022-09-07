package main

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

const (
	progressChar = "#"
	emptyChar    = "."
	progressLen  = 80
)

var OverflowError = errors.New("progress bar is full, can't write more bytes")

type ProgressBar struct {
	bytesWrote int64
	totalBytes int64
}

func NewProgressBar(offset, limit, fileSize int64) *ProgressBar {
	isLimitSet := limit > 0
	isUndefinedSize := fileSize == 0

	if isUndefinedSize {
		return &ProgressBar{totalBytes: limit}
	}
	if isLimitSet {
		return &ProgressBar{totalBytes: Min(limit, fileSize-offset)}
	}
	return &ProgressBar{totalBytes: fileSize}
}

func (pb *ProgressBar) Write(input []byte) (int, error) {
	if pb.bytesWrote >= pb.totalBytes {
		return 0, OverflowError
	}
	pb.bytesWrote += int64(len(input))
	percents := float64(pb.bytesWrote) / float64(pb.totalBytes) * 100
	progressRepeats := int(math.Round(percents * progressLen / 100))
	progress := strings.Repeat(progressChar, progressRepeats)
	empty := strings.Repeat(emptyChar, progressLen-progressRepeats)
	fmt.Printf("\r%v%v\t%.1f%%", progress, empty, percents)
	if pb.bytesWrote >= pb.totalBytes {
		fmt.Println()
	}

	return len(input), nil
}
