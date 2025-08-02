package ft

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/winfsp/cgofuse/fuse"
)

type TagFileSystem struct {
	fuse.FileSystemBase
	tagger *Tagger
}

func NewTagFileSystem(tagger *Tagger) (fuse.FileSystemInterface, error) {
	return &TagFileSystem{
		tagger: tagger,
	}, nil
}

func (fs *TagFileSystem) splitPath(path string) []string {
	path = strings.Trim(path, "/")
	if path == "" {
		return []string{}
	}
	return strings.Split(path, "/")
}

func (fs *TagFileSystem) getAttrForPath(path string) (*fuse.Stat_t, int) {
	parts := fs.splitPath(path)
	
	if len(parts) == 0 {
		return &fuse.Stat_t{
			Mode:  fuse.S_IFDIR | 0755,
			Nlink: 2,
		}, 0
	}
	
	if len(parts) == 1 {
		var count int64
		err := fs.tagger.db.Model(&Tag{}).Where("name = ?", parts[0]).Count(&count).Error
		if err != nil || count == 0 {
			return nil, -fuse.ENOENT
		}
		
		return &fuse.Stat_t{
			Mode:  fuse.S_IFDIR | 0755,
			Nlink: 2,
		}, 0
	}
	
	if len(parts) == 2 {
		var fileEntry FileEntry
		err := fs.tagger.db.Joins("JOIN tags ON tags.file_id = file_entries.id").
			Where("tags.name = ? AND file_entries.path = ?", parts[0], parts[1]).
			First(&fileEntry).Error
		if err != nil {
			return nil, -fuse.ENOENT
		}
		
		fileInfo, err := os.Stat(fileEntry.Path)
		if err != nil {
			return nil, -fuse.ENOENT
		}
		
		stat := &fuse.Stat_t{
			Mode:    fuse.S_IFLNK | 0644,
			Nlink:   1,
			Size:    int64(len(fileEntry.Path)),
			Atim:    fuse.Timespec{Sec: fileInfo.ModTime().Unix(), Nsec: 0},
			Mtim:    fuse.Timespec{Sec: fileInfo.ModTime().Unix(), Nsec: 0},
			Ctim:    fuse.Timespec{Sec: fileInfo.ModTime().Unix(), Nsec: 0},
		}
		return stat, 0
	}
	
	return nil, -fuse.ENOENT
}

func (*TagFileSystem) Init() {
}

func (*TagFileSystem) Destroy() {
}

func (*TagFileSystem) Statfs(path string, stat *fuse.Statfs_t) int {
	return -fuse.ENOSYS
}

func (*TagFileSystem) Mknod(path string, mode uint32, dev uint64) int {
	return -fuse.ENOSYS
}

func (*TagFileSystem) Mkdir(path string, mode uint32) int {
	return -fuse.ENOSYS
}

func (*TagFileSystem) Unlink(path string) int {
	return -fuse.ENOSYS
}

func (*TagFileSystem) Rmdir(path string) int {
	return -fuse.ENOSYS
}

func (*TagFileSystem) Link(oldpath string, newpath string) int {
	return -fuse.ENOSYS
}

func (*TagFileSystem) Symlink(target string, newpath string) int {
	return -fuse.ENOSYS
}

func (fs *TagFileSystem) Readlink(path string) (int, string) {
	parts := fs.splitPath(path)
	
	if len(parts) == 2 {
		var fileEntry FileEntry
		err := fs.tagger.db.Joins("JOIN tags ON tags.file_id = file_entries.id").
			Where("tags.name = ? AND file_entries.path = ?", parts[0], parts[1]).
			First(&fileEntry).Error
		if err != nil {
			return -fuse.ENOENT, ""
		}
		
		return 0, fileEntry.Path
	}
	
	return -fuse.ENOENT, ""
}

func (*TagFileSystem) Rename(oldpath string, newpath string) int {
	return -fuse.ENOSYS
}

func (*TagFileSystem) Chmod(path string, mode uint32) int {
	return -fuse.ENOSYS
}

