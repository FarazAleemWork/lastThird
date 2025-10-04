package main

import (
	"encoding/json"
	"fmt"
	calculate "lastThird/prayertimecalc"
	"log"
	"net/http"
	"strings"
)

func apiGeocode(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	city := q.Get("city")
	state := q.Get("state")
	country := q.Get("country")
	timezone := q.Get("timezone")

	missing := []string{}
	if city == "" {
		missing = append(missing, "city")
	}
	if state == "" {
		missing = append(missing, "state")
	}
	if country == "" {
		missing = append(missing, "country")
	}
	if timezone == "" {
		missing = append(missing, "timezone")
	}

	if len(missing) > 0 {
		msg := fmt.Sprintf(`{"error":"missing required parameter(s): %v"}`, strings.Join(missing, ", "))
		http.Error(w, msg, http.StatusBadRequest)
		log.Printf("Bad request: %v", msg)
		return
	}

	tahajjudStart, err := calculate.GetTahajjud(city, state, country, timezone)
	if err != nil {
		http.Error(w, `{"error":"internal error computing tahajjud"}`, http.StatusInternalServerError)
		log.Println("GetTahajjud error: ", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if err := enc.Encode(map[string]interface{}{"Tahajjud starts at": tahajjudStart}); err != nil {
		http.Error(w, `{"error":"failed to encode response"}`, http.StatusInternalServerError)
		return
	}
}

// remove /app/ when running locally or i guess i can mkdir lastthird in docker
func main() {
	fs := http.FileServer(http.Dir("/app/lastthirdapp/frontend"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Serving index.html for", r.URL.Path)
		http.ServeFile(w, r, "/app/lastthirdapp/frontend/index.html")
	})

	http.HandleFunc("/api/geocode", apiGeocode)

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
