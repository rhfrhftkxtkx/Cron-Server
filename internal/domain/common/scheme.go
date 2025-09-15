package common

type Exhibition struct {
	ExhibitionId     string // Unique ID for the exhibition
	Title            string
	VenueVisitKor2Id string
	Period           struct {
		StartDate string
		EndDate   string
	}
	Price          string
	Description    string
	PosterImageUrl string
	SourceUrl      string
	DataSourceTier int // 1, 2, 3
}
