package main

import "net/http"

func (apiCfg *apiConfig) middleWareToIncreaseHits(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		apiCfg.hits++
		next.ServeHTTP(res, req)
	})
}
