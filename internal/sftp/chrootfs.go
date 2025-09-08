package sftp

import (
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/pkg/sftp"
)

type ChrootFS struct {
	root string
}

func NewChrootFS(root string) *ChrootFS {
	return &ChrootFS{root: root}
}

func (c *ChrootFS) realPath(requestPath string) string {
	clean := path.Clean("/" + requestPath)
	if clean == "/" {
		return c.root
	}
	return filepath.Join(c.root, clean)
}

func (c *ChrootFS) Fileread(r *sftp.Request) (io.ReaderAt, error) {
	return os.Open(c.realPath(r.Filepath))
}

func (c *ChrootFS) Filewrite(r *sftp.Request) (io.WriterAt, error) {
	flags := os.O_RDWR | os.O_CREATE
	if int(r.Flags)&os.O_TRUNC != 0 {
		flags |= os.O_TRUNC
	}
	if int(r.Flags)&os.O_APPEND != 0 {
		flags |= os.O_APPEND
	}
	return os.OpenFile(c.realPath(r.Filepath), flags, 0644)
}

func (c *ChrootFS) Filecmd(r *sftp.Request) error {
	switch r.Method {
	case "Setstat", "Rename":
		if len(r.Target) > 0 {
			return os.Rename(c.realPath(r.Filepath), c.realPath(r.Target))
		}
	case "Rmdir", "Remove":
		return os.Remove(c.realPath(r.Filepath))
	case "Mkdir":
		return os.Mkdir(c.realPath(r.Filepath), 0755)
	}
	return nil
}

func (c *ChrootFS) Filelist(r *sftp.Request) (sftp.ListerAt, error) {
	f, err := os.Open(c.realPath(r.Filepath))
	if err != nil {
		return nil, err
	}
	fis, err := f.Readdir(0)
	if err != nil {
		return nil, err
	}
	return listerAt(fis), nil
}

type listerAt []os.FileInfo

func (l listerAt) ListAt(ls []os.FileInfo, offset int64) (int, error) {
	if int(offset) >= len(l) {
		return 0, io.EOF
	}
	n := copy(ls, l[offset:])
	if n < len(ls) {
		return n, io.EOF
	}
	return n, nil
}
