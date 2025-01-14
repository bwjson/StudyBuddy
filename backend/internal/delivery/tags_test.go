package delivery

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/bwjson/StudyBuddy/internal/dto"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	dsn := "host=localhost user=postgres password=5432 dbname=studybuddy_test port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to %v", err)
	}
	return db
}
func TestGetTagsByUser(t *testing.T) {

	db := setupTestDB(t)

	userID := 1
	tagID := 1

	db.Exec("INSERT INTO users (id, name, username, password_hash, email) VALUES (?, ?, ?, ?, ?)",
		userID, "testname", "test", "hashedpassword", "test@test.com")
	db.Exec("INSERT INTO tags (id, title, description) VALUES (?, ?, ?)",
		tagID, "Golang", "Tech-related tags")
	db.Exec("INSERT INTO user_tags (user_id, tag_id) VALUES (?, ?)", userID, tagID)

	defer func() {

		db.Exec("DELETE FROM user_tags WHERE user_id = ?", userID)
		db.Exec("DELETE FROM tags WHERE title = ?", "Golang")
		db.Exec("DELETE FROM users WHERE username = ?", "test")
	}()

	handler := &Handler{
		db:   db,
		log:  logrus.New(),
		smtp: nil,  
	}

	
	r := gin.Default()
	r.GET("/tags/usertags/:id", handler.getTagsByUser)

	req, _ := http.NewRequest(http.MethodGet, "/tags/usertags/1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	type successResponse struct {
		Status  string      `json:"status"`
		Message string      `json:"message"`
		Data    []map[string]interface{} `json:"data"`
	}
	var response successResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("could not parse data: %v", err)
	}

	var found bool
	for _, tag := range response.Data {
		if tag["title"] == "Golang" {
			found = true
			break
		}
	}

	assert.True(t, found, "Expected to find tag with title 'Golang'")
}

func TestGetUsersByTag(t *testing.T) {

	db := setupTestDB(t)

	tagID := 1
	userID := 1
	user := dto.User{
		ID:       uint(userID),
		Name:     "testname",
		Username: "test",
		Email:    "test@test.com",
	}
	tag := dto.Tag{
		ID:    uint(tagID),
		Title: "Golang",
	}

	db.Create(&user)
	db.Create(&tag)
	db.Exec("INSERT INTO user_tags (user_id, tag_id) VALUES (?, ?)", userID, tagID)

	defer func() {
		db.Exec("DELETE FROM user_tags WHERE tag_id = ?", tagID)
		db.Exec("DELETE FROM users WHERE id = ?", userID)
		db.Exec("DELETE FROM tags WHERE id = ?", tagID)
	}()

	handler := &Handler{
		db:   db,
		log:  logrus.New(),
		smtp: nil,
	}

	r := gin.Default()
	r.GET("/tags/users/:id", handler.getUsersByTag)

	req, _ := http.NewRequest(http.MethodGet, "/tags/users/1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Data dto.UsersWithPagination `json:"data"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("could not parse data: %v", err)
	}
	
	assert.Equal(t, 1, response.Data.TotalCount, "Expected the total number of users - 1")
	assert.Equal(t, "testname", response.Data.User[0].Name, "Expected to find user with name 'testname'")
}
