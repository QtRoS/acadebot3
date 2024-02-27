package searchengine

import (
	"log/slog"
	"strings"

	"github.com/qtros/acadebot3/internal/models"
)

type filteringAdapter struct {
	sourceAdapter SourceAdapter
}

func newFilteringAdapter(adapter SourceAdapter) *filteringAdapter {
	return &filteringAdapter{adapter}
}

func (me *filteringAdapter) Name() string {
	return me.sourceAdapter.Name() + " (Filtered)"
}

func (me *filteringAdapter) Get(query string, limit int) []models.CourseInfo {
	courses := me.sourceAdapter.Get(query, limit)

	slog.Info(me.Name(), "before filter", len(courses))

	infos := make([]models.CourseInfo, 0, limit)
	for i := 0; i < len(courses) && len(infos) < limit; i++ {
		ci := &courses[i]
		if strings.Contains(ci.Name, query) || strings.Contains(ci.Headline, query) {
			infos = append(infos, *ci)
		}
	}

	slog.Info(me.Name(), "after filter", len(infos))

	return infos
}
