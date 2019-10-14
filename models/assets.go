package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"os"
	"time"
)

const dataFile = "sample/data.json"

type JSON map[string]interface{}


type Asset struct {
	gorm.Model
	Type  string `json:"type"`
	Description   string `json:"desc"`
	AssetData JSON `json:"data" sql:"type:json"`
	Users	[]*User	`json:"-" gorm:"many2many:user_asset;"`
}


func (j JSON) Value() (driver.Value, error) {
	valueString, err := json.Marshal(j)
	return string(valueString), err
}


func (j *JSON) Scan(value interface{}) error {
	if err := json.Unmarshal(value.([]byte), &j); err != nil {
		return err
	}
	return nil
}


func (asset *Asset) Validate() error {

	if asset.Type == "" {
		return errors.New("Required Type")
	}
	if asset.Description == "" {
		return errors.New("Required Description")
	}
	return nil

}

func (asset *Asset) SaveAsset() (*Asset, error) {
	var err error
	err = GetDB().Model(&Asset{}).Create(&asset).Error
	if err != nil {
		return &Asset{}, err
	}
	return asset, nil
}

func GetUserAssets(userId uint) (*[]Asset, error) {
	var err error
	user := GetUser(userId)
	assets := []Asset{}
	err = GetDB().Model(&user).Related(&assets, "Assets").Error
	if err != nil {
		fmt.Println(err)
		return &[]Asset{}, err
	}
	return &assets, nil
}

func (asset *Asset) GetAllAssets() (*[]Asset, error) {
	var err error
	assets := []Asset{}
	err = GetDB().Model(&Asset{}).Find(&assets).Error
	if err != nil {
		return &[]Asset{}, err
	}
	return &assets, nil
}

func (asset *Asset) GetAssetById(id uint64) (*Asset, error) {
	var err error
	err = GetDB().Model(&Asset{}).Where("id = ?", id).First(asset).Error
	if err != nil {
		return &Asset{}, err
	}
	return asset, nil
}

func (asset *Asset) UpdateAsset(id uint64) (*Asset, error) {
	var err error
	db = GetDB().Model(&Asset{}).Where("id = ?", id).First(&Asset{}).UpdateColumns(
		map[string]interface{}{
			"description": asset.Description,
			"updated_at": time.Now(),
		},
	)
	err = db.Model(&Asset{}).Where("id = ?", id).First(&asset).Error


	if err != nil {
		return &Asset{}, err
	}
	return asset, nil
}

func (asset *Asset) DeleteAsset(id uint64) (int, error) {
	var usersAffected int
	usersAffected = GetDB().Model(&asset).Association("Users").Count()
	GetDB().Model(&asset).Association("Users").Clear()
	db =  GetDB().Model(&Asset{}).Where("id = ?", id).Delete(&Asset{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Asset not found")
		}
		return 0, db.Error
	}
	return usersAffected, nil

}

func MarkAssetsAsFavoritesForUser(userID uint, assetIDs []int) (error) {
	user := GetUser(userID)
	assets := []Asset{}
	err := GetDB().Model(&Asset{}).Where("id IN (?)", assetIDs).Find(&assets).Error
	if err != nil {
		return err
	}
	for _, asset := range assets {
		err := GetDB().Model(&user).Association("Assets").Append(asset).Error
		if err != nil {
			return err
		}
	}
	return err
}

func UnMarkAssetsAsFavoritesForUser(userID uint, assetIDs []int) (error) {
	user := GetUser(userID)
	assets := []Asset{}
	err := GetDB().Model(&Asset{}).Where("id IN (?)", assetIDs).Find(&assets).Error
	if err != nil {
		return err
	}
	for _, asset := range assets {
		err := GetDB().Model(&user).Association("Assets").Delete(asset).Error
		if err != nil {
			return err
		}
	}
	return err
}

func SeedAssets() ([]*Asset, error) {

	// Open the file
	file, err := os.Open(dataFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var assets []*Asset
	err = json.NewDecoder(file).Decode(&assets)

	for _, asset := range assets {
		GetDB().Create(&asset)
	}

	return assets, err
}
