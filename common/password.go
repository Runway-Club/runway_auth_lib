package common

import (
	"github.com/Runway-Club/auth_lib/domain"
	"golang.org/x/crypto/bcrypt"
	"regexp"
)

type PasswordPolicy string

const (
	PasswordLevel1 PasswordPolicy = "level1"
	PasswordLevel2 PasswordPolicy = "level2"
	PasswordLevel3 PasswordPolicy = "level3"
)

func CheckPasswordPolicy(password string, policy string) error {
	// set default password policy
	policy = string(PasswordLevel1)
	// check password policy
	if policy == string(PasswordLevel1) {
		// minimum 8 characters
		if len(password) < 8 {
			return domain.ErrInvalidPassword
		}
		return nil
	}
	if policy == string(PasswordLevel2) {
		//  minimum 8 any characters and contain at least one number
		regex := regexp.MustCompile(`[0-9]`)
		if len(password) < 8 || !regex.MatchString(password) {
			return domain.ErrInvalidPassword
		}
		return nil
	}
	if policy == string(PasswordLevel3) {
		// minimum 8 any characters and contain at least one number and one uppercase letter and one special character
		regex := regexp.MustCompile(`[0-9]`)
		regex2 := regexp.MustCompile(`[A-Z]`)
		regex3 := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`)

		if len(password) < 8 || !regex.MatchString(password) || !regex2.MatchString(password) || !regex3.MatchString(password) {
			return domain.ErrInvalidPassword
		}
		return nil
	}
	return domain.ErrInvalidPasswordPolicy
}

func GeneratePassword(password string, hashCost string) (string, error) {
	// hash password
	hCost := 0
	if hashCost == "min" {
		hCost = bcrypt.MinCost
	}
	if hashCost == "max" {
		hCost = bcrypt.MaxCost
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), hCost)
	if err != nil {
		return "", domain.ErrInternal
	}
	return string(hashedPassword), nil
}
