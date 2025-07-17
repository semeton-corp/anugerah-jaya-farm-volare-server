package service

import (
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/mapper"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/repository"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/constant"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/util"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type UserService struct {
	log             *zap.Logger
	repository      repository.IUserRepository
	workService     IWorkService
	presenceService IPresenceService
}

type IUserService interface {
	GetUserById(id uuid.UUID) (dto.UserResponse, error)
	UpdateUser(id uuid.UUID, request dto.UpdateUserRequest, userId uuid.UUID) (dto.UserResponse, error)
	GetUsers(filter dto.GetUserListFilter) ([]dto.UserListResponse, error)

	GetUserOverviewList(filter dto.GetUserOverviewListFilter) (dto.UserListOverviewPaginationResponse, error)
	GetOverviewDetailUser(id uuid.UUID, filter dto.GetUserOverviewFilter) (dto.UserOverviewResponse, error)
}

func NewUserService(log *zap.Logger, repository repository.IUserRepository, workService IWorkService, presenceService IPresenceService) IUserService {
	return &UserService{
		log:             log,
		repository:      repository,
		workService:     workService,
		presenceService: presenceService,
	}
}

func (s *UserService) GetUserById(id uuid.UUID) (dto.UserResponse, error) {
	s.repository.UseTx(false)

	user, err := s.repository.GetUserById(id)
	if err != nil {
		s.log.Error("failed to get user by id", zap.Error(err))
		return dto.UserResponse{}, err
	}

	return mapper.UserToResponse(&user), nil
}

func (s *UserService) UpdateUser(id uuid.UUID, request dto.UpdateUserRequest, userId uuid.UUID) (dto.UserResponse, error) {
	s.repository.UseTx(true)
	defer s.repository.Rollback()

	user, err := s.repository.GetUserById(id)
	if err != nil {
		s.log.Error("failed to get user by id", zap.Error(err))
		return dto.UserResponse{}, err
	}

	salary, err := decimal.NewFromString(request.Salary)
	if err != nil {
		s.log.Error("failed to parse salary", zap.Error(err))
		return dto.UserResponse{}, err
	}

	user.Email = request.Email
	user.Username = request.Username
	user.RoleId = request.RoleId
	user.PhotoProfile = request.PhotoProfile
	user.Name = request.Name
	user.PhoneNumber = request.PhoneNumber
	user.Address = request.Address
	user.Salary = salary

	if err := s.repository.UpdateUser(&user); err != nil {
		s.log.Error("failed to update user", zap.Error(err))
		return dto.UserResponse{}, err
	}

	if err := s.repository.Commit(); err != nil {
		s.log.Error("failed to commit transaction", zap.Error(err))
		return dto.UserResponse{}, err
	}

	user, err = s.repository.GetUserById(id)
	if err != nil {
		s.log.Error("failed to get user by id", zap.Error(err))
		return dto.UserResponse{}, err
	}

	return mapper.UserToResponse(&user), nil
}

func (s *UserService) GetUsers(filter dto.GetUserListFilter) ([]dto.UserListResponse, error) {
	s.repository.UseTx(false)

	users, err := s.repository.GetUsers(&filter)
	if err != nil {
		s.log.Error("failed to get users", zap.Error(err))
		return nil, err
	}

	userResponses := make([]dto.UserListResponse, 0)
	for _, user := range users {
		userResponses = append(userResponses, mapper.UserToListResponse(&user))
	}

	return userResponses, nil
}

func (s *UserService) GetOverviewDetailUser(id uuid.UUID, filter dto.GetUserOverviewFilter) (dto.UserOverviewResponse, error) {
	s.repository.UseTx(false)

	weeks := util.GetFourWeekRanges(int(filter.Year), time.Month(filter.Month))

	user, err := s.repository.GetUserById(id)
	if err != nil {
		return dto.UserOverviewResponse{}, nil
	}

	withDeleted := true
	additionalWorkUsers, err := s.workService.GetAdditionalWorkUserByUserId(id,
		dto.GetAdditionalWorkUserFilter{
			Month:       filter.Month,
			Year:        filter.Year,
			WithDeleted: &withDeleted, // In case the user work is done but the work is deleted
		})
	if err != nil {
		return dto.UserOverviewResponse{}, nil
	}

	dailyWorkUsers, err := s.workService.GetDailyWorkUserByUserId(id,
		dto.GetDailyWorkUserFilter{
			Month:       filter.Month,
			Year:        filter.Year,
			WithDeleted: &withDeleted,
		})
	if err != nil {
		return dto.UserOverviewResponse{}, nil
	}

	userPresences, err := s.presenceService.GetUserPresencesByUserId(id,
		dto.GetPresenceFilter{
			Month: filter.Month,
			Year:  filter.Year,
		})
	if err != nil {
		return dto.UserOverviewResponse{}, nil
	}

	// TODO : get cashbon from cashflow tabel later

	totalPresentWeek := make(map[int]uint64)
	totalOvertimeHourWeek := make(map[int]uint64)
	totalWorkHourWeek := make(map[int]uint64)

	var totalPresent uint64 = 0
	var totalOvertime float64 = 0
	var totalWorkHour float64 = 0
	for _, userPresence := range userPresences.Presences {
		week := util.FindWeek(user.CreatedAt, weeks)
		if week == 0 {
			continue
		}

		if userPresence.Status == enum.PresenceStatusPresent.String() {
			totalPresent++
			totalPresentWeek[week]++

			if userPresence.EndTime != "" {
				startTime, err := time.Parse("15:04", userPresence.StartTime)
				if err != nil {
					continue
				}

				endTime, err := time.Parse("15:04", userPresence.EndTime)
				if err != nil {
					continue
				}

				diffHours := endTime.Sub(startTime).Hours()
				if diffHours > 8 {
					totalWorkHour += 8.0
				} else {
					totalWorkHour += diffHours
				}

				totalOvertime += userPresence.Overtime

				totalWorkHourWeek[week] += uint64(diffHours)
				totalOvertimeHourWeek[week] += uint64(userPresence.Overtime)
			}
		}
	}

	overtimeSalary := decimal.NewFromFloat(100000).Mul(decimal.NewFromFloat(totalOvertime))

	bonusWeek := make(map[int]decimal.Decimal)
	totalWorkDoneWeek := make(map[int]uint64)

	var bonusSalary decimal.Decimal
	var totalWorkDone uint64 = 0
	for _, dailyWorkUser := range dailyWorkUsers.DailyWorkUsers {
		week := util.FindWeek(dailyWorkUser.CreatedAt, weeks)
		if week == 0 {
			continue
		}
		if dailyWorkUser.IsDone {
			totalWorkDone++
			totalWorkDoneWeek[week]++
		}
	}

	for _, additionalWorkUser := range additionalWorkUsers.AdditionalWorkUsers {
		week := util.FindWeek(additionalWorkUser.CreatedAt, weeks)
		if week == 0 {
			continue
		}

		if additionalWorkUser.IsDone {
			totalWorkDone++
			totalWorkDoneWeek[week]++

			salary, err := decimal.NewFromString(additionalWorkUser.AdditionalWork.Salary)
			if err != nil {
				return dto.UserOverviewResponse{}, nil
			}

			bonusSalary = bonusSalary.Add(salary)
			bonusWeek[week] = bonusWeek[week].Add(salary)
		}
	}

	var presenceScore float64
	if len(userPresences.Presences) == 0 {
		presenceScore = 0
	} else {
		presenceScore = float64(totalPresent) / float64(len(userPresences.Presences)) * 100
	}

	workPresence := totalWorkHour / float64(8*util.TotalDaysInMonth(int(filter.Year), time.Month(filter.Month)))

	totalSalary := user.Salary.Add(overtimeSalary).Add(bonusSalary) // Todo : need cashbon
	userInformation := dto.UserInformationResponse{
		TotalWorkHour: totalWorkHour,
		TotalSalary:   totalSalary.String(),
		KPIScore:      (presenceScore + workPresence) / 2, // Todo : how about the overtime
	}

	userPresenceInformation := dto.UserPresenceInformationResponse{
		TotalPresent:    totalPresent,
		TotalNotPresent: uint64(len(userPresences.Presences) - int(totalPresent)),
	}

	userWorkInformation := dto.UserWorkInformationResponse{
		TotalWorkDone:    totalWorkDone,
		TotalWorkNotDone: uint64(len(dailyWorkUsers.DailyWorkUsers) + len(additionalWorkUsers.AdditionalWorkUsers) - int(totalWorkDone)),
	}

	userSalaryInformation := dto.UserSalaryInformationResponse{
		BaseSalary:     user.Salary.String(),
		OvertimeSalary: overtimeSalary.String(),
		BonusSalary:    bonusSalary.String(),
		Cashbon:        decimal.Zero.String(),
		TotalSalary:    totalSalary.String(),
	}

	kpiPerformances := make([]dto.KPIPerformanceResponse, 0)
	for key, value := range weeks {

		presenceScore := float64(totalPresentWeek[key]) / float64(value.TotalDays) * 100
		workPresence := float64(totalWorkHourWeek[key]) / float64(8*value.TotalDays) * 100

		kpiPerformance := dto.KPIPerformanceResponse{
			Key:   fmt.Sprintf("Minggu %d", key),
			Value: (presenceScore + workPresence) / 2,
		}

		kpiPerformances = append(kpiPerformances, kpiPerformance)
	}

	overviewResponse := dto.UserOverviewResponse{
		UserInformation:         userInformation,
		KPIPerformances:         kpiPerformances,
		UserPresenceInformation: userPresenceInformation,
		UserSalaryInformation:   userSalaryInformation,
		UserWorkInformation:     userWorkInformation,
	}

	return overviewResponse, nil
}

func (s *UserService) GetUserOverviewList(filter dto.GetUserOverviewListFilter) (dto.UserListOverviewPaginationResponse, error) {
	s.repository.UseTx(false)

	users, err := s.repository.GetUserOverviews(&filter)
	if err != nil {
		s.log.Error("failed to get user overview", zap.Error(err))
		return dto.UserListOverviewPaginationResponse{}, err
	}

	response := make([]dto.UserListOverviewResponse, 0)
	for _, user := range users {
		response = append(response, mapper.UserOverviewToListResponse(&user))
	}

	totalData, err := s.repository.CountTotalUserOverview(&filter)
	if err != nil {
		s.log.Error("failed count user overview", zap.Error(err))
		return dto.UserListOverviewPaginationResponse{}, err
	}

	resp := dto.UserListOverviewPaginationResponse{
		Users: response,
	}

	if filter.Page > 0 {
		resp.TotalData = totalData
		resp.TotalPage = uint64(math.Ceil(float64(totalData) / float64(constant.PaginationDefaultLimit)))
	}

	return resp, nil
}
