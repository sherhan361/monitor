package models

type Metric struct {
	ID    string   `json:"id"`              // Имя метрики
	MType string   `json:"type"`            // Параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // Значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // Значение метрики в случае передачи gauge
	Hash  string   `json:"hash,omitempty"`  // Значение хеш-функции
}

type WriteMetric struct {
	ID    string  `json:"id"`              // Имя метрики
	MType string  `json:"type"`            // Параметр, принимающий значение gauge или counter
	Delta int64   `json:"delta,omitempty"` // Значение метрики в случае передачи counter
	Value float64 `json:"value,omitempty"` // Значение метрики в случае передачи gauge
}
