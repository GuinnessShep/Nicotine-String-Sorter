package user_modules

import (
	"fmt"
	"github.com/klauspost/cpuid/v2"
	"github.com/pbnjay/memory"
	"math"
	"os"
	"runtime"
	"strings"
	"time"
)

func PrintLogoStart(appVersion string) {

	ColorBlue.Print(`
     _   _   _                  _     _                        
    | \ | | (_)                | |   (_)                 
    |  \| |  _    ___    ___   | |_   _   _ __     ___    
    | .   | | |  / __|  / _ \  | __| | | | '_ \   / _ \   
    | |\  | | | | (__  | (_) | | |_  | | | | | | |  __/  
    |_| \_| |_|  \___|  \___/   \__| |_| |_| |_|  \___|
														`)
	time.Sleep(300 * time.Millisecond)
	ColorBlue.Print(`
     _____   _           _                       _____                  _ 
    / ____| | |         (_)                     / ____|                | |              
   | (___   | |_   _ __   _   _ __     __ _    | (___     ___    _ __  | |_    ___   _ __ 
    \___ \  | __| | '__| | | | '_ \   / _  |    \___ \   / _ \  | '__| | __|  / _ \ | '__|
    ____) | | |_  | |    | | | | | | | (_| |    ____) | | (_) | | |    | |_  |  __/ | |   
   |_____/   \__| |_|    |_| |_| |_|  \__, |   |_____/   \___/  |_|     \__|  \___| |_|   
                                       __/ |    
                                      |___/  

`)
	time.Sleep(150 * time.Millisecond)
	ColorMagenta.Print("    v", appVersion)
	fmt.Print(" | ")
	ColorMagenta.Print(runtime.Version())
	ColorBlue.Print("     #")
	fmt.Print(" t.me/rx580_work     ")
	ColorGreen.Print("#")
	fmt.Print(" zelenka.guru/rx580    # НикотиновыйКодер\n\n")
	PrintInfo()
	fmt.Print(cpuid.CPU.BrandName, " @ ", cpuid.CPU.PhysicalCores, "/", cpuid.CPU.LogicalCores, " threads | ")
	fmt.Print(math.Round(float64(memory.FreeMemory()/1073741824)), "/", math.Ceil(float64(memory.TotalMemory()/1073741824)), " GB available memory\n\n")
	isLogoPrinted = true
}

func PrintLogoFast(appVersion string) {

	ColorBlue.Print(`
     _   _   _                  _     _                        
    | \ | | (_)                | |   (_)                 
    |  \| |  _    ___    ___   | |_   _   _ __     ___    
    | .   | | |  / __|  / _ \  | __| | | | '_ \   / _ \   
    | |\  | | | | (__  | (_) | | |_  | | | | | | |  __/  
    |_| \_| |_|  \___|  \___/   \__| |_| |_| |_|  \___|

     _____   _           _                       _____                  _ 
    / ____| | |         (_)                     / ____|                | |              
   | (___   | |_   _ __   _   _ __     __ _    | (___     ___    _ __  | |_    ___   _ __ 
    \___ \  | __| | '__| | | | '_ \   / _  |    \___ \   / _ \  | '__| | __|  / _ \ | '__|
    ____) | | |_  | |    | | | | | | | (_| |    ____) | | (_) | | |    | |_  |  __/ | |   
   |_____/   \__| |_|    |_| |_| |_|  \__, |   |_____/   \___/  |_|     \__|  \___| |_|   
                                       __/ |    
                                      |___/ 
 
`)
	ColorMagenta.Print("    v", appVersion)
	fmt.Print(" | ")
	ColorMagenta.Print(runtime.Version())
	ColorBlue.Print("     #")
	fmt.Print(" t.me/rx580_work     ")
	ColorGreen.Print("#")
	fmt.Print(" zelenka.guru/rx580    # НикотиновыйКодер\n\n")
}

