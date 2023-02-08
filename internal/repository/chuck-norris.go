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
    rsp, err := http.Get(r.endpoint)
    if err != nil {
       return "", fmt.Errorf("get %w", err)
    }
    //We Read the response body on the line below.
    body, err := io.ReadAll(rsp.Body)
    if err != nil {
        return "", fmt.Errorf("read all %w", err)
    }
    type response struct {
        Value string `json:"value"`
    }
    data:=&response{}
    err= json.Unmarshal(body, data)
    if err != nil {
        return "", fmt.Errorf("unmarshal %w", err)
    }
    return data.Value, nil
}