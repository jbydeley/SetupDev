package main

import (
	"io/ioutil"
	"net/http"
)

type Download struct {
	Filename     string
	Url          string
	SaveLocation string
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
	wg.Done()
	return nil
}
