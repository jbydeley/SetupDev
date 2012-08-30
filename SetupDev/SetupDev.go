package main

import (
	"os"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
)

const configFile = "config.xml"

var (
	help      = flag.Bool("?", false, "Display this help menu")
	config    = flag.String("config", "./config.xml", "Set the configuration file to use")
	full      = flag.Bool("f", false, "Do full development setup")
	downloads = flag.Bool("d", false, "Download files only")
	exports   = flag.Bool("e", false, "Setup exports / environment variables only")
	network   = flag.Bool("n", false, "Grab files off networked drives")
	proxy 	  = flag.String("proxy", "", "Use a proxy when accessing the web")
	verbose   = flag.Bool("v", true, "Verbose mode")
)

func main() {
	flag.Parse()

	if *help {
		fmt.Println("Usage: SetupDev [flags]")
		flag.PrintDefaults()
		return
	}

	if *proxy != "" {
		fmt.Printf("Checking the proxy: %v\n", os.Getenv("HTTP_PROXY"))
		err := os.Setenv("HTTP_PROXY", *proxy)
		if err != nil {
			fmt.Printf("Error setting proxy: %v\n", err)
		}
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
		for _, v := range c.LocalFiles {
			body, _ := ioutil.ReadFile(v.Url)
			ioutil.WriteFile(v.SaveLocation+v.Filename, body, 0600)
		}
	}

	if *full || *downloads {
		for v := 0; v < len(c.Downloads); v++ {
			success := <-ch
			fmt.Printf("Download: %v\n", success)
		}
	}
}

func Get(d FileTransfer, ch chan bool) {
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
