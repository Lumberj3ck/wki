package pkg

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	Lang    string
	WikiUrl string
	ApiUrl  string
}

func NewClient(lang string, unformattedWikiUrl string, unformattedApiUrl string) (*Client, error) {
	if _, ok := WikipediaLangs[lang]; !ok {
		return nil, fmt.Errorf("wikipedia language %s does not exist", lang)
	}
	client := &Client{
		Lang:    lang,
		WikiUrl: fmt.Sprintf("https://%s.%s", lang, unformattedWikiUrl),
		ApiUrl:  fmt.Sprintf("https://%s.%s", lang, unformattedApiUrl),
	}
	return client, nil
}

func (c *Client) QueryArticles(queryText string) (map[int]Article, error) {
	if strings.TrimSpace(queryText) == "" {
		return nil, nil
	}

	params := url.Values{}
	params.Add("action", "query")
	params.Add("list", "search")
	params.Add("srsearch", queryText)
	params.Add("utf8", "")
	params.Add("format", "json")
	params.Add("srlimit", "3")
	params.Add("srprop", "snippet")

	apiUrl := c.ApiUrl + params.Encode()

	resp, err := http.Get(apiUrl)
	if err != nil {
		return nil, errors.New("couldn't fetch data from Wikipedia API")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("couldn't read response body")

	}

	var result WikipediaPageQueryJSON
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, errors.New("couldn't decode JSON response")
	}

	if len(result.Query.Search) == 0 {
		return nil, nil
	}

	articles := make(map[int]Article)
	for i, entry := range result.Query.Search {
		articles[i] = Article{
			Title:       entry.Title,
			Description: CleanWikimediaHTML(entry.Snippet),
			Content:     "",
			Url:         fmt.Sprintf("%s/%s", c.WikiUrl, strings.ReplaceAll(entry.Title, " ", "_")),
		}
	}
	return articles, nil
}

func (c *Client) LoadArticle(article Article) (Article, error) {

	params := url.Values{}
	params.Add("action", "query")
	params.Add("formatversion", "2")
	params.Add("prop", "revisions")
	params.Add("rvprop", "content")
	params.Add("rvslots", "*")
	params.Add("titles", article.Title)
	params.Add("format", "json")

	apiUrl := c.ApiUrl + params.Encode()

	resp, err := http.Get(apiUrl)
	if err != nil {
		return article, errors.New("couldn't fetch data from Wikipedia API")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return article, errors.New("couldn't read response body")
	}

	var result WikipediaPageJSON
	err = json.Unmarshal(body, &result)
	if err != nil {
		return article, errors.New("couldn't decode JSON response")
	}

	if len(result.Query.Pages) == 0 {
		return article, errors.New("no pages found")
	}

	page := result.Query.Pages[0]
	article.Title = page.Title
	article.Content = page.Revisions[0].Slots.Main.Content
	return article, nil
}
