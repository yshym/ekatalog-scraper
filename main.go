package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/yevhenshymotiuk/ekatalog-scraper/scraper"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}
}

func run() error {
	product, err := scraper.ScrapeProduct(
		"https://ek.ua/APPLE-MACBOOK-PRO-13--2020--8TH-GEN-INTEL.htm",
	)
	if err != nil {
		return err
	}

	productJSON, err := json.Marshal(product)
	if err != nil {
		return err
	}

	fmt.Printf("%v\n", string(productJSON))

	return nil
}
