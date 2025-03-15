package lgstructs

import (
	"io"
)

// IFileTypeInterrogator is the interface for interrogating file contents to determine proper file types.
type IFileTypeInterrogator interface {
	GetAvailableExtensions() []string
	GetAvailableMimeTypes() []string
	AvailableTypes() []FileTypeInfo
	DetectType(fileContent []byte) *FileTypeInfo
	DetectTypeFromStream(inputStream io.Reader) (*FileTypeInfo, error)
	IsType(fileContent []byte, fileType string) bool
}
