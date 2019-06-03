package configuration

import (
	"regexp"
	"strings"
)
func ParseRegion(zoneUrls *[]string) []string {
	zones := make([]string, 0)
	for _, zone := range *zoneUrls {
		r := *regexp.MustCompile(`\/([^/?]*)$`)
		res := r.FindAllStringSubmatch(zone, -1)
		zones = append(zones, strings.Split(res[0][0], `/`)[1])
	}
	return zones
}