package issue

import "time"

// Issue represents a git-issue with metadata and content
type Issue struct {
	ID       string    `yaml:"id"`
	Status   string    `yaml:"status"` // "open" or "closed"
	Assignee string    `yaml:"assignee"`
	Labels   []string  `yaml:"labels"`
	Created  time.Time `yaml:"created"`
	Updated  time.Time `yaml:"updated"`
	Title    string    `yaml:"-"` // Not in frontmatter, from markdown heading
	Body     string    `yaml:"-"` // Markdown content after frontmatter
}

// IsOpen returns true if the issue is open
func (i *Issue) IsOpen() bool {
	return i.Status == "open"
}

// IsClosed returns true if the issue is closed
func (i *Issue) IsClosed() bool {
	return i.Status == "closed"
}

// HasLabel checks if the issue has a specific label
func (i *Issue) HasLabel(label string) bool {
	for _, l := range i.Labels {
		if l == label {
			return true
		}
	}
	return false
}
