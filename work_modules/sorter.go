package work_modules

import (
	"bufio"
	"fmt"
	"github.com/zeebo/xxh3"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

func RunSorter() {

	PrintInfo()
	fmt.Print("Starting sorter...")

	var (
		compiledRegEx *regexp.Regexp
		err           error
	)

	for _, request := range searchRequests {

		switch saveType {
		case "1":
			compiledRegEx, err = regexp.Compile(".*" + regexp.QuoteMeta(request) + ".*" + delimetr + "(.+" + delimetr + ".+)")
		case "2":
			compiledRegEx, err = regexp.Compile("(.*" + regexp.QuoteMeta(request) + ".*" + delimetr + ".+" + delimetr + ".+)")
		}

		if err != nil {
			PrintErr()
			fmt.Printf("%s: Request compilation error: %s\n", request, err)
			RemoveFromSliceByValue(searchRequests, request)
			continue
		}

		currentStruct := new(Work)
		currentStruct.requestPattern = compiledRegEx
		currentStruct.resultFile = runDir + fileBadSymbolsPattern.ReplaceAllString(request, "_") + ".txt"
		requestStructMap[request] = currentStruct

		if resultFile, err := os.OpenFile(requestStructMap[request].resultFile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0755); err == nil {
			sorterResultFileMap[request] = resultFile
			sorterResultWriterMap[request] = bufio.NewWriter(transform.NewWriter(resultFile, unicode.UTF8.NewDecoder()))
		} else {
			RemoveFromSliceByValue(searchRequests, request)
			PrintResultWriteErr(request, err)
		}
	}

	if len(requestStructMap) == 0 || len(sorterResultFileMap) == 0 || len(sorterResultWriterMap) == 0 {
		PrintZeroRequestsErr()
	}

	fmt.Print("\r")
	PrintSuccess()
	fmt.Print("Sorter started\n\n")
}

func Sorter(path string) {

	_, currPathCut = filepath.Split(path)
	currPath = path
	sorterStringHashMap = make(map[uint64]bool)
	isFileInProcessing = false
	isResultWrited = false
	TMPlinesLen = 0
	currFileDubles = 0
	currFileMatchLines = 0
	for _, req := range searchRequests {
		sorterRequestStatMapCurrFile[req] = 0
	}

	if err := GetCurrentFileSize(currPath); err != nil {
		PrintFileReadErr(currPath, err)
		return
	}

	PrintFileInfo(currPathCut)
	fileDecoder = GetFileProcessInfo(currPath)

	file, err := os.OpenFile(currPath, os.O_RDWR, os.ModePerm)

	if err != nil {
		PrintFileReadErr(currPath, err)
		return
	}

	if GetAviableStringsCount() > currentFileLines {
		sorterPool.Tune(int(math.Round(float64(currentFileLines) / 3)))
	} else {
		 sorterPool.Tune(int(math.Round(float64(GetAviableStringsCount()) / 3)))
	}

	isFileInProcessing = true
	pBar = CreatePBar()
	go PBarUpdater()
	go SorterWriteResult()

	scanner := bufio.NewScanner(transform.NewReader(file, fileDecoder))

	for ; scanner.Scan(); TMPlinesLen++ {
		workWG.Add(1)
		_ = sorterPool.Invoke(scanner.Text())
	}

	workWG.Wait()

	checkedLines += int64(TMPlinesLen)     // Add lines
	_ = pBar.Finish()                      // Finish the progress bar
	_ = pBar.Exit()                        // Close the progress bar
	close(sorterWriteChannelMap[currPath]) // 

	isFileInProcessing = false
	for !isResultWrited {
		time.Sleep(time.Millisecond * 100)
	}

	file.Close() // Close file

	sorterWriteChannelMap[currPath] = nil // Clear the channel
	PrintSortInfo()
	PrintFileDone(currPathCut)       // File sorted
	checkedFiles++                   // Increment processed files count
	matchLines += currFileMatchLines // Sum matched lines
	sorterDubles += currFileDubles   //
}

func SorterProcessLine(line string) {
	defer workWG.Done()
	for _, request := range searchRequests {
		if result := requestStructMap[request].requestPattern.FindStringSubmatch(line); len(result) == 2 {
			sorterWriteChannelMap[currPath] <- [2]string{request, result[1]}
			return
		}
	}
}

func SorterWriteResult() {

	for {
		if data, ok := <-sorterWriteChannelMap[currPath]; ok {
			hash := xxh3.HashString(data[1])
			if _, ok := sorterStringHashMap[hash]; !ok {
				sorterStringHashMap[hash] = true
				sorterResultWriterMap[data[0]].WriteString(data[1] + "\n")
				sorterRequestStatMapCurrFile[data[0]]++
			} else {
				currFileDubles++
			}
			continue
		} else {
			break
		}
	}

	for _, request := range searchRequests {
		if err := sorterResultWriterMap[request].Flush(); err != nil {
			PrintErr()
			fmt.Print("Error writing buffer to file: ", err, "\n")
		}
		currFileMatchLines += sorterRequestStatMapCurrFile[request]
		sorterRequestStatMap[request] += sorterRequestStatMapCurrFile[request]
	}

	isResultWrited = true // Notify that the file has been written
}

func SorterEnd() {
	for _, request := range searchRequests {
		sorterResultFileMap[request].Close()
		if stat, err := os.Stat(requestStructMap[request].resultFile); err == nil && stat.Size() == 0 {
			os.Remove(requestStructMap[request].resultFile)
		}
	}

	if IsDirEmpty(runDir) {
		os.Remove(runDir)
	}
}
