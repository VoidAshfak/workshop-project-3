package main

import (
	"fmt"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	echo "github.com/labstack/echo/v4"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	// DB connection
	var err error
	dsn := "root:ashfak@tcp(localhost:3306)/todo-point?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect to the database!")
	}

	// Create a new instance of the Echo application
	e := echo.New()
	fmt.Println(db)

	// Define routes
	e.POST("/users", createUser)
	e.PATCH("/users/:id", updateUser)
	e.DELETE("/users/:id", deleteUser)

	// Start the server
	err = e.Start(":8080")
	if err != nil {
		panic(err)
	}
}

type User struct {
	ID             uint   `json:"id" gorm:"primaryKey"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Country        string `json:"country"`
	ProfilePicture string `json:"profile_picture"`
}

// Handler for creating a new user
func createUser(c echo.Context) error {
	reqBody := new(User)
	if err := c.Bind(reqBody); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request body")
	}

	// Create a new user object
	user := User{
		ID:             reqBody.ID,
		FirstName:      reqBody.FirstName,
		LastName:       reqBody.LastName,
		Country:        reqBody.Country,
		ProfilePicture: reqBody.ProfilePicture,
	}

	result := db.Create(&user)
	if result.Error != nil {
		// Return an error response if there's an issue with creating the user
		return c.JSON(http.StatusInternalServerError, "Failed to create a new user")
	}

	return c.JSON(http.StatusOK, user)
}

// Handler for updating an existing user
func updateUser(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid user ID")
	}

	// Fetch the existing user from the database
	existingUser := User{}
	result := db.First(&existingUser, id)
	if result.Error != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	reqBody := new(User)
	if err := c.Bind(reqBody); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request body")
	}

	// Update the user fields if provided in the request (conditional updating)
	updateFields := make(map[string]interface{})
	if reqBody.FirstName != "" {
		updateFields["first_name"] = reqBody.FirstName
	}
	if reqBody.LastName != "" {
		updateFields["last_name"] = reqBody.LastName
	}
	if reqBody.Country != "" {
		updateFields["country"] = reqBody.Country
	}
	if reqBody.ProfilePicture != "" {
		updateFields["profile_picture"] = reqBody.ProfilePicture
	}

	result = db.Model(&existingUser).Updates(updateFields)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to update the user")
	}

	return c.JSON(http.StatusOK, "User updated successfully")
}

// Handler for deleting a user
func deleteUser(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid user ID")
	}

	result := db.Delete(&User{}, id)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "User could not be deleted"})
	}

	// Check if the user was deleted successfully
	rowsAffected := result.RowsAffected
	if rowsAffected == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	return c.JSON(http.StatusOK, "User deleted successfully")
}
