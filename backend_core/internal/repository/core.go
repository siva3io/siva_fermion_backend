package repository

import (
	"fermion/backend_core/db"
	model_core "fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/pagination"
	"fermion/backend_core/pkg/util/helpers"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/*
Copyright (C) 2022 Eunimart Omnichannel Pvt Ltd. (www.eunimart.com)
All rights reserved.
This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.
This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.
You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/
type Core interface {
	GetLookupTypes(interface{}, *pagination.Paginatevalue) ([]model_core.Lookuptype, error)
	GetLookupCodes(interface{}, *pagination.Paginatevalue) ([]model_core.Lookupcode, error)
	GetCountries(interface{}, *pagination.Paginatevalue) ([]model_core.Country, error)
	GetStates(interface{}, *pagination.Paginatevalue) ([]model_core.State, error)
	GetCurrencies(interface{}, *pagination.Paginatevalue) ([]model_core.Currency, error)
	SearchLookupCodes(query string) ([]model_core.Lookupcode, error)
	CreateTypes(data *model_core.Lookuptype) error
	CreateCodes(data *[]model_core.Lookupcode) error
	UpdateType(data *model_core.Lookuptype, query interface{}) error
	UpdateCode(data *model_core.Lookupcode, query interface{}) error
	DeleteType(id uint) error
	DeleteCode(query interface{}) error
	FindAllLookupTypes(query interface{}, p *pagination.Paginatevalue) ([]model_core.Lookuptype, error)
	FindAllLookupCodes(query interface{}, p *pagination.Paginatevalue) ([]model_core.Lookupcode, error)

	//apps

	//-----------------------------Store Apps-------------------------------------------------------
	CreateAppStore(data *model_core.AppStore) error
	UpdateAppStore(query map[string]interface{}, data *model_core.AppStore) error
	ListStoreApps(query interface{}, page *pagination.Paginatevalue) ([]model_core.AppStore, error)
	ListStoreAppsInstalled(page *pagination.Paginatevalue) ([]model_core.AppStore, error)
	ListInstalledApps(page *pagination.Paginatevalue) ([]model_core.InstalledApps, error)
	GetApp(query map[string]interface{}) (model_core.AppStore, error)
	SearchApps(query map[string]interface{}) (model_core.AppStore, error)
	CheckState(code string) bool

	//---------------install app-------------------------------------------------------------------
	InstallApp(data *model_core.InstalledApps) error
	UninstallApp(query map[string]interface{}) error

	//-----------------------Channel LookupCodes--------------------------------------------------------
	GetChannelLookupCodes(interface{}, *pagination.Paginatevalue) ([]model_core.ChannelLookupCodes, error)

	//------------------Meta data------------------------------------------------------------------------
	ListMetaData() ([]model_core.IRModel, error)
	ViewMetaData(model_name string) ([]model_core.IRModelFields, error)

	//----------------Company------------------------------------------------------------------------------
	CreateCompany(data *model_core.Company) error
	UpdateCompany(query map[string]interface{}, data *model_core.Company) error
	FindCompany(query map[string]interface{}) (model_core.Company, error)
}

type cores struct {
	db *gorm.DB
}

func NewCore() *cores {
	db := db.DbManager()
	return &cores{db}
}

func (r *cores) GetLookupTypes(query interface{}, p *pagination.Paginatevalue) ([]model_core.Lookuptype, error) {
	var data []model_core.Lookuptype
	res := r.db.Model(&model_core.Lookuptype{}).Scopes(helpers.Paginate(&model_core.Lookuptype{}, p, r.db)).Where(query).Preload("Lookupcodes").Find(&data)
	return data, res.Error
}

func (r *cores) FindAllLookupTypes(query interface{}, p *pagination.Paginatevalue) ([]model_core.Lookuptype, error) {
	var data []model_core.Lookuptype
	err := r.db.Preload(clause.Associations).Model(&model_core.Lookuptype{}).Scopes(helpers.Paginate(&model_core.Lookuptype{}, p, r.db)).Where(query).Find(&data).Error
	if err != nil {
		return data, err
	}
	return data, nil
}

func (r *cores) GetLookupCodes(query interface{}, p *pagination.Paginatevalue) ([]model_core.Lookupcode, error) {
	var data []model_core.Lookupcode
	res := r.db.Model(&model_core.Lookupcode{}).Scopes(helpers.Paginate(&model_core.Lookupcode{}, p, r.db)).Where(query).Find(&data)
	return data, res.Error
}

func (r *cores) FindAllLookupCodes(query interface{}, p *pagination.Paginatevalue) ([]model_core.Lookupcode, error) {
	var data []model_core.Lookupcode
	err := r.db.Preload(clause.Associations).Model(&model_core.Lookupcode{}).Scopes(helpers.Paginate(&model_core.Lookupcode{}, p, r.db)).Where(query).Find(&data).Error
	if err != nil {
		return data, err
	}
	return data, nil
}

func (r *cores) GetCountries(query interface{}, p *pagination.Paginatevalue) ([]model_core.Country, error) {
	var data []model_core.Country
	res := r.db.Model(&model_core.Country{}).Scopes(helpers.Paginate(&model_core.Country{}, p, r.db)).Where(query).Find(&data)
	return data, res.Error
}
func (r *cores) GetStates(query interface{}, p *pagination.Paginatevalue) ([]model_core.State, error) {
	var data []model_core.State
	res := r.db.Model(&model_core.State{}).Scopes(helpers.Paginate(&model_core.State{}, p, r.db)).Where(query).Find(&data)
	return data, res.Error
}
func (r *cores) GetCurrencies(query interface{}, p *pagination.Paginatevalue) ([]model_core.Currency, error) {
	var data []model_core.Currency
	res := r.db.Model(&model_core.Currency{}).Scopes(helpers.Paginate(&model_core.Currency{}, p, r.db)).Where(query).Find(&data)
	return data, res.Error
}
func (r *cores) SearchLookupCodes(query string) ([]model_core.Lookupcode, error) {
	var data []model_core.Lookupcode
	err := r.db.Find(&data, "display_name ILIKE ?", "%"+query+"%").Error
	if err != nil {
		return data, err
	}
	return data, nil
}
func (r *cores) CreateTypes(data *model_core.Lookuptype) error {
	res := r.db.Model(&model_core.Lookuptype{}).Create(&data)
	return res.Error
}
func (r *cores) CreateCodes(data *[]model_core.Lookupcode) error {
	res := r.db.Model(&model_core.Lookupcode{}).Create(&data)
	return res.Error
}
func (r *cores) UpdateType(data *model_core.Lookuptype, query interface{}) error {
	res := r.db.Model(&model_core.Lookuptype{}).Where(query).Updates(&data)
	return res.Error
}
func (r *cores) UpdateCode(data *model_core.Lookupcode, query interface{}) error {
	res := r.db.Model(&model_core.Lookupcode{}).Where(query).Updates(&data)
	return res.Error
}
func (r *cores) DeleteType(id uint) error {
	var data model_core.Lookuptype
	res := r.db.Model(&model_core.Lookuptype{}).Where("id=?", id).Delete(&data)
	return res.Error
}
func (r *cores) DeleteCode(query interface{}) error {
	var data model_core.Lookupcode
	res := r.db.Model(&model_core.Lookupcode{}).Where(query).Delete(&data)
	return res.Error
}

// -------------------------------------------Store apps----------------------------------------------------------
func (r *cores) CreateAppStore(data *model_core.AppStore) error {
	err := r.db.Create(data).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *cores) UpdateAppStore(query map[string]interface{}, data *model_core.AppStore) error {
	err := r.db.Model(&model_core.AppStore{}).Where(query).Updates(data).Error
	return err
}
func (r *cores) ListStoreApps(query interface{}, page *pagination.Paginatevalue) ([]model_core.AppStore, error) {
	var data []model_core.AppStore
	var queryFilterIds []int
	if query != nil && query.(string) == "recommended_apps" {
		queryFilterIds = []int{24, 27, 29, 35, 22, 23}
	}
	if query != nil && query.(string) == "sponsored_apps" {
		queryFilterIds = []int{49, 66, 62, 63}
	}
	if query != nil && query.(string) == "trending_apps" {
		queryFilterIds = []int{55, 63, 49, 18, 35, 24}
	}
	if query != nil && query.(string) == "new_apps" {
		queryFilterIds = []int{31, 22, 37, 46, 39, 65}
	}
	res := r.db.Model(&model_core.AppStore{}).Scopes(helpers.Paginate(&model_core.AppStore{}, page, r.db)).Find(&data, queryFilterIds)
	return data, res.Error
}
func (r *cores) ListStoreAppsInstalled(page *pagination.Paginatevalue) ([]model_core.AppStore, error) {
	var data []model_core.AppStore
	res := r.db.Scopes(helpers.Paginate(&model_core.AppStore{}, page, r.db)).Joins("JOIN installed_apps on installed_apps.code = app_stores.code").Find(&data)
	return data, res.Error
}
func (r *cores) ListInstalledApps(page *pagination.Paginatevalue) ([]model_core.InstalledApps, error) {
	var data []model_core.InstalledApps
	res := r.db.Scopes(helpers.Paginate(&model_core.InstalledApps{}, page, r.db)).Find(&data)
	return data, res.Error
}

func (r *cores) GetApp(query map[string]interface{}) (model_core.AppStore, error) {
	var data model_core.AppStore
	res := r.db.Model(&model_core.AppStore{}).Where(query).Preload("Currency").First(&data)
	return data, res.Error
}
func (r *cores) SearchApps(query map[string]interface{}) (model_core.AppStore, error) {
	var data model_core.AppStore
	res := r.db.Model(&model_core.AppStore{}).Where(query).First(&data)
	return data, res.Error
}

// -----------------------------Channel LookupCodes-----------------------------------------------------------
func (r *cores) GetChannelLookupCodes(query interface{}, p *pagination.Paginatevalue) ([]model_core.ChannelLookupCodes, error) {
	var data []model_core.ChannelLookupCodes
	res := r.db.Model(&model_core.ChannelLookupCodes{}).Scopes(helpers.Paginate(&model_core.ChannelLookupCodes{}, p, r.db)).Where(query).Find(&data)
	return data, res.Error
}

//--------------------------meta data---------------------------------------------------------------------------

func (r *cores) ListMetaData() ([]model_core.IRModel, error) {
	var data []model_core.IRModel
	res := r.db.Model(&model_core.IRModel{}).Preload(clause.Associations).Find(&data)
	return data, res.Error
}
func (r *cores) ViewMetaData(model_name string) ([]model_core.IRModelFields, error) {
	var data []model_core.IRModelFields
	res := r.db.Model(&model_core.IRModelFields{}).Where("table_name = ?", model_name).Find(&data)
	return data, res.Error
}

// --------------------------install app------------------------------------------------------------------------------
func (r *cores) InstallApp(data *model_core.InstalledApps) error {
	res := r.db.Model(&model_core.InstalledApps{}).Create(&data)
	return res.Error
}
func (r *cores) CheckState(code string) bool {
	var resp model_core.InstalledApps
	res := r.db.Model(&model_core.InstalledApps{}).Where("code", code).First(&resp)
	return res.Error == nil
}

func (r *cores) UninstallApp(query map[string]interface{}) error {
	var data model_core.AppStore
	var install_data model_core.InstalledApps
	res := r.db.Model(&model_core.AppStore{}).Where(query).First(&data)
	if res.Error == nil {
		r.db.Model(&model_core.InstalledApps{}).Where("code", data.Code).Delete(&install_data)
	}
	return res.Error
}

// ------------------------core Company----------------------------------------------------------------
func (r *cores) CreateCompany(data *model_core.Company) error {
	err := r.db.Create(data).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *cores) UpdateCompany(query map[string]interface{}, data *model_core.Company) error {
	err := r.db.Model(&model_core.Company{}).Where(query).Updates(data).Error
	return err
}
func (r *cores) FindCompany(query map[string]interface{}) (model_core.Company, error) {
	var data model_core.Company
	err := r.db.Model(&model_core.Company{}).Where(query).First(&data)
	if err.Error != nil {
		return data, err.Error
	}
	return data, nil
}
