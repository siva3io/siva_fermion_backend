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
type RatingCategory struct {
	core.Model
	Name  string `json:"category_name" gorm:"column:category_name"`
	Value int64  `json:"rating_value" gorm:"column:rating_value"`
}
type FeedbackCategory struct {
	core.Model
	Name string `json:"rating_category_name" gorm:"column:rating_category_name"`
}

type FeedbackForm struct {
	core.Model
	FeedBackCategoryId *uint            `json:"feedback_category_id" gorm:"column:feedback_category_id"`
	FeedBackCategory   FeedbackCategory `json:"feedback_category" gorm:"foreignkey:FeedBackCategoryId; references:ID"`
	FeedBackForm       []Feedback       `json:"feedback_form,omitempty" gorm:"foreignkey:FeedbackFormId ;references:ID; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	RatingValue        int64            `json:"rating_value" gorm:"column:rating_value"`
	FeedBackUrl        FeedBackUrl      `json:"feedback_url" gorm:"embedded"`
}

type Rating struct {
	core.Model
	RelatedId uint `json:"related_id" gorm:"column:related_id"`
	//	RelatedType	 string          `json:"related_type"`
	RatingCategoryId uint           `json:"rating_category_id"`
	RatingCategory   RatingCategory `json:"rating_category" gorm:"foreignkey:RatingCategoryId; references:ID" `
	RatingValue      int64          `json:"rating_value" gorm:"column:rating_value"`
	FeedBack         datatypes.JSON `json:"feedback"`
}
type FeedBackUrl struct {
	Url      string `json:"linked_url" gorm:"column:linked_url"`
	TlMethod string `json:"tl_method" gorm:"column:tl_method"`
	Params   Params `json:"params" gorm:"embedded"`
}
type Params struct {
	// FeedBackId      string `json:"feedback_id" gorm:"column:feedback_id"`
	AdditionalProp1 string `json:"additional_prop1" gorm:"column:additional_prop1"`
	AdditionalProp2 string `json:"additional_prop2" gorm:"column:additional_prop2"`
	AdditionalProp3 string `json:"additional_prop3" gorm:"column:additional_prop3"`
}
type Feedback struct {
	core.Model
	FeedbackFormId uint            `json:"feedback_form_id" gorm:"column:feedback_form_id"`
	ParentId       *uint           `json:"parent_id" gorm:"column:parent_id;"`
	ParentFeedback *Feedback       `json:"parent_feedback" gorm:"foreignkey:ParentId; references:ID"`
	Question       string          `json:"question" gorm:"column:question"`
	AnswerTypeId   *uint           `json:"answer_type_id" gorm:"column:answer_type_id"`
	AnswerType     core.Lookupcode `json:"answer_type" gorm:"foreignkey:AnswerTypeId; references:ID"`
}

type Report struct {
	core.Model
	ReportData datatypes.JSON `json:"report_data"`
}
