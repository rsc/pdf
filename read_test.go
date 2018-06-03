package pdf_test

import (
	"image/jpeg"
	"testing"

	"rsc.io/pdf"
)

func TestReaderExtractXObjectDCTDecode(t *testing.T) {
	const (
		testscan = "testdata/testscan.pdf"
	)
	f, err := pdf.Open(testscan)
	if err != nil {
		t.Fatalf("could not open %v: %v", testscan, err)
	}
	x := f.Page(1).Resources().Key("XObject")
	if x.Kind() != pdf.Dict || len(x.Keys()) == 0 {
		t.Fatalf("no xobject dict on page 1")
	}
	k := x.Key(x.Keys()[0])
	if k.IsNull() || k.Kind() != pdf.Stream || k.Key("Subtype").Name() != "Image" {
		t.Fatalf("first xobject child is not an image stream")
	}
	defer func() {
		if r := recover(); r != nil {
			s, ok := r.(string)
			if ok && s == "unknown filter DCTDecode" {
				t.Fatalf("DCTDecode filter handling is not implemented")
			}
			panic(r) // re-panic everything else
		}
	}()
	rc := k.Reader()
	defer rc.Close()
	_, err = jpeg.Decode(rc)
	if err != nil {
		t.Fatalf("could not decode embedded JPEG: %v", err)
	}
}
