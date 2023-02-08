package middleware

import (
    "context"
    "encoding/json"
    "github.com/labstack/echo/v4"
    "io"
    "net/http"
    "sync"
    "time"
    "github.com/bwesterb/go-pow"
    "github.com/google/uuid"
)

type data struct {
    Task string
    Expire time.Time
}

type chalenge struct {
    ID string `json:"id"`
    Proof string `json:"proof"`
    Data string `json:"data"`
}

type auth struct {
    Id string `json:"id"`
    Chalenge string `json:"chalenge"`
}

type Pow struct {
    tasks     map[string]*data
    tasksMu        sync.RWMutex
}


func NewPow(ctx context.Context) *Pow {
    p:= &Pow{tasks: map[string]*data{}}

    go func(ctx context.Context) {
        timer:= time.NewTicker(time.Minute)
        for {
            select {
            case <-ctx.Done():
                return
            case <-timer.C:
                keys:=[]string{}
                p.tasksMu.RLock()
                for k,v:= range p.tasks{
                    if v.Expire.Before(time.Now()){
                       keys = append(keys, k)
                    }
                }
                p.tasksMu.RUnlock()
                p.tasksMu.Lock()
                for i:= range keys {
                    delete(p.tasks, keys[i])
                }
                p.tasksMu.Unlock()
            }
        }
    }(ctx)

    return p
}

func (s *Pow) makeChalenge() *auth  {
    id:=uuid.New().String()
    ch:=pow.NewRequest(30,[]byte(id))
    s.tasksMu.Lock()
    s.tasks[id]= &data{Task: ch, Expire: time.Now().Add(time.Hour)}
    s.tasksMu.Unlock()
    return &auth{Id: id, Chalenge: ch}
}

// Process is the middleware function.
func (s *Pow) Process(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        rqs:=c.Request()
        //We Read the response body on the line below.
        body, err := io.ReadAll(rqs.Body)
        if err != nil {
            return echo.NewHTTPError(http.StatusBadRequest, err)
        }
        ch:=&chalenge{}
        err= json.Unmarshal(body, ch)
        if err != nil {
            return echo.NewHTTPError(http.StatusBadRequest, err)
        }
        s.tasksMu.RLock()
        d, ok:= s.tasks[ch.ID]
        s.tasksMu.RUnlock()
        if !ok {
            return c.JSON(http.StatusUnauthorized, s.makeChalenge())
        }
        if d.Expire.Before(time.Now()) {
            return c.JSON(http.StatusUnauthorized, s.makeChalenge())
        }

        ok, err =pow.Check(d.Task,  ch.Proof, []byte(ch.Data))
        if err!=nil{
            return echo.NewHTTPError(http.StatusBadRequest, err)
        }
        if ok {
            s.tasksMu.Lock()
            delete(s.tasks, ch.ID)
            s.tasksMu.Unlock()
            return next(c)
        } else {
            return echo.NewHTTPError(http.StatusBadRequest, "incorect data")
        }
    }
}