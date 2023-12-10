package main

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStorage(t *testing.T) {
	opts := StorageOpts{
		PathTransformFunc: SHA1PathTransformerFunc,
	}

	s := NewStorage(opts)
	key := "test_key"
	expected := PathKey{
		Pathname: "00942/f4668/670f3/4c594/3cf52/c7ef3/139fe/2b8d6",
		Filename: "00942f4668670f34c5943cf52c7ef3139fe2b8d6",
	}
	assert.Equal(t, expected, s.PathTransformFunc(key))

	data := []byte("some file")
	err := s.writeStream(key, bytes.NewReader(data))
	assert.NoError(t, err)

	r, err := s.Read(key)
	assert.NoError(t, err)

	b, err := io.ReadAll(r)
	assert.NoError(t, err)
	assert.NotNil(t, b)
	assert.Equal(t, b, data)
}

func TestPathTransformFunc(t *testing.T) {
	tCases := []struct {
		key      string
		expected PathKey
	}{
		{
			key: "your_mother's_best_picture",
			expected: PathKey{
				Pathname: "ec94f/cd580/4a16a/7276d/8dff2/20b3e/30249/c9c3a",
				Filename: "ec94fcd5804a16a7276d8dff220b3e30249c9c3a",
			},
		},
		{
			key: "your_sister's_best_picture",
			expected: PathKey{
				Pathname: "c8e9e/33c91/740ec/b17a8/959ad/5fafb/98b8e/acaae",
				Filename: "c8e9e33c91740ecb17a8959ad5fafb98b8eacaae",
			},
		},
	}
	for _, tc := range tCases {
		assert.Equal(t, tc.expected, SHA1PathTransformerFunc(tc.key))
	}
}

func TestStorageDelete(t *testing.T) {
	opts := StorageOpts{
		PathTransformFunc: SHA1PathTransformerFunc,
	}

	s := NewStorage(opts)
	key := "test_key"
	data := []byte("some file")

	err := s.writeStream(key, bytes.NewReader(data))
	assert.NoError(t, err)

	err = s.Delete(key)
	assert.NoError(t, err)
	assert.False(t, s.Has(key))
}
