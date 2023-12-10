package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type (
	StorageOpts struct {
		PathTransformFunc func(string) PathKey
	}

	Storage struct {
		StorageOpts
	}

	PathKey struct {
		Pathname, Filename string
	}
)

func SHA1PathTransformerFunc(key string) PathKey {
	sh := sha1.New()
	sh.Write([]byte(key))

	hexPath := hex.EncodeToString(sh.Sum(nil))

	var (
		blockSize = 5
		sliceLen  = len(hexPath) / blockSize
		paths     = make([]string, sliceLen)
	)

	for i := 0; i < sliceLen; i++ {
		from, to := i*blockSize, i*blockSize+blockSize
		paths[i] = hexPath[from:to]
	}

	return PathKey{
		Pathname: strings.Join(paths, string(os.PathSeparator)),
		Filename: hexPath,
	}
}

func NewStorage(opts StorageOpts) *Storage {
	return &Storage{
		opts,
	}
}

func (s *Storage) Delete(key string) error {
	pathKey := s.PathTransformFunc(key)

	return os.RemoveAll(pathKey.FullPath())
}

func (s *Storage) Read(key string) (io.Reader, error) {
	f, err := s.readStream(key)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	if _, err = io.Copy(buf, f); err != nil {
		return nil, err
	}

	if err = f.Close(); err != nil {
		return nil, err
	}

	return buf, nil
}

func (s *Storage) Has(key string) bool {
	pathKey := s.PathTransformFunc(key)

	f, err := os.Stat(pathKey.FullPath())
	if err != nil {
		return false
	}

	return f.Size() != 0
}

func (s *Storage) readStream(key string) (io.ReadCloser, error) {
	pathKey := s.PathTransformFunc(key)

	return os.Open(pathKey.FullPath())
}

func (s *Storage) writeStream(key string, r io.Reader) error {
	pathKey := s.PathTransformFunc(key)
	if err := os.MkdirAll(pathKey.Pathname, os.ModePerm); err != nil {
		return err
	}

	fullPath := pathKey.FullPath()

	f, err := os.Create(fullPath)
	if err != nil {
		return err
	}

	n, err := io.Copy(f, r)
	if err != nil {
		return err
	}

	log.Printf("written (%d) bytes to disk: %s", n, fullPath)

	return nil
}

func (p PathKey) FullPath() string {
	return fmt.Sprintf("%s/%s", p.Pathname, p.Filename)
}
