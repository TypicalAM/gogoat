package apiclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/TypicalAM/gogoat/config"
	"github.com/charmbracelet/lipgloss"
	"github.com/guptarohit/asciigraph"
)

var ErrApi = errors.New("api error")

var ppHeaderStyle = lipgloss.NewStyle().
	MarginBottom(1).
	Padding(1, 1, 1, 1).
	Background(lipgloss.Color("#a7e2a2")).
	Foreground(lipgloss.Color("#000000")).
	Bold(true)

type singleErrorData struct {
	Error string `json:"error"`
}

type multipleErrorData struct {
	Errors map[string][]string `json:"errors"`
}

type data struct {
	Message interface{} `json:"data"`
}

type Caller struct {
	authToken string
	basesite  string
}

func NewCaller(cfg config.Config) *Caller {
	return &Caller{
		authToken: cfg.Token,
		basesite:  cfg.Site,
	}
}

type TotalPageViews struct {
	Total       int `json:"total"`
	TotalEvents int `json:"total_events"`
	TotalUTC    int `json:"total_utc"`
}

func (c *Caller) GetTotalPageViews() (*TotalPageViews, error) {
	respBytes, err := c.getResult(fmt.Sprintf("%s/api/v0/stats/total", c.basesite), "GET")
	if err != nil {
		return nil, fmt.Errorf("total page views: %w", err)
	}

	var tpv TotalPageViews
	if json.Unmarshal(respBytes, &tpv) != nil {
		return nil, fmt.Errorf("total page views: %w", err)
	}

	return &tpv, nil
}

type HitData struct {
	Day    string `json:"day"`
	Hourly []int  `json:"hourly"`
	Daily  int    `json:"daily"`
}

type Hit struct {
	Path   string    `json:"path"`
	Title  string    `json:"title"`
	Stats  []HitData `json:"stats"`
	Count  int       `json:"count"`
	PathID int       `json:"path_id"`
	Max    int       `json:"max"`
	Event  bool      `json:"event"`
}

type TotalHits struct {
	Hits []Hit `json:"hits"`
}

func (th TotalHits) PrettyPrint() string {
	var b strings.Builder

	if len(th.Hits) == 0 {
		b.WriteString(ppHeaderStyle.Render("No hits found"))
		return b.String()
	}

	for i, hit := range th.Hits {
		b.WriteString(lipgloss.JoinHorizontal(
			lipgloss.Left,
			ppHeaderStyle.Render(fmt.Sprintf("Title: %s, hits: %d, path: %s", hit.Title, hit.Count, hit.Path)),
			th.Plot(i),
		))
		b.WriteRune('\n')
	}

	return b.String()
}

func (c *Caller) GetTotalHits() (*TotalHits, error) {
	respBytes, err := c.getResult(fmt.Sprintf("%s/api/v0/stats/hits", c.basesite), "GET")
	if err != nil {
		return nil, fmt.Errorf("hits: %w", err)
	}

	var th TotalHits
	if json.Unmarshal(respBytes, &th) != nil {
		return nil, fmt.Errorf("hits: %w", err)
	}

	return &th, nil
}

func (th TotalHits) Plot(siteIndex int) string {
	max := func(x []float64) int {
		var maxVal float64
		for _, v := range x {
			if v > maxVal {
				maxVal = v
			}
		}
		return int(maxVal)
	}

	data := make([]float64, len(th.Hits[siteIndex].Stats))
	for i, v := range th.Hits[siteIndex].Stats {
		data[i] = float64(v.Daily)
	}
	_ = max(data)

	return asciigraph.Plot(
		data,
		asciigraph.Precision(0),
		asciigraph.Height(2),
		asciigraph.Width(20),
		asciigraph.CaptionColor(2),
		asciigraph.LabelColor(2),
	)
}

func (c *Caller) getResult(url, method string) ([]byte, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.authToken))

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 300 {
		return body, nil
	}

	var singleError singleErrorData
	var multipleError multipleErrorData

	err = json.Unmarshal(body, &singleError)
	if err != nil {
		err = json.Unmarshal(body, &multipleError)
		if err != nil {
			return nil, ErrApi
		}
	}

	if singleError.Error != "" {
		return nil, fmt.Errorf("api error: %s", singleError.Error)
	} else if multipleError.Errors != nil {
		return nil, fmt.Errorf("api error: %v", multipleError.Errors)
	} else {
		return nil, ErrApi
	}
}
