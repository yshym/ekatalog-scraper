package scraper

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yevhenshymotiuk/ekatalog-scraper/items"
)

func newTestServer() *httptest.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/html", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<title>Test Page</title>
</head>
<body>
<div id="top-page-title">
<b class="ib">Apple MacBook Pro 13 (2020)</b>
</div>
<table>
<tbody>
<tr class="conf-tr">
<td class="conf-td c21"><span title="Серия процессора">Core i5&nbsp;</span></td>
<td class="conf-td c21"><span title="Модель процессора">8257U&nbsp;</span></td>
<td class="conf-td c21"><span title="Объем оперативной памяти">8&nbsp;ГБ</span></td>
<td class="conf-td c21"><span title="Модель видеокарты">Iris Plus Graphics 645&nbsp;</span></td>
<td class="conf-td c21"><span title="Тип накопителя">SSD&nbsp;</span></td>
<td class="conf-td c21"><span title="Емкость накопителя">256&nbsp;ГБ</span></td>
<td class="conf-td conf-price" align="right"><span class="price-int"><span>36&nbsp;949&nbsp;</span>.. <span>43&nbsp;176&nbsp;</span>грн.</span></a></td>
</tr>
</tbody>
</table>
</body>
</html>
		`))
	})

	return httptest.NewServer(mux)
}

func TestScrapeLaptop(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	product, err := ScrapeProduct(ts.URL + "/html")
	if err != nil {
		t.Error("Failed to scrape product")
	}

	wantProduct := items.Product{
		Name: "Apple MacBook Pro 13 (2020)",
		Modifications: []items.ModificationType{
			items.Laptop{
				Processor: items.Processor{
					Series: "Core i5",
					Model:  "8257U",
				},
				RAM: items.RAM{
					Capacity: 8,
				},
				GPU: items.GPU{
					Model: "Iris Plus Graphics 645",
				},
				Drive: items.Drive{
					Type:     "SSD",
					Capacity: 256,
				},
				Price: items.Price{
					Min: 36949,
					Max: 43176,
				},
			},
		},
	}

	productName := product.Name
	wantProductName := wantProduct.Name

	if productName != wantProductName {
		t.Errorf("got: %s, want: %s", productName, wantProductName)
	}

	for i, m := range product.Modifications {
		if m != wantProduct.Modifications[i] {
			t.Errorf("got: %s, want: %s", product, wantProduct)
		}
	}
}