package utils

import "regexp"

var SafeName = regexp.MustCompile(`[^\.\\]`)

func IsNameSafe(name string) bool {
  return SafeName.MatchString(name)
}