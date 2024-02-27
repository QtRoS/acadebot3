package searchengine

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/qtros/acadebot3/internal/models"
	"github.com/qtros/acadebot3/internal/utils"
)

const (
	//openLearningAPIURL = "https://www.openlearning.com/api/courses/list?type=free,paid"
	openLearningAPIURL = "https://api.openlearning.com/v2.1/courses/?featured=false"
	envAPIKey          = "ENV_OL_API_KEY"
)

type openlearningMeta struct {
	Status int    `json:"status"`
	Err    string `json:"error"`
}

type openlearningEnvelope struct {
	Meta openlearningMeta     `json:"meta"`
	Data []openlearningResult `json:"data"`
}

type openlearningResult struct {
	Id        string `json:"id"`
	CourseUrl string `json:"url"`
	Image     string `json:"image"`
	Name      string `json:"title"`
	Summary   string `json:"summary"`
}

type openlearningAdapter struct {
	apiKey string
}

func newOpenlearningAdapter() *openlearningAdapter {
	return &openlearningAdapter{
		apiKey: os.Getenv(envAPIKey),
	}
}

func (me *openlearningAdapter) Name() string {
	return "OpenLearning"
}

func (me *openlearningAdapter) Get(query string, limit int) []models.CourseInfo {

	if me.apiKey == "" {
		slog.Error("API Key is empty")
		return nil
	}

	data, err := utils.MakeRequest(openLearningAPIURL,
		map[string]string{
			"featured":     "false",
			"itemsPerPage": "1000",
		},
		map[string]string{
			"Accept":    "application/json",
			"X-API-Key": me.apiKey,
		})

	if err != nil {
		slog.Error("err", slog.Any("err", err))
		return nil
	}

	response := openlearningEnvelope{}
	err = parseJSON(data, &response)
	if err != nil {
		slog.Error("err", slog.Any("err", err))
		return nil
	}

	slog.Info(fmt.Sprintf("%s results count %d", me.Name(), len(response.Data)))

	var infos = make([]models.CourseInfo, len(response.Data))
	for i, e := range response.Data {
		headline := strings.Split(e.Summary, "\n")[0] //e.Summary[:shared.Min(240, len(e.Summary))]
		info := models.CourseInfo{Name: e.Name, Headline: headline, Link: e.CourseUrl, Art: e.Image}
		infos[i] = info
	}

	return infos
}
