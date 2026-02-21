package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ztcjoe93/ucd/configurations"
	"github.com/ztcjoe93/ucd/records"
	util "github.com/ztcjoe93/ucd/utilities"
)

var (
	configs         configurations.Configuration
	aliasFlag       string
	helpFlag        bool
	clearFlag       bool
	clearStashFlag  bool
	listFlag        bool
	listStashFlag   bool
	historyPathFlag int
	aliasPathFlag   string
	stashFlag       bool
	stashPathFlag   int
	versionFlag     bool
	cachePath       string
	cacheFile       *os.File
)

const APPLICATION_NAME string = "ucd"
const APPLICATION_VERSION string = "1.0.0"

func main() {
	log.SetFlags(0)
	// flags
	flag.BoolVar(&helpFlag, "h", false, "display help")
	flag.BoolVar(&versionFlag, "v", false, "display ucd version")

	flag.BoolVar(&listFlag, "l", false, "List history of cd paths")
	flag.BoolVar(&listStashFlag, "ls", false, "List history of stashed aliases")

	flag.BoolVar(&clearFlag, "c", false, "Clear history list")
	flag.BoolVar(&clearStashFlag, "cs", false, "Clear stash history list")

	flag.StringVar(&aliasFlag, "a", "", "alias for stashed path, used in conjunction with -s")

	flag.IntVar(&historyPathFlag, "p", 0, "Chdir to the # path from the list")
	flag.StringVar(&aliasPathFlag, "pa", "", "Chdir to path with the provided alias from stash list")

	flag.BoolVar(&stashFlag, "s", false, "stash cd path into a separate list")
	flag.IntVar(&stashPathFlag, "ps", 0, "Chdir to the # path from stash list")
	flag.Parse()

	args := flag.Args()
	homeDir, _ := os.UserHomeDir()

	if helpFlag {
		flag.PrintDefaults()
		util.ReturnCwd()
	}

	if versionFlag {
		log.Printf("%v v%v\n", APPLICATION_NAME, APPLICATION_VERSION)
		util.ReturnCwd()
	}

	configs = configs.GetConfigurations()

	cachePath = homeDir + "/.ucd-cache"
	cacheFile, _ := os.Open(cachePath)
	defer cacheFile.Close()
	byteValue, _ := io.ReadAll(cacheFile)

	var r records.Records
	err := json.Unmarshal(byteValue, &r)
	if err != nil {
		r = records.Records{
			PathRecords:  map[string]records.PathRecord{},
			StashRecords: map[string]records.StashRecord{},
		}
	}

	if configs.MaxMRUDisplay < 0 {
		configs.MaxMRUDisplay = len(r.PathRecords)
	}

	if clearFlag {
		r = records.Records{
			PathRecords:  map[string]records.PathRecord{},
			StashRecords: r.StashRecords,
		}
		output, _ := json.Marshal(r)
		os.WriteFile(cachePath, output, 0644)
		util.ReturnCwd()
	}

	if clearStashFlag {
		r = records.Records{
			PathRecords:  r.PathRecords,
			StashRecords: map[string]records.StashRecord{},
		}
		output, _ := json.Marshal(r)
		os.WriteFile(cachePath, output, 0644)
		util.ReturnCwd()
	}

	// exit earlier depending on flag passed in
	if listFlag {
		r.ListRecords("path", configs.MaxMRUDisplay)
		util.ReturnCwd()
	}

	if listStashFlag {
		r.ListRecords("stash", configs.MaxMRUDisplay)
		util.ReturnCwd()
	}

	if len(args) > 1 {
		util.Clog("Only 1 argument is accepted")
		util.ReturnCwd()
	}

	// fmt.Print sends output to stdout, this will be consumed by builtin `cd` command

	var targetPath string

	if aliasPathFlag != "" {
		found := false
		for key, rec := range r.StashRecords {
			if rec.Alias == aliasPathFlag {
				targetPath = key
				found = true
				break
			}
		}

		if !found {
			util.Clog(fmt.Sprintf("Alias %v not found", aliasPathFlag))
			util.ReturnCwd()
		}
	} else if historyPathFlag > 0 {
		mruRecords := records.SortRecords(r.PathRecords)
		if historyPathFlag-1 > len(mruRecords)-1 {
			util.Clog(fmt.Sprintf("Invalid path # provided - there are %v records", len(mruRecords)))
			util.ReturnCwd()
		}
		targetPath = mruRecords[historyPathFlag-1]
	} else if stashPathFlag > 0 {
		stashRecords := records.SortRecords(r.StashRecords)
		targetPath = stashRecords[stashPathFlag-1]
	} else {
		targetPath = homeDir
		if len(args) > 0 {
			var err error
			targetPath, err = filepath.Abs(args[0])
			if err != nil {
				util.Clog(fmt.Sprintf("Error resolving path: %v", err))
				util.ReturnCwd()
			}
		}
	}

	util.Clog(fmt.Sprintf("chdir to %v", targetPath))

	if targetPath == "-" {
		fmt.Print("-")
		os.Exit(0)
	}

	if util.IsInvalidPath(targetPath) {
		if configs.FileFallbackBehavior {
			targetPath = util.GetParentDir(targetPath)
		} else {
			util.ReturnCwd()
		}
	}

	rec, ok := r.PathRecords[targetPath]
	if ok {
		rec.Timestamp = util.TimeNow()
		r.PathRecords[targetPath] = rec
	} else {
		r.PathRecords[targetPath] = records.PathRecord{Timestamp: util.TimeNow()}
	}

	if stashFlag {
		if r.AliasExists(aliasFlag) {
			util.Clog(fmt.Sprintf("Alias `%v` already exists\n", aliasFlag))
			util.ReturnCwd()
		} else {
			log.Printf("Enter an alias for %v:", targetPath)

			_, err := fmt.Scanln(&aliasFlag)
			if err != nil {
				log.Fatalf("%v\n", err)
			}

			log.Printf("Stashed %v as `%v`\n", targetPath, aliasFlag)

			if aliasFlag != "" {
				r.StashRecords[targetPath] = records.StashRecord{Alias: aliasFlag, Timestamp: util.TimeNow()}
			}
		}
		r.StashRecords[targetPath] = records.StashRecord{Alias: aliasFlag, Timestamp: util.TimeNow()}
	}

	util.AutoClear(&r, configs.MaxMRUDisplay)
	targetPath = strings.Replace(targetPath, " ", "\\ ", -1)
	fmt.Print(targetPath)

	output, _ := json.Marshal(r)
	os.WriteFile(cachePath, output, 0644)
}
