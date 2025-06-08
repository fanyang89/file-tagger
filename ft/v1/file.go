package ft

type FileEntry struct {
	ID   uint   `gorm:"primary_key;autoIncrement"`
	Path string `gorm:"unique"`
}
