package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/infra/email"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/repository"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/constant"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/errx"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/jwt"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/util"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthenticationService struct {
	log        *zap.Logger
	email      *email.Email
	repository repository.IAuthenticationRepository
}

type IAuthenticationService interface {
	SignUp(request dto.SignUpRequest, accoundId uuid.UUID) (dto.SignUpResponse, error)
	SignIn(request dto.SignInRequest) (dto.SignInResponse, error)
	ForgotPassword(request dto.ForgotPasswordRequest, accountId uuid.UUID) (dto.ForgotPasswordResponse, error)
	ChangePassword(request dto.ChangePasswordRequest, accountId uuid.UUID) (dto.ChangePasswordResponse, error)
}

func NewAuthenticationService(log *zap.Logger, repository repository.IAuthenticationRepository, email *email.Email) IAuthenticationService {
	return &AuthenticationService{
		log:        log,
		repository: repository,
		email:      email,
	}
}

func (a *AuthenticationService) SignUp(request dto.SignUpRequest, accountId uuid.UUID) (dto.SignUpResponse, error) {
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
		Password:     string(hashedPassword),
		RoleId:       request.RoleId,
		PhotoProfile: "",
		CreatedBy:    accountId,
	}

	staff := entity.Staff{
		AccountId: Id,
		Name:      request.Name,
	}

	if err = a.repository.CreateAccount(&account); err != nil {
		a.log.Error("[SignUp] failed to create account", zap.Error(err))
		return dto.SignUpResponse{}, err
	}

	if err = a.repository.CreateStaff(&staff); err != nil {
		a.log.Error("[SignUp] failed to create staff", zap.Error(err))
		return dto.SignUpResponse{}, err
	}

	account, err = a.repository.GetAccountById(Id)
	if err != nil {
		a.log.Error("[SignUp] failed to get account by id", zap.Error(err))
		return dto.SignUpResponse{}, err
	}

	if err = a.repository.Commit(); err != nil {
		a.log.Error("[SignUp] failed to commit transaction", zap.Error(err))
	}

	return dto.SignUpResponse{
		Id:           account.Id.String(),
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.SignInResponse{}, errx.BadRequest("password or email is incorrect")
		}
		return dto.SignInResponse{}, err
	}

	if bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(request.Password)) != nil {
		a.log.Error("[SignIn] password or email is incorrect")
		return dto.SignInResponse{}, errx.BadRequest("password or email is incorrect")
	}

	token, err := jwt.EncodeToken(&account)
	if err != nil {
		a.log.Error("[SignIn] failed to encode token", zap.Error(err))
		return dto.SignInResponse{}, err
	}

	return dto.SignInResponse{
		TokenType:   constant.TokenType,
		AccessToken: token,
		Role:        account.Role.Name,
		ExpiredAt:   viper.GetDuration("jwt.expiration"),
	}, nil
}

func (a *AuthenticationService) ForgotPassword(request dto.ForgotPasswordRequest, accountId uuid.UUID) (dto.ForgotPasswordResponse, error) {
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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(tempPassord), bcrypt.DefaultCost)
	if err != nil {
		a.log.Error("[ForgotPassword] failed to hash password", zap.Error(err))
		return dto.ForgotPasswordResponse{}, nil
	}

	account.Password = string(hashedPassword)
	account.UpdatedBy = accountId
	if err := a.repository.UpdateAccount(&account); err != nil {
		a.log.Error("[ForgotPassword] failed to update account", zap.Error(err))
		return dto.ForgotPasswordResponse{}, nil
	}

	a.email.SetReciever(account.Email)
	a.email.SetSubject("Forgot Password")
	a.email.SetSender(viper.GetString("email.from"))
	a.email.SetBodyHTML("forgot_password.html", tempPassord)
	if err := a.email.Send(); err != nil {
		a.log.Error("[ForgotPassword] failed to send email", zap.Error(err))
		return dto.ForgotPasswordResponse{}, err
	}

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

	if bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(request.OldPassword)) != nil {
		a.log.Error("[ChangePassword] password is incorrect")
		return dto.ChangePasswordResponse{}, nil
	}

	if request.NewPassword != request.ConfirmPassword {
		a.log.Error("[ChangePassword] new password and confirm password not match")
		return dto.ChangePasswordResponse{}, errx.BadRequest("new password and confirm password not match")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		a.log.Error("[ChangePassword] failed to hash password", zap.Error(err))
		return dto.ChangePasswordResponse{}, nil
	}

	account.Password = string(hashedPassword)
	account.UpdatedBy = accountId

	if err := a.repository.UpdateAccount(&account); err != nil {
		a.log.Error("[ChangePassword] failed to update account", zap.Error(err))
		return dto.ChangePasswordResponse{}, nil
	}

	account, err = a.repository.GetAccountById(accountId)
	if err != nil {
		a.log.Error("[SignUp] failed to get account by id", zap.Error(err))
		return dto.ChangePasswordResponse{}, err
	}

	return dto.ChangePasswordResponse{
		Id:           account.Id.String(),
		PhotoProfile: account.PhotoProfile,
		Email:        account.Email,
		Role: dto.RoleResponse{
			Id:   account.Role.Id,
			Name: account.Role.Name,
		},
	}, nil
}
