package concurrency_management

type ConcurrencyResponseDTO struct {
	Id      uint   `json:"id"`
	Type    string `json:"type"`
	Block   bool   `json:"block"`
	Message string `json:"message"`
}
