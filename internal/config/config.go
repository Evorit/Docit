package config

import (
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Title       string            `yaml:"title"`
	Description string            `yaml:"description"`
	SourceDir   string            `yaml:"sourceDir"`
	OutputDir   string            `yaml:"outputDir"`
	Theme       ThemeConfig       `yaml:"theme"`
	Nav         []NavItem         `yaml:"nav"`
	Sidebar     []SidebarGroup    `yaml:"sidebar"`
}

type ThemeConfig struct {
	PrimaryColor string `yaml:"primaryColor"`
	Logo         string `yaml:"logo"`
}

type NavItem struct {
	Text string `yaml:"text"`
	Link string `yaml:"link"`
}

type SidebarGroup struct {
	Title string   `yaml:"title"`
	Items []string `yaml:"items"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if cfg.SourceDir == "" {
		cfg.SourceDir = "docs"
	}
	if cfg.OutputDir == "" {
		cfg.OutputDir = "dist"
	}
	if cfg.Title == "" {
		cfg.Title = "Documentation"
	}
	if cfg.Theme.PrimaryColor == "" {
		cfg.Theme.PrimaryColor = "#3b82f6"
	}

	return &cfg, nil
}

func (c *Config) GetMarkdownFiles() ([]string, error) {
	var files []string
	err := filepath.Walk(c.SourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".md") {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func (c *Config) GetOutputPath(mdPath string) string {
	relPath, err := filepath.Rel(c.SourceDir, mdPath)
	if err != nil {
		return ""
	}
	htmlPath := strings.TrimSuffix(relPath, ".md") + ".html"
	return filepath.Join(c.OutputDir, htmlPath)
}