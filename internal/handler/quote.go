package handler

import (
    "github.com/hubaxis/pow-server/internal/service"
    "github.com/labstack/echo/v4"
    "github.com/labstack/gommon/log"
    "net/http"
)

type Quote struct {
    quoteService *service.Quote
}

func NewQuote(quoteService *service.Quote) *Quote  {
    return &Quote{quoteService: quoteService}
}

func (h *Quote) Get(c echo.Context) error {
    q, err:= h.quoteService.Get(c.Request().Context())
    if err!=nil{
        log.Errorf("error getting quote %s", err)
        return echo.NewHTTPError(http.StatusBadRequest, err)
    }
    return c.JSON(http.StatusOK, q)
}
