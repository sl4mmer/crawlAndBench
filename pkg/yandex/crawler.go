package yandex

//@see https://github.com/kkhrychikov/kvstorage
//К сожалению оригинальный пакет не экспортирует функции, так что пришлось немножко копипастить
import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"gopkg.in/headzoo/surf.v1"

	"github.com/sl4mmer/crawlAndBench/pkg/common"
)

type Item struct {
	Host string
	Url  string
}

func Search(query string) ([]*Item, error) {
	bow := surf.NewBrowser()
	bow.SetUserAgent(common.DefaultUA)
	err := bow.Open(fmt.Sprintf(baseYandexURL, query))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRequestFailed, err)
	}
	if bow.StatusCode() != 200 {
		return nil, fmt.Errorf("%w : %d", ErrInvalidResponse, bow.StatusCode())
	}

	// в оригинальном парсере в селекторе была опечатка, не знаю уж специально или нет, но факт - было div.serp-Item
	potentialItems := bow.Find("div.serp-item")
	result := make([]*Item, 0, len(potentialItems.Nodes))
	potentialItems.Each(func(_ int, selection *goquery.Selection) {
		item := parseItem(selection)
		if item != nil {
			result = append(result, item)
		}
	})
	return result, nil
}

//вот это взято из исходного парсера, немного уменьшил цикломатическую сложность только #FTGJ
func parseItem(selection *goquery.Selection) *Item {
	_, aExists := selection.Attr("data-fast-name")
	_, cidExists := selection.Attr("data-cid")

	if selection.Is("div.Label") || aExists || selection.Is("span.organic__advLabel") || !cidExists {
		return nil
	}

	link := selection.Find("a.Link").First()
	if link == nil {
		return nil
	}

	urlStr, _ := link.Attr("href")
	dcStr, _ := link.Attr("data-counter")
	if strings.HasPrefix(urlStr, "https://yandex.ru/turbo/") || strings.Contains(urlStr, "turbopages.org") && dcStr != "" {
		var dc []string
		err := json.Unmarshal([]byte(dcStr), &dc)
		if err != nil || len(dc) < 2 {
			return nil
		}
		urlStr = dc[1]
	}

	u, err := url.Parse(urlStr)
	if err != nil {
		return nil
	}

	if u.Host == "" || u.Host == "yabs.yandex.ru" {
		return nil
	}

	return &Item{
		Host: getRootDomain(u.Host),
		Url:  urlStr,
	}
}

func getRootDomain(domain string) string {
	domain = strings.ToLower(domain)

	parts := strings.Split(domain, ".")
	if len(parts) < 3 {
		return domain
	}

	if _, ok := tlds[strings.Join(parts[len(parts)-2:], ".")]; ok {
		return strings.Join(parts[len(parts)-3:], ".")
	}

	return strings.Join(parts[len(parts)-2:], ".")
}
