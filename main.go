package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)


type Cost struct {
	Sum       float64
	Title     string
	CreatedAt string
}

// Временное хранилище прямо в озу
var memoryStorage []Cost
var tpl *template.Template

// Главный хендлер для сайта
func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Если юзер просто открыл страницу  — показываем таблицу
	if r.Method == "GET" {
		// Рендерим каркас base.html 
		err := tpl.ExecuteTemplate(w, "base.html", memoryStorage)
		if err != nil {
			http.Error(w, "Ошибка шаблона", http.StatusInternalServerError)
		}
		return
	}

	// если юзер нажал на кнопку и отправил форму 
	if r.Method == "POST" {
		sumStr := r.FormValue("sum")
		title := r.FormValue("title")
		date := r.FormValue("date")

		// Конвертируем строку с ценой в нормальное число float64
		sum, _ := strconv.ParseFloat(sumStr, 64)

		// Создаем новую запись расходов 
		newRecord := Cost{
			Sum:       sum,
			Title:     title,
			CreatedAt: date,
		}

		// Добавляем запись в наш глобальный массив в памяти
		memoryStorage = append(memoryStorage, newRecord)

		// Возвращаем юзера на главную страницу
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
}

// Служебные страницы из 1 этапа оставляем
func aboutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Прога для ленивых: записываем, сколько бабок улетело в никуда.")
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(405)
		fmt.Fprintf(w, "Сюда можно только через GET заходить, сорян")
		return
	}
	w.WriteHeader(200)
	fmt.Fprintf(w, "pong")
}

func main() {
	// Компилируем наши новые HTML шаблоны из папки web/templates
	tpl = template.Must(template.ParseGlob("web/templates/*.html"))

	// Настраиваем маршруты
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/about", aboutHandler)
	http.HandleFunc("/ping", pingHandler)

	fmt.Println("Сервер Второго этапа (хранение в памяти) запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
