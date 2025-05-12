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

type StaffService struct {
	log             *zap.Logger
	repository      repository.IStaffRepository
	authService     IAuthenticationService
	workService     IWorkService
	presenceService IPresenceService
}

type IStaffService interface {
	GetStaffById(id uuid.UUID) (dto.StaffResponse, error)
	UpdateStaff(id uuid.UUID, request dto.UpdateStaffRequest, accountId uuid.UUID) (dto.StaffResponse, error)
	GetStaffs(filter dto.GetStaffFilter) (dto.StaffListPaginationResponse, error)
	GetOverviewStaff(id uuid.UUID, filter dto.GetStaffOverviewFilter) (dto.StaffOverviewResponse, error)
}

func NewStaffService(log *zap.Logger, repository repository.IStaffRepository, authService IAuthenticationService, workService IWorkService, presenceService IPresenceService) IStaffService {
	return &StaffService{
		log:             log,
		repository:      repository,
		authService:     authService,
		workService:     workService,
		presenceService: presenceService,
	}
}

func (s *StaffService) GetStaffById(id uuid.UUID) (dto.StaffResponse, error) {
	s.repository.UseTx(false)

	staff, err := s.repository.GetStaffById(id)
	if err != nil {
		s.log.Error("[GetStaffById] failed to get staff by id", zap.Error(err))
		return dto.StaffResponse{}, err
	}

	return mapper.StaffToResponse(&staff), nil
}

func (s *StaffService) UpdateStaff(id uuid.UUID, request dto.UpdateStaffRequest, accountId uuid.UUID) (dto.StaffResponse, error) {
	var (
		err error
	)

	account, err := s.authService.GetAccountById(id)
	if err != nil {
		s.log.Error("[UpdateStaff] failed to get account by id")
		return dto.StaffResponse{}, err
	}

	s.repository.UseTx(true)

	defer func() {
		if err != nil {
			s.repository.Rollback()

			_, errS := s.authService.UpdateAccount(id, dto.UpdateAccountRequest{
				Email:        account.Email,
				RoleId:       account.Role.Id,
				PhotoProfile: account.PhotoProfile,
			}, accountId)

			if err != nil {
				s.log.Error("[UpdateStaff] failed to update account")
				err = errS
			}
		}
	}()

	staff, err := s.repository.GetStaffById(id)
	if err != nil {
		s.log.Error("[UpdateStaff] failed to get staff by id", zap.Error(err))
		return dto.StaffResponse{}, err
	}

	_, err = s.authService.UpdateAccount(id, dto.UpdateAccountRequest{
		Email:        request.Email,
		RoleId:       request.RoleId,
		PhotoProfile: request.PhotoProfile,
	}, accountId)
	if err != nil {
		s.log.Error("[UpdateStaff] failed to update account")
		return dto.StaffResponse{}, err
	}

	salary, err := decimal.NewFromString(request.Salary)
	if err != nil {
		s.log.Error("[UpdateStaff] failed to parse salary", zap.Error(err))
		return dto.StaffResponse{}, err
	}

	staff.Name = request.Name
	staff.PhoneNumber = request.PhoneNumber
	staff.Address = request.Address
	staff.Salary = salary

	if err := s.repository.UpdateStaff(&staff); err != nil {
		s.log.Error("[UpdateStaff] failed to update staff", zap.Error(err))
		return dto.StaffResponse{}, err
	}

	if err := s.repository.Commit(); err != nil {
		s.log.Error("[UpdateStaff] failed to commit transaction", zap.Error(err))
		return dto.StaffResponse{}, err
	}

	staff, err = s.repository.GetStaffById(id)
	if err != nil {
		s.log.Error("[UpdateStaff] failed to get staff by id", zap.Error(err))
		return dto.StaffResponse{}, err
	}

	return mapper.StaffToResponse(&staff), nil
}

