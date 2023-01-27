package rating

import (
	"fermion/backend_core/internal/model/core"

	"gorm.io/datatypes"
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
type FeedBackFormDto struct {
	Id                 uint  `json:"id"`
	FeedBackCategoryId *uint `json:"feedback_category_id" `
	// FeedBackCategory   FeedbackCategory `json:"feedback_category" gorm:"foreignkey:FeedBackCategoryId; references:ID"`
	FeedBackForm []FeedBack  `json:"feedback_form"`
	RatingValue  int64       `json:"rating_value"`
	FeedBackUrl  FeedBackUrl `json:"feedback_url"`
}

type FeedBackCategoryDTO struct {
	Id   uint   `json:"id"`
	Name string `json:"rating_category_name"`
}
type FeedBackUrl struct {
	Url      string `json:"linked_url"`
	TlMethod string `json:"tl_method"`
	Params   Params `json:"params"`
}

type Params struct {
	ID string `json:"id"`
	// FeedBackId      string `json:"feedback_id"`
	AdditionalProp1 string `json:"additional_prop1"`
	AdditionalProp2 string `json:"additional_prop2"`
	AdditionalProp3 string `json:"additional_prop3"`
}

type FeedBack struct {
	ID           uint            `json:"id"`
	ParentId     *uint           `json:"parent_id"`
	Question     string          `json:"question"`
	AnswerTypeId *uint           `json:"answer_type_id"`
	Answer       string          `json:"answer"`
	AnswerType   core.Lookupcode `json:"answer_type" `
}

type RatingResponseDTO struct {
	Id               uint `json:"id"`
	RelatedId        uint `json:"related_id" gorm:"column:related_id"`
	RatingCategoryId uint `json:"rating_category_id"`
	// RatingCategory RatingCategoryDto `json:"rating_category"`
	RatingValue  int64          `json:"rating_value" gorm:"column:rating_value"`
	FeedBackForm datatypes.JSON `json:"feedback_form"`
}
type RatingRequestDTO struct {
	Id uint `json:"id"`
	// RelatedId        uint              `json:"related_id" gorm:"column:related_id"`
	// RatingCategoryId uint              `json:"rating_category_id"`
	RatingCategory string `json:"rating_category"`
	// RatingCategory RatingCategoryDto `json:"rating_category"`
	RatingValue  int64          `json:"rating_value" gorm:"column:rating_value"`
	FeedBackForm datatypes.JSON `json:"feedback_form"`
	FeedBackId   uint           `json:"feedback_id"`
}
type RatingCategoryDto struct {
	Id    uint   `json:"id"`
	Name  string `json:"category_name" gorm:"column:category_name"`
	Value int64  `json:"rating_value" gorm:"column:rating_value"`
}
