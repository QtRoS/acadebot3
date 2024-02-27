package searchengine

import (
	"fmt"
	"log/slog"

	"github.com/qtros/acadebot3/internal/models"
	"github.com/qtros/acadebot3/internal/utils"
)

const (
	UdacityApiUrl         = "https://www.udacity.com/public-api/v0/courses"
	UdacityCollectionName = "udacity"
)

type udacityResponse struct {
	Courses []udacityResult `json:"courses"`
}

type udacityResult struct {
	Key          string `json:"key"`
	Homepage     string `json:"homepage"`
	Title        string `json:"title"`
	ShortSummary string `json:"short_summary"`
	Image        string `json:"image"`
}

type udacityAdapter struct {
}

func (me *udacityAdapter) Name() string {
	return "Udacity"
}

func (me *udacityAdapter) Get(query string, limit int) []models.CourseInfo {
	data, err := utils.MakeRequest(UdacityApiUrl, nil, nil)
	if err != nil {
		slog.Error("err", slog.Any("err", err))
		return nil
	}

	response := udacityResponse{}
	err = parseJSON(data, &response)
	if err != nil {
		slog.Error("err", slog.Any("err", err))
		return nil
	}

	uniqueSet := make(map[string]bool)
	var infos = make([]models.CourseInfo, 0, len(response.Courses))
	for _, e := range response.Courses {
		// Check uniqueness.
		if uniqueSet[e.Homepage] {
			slog.Warn("Result dublicate:", slog.Any("homepage", e.Homepage))
			continue
		} else {
			uniqueSet[e.Homepage] = true
		}

		info := models.CourseInfo{Name: e.Title, Headline: e.ShortSummary, Link: e.Homepage, Art: e.Image}
		infos = append(infos, info)
	}

	slog.Info(fmt.Sprintf("Results count %d", len(infos)))

	return infos
}
