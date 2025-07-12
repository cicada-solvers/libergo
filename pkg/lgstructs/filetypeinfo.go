package lgstructs

import (
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

// FileTypeInfo contains information regarding the file type.
type FileTypeInfo struct {
	Name      string
	FileType  string
	MimeType  string
	Header    []byte
	Alias     []string
	Offset    int
	SubHeader []byte
}

// FileTypeInfoModel represents the structure for file type metadata information in the system.
// It includes details like ID, name, file type, MIME type, header, alias, offset, and sub-header.
type FileTypeInfoModel struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	FileType  string `json:"file_type"`
	MimeType  string `json:"mime_type"`
	Header    string `json:"header"`
	Alias     string `json:"alias"`
	Offset    int    `json:"offset"`
	SubHeader string `json:"sub_header"`
}

// fetchFileTypeInfoModels fetches file type metadata from an external API and decodes it into a slice of FileTypeInfoModel structs.
// Returns the slice of FileTypeInfoModel and an error if any occurs during the request or data processing.
func fetchFileTypeInfoModels() ([]FileTypeInfoModel, error) {
	resp, err := http.Get("https://cmbsolver.com/cmbsolver-api/filetypes.php/file_types")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var fileTypeInfoModels []FileTypeInfoModel
	err = json.Unmarshal(body, &fileTypeInfoModels)
	if err != nil {
		return nil, err
	}

	return fileTypeInfoModels, nil
}

// GetAllFileTypeInfo retrieves and processes file type information, returning a slice of FileTypeInfo and an error if any occurs.
func GetAllFileTypeInfo() ([]FileTypeInfo, error) {
	fileTypeInfoModels, err := fetchFileTypeInfoModels()
	if err != nil {
		return nil, err
	}

	var fileTypeInfos []FileTypeInfo
	for _, model := range fileTypeInfoModels {
		header, _ := hex.DecodeString(model.Header)

		var subHeader []byte
		if len(model.SubHeader) > 0 {
			subHeader, _ = hex.DecodeString(model.Header)
		} else {
			subHeader = nil
		}

		var alias []string
		if model.Alias != "" {
			alias = strings.Split(model.Alias, "|")
		} else {
			alias = nil
		}

		fileTypeInfos = append(fileTypeInfos, FileTypeInfo{
			Name:      model.Name,
			FileType:  model.FileType,
			MimeType:  model.MimeType,
			Header:    header,
			Alias:     alias,
			Offset:    model.Offset,
			SubHeader: subHeader,
		})
	}

	return fileTypeInfos, nil
}
