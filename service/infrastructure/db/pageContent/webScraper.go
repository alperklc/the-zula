package pageContent

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/go-shiori/go-readability"
)

type WebScraper interface {
	ScrapPage(URL string) (PageContent, error)
}

type resources struct {
	mdConverter *md.Converter
}

func NewService() WebScraper {
	mdConverter := md.NewConverter("", true, nil)

	return &resources{
		mdConverter: mdConverter,
	}
}

func (r *resources) ScrapPage(rawUrl string) (PageContent, error) {
	client := &http.Client{}
	req, newRequestErr := http.NewRequest("GET", fmt.Sprintf(rawUrl), nil)
	if newRequestErr != nil {

		return PageContent{}, fmt.Errorf("request error: %w", newRequestErr)
	}

	// Some websites won't return content, if the request is not made by an actual browser (see: https://webplatform.news/issues/2020-03-19)
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:79.0) Gecko/20100101 Firefox/79.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")

	resp, httpGetErr := client.Do(req)

	if httpGetErr != nil {
		return PageContent{}, fmt.Errorf("request error: %w", httpGetErr)
	}
	if resp.StatusCode != http.StatusOK {
		return PageContent{}, fmt.Errorf("got status code %d", resp.StatusCode)
	}
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return PageContent{}, fmt.Errorf("read error: %v", readErr)
	}
	defer resp.Body.Close()

	parsedUrl, parseUrlErr := url.ParseRequestURI(rawUrl)
	if parseUrlErr != nil {
		return PageContent{}, parseUrlErr
	}

	article, readErr := readability.FromReader(bytes.NewReader(body), parsedUrl)
	if readErr != nil {
		return PageContent{}, readErr
	}

	markdown, err := r.mdConverter.ConvertString(article.Content)
	if err != nil {
		return PageContent{}, readErr
	}

	pageContent := PageContent{
		URL:         rawUrl,
		Title:       article.Title,
		Length:      article.Length,
		Excerpt:     article.Excerpt,
		SiteName:    article.SiteName,
		Image:       article.Image,
		Favicon:     article.Favicon,
		HTMLContent: article.Content,
		MDContent:   markdown,
	}

	return pageContent, readErr
}
