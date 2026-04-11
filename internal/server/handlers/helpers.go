package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

var validate *validator.Validate

var keyNameRegex = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

func init() {
	validate = validator.New()
	validate.RegisterValidation("keyname", func(fl validator.FieldLevel) bool {
		return keyNameRegex.MatchString(fl.Field().String())
	})
}

func RespondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func RespondError(w http.ResponseWriter, status int, message string) {
	if status >= 500 {
		log.Error().Int("status", status).Str("error", message).Msg("server error")
	}
	RespondJSON(w, status, map[string]string{"error": message})
}

// DecodeAndValidate decodes request body JSON and validates struct tags.
func DecodeAndValidate(r *http.Request, dst interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}
	if err := validate.Struct(dst); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			msgs := make([]string, 0, len(validationErrs))
			for _, e := range validationErrs {
				msgs = append(msgs, fmt.Sprintf("%s: failed on '%s'", e.Field(), e.Tag()))
			}
			return fmt.Errorf("validation failed: %s", strings.Join(msgs, ", "))
		}
		return err
	}
	return nil
}
