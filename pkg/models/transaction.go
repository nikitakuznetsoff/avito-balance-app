package models

type Transaction struct {
	SenderID	int64	`json:"sender"`
	ReceiverID	int64	`json:"receiver"`
	Value		float64	`json:"value"`
}
