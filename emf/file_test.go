package emf

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func Test(t *testing.T) {
	data, err := ioutil.ReadFile("emf_header.bin")
	if err != nil {
		t.Error("file read err=", err)
	}
	emfFile := ReadFile(data)

	fmt.Printf("header type=%x size=%x \n", emfFile.Header.Type, emfFile.Header.Size)
}
