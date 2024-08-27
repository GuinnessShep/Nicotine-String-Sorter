package work_modules

import (
	"fmt"
	"github.com/saintfish/chardet"
	"github.com/schollz/progressbar/v3"
	"os"
	"time"
)

// PrintErr PrintSuccess PrintWarn PrintInfo Icons

func PrintErr() {
	fmt.Print("[")
	ColorRed.Print("-")
	fmt.Print("] ")
}

func PrintSuccess() {
	fmt.Print("[")
	ColorGreen.Print("+")
	fmt.Print("] ")
}

func PrintWarn() {
	fmt.Print("[")
	ColorYellow.Print("*")
	fmt.Print("] ")
}

func PrintInfo() {
	fmt.Print("[")
	ColorMagenta.Print("*")
	fmt.Print("] ")
}

// PrintLinesChunk PrintCheckedFiles PrintFileInfo PrintFileDone Sorter Info

func PrintCheckedFiles() {
	fmt.Print("[")
	ColorBlue.Print(checkedFiles + 1)
	fmt.Print("/")
	ColorBlue.Print(len(filePathList))
	fmt.Print("] ")
}

func PrintFileInfo(path string) {
	PrintInfo()
	PrintCheckedFiles()
	fmt.Print("Processing file ")
	ColorBlue.Print(path)
	fmt.Print(" : ")
	if currentFileSize < 1610612736 {
		ColorBlue.Print(currentFileSize / 1048576)
		fmt.Print(" MB : ")
	} else {
		ColorBlue.Print(currentFileSize / 1073741824)
		fmt.Print(" GB : ")
	}
	ColorBlue.Print("~", currentFileLines)
	fmt.Print(" Lines\n")
}

func PrintFileDone(path string) {
	PrintSuccess()
	PrintCheckedFiles()
	ColorBlue.Print(path)
	fmt.Print(" : File processed\n\n")
}

func PrintSortInfo() {
	fmt.Print("\n")
	switch {
	case reqLen <= 10:
		for _, request := range searchRequests {
			PrintSuccess()
			ColorBlue.Print(request)
			fmt.Print(" : ")
			ColorBlue.Print(sorterRequestStatMapCurrFile[request])
			fmt.Print(" lines\n")

		}
	case reqLen > 10:
		PrintSuccess()
		fmt.Print("Found ")
		ColorBlue.Print(currFileMatchLines)
		fmt.Print(" matching lines for all requests\n")
	}
}

func PrintClearInfo() {
	fmt.Print("\n")
	PrintInfo()
	ColorBlue.Print(TMPlinesLen)
	fmt.Print(" lines : ")
	ColorBlue.Print(currFileWritedString)
	fmt.Print(" Unique : ")
	ColorBlue.Print(currFileDubles)
	fmt.Print(" Duplicates : ")
	ColorBlue.Print(currFileInvalidLen)
	fmt.Print(" Invalid\n")
}

func PrintChunk() {
	PrintInfo()
	fmt.Print("Reading ")
	if GetAviableStringsCount() > currentFileLines {
		ColorBlue.Print(currentFileLines)
	} else {
		ColorBlue.Print(GetAviableStringsCount())
	}
	fmt.Print(" lines : ")
}

func PrintEncoding(result *chardet.Result) {
	ColorBlue.Print(result.Charset)
	fmt.Print(" - ")
	ColorBlue.Print(result.Confidence)
	fmt.Print(" %\n")
}

