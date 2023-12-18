package server

import (
	"errors"
	"fmt"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

const passwordRequirements = "Passoword must be at least 12 characters long," +
	" contain a lower, uppercase letters, a digit, and a special charater."

var errInternalError = errors.New("internal server error")

func (s *Server) validateForm(form any) []error {
	err := s.validator.Struct(form)
	if err == nil {
		return nil
	}

	_, ok := err.(*validator.InvalidValidationError)
	if ok {
		s.router.Logger.Error(err)
		return []error{errInternalError}
	}

	verrs, ok := err.(validator.ValidationErrors)
	if !ok {
		s.router.Logger.Error(err)
		return []error{errInternalError}
	}

	errs := make([]error, 0, len(verrs))
	for _, verr := range verrs {
		errs = append(errs, convertFieldValidationError(verr))
	}
	return errs
}

func convertFieldValidationError(verr validator.FieldError) error {
	switch verr.Field() {
	case "Name":
		if verr.Tag() == "max" {
			return errors.New("Name: must be at most 16 characters long")
		} else if verr.Tag() == "min" {
			return errors.New("Name: must be at least 4 characters long")
		}
	case "Email":
		if verr.Tag() == "required" {
			return errors.New("Email: can't be empty")
		} else if verr.Tag() == "email" {
			return errors.New("Email: invalid format")
		}
	case "Password":
		if verr.Tag() == "required" {
			return errors.New("Password: can't be empty")
		} else if verr.Tag() == "password" {
			return errors.New(passwordRequirements)
		}
	case "PasswordAgain":
		if verr.Tag() == "required_with=Password|eqfield=Password" {
			return errors.New("Passwords mismatch")
		}
	}
	return fmt.Errorf("unexpected validation error: %w", verr)
}

func validatePassword(fl validator.FieldLevel) bool {
	return true
	// TODO: refactor into a loop
	pass := fl.Field().String()
	return len(pass) > 12 &&
		strings.ContainsAny(pass, "!@#$%^*()`") &&
		strings.ContainsFunc(pass, func(r rune) bool {
			return unicode.Is(unicode.Upper, r)
		}) &&
		strings.ContainsFunc(pass, func(r rune) bool {
			return unicode.Is(unicode.Lower, r)
		})
}
