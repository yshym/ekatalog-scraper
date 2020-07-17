package scraper

import (
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/yevhenshymotiuk/ekatalog-scraper/items"
)

func trimCapacitySuffix(s string) string {
	return strings.TrimSuffix(s, "\u00a0ГБ")
}

func scrapeLaptop(row *colly.HTMLElement) (items.Laptop, error) {
	var (
		laptop items.Laptop
		price  items.Price
	)

	ramCapacity, err := strconv.Atoi(
		trimCapacitySuffix(
			row.DOM.Find(
				".conf-td span[title='Объем оперативной памяти']",
			).Text(),
		),
	)
	if err != nil {
		return laptop, err
	}
	driveCapacity, err := strconv.Atoi(
		trimCapacitySuffix(
			row.DOM.Find(".conf-td span[title='Емкость накопителя']").Text(),
		),
	)
	if err != nil {
		return laptop, err
	}

	pricesNode := row.DOM.Find(".price-int")
	pricesSeparator := ".."
	switch {
	// Modification record contains both minimal and maximal price
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
	// Modification record contains only minimal price
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
	// Modification doesn't have price
	default:
		price = items.Price{}
	}

	laptop = items.Laptop{
		CPU: items.CPU{
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

// TODO: Add other categories (not only Laptops)
func scrapeProduct(URL string) (product items.Product, err error) {
	c := colly.NewCollector()

	var (
		name          string
		modifications []items.ModificationType
	)

	c.OnHTML("#top-page-title .ib", func(e *colly.HTMLElement) {
		name = e.DOM.Text()
	})

	// Find row which correspond to modification
	c.OnHTML(".conf-tr", func(e *colly.HTMLElement) {
		laptop := items.Laptop{}
		laptop, err = scrapeLaptop(e)

		modifications = append(modifications, laptop)
	})

	err = c.Visit(URL)

	product = items.Product{Name: name, Modifications: modifications}

	return
}

// ScrapeProducts scrapes products by URLs
func ScrapeProducts(URLs []string) ([]items.Product, error) {
	products := []items.Product{}

	for _, URL := range(URLs) {
		p, err := scrapeProduct(URL)
		if err != nil {
			return products, err
		}

		products = append(products, p)
	}

	return products, nil
}
