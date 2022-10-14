package resources

import (
	"fmt"
	"github.com/go-playground/validator/v10"
)

type Customer struct {
	ID          string `csv:"ref" validate:"required"`                                                         // reference the user’s system.
	Name        string `csv:"name" validate:"required"`                                                        //full name.
	Email       string `csv:"email" validate:"required,email"`                                                 //valid email address.
	Address     string `csv:"address,omitempty" json:"Address,omitempty"`                                      //the customer’s real world address used for shipping.
	CountryCode string `csv:"country_code,omitempty" json:"CountryCode,omitempty" validate:"iso3166_1_alpha2"` // two-letter country code for the customer.
}

func (c *Customer) NewZeroValue() Resource {
	return &Customer{}
}

func (c *Customer) Validate() error {
	validate := validator.New()
	if err := validate.Struct(c); err != nil {
		var errors ValidateErrors
		for _, err := range err.(validator.ValidationErrors) {
			var e ValidateError
			e.Field = err.Field()
			e.Tag = err.Tag()
			e.Value = err.Value()
			e.Record = fmt.Sprintf("%+v", c)
			errors = append(errors, e)
		}
		return errors
	}
	return nil
}
