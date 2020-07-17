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
	URLs := []string{
		"https://ek.ua/APPLE-MACBOOK-PRO-13--2020--8TH-GEN-INTEL.htm",
		"https://ek.ua/APPLE-MACBOOK-PRO-13--2020--10TH-GEN-INTEL.htm",
	}

	products, err := scraper.ScrapeProducts(URLs)
	if err != nil {
		return err
	}

	productsJSON, err := json.Marshal(products)
	if err != nil {
		return err
	}

	fmt.Printf("%v\n", string(productsJSON))

	return nil
}
