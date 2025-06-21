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
	UpdateUser(id uuid.UUID, request dto.UpdateUserRequest, accountId uuid.UUID) (dto.UserResponse, error)
	GetUsers(filter dto.GetUserFilter) (dto.UserListPaginationResponse, error)
	GetOverviewUser(id uuid.UUID, filter dto.GetUserOverviewFilter) (dto.UserOverviewResponse, error)
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

	staff, err := s.repository.GetUserById(id)
	if err != nil {
		s.log.Error("[GetStaffById] failed to get staff by id", zap.Error(err))
		return dto.UserResponse{}, err
	}

	return mapper.UserToResponse(&staff), nil
}

func (s *UserService) UpdateUser(id uuid.UUID, request dto.UpdateUserRequest, accountId uuid.UUID) (dto.UserResponse, error) {
	s.repository.UseTx(true)
	defer s.repository.Rollback()

	user, err := s.repository.GetUserById(id)
	if err != nil {
		s.log.Error("[UpdateStaff] failed to get staff by id", zap.Error(err))
		return dto.UserResponse{}, err
	}

	salary, err := decimal.NewFromString(request.Salary)
	if err != nil {
		s.log.Error("[UpdateStaff] failed to parse salary", zap.Error(err))
		return dto.UserResponse{}, err
	}

	user.Email = request.Email
	user.RoleId = request.RoleId
	user.PhotoProfile = request.PhotoProfile
	user.Name = request.Name
	user.PhoneNumber = request.PhoneNumber
	user.Address = request.Address
	user.Salary = salary

	if err := s.repository.UpdateUser(&user); err != nil {
		s.log.Error("[UpdateStaff] failed to update staff", zap.Error(err))
		return dto.UserResponse{}, err
	}

	if err := s.repository.Commit(); err != nil {
		s.log.Error("[UpdateStaff] failed to commit transaction", zap.Error(err))
		return dto.UserResponse{}, err
	}

	user, err = s.repository.GetUserById(id)
	if err != nil {
		s.log.Error("[UpdateStaff] failed to get staff by id", zap.Error(err))
		return dto.UserResponse{}, err
	}

	return mapper.UserToResponse(&user), nil
}

func (s *UserService) GetUsers(filter dto.GetUserFilter) (dto.UserListPaginationResponse, error) {
	s.repository.UseTx(false)

	staffs, err := s.repository.GetUsers(&filter)
	if err != nil {
		s.log.Error("[GetStaffs] failed to get staffs", zap.Error(err))
		return dto.UserListPaginationResponse{}, err
	}

	// Idea : get with diff method repository, and the get totalOvertime, salary for additional work, and cashbon (using join)
	// Todo : get salary from bonus lembur
	// Todo : get salary from additional work
	// Todo : get salary cashbon

	userResponses := make([]dto.UserListResponse, 0)
	for _, staff := range staffs {
		userResponses = append(userResponses, mapper.UserToListResponse(&staff))
	}

	totalData, err := s.repository.CountTotalUser(&dto.GetUserFilter{
		Keyword: filter.Keyword,
		RoleId:  filter.RoleId,
	})
	if err != nil {
		s.log.Error("[GetStaffs] failed to count total staffs")
		return dto.UserListPaginationResponse{}, err
	}

	resp := dto.UserListPaginationResponse{
		TotalPage: uint64(math.Ceil(float64(totalData) / float64(constant.PaginationDefaultLimit))),
		TotalData: totalData,
		Users:     userResponses,
	}

	return resp, nil
}