func (*TagFileSystem) Chown(path string, uid uint32, gid uint32) int {
	return -fuse.ENOSYS
}

func (*TagFileSystem) Utimens(path string, tmsp []fuse.Timespec) int {
	return -fuse.ENOSYS
}

func (*TagFileSystem) Access(path string, mask uint32) int {
	return -fuse.ENOSYS
}

func (*TagFileSystem) Create(path string, flags int, mode uint32) (int, uint64) {
	return -fuse.ENOSYS, ^uint64(0)
}

func (*TagFileSystem) Open(path string, flags int) (int, uint64) {
	return -fuse.ENOSYS, ^uint64(0)
}

func (fs *TagFileSystem) Getattr(path string, stat *fuse.Stat_t, fh uint64) int {
	result, status := fs.getAttrForPath(path)
	if result != nil {
		*stat = *result
	}
	return status
}

func (*TagFileSystem) Truncate(path string, size int64, fh uint64) int {
	return -fuse.ENOSYS
}

func (*TagFileSystem) Read(path string, buff []byte, ofst int64, fh uint64) int {
	return -fuse.ENOSYS
}

func (*TagFileSystem) Write(path string, buff []byte, ofst int64, fh uint64) int {
	return -fuse.ENOSYS
}

func (*TagFileSystem) Flush(path string, fh uint64) int {
	return -fuse.ENOSYS
}

func (*TagFileSystem) Release(path string, fh uint64) int {
	return -fuse.ENOSYS
}

func (*TagFileSystem) Fsync(path string, datasync bool, fh uint64) int {
	return -fuse.ENOSYS
}

func (fs *TagFileSystem) Opendir(path string) (int, uint64) {
	_, status := fs.getAttrForPath(path)
	return status, 0
}

func (fs *TagFileSystem) Readdir(path string,
	fill func(name string, stat *fuse.Stat_t, ofst int64) bool,
	ofst int64,
	fh uint64) int {
	parts := fs.splitPath(path)
	
	if len(parts) == 0 {
		var tags []Tag
		err := fs.tagger.db.Distinct("name").Find(&tags).Error
		if err != nil {
			return -fuse.EIO
		}
		
		fill(".", &fuse.Stat_t{Mode: fuse.S_IFDIR | 0755}, 0)
		fill("..", &fuse.Stat_t{Mode: fuse.S_IFDIR | 0755}, 0)
		
		for _, tag := range tags {
			if !fill(tag.Name, &fuse.Stat_t{Mode: fuse.S_IFDIR | 0755}, 0) {
				break
			}
		}
		
		return 0
	}
	
	if len(parts) == 1 {
		var files []FileEntry
		err := fs.tagger.db.Joins("JOIN tags ON tags.file_id = file_entries.id").
			Where("tags.name = ?", parts[0]).
			Find(&files).Error
		if err != nil {
			return -fuse.EIO
		}
		
		fill(".", &fuse.Stat_t{Mode: fuse.S_IFDIR | 0755}, 0)
		fill("..", &fuse.Stat_t{Mode: fuse.S_IFDIR | 0755}, 0)
		
		for _, file := range files {
			filename := filepath.Base(file.Path)
			if !fill(filename, &fuse.Stat_t{Mode: fuse.S_IFLNK | 0644}, 0) {
				break
			}
		}
		
		return 0
	}
	
	return -fuse.ENOENT
}

func (*TagFileSystem) Releasedir(path string, fh uint64) int {
	return -fuse.ENOSYS
}

func (*TagFileSystem) Fsyncdir(path string, datasync bool, fh uint64) int {
	return -fuse.ENOSYS
}

func (*TagFileSystem) Setxattr(path string, name string, value []byte, flags int) int {
	return -fuse.ENOSYS
}

func (*TagFileSystem) Getxattr(path string, name string) (int, []byte) {
	return -fuse.ENOSYS, nil
}

func (*TagFileSystem) Removexattr(path string, name string) int {
	return -fuse.ENOSYS
}

func (*TagFileSystem) Listxattr(path string, fill func(name string) bool) int {
	return -fuse.ENOSYS
}
