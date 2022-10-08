package middleware

import (
	"github.com/igilgyrg/betera-test/internal/config"
	"github.com/igilgyrg/betera-test/internal/model"
)

type Manager struct {
	origins []string
	cfg     *config.Config
	userUC  model.UserUsecase
}

func NewMiddlewareManager(origins []string, cfg *config.Config, userUC model.UserUsecase) *Manager {
	return &Manager{origins: origins, cfg: cfg, userUC: userUC}
}
