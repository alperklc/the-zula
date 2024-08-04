package api

import (
	"time"

	useractivity "github.com/alperklc/the-zula/service/infrastructure/db/userActivity"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

func (i *Insights) ConvertInsights(ag []useractivity.ActivityGraphEntry, mv []useractivity.UsageStatisticsEntry, lastVisited []useractivity.UsageStatisticsEntry, nrOfNotes int64, nrOfBookmarks int64) Insights {

	activities := make([]ActivityOnDate, len(ag))
	for i := range ag {
		count := float32(ag[i].Count)
		date := ag[i].Date

		parsedDate, _ := time.Parse(time.DateOnly, date)

		activities[i] = ActivityOnDate{Count: &count, Date: &openapi_types.Date{parsedDate}}
	}

	i.ActivityGraph = &activities

	mostVisited := make([]MostVisited, len(mv))
	for i := range mv {
		count := float32(mv[i].Count)

		mostVisited[i] = MostVisited{
			Count: &count,
			Id:    &mv[i].ObjectID,
			// Name: mv[i].
			// Title: mv[i]
			Typename: &mv[i].ResourceType,
		}
	}

	i.MostVisited = &mostVisited

	numberOfNotes := float32(nrOfNotes)
	numberOfBookmarks := float32(nrOfBookmarks)

	i.NumberOfNotes = &numberOfNotes
	i.NumberOfBookmarks = &numberOfBookmarks

	return *i
}
