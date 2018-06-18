package dal_pivot

import (
	"time"

	"github.com/ghetzel/pivot/dal"
	"github.com/ghetzel/pivot/mapper"
)

var Widgets mapper.Mapper

var WidgetsSchema = &dal.Collection{
	Name:                   `widgets`,
	IdentityFieldType:      dal.StringType,
	IdentityFieldFormatter: dal.GenerateUUID,
	Fields: []dal.Field{
		{
			Name:        `type`,
			Description: `The type of widget.`,
			Type:        dal.StringType,
			Validator:   dal.ValidateIsOneOf(`foo`, `bar`, `baz`),
			Required:    true,
		}, {
			Name:        `usage`,
			Description: `Short description on how to use this widget.`,
			Type:        dal.StringType,
		}, {
			Name:        `created_at`,
			Description: `When the widget was created.`,
			Type:        dal.TimeType,
			Formatter:   dal.CurrentTimeIfUnset,
		}, {
			Name:        `updated_at`,
			Description: `Last time the widget was updated.`,
			Type:        dal.TimeType,
			Formatter:   dal.CurrentTime,
		},
	},
}

type Widget struct {
	ID        string    `pivot:"id,identity"`
	Type      string    `pivot:"type"`
	Usage     string    `pivot:"usage"`
	CreatedAt time.Time `pivot:"created_at"`
	UpdatedAt time.Time `pivot:"updated_at"`
}
