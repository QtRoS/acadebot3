package searchengine

import (
	"fmt"
	"log/slog"

	"github.com/qtros/acadebot3/internal/models"
	"github.com/qtros/acadebot3/internal/utils"
)

const (
	iversityAPIURL = "https://iversity.org/api/v1/courses"
)

type iversityResponse struct {
	Courses []iversityResult `json:"courses"`
}

type iversityResult struct {
	URL      string `json:"url"`
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	Image    string `json:"image"`
}

type iversityAdapter struct {
}

func (me *iversityAdapter) Name() string {
	return "Iversity"
}

func (me *iversityAdapter) Get(query string, limit int) []models.CourseInfo {
	data, err := utils.MakeRequest(iversityAPIURL, nil, nil)
	if err != nil {
		slog.Error("err", slog.Any("err", err))
		return nil
	}

	response := iversityResponse{}
	err = parseJSON(data, &response)
	if err != nil {
		slog.Error("err", slog.Any("err", err))
		return nil
	}

	slog.Info(fmt.Sprintf("Results count %d", len(response.Courses)))

	var infos = make([]models.CourseInfo, len(response.Courses))
	for i, e := range response.Courses {
		info := models.CourseInfo{Name: e.Title, Headline: e.Subtitle, Link: e.URL, Art: e.Image}
		infos[i] = info
	}

	return infos
}