func (s *StaffService) GetStaffs(filter dto.GetStaffFilter) (dto.StaffListPaginationResponse, error) {
	s.repository.UseTx(false)

	staffs, err := s.repository.GetStaffs(&filter)
	if err != nil {
		s.log.Error("[GetStaffs] failed to get staffs", zap.Error(err))
		return dto.StaffListPaginationResponse{}, err
	}

	// Idea : get with diff method repository, and the get totalOvertime, salary for additional work, and cashbon (using join)
	// Todo : get salary from bonus lembur
	// Todo : get salary from additional work
	// Todo : get salary cashbon

	staffResponses := make([]dto.StaffListResponse, 0)
	for _, staff := range staffs {
		staffResponses = append(staffResponses, mapper.StaffToListResponse(&staff))
	}

	totalData, err := s.repository.CountTotalStaff(&dto.GetStaffFilter{
		Keyword: filter.Keyword,
		RoleId:  filter.RoleId,
	})
	if err != nil {
		s.log.Error("[GetStaffs] failed to count total staffs")
		return dto.StaffListPaginationResponse{}, err
	}

	resp := dto.StaffListPaginationResponse{
		TotalPage: uint64(math.Ceil(float64(totalData) / float64(constant.PaginationDefaultLimit))),
		TotalData: totalData,
		Staffs:    staffResponses,
	}

	return resp, nil
}

func (s *StaffService) GetOverviewStaff(id uuid.UUID, filter dto.GetStaffOverviewFilter) (dto.StaffOverviewResponse, error) {
	s.repository.UseTx(false)

	weeks := util.GetFourWeekRanges(int(filter.Year), time.Month(filter.Month))

	staff, err := s.repository.GetStaffById(id)
	if err != nil {
		return dto.StaffOverviewResponse{}, nil
	}

	additionalWorkStaffs, err := s.workService.GetAdditionalWorkStaffByStaffId(id,
		dto.GetAdditionalWorkStaffFilter{
			Month:       filter.Month,
			Year:        filter.Year,
			WithDeleted: true,
		})
	if err != nil {
		return dto.StaffOverviewResponse{}, nil
	}

	dailyWorkStaffs, err := s.workService.GetDailyWorkStaffByStaffId(id,
		dto.GetDailyWorkStaffFilter{
			Month:       filter.Month,
			Year:        filter.Year,
			WithDeleted: true,
		})
	if err != nil {
		return dto.StaffOverviewResponse{}, nil
	}

	staffPresences, err := s.presenceService.GetAllStaffPresences(id,
		dto.GetPresenceFilter{
			Month: filter.Month,
			Year:  filter.Year,
		})
	if err != nil {
		return dto.StaffOverviewResponse{}, nil
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
				return dto.StaffOverviewResponse{}, nil
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
	staffInformation := dto.StaffInformationResponse{
		TotalWorkHour: totalWorkHour,
		TotalSalary:   totalSalary.String(),
		KPIScore:      (presenceScore + workPresence) / 2, // Todo : how about the overtime
	}

	staffPresenceInformation := dto.StaffPresenceInformationResponse{
		TotalPresent:    totalPresent,
		TotalNotPresent: uint64(len(staffPresences.Presences) - int(totalPresent)),
	}

	staffWorkInformation := dto.StaffWorkInformationResponse{
		TotalWorkDone:    totalWorkDone,
		TotalWorkNotDone: uint64(len(dailyWorkStaffs) + len(additionalWorkStaffs) - int(totalWorkDone)),
	}

	staffSalaryInformation := dto.StaffSalaryInformationResponse{
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

	overviewResponse := dto.StaffOverviewResponse{
		StaffInformation:         staffInformation,
		KPIPerformances:          kpiPerformances,
		StaffPresenceInformation: staffPresenceInformation,
		StaffSalaryInformation:   staffSalaryInformation,
		StaffWorkInformation:     staffWorkInformation,
	}

	return overviewResponse, nil
}
