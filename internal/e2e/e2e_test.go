package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/bwesterb/go-pow"
	"github.com/hubaxis/pow-server/internal/config"
	"github.com/hubaxis/pow-server/internal/handler"
	"github.com/hubaxis/pow-server/internal/middleware"
	"github.com/hubaxis/pow-server/internal/repository"
	"github.com/hubaxis/pow-server/internal/service"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

type chalenge struct {
	ID    string `json:"id"`
	Proof string `json:"proof"`
	Data  string `json:"data"`
}

type auth struct {
	Id       string `json:"id"`
	Chalenge string `json:"chalenge"`
}

func TestMain(m *testing.M) {
	cfg, err := config.New()
	if err != nil {
		log.Fatalln(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	chuckNorrisRps := repository.NewChackNoris(cfg.ChuckNorrisEndpoint)
	quoteService := service.NewQuote(chuckNorrisRps)
	quoteHandler := handler.NewQuote(quoteService)
	mw := middleware.NewPow(ctx)
	e := echo.New()
	e.POST("/", quoteHandler.GetQuote, mw.Process)

	go func() {
		e.Start(fmt.Sprintf(":%d", cfg.Port))
	}()
	code := m.Run()
	cancel()
	err = e.Close()
	log.Error(err)
	os.Exit(code)
}

func TestQuotes(t *testing.T) {
	payload := "some bound data"
	rsp, err := http.Post("http://127.0.0.1:44444/", "application/json", bytes.NewBuffer([]byte("{}")))
	require.NoError(t, err)
	body, err := io.ReadAll(rsp.Body)
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, rsp.StatusCode)
	//We Read the response body on the line below.
	data := &auth{}
	err = json.Unmarshal(body, data)
	require.NoError(t, err)
	proof, err := pow.Fulfil(data.Chalenge, []byte(payload))
	require.NoError(t, err)
	ch, err := json.Marshal(&chalenge{
		ID:    data.Id,
		Data:  payload,
		Proof: proof,
	})
	rsp, err = http.Post("http://localhost:44444/", "application/json", bytes.NewBuffer(ch))
	require.NoError(t, err)
	body, err = io.ReadAll(rsp.Body)
	require.NoError(t, err)
	require.Equal(t, rsp.StatusCode, 200)

}
