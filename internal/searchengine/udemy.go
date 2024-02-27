package searchengine

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/qtros/acadebot3/internal/models"
	"github.com/qtros/acadebot3/internal/utils"
)

const (
	udemyAPIURL  = "https://www.udemy.com/api-2.0/courses"
	authHeader   = "Basic MlloUmZ1TXpUSjJLMjJmZWZoSldTeVoyanVtOWx0dkdoWFhFUWZQaTpiNGRIUXhmUDdsODVWa3RHQlM4dUFpdU5ZclpyOEZWY3E3cFpTaWRXbVNMSTBuNm5mWGFyRUxSQ2xqdEtDbjZPcTR3ZkZwWjlqM0RsdU13aUhDN0UxVW1zS1YyQzRtSUlvR2ZEYXpNYVhtbDZjRGtHcjJmOHVqVzVkQ2J5VThaaw=="
	udemyBaseUrl = "https://www.udemy.com"
)

type udemyResponse struct {
	Results []udemyResult `json:"results"`
}

type udemyResult struct {
	ID       int    `json:"id"`
	URL      string `json:"url"`
	Title    string `json:"title"`
	Headline string `json:"headline"`
	Image    string `json:"image_480x270"`
}

type udemyAdapter struct {
}

func (me *udemyAdapter) Name() string {
	return "Udemy"
}

func (me *udemyAdapter) Get(query string, limit int) []models.CourseInfo {
	const fields = "@default,headline"
	data, err := utils.MakeRequest(udemyAPIURL,
		map[string]string{
			"search":         query,
			"page_size":      strconv.Itoa(limit),
			"ordering":       "trending",
			"fields[course]": fields,
		},
		map[string]string{"Authorization": authHeader})

	if err != nil {
		slog.Error("err", slog.Any("err", err))
		return nil
	}

	response := udemyResponse{}
	err = parseJSON(data, &response)
	if err != nil {
		slog.Error("err", slog.Any("err", err))
		return nil
	}

	slog.Info(fmt.Sprintf("Results count %d", len(response.Results)))

	var infos = make([]models.CourseInfo, 0, limit)
	for _, e := range response.Results {
		link := udemyBaseUrl + e.URL
		info := models.CourseInfo{Name: e.Title, Headline: e.Headline, Link: link, Art: e.Image}
		infos = append(infos, info)
	}

	return infos
}
