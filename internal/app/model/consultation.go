package model

import (
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type Consultation struct {
	ID             int           `json:"id"`
	NadzonOrganID  int           `json:"nadzor_organ_id" form:"nadzor_organ_id" validate:"required"`
	ControlTypeID  int           `json:"control_type_id" form:"control_type_id" validate:"required"`
	ConsultTopicID int           `json:"consult_topic_id" form:"consult_topic_id" validate:"required"`
	SlotID         int           `json:"slot_id" form:"slot_id" validate:"required"`
	UserID         int           `json:"user_id" form:"user_id"`
	User           *Account      `json:"user,omitempty"`
	Time           string        `json:"time" form:"time" validate:"required"`
	Date           time.Time     `json:"-" form:"date" validate:"required"`
	DateExport     string        `json:"date"`
	Question       string        `json:"question" form:"question"`
	IsNeedLetter   bool          `json:"is_need_letter" form:"is_need_letter"`
	IsConfirmed    bool          `json:"is_confirmed" form:"is_confirmed"`
	VKSLink        string        `json:"vks_link" form:"vks_link"`
	VideoLink      string        `json:"video_link" form:"video_link"`
	IsDeleted      bool          `json:"is_deleted" form:"is_deleted"`
	Answer         string        `json:"answer" form:"answer"`
	NadzonOrgan    *NadzonOrgan  `json:"nadzor_organ,omitempty"`
	ControlType    *ControlType  `json:"control_type,omitempty"`
	ConsultTopic   *ConsultTopic `json:"consult_topic,omitempty"`
}

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

type Consultations struct {
	Active   []Consultation `json:"active"`
	Finished []Consultation `json:"finished"`
}
