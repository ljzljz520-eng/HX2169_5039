package domain

import "sort"

type ScoreBand struct {
	Name  string
	Min   int
	Max   int
	Order int
}

func DefaultBands() []ScoreBand {
	return []ScoreBand{{Name: "gold", Min: 90, Max: 100, Order: 1}, {Name: "silver", Min: 75, Max: 89, Order: 2}, {Name: "bronze", Min: 60, Max: 74, Order: 3}, {Name: "return", Min: 0, Max: 59, Order: 4}}
}

func BandForScore(score int, bands []ScoreBand) ScoreBand {
	for _, band := range bands {
		if score >= band.Min && score <= band.Max {
			return band
		}
	}
	return ScoreBand{Name: "unrated", Min: 0, Max: 0, Order: 99}
}

func RankRecords(records []Record) []Record {
	copyRecords := append([]Record(nil), records...)
	sort.SliceStable(copyRecords, func(i, j int) bool {
		if copyRecords[i].Score != copyRecords[j].Score {
			return copyRecords[i].Score > copyRecords[j].Score
		}
		return copyRecords[i].ID < copyRecords[j].ID
	})
	return copyRecords
}

func CountByBand(records []Record) map[string]int {
	counts := make(map[string]int)
	for _, record := range records {
		band := BandForScore(record.Score, DefaultBands())
		counts[band.Name]++
	}
	return counts
}

func WeightedScore(jury, audience, technique int) int {
	if jury < 0 {
		jury = 0
	}
	if audience < 0 {
		audience = 0
	}
	if technique < 0 {
		technique = 0
	}
	return (jury*5 + audience*3 + technique*2) / 10
}
