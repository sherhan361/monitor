package service

import (
	"github.com/sherhan361/monitor/internal/server/repository"
	"log"
	"time"
)

type Backuper struct {
	storeInterval time.Duration
	filename      string
	repo          repository.Getter
}

func NewBackuper(duration time.Duration, filename string, repo repository.Getter) (*Backuper, error) {
	return &Backuper{
		storeInterval: duration,
		filename:      filename,
		repo:          repo,
	}, nil
}

func (b *Backuper) Run() {
	for {
		<-time.After(b.storeInterval)
		err := b.repo.WriteMetrics()
		if err != nil {
			log.Println(err)
		}
	}
}
