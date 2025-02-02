package tag

type Tag struct {
	ID   uint
	Name string
}

type EntryTag struct {
	ID      uint `gorm:"primary_key"`
	EntryID uint `json:"entryId"`
	TagID   uint `json:"tagId"`
}
