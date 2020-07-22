package scraper

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/yevhenshymotiuk/ekatalog-scraper/items"
)

func removeSpaces(s string) string {
	return strings.ReplaceAll(s, "\u00a0", "")
}

func trimCapacitySuffix(s string) string {
	return strings.TrimSuffix(s, "\u00a0ГБ")
}

func specificationsURL(URL string) string {
	productNameRegexp := regexp.MustCompilePOSIX(
		`https:\/\/ek\.ua\/(([[:alnum:]]|\-)+)\.htm`,
	)
	groups := productNameRegexp.FindStringSubmatch(URL)

	return fmt.Sprintf(
		"https://ek.ua/ek-item.php?resolved_name_=%s&view_=tbl",
		groups[1],
	)
}

func scrapeCategory(URL string) (string, error) {
	var category string

	c := colly.NewCollector()

	c.OnHTML(".path_lnk", func(e *colly.HTMLElement) {
		category = e.DOM.Text()
	})

	err := c.Visit(URL)

	return category, err
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

func scrapeSmartphone(table *colly.HTMLElement) (items.Smartphone, error) {
	var smartphone items.Smartphone

	tds := table.DOM.Find(".op3")

	ramCapacity, err := strconv.Atoi(
		trimCapacitySuffix(tds.Get(11).FirstChild.Data),
	)
	if err != nil {
		return smartphone, err
	}
	driveCapacity, err := strconv.Atoi(
		trimCapacitySuffix(tds.Get(12).FirstChild.Data),
	)
	if err != nil {
		return smartphone, err
	}

	minPrice, err := strconv.Atoi(
		removeSpaces(
			table.DOM.Find("span[itemprop='lowPrice']").First().Text(),
		),
	)
	if err != nil {
		return smartphone, err
	}
	maxPrice, err := strconv.Atoi(
		removeSpaces(
			table.DOM.Find("span[itemprop='highPrice']").First().Text(),
		),
	)
	if err != nil {
		return smartphone, err
	}

	smartphone = items.Smartphone{
		CPU: items.CPU{
			Model: strings.TrimSpace(tds.Get(7).FirstChild.Data),
		},
		RAM: items.RAM{
			Capacity: ramCapacity,
		},
		GPU: items.GPU{
			Model: strings.TrimSpace(tds.Get(10).FirstChild.Data),
		},
		Drive: items.Drive{
			Type:     tds.Get(13).FirstChild.Data,
			Capacity: driveCapacity,
		},
		Price: items.Price{
			Min: minPrice,
			Max: maxPrice,
		},
	}

	return smartphone, nil
}

// TODO: Add other categories (not only Laptops)
func scrapeProduct(URL string) (product items.Product, err error) {
	var (
		name          string
		modifications []items.ModificationType
	)

	category, err := scrapeCategory(URL)

	c := colly.NewCollector()

	c.OnHTML("#top-page-title b.ib", func(e *colly.HTMLElement) {
		name = e.DOM.Text()
	})

	switch category {
	case "Ноутбуки":
		// Find row which correspond to modification
		c.OnHTML(".conf-tr", func(e *colly.HTMLElement) {
			laptop := items.Laptop{}
			laptop, err = scrapeLaptop(e)

			modifications = append(modifications, laptop)
		})
	default:
		URL = specificationsURL(URL)

		c.OnHTML(".common-table-div", func(e *colly.HTMLElement) {
			smartphone := items.Smartphone{}
			smartphone, err = scrapeSmartphone(e)

			modifications = append(modifications, smartphone)
		})
	}

	err = c.Visit(URL)

	product = items.Product{Name: name, Modifications: modifications}

	return
}

// ScrapeProducts scrapes products by URLs
func ScrapeProducts(URLs []string) ([]items.Product, error) {
	products := []items.Product{}

	for _, URL := range URLs {
		p, err := scrapeProduct(URL)
		if err != nil {
			return products, err
		}

		products = append(products, p)
	}

	return products, nil
}
