package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/Lavos/addata"
)

var (
	config_filename = flag.String("c", "", "filename of json configuration file")
)

func awaitQuitKey() {
	var buf [1]byte
	for {
		_, err := os.Stdin.Read(buf[:])
		if err != nil || buf[0] == 'q' {
			return
		}
	}
}

func main() {
	// configuration JSON
	flag.Parse()
	log.Printf("config filename: %#v", *config_filename)

	config_file, err := os.Open(*config_filename)

	if err != nil {
		log.Fatal(err)
	}

	data := make([]byte, 2049)
	n, err := config_file.Read(data)

	config_file.Close()

	if err != nil {
		log.Fatal(err)
	}

	var config addata.Configuration
	err = json.Unmarshal(data[:n], &config)

	if err != nil {
		log.Fatal(err)
	}

	a := addata.NewApplication(&config)

	a.Run()
	awaitQuitKey()
}
