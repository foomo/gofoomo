package images

import (
	"github.com/foomo/gofoomo/foomo"
	"github.com/foomo/gofoomo/foomo/core"
	"sort"
	"strconv"
)

type MediaServerConfig struct {
	Grid map[string]map[string]int64 `json:"grid"`
}

type SessionConfig struct {
	Name string `json:"name"`
}

func getBreakPoints(f *foomo.Foomo) []int64 {
	c := &MediaServerConfig{}

	err := core.GetConfig(f, &c, "Foomo.Media", "Foomo.Media.Image.server", "")
	if err != nil {
		panic(err)
	}
	breakPointsInt := []int{}

	for breakPointString, _ := range c.Grid {
		breakPointInt, _ := strconv.Atoi(breakPointString)
		if breakPointInt > 0 {
			breakPointsInt = append(breakPointsInt, breakPointInt)
		}
	}
	sort.Ints(breakPointsInt)
	breakPoints := []int64{}

	for _, breakPointInt := range breakPointsInt {
		breakPoints = append(breakPoints, int64(breakPointInt))
	}

	return breakPoints
}

func getFoomoSessionCookieName(f *foomo.Foomo) string {
	c := &SessionConfig{}

	core.GetConfig(f, &c, "Foomo", "Foomo.session", "")

	return c.Name
}
