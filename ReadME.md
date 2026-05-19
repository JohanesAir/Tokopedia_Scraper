Tokopedia Product Scraper

This project scrapes product search results from Tokopedia using Playwright and GraphQL interception.

Features
    -Launches Chromium browser using Playwright-Go
    -Intercepts Tokopedia GraphQL search responses
    -Extracts:
        1. Name of Product
        2. Description
        3. Image Link
        4. Price
        5. Rating (out of 5 stars)
        6. Name of store or merchant
    -Exports results into CSV

Tech Stack
    -Golang
    -Playwright-Go

Run
go run main.go

Output
products.csv


Notes
Tokopedia search responses contain mixed ranking results such as products, accessories, and promoted items. Additional filtering can be implemented depending on business needs.