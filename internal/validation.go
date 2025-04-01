package internal

import (
	"errors"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/gofiber/fiber/v2"
)

var (
	nameRegex      = regexp.MustCompile(`^[a-zA-Z0-9 ]+$`)
	placementRegex = regexp.MustCompile(`^[a-zA-Z]+(_[a-zA-Z]+)*$`)
	categoryRegex  = regexp.MustCompile(`^[a-zA-Z]+(_[a-zA-Z]+)*$`)
	keywordRegex   = regexp.MustCompile(`^[a-zA-Z]+(_[a-zA-Z]+)*$`)
)

var (
	// NameRules the name validation rules.
	NameRules = []validation.Rule{
		validation.Length(2, 50),
		validation.Match(nameRegex),
	}
	// AdvertiserIDRules the advertiser ID rules.
	AdvertiserIDRules = []validation.Rule{
		validation.Length(2, 10),
		is.Alphanumeric,
	}
	// PlacementRules the placement validation rules.
	PlacementRules = []validation.Rule{
		validation.Length(2, 20),
		validation.Match(placementRegex),
	}
	// CategoryRules the category validation rules.
	CategoryRules = []validation.Rule{
		validation.Length(2, 20),
		validation.Match(categoryRegex),
	}
	// KeywordRules the keyword validation rules.
	KeywordRules = []validation.Rule{
		validation.Length(2, 20),
		validation.Match(keywordRegex),
	}
)

// ValidationErrorJSON returns the validation error details (like field and reason) as JSON (fiber.Map).
func ValidationErrorJSON(err error) fiber.Map {
	var (
		errs     validation.Errors
		field    string
		fieldErr error
	)

	// Check is validation error result
	if !errors.As(err, &errs) {
		return nil
	}

	for field, fieldErr = range errs {
		break
	}

	return fiber.Map{
		"field":  field,
		"reason": fieldErr.Error(),
	}
}
