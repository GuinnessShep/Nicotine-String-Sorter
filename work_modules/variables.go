package work_modules

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"github.com/panjf2000/ants/v2"
	"github.com/schollz/progressbar/v3"
	"golang.org/x/text/encoding"
	"os"
	"regexp"
	"runtime"
	"sync"
	"time"
)

var (
	// Colors
	ColorBlue    = color.New(color.FgBlue).Add(color.Bold)
	ColorGreen   = color.New(color.FgGreen).Add(color.Bold)
	ColorRed     = color.New(color.FgRed).Add(color.Bold)
	ColorMagenta = color.New(color.FgMagenta).Add(color.Bold)
	ColorYellow  = color.New(color.FgYellow).Add(color.Bold)

	// General
	isFileInProcessing       bool                                                         // Is the file being processed
	isResultWrited           bool                                                         // Is the result written
	fileBadSymbolsPattern, _                          = regexp.Compile(`[^a-zA-Z0-9\.]+`) // Allowed symbols for the file
	checkedLines             int64                    = 0                                 // Number of processed lines
	checkedFiles                                      = 0                                 // Number of processed files
	currentFileSize          int64                    = 0                                 // Size of the current file in sort
	currentFileLines         int64                    = 0                                 // Number of lines in the current file
	fileDecoder              *encoding.Decoder                                            // File decoder
	cacheMutex               sync.Mutex                                                   // Cache method mutex to get the number of available lines
	cachedStrCount           int64                                                        // Cached available lines count
	lastUpdate               time.Time                                                    // Time since the last update of cachedStrCount
	pBar                     *progressbar.ProgressBar                                     // Progress bar
	runDir                   = GetRunDir()                                                // Run directory
	currPath                 string                                                       // Current file path
	poolerr                  error                                                        // Pool creation error
	workWG                   sync.WaitGroup                                               // Synchronization waitgroup
	TMPlinesLen              = 0                                                          // Chunk of lines in the file
	currPathCut              string                                                       // Current file without the full path

	// Sorter
	currFileMatchLines           int64                             = 0                      // Number of matching lines in the current file
	matchLines                   int64                             = 0                      // Number of matching lines
	reqLen                                                         = 0                      // Number of requests
	sorterDubles                 int64                             = 0                      // Number of duplicate lines
	requestStructMap                                               = make(map[string]*Work) // Map with the structure for each request
	sorterPool                   *ants.MultiPoolWithFunc                                    // Sorter pool
	sorterWriteChannelMap        = make(map[string]chan [2]string)                          // Map of channels
	sorterRequestStatMap         = make(map[string]int64)                                   // Number of found matches for each request
	sorterRequestStatMapCurrFile = make(map[string]int64)                                   // Number of found matches for each request in the current file
	sorterResultWriterMap        = make(map[string]*bufio.Writer)                           // Map of writers for each request
	sorterResultFileMap          = make(map[string]*os.File)                                // Map of files for each request
	sorterStringHashMap          = make(map[uint64]bool)                                    // Map of string hashes

	// Cleaner
	validPattern, _          = regexp.Compile(`^[a-zA-Z0-9\.\,\!\?\:\;\-\'\"\@\/\#\$\%\^\&\*\(\)\_\+\=\~\x60\|\[\]\{\}]{12,256}$`) // Validation pattern
	unknownPattern, _        = regexp.Compile(`(?i)UNKNOWN`)                                                                       // UNKNOWN content pattern
	partsPattern             *regexp.Regexp
	cleanerOutputFilesMap    = make(map[string]string)                              // Map of output files
	cleanerResultChannelMap  = make(map[string]chan string)                         // Map of valid strings
	cleanerWriteFile         *bufio.Writer                                          // Output file writer
	cleanerInvalidLen        int64                          = 0                     // Number of invalid lines
	currFileInvalidLen       int64                          = 0                     // Number of invalid lines in the current file
	cleanerDublesLen         int64                          = 0                     // Number of duplicate lines
	currFileDubles           int64                          = 0                     // Number of duplicate lines in the current file
	cleanerWritedString      int64                          = 0                     // Number of written strings
	currFileWritedString     int64                          = 0                     // Number of written strings in the current file
	cleanerStringHashMap                                    = make(map[uint64]bool) // Map of string hashes
	cleanerPartsPatternIsErr                                = false                 // Error in compiling parts pattern

	// Args
	filePathList   []string
	searchRequests []string
	saveType       string
	workMode       string
	cleanType      string
	delimetr       string
)

type Work struct {
	requestPattern *regexp.Regexp // Request regex
	resultFile     string         // Filename containing found strings
}

func InitVar(_workMode string, _filePathList []string, _searchRequests []string, _saveType string, _cleanType string, _delimetr string) {
	workMode = _workMode
	filePathList = _filePathList
	searchRequests = _searchRequests
	saveType = _saveType
	cleanType = _cleanType
	delimetr = regexp.QuoteMeta(_delimetr)
}

func InitSorter() {

	reqLen = len(searchRequests)

	// Initialize channels
	for _, path := range filePathList {
		sorterWriteChannelMap[path] = make(chan [2]string)
	}

	sorterPool, poolerr = ants.NewMultiPoolWithFunc(
		runtime.NumCPU(),
		100000,
		func(line interface{}) { SorterProcessLine(line.(string)) },
		ants.RoundRobin,
		ants.WithPreAlloc(true),
	)

	if poolerr != nil {
		PrintErr()
		ColorRed.Print("Cannot start sorter: Pool error: \n\n\n		", poolerr, "\n\n\n   Press Enter to exit")
		_, _ = fmt.Scanln()
		os.Exit(1)
	}
}
