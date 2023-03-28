package services

import (
	"compress/gzip"
	"fmt"
	browser "github.com/EDDYCJY/fake-useragent"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"sync"
	"sync/atomic"
	"time"
	"views-servive/config"
	"views-servive/models"
)

type ViewsService interface {
	Boost(linkForBoost string) string
}

type ViewsServiceImp struct {
	parseService ParseServiceImp
	settings     config.Settings
}

func NewViewsServiceImp(parseService ParseServiceImp, settings config.Settings) ViewsServiceImp {
	return ViewsServiceImp{
		parseService: parseService,
		settings:     settings,
	}
}

func (v ViewsServiceImp) Boost(linkForBoost string) models.Response {
	var (
		wg      sync.WaitGroup
		re      = regexp.MustCompile(`data-view="(\w+)"`)
		counter int64
	)

	workingUrls := v.parseService.Parse()

	wg.Add(len(workingUrls))

	for _, workingUrl := range workingUrls {
		go func(workingUrl *url.URL, re *regexp.Regexp) {
			defer wg.Done()
			err := makeViews(workingUrl, re, linkForBoost, v.settings)
			if err == nil {
				atomic.AddInt64(&counter, 1)
				log.Println(workingUrl)
			}

		}(workingUrl, re)
	}

	wg.Wait()

	log.Println("закончил накрутку")

	response := models.Response{
		CountOfProxy: len(workingUrls),
		SuccessCount: counter,
	}

	return response
}

func makeViews(workingUrl *url.URL, re *regexp.Regexp, linkForBoost string, settings config.Settings) error {
	randomUA := browser.NewBrowser(browser.Client{}, browser.Cache{}).Random()

	transport := &http.Transport{
		Proxy: http.ProxyURL(workingUrl),
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   60 * time.Second,
	}

	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s?embed=1", linkForBoost), nil)
	if err != nil {
		log.Println(err)
		return err
	}

	request.Header.Set("Accept", settings.AcceptHeader)
	request.Header.Set("Accept-Encoding", settings.AcceptEncodingHeader)
	request.Header.Set("User-Agent", randomUA)

	resp, err := client.Do(request)
	if err != nil {
		return err
	}

	err = acceptView(resp, re, settings, client, linkForBoost, randomUA)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func acceptView(resp *http.Response, re *regexp.Regexp, settings config.Settings, client *http.Client, linkForBoost, randomUA string) error {
	var (
		dataViewString string
		reader         io.ReadCloser
		err            error
	)

	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		err = func(reader io.ReadCloser) error {
			err = reader.Close()
			if err != nil {
				log.Println(err)
				return err
			}
			return nil
		}(reader)
	default:
		reader = resp.Body
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Println(err)
		return err
	}

	if re.Match(data) {
		dataViewString = string(re.FindSubmatch(data)[1])
	} else {
		fmt.Println("Data-views not found.")
		return err
	}

	request, _ := http.NewRequest(http.MethodGet, "https://t.me/v/?views="+dataViewString, nil)
	request.Close = true
	if len(resp.Cookies()) != 0 {
		request.AddCookie(resp.Cookies()[0])
	}

	request.Header.Add("Accept", "*/*")
	request.Header.Add("Accept-Encoding", settings.AcceptEncodingHeader)
	request.Header.Add("Referer", linkForBoost)
	request.Header.Add("User-Agent", randomUA)
	request.Header.Add("X-Requested-With", settings.XRequestedWithHeader)

	resp, _ = client.Do(request)

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	return fmt.Errorf("error")
}
