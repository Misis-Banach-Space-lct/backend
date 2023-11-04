package router

import (
	"lct/internal/config"

	"github.com/gofiber/fiber/v2"
)

type Router struct {
	engine *fiber.App
}

func NewRouter() (*Router, error) {
	r := &Router{
		engine: fiber.New(),
	}

	if err := r.setup(); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *Router) setup() error {
	return nil
}

func (r *Router) Run() error {
	return r.engine.Listen(":" + config.Cfg.ServerPort)
}
