package main

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type Message struct {
	User    string `json:"user"`
	Content string `json:"content"`
}

var (
	messages []Message  // Хранилище сообщений (в памяти)
	mu       sync.Mutex // Для потокобезопасной работы
)

func main() {
	r := gin.Default()

	// Маршрут для получения всех сообщений
	r.GET("/messages", getMessages)
	// Маршрут для отправки нового сообщения
	r.POST("/messages", postMessage)
	// Проверка состояния сервера
	r.GET("/status", checkStatus)
	r.Run(":8080")
}

func getMessages(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()

	c.JSON(http.StatusOK, gin.H{
		"messages": messages,
	})
}

func postMessage(c *gin.Context) {
	var newMessage Message

	// Считывание тела запроса
	if err := c.ShouldBindJSON(&newMessage); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Отсутствует JSON"})
		return
	}

	mu.Lock()
	messages = append(messages, newMessage)
	mu.Unlock()

	c.JSON(http.StatusOK, gin.H{
		"message": "Сообщение отправлено",
	})
}

func checkStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "Сервер запущен",
	})
}
