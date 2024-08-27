package user_modules

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func GetWorkMode() (work string) {
	PrintWorkModes()

LoopWork:
	for {
		PrintInput()
		fmt.Print("Select a work mode: ")  // Reworded
		wmraw, _ := userInputReader.ReadString('\n')
		wmraw = strings.TrimSpace(wmraw)

		switch wmraw {
		case "1":
			work = "sorter"
			break LoopWork
		case "2":
			work = "cleaner"
			break LoopWork
		case "4":
			os.Exit(0)
		default:
			continue LoopWork
		}
	}
	return work
}

func GetFilesInput() (result []string) {

Loop:
	for {
		PrintInput()
		fmt.Print("Enter the file or folder path for processing: ")  // Reworded

		rawPath, _ := userInputReader.ReadString('\n')
		rawPath = strings.TrimSpace(rawPath)
		if rawPath == "" {
			continue Loop
		}

		rawPath = filepath.Clean(rawPath)

		if fileInfo, fierr := os.Stat(rawPath); fierr == nil {

			if fileInfo.IsDir() {
				PrintSuccess()
				fmt.Printf("Folder '")  // Reworded
				ColorBlue.Print(rawPath)
				fmt.Print("' found:\n\n")  // Reworded

				_ = filepath.Walk(rawPath, func(path string, info os.FileInfo, fwerr error) error {

					if fwerr != nil {
						PrintErr()
						fmt.Print(fwerr, "\n")
						return fwerr
					}

					if !info.IsDir() {
						if filepath.Ext(path) == ".txt" {
							fmt.Printf("    %s\n", path)
							result = append(result, path)
						}
					}
					return nil
				})

				if len(result) >= 1 {
					fmt.Print("\n")
					break Loop
				} else {
					PrintErr()
					fmt.Print("No files found for processing\n")  // Reworded
					continue Loop
				}

			} else {
				PrintSuccess()
				fmt.Print("File found\n\n")  // Reworded
				result = append(result, rawPath)
				break Loop
			}

		} else {
			PrintErr()
			fmt.Printf("Path '%s' does not exist\n", rawPath)  // Reworded
			continue Loop
		}
	}

	result = Unique(result)
	GetFilesSize(result)
	return result
}

func GetRequestsInput() (requests []string) {

	PrintInfo()
	fmt.Print("Supported input methods:\n\n")  // Reworded
	ColorBlue.Print("       1")
	fmt.Print(" - Enter from terminal\n")  // Reworded
	ColorBlue.Print("       2")
	fmt.Print(" - Load from file\n\n")  // Reworded

LoopA:
	for {

		PrintInput()
		fmt.Print("Select a request input method: ")  // Reworded

		inputType, _ := userInputReader.ReadString('\n')

		switch strings.TrimSpace(inputType) {
		case "1":
		LoopB:
			for true {
				PrintInput()
				fmt.Print("Enter requests separated by space: ")  // Reworded
				rawRequests, _ := userInputReader.ReadString('\n')
				rawRequests = strings.TrimSpace(rawRequests)
				if rawRequests == "" {
					continue LoopB
				}
				for _, request := range strings.Split(rawRequests, " ") {
					request = strings.TrimSpace(strings.ToLower(request))
					_, err := regexp.Compile(".*" + request + ".*:(.+:.+)")
					if err != nil {
						PrintErr()
						fmt.Printf("%s : Failed to create regular expression : %s\n", request, err)  // Reworded
						continue LoopB
					}
					requests = append(requests, request)
				}

				if len(requests) == 0 {
					PrintErr()
					fmt.Print("No requests entered for search\n")  // Reworded
					continue LoopB
				}
				fmt.Print("\n")
				break LoopA
			}
		case "2":
		LoopC:
			for true {
				PrintInput()
				fmt.Print("Enter the file path: ")  // Reworded
				rawRequests, _ := userInputReader.ReadString('\n')
				rawRequests = strings.TrimSpace(rawRequests)
				_, sterr := os.Stat(rawRequests)
				if sterr != nil {
					PrintErr()
					fmt.Print("File does not exist\n")  // Reworded
					continue LoopC
				}
				file, operr := os.Open(rawRequests)
				if operr != nil {
					PrintErr()
					fmt.Printf("Error reading request file : %s\n", operr)  // Reworded
					fmt.Println(operr)
					continue LoopC
				}

				defer file.Close()

				scanner := bufio.NewScanner(file)
				scanner.Split(bufio.ScanLines)

				for scanner.Scan() {
					request := strings.TrimSpace(strings.ToLower(scanner.Text()))
					_, err := regexp.Compile(regexp.QuoteMeta(request) + ".*:.+:.+")
					if err != nil {
						PrintErr()
						fmt.Printf("%s : Failed to create regular expression : %s\n", request, err)  // Reworded
						continue LoopC
					}
					requests = append(requests, request)

				}

				PrintSuccess()
				fmt.Print("Request file found: ")  // Reworded
				ColorBlue.Print(len(requests))
				fmt.Print(" requests\n")  // Reworded

				if len(requests) == 0 {
					PrintErr()
					fmt.Print("No requests found for search\n")  // Reworded
					continue LoopA
				}
				fmt.Print("\n")
				break LoopA
			}
		default:
			continue LoopA
		}
	}
	return Unique(requests)
}

