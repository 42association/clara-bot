package main

import (
	"database/sql"
//	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	_ "github.com/mattn/go-sqlite3"
)

const (
	dbPath = "../db/clara.sqlite3"
)

type WaterServer struct {
	BottleID		int `db:"bottle_id"`
	BottleStatus    string `db:"bottle_status"`
	UserName		string `db:"user_name"`
	ExchageStatus	string `db:"exchange_status"`
}

type BottleResponse struct {
	BottleID     int    `json:"BottleID"`
	BottleStatus string `json:"BottleStatus"`
}

type Response struct {
	Message string `json:"message"`
}

func openDB() *sql.DB {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	return db
}

func root(c echo.Context) error {
	res := Response{Message: "Hello, world!"}
	return c.JSON(http.StatusOK, res)
}

func addUser(c echo.Context) error {
    userName := c.FormValue("user_name")

    db := openDB()
    defer db.Close()

    // 最大のbottle_idを見つける
    var maxID int
    err := db.QueryRow("SELECT MAX(bottle_id) FROM clara").Scan(&maxID)
    if err != nil {
        log.Errorf("Error finding max bottle_id: %v", err)
        return c.JSON(http.StatusInternalServerError, Response{Message: "Failed to find max bottle_id"})
    }

    // 最新のbottle_idのレコードのuser_nameを更新し、exchange_statusをtrueにする
    _, err = db.Exec("UPDATE clara SET user_name = ?, exchange_status = 'true' WHERE bottle_id = ?", userName, maxID)
    if err != nil {
        log.Errorf("Error updating user_name and exchange_status: %v", err)
        return c.JSON(http.StatusInternalServerError, Response{Message: "Failed to update user_name and exchange_status"})
    }

    // 新しいbottle_idで新しい行を追加する（bottle_idは自動インクリメントなので指定しない）
    // 初期値: bottle_status = 'true', user_name = '', exchange_status = 'false'
    _, err = db.Exec("INSERT INTO clara (bottle_status, user_name, exchange_status) VALUES ('true', '', 'false')")
    if err != nil {
        log.Errorf("Error inserting new bottle with initial values: %v", err)
        return c.JSON(http.StatusInternalServerError, Response{Message: "Failed to insert new bottle with initial values"})
    }

    return c.JSON(http.StatusOK, Response{Message: "User name updated successfully and new bottle added"})
}

func getBottleStatus(c echo.Context) error {
	db := openDB()
	defer db.Close()

	var maxID int
	err1 := db.QueryRow("SELECT MAX(bottle_id) FROM clara").Scan(&maxID)
	if err1 != nil {
		log.Errorf("Error finding max bottle_id: %v", err1)
		return c.JSON(http.StatusInternalServerError, Response{Message: "Failed to find max bottle_id"})
	}
	var br BottleResponse
	err := db.QueryRow("SELECT bottle_id, bottle_status FROM clara WHERE bottle_id = ?", maxID).Scan(&br.BottleID, &br.BottleStatus)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, Response{Message: "Bottle not found"})
		}
		return c.JSON(http.StatusInternalServerError, Response{Message: "Failed to query database"})
	}

	return c.JSON(http.StatusOK, br)
}

func emptyBottle(c echo.Context) error {
    db := openDB()
    defer db.Close()

    // 更新対象のbottle_idを決定するために最大のbottle_idを見つける
    var maxID int
    err := db.QueryRow("SELECT MAX(bottle_id) FROM clara").Scan(&maxID)
    if err != nil {
        log.Errorf("Error finding max bottle_id: %v", err)
        return c.JSON(http.StatusInternalServerError, Response{Message: "Failed to find max bottle_id"})
    }

    // 最新のbottle_idに紐づくbottle_statusをfalseに更新する
    _, err = db.Exec("UPDATE clara SET bottle_status = 'false' WHERE bottle_id = ?", maxID)
    if err != nil {
        log.Errorf("Error updating bottle_status: %v", err)
        return c.JSON(http.StatusInternalServerError, Response{Message: "Failed to update bottle_status"})
    }

    return c.JSON(http.StatusOK, Response{Message: "Bottle status updated to empty successfully"})
}


func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Logger.SetLevel(log.DEBUG)
	frontURL := os.Getenv("FRONT_URL")
	if frontURL == "" {
		frontURL = "http://localhost:3000"
	}
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{frontURL},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	// Routes
	e.GET("/", root)
	e.POST("/adduser", addUser)
	e.GET("/bottlestatus", getBottleStatus)
	e.POST("/empty", emptyBottle)

	// Start server
	e.Logger.Fatal(e.Start(":9000"))
}
