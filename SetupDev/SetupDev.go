package main

import (
	"flag"
	"fmt"
	"sync"
)

const configFile = "config.xml"

var (
	help      = flag.Bool("?", false, "Display this help menu")
	config    = flag.String("config", "./config.xml", "Set the configuration file to use.")
	full      = flag.Bool("f", false, "Do full development setup")
	downloads = flag.Bool("d", false, "Download files only")
	exports   = flag.Bool("e", false, "Setup exports / environment variables only.")
	verbose   = flag.Bool("v", true, "Verbose mode")

	wg sync.WaitGroup
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
	if *full || *downloads {
		wg.Add(len(c.Downloads))
		for _, v := range c.Downloads {
			fmt.Printf("Downloading '%v' to '%v'\n", v.Filename, v.SaveLocation)
			go v.Download()
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
	wg.Wait()
}
