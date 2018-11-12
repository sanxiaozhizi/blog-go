package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"github.com/sanxiaozhizi/blog-go/model"

	"github.com/go-chi/chi"

	_ "github.com/sanxiaozhizi/blog-go/model"
)

var (
	config struct {
		Addr      string
		AccessKey string
	}
)

func init() {
	bytes, err := ioutil.ReadFile(filepath.Join(model.ROOT, "config/config.yaml"))
	if err != nil {
		panic(err)
	}
	if err := yaml.Unmarshal(bytes, &config); err != nil {
		panic(err)
	}
}

func main() {
	r := chi.NewRouter()
	r.Post("/graphql", func(w http.ResponseWriter, r *http.Request) {
		query := r.PostFormValue("query")
		if query == "" {
			r.ParseForm()
			query = r.PostFormValue("query")
		}
		json.NewEncoder(w).Encode(model.ExecuteQuery(query))
	})
	r.Post("/update", func(w http.ResponseWriter, r *http.Request) {
		if r.FormValue("access_key") != config.AccessKey {
			w.WriteHeader(401)
			return
		}
		model.PostLoading()
		model.LinkLoading()
		model.AboutLoading()
	})
	r.Get("/about", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(model.About)
	})
	http.ListenAndServe(config.Addr, r)
}
