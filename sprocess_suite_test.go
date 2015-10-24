package sprocess_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"log"
	"os"
	"io/ioutil"
	"testing"
)

func TestSprocess(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Sprocess Suite")
}

var TMP string

var _ = BeforeSuite(func() {
	var err error
    TMP, err = ioutil.TempDir("", "")
	Ω(err).To(BeNil())
})

var _ = AfterSuite(func() {
    os.RemoveAll(TMP)
})

var genBlob = func(size int) []byte {
	blob := make([]byte, size)
	for i := 0; i < size; i++ {
		blob[i] = 65 // ascii 'A'
	}
	return blob
}

var testFileReader = func() *os.File {
	f, err := os.Open("./tests/test.jpg")
	if err != nil {
		log.Fatal(err)
	}
	return f
}

var tmpDir = func() string {
	 d, err := ioutil.TempDir(TMP, "")
	 if err != nil {
		log.Fatal(err)
	}
	return d
}
