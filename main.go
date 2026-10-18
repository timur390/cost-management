package main

import (
	"fmt"
	"net/http"
)

// хендлер для главной, 
func home(w http.ResponseWriter, r *http.Request) {
	// если оставить пустой корень, сюда будут лететь все левые запросы, ну и ладно пока
	fmt.Fprintf(w, "Здарова! Это стартовая страница журнала расходов.")
}

// страница про сам проект
func about(w http.ResponseWriter, r *http.Request) {
	// написал чисто отсебятину, чтобы была хоть какая-то инфа на вкладке
	fmt.Fprintf(w, "Прога для ленивых: записываем, сколько бабок улетело в никуда.")
}

// тут проверяем статус сервера (пинг-понг)
func ping(w http.ResponseWriter, r *http.Request) {
	// костыль для проверки метода
	if r.Method != "GET" {
		w.WriteHeader(405) // кидаем 405 код, если 
		fmt.Fprintf(w, "Сюда можно только через GET заходить, сорян")
		return
	}

	// если дошли сюда, пишем ответ
	w.WriteHeader(200)
	fmt.Fprintf(w, "pong")
}

func main() {
	// распределяем функции по адресам, стандартный роутер net/http
	http.HandleFunc("/", home)
	http.HandleFunc("/about", about)
	http.HandleFunc("/ping", ping)

	// маяк в консоль, чтобы видеть, что процесс не завис при старте
	fmt.Println("Сервер вроде поднялся, чекай http://localhost:8080")
	
	// врубаем бесконечный цикл сервера на 8080 порту
	http.ListenAndServe(":8080", nil)
}