func PrintInputInfo(appVersion string) {
	ClearTerm()
	PrintLogoFast(appVersion)

	PrintInfo()
	fmt.Print("Total files : ")
	ColorBlue.Print(len(filePathList))
	fmt.Print(" : Size : ")
	ColorBlue.Print(filesSize / 1048576)
	fmt.Print(" MB ")
	fmt.Print(": Lines : ")
	ColorBlue.Print("~", filesSize/80, "\n")
}
func PrintSorterData() {
	PrintInfo()
	fmt.Printf("Total requests : ")

	reqLen := len(searchRequests)

	switch {
	case reqLen <= 3:
		ColorBlue.Print(reqLen)
		fmt.Print(" : ")
		for i, req := range searchRequests {
			ColorBlue.Print(req)
			if i != reqLen-1 {
				fmt.Print(", ")
			}
		}
		fmt.Print("\n")
	case reqLen > 3 && reqLen <= 10:
		ColorBlue.Print(reqLen, "\n")
		for _, request := range searchRequests {
			fmt.Println("    ", request)
		}
		fmt.Print("\n")
	case reqLen > 10:
		ColorBlue.Print(reqLen, "\n\n")

	}
	PrintInfo()
	fmt.Print("Save format: ")
	switch saveType {
	case "1":
		ColorBlue.Print("Log:Pass\n")
	case "2":
		ColorBlue.Print("Url:Log:Pass\n")
	}
}

func PrintInputData(appVersion string) (uSelect string) {
LoopData:
	for true {
		PrintInputInfo(appVersion)

		switch workMode {
		case "sorter":
			PrintSorterData()
		case "cleaner":
		}

		PrintInput()
		fmt.Print("Choose an action:\n\n")

		ColorBlue.Print("	1")
		fmt.Print(" - Launch\n")

		ColorBlue.Print("	2")
		fmt.Print(" - Select string delimiter - '")
		ColorBlue.Print(":")
		fmt.Print("' default\n")

		ColorBlue.Print("	3")
		fmt.Print(" - Enter data again\n\n")

	LoopMode:
		for true {
			fmt.Print("> ")
			userSelect, _ := userInputReader.ReadString('\n')
			userSelect = strings.TrimSpace(userSelect)

			switch userSelect {
			case "1":
				uSelect = "continue"
				break LoopMode
			case "2":
				delimetr = GetDelimetrInput()
				continue LoopData
			case "3":
				uSelect = "restart"
				break LoopMode
			default:
				continue LoopMode
			}
		}
		break LoopData
	}
	ClearTerm()
	return uSelect
}

func PrintTimeDuration(duration time.Duration) {
	fmt.Print("\n")
	PrintSuccess()
	fmt.Print("Sort time: ")
	ColorBlue.Print(duration, "\n\n\n")

	PrintInfo()
	fmt.Print("Press ")
	ColorBlue.Print("Enter")
	fmt.Print(" to exit")
	fmt.Scanln()
	os.Exit(0)
}

func PrintInput() {
	fmt.Print("[")
	ColorBlue.Print("#")
	fmt.Print("] ")
}

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

func _() {
	fmt.Print("[")
	ColorYellow.Print("*")
	fmt.Print("] ")
}

func PrintInfo() {
	fmt.Print("[")
	ColorMagenta.Print("*")
	fmt.Print("] ")
}

func PrintWorkModes() {
	PrintInfo()
	fmt.Print("Supported work modes:\n\n")
	ColorBlue.Print("       1")
	fmt.Print(" - String sorter\n")
	fmt.Print("       Search for lines in the database matching requests and write to a ")
	ColorBlue.Print("separate file")
	fmt.Print(" without ")
	ColorBlue.Print("duplicates\n")
	fmt.Print("       The request should be in the format '")
	ColorBlue.Print("google.com")
	fmt.Print("' or '")
	ColorBlue.Print("google")
	fmt.Print("'\n\n")

	ColorBlue.Print("       2")
	fmt.Print(" - Database cleaner from invalid lines and duplicates\n")
	fmt.Print("       Delete duplicates and lines not matching '")
	ColorBlue.Print("A-z")
	fmt.Print(" | ")
	ColorBlue.Print("0-9")
	fmt.Print(" | ")
	ColorBlue.Print("Special characters")
	fmt.Print(" | ")
	ColorBlue.Print("10-256")
	fmt.Print(" characters | without ")
	ColorBlue.Print("UNKNOWN")
	fmt.Print("'\n\n")

	ColorBlue.Print("       4")
	fmt.Print(" - Close the program\n\n")
}
