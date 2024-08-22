package main

import (
	"bufio"
	"keepdata/internal/logger"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true }, // Важно для тестов на localhost, чтобы разрешить любые источники
}

func homeHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", "")
}

func aboutHandler(c *gin.Context) {
	c.Redirect(http.StatusFound, "/")
}

func streamHandler(c *gin.Context) {

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Log.Error("ошибка при запуске websocket", zap.Error(err))
		return
	}
	defer ws.Close()

	file, err := os.OpenFile("app.log", os.O_RDONLY, 0644)
	if err != nil {
		logger.Log.Error("ошибка при открытии файла", zap.Error(err))
	}
	defer file.Close()

	for { // стрим читаемых из файла логов по веб сокету
		sc := bufio.NewScanner(file)
		for sc.Scan() {
			if err := ws.WriteMessage(websocket.TextMessage, []byte(sc.Text())); err != nil {
				logger.Log.Error("ошибка при отправке сообщений по websocket", zap.Error(err))
				return
			}
		}
	}
}

func filterHandler(c *gin.Context) {
	user := c.Query("user")
	action := c.Query("action")
	resource := c.Query("resource")

	var err error

	logs, err := getFilteredLogs(user, action, resource)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "filter.tmpl", logs)
}
