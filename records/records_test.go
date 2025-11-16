package records

import (
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

func TestSortPathRecords(t *testing.T) {
	pr := map[string]PathRecord{
		"path1": PathRecord{Timestamp: "1990-01-01 00:00:01 +08"},
		"path2": PathRecord{Timestamp: "1991-01-01 00:00:01 +08"},
		"path3": PathRecord{Timestamp: "2024-04-01 12:34:12 +08"},
		"path4": PathRecord{Timestamp: "2016-02-23 22:10:10 +08"},
	}

	records := SortRecords(pr)

	if records[0] != "path3" {
		t.Fatalf(`path3 is not the latest record`)
	}

	if records[1] != "path4" {
		t.Fatalf(`path4 is not the 2nd record`)
	}

	if records[2] != "path2" {
		t.Fatalf(`path2 is not the 3rd record`)
	}

	if records[3] != "path1" {
		t.Fatalf(`path1 is not the last record`)
	}
}

func TestSortStashRecords(t *testing.T) {
	pr := map[string]StashRecord{
		"path1": StashRecord{Timestamp: "1990-01-01 00:00:01 +08", Alias: "apple"},
		"path2": StashRecord{Timestamp: "1991-01-01 00:00:01 +08", Alias: "banana"},
		"path3": StashRecord{Timestamp: "2024-04-01 12:34:12 +08", Alias: "pear"},
		"path4": StashRecord{Timestamp: "2016-02-23 22:10:10 +08", Alias: "orange"},
	}

	records := SortRecords(pr)

	if records[0] != "path3" && pr[records[0]].Alias != "pear" {
		t.Fatalf(`path3 is not the latest record`)
	}

	if records[1] != "path4" && pr[records[1]].Alias != "orange" {
		t.Fatalf(`path4 is not the 2nd record`)
	}

	if records[2] != "path2" && pr[records[2]].Alias != "banana" {
		t.Fatalf(`path2 is not the 3rd record`)
	}

	if records[3] != "path1" && pr[records[3]].Alias != "apple" {
		t.Fatalf(`path1 is not the last record`)
	}
}

func TestAliasExists(t *testing.T) {
	r := Records{
		PathRecords: map[string]PathRecord{},
		StashRecords: map[string]StashRecord{
			"path1": StashRecord{Timestamp: "1991-01-01 01:01:01 +08", Alias: "ape"},
			"path2": StashRecord{Timestamp: "1991-01-01 01:01:01 +08", Alias: "bear"},
		},
	}

	if !r.AliasExists("ape") {
		t.Fatalf(`ape alias is not detected in StashRecords`)
	}

	if !r.AliasExists("bear") {
		t.Fatalf(`bear alias is not detected in StashRecords`)
	}

	if r.AliasExists("cat") {
		t.Fatalf(`cat alias is detected in StashRecords but does not exists`)
	}
}
