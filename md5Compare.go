package main

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

// CompareMD5 receives two filenames and generates checksums for both to check if they are identical
func CompareMD5(filename1 string, filename2 string) bool {
	file, err := os.Open(filename1)

	if err != nil {
		fmt.Println(err)
	}

	defer file.Close()

	file2, err := os.Open(filename2)
	if err != nil {
		fmt.Println(err)
	}
	defer file2.Close()

	hash := md5.New()
	_, err = io.Copy(hash, file)

	if err != nil {
		fmt.Println(err)
	}

	hash2 := md5.New()
	_, err = io.Copy(hash2, file2)

	return bytes.Equal(hash.Sum(nil), hash2.Sum(nil))
}
