package sprocess_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
	"log"
	"os"
)

func TestSprocess(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Sprocess Suite")
}

var genBlob = func(size int) []byte {
	blob := make([]byte, size)
	for i := 0; i < size; i++ {
		blob[i] = 65 // ascii 'A'
	}
	return blob
}

var testFileReader = func() (*os.File) {
	f, err := os.Open("./tests/test.jpg")
	if err != nil {
		log.Fatal(err)
	}
	return f
}
