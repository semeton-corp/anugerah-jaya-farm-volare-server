package service

import (
	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/infra/email"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/mapper"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/repository"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/constant"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/errx"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/jwt"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/util"
	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthenticationService struct {
	log          *zap.Logger
	emailService email.IEmail
	repository   repository.IAuthenticationRepository
}

type IAuthenticationService interface {
	SignUp(request dto.SignUpRequest, accoundId uuid.UUID) (dto.SignUpResponse, error)
	SignIn(request dto.SignInRequest) (dto.SignInResponse, error)
	ForgotPassword(request dto.ForgotPasswordRequest) (dto.ForgotPasswordResponse, error)
	ChangePassword(request dto.ChangePasswordRequest, accountId uuid.UUID) (dto.ChangePasswordResponse, error)
	UpdateAccount(id uuid.UUID, request dto.UpdateAccountRequest, accountId uuid.UUID) (dto.ChangePasswordResponse, error)
	DeleteAccount(id uuid.UUID) error
	GetAccountById(id uuid.UUID) (dto.AccountResponse, error)
}

func NewAuthenticationService(log *zap.Logger, repository repository.IAuthenticationRepository, emailService email.IEmail) IAuthenticationService {
	return &AuthenticationService{
		log:          log,
		repository:   repository,
		emailService: emailService,
	}
}

func (a *AuthenticationService) SignUp(request dto.SignUpRequest, accountId uuid.UUID) (dto.SignUpResponse, error) {
	a.repository.UseTx(true)
	defer a.repository.Rollback()

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
		PhotoProfile: "https://www.gravatar.com/avatar/?d=mp",
		CreatedBy:    uuid.NullUUID{UUID: accountId, Valid: true},
	}

	if request.PhotoProfile != "" {
		account.PhotoProfile = request.PhotoProfile
	}

	salary, err := decimal.NewFromString(request.Salary)
	if err != nil {
		a.log.Error("[SignUp] failed to parse salary from string")
		return dto.SignUpResponse{}, err
	}

	staff := entity.Staff{
		Id:          Id,
		AccountId:   Id,
		Name:        request.Name,
		PhoneNumber: request.PhoneNumber,
		Address:     request.Address,
		Salary:      salary,
		CreatedBy:   uuid.NullUUID{UUID: accountId, Valid: true},
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

	staff, err := a.repository.GetStaffById(account.Id)
	if err != nil {
		a.log.Error("[SignIn] failed to get staff by id", zap.Error(err))
		return dto.SignInResponse{}, err
	}

	return dto.SignInResponse{
		TokenType:    constant.TokenType,
		Name:         staff.Name,
		PhotoProfile: account.PhotoProfile,
		AccessToken:  token,
		Role: dto.RoleResponse{
			Id:   account.Role.Id,
			Name: account.Role.Name,
		},
		ExpiredAt: viper.GetDuration("jwt.expiration"),
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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(tempPassord), bcrypt.DefaultCost)
	if err != nil {
		a.log.Error("[ForgotPassword] failed to hash password", zap.Error(err))
		return dto.ForgotPasswordResponse{}, nil
	}

	account.Password = string(hashedPassword)

	if err := a.repository.UpdateAccount(&account); err != nil {
		a.log.Error("[ForgotPassword] failed to update account", zap.Error(err))
		return dto.ForgotPasswordResponse{}, nil
	}

	a.emailService.SetReciever(account.Email)
	a.emailService.SetSubject("Forgot Password")
	a.emailService.SetSender(viper.GetString("email.from"))
	a.emailService.SetBodyHTML("forgot_password.html", tempPassord)
	if err := a.emailService.Send(); err != nil {
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
		return dto.ChangePasswordResponse{}, errx.BadRequest("old password is incorrect")
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
	account.UpdatedBy = uuid.NullUUID{UUID: accountId, Valid: true}

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

func (a *AuthenticationService) UpdateAccount(id uuid.UUID, request dto.UpdateAccountRequest, accountId uuid.UUID) (dto.ChangePasswordResponse, error) {
	a.repository.UseTx(false)

	account, err := a.repository.GetAccountById(id)
	if err != nil {
		a.log.Error("[UpdateAccount] failed to get account by id", zap.Error(err))
		return dto.ChangePasswordResponse{}, nil
	}

	account.Email = request.Email
	account.RoleId = request.RoleId
	account.PhotoProfile = request.PhotoProfile
	account.UpdatedBy = uuid.NullUUID{UUID: accountId, Valid: true}

	if err := a.repository.UpdateAccount(&account); err != nil {
		a.log.Error("[UpdateAccount] failed to update account", zap.Error(err))
		return dto.ChangePasswordResponse{}, nil
	}

	account, err = a.repository.GetAccountById(accountId)
	if err != nil {
		a.log.Error("[SignUp] failed to get account by id", zap.Error(err))
		return dto.ChangePasswordResponse{}, nil
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

func (a *AuthenticationService) DeleteAccount(id uuid.UUID) error {
	a.repository.UseTx(false)

	if err := a.repository.DeleteAccount(id); err != nil {
		a.log.Error("[DeleteAccount] failed to delete account", zap.Error(err))
		return nil
	}

	return nil
}

func (a *AuthenticationService) GetAccountById(id uuid.UUID) (dto.AccountResponse, error) {
	a.repository.UseTx(false)

	account, err := a.repository.GetAccountById(id)
	if err != nil {
		a.log.Error("[GetAccountById] failed to get account by id")
		return dto.AccountResponse{}, err
	}

	return mapper.AccountToResponse(&account), nil
}
