package markdown

import "fmt"

// Code returns 's' as code
func Code(s string) string {
	return fmt.Sprintf("`%s`", s)
}

// Chapter returns a Level 2 mark
func Chapter(name string) string {
	return fmt.Sprintf("## %s\n", name)
}

// Section returns a Level 3 mark
func Section(name string) string {
	return fmt.Sprintf("### %s\n", name)
}

// ListEntry returns a list entry
func ListEntry(title string, content string) string {
	return fmt.Sprintf("- %s\n  %s\n", title, content)
}
