package database

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"gorm.io/gorm"
	"os"
	"programmerjournal-backend/model/entry"
	"programmerjournal-backend/model/tag"
	"testing"
)

func TestSaveTags(t *testing.T) {
	testCases := []struct {
		name           string
		initEntryTitle string
		initTags       []string
		initTagEntries []string
		wantTags       []string
		wantTagEntries []string
	}{
		{
			name:           "new_entry",
			initEntryTitle: "Something something #tagone",
			initTags:       []string{},
			initTagEntries: []string{},
			wantTags:       []string{"tagone"},
			wantTagEntries: []string{"tagone"},
		},
		{
			name:           "new_entry_with_capitals",
			initEntryTitle: "Something something #tagOne",
			initTags:       []string{},
			initTagEntries: []string{},
			wantTags:       []string{"tagone"},
			wantTagEntries: []string{"tagone"},
		},
		{
			name:           "new_entry_two_tags",
			initEntryTitle: "Something something #tagone #tagtwo",
			initTags:       []string{},
			initTagEntries: []string{},
			wantTags:       []string{"tagone", "tagtwo"},
			wantTagEntries: []string{"tagone", "tagtwo"},
		},
		{
			name:           "existing_entry",
			initEntryTitle: "Something something #tagone",
			initTags:       []string{"tagone"},
			initTagEntries: []string{"tagone"},
			wantTags:       []string{"tagone"},
			wantTagEntries: []string{"tagone"},
		},
		{
			name:           "existing_entry_with_capitals",
			initEntryTitle: "Something something #tagOne",
			initTags:       []string{"tagone"},
			initTagEntries: []string{"tagone"},
			wantTags:       []string{"tagone"},
			wantTagEntries: []string{"tagone"},
		},
		{
			name:           "old_tags_are_removed",
			initEntryTitle: "Something something #tagone #tagtwo",
			initTags:       []string{"tagold", "tagtwo"},
			initTagEntries: []string{"tagold", "tagtwo"},
			wantTags:       []string{"tagold", "tagtwo", "tagone"},
			wantTagEntries: []string{"tagone", "tagtwo"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dbTestPath, db := initDB()
			defer os.Remove(dbTestPath)
			e, err := initializeDBObjects(t, db, tc.initEntryTitle, tc.initTags, tc.initTagEntries)

			err = saveTags(db, e)
			if err != nil {
				t.Fatalf("Error saving tags: %v", err)
			}

			tags := listTags(db)
			err = verifyTags(tags, tc.wantTags)
			if err != nil {
				t.Fatalf("Error verifying tags: %v", err)
			}
			err = verifyTagEntries(db, tags, e, tc.wantTagEntries)
			if err != nil {
				t.Fatalf("Error verifying tags: %v", err)
			}
		})
	}
}

func TestDeleteTagsFromEntry(t *testing.T) {
	testCases := []struct {
		name           string
		initEntryTitle string
		initTags       []string
		initTagEntries []string
		wantTags       []string
		wantTagEntries []string
	}{
		{
			name:           "new_entry_two_tags",
			initEntryTitle: "Something something #tagone #tagtwo",
			initTags:       []string{"tagone", "tagtwo"},
			initTagEntries: []string{"tagone", "tagtwo"},
			wantTags:       []string{"tagone", "tagtwo"},
			wantTagEntries: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dbTestPath, db := initDB()
			defer os.Remove(dbTestPath)
			e, err := initializeDBObjects(t, db, tc.initEntryTitle, tc.initTags, tc.initTagEntries)

			err = deleteTagsFromEntry(db, e.ID)
			if err != nil {
				t.Fatalf("Error saving tags: %v", err)
			}

			tags := listTags(db)
			err = verifyTags(tags, tc.wantTags)
			if err != nil {
				t.Fatalf("Error verifying tags: %v", err)
			}
			err = verifyTagEntries(db, tags, e, tc.wantTagEntries)
			if err != nil {
				t.Fatalf("Error verifying tags: %v", err)
			}
		})
	}
}

func initDB() (string, *gorm.DB) {
	dbTestPath := "./test.db"
	db, _ := InitDB(dbTestPath)

	db.Exec("DELETE FROM entries")
	db.Exec("DELETE FROM tags")
	db.Exec("DELETE FROM entrytags")
	return dbTestPath, db
}

func initializeDBObjects(t *testing.T, db *gorm.DB, initEntryTitle string, initTags []string, initTagEntries []string) (*entry.Entry, error) {
	// Initialize EntryService
	e := entry.Entry{
		Title: initEntryTitle,
	}
	err := db.Save(&e).Error
	if err != nil {
		return nil, fmt.Errorf("Error saving entry: %v", err)
	}

	// Initialize Tags
	var initT []tag.Tag
	for _, t := range initTags {
		initT = append(initT, tag.Tag{Name: t})
	}
	for _, initTag := range initT {
		err = db.Save(&initTag).Error
	}
	if err != nil {
		t.Fatalf("error saving initial tags: %v", err)
	}

	// Initialize EntryTags
	var initET []tag.EntryTag
	for _, teStr := range initTagEntries {
		newET := tag.EntryTag{
			EntryID: e.ID,
			TagID:   getTagByName(initT, teStr).ID,
		}
		initET = append(initET, newET)
	}
	for _, et := range initET {
		err = db.Save(&et).Error
	}
	if err != nil {
		t.Fatalf("error saving initial entryTags: %v", err)
	}
	return &e, err
}

func listTags(db *gorm.DB) []tag.Tag {
	var tags []tag.Tag
	db.Model(tag.Tag{}).Find(&tags)
	return tags
}

func verifyTags(tags []tag.Tag, wantTagsStrs []string) error {
	wantTags := []tag.Tag{}
	for _, tStr := range wantTagsStrs {
		wt := tag.Tag{Name: tStr}
		wantTags = append(wantTags, wt)
	}
	if diff := cmp.Diff(wantTags, tags, cmpopts.IgnoreFields(tag.Tag{}, "ID")); diff != "" {
		return fmt.Errorf("tags mismatch (-want +got):\n%s", diff)
	}
	return nil
}

func verifyTagEntries(db *gorm.DB, tags []tag.Tag, e *entry.Entry, wantTagEntriesStrs []string) error {
	entryTags := []tag.EntryTag{}
	db.Model(tag.EntryTag{}).Find(&entryTags)
	wantEntryTag := []tag.EntryTag{}
	for _, teStr := range wantTagEntriesStrs {
		newET := tag.EntryTag{
			EntryID: e.ID,
			TagID:   getTagByName(tags, teStr).ID,
		}
		wantEntryTag = append(wantEntryTag, newET)
	}
	if diff := cmp.Diff(wantEntryTag, entryTags, cmpopts.IgnoreFields(tag.EntryTag{}, "ID")); diff != "" {
		return fmt.Errorf("EntryTags mismatch (-want +got):\n%s", diff)
	}

	return nil
}

func getTagByName(tags []tag.Tag, s string) tag.Tag {
	for _, t := range tags {
		if t.Name == s {
			return t
		}
	}
	return tag.Tag{}
}
