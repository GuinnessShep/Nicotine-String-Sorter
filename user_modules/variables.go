package user_modules

import (
	"bufio"
	"github.com/fatih/color"
	"os"
	"runtime"
	"strings"
	"sync"
)

var (

	// Colors
	ColorBlue    = color.New(color.FgBlue).Add(color.Bold)
	ColorGreen   = color.New(color.FgGreen).Add(color.Bold)
	ColorRed     = color.New(color.FgRed).Add(color.Bold)
	ColorMagenta = color.New(color.FgMagenta).Add(color.Bold)
	ColorYellow  = color.New(color.FgYellow).Add(color.Bold)

	filesSize       int64                       // Size of all input files
	userInputReader = bufio.NewReader(os.Stdin) // Alternative input reader with space support
	userOs          = runtime.GOOS              // User's OS
	updateWG        sync.WaitGroup              // Update wait group
	isLogoPrinted   = false                     // Is the logo printed

	workMode       string
	filePathList   []string
	searchRequests []string
	saveType       string
	cleanType      string
	delimetr       = strings.TrimSpace(":")
)
