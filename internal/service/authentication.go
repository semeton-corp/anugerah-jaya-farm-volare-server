package service

import (
	"slices"

	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/infra/email"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
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
	roleService  IRoleService
}

type IAuthenticationService interface {
	SignUp(request dto.SignUpRequest, userId uuid.UUID) (dto.SignUpResponse, error)
	SignIn(request dto.SignInRequest) (dto.SignInResponse, error)
	ForgotPassword(request dto.ForgotPasswordRequest) (dto.ForgotPasswordResponse, error)
	ChangePassword(request dto.ChangePasswordRequest, userId uuid.UUID) (dto.ChangePasswordResponse, error)
	DeleteUser(id uuid.UUID) error
}

func NewAuthenticationService(log *zap.Logger, repository repository.IAuthenticationRepository, emailService email.IEmail, roleService IRoleService) IAuthenticationService {
	return &AuthenticationService{
		log:          log,
		repository:   repository,
		emailService: emailService,
		roleService:  roleService,
	}
}

func (s *AuthenticationService) SignUp(request dto.SignUpRequest, userId uuid.UUID) (dto.SignUpResponse, error) {
	s.repository.UseTx(true)
	defer s.repository.Rollback()

	Id, err := uuid.NewUUID()
	if err != nil {
		s.log.Error("[SignUp] failed to generate UUID", zap.Error(err))
		return dto.SignUpResponse{}, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		s.log.Error("[SignUp] failed to hash password", zap.Error(err))
		return dto.SignUpResponse{}, err
	}

	salary, err := decimal.NewFromString(request.Salary)
	if err != nil {
		s.log.Error("[SignUp] failed to parse salary from string", zap.Error(err))
		return dto.SignUpResponse{}, err
	}

	role, err := s.roleService.GetRoleById(request.RoleId)
	if err != nil {
		return dto.SignUpResponse{}, err
	}

	user := entity.User{
		Id:           Id,
		Email:        request.Email,
		Username:     request.Username,
		Password:     string(hashedPassword),
		RoleId:       request.RoleId,
		LocationId:   request.LocationId,
		PhotoProfile: "https://www.gravatar.com/avatar/?d=mp",
		Name:         request.Name,
		PhoneNumber:  request.PhoneNumber,
		Address:      request.Address,
		Salary:       salary,
		CreatedBy:    uuid.NullUUID{UUID: userId, Valid: true},
	}

	if request.PhotoProfile != "" {
		user.PhotoProfile = request.PhotoProfile
	}

	if err = s.repository.CreateUser(&user); err != nil {
		s.log.Error("[SignUp] failed to create user", zap.Error(err))
		return dto.SignUpResponse{}, err
	}

	if slices.Contains(entity.CageLocationTypeList, role.Name) {
		cagePlacement := make([]entity.CagePlacement, 0)
		for _, e := range request.PlacementIds {
			cagePlacement = append(cagePlacement, entity.CagePlacement{
				UserId:    Id,
				CageId:    e,
				CreatedBy: uuid.NullUUID{UUID: userId, Valid: true},
			})
		}
	} else if slices.Contains(entity.StoreLocationTypeList, role.Name) {
		warehousePlacement := make([]entity.WarehousePlacement, 0)
		for _, e := range request.PlacementIds {
			warehousePlacement = append(warehousePlacement, entity.WarehousePlacement{
				UserId:     Id,
				WarehousId: e,
				CreatedBy:  uuid.NullUUID{UUID: userId, Valid: true},
			})
		}
	} else if slices.Contains(entity.WarehouseLocationTypeList, role.Name) {
		storePlacement := make([]entity.StorePlacement, 0)
		for _, e := range request.PlacementIds {
			storePlacement = append(storePlacement, entity.StorePlacement{
				UserId:    Id,
				StoreId:   e,
				CreatedBy: uuid.NullUUID{UUID: userId, Valid: true},
			})
		}
	}

	user, err = s.repository.GetUserById(Id)
	if err != nil {
		s.log.Error("[SignUp] failed to get user by id", zap.Error(err))
		return dto.SignUpResponse{}, err
	}

	if err = s.repository.Commit(); err != nil {
		s.log.Error("[SignUp] failed to commit transaction", zap.Error(err))
	}

	return dto.SignUpResponse{
		Id:           user.Id.String(),
		PhotoProfile: user.PhotoProfile,
		Email:        user.Email,
		Role: dto.RoleResponse{
			Id:   user.Role.Id,
			Name: user.Role.Name,
		},
		Location: dto.LocationResponse{
			Id:   user.Location.Id,
			Name: user.Location.Name,
		},
	}, nil
}

