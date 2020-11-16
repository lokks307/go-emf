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
	emfFile, err := ReadFile(data)
	if err != nil {
		t.Error("emf read err=", err)
	}

	fmt.Printf("header type=%x size=%x \n", emfFile.Header.Type, emfFile.Header.Size)
}
