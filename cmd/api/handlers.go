package main

import (
	"net/http"
)

func ampInfoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	//result := ampStatus()
	//fmt.Fprint(w, result)
}
