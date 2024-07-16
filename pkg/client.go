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

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, errors.New("couldn't decode JSON response")

	}

	searchResults, ok := result["query"].(map[string]interface{})["search"].([]interface{})
	if !ok || len(searchResults) == 0 {
		return nil, nil
	}

	articles := make(map[int]Article)
	for i, item := range searchResults {
		entry := item.(map[string]interface{})
		title := entry["title"].(string)
		description := entry["snippet"].(string)
		articles[i] = Article{
			Title:       title,
			Description: CleanWikimediaHTML(description),
			Content:     "",
			Url:         fmt.Sprintf("%s/%s", c.WikiUrl, strings.ReplaceAll(title, " ", "_")),
		}
	}
	return articles, nil
}

func (c *Client) LoadArticle(article Article) {
}