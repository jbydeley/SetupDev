package main

import (
	"encoding/xml"
	"io"
	"io/ioutil"
	"os"
)

type Config struct {
	Downloads []Download     `xml:"Downloads>Download"`
	Exports   []Export `xml:"Exports>Export"`
}

func (c *Config) Load(fileName string) error {
	results, err := ioutil.ReadFile(fileName)
	if err != nil {
		if _, ok := err.(*os.PathError); ok {
			c.SaveHelpConfig()
		} else {
			return err
		}
	}

	err = xml.Unmarshal(results, &c)
	if err != nil && err != io.EOF {
		return err
	}
	return nil
}

func (c *Config) Save() error {
	body, _ := xml.MarshalIndent(c, "", "  ")
	return ioutil.WriteFile(configFile, body, 0600)
}

func (c *Config) SaveHelpConfig() error {

	c.Downloads = []Download{
		Download{
			Filename:     "go1.0.2.windows-amd64.msi",
			Url:          "http://go.googlecode.com/files/go1.0.2.windows-amd64.msi",
			SaveLocation: "./"}}

	c.Exports = []Export{
		Export{
			Key:   "GOPATH",
			Value: "./"}}

	return c.Save()
}
