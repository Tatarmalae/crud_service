package main

import (
	"encoding/json"
	"fmt"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"os"
	"strconv"
)

type User struct {
	UserID int64  `json:"user_id"`
	Name   string `json:"name"`
	Phone  int64  `json:"phone"`
	Email  string `json:"email"`
}

var users []User
var dataFile = "users.json"

func main() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Gophers!")
	})
	e.POST("/user", saveUser)
	e.GET("/user/:id", getUser)
	e.PUT("/user/:id", updateUser)
	e.DELETE("/user/:id", deleteUser)
	e.GET("/users", getUsers)

	loadUsersFromFile()

	e.Logger.Fatal(e.Start(":1323"))
}

// Создание нового пользователя
func saveUser(c echo.Context) error {
	user := new(User)
	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	// Генерация фейковых данных
	user.UserID = int64(len(users) + 1)
	user.Name = gofakeit.Name()
	user.Phone = gofakeit.Int64()
	user.Email = gofakeit.Email()

	users = append(users, *user)
	saveUsersToFile()
	return c.JSON(http.StatusCreated, user)
}

// Получение пользователя по ID
func getUser(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	for _, user := range users {
		if user.UserID == id {
			return c.JSON(http.StatusOK, user)
		}
	}
	return c.JSON(http.StatusNotFound, "Пользователь не найден")
}

// Обновление пользователя по ID
func updateUser(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	userToUpdate := new(User)
	if err := c.Bind(userToUpdate); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	for i, user := range users {
		if user.UserID == id {
			users[i] = *userToUpdate
			saveUsersToFile()
			return c.JSON(http.StatusOK, userToUpdate)
		}
	}
	return c.JSON(http.StatusNotFound, "Пользователь не найден")
}

// Удаление пользователя по ID
func deleteUser(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	for i, user := range users {
		if user.UserID == id {
			users = append(users[:i], users[i+1:]...)
			saveUsersToFile()
			return c.JSON(http.StatusOK, "Пользователь успешно удален")
		}
	}
	return c.JSON(http.StatusNotFound, "Пользователь не найден")
}

// Получение пользователей по массиву ID
func getUsers(c echo.Context) error {
	var userIDs []int64
	if err := c.Bind(&userIDs); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	var foundUsers []User
	for _, id := range userIDs {
		for _, user := range users {
			if user.UserID == id {
				foundUsers = append(foundUsers, user)
				break
			}
		}
	}
	return c.JSON(http.StatusOK, foundUsers)
}

// Загрузка данных пользователей из файла
func loadUsersFromFile() {
	file, err := os.ReadFile(dataFile)
	if err != nil {
		log.Println("Ошибка чтения данных из файла:", err)
		return
	}
	err = json.Unmarshal(file, &users)
	if err != nil {
		log.Println("Ошибка декодирования данны:", err)
		return
	}
	fmt.Println("Пользователи загружены из файла")
}

// Сохранение пользователя в файл
func saveUsersToFile() {
	fileData, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		log.Println("Ошибка кодирования данных:", err)
		return
	}
	err = os.WriteFile(dataFile, fileData, 0644)
	if err != nil {
		log.Println("Ошибка записи данных в файл:", err)
		return
	}
	fmt.Println("Пользователь сохранен в файл")
}
