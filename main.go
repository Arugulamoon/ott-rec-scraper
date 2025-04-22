package main

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"

	"eden-walker.com/home/ott-rec-scraper/pkg/timefmt"
)

func main() {
	c := colly.NewCollector(
		colly.AllowedDomains("ottawa.ca"),
		colly.CacheDir("./cache"),
	)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println()
		fmt.Println("Visiting", r.URL)
		fmt.Println()
	})

	type Facility struct {
		Name string
	}

	type Event struct {
		Day  string
		Name string
		Time timefmt.TimeFmt
	}

	c.OnHTML("h1[class=page-title]", func(h *colly.HTMLElement) {
		facility := Facility{
			Name: h.DOM.Find("span").Text(),
		}
		fmt.Printf("Facility: %v\n", facility)
	})

	c.OnHTML("div[class=field__item]", func(h *colly.HTMLElement) {
		if strings.HasPrefix(h.DOM.Find("button").Text(), "Drop-in schedule") {
			tbl := h.DOM.Find("table")

			// Get Start and End Date of Recurring Event
			caption := tbl.Find("caption")
			fmt.Println(caption.Text())
			fmt.Println()

			// Get Days
			var days []string
			days = make([]string, 0)
			tbl.Find("thead > tr > th").Each(func(_ int, s *goquery.Selection) {
				days = append(days, strings.TrimSpace(s.Text()))
			})
			if len(days) == 8 {
				days = days[1:]
			}

			var events []Event
			events = make([]Event, 0)
			tbl.Find("tbody > tr").Each(func(_ int, s *goquery.Selection) {
				name := timefmt.SanitizeName(s.Find("th").Text())
				s.Find("td").Each(func(i int, s *goquery.Selection) {
					if !strings.Contains(s.Text(), "n/a") {
						times := timefmt.TranslateEvents(s.Text())
						for _, time := range times {
							events = append(events, Event{
								Day:  days[i],
								Name: name,
								Time: time,
							})
						}
					}
				})
			})
			// fmt.Printf("%#v", events)
			for _, event := range events {
				fmt.Printf("%s %s-%s: %s\n",
					event.Day, event.Time.Start, event.Time.End, event.Name)
			}
		}
	})

	// c.OnResponse(func(r *colly.Response) {
	// 	r.Save(strings.ReplaceAll(r.FileName(), "unknown", "html"))
	// })

	c.Visit("https://ottawa.ca/en/recreation-and-parks/facilities/place-listing/walter-baker-sports-centre")
	c.Visit("https://ottawa.ca/en/recreation-and-parks/facilities/place-listing/minto-recreation-complex-barrhaven")
}
