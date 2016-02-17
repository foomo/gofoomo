package images

import (
	"testing"
)

func TestReadFoomoMediaClientInfo(t *testing.T) {
	// 2560x1440@0.8999999761581421
	info, err := readFoomoMediaClientInfo("2560x1440@0.8999999761581421")
	if err != nil {
		t.Log("parse error", err)
		t.Fail()
	} else {
		if !(info.screenHeight == 1440 && info.screenWidth == 2560) {
			t.Error("wrong screen size")
		}
		pixelRatioDiff := info.pixelRatio - 0.8999999761581421
		if pixelRatioDiff > 0.0000000000000000000001 {
			t.Error("fuck those floats", pixelRatioDiff)
		}
	}
}

func TestClampScreenWidthToGrid(t *testing.T) {
	breakPoints := []int64{320, 768, 1024, 1440}
	assertSize := func(t *testing.T, expected int64, actual int64) {
		if expected != actual {
			t.Error("expected", expected, "actual", actual)
		}
	}
	assertSize(t, 320, clampScreenWidthToGrid(298, breakPoints))
	assertSize(t, 768, clampScreenWidthToGrid(321, breakPoints))
	assertSize(t, 768, clampScreenWidthToGrid(766, breakPoints))
	assertSize(t, 1024, clampScreenWidthToGrid(1024, breakPoints))
	assertSize(t, 1440, clampScreenWidthToGrid(1440, breakPoints))
	assertSize(t, 1440, clampScreenWidthToGrid(1441, breakPoints))
	assertSize(t, 1440, clampScreenWidthToGrid(10000, breakPoints))
}
