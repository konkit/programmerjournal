package database

import (
	"errors"
	"gorm.io/gorm"
	"programmerjournal-backend/model/entry"
	"programmerjournal-backend/model/tag"
	"regexp"
	"strings"
)

func saveTags(db *gorm.DB, entry *entry.Entry) error {
	tags, err := parseTags(db, entry)
	if err != nil {
		return err
	}

	entryTags, err := createEntryTags(db, entry, tags)
	if err != nil {
		return err
	}

	err = deleteOldEntryTags(db, entry, entryTags)
	if err != nil {
		return err
	}

	return nil
}

func deleteTagsFromEntry(db *gorm.DB, entryID uint) error {
	etArr := []tag.EntryTag{}
	err := db.Model(&tag.EntryTag{}).
		Where("entry_id = ?", entryID).
		Find(&etArr).
		Error

	if err != nil {
		return err
	}

	for _, et := range etArr {
		err = db.Delete(&et).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func containsTag(tags []tag.EntryTag, et tag.EntryTag) bool {
	for _, t := range tags {
		if t.ID == et.ID {
			return true
		}
	}
	return false
}

func parseTags(db *gorm.DB, e *entry.Entry) ([]tag.Tag, error) {
	hashtagRegex := regexp.MustCompile(`^#[a-zA-Z0-9_]+$`)

	tokens := strings.Split(e.Title, " ")

	var tagTokens []string
	for _, token := range tokens {
		if hashtagRegex.MatchString(token) {
			tagToAdd := token[1:]
			tagToAdd = strings.ToLower(tagToAdd)
			tagTokens = append(tagTokens, tagToAdd)
		}
	}

	var tags []tag.Tag
	for _, tagToken := range tagTokens {
		t := tag.Tag{}
		err := db.Model(tag.Tag{}).Where("name = ?", tagToken).First(&t).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				t = tag.Tag{
					Name: tagToken,
				}
				err = db.Save(&t).Error
				if err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		}

		tags = append(tags, t)
	}

	return tags, nil
}

func createEntryTags(db *gorm.DB, e *entry.Entry, tags []tag.Tag) ([]tag.EntryTag, error) {
	etArr := []tag.EntryTag{}

	for _, t := range tags {
		et := tag.EntryTag{}
		err := db.Model(&tag.EntryTag{}).
			Where("entry_id = ?", e.ID).
			Where("tag_id = ?", t.ID).
			First(&et).
			Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				et = tag.EntryTag{
					EntryID: e.ID,
					TagID:   t.ID,
				}
				err = db.Save(&et).Error
				if err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		}

		etArr = append(etArr, et)
	}

	return etArr, nil
}

func deleteOldEntryTags(db *gorm.DB, e *entry.Entry, tags []tag.EntryTag) error {
	etArr := []tag.EntryTag{}
	err := db.Model(&tag.EntryTag{}).
		Where("entry_id = ?", e.ID).
		Find(&etArr).
		Error

	if err != nil {
		return err
	}

	for _, et := range etArr {
		if !containsTag(tags, et) {
			err = db.Delete(&et).Error
			if err != nil {
				return err
			}
		}
	}

	return nil
}
