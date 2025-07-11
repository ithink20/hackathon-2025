package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"hackathon-2025/pkg/models"
)

type ConfluenceService struct {
	BaseURL string
	Headers map[string]string
}

func NewConfluenceService() *ConfluenceService {
	return &ConfluenceService{
		BaseURL: "https://confluence.shopee.io",
		Headers: map[string]string{
			"accept":             "*/*",
			"accept-language":    "en-US,en;q=0.9",
			"cache-control":      "no-cache, no-store, must-revalidate",
			"expires":            "0",
			"pragma":             "no-cache",
			"priority":           "u=1, i",
			"sec-ch-ua":          `"Google Chrome";v="137", "Chromium";v="137", "Not/A)Brand";v="24"`,
			"sec-ch-ua-mobile":   "?0",
			"sec-ch-ua-platform": `"macOS"`,
			"sec-fetch-dest":     "empty",
			"sec-fetch-mode":     "cors",
			"sec-fetch-site":     "same-origin",
			"user-agent":         "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36",
		},
	}
}

func (cs *ConfluenceService) GetPagesByUser(email string) ([]models.PageInfo, error) {
	var allPages []models.PageInfo
	limit := 100
	start := 0

	for {
		searchURL := fmt.Sprintf("%s/rest/api/search", cs.BaseURL)

		params := url.Values{}
		params.Set("cql", fmt.Sprintf(`contributor in ("%s") AND type in ("page")`, email))
		params.Set("start", fmt.Sprintf("%d", start))
		params.Set("limit", fmt.Sprintf("%d", limit))
		params.Set("excerpt", "highlight")
		params.Set("expand", "space.icon")
		params.Set("includeArchivedSpaces", "false")
		params.Set("src", "next.ui.search")

		req, err := http.NewRequest("GET", searchURL+"?"+params.Encode(), nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		for key, value := range cs.Headers {
			req.Header.Set(key, value)
		}

		req.Header.Set("Cookie", "mywork.tab.tasks=false; confluence.list.pages.cookie=list-content-tree; confluence.last-web-item-clicked=system.space.tools%2Fcontenttools%2Fbrowse; mo.confluence-oauth.FORM_COOKIE=loginform; confluence-language=en_US; _gid=GA1.2.467484768.1752124827; space_auth_live=MTc1MjE1NTk4NnxOd3dBTkRkRlJFeExUVFkwUVRVMFRVZFVVMVl6U1VjMVdVbE5ObFEzVTFoUVdrNDJTbFUzVUZnMVdqUXlRa1ZMTjFJMFdUZFpSVkU9fJbZu3M5RkSBbatR9GVsapnSzcDxmcofYubL198JTdwh; JSESSIONID=0CD0585B00B37625557A1482C2C4FDEC; _ga=GA1.1.1600566590.1752124827; _ga_VPBMX0QP83=GS2.1.s1752236872$o5$g1$t1752236884$j48$l0$h0")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to make request: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
		}

		var confluenceResp models.ConfluenceSearchResponse
		if err := json.Unmarshal(body, &confluenceResp); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}

		for _, result := range confluenceResp.Results {
			allPages = append(allPages, models.PageInfo{
				ID:    result.Content.ID,
				Type:  result.Content.Type,
				Title: result.Content.Title,
			})
		}

		if start+len(confluenceResp.Results) >= confluenceResp.TotalSize {
			break
		}

		start += limit
	}

	return allPages, nil
}

func (cs *ConfluenceService) GetPageContent(pageID string) (string, error) {
	contentURL := fmt.Sprintf("%s/plugins/viewstorage/viewpagestorage.action?pageId=%s", cs.BaseURL, pageID)

	req, err := http.NewRequest("GET", contentURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("accept-language", "en-US,en;q=0.9")
	req.Header.Set("cache-control", "max-age=0")
	req.Header.Set("priority", "u=0, i")
	req.Header.Set("sec-ch-ua", `"Google Chrome";v="137", "Chromium";v="137", "Not/A)Brand";v="24"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"macOS"`)
	req.Header.Set("sec-fetch-dest", "document")
	req.Header.Set("sec-fetch-mode", "navigate")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("sec-fetch-user", "?1")
	req.Header.Set("upgrade-insecure-requests", "1")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36")

	req.Header.Set("Cookie", "mywork.tab.tasks=false; confluence.list.pages.cookie=list-content-tree; confluence.last-web-item-clicked=system.space.tools%2Fcontenttools%2Fbrowse; mo.confluence-oauth.FORM_COOKIE=loginform; confluence-language=en_US; _gid=GA1.2.467484768.1752124827; space_auth_live=MTc1MjE1NTk4NnxOd3dBTkRkRlJFeExUVFkwUVRVMFRVZFVVMVl6U1VjMVdVbE5ObFEzVTFoUVdrNDJTbFUzVUZnMVdqUXlRa1ZMTjFJMFdUZFpSVkU9fJbZu3M5RkSBbatR9GVsapnSzcDxmcofYubL198JTdwh; JSESSIONID=25A05E649B6ECA7FD71A2B50D0F9BFE2; _ga_VPBMX0QP83=GS2.1.s1752221241$o3$g1$t1752221612$j60$l0$h0; _ga=GA1.2.1600566590.1752124827; _gat_gtag_UA_156269607_2=1")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return string(body), nil
}

func (cs *ConfluenceService) GetPagesByUserWithContent(email string) ([]models.PageInfo, error) {
	pages, err := cs.GetPagesByUser(email)
	if err != nil {
		return nil, err
	}

	for i := range pages {
		content, err := cs.GetPageContent(pages[i].ID)
		if err != nil {
			log.Printf("Failed to fetch content for page %s: %v", pages[i].ID, err)
			continue
		}
		pages[i].Content = content
	}

	return pages, nil
}
