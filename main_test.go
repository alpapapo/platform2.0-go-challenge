package main

import (
	"database/sql"
	"fmt"
	"github.com/GlobalWebIndex/platform2.0-go-challenge/app"
	c "github.com/GlobalWebIndex/platform2.0-go-challenge/common"
	"github.com/GlobalWebIndex/platform2.0-go-challenge/models"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var a app.App

func TestMain(m *testing.M) {
	a = app.App{}
	a.Initialize("test")
	// Setup the cleaner
	cleaner := cleanupTables(a.DB)
	//_ = getTestidis()

	defer cleaner()

	m.Run()
}

func cleanupTables(db *gorm.DB) func() {
	type entity struct {
		table   string
		keyname string
		key     interface{}
	}
	var entries []entity
	hookName := "cleanupHook"
	//hookAssociations := "cleanupAsHook"

	// Setup the onCreate Hook
	db.Callback().Create().After("gorm:create").Register(hookName, func(scope *gorm.Scope) {
		fmt.Printf("Inserted entities of %s with %s=%v\n", scope.TableName(), scope.PrimaryKey(), scope.PrimaryKeyValue())
		entries = append(entries, entity{table: scope.TableName(), keyname: scope.PrimaryKey(), key: scope.PrimaryKeyValue()})
	})
	return func() {
		// Remove the hook once we're done
		defer db.Callback().Create().Remove(hookName)
		// Find out if the current db object is already a transaction
		_, inTransaction := db.CommonDB().(*sql.Tx)
		tx := db
		if !inTransaction {
			tx = db.Begin()
		}
		// Loop from the end. It is important that we delete the entries in the
		// reverse order of their insertion
		for i := len(entries) - 1; i >= 0; i-- {
			entry := entries[i]
			fmt.Printf("Deleting entities from '%s' table with key %v\n", entry.table, entry.key)
			tx.Table(entry.table).Where(entry.keyname+" = ?", entry.key).Delete("")
		}

		if !inTransaction {
			tx.Commit()
		}
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func getTestidis() (*models.User) {
	user := &models.User{}
	models.GetDB().Model(&models.User{}).Where(models.User{Email: "takis@testidis.com"}).Attrs(models.User{Password: c.SHA256OfString("secret")}).FirstOrCreate(&user)
	return user
}

func seedOneUser(id int) (models.User, error) {

	user := models.User{
		Email:    fmt.Sprintf("tester%d@gmail.com", id),
		Password: c.SHA256OfString(fmt.Sprintf("secret%d", id)),
	}

	err := models.GetDB().Model(&models.User{}).Create(&user).Error
	if err != nil {
		log.Fatalf("cannot seed users table: %v", err)
	}
	return user, nil
}



