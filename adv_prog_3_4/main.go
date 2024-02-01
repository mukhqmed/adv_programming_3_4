package main

import (
	"database/sql"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Barber struct {
	ID         int
	Name       string
	BasicInfo  string
	Price      int
	Experience string
	Status     string
	ImagePath  string
}

var db *sql.DB
var tpl *template.Template

func init() {
	var err error
	db, err = sql.Open("postgres", "postgres://postgres:1234@localhost/adv_prog?sslmode=disable")
	if err != nil {
		panic(err)
	}

	tpl = template.Must(template.ParseGlob("templates/*.html"))
}

func main() {
	router := gin.Default()

	router.LoadHTMLGlob("templates/*.html")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	router.GET("/barbers", func(c *gin.Context) {
		barbers := getBarbersFromDB()

		c.HTML(http.StatusOK, "barbers.html", gin.H{
			"Barbers": barbers,
		})
	})

	router.GET("/filtered-barbers", func(c *gin.Context) {
		statusFilter := c.Query("status")
		experienceFilter := c.Query("experience")
		sortBy := c.Query("sort")
		pageStr := c.Query("page")
		itemsPerPage := 3

		barbers := getFilteredBarbersFromDB(statusFilter, experienceFilter, sortBy, pageStr, itemsPerPage)

		c.HTML(http.StatusOK, "barbers.html", gin.H{
			"Barbers": barbers,
		})
	})

	router.Run(":8080")
}

func getBarbersFromDB() []Barber {
	rows, err := db.Query("SELECT id, name, basic_info, price, experience, status, image_path FROM barbers")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var barbers []Barber
	for rows.Next() {
		var b Barber
		err := rows.Scan(&b.ID, &b.Name, &b.BasicInfo, &b.Price, &b.Experience, &b.Status, &b.ImagePath)
		if err != nil {
			panic(err)
		}
		barbers = append(barbers, b)
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}

	return barbers
}

func getFilteredBarbersFromDB(statusFilter, experienceFilter, sortBy, pageStr string, itemsPerPage int) []Barber {
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	// Construct SQL query based on filter, sort, and pagination parameters
	query := "SELECT id, name, basic_info, price, experience, status, image_path FROM barbers WHERE true"
	if statusFilter != "" {
		query += " AND status = '" + statusFilter + "'"
	}
	if experienceFilter != "" {
		query += " AND experience = '" + experienceFilter + "'"
	}
	switch sortBy {
	case "name":
		query += " ORDER BY name"
	case "price":
		query += " ORDER BY price"
	}
	query += " LIMIT " + strconv.Itoa(itemsPerPage) + " OFFSET " + strconv.Itoa((page-1)*itemsPerPage)

	// Execute SQL query and parse results...

	// Execute SQL query
	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	// Parse query results into []Barber
	var barbers []Barber
	for rows.Next() {
		var b Barber
		err := rows.Scan(&b.ID, &b.Name, &b.BasicInfo, &b.Price, &b.Experience, &b.Status, &b.ImagePath)
		if err != nil {
			panic(err)
		}
		barbers = append(barbers, b)
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}

	return barbers
}
