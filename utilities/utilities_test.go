package utilities

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ztcjoe93/ucd/records"
)

func TestRepeatFn(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		times    int
		expected string
	}{
		{
			name:     "repeat double dot 4 times",
			str:      "..",
			times:    4,
			expected: "../../../..",
		},
		{
			name:     "repeat dot 2 times",
			str:      ".",
			times:    2,
			expected: "./.",
		},
		{
			name:     "repeat zero times",
			str:      "test",
			times:    0,
			expected: "",
		},
		{
			name:     "repeat single time",
			str:      "abc",
			times:    1,
			expected: "abc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Repeat(tt.str, tt.times)
			if result != tt.expected {
				t.Fatalf(`Repeat(%q, %d) = %q, want %q`, tt.str, tt.times, result, tt.expected)
			}
		})
	}
}

func TestIsInvalidPath(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "invalid random path",
			path:     "doanfkjzx/sdfj931/sdfkjal",
			expected: true,
		},
		{
			name:     "empty path",
			path:     "",
			expected: true,
		},
		{
			name:     "nonexistent deep path",
			path:     "/nonexistent/path/that/does/not/exist",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsInvalidPath(tt.path)
			if result != tt.expected {
				t.Fatalf(`IsInvalidPath(%q) = %v, want %v`, tt.path, result, tt.expected)
			}
		})
	}
}

func TestPrependStrSlice(t *testing.T) {
	tests := []struct {
		name     string
		slice    []string
		prepend  string
		expected []string
	}{
		{
			name:     "prepend to existing slice",
			slice:    []string{"apple", "banana"},
			prepend:  "cookie",
			expected: []string{"cookie", "apple", "banana"},
		},
		{
			name:     "prepend to empty slice",
			slice:    []string{},
			prepend:  "first",
			expected: []string{"first"},
		},
		{
			name:     "prepend to single element",
			slice:    []string{"only"},
			prepend:  "new",
			expected: []string{"new", "only"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := prependStrSlice(tt.slice, tt.prepend)
			if len(result) != len(tt.expected) {
				t.Fatalf(`slice length = %d, want %d`, len(result), len(tt.expected))
			}
			for i, v := range tt.expected {
				if result[i] != v {
					t.Fatalf(`result[%d] = %q, want %q`, i, result[i], v)
				}
			}
		})
	}
}

func TestTimeNow(t *testing.T) {
	result := TimeNow()

	if len(result) == 0 {
		t.Fatalf(`TimeNow() returned empty string`)
	}

	if !strings.Contains(result, "-") || !strings.Contains(result, ":") {
		t.Fatalf(`TimeNow() = %q, expected format "2006-01-02 15:04:05 MST"`, result)
	}

	parts := strings.Split(result, " ")
	if len(parts) < 3 {
		t.Fatalf(`TimeNow() = %q, expected at least 3 space-separated parts`, result)
	}
}

func TestAutoClear(t *testing.T) {
	t.Run("function returns early when limit is not -1", func(t *testing.T) {
		r := &records.Records{
			PathRecords:  make(map[string]records.PathRecord),
			StashRecords: make(map[string]records.StashRecord),
		}

		for i := 0; i < 5; i++ {
			path := "/path/" + string(rune('a'+i))
			r.PathRecords[path] = records.PathRecord{Timestamp: "2024-01-01 00:00:00 +08"}
		}

		AutoClear(r, 3)

		if len(r.PathRecords) != 5 {
			t.Fatalf(`len(PathRecords) = %d, want 5 (function does nothing when limit != -1)`, len(r.PathRecords))
		}
	})

	t.Run("empty records", func(t *testing.T) {
		r := &records.Records{
			PathRecords:  make(map[string]records.PathRecord),
			StashRecords: make(map[string]records.StashRecord),
		}

		AutoClear(r, 5)

		if len(r.PathRecords) != 0 {
			t.Fatalf(`len(PathRecords) = %d, want 0`, len(r.PathRecords))
		}
	})
}

func TestGetParentDir(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() (string, func())
		expected string
	}{
		{
			name: "get parent of nested directory",
			setup: func() (string, func()) {
				dir, _ := os.MkdirTemp("", "test")
				subdir := filepath.Join(dir, "subdir")
				os.Mkdir(subdir, 0755)
				return subdir, func() { os.RemoveAll(dir) }
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			targetPath, cleanup := tt.setup()
			defer cleanup()

			parentPath := filepath.Dir(targetPath)

			originalWd, _ := os.Getwd()
			result := GetParentDir(targetPath)
			os.Chdir(originalWd)

			if !strings.Contains(result, parentPath) {
				t.Fatalf(`GetParentDir(%q) = %q, expected to contain %q`, targetPath, result, parentPath)
			}
		})
	}
}

func TestClog(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(nil)

	Clog("test message")
	output := buf.String()

	if !strings.Contains(output, "[ucd]") {
		t.Fatalf(`Clog output missing [ucd] prefix, got: %s`, output)
	}

	if !strings.Contains(output, "test message") {
		t.Fatalf(`Clog output missing message, got: %s`, output)
	}
}

func TestDynamicPathSwap(t *testing.T) {
	originalWd, _ := os.Getwd()

	dir, _ := os.MkdirTemp("", "ucd-test")
	subdir1 := filepath.Join(dir, "level1")
	subdir2 := filepath.Join(subdir1, "level2")
	os.MkdirAll(subdir2, 0755)

	os.Chdir(subdir2)

	result := DynamicPathSwap("newdir", 1)

	os.Chdir(originalWd)
	os.RemoveAll(dir)

	if !strings.Contains(result, "newdir") {
		t.Fatalf(`DynamicPathSwap("newdir", 1) = %q, expected to contain "newdir"`, result)
	}

	if !strings.Contains(result, "level2") {
		t.Fatalf(`DynamicPathSwap("newdir", 1) = %q, expected to contain "level2"`, result)
	}
}
