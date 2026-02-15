package records

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

func TestRecordHasCorrectCount(t *testing.T) {
	var records []interface{}

	records = append(records, PathRecord{})
	records = append(records, StashRecord{})

	if !records[0].(PathRecord).HasCount() {
		t.Fatalf(`PathRecord does not return true when checking for count`)
	}

	if records[1].(StashRecord).HasCount() {
		t.Fatalf(`StashRecord returns true when checking for count`)
	}
}

func TestGetTimestamp(t *testing.T) {
	tests := []struct {
		name   string
		record interface {
			GetTimestamp() string
		}
		expected string
	}{
		{
			name:     "PathRecord returns timestamp",
			record:   PathRecord{Timestamp: "2024-01-01 12:00:00 MST"},
			expected: "2024-01-01 12:00:00 MST",
		},
		{
			name:     "StashRecord returns timestamp",
			record:   StashRecord{Timestamp: "2024-06-15 08:30:00 MST", Alias: "test"},
			expected: "2024-06-15 08:30:00 MST",
		},
		{
			name:     "Empty timestamp returns empty string",
			record:   PathRecord{Timestamp: ""},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.record.GetTimestamp()
			if result != tt.expected {
				t.Fatalf(`GetTimestamp() = %q, want %q`, result, tt.expected)
			}
		})
	}
}

func TestSortRecords(t *testing.T) {
	t.Run("SortPathRecords", func(t *testing.T) {
		pr := map[string]PathRecord{
			"path1": {Timestamp: "1990-01-01 00:00:01 +08"},
			"path2": {Timestamp: "1991-01-01 00:00:01 +08"},
			"path3": {Timestamp: "2024-04-01 12:34:12 +08"},
			"path4": {Timestamp: "2016-02-23 22:10:10 +08"},
		}

		records := SortRecords(pr)

		if records[0] != "path3" {
			t.Fatalf(`path3 is not the latest record, got %s`, records[0])
		}

		if records[1] != "path4" {
			t.Fatalf(`path4 is not the 2nd record, got %s`, records[1])
		}

		if records[2] != "path2" {
			t.Fatalf(`path2 is not the 3rd record, got %s`, records[2])
		}

		if records[3] != "path1" {
			t.Fatalf(`path1 is not the last record, got %s`, records[3])
		}
	})

	t.Run("SortStashRecords", func(t *testing.T) {
		sr := map[string]StashRecord{
			"path1": {Timestamp: "1990-01-01 00:00:01 +08", Alias: "apple"},
			"path2": {Timestamp: "1991-01-01 00:00:01 +08", Alias: "banana"},
			"path3": {Timestamp: "2024-04-01 12:34:12 +08", Alias: "pear"},
			"path4": {Timestamp: "2016-02-23 22:10:10 +08", Alias: "orange"},
		}

		records := SortRecords(sr)

		if records[0] != "path3" || sr[records[0]].Alias != "pear" {
			t.Fatalf(`path3 with alias pear is not the latest record, got %s with alias %s`, records[0], sr[records[0]].Alias)
		}

		if records[1] != "path4" || sr[records[1]].Alias != "orange" {
			t.Fatalf(`path4 with alias orange is not the 2nd record, got %s with alias %s`, records[1], sr[records[1]].Alias)
		}

		if records[2] != "path2" || sr[records[2]].Alias != "banana" {
			t.Fatalf(`path2 with alias banana is not the 3rd record, got %s with alias %s`, records[2], sr[records[2]].Alias)
		}

		if records[3] != "path1" || sr[records[3]].Alias != "apple" {
			t.Fatalf(`path1 with alias apple is not the last record, got %s with alias %s`, records[3], sr[records[3]].Alias)
		}
	})

	t.Run("EmptyRecords", func(t *testing.T) {
		pr := map[string]PathRecord{}
		records := SortRecords(pr)

		if len(records) != 0 {
			t.Fatalf(`expected empty slice, got %d records`, len(records))
		}
	})

	t.Run("SingleRecord", func(t *testing.T) {
		pr := map[string]PathRecord{
			"only": {Timestamp: "2024-01-01 00:00:00 +08"},
		}
		records := SortRecords(pr)

		if len(records) != 1 || records[0] != "only" {
			t.Fatalf(`expected single record "only", got %v`, records)
		}
	})
}

