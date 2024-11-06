package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	db "github.com/devphaseX/go-microservice/authenication-service/db/sqlc"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type SigninParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		ID:        uuid.UUID(user.ID.Bytes).String(),
		FirstName: user.FirstName.String,
		LastName:  user.LastName.String,
		Email:     user.Email,
		UpdatedAt: user.UpdatedAt.Time,
		CreatedAt: user.CreatedAt.Time,
	}
}

type signinResponse struct {
	AccessToken          string       `json:"access_token"`
	AccessTokenExpiredAt time.Time    `json:"access_token_expired_at"`
	User                 userResponse `json:"user"`
}

func (p *SigninParams) Validate() error {
	err := validation.ValidateStruct(p,
		validation.Field(&p.Email, validation.Required, is.Email),
		validation.Field(&p.Password, validation.Required, is.Alphanumeric, validation.Min(8)),
	)

	return NewError(err, ErrBadRequest)
}

func (c *Config) Signin(w http.ResponseWriter, r *http.Request) {
	var signinParams SigninParams

	if err := json.NewDecoder(r.Body).Decode(&signinParams); err != nil {
		c.errorJSON(w, NewError(err, ErrBadRequest))
		return
	}

	user, err := c.store.GetUserByUnqiueID(context.Background(), db.GetUserByUnqiueIDParams{
		Email: pgtype.Text{
			Valid:  true,
			String: signinParams.Email,
		},
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.errorJSON(w, NewError(errors.New("invalid credential"), ErrNotFound))
			return
		}
	}

	passwordByte, passwordSaltByte := ArgonStringDecode(user.Password, user.PasswordSalt)

	if err := c.hash.Compare(passwordByte, passwordSaltByte, []byte(signinParams.Password)); err != nil {
		c.errorJSON(w, NewError(errors.New("invalid credential"), ErrNotFound))
		return
	}

	accessToken, payload, err := c.tokenMaker.CreateToken(uuid.UUID(user.ID.Bytes).String(), user.Email, c.env.AccessTokenTime)

	if err != nil {
		c.errorJSON(w, NewError(errors.New("failed to sign in user"), ErrInternalFailure))
		return
	}

	resp := signinResponse{
		AccessToken:          accessToken,
		AccessTokenExpiredAt: payload.ExpiresAt.Time,
		User:                 newUserResponse(user),
	}

	c.successJSON(w, http.StatusOK, resp, nil, "sign in successfully")
}

type RegisterUserParams struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (p *RegisterUserParams) Validate() error {
	err := validation.ValidateStruct(p,
		validation.Field(&p.FirstName, validation.Required, validation.Min(1)),
		validation.Field(&p.LastName, validation.Required, validation.Min(1)),
		validation.Field(&p.Email, validation.Required, is.Email),
		validation.Field(&p.Password, validation.Required, is.Alphanumeric, validation.Min(8)),
	)

	return NewError(err, ErrBadRequest)
}

func (c *Config) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var registerUserParams RegisterUserParams

	if err := json.NewDecoder(r.Body).Decode(&registerUserParams); err != nil {
		c.errorJSON(w, NewError(err, ErrBadRequest))
		return
	}

	hashSalt, err := c.hash.GenerateHash([]byte(registerUserParams.Password), nil)

	if err != nil {
		fmt.Println("hash salt err", err)
		c.errorJSON(w, NewError(errors.New("failed to create user"), ErrInternalFailure))
		return
	}

	passwordStr, passwordSaltStr := ArgonStringEncode(hashSalt)

	newUser, err := c.store.CreateUser(context.Background(), db.CreateUserParams{
		FirstName: pgtype.Text{
			String: registerUserParams.FirstName,
			Valid:  true,
		},

		LastName: pgtype.Text{
			String: registerUserParams.LastName,
			Valid:  true,
		},

		Email:        registerUserParams.Email,
		Password:     passwordStr,
		PasswordSalt: passwordSaltStr,
	})

	if err != nil {
		var dbErr *pgconn.PgError
		if errors.As(err, &dbErr) {
			if dbErr.ConstraintName == "users_email_key" {
				c.errorJSON(w, NewError(errors.New("email not available"), ErrConflictRequest))
				return
			}
		}
		fmt.Println("db err", err)

		c.errorJSON(w, NewError(errors.New("failed to create user"), ErrInternalFailure))
		return
	}

	_ = newUser

	c.successJSON(w, http.StatusCreated, nil, nil, "account created successfully")
}
