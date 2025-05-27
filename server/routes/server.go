package routes

import (
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Product struct {
	Display     string `json:"label,omitempty"`
	URL         string `json:"url,omitempty"`
	Price       string `json:"price,omitempty"`
	InStock     bool   `json:"inStock,omitempty"`
	Description string `json:"description,omitempty"`
}

func RunServer() {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:5173"},
	}))

	r.GET("/products", getProducts)
	r.GET("/parts", getParts)
	r.GET("/partinfo", getPartInfo)
	r.Run()
}

func getProducts(c *gin.Context) {
	products := []Product{}
	resp, err := http.Get("https://rikontools.com/parts-menu/")
	if err != nil {
		log.Println(err)
		c.Status(400)
	}
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode > 299 {
		log.Printf("Response failed with status code: %d and\nbody: %s\n", resp.StatusCode, body)
		c.Status(400)
	}
	if err != nil {
		log.Println(err)
		c.Status(400)
	}

	bodyString := string(body)
	var re = regexp.MustCompile(`(?m)<option value='https://rikontools.com/product/([%a-zA-Z\d-/]*)'>([a-zA-Z\d-/]*)</option>`)
	results := re.FindAllStringSubmatch(bodyString, -1)
	for _, r := range results {
		if len(r) == 3 {
			product := Product{
				Display: r[2],
				URL:     r[1],
			}

			products = append(products, product)
		}
	}

	c.JSON(200, products)
}

func getParts(c *gin.Context) {
	products := []Product{}
	url := c.Query("partNumber")
	resp, err := http.Get("https://rikontools.com/product/" + url)
	if err != nil {
		log.Println(err)
		c.Status(400)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("Request to %s failed with status code: %d", resp.Request.URL, resp.StatusCode)
		c.Status(400)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Printf("Error while reading the response body, %s", err)
		c.Status(400)
	}

	partsTableRows := doc.Find("tr")
	rows := partsTableRows.Nodes
	for _, r := range rows {
		product := Product{}
		for n := range r.ChildNodes() {
			if n.FirstChild != nil && n.FirstChild.Data == "a" {
				linkNode := n.FirstChild
				if len(linkNode.Attr) > 0 && linkNode.Attr[0].Key == "href" {
					product.URL = linkNode.Attr[0].Val
					product.URL = strings.TrimPrefix(product.URL, "https://rikontools.com/product/")
				}
				if linkNode.FirstChild != nil {
					product.Display = linkNode.FirstChild.Data
				}
				products = append(products, product)
				break
			}
		}
	}

	c.JSON(200, products)
}

func getPartInfo(c *gin.Context) {
	url := c.Query("url")
	resp, err := http.Get("https://rikontools.com/product/" + url)
	if err != nil {
		log.Printf("Failed to get info from rikon site: %+v", err)
		c.Status(400)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("Request to %s failed with status code: %d", resp.Request.URL, resp.StatusCode)
		c.Status(400)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Printf("Error while reading the response body, %s", err)
		c.Status(400)
	}

	productInfoNodes := doc.Find(".entry-summary")
	product := Product{
		URL: url,
	}
	productInfoNode := productInfoNodes.Nodes[0]
	for n := range productInfoNode.ChildNodes() {
		switch n.Data {
		case "p":
			if len(n.Attr) > 0 && n.Attr[0].Key == "class" && n.Attr[0].Val == "stock in-stock" {
				product.InStock = n.FirstChild.Data == "In stock"
			}
			if len(n.FirstChild.Attr) > 0 && n.FirstChild.Attr[0].Val == "woocommerce-Price-amount amount" &&
				n.FirstChild.FirstChild != nil &&
				n.FirstChild.FirstChild.ChildNodes() != nil {
				for cn := range n.FirstChild.FirstChild.ChildNodes() {
					if cn.Type == 1 {
						product.Price = cn.Data
						break
					}
				}
			}
		case "h1":
			product.Display = n.FirstChild.Data
		case "div":
			if len(n.Attr) > 0 && n.Attr[0].Val == "woocommerce-product-details__short-description" {
				log.Printf("Node: %+v", n)
				for cn := range n.ChildNodes() {
					if cn.Data == "p" {
						for ccn := range cn.ChildNodes() {
							if ccn.Type == 1 {
								product.Description = ccn.Data
								break
							}
						}
					}
					log.Printf("Node: %+v", cn)
				}
			}
		}
	}

	c.JSON(200, product)
}
