package config

import (
	"html/template"
	"log"
)

// Values interface variable allowing access to config data
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

// Get values from config file
func (m *conf) Get() *conf {
	log.Println("=== model Values Get ===")
	return m
}

// Set changes in config file values based on input
func (m *conf) Set(c conf) {
	log.Println("=== model Values Set ===")
	m.UseCache = c.UseCache
	m.InProduction = c.InProduction
	m.Session = c.Session
	m.TemplateCache = c.TemplateCache
	m.Toast = c.Toast
}
