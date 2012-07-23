package main

import (
	"flag"
	"fmt"
	"net/http"
	"io/ioutil"
)

const configFile = "config.xml"

var (
	help      = flag.Bool("?", false, "Display this help menu")
	config    = flag.String("config", "./config.xml", "Set the configuration file to use")
	full      = flag.Bool("f", false, "Do full development setup")
	downloads = flag.Bool("d", false, "Download files only")
	exports   = flag.Bool("e", false, "Setup exports / environment variables only")
	network   = flag.Bool("n", false, "Grab files off networked drives")
	verbose   = flag.Bool("v", true, "Verbose mode")
)

func main() {
	flag.Parse()

	if *help {
		fmt.Println("Usage: SetupDev [flags]")
		flag.PrintDefaults()
		return
	}

	c := new(Config)
	err := c.Load(configFile)
	if err != nil {
		fmt.Printf("Error: %v\n", err.Error())
	}

	fmt.Printf("Config set to: %v\n", flag.Lookup("config").Value)
	ch := make(chan bool, len(c.Downloads))

	if *full || *downloads {
		for _, v := range c.Downloads {
			fmt.Printf("Downloading '%v' to '%v'\n", v.Filename, v.SaveLocation)
			go Get(v, ch)
		}
	}

	if *full || *exports {
		fmt.Println("* NOTE *")
		fmt.Println("Exports require you to restart your command prompt")
		for _, v := range c.Exports {
			if *verbose {
				fmt.Printf("Exporting %v=%v\n", v.Key, v.Value)
			}
			print("Outside\n")

			if err := v.Set(); err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		}
	}
	if *full || *network {
		body, _ := ioutil.ReadFile("\\\\WTLWF046\\c$\\certreq.txt")
		ioutil.WriteFile("d:/Testthis.txt", body, 0600)
		print(string(body), "\n")
	}

	if *full || *downloads {
	for v := 0; v < len(c.Downloads); v++ {
		success := <-ch
		fmt.Printf("Download: %v\n", success)
	}
}
}

func Get(d Download, ch chan bool) {
	results, err := http.Get(d.Url)
	if err != nil {
		ch <- false
		return
	}
	defer results.Body.Close()
	body, err := ioutil.ReadAll(results.Body)
	if err != nil {
		ch <- false
		return
	}
	ioutil.WriteFile(d.SaveLocation+d.Filename, body, 0600)
	if d.ZipLocation != "" {
		unZip(d)
	}
	ch <- true
}

