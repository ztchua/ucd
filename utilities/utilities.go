package utilities

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ztcjoe93/ucd/records"
)

func ReturnCwd() {
	fmt.Print(".")
	os.Exit(0)
}

func DynamicPathSwap(swapArg string, upCount int) string {
	paths := make([]string, 0)
	for i := 0; i < upCount; i++ {
		wd, _ := os.Getwd()
		wdArr := strings.Split(wd, "/")
		paths = prependStrSlice(paths, wdArr[len(wdArr)-1])
		os.Chdir("..")
	}

	os.Chdir("..")
	paths = prependStrSlice(paths, swapArg)
	targetPath := strings.Join(paths, "/")

	return targetPath
}

func prependStrSlice(x []string, y string) []string {
	x = append(x, "")
	copy(x[1:], x)
	x[0] = y
	return x
}

func IsInvalidPath(targetPath string) bool {
	err := os.Chdir(targetPath)
	if err != nil {
		Clog(fmt.Sprintf("%v is not a valid path", targetPath))
		return true
	}

	return false
}

func GetParentDir(targetPath string) string {
	parentPath := filepath.Dir(targetPath)
	err := os.Chdir(parentPath)
	if err != nil {
		Clog(fmt.Sprintf("%v is not a valid path", targetPath))
		os.Exit(0)
	}

	// fallback to current working directory (which is parent directory of target file)
	parentPath, _ = os.Getwd()
	Clog(fmt.Sprintf("Falling back to parent directory %v", parentPath))
	return parentPath
}

func Repeat(str string, times int) string {
	s := make([]string, times)
	for i := range s {
		s[i] = str
	}

	return strings.Join(s, "/")
}

func TimeNow() string {
	return time.Now().Format("2006-01-02 15:04:05 MST")
}

func AutoClear(r *records.Records, limit int) {
	rk := records.SortRecords(r.PathRecords)

	if limit != -1 {
		return
	}

	if len(rk) > limit {
		for i := limit; i < len(rk); i++ {
			delete(r.PathRecords, rk[i])
		}
	}
}

func Clog(msg string) {
	log.Printf("[ucd] %v\n", msg)
}
