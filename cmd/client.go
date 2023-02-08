package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "github.com/bwesterb/go-pow"
    "io"
    "net/http"
    "time"
)
type chalenge struct {
    ID string `json:"id"`
    Proof string `json:"proof"`
    Data string `json:"data"`
}

type auth struct {
    Id string `json:"id"`
    Chalenge string `json:"chalenge"`
}

func main()  {
    payload:="some bound data"
    for {
        rsp, err := http.Post("http://server:44444/", "application/json",bytes.NewBuffer([]byte("{}")))
        if err != nil {
            fmt.Println(fmt.Errorf("post %w", err))
            <- time.After(time.Second)
            continue
        }
        body, err := io.ReadAll(rsp.Body)
        if err != nil {
            fmt.Println( fmt.Errorf("read all %w", err))
            continue
        }
        fmt.Println(string(body))
        if rsp.StatusCode == http.StatusUnauthorized {
            //We Read the response body on the line below.
            data:= &auth{}
            err:= json.Unmarshal(body, data)
            if err != nil {
                fmt.Println( fmt.Errorf("unmarshal %w", err))
                continue
            }
            proof, err := pow.Fulfil(data.Chalenge, []byte(payload))
            if err != nil {
                fmt.Println( fmt.Errorf("fulfiil %w", err))
                continue
            }
            ch, err:=json.Marshal(&chalenge{
                ID: data.Id,
                Data: payload,
                Proof: proof,
            })
            rsp, err := http.Post("http://localhost:44444/", "application/json",bytes.NewBuffer(ch))
            if err != nil {
                fmt.Println(fmt.Errorf("post auth %w", err))
                continue
            }
            body, err := io.ReadAll(rsp.Body)
            if err != nil {
                fmt.Println( fmt.Errorf("read all %w", err))
                continue
            }
            fmt.Print(string(body))
        }
    }
}