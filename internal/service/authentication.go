package service

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/repository"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/constant"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/util"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthenticationService struct {
	log        *zap.Logger
	repository repository.IAuthenticationRepository
}

type IAuthenticationService interface {
	SignUp(request dto.SignUpRequest) (dto.SignUpResponse, error)
	SignIn(request dto.SignInRequest) (dto.SignInResponse, error)
	ForgotPassword(request dto.ForgotPasswordRequest) (dto.ForgotPasswordResponse, error)
	ChangePassword(request dto.ChangePasswordRequest, accountId uuid.UUID) (dto.ChangePasswordResponse, error)
}

func NewAuthenticationService(log *zap.Logger, repository repository.IAuthenticationRepository) IAuthenticationService {
	return &AuthenticationService{
		log:        log,
		repository: repository,
	}
}

func (a *AuthenticationService) SignUp(request dto.SignUpRequest) (dto.SignUpResponse, error) {
	a.repository.UseTx(true)

	var (
		err error
	)

	defer func() {
		if err != nil {
			if err := a.repository.Rollback(); err != nil {
				a.log.Error("[SignUp] failed to rollback transaction", zap.Error(err))
			}
			return
		}
	}()

	Id, err := uuid.NewUUID()
	if err != nil {
		a.log.Error("[SignUp] failed to generate UUID", zap.Error(err))
		return dto.SignUpResponse{}, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		a.log.Error("[SignUp] failed to hash password", zap.Error(err))
		return dto.SignUpResponse{}, err
	}

	account := entity.Account{
		Id:           Id,
		Email:        request.Email,
		Name:         request.Name,
		Password:     string(hashedPassword),
		RoleId:       request.RoleId,
		PhotoProfile: "",
	}

	// Todo : create staff

	if err := a.repository.CreateAccount(&account); err != nil {
		a.log.Error("[SignUp] failed to create account", zap.Error(err))
		return dto.SignUpResponse{}, err
	}

	if err := a.repository.Commit(); err != nil {
		a.log.Error("[SignUp] failed to commit transaction", zap.Error(err))
	}

	return dto.SignUpResponse{
		Id:           account.Id.String(),
		Name:         account.Name,
		PhotoProfile: account.PhotoProfile,
		Email:        account.Email,
		Role: dto.RoleResponse{
			Id:   account.Role.Id,
			Name: account.Role.Name,
		},
	}, nil
}

func (a *AuthenticationService) SignIn(request dto.SignInRequest) (dto.SignInResponse, error) {
	a.repository.UseTx(false)

	account, err := a.repository.GetAccountByEmail(request.Email)
	if err != nil {
		a.log.Error("[SigIn] failed to get account by email", zap.Error(err))
		return dto.SignInResponse{}, nil
	}

	if bcrypt.CompareHashAndPassword([]byte(request.Password), []byte(account.Password)) != nil {
		a.log.Error("[SignIn] password is incorrect")
		return dto.SignInResponse{}, nil
	}

	return dto.SignInResponse{
		TokenType:   constant.TokenType,
		AccessToken: "",
		ExpiredAt:   0,
	}, nil
}

func (a *AuthenticationService) ForgotPassword(request dto.ForgotPasswordRequest) (dto.ForgotPasswordResponse, error) {
	a.repository.UseTx(false)

	account, err := a.repository.GetAccountByEmail(request.Email)
	if err != nil {
		a.log.Error("[ForgotPassword] failed to get account by email", zap.Error(err))
		return dto.ForgotPasswordResponse{}, nil
	}

	tempPassord, err := util.RandomString(8)
	if err != nil {
		a.log.Error("[ForgotPassword] failed to generate random string", zap.Error(err))
		return dto.ForgotPasswordResponse{}, nil
	}

	// Todo : send email
	fmt.Printf("Send email to %s with temp password %s", account.Email, tempPassord)

	return dto.ForgotPasswordResponse{
		Id:    account.Id.String(),
		Email: account.Email,
	}, nil
}

func (a *AuthenticationService) ChangePassword(request dto.ChangePasswordRequest, accountId uuid.UUID) (dto.ChangePasswordResponse, error) {
	a.repository.UseTx(false)

	account, err := a.repository.GetAccountById(accountId)
	if err != nil {
		a.log.Error("[ChangePassword] failed to get account by id", zap.Error(err))
		return dto.ChangePasswordResponse{}, nil
	}

	if bcrypt.CompareHashAndPassword([]byte(request.OldPassword), []byte(account.Password)) != nil {
		a.log.Error("[ChangePassword] password is incorrect")
		return dto.ChangePasswordResponse{}, nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		a.log.Error("[ChangePassword] failed to hash password", zap.Error(err))
		return dto.ChangePasswordResponse{}, nil
	}

	account.Password = string(hashedPassword)

	if err := a.repository.UpdateAccount(&account); err != nil {
		a.log.Error("[ChangePassword] failed to update account", zap.Error(err))
		return dto.ChangePasswordResponse{}, nil
	}

	return dto.ChangePasswordResponse{
		Id:           account.Id.String(),
		Name:         account.Name,
		PhotoProfile: account.PhotoProfile,
		Email:        account.Email,
		Role: dto.RoleResponse{
			Id:   account.Role.Id,
			Name: account.Role.Name,
		},
	}, nil
}