func (s *AuthenticationService) SignIn(request dto.SignInRequest) (dto.SignInResponse, error) {
	s.repository.UseTx(false)

	user, err := s.repository.GetUserByUsername(request.Username)
	if err != nil {
		s.log.Error("[SigIn] failed to get user by email", zap.Error(err))
		return dto.SignInResponse{}, err
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)) != nil {
		s.log.Error("[SignIn] password or email is incorrect")
		return dto.SignInResponse{}, errx.BadRequest("password or email is incorrect")
	}

	token, err := jwt.EncodeToken(&user)
	if err != nil {
		s.log.Error("[SignIn] failed to encode token", zap.Error(err))
		return dto.SignInResponse{}, err
	}

	user, err = s.repository.GetUserById(user.Id)
	if err != nil {
		s.log.Error("[SignIn] failed to get staff by id", zap.Error(err))
		return dto.SignInResponse{}, err
	}

	return dto.SignInResponse{
		Id:           user.Id.String(),
		TokenType:    constant.TokenType,
		Name:         user.Name,
		PhotoProfile: user.PhotoProfile,
		AccessToken:  token,
		Role: dto.RoleResponse{
			Id:   user.Role.Id,
			Name: user.Role.Name,
		},
		Location: dto.LocationResponse{
			Id:   user.Location.Id,
			Name: user.Location.Name,
		},
		ExpiredAt: viper.GetDuration("jwt.expiration"),
	}, nil
}

func (s *AuthenticationService) ForgotPassword(request dto.ForgotPasswordRequest) (dto.ForgotPasswordResponse, error) {
	s.repository.UseTx(false)

	user, err := s.repository.GetUserByEmail(request.Email)
	if err != nil {
		s.log.Error("[ForgotPassword] failed to get user by email", zap.Error(err))
		return dto.ForgotPasswordResponse{}, err
	}

	tempPassword, err := util.RandomString(8)
	if err != nil {
		s.log.Error("[ForgotPassword] failed to generate random string", zap.Error(err))
		return dto.ForgotPasswordResponse{}, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(tempPassword), bcrypt.DefaultCost)
	if err != nil {
		s.log.Error("[ForgotPassword] failed to hash password", zap.Error(err))
		return dto.ForgotPasswordResponse{}, err
	}

	user.Password = string(hashedPassword)

	if err := s.repository.UpdateUser(&user); err != nil {
		s.log.Error("[ForgotPassword] failed to update user", zap.Error(err))
		return dto.ForgotPasswordResponse{}, err
	}

	s.emailService.SetReciever(user.Email)
	s.emailService.SetSubject("Forgot Password")
	s.emailService.SetSender(viper.GetString("email.from"))
	s.emailService.SetBodyHTML("forgot_password.html", struct {
		TempPassword string
	}{
		TempPassword: tempPassword,
	})
	if err := s.emailService.Send(); err != nil {
		s.log.Error("[ForgotPassword] failed to send email", zap.Error(err))
		return dto.ForgotPasswordResponse{}, err
	}

	return dto.ForgotPasswordResponse{
		Id:    user.Id.String(),
		Email: user.Email,
	}, nil
}

func (s *AuthenticationService) ChangePassword(request dto.ChangePasswordRequest, userId uuid.UUID) (dto.ChangePasswordResponse, error) {
	s.repository.UseTx(false)

	user, err := s.repository.GetUserById(userId)
	if err != nil {
		s.log.Error("[ChangePassword] failed to get user by id", zap.Error(err))
		return dto.ChangePasswordResponse{}, nil
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.OldPassword)) != nil {
		s.log.Error("[ChangePassword] password is incorrect")
		return dto.ChangePasswordResponse{}, errx.BadRequest("old password is incorrect")
	}

	if request.NewPassword != request.ConfirmPassword {
		s.log.Error("[ChangePassword] new password and confirm password not match")
		return dto.ChangePasswordResponse{}, errx.BadRequest("new password and confirm password not match")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		s.log.Error("[ChangePassword] failed to hash password", zap.Error(err))
		return dto.ChangePasswordResponse{}, nil
	}

	user.Password = string(hashedPassword)
	user.UpdatedBy = uuid.NullUUID{UUID: userId, Valid: true}

	if err := s.repository.UpdateUser(&user); err != nil {
		s.log.Error("[ChangePassword] failed to update user", zap.Error(err))
		return dto.ChangePasswordResponse{}, nil
	}

	user, err = s.repository.GetUserById(userId)
	if err != nil {
		s.log.Error("[SignUp] failed to get user by id", zap.Error(err))
		return dto.ChangePasswordResponse{}, err
	}

	return dto.ChangePasswordResponse{
		Id:           user.Id.String(),
		PhotoProfile: user.PhotoProfile,
		Email:        user.Email,
		Role: dto.RoleResponse{
			Id:   user.Role.Id,
			Name: user.Role.Name,
		},
		Location: dto.LocationResponse{
			Id:   user.Location.Id,
			Name: user.Location.Name,
		},
	}, nil
}

func (s *AuthenticationService) DeleteUser(id uuid.UUID) error {
	s.repository.UseTx(false)

	if err := s.repository.DeleteUser(id); err != nil {
		s.log.Error("[DeleteUser] failed to delete user", zap.Error(err))
		return nil
	}

	return nil
}
