package configuration

import (
	"context"
	"encoding/json"
	"log"
	"regexp"
	"strings"

	"google.golang.org/api/compute/v1"
	_ "google.golang.org/api/compute/v1"
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

func GetZones(ctx context.Context, comService *compute.Service, project string, region string) []string {
	resp, err := comService.Regions.Get(project, region).Context(ctx).Do()
	if err != nil {
		log.Fatal(err)
	}
	jsonRegions, _ := json.Marshal(resp)
	var regions Region
	if unmarshallError := json.Unmarshal(jsonRegions, &regions); unmarshallError != nil {
		log.Fatal(err)
		return nil
	}
	return ParseRegion(&regions.Zones)
}
