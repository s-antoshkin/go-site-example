package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"

	"github.com/jackc/pgx/v5"
)

type Rsvp struct {
	Name, Email, Phone string
	WillAttend         bool
}

var responses = make([]*Rsvp, 0, 10)
var templates = make(map[string]*template.Template, 3)

func loadTemplates() {
	templateNames := [5]string{"welcome", "form", "thanks", "sorry", "list"}
	for index, name := range templateNames {
		t, err := template.ParseFiles("./templates/base.html", "./templates/"+name+".html")
		if err == nil {
			templates[name] = t
			fmt.Println("Loaded template", index, name)
		} else {
			panic(err)
		}
	}
}

func welcomeHandler(writer http.ResponseWriter, request *http.Request) {
	templates["welcome"].Execute(writer, nil)
}

type formData struct {
	*Rsvp
	Errors []string
}

func checkDB(email string) bool {
	var count int
	connURL := "postgres://postgres:postgres@localhost:5432/rsvp_db"
	conn, err := pgx.Connect(context.Background(), connURL)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer conn.Close(context.Background())

	row := conn.QueryRow(context.Background(), "SELECT COUNT(*) FROM rsvp WHERE email=$1", email)
	err = row.Scan(&count)
	if err != nil {
		fmt.Println(err)
		return false
	}
	if count > 0 {
		return true
	}

	return false

}

func updadeRecord(conn *pgx.Conn, data *Rsvp) {
	query := "UPDATE rsvp SET name=$1, phone=$2, will_attend=$3 WHERE email=$4"
	_, err := conn.Exec(context.Background(), query, data.Name, data.Phone, data.WillAttend, data.Email)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func insertRecord(data *Rsvp) {
	connURL := "postgres://postgres:postgres@localhost:5432/rsvp_db"
	conn, err := pgx.Connect(context.Background(), connURL)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close(context.Background())

	if checkDB(data.Email) {
		updadeRecord(conn, data)
	} else {
		query := "INSERT INTO rsvp (name, email, phone, will_attend) VALUES ($1, $2, $3, $4)"
		_, err := conn.Query(context.Background(), query, data.Name, data.Email, data.Phone, data.WillAttend)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func formHandler(writer http.ResponseWriter, request *http.Request) {

	if request.Method == http.MethodGet {
		templates["form"].Execute(writer, formData{
			Rsvp:   &Rsvp{},
			Errors: []string{},
		})
	} else if request.Method == http.MethodPost {
		request.ParseForm()
		responseData := Rsvp{
			Name:       request.Form["name"][0],
			Email:      request.Form["email"][0],
			Phone:      request.Form["phone"][0],
			WillAttend: request.Form["willattend"][0] == "true",
		}

		errors := []string{}
		if responseData.Name == "" {
			errors = append(errors, "Пожалуйста, введите Ваше имя")
		}

		if responseData.Email == "" {
			errors = append(errors, "Пожалуйста, введите Ваш адрес электронной почты")
		}

		if responseData.Phone == "" {
			errors = append(errors, "Пожалуйста, введите Ваш номер телефона")
		}

		if len(errors) > 0 {
			templates["form"].Execute(writer, formData{
				Rsvp: &responseData, Errors: errors,
			})
		} else {
			insertRecord(&responseData)

			if responseData.WillAttend {
				templates["thanks"].Execute(writer, responseData.Name)
			} else {
				templates["sorry"].Execute(writer, responseData.Name)
			}
		}
	}
}

type Scanner interface {
	Scan(...interface{}) error
}

func (r *Rsvp) Scan(row Scanner) error {
	err := row.Scan(
		&r.Name,
		&r.Phone,
		&r.Email,
		&r.WillAttend,
	)
	if err != nil {
		return err
	}

	return nil
}

func listHandler(writer http.ResponseWriter, request *http.Request) {
	responses = nil
	connURL := "postgres://postgres:postgres@localhost:5432/rsvp_db"
	conn, err := pgx.Connect(context.Background(), connURL)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close(context.Background())

	query := "SELECT name, phone, email, will_attend FROM rsvp"
	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		fmt.Println(err)
		return
	}
	for rows.Next() {
		rsvp := new(Rsvp)
		if err = rsvp.Scan(rows); err != nil {
			fmt.Println(err)
		}
		responses = append(responses, rsvp)
	}

	templates["list"].Execute(writer, responses)
}

func main() {
	loadTemplates()

	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/form", formHandler)
	http.HandleFunc("/list", listHandler)

	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		fmt.Println(err)
	}
}
