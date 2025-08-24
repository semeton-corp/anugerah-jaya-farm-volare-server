package service

import (
	"database/sql"
	"slices"

	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/infra/email"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/repository"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/constant"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/errx"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/jwt"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/util"
	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthenticationService struct {
	log              *zap.Logger
	emailService     email.IEmail
	repository       repository.IAuthenticationRepository
	roleService      IRoleService
	placementService IPlacementService
}

type IAuthenticationService interface {
	SignUp(request dto.SignUpRequest, userId uuid.UUID) (dto.SignUpResponse, error)
	SignIn(request dto.SignInRequest) (dto.SignInResponse, error)
	ForgotPassword(request dto.ForgotPasswordRequest) (dto.ForgotPasswordResponse, error)
	ChangePassword(request dto.ChangePasswordRequest, userId uuid.UUID) (dto.ChangePasswordResponse, error)
	DeleteUser(id uuid.UUID) error
}

func NewAuthenticationService(log *zap.Logger, repository repository.IAuthenticationRepository, emailService email.IEmail, roleService IRoleService, placementService IPlacementService) IAuthenticationService {
	return &AuthenticationService{
		log:              log,
		repository:       repository,
		emailService:     emailService,
		roleService:      roleService,
		placementService: placementService,
	}
}

func (s *AuthenticationService) SignUp(request dto.SignUpRequest, userId uuid.UUID) (dto.SignUpResponse, error) {
	s.repository.UseTx(true)
	defer s.repository.Rollback()

	Id, err := uuid.NewUUID()
	if err != nil {
		s.log.Error("failed to generate UUID", zap.Error(err))
		return dto.SignUpResponse{}, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		s.log.Error("failed to hash password", zap.Error(err))
		return dto.SignUpResponse{}, err
	}

	salary, err := decimal.NewFromString(request.Salary)
	if err != nil {
		s.log.Error("failed to parse salary from string", zap.Error(err))
		return dto.SignUpResponse{}, err
	}

	salaryInterval := enum.ValueOfSalaryInterval(request.SalaryInterval)
	if !salaryInterval.IsValid() {
		s.log.Warn("invalid salary interval", zap.String("salaryInterval", request.SalaryInterval))
		return dto.SignUpResponse{}, err
	}

	user := entity.User{
		Id:             Id,
		Email:          request.Email,
		Username:       request.Username,
		Password:       string(hashedPassword),
		RoleId:         request.RoleId,
		PhotoProfile:   "https://www.gravatar.com/avatar/?d=mp",
		Name:           request.Name,
		PhoneNumber:    request.PhoneNumber,
		Address:        request.Address,
		Salary:         salary,
		SalaryInterval: salaryInterval,
		CreatedByOwner: uuid.NullUUID{UUID: userId, Valid: true},
	}

	if request.LocationId != nil {
		user.LocationId = sql.NullInt64{Int64: int64(*request.LocationId), Valid: true}
	}

	if request.PhotoProfile != "" {
		user.PhotoProfile = request.PhotoProfile
	}

	if err = s.repository.CreateUser(&user); err != nil {
		s.log.Error("failed to create user", zap.Error(err))
		return dto.SignUpResponse{}, err
	}

	if err = s.repository.Commit(); err != nil {
		s.log.Error("failed to commit transaction", zap.Error(err))
	}

	if request.PlacementIds != nil {
		role, err := s.roleService.GetRoleById(request.RoleId)
		if err != nil {
			return dto.SignUpResponse{}, err
		}

		if slices.Contains(entity.CageLocationTypeList, role.Name) {
			createCagePlacementRequests := make([]dto.CreateCagePlacementRequest, 0)
			for _, id := range request.PlacementIds {
				createCagePlacementRequests = append(createCagePlacementRequests, dto.CreateCagePlacementRequest{
					UserId: Id.String(),
					CageId: id,
				})
			}

			_, err := s.placementService.CreateCagePlacementForAuthentication(createCagePlacementRequests, userId)
			if err != nil {
				// Saga Pattern
				s.repository.DeleteUser(Id)
				return dto.SignUpResponse{}, err
			}
		} else if slices.Contains(entity.StoreLocationTypeList, role.Name) {
			if len(request.PlacementIds) > 1 {
				return dto.SignUpResponse{}, errx.BadRequest("store type must be only 1 placement")
			}

			_, err := s.placementService.CreateStorePlacementForAuthentication(dto.CreateStorePlacementRequest{
				UserId:  Id.String(),
				StoreId: request.PlacementIds[0],
			}, userId)
			if err != nil {
				// Saga Pattern
				s.repository.DeleteUser(Id)
				return dto.SignUpResponse{}, err
			}
		} else if slices.Contains(entity.WarehouseLocationTypeList, role.Name) {
			if len(request.PlacementIds) > 1 {
				return dto.SignUpResponse{}, errx.BadRequest("warehouse type must be only 1 placement")
			}

			_, err := s.placementService.CreateWarehousePlacementForAuthentication(dto.CreateWarehousePlacementRequest{
				UserId:      Id.String(),
				WarehouseId: request.PlacementIds[0],
			}, userId)
			if err != nil {
				// Saga Pattern
				s.repository.DeleteUser(Id)
				return dto.SignUpResponse{}, err
			}
		}
	}

	user, err = s.repository.GetUserById(Id)
	if err != nil {
		s.log.Error("failed to get user by id", zap.Error(err))
		return dto.SignUpResponse{}, err
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
		s.log.Error("failed to get user by username", zap.Error(err))
		return dto.SignInResponse{}, err
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)) != nil {
		s.log.Error("password or email is incorrect")
		return dto.SignInResponse{}, errx.BadRequest("password or email is incorrect")
	}

	token, err := jwt.EncodeToken(&user)
	if err != nil {
		s.log.Error("failed to encode token", zap.Error(err))
		return dto.SignInResponse{}, err
	}

	user, err = s.repository.GetUserById(user.Id)
	if err != nil {
		s.log.Error("failed to get user by id", zap.Error(err))
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
		s.log.Error("failed to get user by email", zap.Error(err))
		return dto.ForgotPasswordResponse{}, err
	}

	tempPassword, err := util.RandomString(8)
	if err != nil {
		s.log.Error("failed to generate random string", zap.Error(err))
		return dto.ForgotPasswordResponse{}, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(tempPassword), bcrypt.DefaultCost)
	if err != nil {
		s.log.Error("failed to hash password", zap.Error(err))
		return dto.ForgotPasswordResponse{}, err
	}

	user.Password = string(hashedPassword)

	if err := s.repository.UpdateUser(&user); err != nil {
		s.log.Error("failed to update user", zap.Error(err))
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
		s.log.Error("failed to send email", zap.Error(err))
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
		s.log.Error("failed to get user by id", zap.Error(err))
		return dto.ChangePasswordResponse{}, nil
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.OldPassword)) != nil {
		s.log.Error("password is incorrect")
		return dto.ChangePasswordResponse{}, errx.BadRequest("old password is incorrect")
	}

	if request.NewPassword != request.ConfirmPassword {
		s.log.Error("new password and confirm password not match")
		return dto.ChangePasswordResponse{}, errx.BadRequest("new password and confirm password not match")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		s.log.Error("failed to hash password", zap.Error(err))
		return dto.ChangePasswordResponse{}, nil
	}

	user.Password = string(hashedPassword)
	user.UpdatedBy = uuid.NullUUID{UUID: userId, Valid: true}

	if err := s.repository.UpdateUser(&user); err != nil {
		s.log.Error("failed to update user", zap.Error(err))
		return dto.ChangePasswordResponse{}, nil
	}

	user, err = s.repository.GetUserById(userId)
	if err != nil {
		s.log.Error("failed to get user by id", zap.Error(err))
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
		s.log.Error("failed to delete user", zap.Error(err))
		return nil
	}

	return nil
}
