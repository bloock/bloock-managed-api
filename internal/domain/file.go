package domain

import (
	"mime"
	"path/filepath"
	"strings"
)

type File struct {
	File        []byte
	filename    string
	extension   string
	contentType string
}

func NewFile(file []byte, filename, contentType string) File {
	extension := filepath.Ext(filename)
	filename = strings.TrimSuffix(filename, filepath.Ext(filename))

	return File{
		File:        file,
		filename:    filename,
		extension:   extension,
		contentType: contentType,
	}
}

func (f File) Bytes() []byte {
	return f.File
}

func (f File) Filename() string {
	return f.filename
}

func (f File) ContentType() string {
	return f.contentType
}

func (f File) FileExtension() string {
	if f.extension != "" {
		return f.extension
	}

	ext := ""
	exts, err := mime.ExtensionsByType(f.ContentType())
	if err == nil {
		ext = exts[0]
	}
	return ext
}
