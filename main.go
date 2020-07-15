package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/yevhenshymotiuk/ekatalog-scraper/items"
)

func main() {
	var (
		title   string
		laptops []items.Laptop
	)
	c := colly.NewCollector()

	c.OnHTML("#top-page-title .ib", func(e *colly.HTMLElement) {
		title = e.DOM.Text()
	})

	c.OnHTML(".conf-tr", func(e *colly.HTMLElement) {
		ramCapacity, err := strconv.Atoi(
			strings.TrimSuffix(
				e.DOM.Find(
					".conf-td span[title='Объем оперативной памяти']",
				).Text(),
				"\u00a0ГБ",
			),
		)
		if err != nil {
			log.Fatal(err)
		}
		driveCapacity, err := strconv.Atoi(
			strings.TrimSuffix(
				e.DOM.Find(".conf-td span[title='Емкость накопителя']").Text(),
				"\u00a0ГБ",
			),
		)
		if err != nil {
			log.Fatal(err)
		}

		pricesNode := e.DOM.Find(".price-int")
		pricesSeparator := ".."
		var price items.Price
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
				log.Fatal(err)
			}
			maxPrice, err = strconv.Atoi(priceTexts[1])
			if err != nil {
				log.Fatal(err)
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
				log.Fatal(err)
			}

			price = items.Price{
				Min: minPrice,
			}
		default:
			price = items.Price{}
		}

		laptop := items.Laptop{
			Processor: items.Processor{
				Series: strings.TrimSpace(
					e.DOM.Find(
						".conf-td span[title='Серия процессора']",
					).Text(),
				),
				Model: strings.TrimSpace(
					e.DOM.Find(
						".conf-td span[title='Модель процессора']",
					).Text(),
				),
			},
			RAM: items.RAM{
				Capacity: ramCapacity,
			},
			GPU: items.GPU{
				Model: e.DOM.Find(
					".conf-td span[title='Модель видеокарты']",
				).Text(),
			},
			Drive: items.Drive{
				Type: strings.TrimSpace(
					e.DOM.Find(".conf-td span[title='Тип накопителя']").Text(),
				),
				Capacity: driveCapacity,
			},
			Price: price,
		}
		laptops = append(laptops, laptop)
	})

	c.Visit("https://ek.ua/APPLE-MACBOOK-PRO-13--2020--8TH-GEN-INTEL.htm")

	fmt.Println(title)
	fmt.Printf("%v\n", laptops)
}
