package work_modules

import (
	"bufio"
	"fmt"
	"github.com/zeebo/xxh3"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func RunCleaner() {

	PrintInfo()
	fmt.Print("Starting Cleaner...")

	for _, path := range filePathList {
		switch cleanType {
		case "1":
			cleanerOutputFilesMap[path] = GetRunDir() + strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)) + "_cleaned.txt"
		case "2":
			cleanerOutputFilesMap[path] = GetRunDir() + "cleaned.txt"
		}
	}

	var err error
	if partsPattern, err = regexp.Compile(`.+` + delimetr + `.+` + delimetr + `.+`); err != nil {
		PrintErr()
		fmt.Print("Error compiling part pattern regex: ", err, ": Cleaning without part check\n")
		cleanerPartsPatternIsErr = true
	}

	cleanerStringHashMap = make(map[uint64]bool)

	fmt.Print("\r")
	PrintSuccess()
	fmt.Print("Cleaner started\n\n")

}

func Cleaner(path string) {
	currPath = path
	_, currPathCut = filepath.Split(currPath)
	cleanerResultChannelMap[currPath] = make(chan string)
	TMPlinesLen = 0
	currFileDubles = 0
	currFileWritedString = 0
	currFileInvalidLen = 0

	if cleanType == "1" {
		cleanerStringHashMap = make(map[uint64]bool)
	}

	if err := GetCurrentFileSize(currPath); err != nil {
		PrintFileReadErr(currPath, err)
		return
	}

	PrintFileInfo(currPathCut)
	fileDecoder = GetFileProcessInfo(currPath)

	cleanerWrFile, err := os.OpenFile(cleanerOutputFilesMap[currPath], os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		PrintFileReadErr(currPath, err)
		return
	} else {
		cleanerWriteFile = bufio.NewWriter(transform.NewWriter(cleanerWrFile, unicode.UTF8.NewDecoder()))
	}

	cleanerReadFile, err := os.OpenFile(currPath, os.O_RDWR, os.ModePerm)
	if err != nil {
		PrintFileReadErr(currPath, err)
		return
	}

	scanner := bufio.NewScanner(transform.NewReader(cleanerReadFile, fileDecoder))
	isFileInProcessing = true
	pBar = CreatePBar()

	go PBarUpdater()
	go CleanerWriteLine()

	if cleanerPartsPatternIsErr {
		for ; scanner.Scan(); TMPlinesLen++ {
			line := scanner.Text()
			if validPattern.MatchString(line) && !unknownPattern.MatchString(line) {
				hash := xxh3.HashString(line)
				if _, ok := cleanerStringHashMap[hash]; !ok {
					cleanerStringHashMap[hash] = true
					cleanerResultChannelMap[currPath] <- line
				} else {
					currFileDubles++
				}
			} else {
				currFileInvalidLen++
			}
		}
	} else {
		for ; scanner.Scan(); TMPlinesLen++ {
			line := scanner.Text()
			if validPattern.MatchString(line) && !unknownPattern.MatchString(line) && partsPattern.MatchString(line) {
				hash := xxh3.HashString(line)
				if _, ok := cleanerStringHashMap[hash]; !ok {
					cleanerStringHashMap[hash] = true
					cleanerResultChannelMap[currPath] <- line
				} else {
					currFileDubles++
				}
			} else {
				currFileInvalidLen++
			}
		}
	}

	isFileInProcessing = false              // Stop progress bar
	close(cleanerResultChannelMap[currPath]) //
	checkedLines += int64(TMPlinesLen)      // Add lines
	cleanerDublesLen += currFileDubles      //
	cleanerWritedString += currFileWritedString //
	cleanerInvalidLen += currFileInvalidLen //
	_ = pBar.Finish()                       // Finish bar
	_ = pBar.Exit()                         // Close bar
	cleanerReadFile.Close()                 // Close file
	cleanerWrFile.Close()                   // Close file
	cleanerResultChannelMap[currPath] = nil //
	if cleanType == "1" {
		cleanerStringHashMap = nil // Clear map if databases are sorted separately
	}
	PrintClearInfo()           //
	PrintFileDone(currPathCut) // Output file sorted
	checkedFiles++             // Increment processed files count
}

func CleanerWriteLine() {
	for {
		if data, ok := <-cleanerResultChannelMap[currPath]; ok {
			_, _ = cleanerWriteFile.WriteString(data + "\n")
			currFileWritedString++
			continue
		} else {
			if err := cleanerWriteFile.Flush(); err != nil {
				PrintErr()
				fmt.Print("Error flushing buffer to file: ", err, "\n")
			}
			break
		}
	}
}
