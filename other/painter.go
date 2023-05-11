package other

import (
	"github.com/fogleman/gg"
)

func DrawENDLeaderboard(rest *MarketStackRest, indexes []string, bgImagePath string, saveFilepath string) error {
	bg, err := gg.LoadPNG(bgImagePath)
	if err != nil {
		return err
	}

	imgWidth := bg.Bounds().Dx()
	imgHeight := bg.Bounds().Dy()

	dc := gg.NewContext(imgWidth, imgHeight)
	dc.DrawImage(bg, 0, 0)

	if err := dc.LoadFontFace("./assets/rubik.ttf", 43); err != nil {
		return err
	}

	var x, y float64 = 540, 255
	for _, index := range indexes {
		stock, err := rest.RequestEOD(index)
		if err != nil {
			return err
		}

		value, color := stock.DiffPaint()
		dc.SetColor(color)
		dc.DrawString(value, x, y)
		y += 105
	}

	return dc.SavePNG(saveFilepath)
}

func DrawYTDLeaderboard(rest *MarketStackRest, indexes []string, bgImagePath string, saveFilepath string) error {
	bg, err := gg.LoadPNG(bgImagePath)
	if err != nil {
		return err
	}

	imgWidth := bg.Bounds().Dx()
	imgHeight := bg.Bounds().Dy()

	dc := gg.NewContext(imgWidth, imgHeight)
	dc.DrawImage(bg, 0, 0)

	if err := dc.LoadFontFace("./assets/rubik.ttf", 43); err != nil {
		return err
	}

	var x, y float64 = 540, 255
	for _, index := range indexes {
		stock, err := rest.RequestYTD(index)
		if err != nil {
			return err
		}

		value, color := stock.DiffPaint()
		dc.SetColor(color)
		dc.DrawString(value, x, y)
		y += 105
	}

	return dc.SavePNG(saveFilepath)
}
