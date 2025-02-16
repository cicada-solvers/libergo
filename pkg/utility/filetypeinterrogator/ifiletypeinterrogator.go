package filetypeinterrogator

import (
	"io"
	"liberdatabase"
)

// IFileTypeInterrogator is the interface for interrogating file contents to determine proper file types.
type IFileTypeInterrogator interface {
	GetAvailableExtensions() []string
	GetAvailableMimeTypes() []string
	AvailableTypes() []liberdatabase.FileTypeInfo
	DetectType(fileContent []byte) *liberdatabase.FileTypeInfo
	DetectTypeFromStream(inputStream io.Reader) (*liberdatabase.FileTypeInfo, error)
	IsType(fileContent []byte, fileType string) bool
}
