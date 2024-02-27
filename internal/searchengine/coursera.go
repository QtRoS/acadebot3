package searchengine

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/qtros/acadebot3/internal/models"
	"github.com/qtros/acadebot3/internal/utils"
)

const (
	courseraApiUrl  = "https://api.coursera.org/api/courses.v1"
	courseraBaseUrl = "http://www.coursera.org/learn/"
)

type courseraResponse struct {
	Elements []courseraElement `json:"elements"`
}

type courseraElement struct {
	ID          string `json:"id"`
	Slug        string `json:"slug"`
	Name        string `json:"name"`
	Description string `json:"description"`
	PhotoUrl    string `json:"photoUrl"`
	Link        string `json:"link"`
}

type courseraAdapter struct {
}

func newCourseraAdapter() *courseraAdapter {
	return &courseraAdapter{}
}

func (me *courseraAdapter) Name() string {
	return "Coursera"
}

func (me *courseraAdapter) Get(query string, limit int) []models.CourseInfo {
	data, err := utils.MakeRequest(courseraApiUrl,
		map[string]string{
			"q":      "search",
			"fields": "description,photoUrl",
			"query":  query,
			//"limit":  strconv.Itoa(limit)
		},
		nil)

	if err != nil {
		slog.Error("marshal error", slog.Any("err", err))
		return nil
	}

	response := courseraResponse{}
	err = parseJSON(data, &response)
	if err != nil {
		slog.Error("marshal error", slog.Any("err", err))
		return nil
	}

	slog.Info(fmt.Sprintf("Results count: %d", len(response.Elements)))

	var infos = make([]models.CourseInfo, 0, limit)
	for _, e := range response.Elements {
		if len(infos) >= limit {
			break
		}

		link := courseraBaseUrl + e.Slug
		desc := strings.Split(e.Description, "\n")[0] //e.Description[:shared.Min(240, len(e.Description))]
		info := models.CourseInfo{Name: e.Name, Headline: desc, Link: link, Art: e.PhotoUrl}
		infos = append(infos, info)
	}

	return infos
}
