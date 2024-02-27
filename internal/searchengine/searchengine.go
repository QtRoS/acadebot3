package searchengine

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/qtros/acadebot3/internal/models"
)

const (
	emptyResult = "[]"
)

type SourceAdapter interface {
	Get(query string, limit int) []models.CourseInfo
	Name() string
}

// TODO Remove from init section.
var adapters = []SourceAdapter{
	newCourseraAdapter(),
	newFuzzyFilteringAdapter(newCachingAdapter(newOpenlearningAdapter(), time.Hour*24)),
	// &udemyAdapter{}, // That's sad, looks like it's blocked
	// newFuzzyFilteringAdapter(newCachingAdapter(&udacityAdapter{}, time.Hour*6)),
	//newFuzzyFilteringAdapter(newCachingAdapter(&iversityAdapter{}, time.Hour*6)),
}

// Search for courses in all services.
func Search(query string, perSourceLimit int) string {
	if query == "" || perSourceLimit <= 0 {
		return emptyResult
	}

	slog.Info(fmt.Sprintf("Gonna search for: %s", query))

	mergedResults := callAdapters(query, perSourceLimit)
	jsonData, err := toJSON(mergedResults)
	if err != nil {
		slog.Error("marshal error: ", slog.Any("err", err))
		return emptyResult
	}

	// slog.Debug(string(json))
	return string(jsonData)
}

func callAdapters(query string, perSourceLimit int) []models.CourseInfo {
	results := make([]models.CourseInfo, 0, perSourceLimit)

	adaptersChunks := make(chan []models.CourseInfo)
	defer close(adaptersChunks)

	var wg sync.WaitGroup
	wg.Add(len(adapters) + 1)

	slog.Info("Before calling adapters...")

	for _, adapter := range adapters {
		go func(adapt SourceAdapter) {
			defer wg.Done()
			adaptersChunks <- adapt.Get(query, perSourceLimit)
		}(adapter)
	}

	go func() {
		defer wg.Done()
		for i := 0; i < len(adapters); i++ {
			chunk := <-adaptersChunks
			results = append(results, chunk...)
		}
	}()

	wg.Wait()
	slog.Info(fmt.Sprintf("Merged result len: %d", len(results)))

	return results
}

// func callAdapters(query string, perSourceLimit int) []CourseInfo {
// 	var results []CourseInfo
// 	for _, adapter := range adapters {
// 		courses := adapter(query, perSourceLimit)
// 		results = append(results, courses...)
// 	}

// 	return results
// }

func toJSON(infos []models.CourseInfo) ([]byte, error) {
	return json.Marshal(infos)
}

func parseJSON(data []byte, target interface{}) error {
	return json.Unmarshal(data, target)
}
