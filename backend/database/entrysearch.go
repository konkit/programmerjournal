package database

import (
	"programmerjournal-backend/model/entry"
	"programmerjournal-backend/model/tag"
)

func (es *EntryService) FindEntriesByStatus(status entry.Status) ([]entry.Entry, error) {
	var entriesFromDB []entry.Entry
	err := es.db.Model(entry.Entry{}).
		Where("status = ?", status).
		Find(&entriesFromDB).
		Error

	return entriesFromDB, err
}

// SearchEntries matches the query case-insensitively against the title, description and task update.
func (es *EntryService) SearchEntries(query string) ([]entry.Entry, error) {
	var entriesFromDB []entry.Entry
	pattern := "%" + query + "%"
	err := es.db.Model(entry.Entry{}).
		Where("title LIKE ? OR description LIKE ? OR task_update LIKE ?", pattern, pattern, pattern).
		Find(&entriesFromDB).
		Error

	return entriesFromDB, err
}

func (es *EntryService) FindEntryIDsByTag(tagName string) (map[uint]bool, error) {
	var entryTags []tag.EntryTag
	err := es.db.Model(tag.EntryTag{}).
		Joins("JOIN tags ON tags.id = entry_tags.tag_id").
		Where("tags.name = ?", tagName).
		Find(&entryTags).
		Error
	if err != nil {
		return nil, err
	}

	ids := make(map[uint]bool, len(entryTags))
	for _, et := range entryTags {
		ids[et.EntryID] = true
	}
	return ids, nil
}

func (es *EntryService) ListTagNames() ([]string, error) {
	names := []string{}
	err := es.db.Model(tag.Tag{}).Order("name").Pluck("name", &names).Error
	return names, err
}
