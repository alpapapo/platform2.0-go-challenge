package controllers

import (
	"encoding/json"
	c "github.com/GlobalWebIndex/platform2.0-go-challenge/common"
	"github.com/GlobalWebIndex/platform2.0-go-challenge/models"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
)


var GetAssets = func(w http.ResponseWriter, r *http.Request) {
	asset := models.Asset{}

	data, err := asset.GetAllAssets()
	if err != nil {
		resp := c.Message(false, "fail to retrieve assets")
		c.Respond(w, http.StatusNotFound, resp)
	} else {
		resp := c.Message(true, "success")
		resp["data"] = data
		c.Respond(w, http.StatusOK, resp)
	}
}


var GetAssetsFavor = func(w http.ResponseWriter, r *http.Request) {

	id := r.Context().Value("user").(uint)
	data, err := models.GetUserAssets(id)
	if err != nil {
		resp := c.Message(false, "fail to retrieve assets")
		c.Respond(w, http.StatusNotFound, resp)
	} else {
		resp := c.Message(true, "success")
		resp["data"] = data
		c.Respond(w, http.StatusOK, resp)
	}
}


var GetAsset = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	assetID, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		resp := c.Message(false, "assetID is not recognisable")
		c.Respond(w, http.StatusBadRequest, resp)
		return
	}
	asset := models.Asset{}

	assetRequested, err := asset.GetAssetById(assetID)
	if err != nil {
		resp := c.Message(false, "Asset Not Found")
		c.Respond(w, http.StatusNotFound, resp)
		return
	}
	resp := c.Message(true, "success")
	resp["data"] = assetRequested
	c.Respond(w, http.StatusOK, resp)
}


var UpdateAsset = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	assetID, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		resp := c.Message(false, "assetID is not recognisable")
		c.Respond(w, http.StatusBadRequest, resp)
		return
	}
	asset := models.Asset{}

	// Check if Asset exists
	err = models.GetDB().Model(models.Asset{}).Where("id = ?", assetID).Take(&asset).Error
	if err != nil {
		resp := c.Message(false, "Asset Not Found")
		c.Respond(w, http.StatusNotFound, resp)
		return
	}

	// Read the data posted
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		resp := c.Message(false, "Check Data to be Updated")
		c.Respond(w, http.StatusUnprocessableEntity, resp)
		return
	}
	// Start processing the request data
	assetUpdate := models.Asset{}
	err = json.Unmarshal(body, &assetUpdate)
	if err != nil {
		resp := c.Message(false, "Check Data to be Updated")
		c.Respond(w, http.StatusUnprocessableEntity, resp)
		return
	}

	err = assetUpdate.Validate()

	assetUpdated, err := assetUpdate.UpdateAsset(assetID)

	if err != nil {
		resp := c.Message(false, "Internal Server Error")
		c.Respond(w, http.StatusInternalServerError, resp)
		return
	}
	resp := c.Message(true, "Asset Updated!")
	resp["data"] = assetUpdated
	c.Respond(w, http.StatusOK, resp)

}


var DeleteAsset = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	assetID, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		resp := c.Message(false, "assetID is not recognisable")
		c.Respond(w, http.StatusBadRequest, resp)
		return
	}
	asset := models.Asset{}

	// Check if Asset exists
	err = models.GetDB().Model(models.Asset{}).Where("id = ?", assetID).Take(&asset).Error
	if err != nil {
		resp := c.Message(false, "Asset Not Found")
		c.Respond(w, http.StatusNotFound, resp)
		return
	}

	usersAffected, err := asset.DeleteAsset(assetID)
	if err != nil {
		resp := c.Message(false, "Fail to delete Asset")
		c.Respond(w, http.StatusInternalServerError, resp)
		return
	}

	resp := c.Message(true, "Asset Deleted!")
	resp["usersAffected"] = usersAffected
	c.Respond(w, http.StatusOK, resp)
}


var SetAssetsFavor = func(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user").(uint)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		resp := c.Message(false, "Check Data to be Updated")
		c.Respond(w, http.StatusUnprocessableEntity, resp)
		return
	}
	var ids []int
	err = json.Unmarshal(body, &ids)
	if err != nil {
		resp := c.Message(false, "Check Data to be Updated")
		c.Respond(w, http.StatusUnprocessableEntity, resp)
		return
	}
	err = models.MarkAssetsAsFavoritesForUser(userID, ids)
	if err != nil {
		resp := c.Message(false, "Fail to Mark Assets as Favorites")
		c.Respond(w, http.StatusInternalServerError, resp)
		return
	}

	resp := c.Message(true, "Assets Marked As Favorites!")
	c.Respond(w, http.StatusOK, resp)
}


var UnsetAssetsFavor = func(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user").(uint)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		resp := c.Message(false, "Check Data to be Updated")
		c.Respond(w, http.StatusUnprocessableEntity, resp)
		return
	}
	var ids []int
	err = json.Unmarshal(body, &ids)
	if err != nil {
		resp := c.Message(false, "Check Data to be Updated")
		c.Respond(w, http.StatusUnprocessableEntity, resp)
		return
	}
	err = models.UnMarkAssetsAsFavoritesForUser(userID, ids)

	if err != nil {
		resp := c.Message(false, "Fail to Remove Assets from Favorites")
		c.Respond(w, http.StatusInternalServerError, resp)
		return
	}

	resp := c.Message(true, "Assets Removed From Favorites!")
	c.Respond(w, http.StatusOK, resp)

}


var PopulateAssets = func(w http.ResponseWriter, r *http.Request) {
	_,err := models.SeedAssets()
	if err != nil {
		resp := c.Message(false, "Failed to Seed Assets")
		c.Respond(w, http.StatusInternalServerError, resp)
		return
	}
	resp := c.Message(true, "Assets Added!")
	c.Respond(w, http.StatusOK, resp)

}