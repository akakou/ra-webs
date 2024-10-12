package goutils

import (
	"io"
	"os"
	"strconv"
)

type FileMarshal[T any] func(*T) ([]byte, error)
type FileUnmarshal[T any] func([]byte) (*T, error)

type File[T any] struct {
	*os.File
	FileMarshal[T]
	FileUnmarshal[T]
}

func OpenFile[T any](path string, marshal FileMarshal[T], unmarshal FileUnmarshal[T]) (*File[T], error) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	return &File[T]{
		File:          f,
		FileMarshal:   marshal,
		FileUnmarshal: unmarshal,
	}, err
}

func OpenBytesFile(path string) (*File[[]byte], error) {
	marshal := func(b *[]byte) ([]byte, error) { return *b, nil }
	unmarshal := func(b []byte) (*[]byte, error) { return &b, nil }

	return OpenFile(path, marshal, unmarshal)
}

func OpenStringFile(path string) (*File[string], error) {
	marshal := func(s *string) ([]byte, error) {
		data := []byte(*s)
		return data, nil
	}

	unmarshal := func(b []byte) (*string, error) {
		data := string(b)
		return &data, nil
	}

	return OpenFile(path, marshal, unmarshal)
}

func OpenIntFile(path string) (*File[int], error) {
	marshal := func(i *int) ([]byte, error) {
		s := strconv.Itoa(*i)
		return []byte(s), nil
	}

	unmarshal := func(b []byte) (*int, error) {
		i, err := strconv.Atoi(string(b))
		return &i, err
	}

	return OpenFile(path, marshal, unmarshal)
}

func (f *File[T]) Store(data *T) error {
	res, err := f.FileMarshal(data)
	if err != nil {
		return err
	}

	_, err = f.Write(res)
	if err != nil {
		return err
	}

	_, err = f.Seek(0, 0)
	if err != nil {
		return err
	}

	return err

}

func (f *File[T]) Restore() (*T, error) {
	buf, err := io.ReadAll(f.File)
	if err != nil {
		return nil, err
	}

	res, err := f.FileUnmarshal(buf)
	if err != nil {
		return nil, err
	}

	_, err = f.Seek(0, 0)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (f *File[T]) Close() error {
	err := f.File.Close()
	if err != nil {
		return err
	}

	return nil
}
