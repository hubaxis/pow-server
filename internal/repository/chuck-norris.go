package repository

import (
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
)

type ChuckNorris struct {
    endpoint string
}

func NewChackNoris(endpoint string) *ChuckNorris {
    return &ChuckNorris{endpoint: endpoint}
}

func (r *ChuckNorris) Get(ctx context.Context) (string, error){
    resp, err := http.Get(r.endpoint)
    if err != nil {
       return "", fmt.Errorf("get %w", err)
    }
    //We Read the response body on the line below.
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", fmt.Errorf("read all %w", err)
    }
    type response struct {
        Value string `json:"value"`
    }
    rsp:=&response{}
    err= json.Unmarshal(body, rsp)
    if err != nil {
        return "", fmt.Errorf("unmarshal %w", err)
    }
    return rsp.Value, nil
}