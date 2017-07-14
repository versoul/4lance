package routes

import (
	"encoding/json"
	"net/http"
)

// Routes creates a REST router for the todos resource
/*func (rs dashboardResource) SettingsRoutes() chi.Router {
	r := chi.NewRouter()
	// r.Use() // some middleware..

	r.Post("/saveFilter", rs.saveFilter) // GET /todos - read a list of todos
	return r
}*/
type Data struct {
	Categories []string `json'categories,omitempty'`
	Keywords   []string `json'keywords,omitempty'`
}

func saveFilter(w http.ResponseWriter, r *http.Request) {
	sess, _ := globalSessions.SessionStart(w, r)
	defer sess.SessionRelease(w)
	data := Data{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		panic(err)
	}

	sess.Set("categories", data.Categories)
	w.Write([]byte("{}"))
}
func saveKeyWords(w http.ResponseWriter, r *http.Request) {
	sess, _ := globalSessions.SessionStart(w, r)
	defer sess.SessionRelease(w)
	data := Data{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		panic(err)
	}

	sess.Set("keywords", data.Keywords)
	w.Write([]byte("{}"))
}
