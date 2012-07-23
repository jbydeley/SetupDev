package main

import (
	"archive/zip"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Download struct {
	Filename     string
	Url          string
	SaveLocation string
	ZipLocation  string
}

func unZip(d Download) error {
	r, _ := zip.OpenReader(d.SaveLocation + d.Filename)
	defer r.Close()

	for _, f := range r.File {
		f.Open()
		var buf []byte
		os.MkdirAll(filepath.Dir(d.ZipLocation+f.Name), 0600)

		ioutil.WriteFile(d.ZipLocation+f.Name, buf, 0600)
	}
	return nil
}
