package sprocess_test

import (
	"bytes"
	. "github.com/hyperboloide/sprocess"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io"
)

var _ = Describe("Null", func() {

	testBin := genBlob(1 << 22)
	data := NewData()
	n := &Null{"null"}

	id := GenId()

	It("should Write", func() {
		Ω(n.GetName()).To(Equal("null"))
		Ω(n.Start()).To(BeNil())
		w, err := n.NewWriter(id, data)
		Ω(err).To(BeNil())
		Ω(w).ToNot(BeNil())
		l, err := io.Copy(w, bytes.NewReader(testBin))
		w.Close()
		Ω(err).To(BeNil())
		Ω(len(testBin) == int(l)).To(BeTrue())
	})

})
