package utils

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"gorm.io/gorm"
	"os"
	"programmerjournal-backend/database"
	"programmerjournal-backend/model/entry"
	"programmerjournal-backend/model/tag"
	"testing"
)

func TestSaveTags(t *testing.T) {
	dbTestPath, db := initDB()
	defer os.Remove(dbTestPath)

	initEntry := entry.Entry{
		Title: "Something something #tagone",
	}
	err := initializeEntry(db, initEntry)
	if err != nil {
		t.Fatalf("error saving initial entry: %v", err)
	}

	initTags := []tag.Tag{}
	err = initializeTags(db, initTags, err)
	if err != nil {
		t.Fatalf("error saving initial tags: %v", err)
	}

	initEntryTags := []tag.EntryTag{}
	err = initializeEntryTags(db, initEntryTags, err)
	if err != nil {
		t.Fatalf("error saving initial entryTags: %v", err)
	}

	err = SaveTags(db, initEntry)
	if err != nil {
		t.Fatalf("Error saving tags: %v", err)
	}

	tags := []tag.Tag{}
	db.Model(tag.Tag{}).Find(&tags)
	wantTags := []tag.Tag{
		{
			Name: "tagone",
		},
	}
	if diff := cmp.Diff(wantTags, tags, cmpopts.IgnoreFields(tag.Tag{}, "ID")); diff != "" {
		t.Errorf("Tags mismatch (-want +got):\n%s", diff)
	}

	entryTags := []tag.EntryTag{}
	db.Model(tag.EntryTag{}).Find(&entryTags)
	wantEntryTag := []tag.EntryTag{
		{
			EntryID: initEntry.ID,
			TagID:   tags[0].ID,
		},
	}
	if diff := cmp.Diff(wantEntryTag, entryTags, cmpopts.IgnoreFields(tag.EntryTag{}, "ID")); diff != "" {
		t.Errorf("EntryTags mismatch (-want +got):\n%s", diff)
	}
}

func TestSaveTagsTagExists(t *testing.T) {
	dbTestPath, db := initDB()
	defer os.Remove(dbTestPath)

	initEntry := entry.Entry{
		Title: "Something something #tagone",
	}
	err := initializeEntry(db, initEntry)
	if err != nil {
		t.Fatalf("error saving initial entry: %v", err)
	}

	initTags := []tag.Tag{
		{
			Name: "tagone",
		},
	}
	err = initializeTags(db, initTags, err)
	if err != nil {
		t.Fatalf("error saving initial tags: %v", err)
	}

	initEntryTags := []tag.EntryTag{
		{
			EntryID: initEntry.ID,
			TagID:   initTags[0].ID,
		},
	}
	err = initializeEntryTags(db, initEntryTags, err)
	if err != nil {
		t.Fatalf("error saving initial entryTags: %v", err)
	}

	err = SaveTags(db, initEntry)
	if err != nil {
		t.Fatalf("Error saving tags: %v", err)
	}

	tags := []tag.Tag{}
	db.Model(tag.Tag{}).Find(&tags)
	wantTags := []tag.Tag{
		{
			Name: "tagone",
		},
	}
	if diff := cmp.Diff(wantTags, tags, cmpopts.IgnoreFields(tag.Tag{}, "ID")); diff != "" {
		t.Errorf("Tags mismatch (-want +got):\n%s", diff)
	}

	entryTags := []tag.EntryTag{}
	db.Model(tag.EntryTag{}).Find(&entryTags)
	wantEntryTag := []tag.EntryTag{
		{
			EntryID: initEntry.ID,
			TagID:   tags[0].ID,
		},
	}
	if diff := cmp.Diff(wantEntryTag, entryTags, cmpopts.IgnoreFields(tag.EntryTag{}, "ID")); diff != "" {
		t.Errorf("EntryTags mismatch (-want +got):\n%s", diff)
	}
}

func TestSaveTagsOldTags(t *testing.T) {
	dbTestPath, db := initDB()
	defer os.Remove(dbTestPath)

	initEntry := entry.Entry{
		Title: "Something something #tagone",
	}
	err := initializeEntry(db, initEntry)
	if err != nil {
		t.Fatalf("error saving initial entry: %v", err)
	}

	initTags := []tag.Tag{
		{
			Name: "tagold",
		},
	}
	err = initializeTags(db, initTags, err)
	if err != nil {
		t.Fatalf("error saving initial tags: %v", err)
	}

	initEntryTags := []tag.EntryTag{
		{
			EntryID: initEntry.ID,
			TagID:   initTags[0].ID,
		},
	}
	err = initializeEntryTags(db, initEntryTags, err)
	if err != nil {
		t.Fatalf("error saving initial entryTags: %v", err)
	}

	err = SaveTags(db, initEntry)
	if err != nil {
		t.Fatalf("Error saving tags: %v", err)
	}

	tags := []tag.Tag{}
	db.Model(tag.Tag{}).Find(&tags)
	wantTags := []tag.Tag{
		{
			Name: "tagold",
		},
		{
			Name: "tagone",
		},
	}
	if diff := cmp.Diff(wantTags, tags, cmpopts.IgnoreFields(tag.Tag{}, "ID")); diff != "" {
		t.Errorf("Tags mismatch (-want +got):\n%s", diff)
	}

	entryTags := []tag.EntryTag{}
	db.Model(tag.EntryTag{}).Find(&entryTags)
	wantEntryTag := []tag.EntryTag{
		{
			EntryID: initEntry.ID,
			TagID:   getTagByName(tags, "tagone").ID,
		},
	}
	if diff := cmp.Diff(wantEntryTag, entryTags, cmpopts.IgnoreFields(tag.EntryTag{}, "ID")); diff != "" {
		t.Errorf("EntryTags mismatch (-want +got):\n%s", diff)
	}
}

func initDB() (string, *gorm.DB) {
	dbTestPath := "./test.db"
	db, _ := database.InitDB(dbTestPath)

	db.Exec("DELETE FROM entries")
	db.Exec("DELETE FROM tags")
	db.Exec("DELETE FROM entrytags")
	return dbTestPath, db
}

func initializeEntry(db *gorm.DB, initEntry entry.Entry) error {
	return db.Save(&initEntry).Error
}

func initializeTags(db *gorm.DB, initTags []tag.Tag, err error) error {
	for _, initTag := range initTags {
		err = db.Save(&initTag).Error
	}
	return err
}

func initializeEntryTags(db *gorm.DB, initEntryTags []tag.EntryTag, err error) error {
	for _, initET := range initEntryTags {
		err = db.Save(&initET).Error
	}
	return err
}

func getTagByName(tags []tag.Tag, s string) tag.Tag {
	for _, t := range tags {
		if t.Name == s {
			return t
		}
	}
	return tag.Tag{}
}
