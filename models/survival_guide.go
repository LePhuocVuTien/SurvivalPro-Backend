package models

type SurvivalGuide struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	Category   string `json:"category"`
	Difficulty string `json:"difficulty"`
	Icon       string `json:"icon"`
	Content    string `json:"content"`
	ImageURL   string `json:"image_url"`
	Views      int    `json:"views"`
	CreatedAt  string `json:"created_at"`
}
