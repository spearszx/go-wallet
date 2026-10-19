package api

import "net/http"

type HttpFunc func(w http.ResponseWriter, r *http.Request) error
