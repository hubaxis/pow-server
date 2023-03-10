package main

import (
    "fmt"
    "github.com/hubaxis/pow-server/internal/config"
    "github.com/hubaxis/pow-server/internal/handler"
    "github.com/hubaxis/pow-server/internal/middleware"
    "github.com/hubaxis/pow-server/internal/repository"
    "github.com/hubaxis/pow-server/internal/service"

    "github.com/labstack/echo/v4"
    "context"
    "syscall"
    "os/signal"
    "os"
    log "github.com/sirupsen/logrus"
)

func main() {
    cfg, err := config.New()
    if err != nil {
        log.Fatalln(err)
    }

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGTERM)
    chuckNorrisRps:= repository.NewChackNoris(cfg.ChuckNorrisEndpoint)
    quoteService:= service.NewQuote(chuckNorrisRps)
    quoteHandler:= handler.NewQuote(quoteService)
    mw:= middleware.NewPow(ctx)
    e := echo.New()
    e.POST("/", quoteHandler.GetQuote, mw.Process)

    go func() {
    e.Logger.Fatal(e.Start(fmt.Sprintf(":%d",cfg.Port)))
    }()
        <-sigChan
        cancel()
        err= e.Close()
        if err != nil {
            log.Fatal(err)
        }

}