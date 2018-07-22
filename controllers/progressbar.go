package controllers

import (
	"github.com/labstack/echo"
	"net/http"
	"strings"
	"strconv"
)

type Progress struct {
	Black rune
	White rune
}

func (p *Progress) toBar(percentage, size int) string {
	n := size * percentage / 100;
	return strings.Repeat(string(p.Black), n) + strings.Repeat(string(p.White), size-n)
}

const (
	dot = iota
	rectangle1
	rectangle2
	squareEmoji
	circleEmoji
)

func createProgress(pType int) Progress {
	switch pType {
	case dot:
		return Progress{Black: '⣿', White: '⣀'}
	case rectangle1:
		return Progress{Black: '█', White: '▁'}
	case rectangle2:
		return Progress{Black: '█', White: '▒'}
	case squareEmoji:
		return Progress{Black: '\u2B1B', White: '\u2B1C'}
	case circleEmoji:
		return Progress{Black: '\u26AB', White: '\u26AA'}
	}
	return Progress{}
}

func ProgressBar(c echo.Context) error {
	percentage, err := strconv.Atoi(c.Param("percentage"))
	pType := 1
	if len(c.QueryParam("type")) != 0 {
		pType, err = strconv.Atoi(c.QueryParam("type"))
		if err != nil {
			return c.String(http.StatusBadRequest, "")
		}
	}
	progress := createProgress(pType)
	return c.String(http.StatusOK, progress.toBar(percentage, 13))
}
