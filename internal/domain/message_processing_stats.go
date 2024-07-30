package domain

type MessageProcessingStats struct {
	SavedTotal     int `json:"savedTotal"`
	ProcessedTotal int `json:"processedTotal"`
}
