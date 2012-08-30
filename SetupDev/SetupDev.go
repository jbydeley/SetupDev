package main

import (
	"os"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
)



var (
	help      = flag.Bool("?", false, "Display this help menu")
	config    = flag.String("config", "config.xml", "Set the configuration file to use")
	full      = flag.Bool("f", false, "Do full development setup")
	downloads = flag.Bool("d", false, "Download files only")
	exports   = flag.Bool("e", false, "Setup exports / environment variables only")
	network   = flag.Bool("n", false, "Grab files off networked drives")
	proxy 	  = flag.String("proxy", "", "Use a proxy when accessing the web")
	verbose   = flag.Bool("v", false, "Verbose mode")
	c 		  = new(Config)
)

func main() {
	flag.Parse()

	if *help {
		fmt.Println("Usage: SetupDev [flags]")
		flag.PrintDefaults()
		return
	}

	c = HandleConfig()
	
	if *proxy != "" {
		HandleProxy()
	}
	if *full {
		HandleFull()
	} else {
		if *downloads {
			ch := HandleDownloads()
			FinishDownloads(ch)
		}

		if *exports {
			HandleExports()
		}

		if *network {
			HandleLocal()
		}
	}
	
	fmt.Println(c.Instructions)
	
}

func HandleConfig() *Config {
	co := new(Config)
	err := co.Load(*config)
	if err != nil {
		fmt.Printf("Error: %v\n", err.Error())
	}

	Display(fmt.Sprintf("Config set to: %v", *config))
	return co
}

func HandleProxy() {
	fmt.Printf("Checking the proxy: %v\n", os.Getenv("HTTP_PROXY"))
	err := os.Setenv("HTTP_PROXY", *proxy)
	if err != nil {
		fmt.Printf("Error setting proxy: %v\n", err)
	}
}

func HandleFull() {
	ch := HandleDownloads();
	HandleExports();
	HandleLocal();
	FinishDownloads(ch);
}

func HandleDownloads() chan bool {
	ch := make(chan bool, len(c.Downloads))
	for _, v := range c.Downloads {
		Display(fmt.Sprintf("Downloading '%v' to '%v'", v.Filename, v.SaveLocation))
		go Get(v, ch)
	}
	return ch
}

func FinishDownloads( ch chan bool ) {
	for v := 0; v < len(c.Downloads); v++ {
		success := <-ch
		
		if success == false {
			fmt.Println("Download failed. Check that the location exists and that your http_proxy environment variable is setup correctly (if you use a proxy).")
		} else {
			Display("Download: Success")
		}
	}
}

func HandleExports() {
	fmt.Println("* NOTE *")
	fmt.Println("Exports require you to restart your command prompt")
	for _, v := range c.Exports {
		if *verbose {
			fmt.Printf("Exporting %v=%v\n", v.Key, v.Value)
		}

		if err := v.Set(); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}
}

func HandleLocal() {
	for _, v := range c.LocalFiles {
		body, err := ioutil.ReadFile(v.Url)
		if err != nil {
			fmt.Printf("%v Error:\n -> %v\n", v.Filename, err)
		}
		ioutil.WriteFile(v.SaveLocation+v.Filename, body, 0600)
	}
}

func Display(text string) {
	if *verbose {
		fmt.Println(text)
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