func CreatePBar() *progressbar.ProgressBar {
	return progressbar.NewOptions(
		int(currentFileLines),
		progressbar.OptionSetWidth(50),
		progressbar.OptionSetItsString("Str"),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetDescription("[*]"),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[blue]█[reset]",
			SaucerHead:    "[green]░[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)
}

func PBarUpdater() {
	for isFileInProcessing {
		if TMPlinesLen > int(currentFileLines) {
			_ = pBar.Set64(currentFileLines)
		} else {
			_ = pBar.Set(TMPlinesLen)
		}
		time.Sleep(time.Millisecond * 250)
	}
}

// Errors

func PrintFileReadErr(path string, err error) {
	PrintErr()
	fmt.Printf("%s : File read error : %s\n\n", path, err)
}

func PrintZeroRequestsErr() {
	PrintErr()
	fmt.Print("No requests for sort : Restart the sorter\n")
	PrintErr()
	fmt.Print("Press ")
	ColorBlue.Print("Enter")
	fmt.Print(" to exit")
	fmt.Scanln()
	os.Exit(1)
}

func PrintResultWriteErr(request string, err error) {
	PrintErr()
	ColorBlue.Print(request)
	fmt.Print(" : Error writing found lines : ")
	ColorRed.Print(err, "\n")
	PrintInfo()
	fmt.Print("Run the sorter with Administrator rights if the error is related to access\n")
}

func PrintEncodingErr(err error) {
	PrintErr()
	fmt.Printf(" Encoding detection error : %s : Using ", err)
	ColorBlue.Print("UTF-8\n")
}

func PrintEndodingLinesEnd() {
	PrintWarn()
	fmt.Print(" Not enough lines to determine encoding : Using ")
	ColorBlue.Print("UTF-8\n")
}

func PrintSorterResult() {

	fmt.Print("\n")
	for _, request := range searchRequests {
		PrintSuccess()
		ColorBlue.Print(request)
		fmt.Print(" : ")
		ColorBlue.Print(sorterRequestStatMap[request])
		fmt.Print(" lines : ")
		if fi, err := os.Stat(requestStructMap[request].resultFile); err == nil {
			fsize := fi.Size()
			switch {
			case fsize < 1048576:
				ColorBlue.Print(fi.Size() / 1024)
				fmt.Print(" KB : ")
			case fsize >= 1048576:
				ColorBlue.Print(fi.Size() / 1048576)
				fmt.Print(" MB : ")
			}

		} else {
			ColorBlue.Print("?")
			fmt.Print(" MB : ")
		}
		ColorBlue.Print(requestStructMap[request].resultFile, "\n")
	}
	fmt.Print("\n\n")

	PrintSuccess()
	fmt.Print("Files sorted : ")
	ColorBlue.Print(checkedFiles)
	fmt.Print(" out of ")
	ColorBlue.Print(len(filePathList), "\n")

	PrintSuccess()
	fmt.Print("Lines sorted : ")
	ColorBlue.Print(checkedLines, "\n")

	PrintSuccess()
	fmt.Print("Matching lines : ")
	ColorGreen.Print(matchLines, "\n")

	PrintWarn()
	fmt.Print("Duplicate lines : ")
	ColorGreen.Print(sorterDubles, "\n\n")
}

func PrintCleanerResult() {
	fmt.Print("\n\n")
	PrintSuccess()
	fmt.Print("Files cleaned : ")
	ColorBlue.Print(checkedFiles)
	fmt.Print(" out of ")
	ColorBlue.Print(len(filePathList), "\n")
	PrintSuccess()
	fmt.Print("Duplicates removed : ")
	ColorBlue.Print(cleanerDublesLen, "\n")
	PrintSuccess()
	fmt.Print("Invalid removed : ")
	ColorBlue.Print(cleanerInvalidLen, "\n")
	PrintSuccess()
	fmt.Print("Unique lines written : ")
	ColorBlue.Print(cleanerWritedString, "\n\n")

	switch cleanType {
	case "1":
		for _, path := range filePathList {
			PrintSuccess()
			if fi, err := os.Stat(cleanerOutputFilesMap[path]); err == nil {
				fsize := fi.Size()
				switch {
				case fsize < 1048576:
					ColorBlue.Print(fi.Size() / 1024)
					fmt.Print(" KB : ")
				case fsize >= 1048576:
					ColorBlue.Print(fi.Size() / 1048576)
					fmt.Print(" MB : ")
				}
			} else {
				ColorBlue.Print("?")
				fmt.Print(" MB : ")
			}
			fmt.Print(cleanerOutputFilesMap[path] + "\n")
		}
	case "2":
		PrintSuccess()
		PrintSuccess()
		if fi, err := os.Stat(cleanerOutputFilesMap[filePathList[0]]); err == nil {
			fsize := fi.Size()
			switch {
			case fsize < 1048576:
				ColorBlue.Print(fi.Size() / 1024)
				fmt.Print(" KB : ")
			case fsize >= 1048576:
				ColorBlue.Print(fi.Size() / 1048576)
				fmt.Print(" MB : ")
			}
		} else {
			ColorBlue.Print("?")
			fmt.Print(" MB : ")
		}
		fmt.Print(cleanerOutputFilesMap[filePathList[0]] + "\n")
	}
}
