package products

import (
	"context"
	"io"
	"log"
	"net/http"
	"regexp"
	"rikonscraper/pkg/config"
	"rikonscraper/pkg/types"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robfig/cron/v3"
)

var (
	pool *pgxpool.Pool
	cfg  = config.LoadConfig()
)

func HandleUpdateProducts(dbpool *pgxpool.Pool) {
	pool = dbpool

	c := *cron.New()
	c.AddFunc("@every "+"01h", func() {
		log.Printf("Updating all products")
		if err := getAllProductsFromSite(); err != nil {
			log.Println(err)
		}
		log.Printf("Updating all parts")
		if err := getPartsForAllProducts(); err != nil {
			log.Println(err)
		}
	})
	log.Println("Updating products every 01h")
	c.Start()
}

func getAllProductsFromSite() error {
	products := []types.Product{}
	resp, err := http.Get("https://rikontools.com/parts-menu/")
	if err != nil {
		log.Println(err)
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("Request to %s failed with status code: %d", resp.Request.URL, resp.StatusCode)
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return err
	}

	bodyString := string(body)
	var re = regexp.MustCompile(`(?m)<option value='https://rikontools.com/product/([%a-zA-Z\d-/]*)'>([a-zA-Z\d-/]*)</option>`)
	results := re.FindAllStringSubmatch(bodyString, -1)
	for _, r := range results {
		if len(r) == 3 {
			product := types.Product{
				Display:     r[2],
				ProductCode: strings.TrimSuffix(r[1], "/"),
				URL:         cfg.RikonURL + r[1],
			}

			products = append(products, product)
		}
	}

	for _, p := range products {
		if _, err := pool.Exec(context.Background(), "INSERT INTO rikonscraper.products (display, url, product_code) VALUES ($1, $2, $3) ON CONFLICT (product_code) DO NOTHING;", p.Display, p.URL, p.ProductCode); err != nil {
			log.Printf("Failed to insert alert into db, %s", err)
		}
	}

	log.Printf("Done Updating Products")
	return nil
}

func getPartsForAllProducts() error {
	products, err := GetAllProducts()
	if err != nil {
		log.Printf("Failed to get products from database, %s", err)
		return err
	}

	for _, p := range products {
		log.Printf("Updating parts for product %s", p.Display)
		parts, err := getPartsForProduct(p.URL)
		if err != nil {
			log.Printf("Failed to get parts for product %s, %s", p.Display, err)
		}

		for _, part := range parts {
			if _, err := pool.Exec(context.Background(), "INSERT INTO rikonscraper.parts (parent, display, url, product_code) VALUES ($1, $2, $3, $4) ON CONFLICT (parent, product_code) DO NOTHING;", p.ID, part.Display, part.URL, part.ProductCode); err != nil {
				log.Printf("Failed to insert part into db, %s", err)
			}
		}
	}

	log.Printf("Done Updating Parts")
	return nil
}

func getPartsForProduct(url string) ([]types.Product, error) {
	parts := []types.Product{}
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return parts, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("Request to %s failed with status code: %d", resp.Request.URL, resp.StatusCode)
		return parts, err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Printf("Error while reading the response body, %s", err)
		return parts, err
	}

	partsTableRows := doc.Find("tr")
	rows := partsTableRows.Nodes
	for _, r := range rows {
		part := types.Product{}
		for n := range r.ChildNodes() {
			if n.FirstChild != nil && n.FirstChild.Data == "a" {
				linkNode := n.FirstChild
				if len(linkNode.Attr) > 0 && linkNode.Attr[0].Key == "href" {
					part.URL = linkNode.Attr[0].Val
					part.ProductCode = linkNode.Attr[0].Val
					part.ProductCode = strings.TrimPrefix(part.ProductCode, cfg.RikonURL)
					part.ProductCode = strings.TrimSuffix(part.ProductCode, "/")
				}
				if linkNode.FirstChild != nil {
					part.Display = linkNode.FirstChild.Data
				}
				parts = append(parts, part)
				break
			}
		}
	}

	return parts, nil
}
