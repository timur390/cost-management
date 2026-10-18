package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	_ "modernc.org/sqlite" // без этой штуки бд не запустится вообще 
)

// простая структура для одной строчки
type Cost struct {
	ID    int
	Sum   float64
	Title string
	Date  string
}

var db *sql.DB
var tpl *template.Template


func home(w http.ResponseWriter, r *http.Request) {
	// если просто открыли сайт (GET), то грузим данные из базы в таблицу
	if r.Method == "GET" {
		rows, _ := db.Query("SELECT id, sum, title, date FROM costs ORDER BY id DESC")
		defer rows.Close()

		var list []Cost
		for rows.Next() {
			var c Cost
			_ = rows.Scan(&c.ID, &c.Sum, &c.Title, &c.Date)
			list = append(list, c)
		}

		
		_ = tpl.ExecuteTemplate(w, "base.html", list)
		return
	}

	// если отправили форму (POST), то тупо пихаем новую запись в бд
	if r.Method == "POST" {
		sumStr := r.FormValue("sum")
		title := r.FormValue("title")
		date := r.FormValue("date")

		sum, _ := strconv.ParseFloat(sumStr, 64)

		
		_, _ = db.Exec("INSERT INTO costs (sum, title, date) VALUES (?, ?, ?)", sum, title, date)

		// страница обновилась
		http.Redirect(w, r, "/", 303)
		return
	}
}

// хендлер для кнопки удалить
func cancel(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		idStr := r.FormValue("id")
		id, _ := strconv.Atoi(idStr)

		// удаляем из таблицы по id
		_, _ = db.Exec("DELETE FROM costs WHERE id = ?", id)
	}
	http.Redirect(w, r, "/", 303)
}

func main() {
	// открываем файл базы данных, он сам создается в папке
	db, _ = sql.Open("sqlite", "finance.db")
	defer db.Close()

	// создаем самую обычную таблицу если её еще нет
	_, _ = db.Exec(`CREATE TABLE IF NOT EXISTS costs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		sum REAL,
		title TEXT,
		date TEXT
	);`)

	
	tpl = template.Must(template.ParseGlob("web/templates/*.html"))

	// маршруты для страниц
	http.HandleFunc("/", home)
	http.HandleFunc("/cancel", cancel)

	fmt.Println("Сервер третьего этапа пашет на http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
