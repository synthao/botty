package telegram

import (
	"bytes"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

type MultipartForm interface {
	FormDataContentType() string
	Form() *bytes.Buffer
	AddField(name, value string) error
	AddFile(name string, f *os.File) error
}

type multipartForm struct {
	buffer *bytes.Buffer
	writer *multipart.Writer
}

func NewMultipartForm() MultipartForm {
	body := &bytes.Buffer{}

	return &multipartForm{
		buffer: body,
		writer: multipart.NewWriter(body),
	}
}

func (m *multipartForm) AddField(name, value string) error {
	return m.writer.WriteField(name, value)
}

func (m *multipartForm) AddFile(name string, f *os.File) error {
	part, err := m.writer.CreateFormFile(name, filepath.Base(f.Name()))
	if err != nil {
		return err
	}
	_, err = io.Copy(part, f)

	return err
}

func (m *multipartForm) Form() *bytes.Buffer {
	defer m.writer.Close()
	return m.buffer
}

func (m *multipartForm) FormDataContentType() string {
	return m.writer.FormDataContentType()
}
