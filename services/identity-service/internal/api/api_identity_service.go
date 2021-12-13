/*
 * Identity Service
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"identity-service/internal/users"

	"github.com/muonsoft/validation/validator"
)

// IdentityApiService is a service that implents the logic for the IdentityApiServicer
// This service should implement the business logic for every endpoint for the IdentityApi API.
// Include any external packages or services that will be required by this service.
type IdentityApiService struct {
	users          users.Repository
	passwordHasher PasswordHasher
	tokenIssuer    TokenIssuer
}

// NewIdentityApiService creates a default api service
func NewIdentityApiService(
	users users.Repository,
	passwordHasher PasswordHasher,
	tokenIssuer TokenIssuer,
) IdentityApiServicer {
	return &IdentityApiService{users: users, passwordHasher: passwordHasher, tokenIssuer: tokenIssuer}
}

// GetCurrentUser -
func (s *IdentityApiService) GetCurrentUser(ctx context.Context, id int64) (ImplResponse, error) {
	user, err := s.users.FindByID(ctx, id)
	if errors.Is(err, users.ErrUserNotFound) {
		return Response(http.StatusForbidden, Error{
			Code:    http.StatusForbidden,
			Message: "access denied",
		}), nil
	}
	if err != nil {
		return Response(http.StatusInternalServerError, nil), err
	}

	return Response(200, user), nil
}

// LoginUser - LoginForm to the system
func (s *IdentityApiService) LoginUser(ctx context.Context, form LoginForm) (ImplResponse, error) {
	user, err := s.users.FindByEmail(ctx, form.Email)
	if errors.Is(err, users.ErrUserNotFound) {
		return Response(http.StatusUnprocessableEntity, Error{
			Code:    http.StatusUnprocessableEntity,
			Message: "email or password is incorrect",
		}), nil
	}
	if err != nil {
		return Response(http.StatusInternalServerError, nil), err
	}

	isValid, err := s.passwordHasher.Verify(form.Password, user.Password)
	if err != nil {
		return Response(http.StatusInternalServerError, nil), err
	}
	if !isValid {
		return Response(http.StatusUnprocessableEntity, Error{
			Code:    http.StatusUnprocessableEntity,
			Message: "email or password is incorrect",
		}), nil
	}

	token, err := s.tokenIssuer.Issue(user)
	if err != nil {
		return Response(http.StatusInternalServerError, nil), err
	}

	return Response(200, LoginResponse{AccessToken: token}), nil
}

// LogoutUser - Logout from the system
func (s *IdentityApiService) LogoutUser(ctx context.Context) (ImplResponse, error) {
	return Response(204, nil), nil
}

// RegisterUser - RegistrationForm a user
func (s *IdentityApiService) RegisterUser(ctx context.Context, form RegistrationForm) (ImplResponse, error) {
	err := validator.ValidateValidatable(ctx, form)
	if err != nil {
		return Response(
			http.StatusUnprocessableEntity,
			Error{Code: http.StatusUnprocessableEntity, Message: err.Error()},
		), nil
	}

	count, err := s.users.CountByEmail(ctx, form.Email)
	if err != nil {
		return Response(http.StatusInternalServerError, nil), err
	}
	if count > 0 {
		return Response(
			http.StatusUnprocessableEntity,
			Error{
				Code:    http.StatusUnprocessableEntity,
				Message: fmt.Sprintf(`user with email "%s" already exists`, form.Email),
			},
		), nil
	}

	password, err := s.passwordHasher.Hash(form.Password)
	if err != nil {
		return Response(http.StatusInternalServerError, nil), err
	}

	user := &users.User{
		Email:     form.Email,
		Password:  password,
		FirstName: form.FirstName,
		LastName:  form.LastName,
	}

	err = s.users.Save(ctx, user)
	if err != nil {
		return Response(http.StatusInternalServerError, user), err
	}

	return Response(http.StatusCreated, user), nil
}

// UpdateCurrentUser -
func (s *IdentityApiService) UpdateCurrentUser(ctx context.Context, id int64, form UpdateForm) (ImplResponse, error) {
	err := validator.ValidateValidatable(ctx, form)
	if err != nil {
		return Response(
			http.StatusUnprocessableEntity,
			Error{Code: http.StatusUnprocessableEntity, Message: err.Error()},
		), nil
	}

	user, err := s.users.FindByID(ctx, id)
	if errors.Is(err, users.ErrUserNotFound) {
		return Response(http.StatusForbidden, Error{
			Code:    http.StatusForbidden,
			Message: "access denied",
		}), nil
	}
	if err != nil {
		return Response(http.StatusInternalServerError, nil), err
	}

	user.FirstName = form.FirstName
	user.LastName = form.LastName
	user.Phone = form.Phone

	err = s.users.Save(ctx, user)
	if err != nil {
		return Response(http.StatusInternalServerError, user), err
	}

	return Response(http.StatusOK, user), nil
}
