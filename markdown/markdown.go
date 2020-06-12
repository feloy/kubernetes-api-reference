package markdown

import "fmt"

// Code returns 's' as code
func Code(s string) string {
	return fmt.Sprintf("`%s`", s)
}
