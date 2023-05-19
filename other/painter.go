package other

import (
	"github.com/fogleman/gg"
)

func DrawUSEODLeaderboard(rest *MarketStackRest, filepath string) error {
	bg, err := gg.LoadPNG("./assets/us-eod.png")
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

	indexes := []string{"IXIC", "GSPC", "DJI"}
	var x, y float64 = 520, 255
	for _, index := range indexes {
		stock, err := rest.RequestEOD(index)
		if err != nil {
			return err
		}

		value, color := stock.DiffEODPaint()
		dc.SetColor(color)
		dc.DrawString(value, x, y)
		y += 105
	}

	return dc.SavePNG(filepath)
}

func DrawUSYTDLeaderboard(rest *MarketStackRest, filepath string) error {
	bg, err := gg.LoadPNG("./assets/us-ytd.png")
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

	indexes := []string{"IXIC", "GSPC", "DJI"}
	var x, y float64 = 520, 255
	for _, index := range indexes {
		stock, err := rest.RequestYearly(index)
		if err != nil {
			return err
		}

		value, color := stock.DiffYTDPaint()
		dc.SetColor(color)
		dc.DrawString(value, x, y)
		y += 105
	}

	return dc.SavePNG(filepath)
}

func DrawInternationalYTDLeaderboard(rest *MarketStackRest, filepath string) error {
	bg, err := gg.LoadPNG("./assets/international-ytd.png")
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

	indexes := []string{"DAXI", "FTSE", "FCHI", "000001.SS", "WSI", "N225"}
	var x, y float64 = 1000, 275
	for pos, index := range indexes {
		stock, err := rest.RequestYearly(index)
		if err != nil {
			return err
		}

		value, color := stock.DiffYTDPaint()
		dc.SetColor(color)
		dc.DrawString(value, x, y)
		y += 120 + float64(pos)
	}

	return dc.SavePNG(filepath)
}
