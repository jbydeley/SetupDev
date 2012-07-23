package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"archive/zip"
)

type Download struct {
	Filename     string
	Url          string
	SaveLocation string
	ZipLocation  string
}

func (d *Download) Download() error {
	results, err := http.Get(d.Url)
	if err != nil {
		return err
	}
	defer results.Body.Close()
	body, err := ioutil.ReadAll(results.Body)
	if err != nil {
		return err
	}
	ioutil.WriteFile(d.SaveLocation+d.Filename, body, 0600)
	if d.ZipLocation != "" {
		unZip(d)
	}
	wg.Done()
	return nil
}

func unZip(d *Download) error {
	r, err := zip.OpenReader(d.SaveLocation + d.Filename)
	if err != nil {
		print(err.Error(), "\n")
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			print(err.Error(), "\n")
		}
		var buf []byte

		_, err = io.ReadFull(rc, buf)
		err = ioutil.WriteFile(d.ZipLocation+f.Name, buf, 0600)
		if err != nil {
			print(err.Error(), "\n")
		}
		rc.Close()
	}
	return nil
}
