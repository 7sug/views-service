package browser

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/EDDYCJY/fake-useragent/setting"
	"github.com/EDDYCJY/fake-useragent/spiders"
	"github.com/EDDYCJY/fake-useragent/useragent"
	"github.com/EDDYCJY/fake-useragent/useragent/cache"
)

type browser struct {
	Client
	Cache
}

type Client struct {
	MaxPage int
	Delay   time.Duration
	Timeout time.Duration
}

type Cache struct {
	UpdateFile bool
}

var defaultBrowser = NewBrowser(Client{
	MaxPage: setting.BROWSER_MAX_PAGE,
	Delay:   setting.HTTP_DELAY,
	Timeout: setting.HTTP_TIMEOUT,
}, Cache{})

func NewBrowser(client Client, cache Cache) *browser {
	maxPage := setting.GetMaxPage(client.MaxPage)
	delay := setting.GetDelay(client.Delay)
	timeout := setting.GetTimeout(client.Timeout)

	b := browser{
		Client: Client{
			MaxPage: maxPage,
			Delay:   delay,
			Timeout: timeout,
		},
		Cache: Cache{
			UpdateFile: cache.UpdateFile,
		},
	}
	return b.load()
}

func (b *browser) load() *browser {
	fileCache := cache.NewFileCache(cache.GetTempDir(), fmt.Sprintf(setting.TEMP_FILE_NAME, setting.VERSION))
	fileExist, err := fileCache.IsExist()
	if err != nil {
		log.Fatalf("fileCache.IsExist err: %v", err)
	}

	// handle cache.
	if b.UpdateFile == false {
		var (
			isCache      bool
			cacheContent []byte
			m            map[string][]string
		)

		if fileExist == true {
			cacheContent, err = fileCache.Read()
			if err != nil {
				log.Fatalf("fileCache.Read err: %v", err)
			}
			isCache = true
		} else {
			rawCache := cache.NewRawCache(setting.CACHE_URL, fmt.Sprintf(setting.TEMP_FILE_NAME, setting.CACHE_VERSION))
			rawResp, rawExist, err := rawCache.Get()
			if err == nil && rawExist == true {
				defer rawResp.Body.Close()
				rawRead, err := rawCache.Read(rawResp.Body)
				if err == nil && len(rawRead) > 0 {
					cacheContent = rawRead
					isCache = true
				}
			}
		}

		if isCache == true {
			json.Unmarshal(cacheContent, &m)
			useragent.UA.SetData(m)
			if fileExist == false {
				fileCache.WriteJson(useragent.UA.GetAll())
			}
			return b
		}
	}

	// handle origin.
	s := spiders.NewBrowserSpider()
	s.AppendBrowser(b.MaxPage)
	s.StartBrowser(b.Delay, b.Timeout)
	if fileExist == true && b.UpdateFile == true {
		err := fileCache.Remove()
		if err != nil {
			log.Fatalf("fileCache.Remove err: %v", err)
		}
	}

	fileCache.WriteJson(useragent.UA.GetAll())
	return b
}
