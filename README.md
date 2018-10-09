
# InfluxDB with Go

InfluxDB is one of the popular time series database out there. In this article, I’ll try to explain some of the key concepts needed to get you up-and-running with InfluxDB. This article assumes that you have InfluxDB already installed on your machine. if you don't, you can install it by following [this instructions.](https://influxdbcom.readthedocs.io/en/latest/content/site/download/)

## Primer

A time series database should not be confused with relational databases as they are totally different. A time series is basically a series of points indexed in time order. Examples of Time Series Data include:

* Measuring the level of unemployment each month of the year 

* Tracking the product that sold more at a particular time in an e-commerce site.

* Tracking the most visited product in an e-commerce website

* Price of stocks and shares taken at a regular interval.

Essentially, the root of everything is time in InfluxDB, as all columns have a special time column.

### Key Concepts

InfluxDB doesn’t impose any schema structure on the data you are going to store. Developers that are already familiar with SQL can quickly pick up InfluxDB as it provides an SQL-like syntax. This has particularly led to an increase adoption of influxdb. 

* **Database**— InfluxDB has a notion of database

* **Measurement**— Is a description of the data that are stored. It is similar to a table in SQL.

* **Fields**—Comprised of fieldKeys and fieldValues. Fields are required and are not indexed. This means that Queries using Fields to filter data are relatively slow.

* **Tags**— Comprised of tagKeys and tagValues. Tags are optional and can be indexed. This means that Queries using Tags are much faster than their Fields counterpart.

In this article We will be building a simple shop to display our cool products and take a measurement of how much visits each product get. we need to install the following;

#### **InfluxDB Go Client Library**

