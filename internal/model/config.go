package model

import "time"

// Config represents the root .prm/prm.json configuration.
type Config struct {
	Version         string   `json:"version"`
	ProjectName     string   `json:"project_name"`
	DefaultPriority Priority `json:"default_priority"`
	DefaultStatus   Status   `json:"default_status"`
	Tags            []string `json:"tags,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

// DefaultConfig returns a new config with sensible defaults.
func DefaultConfig(projectName string) *Config {
	return &Config{
		Version:         "1.0.0",
		ProjectName:     projectName,
		DefaultPriority: PriorityMedium,
		DefaultStatus:   StatusBacklog,
		Tags:            []string{"backend", "frontend", "security", "reviewer", "devops"},
		CreatedAt:       time.Now().UTC(),
	}
}
