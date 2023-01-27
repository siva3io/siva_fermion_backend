package rating

import (
	"errors"
	"os"
	"time"

	"fermion/backend_core/db"
	"fermion/backend_core/internal/model/pagination"
	model "fermion/backend_core/internal/model/rating"
	"fermion/backend_core/pkg/util/helpers"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/*
Copyright (C) 2022 Eunimart Omnichannel Pvt Ltd. (www.eunimart.com)
All rights reserved.
This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Lesser General Public License v3.0 as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.
This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Lesser General Public License v3.0 for more details.
You should have received a copy of the GNU Lesser General Public License v3.0
along with this program.  If not, see <https://www.gnu.org/licenses/lgpl-3.0.html/>.
*/

type RatingRepo interface {
	SaveRatingTable(data *model.Rating) error
	UpdateRating(query map[string]interface{}, data *model.Rating) error
	GetRating(query map[string]interface{}) (model.Rating, error)
	GetRatingList(query map[string]interface{}, p *pagination.Paginatevalue) ([]model.Rating, error)
	DeleteRating(query map[string]interface{}) error

	CreateFeedBackForm(data *model.FeedbackForm) error
	RatingCategoriesList(query map[string]interface{}, p *pagination.Paginatevalue) ([]model.RatingCategory, error)
	FeedBackCategoryList(query map[string]interface{}, p *pagination.Paginatevalue) ([]model.FeedbackCategory, error)
	GetFeedBackCategory(query map[string]interface{}) (model.FeedbackCategory, error)
	GetFeedBacKForm(query map[string]interface{}) (model.FeedbackForm, error)

	CreateFeedback(data *model.Feedback) error
	UpdateFeedback(query map[string]interface{}, data *model.Feedback) error
	GetFeedback(query map[string]interface{}) (model.Feedback, error)
	GetFeedbackList(query map[string]interface{}, p *pagination.Paginatevalue) ([]model.Feedback, error)
	DeleteFeedback(query map[string]interface{}) error

	SaveReport(data *model.Report) error
	UpdateReport(query map[string]interface{}, data *model.Report) error
	GetReport(query map[string]interface{}) (model.Report, error)
	GetReportList(query map[string]interface{}, p *pagination.Paginatevalue) ([]model.Report, error)
	DeleteReport(query map[string]interface{}) error

	ListFeedBackForm(query map[string]interface{}, p *pagination.Paginatevalue) (model.FeedbackForm, error)
}

type rating struct {
	db *gorm.DB
}

var ratingRespository *rating

// singleton function
func NewRating() *rating {
	if ratingRespository != nil {
		return ratingRespository
	}
	db := db.DbManager()
	ratingRespository = &rating{db}
	return ratingRespository

}

func (r *rating) RatingCategoriesList(query map[string]interface{}, p *pagination.Paginatevalue) ([]model.RatingCategory, error) {

	//err:= r.db.Model(RatingCategories{})
	var data []model.RatingCategory
	err := r.db.Model(&model.RatingCategory{}).Scopes(helpers.Paginate(&model.RatingCategory{}, p, r.db)).Where(query).Find(&data).Error
	if err != nil {
		return data, err
	}

	return data, nil
}

func (r *rating) FeedBackCategoryList(query map[string]interface{}, p *pagination.Paginatevalue) ([]model.FeedbackCategory, error) {
	var data []model.FeedbackCategory
	err := r.db.Model(&model.FeedbackCategory{}).Scopes(helpers.Paginate(&model.FeedbackCategory{}, p, r.db)).Where(query).Find(&data)
	if err != nil {
		return data, err.Error
	}

	return data, nil
}

func (r *rating) SaveRatingTable(data *model.Rating) error {
	err := r.db.Model(&model.Rating{}).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *rating) UpdateRating(query map[string]interface{}, data *model.Rating) error {

	err := r.db.Model(&model.Rating{}).Where(query).Updates(data)
	if err.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if err != nil {
		return err.Error
	}
	return nil
}

func (r *rating) GetRating(query map[string]interface{}) (model.Rating, error) {
	var data model.Rating
	err := r.db.Model(&model.Rating{}).Where(query).First(&data)
	if err.RowsAffected == 0 {
		return data, errors.New("oops! record not found")
	}
	if err != nil {
		return data, err.Error
	}
	return data, nil
}

func (r *rating) GetRatingList(query map[string]interface{}, p *pagination.Paginatevalue) ([]model.Rating, error) {
	var data []model.Rating
	err := r.db.Preload(clause.Associations).Model(&model.Rating{}).Scopes(helpers.Paginate(&model.Rating{}, p, r.db)).Where(query).Find(&data).Error
	if err != nil {
		return data, err
	}
	return data, nil
}

func (r *rating) DeleteRating(query map[string]interface{}) error {
	timeZone := os.Getenv("DB_TZ")
	timeLocation, _ := time.LoadLocation(timeZone)
	data := map[string]interface{}{
		"deleted_by": query["user_id"],
		"deleted_at": time.Now().In(timeLocation),
	}
	delete(query, "user_id")
	res := r.db.Model(&model.Rating{}).Where(query).Updates(data)
	if res.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r *rating) CreateFeedBackForm(data *model.FeedbackForm) error {
	err := r.db.Model(&model.FeedbackForm{}).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *rating) GetFeedBackCategory(query map[string]interface{}) (model.FeedbackCategory, error) {
	var data model.FeedbackCategory
	err := r.db.Model(&model.FeedbackCategory{}).Where(query).First(&data)
	if err.RowsAffected == 0 {
		return data, errors.New("oops! record not found")
	}
	if err != nil {
		return data, err.Error
	}
	return data, nil
}

func (r *rating) GetFeedBacKForm(query map[string]interface{}) (model.FeedbackForm, error) {
	var data model.FeedbackForm
	err := r.db.Model(&model.FeedbackForm{}).Where(query).First(&data)
	if err.RowsAffected == 0 {
		return data, errors.New("oops! record not found")
	}
	if err != nil {
		return data, err.Error
	}
	return data, nil
}

func (r *rating) CreateFeedback(data *model.Feedback) error {
	err := r.db.Model(&model.Feedback{}).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *rating) GetFeedback(query map[string]interface{}) (model.Feedback, error) {
	var data model.Feedback
	err := r.db.Model(&model.Feedback{}).Where(query).First(&data)
	if err.RowsAffected == 0 {
		return data, errors.New("oops! record not found")
	}
	if err != nil {
		return data, err.Error
	}
	return data, nil
}

func (r *rating) GetFeedbackList(query map[string]interface{}, p *pagination.Paginatevalue) ([]model.Feedback, error) {
	var data []model.Feedback
	err := r.db.Preload(clause.Associations).Model(&model.Feedback{}).Scopes(helpers.Paginate(&model.Feedback{}, p, r.db)).Where(query).Find(&data).Error
	if err != nil {
		return data, err
	}
	return data, nil
}

func (r *rating) UpdateFeedback(query map[string]interface{}, data *model.Feedback) error {

	err := r.db.Model(&model.Feedback{}).Where(query).Updates(data)
	if err.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if err != nil {
		return err.Error
	}
	return nil
}

func (r *rating) DeleteFeedback(query map[string]interface{}) error {
	timeZone := os.Getenv("DB_TZ")
	timeLocation, _ := time.LoadLocation(timeZone)
	data := map[string]interface{}{
		"deleted_by": query["user_id"],
		"deleted_at": time.Now().In(timeLocation),
	}
	delete(query, "user_id")
	res := r.db.Model(&model.Feedback{}).Where(query).Updates(data)
	if res.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r *rating) SaveReport(data *model.Report) error {
	err := r.db.Model(&model.Report{}).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *rating) UpdateReport(query map[string]interface{}, data *model.Report) error {

	err := r.db.Model(&model.Report{}).Where(query).Updates(data)
	if err.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if err != nil {
		return err.Error
	}
	return nil
}

func (r *rating) GetReport(query map[string]interface{}) (model.Report, error) {
	var data model.Report
	err := r.db.Model(&model.Report{}).Where(query).First(&data)
	if err.RowsAffected == 0 {
		return data, errors.New("oops! record not found")
	}
	if err != nil {
		return data, err.Error
	}
	return data, nil
}

func (r *rating) GetReportList(query map[string]interface{}, p *pagination.Paginatevalue) ([]model.Report, error) {
	var data []model.Report
	err := r.db.Model(&model.Report{}).Scopes(helpers.Paginate(&model.Report{}, p, r.db)).Where(query).Find(&data).Error
	if err != nil {
		return data, err
	}
	return data, nil
}

func (r *rating) DeleteReport(query map[string]interface{}) error {
	timeZone := os.Getenv("DB_TZ")
	timeLocation, _ := time.LoadLocation(timeZone)
	data := map[string]interface{}{
		"deleted_by": query["user_id"],
		"deleted_at": time.Now().In(timeLocation),
	}
	delete(query, "user_id")
	res := r.db.Model(&model.Report{}).Where(query).Updates(data)
	if res.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r *rating) ListFeedBackForm(query map[string]interface{}, p *pagination.Paginatevalue) (model.FeedbackForm, error) {
	var data model.FeedbackForm
	err := r.db.Preload(clause.Associations + "." + clause.Associations + "." + clause.Associations).Model(&model.FeedbackForm{}).Scopes(helpers.Paginate(&model.FeedbackForm{}, p, r.db)).Where(query).First(&data).Error
	if err != nil {
		return data, err
	}
	return data, nil
}
