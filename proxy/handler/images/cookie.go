package images

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	DefaultScreenWidth  = 1024
	DefaultScreenHeight = 768
	DefaultPixelRatio   = 1.0
)

type ClientInfo struct {
	screenWidth  int64
	screenHeight int64
	pixelRatio   float64
}

func NewClientInfo() *ClientInfo {
	info := new(ClientInfo)
	info.screenWidth = DEFAULT_SCREEN_WIDTH
	info.screenHeight = DEFAULT_SCREEN_HEIGHT
	info.pixelRatio = DEFAULT_PIXEL_RATIO
	return info
}

func ReadFoomoMediaClientInfo(cookie string) (clientInfo *ClientInfo, err error) {
	// screenWidthxscreenHeight@pixelRatio
	parts := strings.Split(cookie, "@")
	if len(parts) != 2 {
		err = errors.New(fmt.Sprint("could not separate screen size from pixel ratio", parts))
		return clientInfo, err
	} else {
		screenSizeParts := strings.Split(parts[0], "x")
		if len(screenSizeParts) != 2 {
			err = errors.New(fmt.Sprint("could not find screen size components ", len(screenSizeParts), " in ", parts[0]))
		} else {
			clientInfo = NewClientInfo()
			clientInfo.pixelRatio, _ = strconv.ParseFloat(parts[1], 32)
			clientInfo.screenWidth, _ = strconv.ParseInt(screenSizeParts[0], 0, 32)
			clientInfo.screenHeight, _ = strconv.ParseInt(screenSizeParts[1], 0, 32)
		}
		return clientInfo, err
	}
}

func ClampScreenWidthToGrid(screenWidth int64, breakPoints []int64) int64 {
	// the last breakpoint
	distance := breakPoints[len(breakPoints)-1]
	clampedValue := distance
	for _, breakPoint := range breakPoints {
		currentDistance := breakPoint - screenWidth
		if currentDistance < 0 {
			currentDistance *= -1
		}
		if screenWidth <= breakPoint && currentDistance < distance {
			distance = currentDistance
			clampedValue = breakPoint
		}
	}
	return clampedValue
}