func (s *UserService) GetOverviewUser(id uuid.UUID, filter dto.GetUserOverviewFilter) (dto.UserOverviewResponse, error) {
	s.repository.UseTx(false)

	weeks := util.GetFourWeekRanges(int(filter.Year), time.Month(filter.Month))

	staff, err := s.repository.GetUserById(id)
	if err != nil {
		return dto.UserOverviewResponse{}, nil
	}

	additionalWorkStaffs, err := s.workService.GetAdditionalWorkStaffByStaffId(id,
		dto.GetAdditionalWorkStaffFilter{
			Month:       filter.Month,
			Year:        filter.Year,
			WithDeleted: true,
		})
	if err != nil {
		return dto.UserOverviewResponse{}, nil
	}

	dailyWorkStaffs, err := s.workService.GetDailyWorkStaffByStaffId(id,
		dto.GetDailyWorkStaffFilter{
			Month:       filter.Month,
			Year:        filter.Year,
			WithDeleted: true,
		})
	if err != nil {
		return dto.UserOverviewResponse{}, nil
	}

	staffPresences, err := s.presenceService.GetAllStaffPresences(id,
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
	for _, staffPresence := range staffPresences.Presences {
		week := util.FindWeek(staff.CreatedAt, weeks)
		if week == 0 {
			continue
		}

		if staffPresence.IsPresent {
			totalPresent++
			totalPresentWeek[week]++

			if staffPresence.EndTime != "" {
				startTime, err := time.Parse("15:04", staffPresence.StartTime)
				if err != nil {
					continue
				}

				endTime, err := time.Parse("15:04", staffPresence.EndTime)
				if err != nil {
					continue
				}

				diffHours := endTime.Sub(startTime).Hours()
				if diffHours > 8 {
					totalWorkHour += 8.0
				} else {
					totalWorkHour += diffHours
				}

				totalOvertime += staffPresence.Overtime

				totalWorkHourWeek[week] += uint64(diffHours)
				totalOvertimeHourWeek[week] += uint64(staffPresence.Overtime)
			}
		}
	}

	overtimeSalary := decimal.NewFromFloat(100000).Mul(decimal.NewFromFloat(totalOvertime))

	bonusWeek := make(map[int]decimal.Decimal)
	totalWorkDoneWeek := make(map[int]uint64)

	var bonusSalary decimal.Decimal
	var totalWorkDone uint64 = 0
	for _, dailyWorkStaff := range dailyWorkStaffs {
		week := util.FindWeek(dailyWorkStaff.CreatedAt, weeks)
		if week == 0 {
			continue
		}
		if dailyWorkStaff.IsDone {
			totalWorkDone++
			totalWorkDoneWeek[week]++
		}
	}

	for _, additionalWorkStaff := range additionalWorkStaffs {
		week := util.FindWeek(additionalWorkStaff.CreatedAt, weeks)
		if week == 0 {
			continue
		}

		if additionalWorkStaff.IsDone {
			totalWorkDone++
			totalWorkDoneWeek[week]++

			salary, err := decimal.NewFromString(additionalWorkStaff.AdditionalWork.Salary)
			if err != nil {
				return dto.UserOverviewResponse{}, nil
			}

			bonusSalary = bonusSalary.Add(salary)
			bonusWeek[week] = bonusWeek[week].Add(salary)
		}
	}

	var presenceScore float64
	if len(staffPresences.Presences) == 0 {
		presenceScore = 0
	} else {
		presenceScore = float64(totalPresent) / float64(len(staffPresences.Presences)) * 100
	}

	workPresence := totalWorkHour / float64(8*util.TotalDaysInMonth(int(filter.Year), time.Month(filter.Month)))

	totalSalary := staff.Salary.Add(overtimeSalary).Add(bonusSalary) // Todo : need cashbon
	userInformation := dto.UserInformationResponse{
		TotalWorkHour: totalWorkHour,
		TotalSalary:   totalSalary.String(),
		KPIScore:      (presenceScore + workPresence) / 2, // Todo : how about the overtime
	}

	userPresenceInformation := dto.UserPresenceInformationResponse{
		TotalPresent:    totalPresent,
		TotalNotPresent: uint64(len(staffPresences.Presences) - int(totalPresent)),
	}

	userWorkInformation := dto.UserWorkInformationResponse{
		TotalWorkDone:    totalWorkDone,
		TotalWorkNotDone: uint64(len(dailyWorkStaffs) + len(additionalWorkStaffs) - int(totalWorkDone)),
	}

	userSalaryInformation := dto.UserSalaryInformationResponse{
		BaseSalary:     staff.Salary.String(),
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
