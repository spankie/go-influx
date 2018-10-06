package shops

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/fatih/structs"
	"github.com/go-chi/chi"
)

// GetAllProducts fetches all products in my shop
func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	t := template.New("index.html") // Create a template.
	cwd, _ := os.Getwd()
	p := path.Join(cwd, "public", "index.html")
	t, err := t.ParseFiles(p) // Parse template file.
	if err != nil {
		log.Println(err)
	}
	products := Product{}.GetAll() // Get current user infomration.
	err = t.Execute(w, products)   // merge.
	if err != nil {
		log.Println(err)
	}
}

// GetProduct handler for fetching a particular product
func GetProduct(w http.ResponseWriter, r *http.Request) {
	ID, err := strconv.Atoi(chi.URLParam(r, "ID"))
	if err != nil {
		log.Println(err)
		return
	}
	product, err := Product{}.Get(ID)
	if err != nil {
		log.Println(err)
	}
	t := template.New("product.html") // Create a template.
	cwd, _ := os.Getwd()
	p := path.Join(cwd, "public", "product.html")
	t, err = t.ParseFiles(p) // Parse template file.
	if err != nil {
		log.Println(err)
	}
	res, err := queryDB(fmt.Sprintf("select count(ProductID) from products where ProductID = %d", product.ID))
	if err != nil {
		log.Println(err)
	}
	result := res[0]
	if len(result.Series) > 0 {
		product.Views = string(result.Series[0].Values[0][1].(json.Number))
	} else {
		product.Views = "0"
	}
	err = t.Execute(w, product) // merge.
	if err != nil {
		log.Println(err)
	}

	pm := &ProductMeasurement{product.ID, product.Name, time.Now()}

	Insert(structs.Map(pm))
}
