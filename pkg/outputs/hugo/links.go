package hugo

import (
	"fmt"
	"strings"
)

// LinkEnd returns a link to a section in a part/chapter
// s is an array containing partname / chaptername
func (o *Hugo) LinkEnd(s []string, name string) string {
	typename := name
	array := ""
	if strings.HasPrefix(name, "[]") {
		array = "[]"
		typename = strings.TrimPrefix(name, array)
	}
	return fmt.Sprintf("%s<a href=\"{{< ref \"/docs/%s/%s#%s\" >}}\">%s</a>", array, escapeName(s[0]), escapeName(s[1]), escapeName(typename), typename)
}
