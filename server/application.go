package main

import (
	"log"
	"os"

	"rikonscraper/routes"
)

// const rikonURL = "https://rikontools.com/product/"
// const gmailPass = "dshaktakukekkhfk"
// const emailTo = "zkerikson96@gmail.com"
// const emailFrom = "zane95782@gmail.com"

func main() {
	err := os.Remove("rikonscraper.log")
	if err != nil {
		log.Fatalf("error removing file: %v", err)
	}
	f, err := os.OpenFile("rikonscraper.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	routes.RunServer()

	// productCode = os.Args[1]
	// log.Printf("Product Code: %s", productCode)
	// productURL = rikonURL + productCode
}

// func updateStock() error {
// 	inStock, err := fetch()
// 	if err != nil {
// 		return fmt.Errorf("error getting stock status: %s", err)
// 	}

// 	stockStatus := "out of stock"
// 	if inStock {
// 		stockStatus = "in stock"
// 	}

// 	log.Printf("Part number %s is %s\n", productCode, stockStatus)

// 	err = send(productURL, productCode, inStock)
// 	if err != nil {
// 		return fmt.Errorf("error sending stock alert: %s", err)
// 	}

// 	return nil
// }

// func fetch() (bool, error) {
// 	resp, err := http.Get(productURL)
// 	if err != nil {
// 		log.Println(err)
// 		return false, err
// 	}
// 	body, err := io.ReadAll(resp.Body)
// 	resp.Body.Close()
// 	if resp.StatusCode > 299 {
// 		log.Printf("Response failed with status code: %d and\nbody: %s\n", resp.StatusCode, body)
// 		return false, err
// 	}
// 	if err != nil {
// 		log.Println(err)
// 		return false, err
// 	}

// 	bodyString := string(body)

// 	return strings.Contains(bodyString, "In stock"), nil
// }

// func send(url, productCode string, inStock bool) error {
// 	subjectLine := fmt.Sprintf("Stock alert for %s, ", productCode)
// 	if inStock {
// 		subjectLine += "In Stock!"
// 	} else {
// 		subjectLine += "Out of Stock"
// 	}

// 	m := gomail.NewMessage()
// 	m.SetHeader("From", emailFrom)
// 	m.SetHeader("To", emailTo)
// 	m.SetHeader("Subject", subjectLine)
// 	m.SetBody("text/plain", url)

// 	d := gomail.NewDialer("smtp.gmail.com", 587, emailFrom, gmailPass)

// 	if err := d.DialAndSend(m); err != nil {
// 		log.Println(err)
// 		return err
// 	}

// 	log.Printf("Sent stock alert email to %s", emailTo)
// 	return nil
// }

// func onReady() {
// 	var err error

// 	go func() {
// 		for {
// 			select {
// 			case <-mQuitOrig.ClickedCh:
// 				log.Println("Requesting quit")
// 				cmd.Process.Kill()
// 				log.Println("Finished quitting")
// 			case <-checkStock.ClickedCh:
// 				log.Println("Checking stock status now")
// 				if err = updateStock(); err != nil {
// 					log.Println(err)
// 				}
// 			case <-viewLogs.ClickedCh:
// 				go func() {
// 					command := "Get-Content 'C:\\Users\\Zane\\Coding\\Rikon Scraper\\rikonscraper.log' -Wait -Tail 30"
// 					cmd = *exec.Command("PowerShell", "-Command", command)
// 					cmd.Stdout = os.Stdout
// 					cmd.Stderr = os.Stderr
// 					err = cmd.Run()
// 					if err != nil {
// 						log.Println(err)
// 					}
// 				}()
// 			}
// 		}
// 	}()

// 	if err = updateStock(); err != nil {
// 		log.Println(err)
// 	}

// 	c = *cron.New()
// 	c.AddFunc("30 7 * * *", func() {
// 		if err := updateStock(); err != nil {
// 			log.Println(err)
// 		}
// 	})
// 	c.Start()
// }

// func onExit() {
// 	c.Stop()
// 	os.Exit(1)
// }
