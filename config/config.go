package config

import (
	"html/template"
	"log"
)

var Values ci = &conf{}

type ci interface {
	Get() *conf
	Set(conf)
}

type conf struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InProduction  bool
	Session       string
	Toast         toast
}

type toast struct {
	Active  bool
	Message string
	Success bool
}

func (m *conf) Get() *conf {
	log.Println("=== model Values Get ===")
	return m
}

func (m *conf) Set(c conf) {
	log.Println("=== model Values Set ===")
	m.UseCache = c.UseCache
	m.InProduction = c.InProduction
	m.Session = c.Session
	m.TemplateCache = c.TemplateCache
	m.Toast = c.Toast
}
