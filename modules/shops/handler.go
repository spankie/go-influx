package shops

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/go-chi/chi"
)

func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	t := template.New("index.html") // Create a template.
	cwd, _ := os.Getwd()
	p := path.Join(cwd, "public", "index.html")
	t, err := t.ParseFiles(p) // Parse template file.
	if err != nil {
		log.Println(err)
	}
	shops := Shop{}.GetAll()  // Get current user infomration.
	err = t.Execute(w, shops) // merge.
	if err != nil {
		log.Println(err)
	}
}

func GetProduct(w http.ResponseWriter, r *http.Request) {
	shops := Shop{}.GetAll()
	ID, err := strconv.Atoi(chi.URLParam(r, "ID"))
	if err != nil {
		log.Println(err)
		return
	}
	if ID < 0 && len(shops) >= ID {
		log.Println(ID)
	}
	t := template.New("product.html") // Create a template.
	cwd, _ := os.Getwd()
	p := path.Join(cwd, "public", "product.html")
	t, err = t.ParseFiles(p) // Parse template file.
	if err != nil {
		log.Println(err)
	}
	err = t.Execute(w, shops[ID]) // merge.
	if err != nil {
		log.Println(err)
	}
}
