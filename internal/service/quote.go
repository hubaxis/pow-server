package service

import (
    "context"
    "fmt"
)

type QuoteProvider interface {
    Get(ctx context.Context) (string, error)
}

type Quote struct {
    provider QuoteProvider
}

func NewQuote(provider QuoteProvider) *Quote {
    return &Quote{provider: provider}
}

func (s *Quote) Get(ctx context.Context) (string, error){
   q, err:= s.provider.Get(ctx)
   if err!=nil {
       return "", fmt.Errorf("get %w", err)
   }
   return q, nil
}