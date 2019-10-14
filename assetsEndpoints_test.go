package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/GlobalWebIndex/platform2.0-go-challenge/models"
	"net/http"
	"testing"
)

func TestGetAssets(t *testing.T) {

	_, err := models.SeedAssets()
	user := getTestidis()
	user.SetToken()

	jwt := user.Token
	request, err := http.NewRequest("GET", "/api/assets", nil)
	if err != nil {
		t.Fatal(err)
	}
	request.Header = map[string][]string{
		"Authorization": {fmt.Sprintf("Bearer %s", jwt)},
	}

	response := executeRequest(request)
	var response_data map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &response_data)
	assetsCounter := len(response_data["data"].([]interface{}))
	if assetsCounter != 5 {
		t.Errorf("Expected 5 assets instead of %d", assetsCounter)
	}
	checkResponseCode(t, http.StatusOK, response.Code)
}


func TestGetAssetsFavor(t *testing.T) {
	assets, err := models.SeedAssets()
	user := getTestidis()
	user.SetToken()

	for _, asset := range assets[:2] {
		models.GetDB().Model(&user).Association("Assets").Append(asset)
	}

	jwt := user.Token
	request, err := http.NewRequest("GET", "/api/assets/favorites", nil)
	if err != nil {
		t.Fatal(err)
	}
	request.Header = map[string][]string{
		"Authorization": {fmt.Sprintf("Bearer %s", jwt)},
	}

	response := executeRequest(request)
	var response_data map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &response_data)
	assetsCounter := len(response_data["data"].([]interface{}))
	if assetsCounter != 2 {
		t.Errorf("Expected 2 assets instead of %d", assetsCounter)
	}
	checkResponseCode(t, http.StatusOK, response.Code)

	models.GetDB().Model(&user).Association("Assets").Clear()
}

func TestGetAsset(t *testing.T) {
	_, err := models.SeedAssets()
	user := getTestidis()
	user.SetToken()
	asset := models.Asset{}
	models.GetDB().Model(&models.Asset{}).Last(&asset)

	jwt := user.Token
	requestUrl := fmt.Sprintf("/api/assets/%d", asset.ID)
	request, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		t.Fatal(err)
	}
	request.Header = map[string][]string{
		"Authorization": {fmt.Sprintf("Bearer %s", jwt)},
	}
	//requestAssetID := fmt.Sprint(asset.ID)
	//request = mux.SetURLVars(request, map[string]string{"id": requestAssetID})

	response := executeRequest(request)
	var response_data map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &response_data)
	if response_data["message"] != "success" {
		t.Errorf("Expected succesful message")
	}
	checkResponseCode(t, http.StatusOK, response.Code)
}


func TestUpdateAsset(t *testing.T) {
	_, err := models.SeedAssets()
	user := getTestidis()
	user.SetToken()
	asset := models.Asset{}
	models.GetDB().Model(&models.Asset{}).Last(&asset)

	jwt := user.Token
	requestUrl := fmt.Sprintf("/api/assets/%d", asset.ID)
	updateJSON := `{"desc": "other_description"}`

	request, err := http.NewRequest("PUT", requestUrl, bytes.NewBufferString(updateJSON))
	if err != nil {
		t.Fatal(err)
	}
	request.Header = map[string][]string{
		"Authorization": {fmt.Sprintf("Bearer %s", jwt)},
	}

	response := executeRequest(request)
	var response_data map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &response_data)
	if response_data["message"] != "Asset Updated!" {
		t.Errorf("Expected succesful message")
	}
	checkResponseCode(t, http.StatusOK, response.Code)
}


func TestDeleteAsset(t *testing.T) {
	_, err := models.SeedAssets()
	user := getTestidis()
	user.SetToken()
	asset := models.Asset{}
	models.GetDB().Model(&models.Asset{}).Last(&asset)

	// make it favorite
	models.GetDB().Model(&user).Association("Assets").Append(asset)

	jwt := user.Token
	requestUrl := fmt.Sprintf("/api/assets/%d", asset.ID)

	request, err := http.NewRequest("DELETE", requestUrl, nil)
	if err != nil {
		t.Fatal(err)
	}
	request.Header = map[string][]string{
		"Authorization": {fmt.Sprintf("Bearer %s", jwt)},
	}

	response := executeRequest(request)
	var response_data map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &response_data)
	if response_data["message"] != "Asset Deleted!" {
		t.Errorf("Expected succesful message")
	}
	checkResponseCode(t, http.StatusOK, response.Code)
	models.GetDB().Model(&models.Asset{}).Association("Assets").Clear()

}

func TestSetAssetsFavor(t *testing.T) {
	assets, err := models.SeedAssets()
	user := getTestidis()
	user.SetToken()

	jwt := user.Token

	var assetIds []int
	for _, asset := range assets [:3] {
		assetIds = append(assetIds, int(asset.ID))
	}
	assetSliceJSON, _ := json.Marshal(assetIds)

	request, err := http.NewRequest("POST", "/api/assets/favorites", bytes.NewBufferString(string(assetSliceJSON)))
	if err != nil {
		t.Fatal(err)
	}
	request.Header = map[string][]string{
		"Authorization": {fmt.Sprintf("Bearer %s", jwt)},
	}

	response := executeRequest(request)
	checkResponseCode(t, http.StatusOK, response.Code)

	models.GetDB().Model(&user).Association("Assets").Clear()
}


func TestUnsetAssetsFavor(t *testing.T) {
	assets, err := models.SeedAssets()
	user := getTestidis()
	user.SetToken()
	for _, asset := range assets[:3] {
		models.GetDB().Model(&user).Association("Assets").Append(asset)
	}


	jwt := user.Token

	var assetIds []int
	for _, asset := range assets [:1] {
		assetIds = append(assetIds, int(asset.ID))
	}
	assetSliceJSON, _ := json.Marshal(assetIds)

	request, err := http.NewRequest("DELETE", "/api/assets/favorites", bytes.NewBufferString(string(assetSliceJSON)))
	if err != nil {
		t.Fatal(err)
	}
	request.Header = map[string][]string{
		"Authorization": {fmt.Sprintf("Bearer %s", jwt)},
	}

	response := executeRequest(request)
	checkResponseCode(t, http.StatusOK, response.Code)

	models.GetDB().Model(&user).Association("Assets").Clear()
}
