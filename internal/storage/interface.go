package storage

import "github.com/sherhan361/monitor/internal/storage/ram"

type Getter interface {
	Set(typ, name, value string) error
}

func NewGetter() (Getter, error) {
	return ram.New(), nil
}
