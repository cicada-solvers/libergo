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
