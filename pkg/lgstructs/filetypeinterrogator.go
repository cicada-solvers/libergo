package lgstructs

import (
	"bytes"
	"io"
)

// FileTypeInterrogator is the default implementation with updated definitions.
type FileTypeInterrogator struct {
	definitions []FileTypeInfo
}

// NewFileTypeInterrogator initializes a new instance of FileTypeInterrogator with the default definitions.
func NewFileTypeInterrogator(definitions []FileTypeInfo) *FileTypeInterrogator {
	return &FileTypeInterrogator{definitions: definitions}
}

// GetAvailableExtensions retrieves extensions that are supported based on the current definitions.
func (fti *FileTypeInterrogator) GetAvailableExtensions() []string {
	extensions := make([]string, 0)
	for _, def := range fti.definitions {
		extensions = append(extensions, def.FileType)
	}
	return extensions
}

// GetAvailableMimeTypes retrieves mime types that are supported based on the current definitions.
func (fti *FileTypeInterrogator) GetAvailableMimeTypes() []string {
	mimeTypes := make([]string, 0)
	for _, def := range fti.definitions {
		mimeTypes = append(mimeTypes, def.MimeType)
	}
	return mimeTypes
}

// AvailableTypes retrieves available types that are supported based on the current definitions.
func (fti *FileTypeInterrogator) AvailableTypes() []FileTypeInfo {
	return fti.definitions
}

// DetectType detects the file type based on the file content.
func (fti *FileTypeInterrogator) DetectType(fileContent []byte) *FileTypeInfo {
	for _, def := range fti.definitions {
		if bytes.HasPrefix(fileContent, def.Header) {
			return &def
		}
	}
	return nil
}

// DetectTypeFromStream detects the file type based on the input stream.
func (fti *FileTypeInterrogator) DetectTypeFromStream(inputStream io.Reader) (*FileTypeInfo, error) {
	buffer := make([]byte, 512)
	_, err := inputStream.Read(buffer)
	if err != nil {
		return nil, err
	}
	return fti.DetectType(buffer), nil
}

// IsType determines if the file contents are of a specified type.
func (fti *FileTypeInterrogator) IsType(fileContent []byte, fileType string) bool {
	fileTypeInfo := fti.DetectType(fileContent)
	return fileTypeInfo != nil && fileTypeInfo.FileType == fileType
}
