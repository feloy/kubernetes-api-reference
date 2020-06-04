package kubernetes

import "strings"

// APIGroup represents the group of a Kubernetes API
type APIGroup string

// GoImportPrefix returns the path to use for this group in go import
func (o APIGroup) GoImportPrefix() string {
	parts := strings.Split(o.String(), ".")
	return parts[1] + "." + parts[0] + "/" + strings.Join(parts[2:], "/")
}

func (o APIGroup) String() string {
	return string(o)
}

// Replaces returns true if 'o' group is replaced by 'p' group
func (o APIGroup) Replaces(p APIGroup) bool {
	// * replaces extensions
	if o.String() != "extensions" && p.String() == "extensions" {
		return true
	}

	// core replaces events
	if o.String() == "" && p.String() == "events.k8s.io" {
		return true
	}

	return false
}
