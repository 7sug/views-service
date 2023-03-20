package services

import (
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"net/url"
	"strings"
	"views-servive/config"
	"views-servive/models"
)

type ParseService interface {
	Parse() []*url.URL
}

type ParseServiceImp struct {
	settings config.Settings
}

func NewParseServiceImp(settings config.Settings) ParseServiceImp {
	return ParseServiceImp{
		settings: settings,
	}
}

func (p ParseServiceImp) Parse() []*url.URL {
	proxies := make([]models.Proxy, 0)
	collector := colly.NewCollector()
	c1 := make(chan []models.Proxy)
	c2 := make(chan []models.Proxy)
	c3 := make(chan []models.Proxy)
	c4 := make(chan []models.Proxy)
	c5 := make(chan []models.Proxy)
	c6 := make(chan []models.Proxy)

	go parseFirstSource(collector, c1)
	go parseSecondSource(collector, c2)
	go parseThirdSource(collector, c3)
	go parseFourthSource(collector, c4)
	go parseFifthSource(collector, c5)
	go parseSixthSource(collector, c6)

	proxies = concat([][]models.Proxy{<-c1, <-c2, <-c3, <-c4, <-c5, <-c6})
	urls := parseToUrl(proxies)

	return urls
}

func parseFirstSource(collector *colly.Collector, c1 chan []models.Proxy) {
	proxies := make([]models.Proxy, 0)
	collector.OnHTML(".table.table-striped.table-bordered > tbody", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(_ int, el *colly.HTMLElement) {
			table := models.Proxy{
				Host: el.ChildText("td:nth-child(1)"),
				Port: el.ChildText("td:nth-child(2)"),
			}
			proxies = append(proxies, table)
		})
	})
	collector.Visit("https://free-proxy-list.net")
	c1 <- proxies
}

func parseSecondSource(collector *colly.Collector, c2 chan []models.Proxy) {
	proxies := make([]models.Proxy, 0)
	collector.OnHTML("table#table_1 > tbody", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(_ int, el *colly.HTMLElement) {
			table := models.Proxy{
				Host: el.ChildText("td:nth-child(1)"),
				Port: el.ChildText("td:nth-child(2)"),
			}
			proxies = append(proxies, table)
		})
	})
	collector.Visit("https://www.vpnside.com/proxy/list/")
	c2 <- proxies
}

func parseThirdSource(collector *colly.Collector, c3 chan []models.Proxy) {
	proxies := make([]models.Proxy, 0)
	collector.OnHTML("table > tbody", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(_ int, el *colly.HTMLElement) {
			table := models.Proxy{
				Host: el.ChildText("td:nth-child(1)"),
				Port: el.ChildText("td:nth-child(2)"),
			}
			proxies = append(proxies, table)
		})
	})
	collector.Visit("https://hidemy.name/en/proxy-list")
	c3 <- proxies
}

func parseFourthSource(collector *colly.Collector, c4 chan []models.Proxy) {
	proxies := make([]models.Proxy, 0)
	collector.OnHTML("div.table__body.js_table__body", func(e *colly.HTMLElement) {
		e.ForEach("div.table__row", func(_ int, el *colly.HTMLElement) {
			url := el.ChildText("div.table__td:nth-child(1)")
			urlList := strings.Split(url, ":")
			table := models.Proxy{
				Host: urlList[0],
				Port: urlList[1],
			}
			proxies = append(proxies, table)
		})
	})
	collector.Visit("https://proxy-store.com/ru/free-proxy")
	c4 <- proxies
}

func parseFifthSource(collector *colly.Collector, c5 chan []models.Proxy) {
	proxies := make([]models.Proxy, 0)
	collector.OnHTML("table.layui-table > tbody", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(_ int, el *colly.HTMLElement) {
			table := models.Proxy{
				Host: el.ChildText("td:nth-child(1)"),
				Port: el.ChildText("td:nth-child(2)"),
			}
			proxies = append(proxies, table)
		})
	})

	for i := 0; i < 5; i++ {
		collector.Visit(fmt.Sprintf("https://www.freeproxy.world/?type=socks5&anonymity=&country=&speed=&port=&page=%v", i))
	}

	c5 <- proxies
}

func parseSixthSource(collector *colly.Collector, c6 chan []models.Proxy) {
	proxies := make([]models.Proxy, 0)
	collector.OnHTML("table.table.table-bordered > tbody", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(_ int, el *colly.HTMLElement) {
			table := models.Proxy{
				Host: el.ChildText("td:nth-child(1)"),
				Port: el.ChildText("td:nth-child(2)"),
			}
			proxies = append(proxies, table)
		})
	})
	collector.Visit("https://proxyhub.me/en/ma-sock5-proxy-list.html")

	c6 <- proxies
}

func concat[T any](slices [][]T) []T {
	var totalLen int

	for _, s := range slices {
		totalLen += len(s)
	}

	result := make([]T, totalLen)

	var i int

	for _, s := range slices {
		i += copy(result[i:], s)
	}

	return result
}

func parseToUrl(proxies []models.Proxy) []*url.URL {
	urls := make([]*url.URL, 0)

	proxiesString := parseToString(proxies)
	for _, proxy := range proxiesString {
		newUrl, err := url.Parse(proxy)
		if err != nil {
			log.Println(err)
			continue
		}

		urls = append(urls, newUrl)
	}

	return urls
}

func parseToString(proxies []models.Proxy) []string {
	proxiesString := make([]string, 0)

	for _, proxy := range proxies {
		correctProxy := strings.Replace(fmt.Sprintf("socks5://%s:%s", proxy.Host, proxy.Port), " ", "", 1)
		if correctProxy != "" {
			proxiesString = append(proxiesString, correctProxy)
		}
	}

	return proxiesString
}
