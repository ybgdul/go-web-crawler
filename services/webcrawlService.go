package services

import (
	"net/http"
	"sync"
	"web-crawler/models"

	"github.com/PuerkitoBio/goquery"
)

type CrawlService struct{ 
	client *http.Client
}

func NewCrawlService(http *http.Client) *CrawlService{ 
	return &CrawlService{
		client: http,
	}
}

func (s *CrawlService) Crawl(urls []string) ([]models.Result, error) { 

	ch := make(chan models.Result, len(urls))
	var wg sync.WaitGroup

	for _, url := range urls { 
		wg.Add(1)
		go s.fetchURL(ch, url, &wg)
	}

	wg.Wait()
	close(ch)

	result := make([]models.Result,0, len(urls))

	for res := range ch { 
		if res.Err != nil {
			return nil, res.Err
		}
		result = append(result, res)
	}

	return result, nil
}

func (s *CrawlService) fetchURL(ch chan<- models.Result, url string, wg *sync.WaitGroup)  { 
	defer wg.Done()

	resp, err := s.client.Get(url)

	if err != nil { 
		ch <- models.Result{
			URL: url,
			Err: err,
		}
		return
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil { 
		ch <- models.Result{
			URL: url,
			Err: err,
		}
		return
	}

	title := doc.Find("title").Text()

	ch <- models.Result{Title: title, URL: url, Err: nil}
}