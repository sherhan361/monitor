package repository

import (
	"github.com/sherhan361/monitor/internal/server/repository/memory"
)

type Getter interface {
	Set(typ, name, value string) error
}

func NewGetter() (Getter, error) {
	return memory.New(), nil
}
