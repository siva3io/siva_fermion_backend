package rating

import (
	"fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/pagination"
	model "fermion/backend_core/internal/model/rating"
	repo "fermion/backend_core/internal/repository/rating"
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
type Service interface {
	SaveRating(metaData core.MetaData, data *model.Rating) error
	SaveFeedBackForm(metaData core.MetaData, data *model.FeedbackForm) error
	RatingCategoriesList(metaData core.MetaData, p *pagination.Paginatevalue) (interface{}, error)
	FeedBackCategoryList(metaData core.MetaData, p *pagination.Paginatevalue) (interface{}, error)

	GetRatingList(metaData core.MetaData, p *pagination.Paginatevalue) (interface{}, error)
	UpdateRating(metaData core.MetaData, data *model.Rating) error
	GetRating(MetaData core.MetaData) (interface{}, error)
	DeleteRating(metaData core.MetaData) error

	SaveFeedBack(metaData core.MetaData, data *model.Feedback) error
	GetFeedbackList(metaData core.MetaData, p *pagination.Paginatevalue) (interface{}, error)
	UpdateFeedback(metaData core.MetaData, data *model.Feedback) error
	GetFeedback(MetaData core.MetaData) (interface{}, error)
	DeleteFeedback(metaData core.MetaData) error

	CreateReport(metaData core.MetaData, data *model.Report) error
	UpdateReport(metaData core.MetaData, data *model.Report) error
	GetReport(metaData core.MetaData) (interface{}, error)
	GetReportList(metaData core.MetaData, p *pagination.Paginatevalue) (interface{}, error)
	DeleteReport(metaData core.MetaData) error
	GetFeedbackCategory(metaData core.MetaData) (interface{}, error)
	ListFeedBackForm(metaData core.MetaData, p *pagination.Paginatevalue) (interface{}, error)
}
type service struct {
	rating_repo repo.RatingRepo
}

var newServiceObj *service

func NewService() *service {
	if newServiceObj != nil {
		return newServiceObj
	}
	newServiceObj = &service{
		repo.NewRating(),
	}
	return newServiceObj
}

func (s *service) SaveRating(metaData core.MetaData, data *model.Rating) error {
	data.CompanyId = metaData.CompanyId
	data.CreatedByID = &metaData.TokenUserId

	err := s.rating_repo.SaveRatingTable(data)
	if err != nil {
		return err
	}
	return nil

}

func (s *service) GetFeedbackCategory(metaData core.MetaData) (interface{}, error) {
	data, err := s.rating_repo.GetFeedBackCategory(metaData.Query)
	if err != nil {
		return data, err
	}
	return data, nil
}

func (s *service) UpdateRating(metaData core.MetaData, data *model.Rating) error {
	data.CompanyId = metaData.CompanyId
	data.UpdatedByID = &metaData.TokenUserId
	err := s.rating_repo.UpdateRating(metaData.Query, data)
	if err != nil {
		return err
	}
	return nil

}

func (s *service) GetRating(metaData core.MetaData) (interface{}, error) {
	result, err := s.rating_repo.GetRating(metaData.Query)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (s *service) GetRatingList(metaData core.MetaData, p *pagination.Paginatevalue) (interface{}, error) {
	result, err := s.rating_repo.GetRatingList(metaData.Query, p)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *service) DeleteRating(metaData core.MetaData) error {
	err := s.rating_repo.DeleteRating(metaData.Query)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) SaveFeedBackForm(metaData core.MetaData, data *model.FeedbackForm) error {
	data.CompanyId = metaData.CompanyId
	data.CreatedByID = &metaData.TokenUserId

	err := s.rating_repo.CreateFeedBackForm(data)
	if err != nil {
		return err
	}
	return nil

}

func (s *service) RatingCategoriesList(metaData core.MetaData, p *pagination.Paginatevalue) (interface{}, error) {
	result, err := s.rating_repo.RatingCategoriesList(metaData.Query, p)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *service) FeedBackCategoryList(metaData core.MetaData, p *pagination.Paginatevalue) (interface{}, error) {
	result, err := s.rating_repo.FeedBackCategoryList(metaData.Query, p)
	if err != nil {
		return result, err
	}
	return result, nil

}

func (s *service) SaveFeedBack(metaData core.MetaData, data *model.Feedback) error {
	data.CompanyId = metaData.CompanyId
	data.CreatedByID = &metaData.TokenUserId

	err := s.rating_repo.CreateFeedback(data)
	if err != nil {
		return err
	}
	return nil

}

func (s *service) UpdateFeedback(metaData core.MetaData, data *model.Feedback) error {
	data.CompanyId = metaData.CompanyId
	data.UpdatedByID = &metaData.TokenUserId
	err := s.rating_repo.UpdateFeedback(metaData.Query, data)
	if err != nil {
		return err
	}
	return nil

}

func (s *service) GetFeedback(metaData core.MetaData) (interface{}, error) {
	result, err := s.rating_repo.GetFeedback(metaData.Query)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *service) GetFeedbackList(metaData core.MetaData, p *pagination.Paginatevalue) (interface{}, error) {
	result, err := s.rating_repo.GetFeedbackList(metaData.Query, p)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *service) DeleteFeedback(metaData core.MetaData) error {
	err := s.rating_repo.DeleteFeedback(metaData.Query)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) CreateReport(metaData core.MetaData, data *model.Report) error {
	data.CompanyId = metaData.CompanyId
	data.CreatedByID = &metaData.TokenUserId

	err := s.rating_repo.SaveReport(data)
	if err != nil {
		return err
	}
	return nil

}

func (s *service) UpdateReport(metaData core.MetaData, data *model.Report) error {
	data.CompanyId = metaData.CompanyId
	data.UpdatedByID = &metaData.TokenUserId
	err := s.rating_repo.UpdateReport(metaData.Query, data)
	if err != nil {
		return err
	}
	return nil

}

func (s *service) GetReport(metaData core.MetaData) (interface{}, error) {
	result, err := s.rating_repo.GetReport(metaData.Query)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *service) GetReportList(metaData core.MetaData, p *pagination.Paginatevalue) (interface{}, error) {
	result, err := s.rating_repo.GetReportList(metaData.Query, p)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *service) DeleteReport(metaData core.MetaData) error {
	err := s.rating_repo.DeleteReport(metaData.Query)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) ListFeedBackForm(metaData core.MetaData, p *pagination.Paginatevalue) (interface{}, error) {
	result, err := s.rating_repo.ListFeedBackForm(metaData.Query, p)
	if err != nil {
		return result, err
	}
	return result, nil
}
