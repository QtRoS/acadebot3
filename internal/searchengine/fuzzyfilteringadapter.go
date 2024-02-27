package searchengine

import (
	"strings"

	"github.com/qtros/acadebot3/internal/models"
	"github.com/qtros/acadebot3/internal/searchengine/fuzzy"
)

type fuzzyFilteringAdapter struct {
	sourceAdapter SourceAdapter
}

func newFuzzyFilteringAdapter(adapter SourceAdapter) *fuzzyFilteringAdapter {
	return &fuzzyFilteringAdapter{adapter}
}

func (me *fuzzyFilteringAdapter) Name() string {
	return me.sourceAdapter.Name() + " (Fuzzy)"
}

func (me *fuzzyFilteringAdapter) Get(query string, limit int) []models.CourseInfo {
	courses := me.sourceAdapter.Get(query, limit)

	//slog.Info(me.Name(), "before filter", len(courses))
	queryLower := strings.ToLower(query)
	infos := make([]models.CourseInfo, 0, limit)
	for i := 0; i < len(courses) && len(infos) < limit; i++ {
		ci := &courses[i]
		if fuzzy.Match(queryLower, strings.ToLower(ci.Name)) ||
			fuzzy.Match(queryLower, strings.ToLower(ci.Headline)) {
			infos = append(infos, *ci)
		}
	}
	//slog.Info(me.Name(), "after filter", len(infos))

	return infos
}