func GetSaveTypeInput() (saveType string) {

	PrintInfo()
	fmt.Print("Supported save formats:\n\n")  // Reworded
	ColorBlue.Print("       1")
	fmt.Print(" - Log:Pass\n")
	ColorBlue.Print("       2")
	fmt.Print(" - Url:Log:Pass\n\n")

Loop:
	for true {
		PrintInput()
		fmt.Print("Select save format: ")  // Reworded
		rawSaveType, _ := userInputReader.ReadString('\n')
		rawSaveType = strings.TrimSpace(rawSaveType)

		switch rawSaveType {
		case "1", "2":
			saveType = rawSaveType
			fmt.Print("\n")
			break Loop
		default:
			continue Loop
		}
	}
	return saveType
}

func GetCleanTypeInput() (cleanType string) {

	PrintInfo()
	fmt.Print("Supported cleaner modes:\n\n")  // Reworded
	ColorBlue.Print("       1")
	fmt.Print(" - Clean and save each database separately\n")  // Reworded
	ColorBlue.Print("       2")
	fmt.Print(" - Clean all databases together and save to one file\n\n")  // Reworded

Loop:
	for true {
		PrintInput()
		fmt.Print("Select cleaner mode: ")  // Reworded
		rawcleanType, _ := userInputReader.ReadString('\n')
		rawcleanType = strings.TrimSpace(rawcleanType)

		switch rawcleanType {
		case "1", "2":
			cleanType = rawcleanType
			fmt.Print("\n")
			break Loop
		default:
			continue Loop
		}
	}
	return cleanType
}

func GetDelimiterInput() (delimiter string) {

LoopDel:
	for true {
		PrintInput()
		fmt.Print("Enter the row delimiter: ") // Reworded
		var rawDelTrim string

		rawDel, _ := userInputReader.ReadString('\n')

		switch rawDel {
		case "":
			continue LoopDel
		case " ":
			rawDelTrim = rawDel
		default:
			rawDelTrim = strings.TrimSpace(rawDel)
		}

		PrintInfo()
		fmt.Print("Row delimiter: '")  // Reworded
		ColorBlue.Print(rawDelTrim)
		fmt.Print("'\n\n")
		ColorBlue.Print("       1")
		fmt.Print(" - Continue\n")  // Reworded
		ColorBlue.Print("       2")
		fmt.Print(" - Enter again\n\n")  // Reworded
	LoopAction:
		for true {
			PrintInput()
			fmt.Print("Select action: ")  // Reworded
			action, _ := userInputReader.ReadString('\n')
			action = strings.TrimSpace(action)
			switch action {
			case "1":
				delimiter = rawDelTrim
				break LoopDel
			case "2":
				continue LoopDel
			default:
				continue LoopAction
			}
		}
	}

	return
}
