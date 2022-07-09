package input

import (
	"bytes"
	"io"
	"io/fs"
	"os"
)

func NewMockFile(data string, mode fs.FileMode) *MockFile {
	return &MockFile{
		Reader: bytes.NewBufferString(data),
		Size:   int64(len(data)),
		Mode:   mode,
	}
}

type MockFile struct {
	Reader    io.Reader
	Mode      fs.FileMode
	Size      int64
	InputRead bool
}

func (m *MockFile) Stat() (fs.FileInfo, error) { return &mockFileInfo{size: m.Size, mode: m.Mode}, nil }
func (m *MockFile) Close() error               { return nil }
func (m *MockFile) Read(buf []byte) (int, error) {
	m.InputRead = true
	return m.Reader.Read(buf)
}

type mockFileInfo struct {
	fs.FileInfo
	size int64
	mode fs.FileMode
}

func (fi *mockFileInfo) Mode() os.FileMode { return fi.mode }
func (fi *mockFileInfo) Size() int64       { return fi.size }
