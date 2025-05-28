package routes

import (
	"log"
	"net/http"
	"net/mail"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"rikonscraper/pkg/alerts"
	"rikonscraper/pkg/products"
	"rikonscraper/pkg/types"
)

func RunServer() {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
	}))

	v0 := r.Group("/v0")

	v0.GET("/products", getProducts)
	v0.GET("/parts/:id", getParts)
	v0.GET("/partinfo/:id", getPartInfo)

	alerts := v0.Group("/alerts")
	alerts.PUT("/", createPartAlert)
	alerts.GET("/:email", getAlertsByEmail)

	r.Run()
}

func getProducts(c *gin.Context) {
	products, err := products.GetAllProducts()
	if err != nil {
		log.Printf("Failed to get products from database, %s", err)
		c.Status(400)
		return
	}

	c.JSON(200, products)
}

func getParts(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		log.Printf("Invalid id parameter, %s", err)
		c.Status(400)
		return
	}

	parts, err := products.GetAllPartsForProduct(id)
	if err != nil {
		log.Printf("Failed to get parts for product %s, %s", id, err)
		c.Status(400)
		return
	}

	c.JSON(200, parts)
}

func getPartInfo(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		log.Printf("Invalid id parameter, %s", err)
		c.Status(400)
		return
	}

	part, err := products.GetPartByID(id)
	if err != nil {
		log.Printf("Failed to get part from db: %+v", err)
		c.Status(400)
		return
	}

	resp, err := http.Get(part.URL)
	if err != nil {
		log.Printf("Failed to get info from rikon site: %+v", err)
		c.Status(400)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("Request to %s failed with status code: %d", resp.Request.URL, resp.StatusCode)
		c.Status(400)
		return
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Printf("Error while reading the response body, %s", err)
		c.Status(400)
		return
	}

	productInfoNodes := doc.Find(".entry-summary")
	product := types.Product{
		ProductCode: part.ProductCode,
		URL:         part.URL,
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

func createPartAlert(c *gin.Context) {
	var newAlert types.AlertRecord

	if err := c.BindJSON(&newAlert); err != nil {
		log.Printf("Failed to bind json, %s", err)
		c.Status(400)
		return
	}

	if !validEmail(newAlert.Email) {
		log.Printf("Invalid email address %s", newAlert.Email)
		c.Status(400)
		return
	}

	if err := alerts.CreateAlert(newAlert.PartID, newAlert.Email); err != nil {
		log.Printf("Failed to insert alert into db, %s", err)
		c.Status(400)
		return
	}

	part, err := products.GetPartByID(newAlert.PartID)
	if err != nil {
		log.Printf("Failed to get part from db, %s", err)
		c.Status(400)
		return
	}

	if err := alerts.SendSubscribeEmail(newAlert.Email, part); err != nil {
		log.Printf("Failed to send subscription email %s", err)
		c.Status(400)
		return
	}

	c.Status(200)
}

func getAlertsByEmail(c *gin.Context) {
	email := c.Param("email")
	if email == "" {
		log.Printf("Invalid email parameter")
		c.Status(400)
		return
	}

	alerts, err := alerts.GetAlertsByEmail(email)
	if err != nil {
		log.Printf("Failed to get alerts from db, %s", err)
		c.Status(400)
		return
	}

	c.JSON(200, alerts)
}

func validEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
