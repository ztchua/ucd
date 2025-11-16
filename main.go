package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/ztcjoe93/ucd/configurations"
	"github.com/ztcjoe93/ucd/records"
	"github.com/ztcjoe93/ucd/utilities"
)

var (
	configs         configurations.Configuration
	aliasFlag       string
	helpFlag        bool
	clearFlag       bool
	clearStashFlag  bool
	dynamicSwapFlag int
	numRepeatFlag   int
	listFlag        bool
	listStashFlag   bool
	historyPathFlag int
	aliasPathFlag   string
	modifyAliasFlag int
	stashPathFlag   int
	stashFlag       bool
	versionFlag     bool
	cachePath       string
	cacheFile       *os.File

	invalidPath bool = false
)

const APPLICATION_NAME string = "ucd"
const APPLICATION_VERSION string = "1.0.0"

func main() {
	log.SetFlags(0)
	// flags
	flag.BoolVar(&helpFlag, "h", false, "display help")
	flag.StringVar(&aliasFlag, "a", "", "alias for stashed path, used in conjunction with -s")
	flag.BoolVar(&versionFlag, "v", false, "display ucd version")
	flag.BoolVar(&clearFlag, "c", false, "clear history list")
	flag.BoolVar(&clearStashFlag, "cs", false, "clear stash list")
	flag.IntVar(&dynamicSwapFlag, "d", 0, "swap out directory to arg after -d parent directories")
	flag.BoolVar(&listFlag, "l", false, "display Most Recently Used (MRU) list of paths chdir-ed into")
	flag.BoolVar(&listStashFlag, "ls", false, "display list of stashed cd commands")
	flag.IntVar(&modifyAliasFlag, "ma", 0, "modify alias of indicated # from the stash list")
	flag.IntVar(&numRepeatFlag, "n", 1, "no. of times to execute chdir")
	flag.IntVar(&historyPathFlag, "p", 0, "chdir to the indicated # from the MRU list")
	flag.IntVar(&stashPathFlag, "ps", 0, "chdir to the indicated # from the stash list")
	flag.StringVar(&aliasPathFlag, "pa", "", "chdir to path with matching alias from stash list")
	flag.BoolVar(&stashFlag, "s", false, "stash cd path into a separate list")
	flag.Parse()

	args := flag.Args()
	homeDir, _ := os.UserHomeDir()

	if helpFlag {
		flag.PrintDefaults()
		utilities.ReturnCwd()
	}

	if versionFlag {
		log.Printf("%v v%v\n", APPLICATION_NAME, APPLICATION_VERSION)
		utilities.ReturnCwd()
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
		utilities.ReturnCwd()
	}

	if clearStashFlag {
		r = records.Records{
			PathRecords:  r.PathRecords,
			StashRecords: map[string]records.StashRecord{},
		}
		output, _ := json.Marshal(r)
		os.WriteFile(cachePath, output, 0644)
		utilities.ReturnCwd()
	}

	// exit earlier depending on flag passed in
	if listFlag {
		r.ListRecords("path", configs.MaxMRUDisplay)
		utilities.ReturnCwd()
	}

	if listStashFlag {
		r.ListRecords("stash", configs.MaxMRUDisplay)
		utilities.ReturnCwd()
	}

	if len(args) > 1 {
		log.Printf("Only < 1 arguments, found %v args can be passed to ucd\n", len(args))
		utilities.ReturnCwd()
	}

	if modifyAliasFlag > 0 {
		srk := records.SortRecords(r.StashRecords)
		sr := r.StashRecords[srk[modifyAliasFlag-1]]
		sr.Alias = args[0]
		r.StashRecords[srk[modifyAliasFlag-1]] = sr

		output, _ := json.Marshal(r)
		os.WriteFile(cachePath, output, 0644)

		r.ListRecords("stash", configs.MaxMRUDisplay)
		utilities.ReturnCwd()
	}

	// fmt.Print sends output to stdout, this will be consumed by builtin `cd` command

	var targetPath string

	if dynamicSwapFlag > 0 {
		targetPath = utilities.DynamicPathSwap(args[0], dynamicSwapFlag)
	} else if aliasPathFlag != "" {
		found := false
		for key, rec := range r.StashRecords {
			if rec.Alias == aliasPathFlag {
				targetPath = key
				found = true
				break
			}
		}

		if !found {
			log.Printf("unable to cd -- alias ``%v` not found\n", aliasPathFlag)
			utilities.ReturnCwd()
		}
	} else if historyPathFlag > 0 {
		mruRecords := records.SortRecords(r.PathRecords)
		if historyPathFlag-1 > len(mruRecords)-1 {
			log.Printf("invalid #, there are only %v records\n", len(mruRecords))
			utilities.ReturnCwd()
		}
		targetPath = mruRecords[historyPathFlag-1]
	} else if stashPathFlag > 0 {
		stashRecords := records.SortRecords(r.StashRecords)
		targetPath = stashRecords[stashPathFlag-1]
	} else {
		if len(args) > 0 {
			targetPath = utilities.Repeat(args[0], numRepeatFlag)
		} else {
			targetPath = homeDir
		}
	}

	if targetPath == "-" {
		fmt.Print("-")
		os.Exit(0)
	}

	if utilities.IsInvalidPath(targetPath) {
		if configs.FileFallbackBehavior {
			targetPath = utilities.GetParentDir(targetPath)
		} else {
			utilities.ReturnCwd()
		}
	} else {
		targetPath, _ = os.Getwd()
	}

	rec, ok := r.PathRecords[targetPath]
	if ok {
		rec.Timestamp = utilities.TimeNow()
		r.PathRecords[targetPath] = rec
	} else {
		r.PathRecords[targetPath] = records.PathRecord{Timestamp: utilities.TimeNow()}
	}

	if stashFlag {
		if r.AliasExists(aliasFlag) {
			log.Printf("Alias `%v` already exists\n", aliasFlag)
			utilities.ReturnCwd()
		}
		r.StashRecords[targetPath] = records.StashRecord{Alias: aliasFlag, Timestamp: utilities.TimeNow()}
	}

	utilities.AutoClear(&r, configs.MaxMRUDisplay)
	strings.Replace(targetPath, " ", "\\ ", -1)
	fmt.Print(targetPath)

	output, _ := json.Marshal(r)
	os.WriteFile(cachePath, output, 0644)
}
