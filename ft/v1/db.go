package ft

import (
	"github.com/cockroachdb/errors"
	"gorm.io/gorm"
)

type Tagger struct {
	db *gorm.DB
}

func NewTagger(db *gorm.DB) Tagger {
	return Tagger{
		db: db,
	}
}

func (t *Tagger) Migrate() error {
	err := t.db.AutoMigrate(&Tag{}, &FileEntry{})
	return errors.Wrap(err, "db migrate")
}

func (t *Tagger) ensureFileEntry(filePath string) (*FileEntry, error) {
	var entry FileEntry
	err := t.db.First(&entry, "path = ?", filePath).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.Wrap(err, "find file entry")
		}

		entry.Path = filePath
		err = t.db.Create(&entry).Error
		if err != nil {
			return nil, errors.Wrap(err, "create file entry")
		}
	}
	return &entry, nil
}

func (t *Tagger) Tag(filePath string, name string, value string) error {
	fileEntry, err := t.ensureFileEntry(filePath)
	if err != nil {
		return err
	}

	var tag Tag
	err = t.db.First(&tag, "file_id = ? AND name = ?", fileEntry.ID, name).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.Wrap(err, "find tag")
		}

		tag.FileID = fileEntry.ID
		tag.Name = name
		tag.Value = value
		err = t.db.Create(&tag).Error
		if err != nil {
			return errors.Wrap(err, "create tag")
		}
	} else {
		tag.FileID = fileEntry.ID
		tag.Name = name
		tag.Value = value
		err = t.db.Save(&tag).Error
		if err != nil {
			return errors.Wrap(err, "update tag")
		}
	}

	return nil
}
