package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/playwright-community/playwright-go"
)

type Wrapper struct {
	Data struct {
		SearchProductV5 struct {
			Data struct {
				Products []Product `json:"products"`
			} `json:"data"`
		} `json:"searchProductV5"`
	} `json:"data"`
}

type Product struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       struct {
		Text string `json:"text"`
	} `json:"price"`
	Rating string `json:"rating"`
	Shop   struct {
		Name string `json:"name"`
	} `json:"shop"`
	MediaURL struct {
		Image string `json:"image"`
	} `json:"mediaURL"`
}

type Response struct {
	Data struct {
		SearchProductV5 struct {
			Products []Product `json:"products"`
		} `json:"searchProductV5"`
	} `json:"data"`
}

func main() {
	var products []Product
	var result []Wrapper
	pw, err := playwright.Run()
	if err != nil {
		log.Fatal(err)
	}
	defer pw.Stop()

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false),
	})
	if err != nil {
		log.Fatal(err)
	}
	defer browser.Close()

	page, err := browser.NewPage()
	if err != nil {
		log.Fatal(err)
	}

	// CHANNEL (ANTI-DEADLOCK)
	dataChannel := make(chan string, 1)

	// LISTENER (NO BLOCKING CALL)
	page.OnResponse(func(resp playwright.Response) {

		if resp.URL() == "https://gql.tokopedia.com/graphql/SearchProductV5Query" {

			fmt.Println("🔥 GOT GRAPHQL RESPONSE (EVENT)")

			go func() {
				body, err := resp.Body()
				if err != nil {
					log.Println("body error:", err)
					return
				}

				dataChannel <- string(body)

				select {
				case data := <-dataChannel:
					fmt.Println("DATA RECEIVED SIZE:", len(data))

					err := json.Unmarshal([]byte(data), &result)
					if err != nil {
						log.Fatal(err)
					}

					if len(result) != 0 {
						products = append(products, result[0].Data.SearchProductV5.Data.Products...)
						fmt.Println("🔥 Tambah data")
						fmt.Println("TOTAL:", len(products))
					}

					if len(products) < 100 {
						page.Evaluate(`window.scrollTo(0, document.body.scrollHeight)`)

					} else {
						// CREATE CSV
						file, err := os.Create("products.csv")
						if err != nil {
							log.Fatal("could not create file:", err)
						}
						defer file.Close()

						writer := csv.NewWriter(file)
						defer writer.Flush()
						writer.Write([]string{
							"Product Name",
							"Description",
							"Image Link",
							"Price",
							"Rating",
							"Store",
						})

						// WRITE CSV

						for i := 0; i < len(products) && i < 100; i++ {
							p := products[i]
							writer.Write([]string{
								p.Name,
								p.Description,
								p.MediaURL.Image,
								p.Price.Text,
								p.Rating,
								p.Shop.Name,
							})
						}
						fmt.Println("CSV CREATED")
						fmt.Println("DONE XD")
					}

				case <-time.After(15 * time.Second):
					log.Fatal("timeout waiting graphql")
				}
			}()
		}
	})

	// OPEN PAGE
	_, err = page.Goto(
		"https://www.tokopedia.com",
		playwright.PageGotoOptions{
			Timeout: playwright.Float(0),
		},
	)
	page.Goto("https://www.tokopedia.com/search?ob=5&related=true&srp_component_id=04.06.00.00&st=product&q=handphone")
	fmt.Scanln()
}
