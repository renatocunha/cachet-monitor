package cachet

import (
	"os"
	"fmt"
	"flag"
	"net/url"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

var Config CachetConfig

type CachetConfig struct {
	API_Url string `json:"api_url"`
	API_Token string `json:"api_token"`
	Monitors []*Monitor `json:"monitors"`
}

func init() {
	var configPath string
	flag.StringVar(&configPath, "c", "/etc/cachet-monitor.config.json", "Config path")
	flag.Parse()

	var data []byte

	// test if its a url
	_, err := url.ParseRequestURI(configPath)
	if err == nil {
		// download config
		response, err := http.Get(configPath)
		if err != nil {
			fmt.Printf("Cannot download network config: %v\n", err)
			os.Exit(1)
		}

		defer response.Body.Close()

		data, _ = ioutil.ReadAll(response.Body)

		fmt.Println("Downloaded network configuration.")
	} else {
		data, err = ioutil.ReadFile(configPath)
		if err != nil {
			fmt.Println("Config file '" + configPath + "' missing!")
			os.Exit(1)
		}
	}

	err = json.Unmarshal(data, &Config)

	if err != nil {
		fmt.Println("Cannot parse config!")
		os.Exit(1)
	}

	if len(os.Getenv("CACHET_API")) > 0 {
		Config.API_Url = os.Getenv("CACHET_API")
	}
	if len(os.Getenv("CACHET_TOKEN")) > 0 {
		Config.API_Token = os.Getenv("CACHET_TOKEN")
	}

	if len(Config.API_Token) == 0 || len(Config.API_Url) == 0 {
		fmt.Printf("API URL or API Token not set. cachet-monitor won't be able to report incidents.\n\nPlease set:\n CACHET_API and CACHET_TOKEN environment variable to override settings.\n\nGet help at https://github.com/CastawayLabs/cachet-monitor\n")
		os.Exit(1)
	}

	if len(Config.Monitors) == 0 {
		fmt.Printf("No monitors defined!\nSee sample configuration: https://github.com/CastawayLabs/cachet-monitor/blob/master/example.config.json\n")
		os.Exit(1)
	}
}