package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/common-nighthawk/go-figure"
)

type Data struct {
	Title       string
	Description string
	Date        string
}

func main() {
	myFigure := figure.NewFigure("Yavuzlar Web", "", true)
	myFigure.Print()

	var defaultURL1 = "https://thehackernews.com"
	var defaultURL2 = "https://shiftdelete.net/"
	var defaultURL3 = "https://www.webtekno.com/"

	for _, arg := range os.Args {
		if arg == "-h" || arg == "--help" {
			printHelp()
			return
		}
	}

	urls := []string{}

	for _, v := range os.Args {
		if v == "-1" {
			urls = append(urls, defaultURL1)
		}

		if v == "-2" {
			urls = append(urls, defaultURL2)
		}

		if v == "-3" {
			urls = append(urls, defaultURL3)
		}
	}

	var descriptionFilter, dateFilter bool
	for _, arg := range os.Args[1:] {
		if arg == "-description" {
			descriptionFilter = true
		} else if arg == "-date" {
			dateFilter = true
		}
	}

	for _, url := range urls {
		res, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}

		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			fmt.Println("URL alınamadı:", err)
			continue
		}

		if strings.Contains(url, "thehackernews.com") {
			doc.Find(".body-post").Each(func(i int, sa *goquery.Selection) {
				title := sa.Find("h2").Text()
				description := sa.Find("div.home-desc").Text()
				date := sa.Find("div .item-label").Find("span").Text()

				printData(title, description, date, descriptionFilter, dateFilter)
			})
		}

		if strings.Contains(url, "shiftdelete.net") {
			fmt.Println(url)
			doc.Find("div .post-inner-content").Each(func(i int, sa *goquery.Selection) {
				title := sa.Find("div .post-title").Find("h4").Find("span").Text()
				description := sa.Find("div .post-excerpt").Find("p").Text()
				date := sa.Find("div.thb-date").Text()

				printData(title, description, date, descriptionFilter, dateFilter)

			})

		}

		if strings.Contains(url, "webtekno.com") {
			fmt.Println(url)
			doc.Find("div .content-timeline--right").Each(func(i int, sa *goquery.Selection) {
				title := sa.Find("div .content-timeline__detail__container").Find("h3").Text()
				description := sa.Find("span.content-timeline--underline").Text()
				date := sa.Find("content-timeline__detail__category").Text()

				printData(title, description, date, descriptionFilter, dateFilter)
			})
		}
	}

}

func printData(title, description, date string, descriptionFilter, dateFilter bool) {
	if title != "" {
		fmt.Printf("Haber başlığı : %s\n", title)
	}

	if descriptionFilter {
		fmt.Printf("Haber içeriği : %s\n", description)
	}

	if dateFilter {
		fmt.Printf("Haber tarihi: %s\n", date)
	}

	fmt.Println()
}

func printHelp() {
	fmt.Println("Kullanım: go run main.go [SEÇENEKLER] [PARAMETRELER]")
	fmt.Println("Seçenekler:")
	fmt.Println("-1, -2, -3 : '-' ile takip edilen bir sayı (1, 2 ve 3) bir web sitesini temsil eder:")
	fmt.Println("  -1: thehackernews.com")
	fmt.Println("  -2: shiftdelete.net")
	fmt.Println("  -3: webtekno.com")
	fmt.Println("paramereler:")
	fmt.Println("date: '-date' kullanarak haberin yayınlanma saatinin görüntülenmesi.")
	fmt.Println("description: '-description' kullanarak haberin açıklama kısmının görüntülenmesi.")
	fmt.Println("Parametreler aynı anda kullanılabilir.")
}
