package common

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"github.com/sherhan361/monitor/internal/models"
)

func GetHash(metric models.Metric, key string) string {
	var hashMetric string
	switch metric.MType {
	case "gauge":
		hashMetric = fmt.Sprintf("%s:gauge:%f", metric.ID, *metric.Value)
	case "counter":
		hashMetric = fmt.Sprintf("%s:counter:%d", metric.ID, *metric.Delta)
	}
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(hashMetric))
	return fmt.Sprintf("%x", h.Sum(nil))
}
