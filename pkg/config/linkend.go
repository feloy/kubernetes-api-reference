package config

import "github.com/feloy/kubernetes-api-reference/pkg/kubernetes"

// LinkEnds maps definition key to a link-end
type LinkEnds map[kubernetes.Key][]string

// Add a new map between key and linkend
func (o LinkEnds) Add(key kubernetes.Key, linkend []string) {
	o[key] = linkend
}
