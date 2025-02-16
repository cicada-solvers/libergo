package liberdatabase

import (
	"bufio"
	"config"
	"encoding/hex"
	"fmt"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"strings"
)

// FileTypeInfo contains information regarding the file type.
type FileTypeInfo struct {
	gorm.Model
	Name      string
	FileType  string
	MimeType  string
	Header    []byte
	Alias     []string
	Offset    int
	SubHeader []byte
}

type FileTypeInfoModel struct {
	gorm.Model
	Name      string `gorm:"column:name"`
	FileType  string `gorm:"column:file_type"`
	MimeType  string `gorm:"column:mime_type"`
	Header    string `gorm:"column:header"` //Actually a byte[]
	Alias     string `gorm:"column:alias"`
	Offset    int    `gorm:"column:offset"`
	SubHeader string `gorm:"column:sub_header"` // Actually a byte[]
}

func (FileTypeInfoModel) TableName() string {
	return "public.file_type_info"
}

// LoadDefinitions loads the file type definitions from a file.
func LoadDefinitions() error {
	db, err := InitConnection()
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer func(db *gorm.DB) {
		err := CloseConnection(db)
		if err != nil {
			fmt.Println(err)
		}
	}(db)

	configuration, configError := config.GetConfigFolderPath()
	if configError != nil {
		fmt.Println(err)
		return err
	}

	filePath := filepath.Join(configuration, "definitions_flat")

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(file)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "\t")
		if len(parts) < 6 {
			continue
		}

		offset := 0
		_, err := fmt.Sscanf(parts[0], "%d", &offset)
		if err != nil {
			fmt.Println(err)
			return err
		}
		header := parts[2]
		subHeader := parts[3]
		name := parts[4]
		fileType := parts[5]
		mimeType := parts[6]
		alias := ""
		if len(parts) > 7 {
			alias = parts[7]
		}

		definition := FileTypeInfoModel{
			Name:      name,
			FileType:  fileType,
			MimeType:  mimeType,
			Header:    header,
			Alias:     alias,
			Offset:    offset,
			SubHeader: subHeader,
		}

		// Insert the definition into the database
		if err := db.Create(&definition).Error; err != nil {
			fmt.Println(err)
			return err
		}

		fmt.Printf("Inserted definition: %s\n", name)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func GetAllFileTypeInfo() ([]FileTypeInfo, error) {
	db, err := InitConnection()
	if err != nil {
		return nil, err
	}
	defer func(db *gorm.DB) {
		err := CloseConnection(db)
		if err != nil {
			fmt.Println(err)
		}
	}(db)

	var fileTypeInfoModels []FileTypeInfoModel
	if err := db.Find(&fileTypeInfoModels).Error; err != nil {
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
			Model:     model.Model,
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
