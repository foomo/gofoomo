package images

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const (
	DefaultScreenWidth             = 1024
	DefaultScreenHeight            = 768
	DefaultPixelRatio              = 1.0
	FoomoMediaClientInfoCookieName = "foomoMediaClientInfo"
)

type ClientInfo struct {
	screenWidth  int64
	screenHeight int64
	pixelRatio   float64
}

func NewClientInfo() *ClientInfo {
	info := new(ClientInfo)
	info.screenWidth = DefaultScreenWidth
	info.screenHeight = DefaultScreenHeight
	info.pixelRatio = DefaultPixelRatio
	return info
}

func readFoomoMediaClientInfo(cookie string) (clientInfo *ClientInfo, err error) {
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

func clampScreenWidthToGrid(screenWidth int64, breakPoints []int64) int64 {
	// the last breakpoint
	log.Println("clampScreenWidthToGrid", screenWidth, breakPoints)
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

func getFoomoMediaClientInfoCookie(incomingCookies []*http.Cookie, breakPoints []int64) *http.Cookie {
	clientInfo := NewClientInfo()
	clientInfoCookie := getCookieByName(incomingCookies, FoomoMediaClientInfoCookieName)
	if clientInfoCookie != nil {
		cookieClientInfo, cookieReadError := readFoomoMediaClientInfo(clientInfoCookie.Value)
		if cookieReadError == nil {
			clientInfo = cookieClientInfo
		}
	}
	pixelRatio := clientInfo.pixelRatio
	if pixelRatio > 1.5 {
		pixelRatio = 2.0
	} else {
		pixelRatio = 1.0
	}
	cookieValue := fmt.Sprintf("%dx%d@%f", clampScreenWidthToGrid(clientInfo.screenWidth, breakPoints), 1000, pixelRatio)
	cookie := &http.Cookie{
		Name:  FoomoMediaClientInfoCookieName,
		Value: cookieValue,
	}
	return cookie
}

func getCookieByName(cookies []*http.Cookie, name string) (cookie *http.Cookie) {
	for _, cookie := range cookies {
		if cookie.Name == name {
			return cookie
		}
	}
	return nil
}
