package kubernetes

// APIGroup represents the group of a Kubernetes API
type APIGroup string

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
