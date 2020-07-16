package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/yevhenshymotiuk/ekatalog-scraper/items"
)

func scrapeLaptop(row *colly.HTMLElement) (items.Laptop, error) {
	var (
		laptop items.Laptop
		price  items.Price
	)

	ramCapacity, err := strconv.Atoi(
		strings.TrimSuffix(
			row.DOM.Find(
				".conf-td span[title='Объем оперативной памяти']",
			).Text(),
			"\u00a0ГБ",
		),
	)
	if err != nil {
		return laptop, err
	}
	driveCapacity, err := strconv.Atoi(
		strings.TrimSuffix(
			row.DOM.Find(".conf-td span[title='Емкость накопителя']").Text(),
			"\u00a0ГБ",
		),
	)
	if err != nil {
		return laptop, err
	}

	pricesNode := row.DOM.Find(".price-int")
	pricesSeparator := ".."
	switch {
	case strings.Contains(pricesNode.Text(), pricesSeparator):
		var minPrice, maxPrice int

		priceTexts := pricesNode.Find(
			"span",
		).Map(
			func(_ int, s *goquery.Selection) string {
				return strings.Replace(
					strings.TrimSpace(s.Text()),
					"\u00a0",
					"",
					-1,
				)
			},
		)

		minPrice, err = strconv.Atoi(priceTexts[0])
		if err != nil {
			return laptop, err
		}
		maxPrice, err = strconv.Atoi(priceTexts[1])
		if err != nil {
			return laptop, err
		}

		price = items.Price{
			Min: minPrice,
			Max: maxPrice,
		}
	case strings.Contains(pricesNode.Text(), "грн"):
		minPrice, err := strconv.Atoi(
			strings.Replace(
				strings.TrimSpace(pricesNode.Find("span").Text()),
				"\u00a0",
				"",
				-1,
			),
		)
		if err != nil {
			return laptop, err
		}

		price = items.Price{
			Min: minPrice,
		}
	default:
		price = items.Price{}
	}

	laptop = items.Laptop{
		Processor: items.Processor{
			Series: strings.TrimSpace(
				row.DOM.Find(
					".conf-td span[title='Серия процессора']",
				).Text(),
			),
			Model: strings.TrimSpace(
				row.DOM.Find(
					".conf-td span[title='Модель процессора']",
				).Text(),
			),
		},
		RAM: items.RAM{
			Capacity: ramCapacity,
		},
		GPU: items.GPU{
			Model: strings.TrimSpace(
				row.DOM.Find(
					".conf-td span[title='Модель видеокарты']",
				).Text(),
			),
		},
		Drive: items.Drive{
			Type: strings.TrimSpace(
				row.DOM.Find(".conf-td span[title='Тип накопителя']").Text(),
			),
			Capacity: driveCapacity,
		},
		Price: price,
	}

	return laptop, nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}
}

func run() (err error) {
	var (
		name          string
		modifications []interface{}
	)
	c := colly.NewCollector()

	c.OnHTML("#top-page-title .ib", func(e *colly.HTMLElement) {
		name = e.DOM.Text()
	})

	c.OnHTML(".conf-tr", func(e *colly.HTMLElement) {
		laptop := items.Laptop{}
		laptop, err = scrapeLaptop(e)

		modifications = append(modifications, laptop)
	})

	err = c.Visit("https://ek.ua/APPLE-MACBOOK-PRO-13--2020--8TH-GEN-INTEL.htm")

	productJSON, err := json.Marshal(
		items.Product{Name: name, Modifications: modifications},
	)

	fmt.Printf("%v\n", string(productJSON))

	return
}