func TestAliasExists(t *testing.T) {
	tests := []struct {
		name     string
		records  Records
		alias    string
		expected bool
	}{
		{
			name: "alias exists",
			records: Records{
				PathRecords: map[string]PathRecord{},
				StashRecords: map[string]StashRecord{
					"path1": {Timestamp: "1991-01-01 01:01:01 +08", Alias: "ape"},
					"path2": {Timestamp: "1991-01-01 01:01:01 +08", Alias: "bear"},
				},
			},
			alias:    "ape",
			expected: true,
		},
		{
			name: "alias does not exist",
			records: Records{
				PathRecords: map[string]PathRecord{},
				StashRecords: map[string]StashRecord{
					"path1": {Timestamp: "1991-01-01 01:01:01 +08", Alias: "ape"},
				},
			},
			alias:    "cat",
			expected: false,
		},
		{
			name: "empty stash records",
			records: Records{
				PathRecords:  map[string]PathRecord{},
				StashRecords: map[string]StashRecord{},
			},
			alias:    "any",
			expected: false,
		},
		{
			name: "case sensitive match",
			records: Records{
				PathRecords: map[string]PathRecord{},
				StashRecords: map[string]StashRecord{
					"path1": {Timestamp: "1991-01-01 01:01:01 +08", Alias: "Ape"},
				},
			},
			alias:    "ape",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.records.AliasExists(tt.alias)
			if result != tt.expected {
				t.Fatalf(`AliasExists(%q) = %v, want %v`, tt.alias, result, tt.expected)
			}
		})
	}
}

func TestListRecords(t *testing.T) {
	t.Run("ListPathRecords", func(t *testing.T) {
		var buf bytes.Buffer
		log.SetOutput(&buf)
		defer log.SetOutput(nil)

		r := Records{
			PathRecords: map[string]PathRecord{
				"/path/a": {Timestamp: "2024-01-01 12:00:00 +08"},
				"/path/b": {Timestamp: "2024-01-02 12:00:00 +08"},
			},
			StashRecords: map[string]StashRecord{},
		}

		r.ListRecords("path", 10)
		output := buf.String()

		if !strings.Contains(output, "#") || !strings.Contains(output, "PATH") || !strings.Contains(output, "TIMESTAMP") {
			t.Fatalf(`ListRecords("path") output missing expected headers, got: %s`, output)
		}
	})

	t.Run("ListStashRecords", func(t *testing.T) {
		var buf bytes.Buffer
		log.SetOutput(&buf)
		defer log.SetOutput(nil)

		r := Records{
			PathRecords: map[string]PathRecord{},
			StashRecords: map[string]StashRecord{
				"/stash/a": {Timestamp: "2024-01-01 12:00:00 +08", Alias: "work"},
			},
		}

		r.ListRecords("stash", 10)
		output := buf.String()

		if !strings.Contains(output, "#") || !strings.Contains(output, "ALIAS") || !strings.Contains(output, "PATH") {
			t.Fatalf(`ListRecords("stash") output missing expected headers, got: %s`, output)
		}
	})

	t.Run("ListRecordsWithLimit", func(t *testing.T) {
		var buf bytes.Buffer
		log.SetOutput(&buf)
		defer log.SetOutput(nil)

		r := Records{
			PathRecords: map[string]PathRecord{
				"/path/a": {Timestamp: "2024-01-01 12:00:00 +08"},
				"/path/b": {Timestamp: "2024-01-02 12:00:00 +08"},
				"/path/c": {Timestamp: "2024-01-03 12:00:00 +08"},
			},
			StashRecords: map[string]StashRecord{},
		}

		r.ListRecords("path", 2)
		output := buf.String()

		lineCount := strings.Count(output, "\n")
		headerAndSeparator := 2
		dataLines := 2
		expectedLines := headerAndSeparator + dataLines

		if lineCount < dataLines {
			t.Fatalf(`ListRecords with limit 2 should show at most 2 data rows, got %d lines`, lineCount)
		}
		if lineCount > expectedLines+2 {
			t.Fatalf(`ListRecords with limit 2 showing more than expected, got %d lines`, lineCount)
		}
	})

	t.Run("ListEmptyRecords", func(t *testing.T) {
		var buf bytes.Buffer
		log.SetOutput(&buf)
		defer log.SetOutput(nil)

		r := Records{
			PathRecords:  map[string]PathRecord{},
			StashRecords: map[string]StashRecord{},
		}

		r.ListRecords("path", 10)
		output := buf.String()

		if !strings.Contains(output, "#") {
			t.Fatalf(`ListRecords on empty records should still show header, got: %s`, output)
		}
	})
}
