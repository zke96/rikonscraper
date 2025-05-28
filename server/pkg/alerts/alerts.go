package alerts

import (
	"fmt"
	"log"
	"net/http"
	"rikonscraper/pkg/config"
	"rikonscraper/pkg/types"

	"github.com/PuerkitoBio/goquery"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robfig/cron/v3"
	"gopkg.in/gomail.v2"
)

var (
	pool   *pgxpool.Pool
	dialer *gomail.Dialer
	cfg    = config.LoadConfig()
)

func HandleAlerts(dbpool *pgxpool.Pool) {
	pool = dbpool

	dialer = gomail.NewDialer("smtp.gmail.com", 587, cfg.EmailFrom, cfg.GmailPass)

	c := *cron.New()
	c.AddFunc("@every "+cfg.DurationString, func() {
		if err := handleAlerts(); err != nil {
			log.Println(err)
		}
	})
	log.Printf("Handling all stored alerts every %s", cfg.DurationString)
	c.Start()
}

func handleAlerts() error {
	var allAlerts []types.AlertRecord
	allAlerts, err := GetAllAlerts()
	if err != nil {
		log.Printf("Failed to get all alerts, %s", err)
	}

	for _, a := range allAlerts {
		stockStatus, err := getProductStatus(a.URL)
		if err != nil {
			log.Printf("Failed to get stock status for product: %s, %s", a.Display, err)
			break
		}

		if err := sendAlert(a, stockStatus); err != nil {
			log.Printf("Failed to send email alert to %s for product %s, %s", a.Email, a.Display, err)
		}
	}

	return nil
}

func getProductStatus(url string) (bool, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to get info from rikon site: %+v", err)
		return false, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("Request to %s failed with status code: %d", resp.Request.URL, resp.StatusCode)
		return false, err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Printf("Error while reading the response body, %s", err)
		return false, err
	}

	selection := doc.Find(".stock")
	if len(selection.Nodes) > 0 && selection.Nodes[0].Data == "p" {
		val := selection.Nodes[0].FirstChild.Data
		return val == "In stock", nil
	}

	return false, fmt.Errorf("failed to find stock value for product: %s", url)
}

func sendAlert(alert types.AlertRecord, inStock bool) error {
	subjectLine := fmt.Sprintf("Stock alert for %s, ", alert.Display)
	stockStatus := ""
	if inStock {
		stockStatus = "In Stock"
	} else {
		stockStatus = "Out of Stock"
	}
	subjectLine += stockStatus

	m := gomail.NewMessage()
	m.SetHeader("From", cfg.EmailFrom)
	m.SetHeader("To", alert.Email)
	m.SetHeader("Subject", subjectLine)
	m.SetBody("text/plain", fmt.Sprintf("Rikon Part %s is %s.\n%s\n\nTo manage your subscriptions visit %salerts/%s", alert.Display, stockStatus, alert.URL, cfg.WebHost, alert.Email))

	if err := dialer.DialAndSend(m); err != nil {
		log.Println(err)
		return err
	}

	log.Printf("Sent stock alert email to %s", alert.Email)
	return nil
}

func SendSubscribeEmail(email string, part types.Product) error {
	subjectLine := fmt.Sprintf("Subscribed to Stock Alerts for Rikon %s, ", part.Display)

	m := gomail.NewMessage()
	m.SetHeader("From", cfg.EmailFrom)
	m.SetHeader("To", email)
	m.SetHeader("Subject", subjectLine)
	m.SetBody("text/plain", fmt.Sprintf("You have subscribed to email alerts for Rikon part number %s. You will receive stock status alerts at this email once a day.\n\nTo manage your subscriptions visit %salerts/%s", part.Display, cfg.WebHost, email))

	if err := dialer.DialAndSend(m); err != nil {
		log.Println(err)
		return err
	}

	log.Printf("Sent subscription email to %s", email)
	return nil
}
