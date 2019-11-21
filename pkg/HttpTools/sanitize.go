package HttpTools

import "github.com/microcosm-cc/bluemonday"

func Sanitizer (s [](*string)){
	sanitizer := bluemonday.UGCPolicy()
	for _, st := range s{
		*st = sanitizer.Sanitize(*st)
	}
}
