package ft

import "github.com/winfsp/cgofuse/fuse"

type TagFileSystem struct {
	fuse.FileSystemBase
}

func NewTagFileSystem() (fuse.FileSystemInterface, error) {
	return &TagFileSystem{}, nil
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

func (*TagFileSystem) Readlink(path string) (int, string) {
	return -fuse.ENOSYS, ""
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

func (*TagFileSystem) Getattr(path string, stat *fuse.Stat_t, fh uint64) int {
	return -fuse.ENOSYS
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

func (*TagFileSystem) Opendir(path string) (int, uint64) {
	return -fuse.ENOSYS, ^uint64(0)
}

func (*TagFileSystem) Readdir(path string,
	fill func(name string, stat *fuse.Stat_t, ofst int64) bool,
	ofst int64,
	fh uint64) int {
	return -fuse.ENOSYS
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