InfluxDB is written completely with Go with no external dependencies and it makes sense that the team provides a [Go client](https://github.com/influxdata/influxdb) to make interfacing with InfluxDB easy from a Go application.


#### **Graphana**

Graphana is an open-source, general purpose dashboard and graph composer, which runs as a web application. It is a Data visualization & Monitoring platform and it supports graphite, InfluxDB or opentsdb as backends

N/B: This article assumes you already have **Golang** installed on your OS, if you don't you can follow the instructions [here](https://golang.org/dl/)

#### Other Dependencies:
- **[go-chi chi router](https://github.com/go-chi/chi)** This is a lightweight, idiomatic and composable router for building Go HTTP services.

- [**InfluxDB Go Client**](github.com/influxdata/influxdb/client/v2)

- **UI Library [tachyons.io](http://tachyons.io)** is a library for creating fast loading, highly readable, and 100% responsive interfaces with as little css as possible.

Now we have our dependencies set let's write some code. Thankfully, getting a web application up and running in golang is very simple! All we need is a server and some html to serve right?

Our final website folder structure will look like this:
```
.
+-- modules
|   +-- shops
|       +-- db.go
|       +-- handler.go
|       +-- model.go
+-- public
|   +-- assets
|       +-- css
|           +-- tachyons.min.css
|       +-- img
|           +-- camera.jpeg
|           +-- glass.jpeg
|           +-- toy.jpeg
|           +-- watch.jpeg
|   +-- index.html
|   +-- products.html
|   +-- 404.html
+-- main.go
+-- README.md
+-- .gitignore
```

Lets create a folder called **GO__INFLUX_APP**
First we create a **public** folder and add [this files](https://github.com/spankie/go-influx/tree/develop/public) to it.

Next Create a Database. In your terminal type:
```
 step1: influx //to enter influx DB mode
 step2: create database go_influx //to create the db
```
Inside the folder create a file called main.go in the root directory and add the following lines of code to it.
```
package main

import (
	"net/http"
	"github.com/go-chi/chi"                // http router
	"github.com/go-chi/chi/middleware"    // package for commonly used middlewares
	log "github.com/sirupsen/logrus"     // cool logger package
	"github.com/spankie/go-influx/handlers" // handler packages
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	fileServer := http.StripPrefix("/assets/", http.FileServer(http.Dir("./public/assets")))
	r.Get("/assets/*", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Vary", "Accept-Encoding")
		w.Header().Set("Cache-Control", "public, max-age=7776000")
		fileServer.ServeHTTP(w, r)
	})

	// Define App routes
	
	//Homepage
	r.Get("/", shops.GetAllProducts)
	
	// Product page handler
	r.Get("/product/{ID}", shops.GetProduct)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./public/404.html")
	})

	log.Print("👉  Server started at 127.0.0.1:3333")
	http.ListenAndServe(":3333", r)
}
```
Here we have defined three routes to display the home page with the listing of all the products (/), the product details page (/product/{id}), where ID is the id of the product we want to view. The last route is any route not explicitly defined which will render a 404 page to the user.

N/B: We will use an array of static products for this article.

Create a folder called **modules/shops** and add a file called **model.go** (which contains the array of products we will use and other methods):
```
package shops

import (
	"errors"
	"log"
	"time"
)

// Product structure of my cool shop product
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Price       float32 `json:"price"`
	Image       string  `json:"image"`
	Description string  `json:"description"`
	Views       string
}

// ProductMeasurement schema of product measurement
type ProductMeasurement struct {
	ProductID   int `json:"id"`
	ProductName string
	time        time.Time
}

var allProducts = []Product{
	Product{ID: 0, Name: "Watch", Price: 56.2, Image: "/assets/img/watch.jpeg", Description: "Nice Watch"},
	Product{ID: 1, Name: "Camera", Price: 34.2, Image: "/assets/img/camera.jpeg", Description: "Nice Camera"},
	Product{ID: 2, Name: "Glass", Price: 24.2, Image: "/assets/img/glass.jpeg", Description: "Nice Glass"},
	Product{ID: 3, Name: "Toy", Price: 56.2, Image: "/assets/img/toy.jpeg", Description: "Nice Toy"},
}

// GetAll fetches all products in my shop
func (product Product) GetAll() []Product {
	return allProducts
}

// Get fetches a particular product in my shop identified By its ID
func (product Product) Get(id int) (Product, error) {
	if id < 0 && len(allProducts) >= id {
		log.Println(id)
		return Product{}, errors.New("No product found")
	}
	return allProducts[id], nil
}
```
we need to create a file called DB.go because => (This is to abstract all Database access from the other parts of the app)

The `Insert` function takes a productMeasurement type and saves it in the Database. While the queryDB function is a conevinience function that receives string query, executes the query and return the result.
```
package shops

import (
	"log"
	"time"

	"github.com/influxdata/influxdb/client/v2"
)

const (
	// MyDB spcifies name of database
	MyDB = "go_influx"
)

// Insert saves points to database
func Insert(productMeasurement map[string]interface{}) {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: "http://localhost:8086",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  MyDB,
		Precision: "s",
	})
	if err != nil {
		log.Fatal(err)
	}

	// Create a point and add to batch
	tags := map[string]string{"productView": productMeasurement["ProductName"].(string)}
	fields := productMeasurement

	pt, err := client.NewPoint("products", tags, fields, time.Now())
	if err != nil {
		log.Fatal(err)
	}
	bp.AddPoint(pt)

	// Write the batch
	if err := c.Write(bp); err != nil {
		log.Fatal(err)
	}

	// Close client resources
	if err := c.Close(); err != nil {
		log.Fatal(err)
	}
}

// queryDB convenience function to query the database
func queryDB(cmd string) (res []client.Result, err error) {
	q := client.Query{
		Command:  cmd,
		Database: MyDB,
	}
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: "http://localhost:8086",
	})
	if response, err := c.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	} else {
		return res, err
	}
	return res, nil
}
```
Now lets analyze each function:

The `Insert()` function creates a client to connect to the database at `localhost` port `8086` which is the default port for the influxDB server

Next a point batch is created to store points to be saved to the DB. A point is then created to store the ProductMeasurement details which includes the time of visit, the id of the product as well as the Name of the product being viewed.

Here is how the DB looks like after insert:

![db.png](https://raw.githubusercontent.com/spankie/go-influx/develop/public/screenshots/db.png)


The time field is indexed and stores the exact time the product was viewed. This time can be used to query the DB for analysis.

The `queryDB()` function queries the influxDB and return the result as an array of client.Result type.

Finally lets create a file called `handler.go` and add the following lines of code to it.

**INSERT REASON FOR HANDLER.GO**
Each handler function will perform a specific task required for rendering a page to the route being accessed. e.g GetAllProducts() gets all products and renders it in the index.html template.

```
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

```

Now lets see what we have so far; To Start your app, in your terminal type:

```
go run main.go
```

### Included here is the screenshot of what it looks like when you start the browser.

Visit `localhost:3333` on your browser, and you should see this:


![index.html](https://raw.githubusercontent.com/spankie/go-influx/develop/public/screenshots/screenshot1.png)


Clicking on a product should take you to a page that looks like this:

![product.html](https://raw.githubusercontent.com/spankie/go-influx/develop/public/screenshots/screenshot2.png)


We need a way to graphically view this data we just saved in our awesome influxDB

These views can be analyzed/viewed using an awesome tool called [grafana](http://docs.grafana.org/).

After installing grafana from [here](http://docs.grafana.org/), set it up by following their straightfowarad guide on their website. Then do the following:
	
- Create a new dashboard with a graph panel

- Make a query as specified in this screenshot below:
	![query](https://raw.githubusercontent.com/spankie/go-influx/develop/public/screenshots/query.png)

- The axes tab should look like this:
	![axes](https://raw.githubusercontent.com/spankie/go-influx/develop/public/screenshots/axes.png)

- Display tab should also look like this:
	![display](https://raw.githubusercontent.com/spankie/go-influx/develop/public/screenshots/display.png)

- Your graph should finally look like this:
	![display](https://raw.githubusercontent.com/spankie/go-influx/develop/public/screenshots/graph.png)

Screeshots:

![index.html](https://raw.githubusercontent.com/spankie/go-influx/develop/public/screenshots/screenshot1.png)

![product.html](https://raw.githubusercontent.com/spankie/go-influx/develop/public/screenshots/screenshot2.png)


# Conclusion

_In this article we looked at how to work with InfluxDB to monitor time series data, I hope you have learnt alot to start building awesome stuff. Hack ON!_