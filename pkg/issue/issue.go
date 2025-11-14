package issue

import "time"

// Issue represents a git-issue with metadata and content
type Issue struct {
	ID       string    `yaml:"id"`
	Assignee string    `yaml:"assignee"`
	Labels   []string  `yaml:"labels"`
	Created  time.Time `yaml:"created"`
	Updated  time.Time `yaml:"updated"`
	Title    string    `yaml:"-"` // Not in frontmatter, from markdown heading
	Body     string    `yaml:"-"` // Markdown content after frontmatter
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
