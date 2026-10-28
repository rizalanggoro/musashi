package costestimate

type CostEstimateDTO struct {
	ProjectID string `json:"project_id"`
	Budget    string `json:"budget"`
	Rab       string `json:"rab"`
	Spb       string `json:"spb"`
	Variance  string `json:"variance"`
	Sisa      string `json:"sisa"`
}
