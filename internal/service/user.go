package service

import (
	"fmt"
	"math"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/mapper"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/repository"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/constant"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/param"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/util"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type UserService struct {
	log              *zap.Logger
	repository       repository.IUserRepository
	workService      IWorkService
	roleService      IRoleService
	presenceService  IPresenceService
	chickenService   IChickenService
	placementService IPlacementService
	cashflowService  ICashflowService
}

type IUserService interface {
	GetUserById(id uuid.UUID) (dto.UserResponse, error)
	UpdateUser(id uuid.UUID, request dto.UpdateUserRequest, userId uuid.UUID) (dto.UserResponse, error)
	GetUsers(filter dto.GetUserListFilter) ([]dto.UserListResponse, error)

	GetUserOverviewList(filter dto.GetUserOverviewListFilter) (dto.UserListOverviewPaginationResponse, error)
	GetUserOverview(id uuid.UUID, filter dto.GetUserOverviewFilter) (dto.UserOverviewResponse, error)

	GetUserPerformanceOverview(filter dto.GetUserPerformanceOverviewFilter) (dto.UserPerformanceOverviewResponse, error)
}

func NewUserService(log *zap.Logger, repository repository.IUserRepository, workService IWorkService, presenceService IPresenceService, chickenService IChickenService, placementService IPlacementService, cashflowService ICashflowService, roleService IRoleService) IUserService {
	return &UserService{
		log:              log,
		repository:       repository,
		workService:      workService,
		presenceService:  presenceService,
		chickenService:   chickenService,
		placementService: placementService,
		cashflowService:  cashflowService,
		roleService:      roleService,
	}
}

func (s *UserService) GetUserById(id uuid.UUID) (dto.UserResponse, error) {
	s.repository.UseTx(false)

	user, err := s.repository.GetUserById(id)
	if err != nil {
		s.log.Error("failed to get user by id", zap.Error(err))
		return dto.UserResponse{}, err
	}

	placements := make([]dto.PlacementResponse, 0)
	if slices.Contains(entity.CageLocationTypeList, user.Role.Name) {
		data, err := s.placementService.GetCagePlacementByUserId(user.Id)
		if err != nil {
			return dto.UserResponse{}, err
		}

		for _, e := range data {
			placements = append(placements, dto.PlacementResponse{
				PlaceId:   e.Cage.Id,
				PlaceName: e.Cage.Name,
			})
		}
	} else if slices.Contains(entity.WarehouseLocationTypeList, user.Role.Name) {
		data, err := s.placementService.GetWarehousePlacementByUserId(user.Id)
		if err != nil {
			return dto.UserResponse{}, err
		}

		for _, e := range data {
			placements = append(placements, dto.PlacementResponse{
				PlaceId:   e.Warehouse.Id,
				PlaceName: e.Warehouse.Name,
			})
		}
	} else if slices.Contains(entity.StoreLocationTypeList, user.Role.Name) {
		data, err := s.placementService.GetStorePlacementByUserId(user.Id)
		if err != nil {
			return dto.UserResponse{}, err
		}

		for _, e := range data {
			placements = append(placements, dto.PlacementResponse{
				PlaceId:   e.Store.Id,
				PlaceName: e.Store.Name,
			})
		}
	}

	response := mapper.UserToResponse(&user)
	response.Placement = placements

	return response, nil
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

func (s *UserService) GetUserOverview(id uuid.UUID, filter dto.GetUserOverviewFilter) (dto.UserOverviewResponse, error) {
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
			WithDeleted: &withDeleted,
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

	totalPresentWeek := make(map[int]uint64)
	totalWorkHourWeek := make(map[int]uint64)

	var totalPresent uint64 = 0
	var totalWorkHour float64 = 0
	for _, userPresence := range userPresences.Presences {
		week := util.FindWeek(userPresence.CreatedAt, weeks)
		if week == 0 {
			continue
		}
		if userPresence.Status == enum.PresenceStatusPresent.String() {
			totalPresent++
			totalPresentWeek[week]++

			if userPresence.EndTime != "" {
				startTime, err := time.Parse("15:04", userPresence.StartTime)
				if err != nil {
					s.log.Error("invalid start time", zap.Error(err))
					return dto.UserOverviewResponse{}, err
				}

				endTime, err := time.Parse("15:04", userPresence.EndTime)
				if err != nil {
					s.log.Error("invalid end time", zap.Error(err))
					return dto.UserOverviewResponse{}, err
				}

				diffHours := endTime.Sub(startTime).Hours()
				if diffHours > 8.0 {
					totalWorkHour += 8.0
				} else {
					totalWorkHour += diffHours
				}

				totalWorkHourWeek[week] += uint64(diffHours)
			} else {
				totalWorkHour += 8.0
				totalWorkHourWeek[week] += 8.0
			}
		}
	}

	totalWorkDoneWeek := make(map[int]uint64)

	var additionalWorkSalary decimal.Decimal
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

			additionalWorkSalary = additionalWorkSalary.Add(salary)
		}
	}

	cageStaffRole, err := s.repository.GetRoleByName(constant.RolePekerjaKandang)
	if err != nil {
		s.log.Error("failed get role by name", zap.Error(err))
		return dto.UserOverviewResponse{}, err
	}

	var presenceScore float64
	if len(userPresences.Presences) == 0 {
		presenceScore = 0
	} else {
		presenceScore = float64(totalPresent) / float64(len(userPresences.Presences)) * 100
	}

	workPresence := totalWorkHour / float64(8*util.TotalDaysInMonth(int(filter.Year), time.Month(filter.Month)))

	totalSalary := user.Salary.Add(additionalWorkSalary)
	userInformation := dto.UserInformationResponse{
		TotalWorkHour: totalWorkHour,
		WorkKpiScore:  (presenceScore + workPresence) / 2,
	}

	userPresenceInformation := dto.UserPresenceInformationResponse{
		TotalPresent:    totalPresent,
		TotalNotPresent: uint64(len(userPresences.Presences) - int(totalPresent)),
	}

	userWorkInformation := dto.UserWorkInformationResponse{
		TotalWorkDone:    totalWorkDone,
		TotalWorkNotDone: uint64(len(dailyWorkUsers.DailyWorkUsers) + len(additionalWorkUsers.AdditionalWorkUsers) - int(totalWorkDone)),
	}

	var userSalaryInformation dto.UserSalaryInformationResponse
	if time.Month(filter.Month.Value()) == time.Now().Month() && time.Now().Year() == int(filter.Year) {
		userSalaryInformation = dto.UserSalaryInformationResponse{
			BaseSalary:           user.Salary.String(),
			AdditionalWorkSalary: additionalWorkSalary.String(),
			BonusSalary:          decimal.Zero.String(),
			CompentationSalary:   decimal.Zero.String(),
			Cashbond:             decimal.Zero.String(),
			IsPaid:               false,
			TotalSalary:          totalSalary.String(),
		}
	} else {
		userSalary, err := s.repository.GetUserSalaryPaymentSpesificMonth(user.Id, time.Month(filter.Month.Value()), filter.Year)
		if err != nil {
			s.log.Error("failed get user salary payment spesific month", zap.Error(err))
			return dto.UserOverviewResponse{}, err
		}

		userSalaryInformation = dto.UserSalaryInformationResponse{
			BaseSalary:           userSalary.BaseSalary.String(),
			AdditionalWorkSalary: userSalary.AdditionalWorkSalary.String(),
			BonusSalary:          userSalary.BonusSalary.String(),
			CompentationSalary:   userSalary.CompentationSalary.String(),
			Cashbond:             userSalary.Cashbond.String(),
			IsPaid:               true,
			TotalSalary:          userSalary.BaseSalary.Add(userSalary.BonusSalary).Add(userSalary.CompentationSalary).Sub(userSalary.Cashbond).String(),
		}
	}

	keys := util.GetSortedKeysInt(weeks)
	kpiPerformances := make([]dto.KPIPerformanceResponse, 0)
	for _, key := range keys {
		presenceScore := float64(totalPresentWeek[key]) / float64(weeks[key].TotalDays) * 100
		workScore := float64(totalWorkHourWeek[key]) / float64(8*weeks[key].TotalDays) * 100
		kpiPerformance := dto.KPIPerformanceResponse{
			Key:          fmt.Sprintf("Minggu %d", key),
			WorkKpiScore: (presenceScore * 0.6) + (workScore * 0.4),
		}

		kpiPerformances = append(kpiPerformances, kpiPerformance)
	}

	if user.RoleId != cageStaffRole.Id {
		chickenKpi, err := s.chickenService.GetKPIScoreChickenInMonth(uint64(user.LocationId.Int64), filter.Month.Value(), filter.Year)
		if err != nil {
			return dto.UserOverviewResponse{}, err
		}

		userInformation.ChickenKpiScore = chickenKpi

		chickenKpiPerWeek, err := s.chickenService.GetKPIScoreChickenPerWeek(uint64(user.LocationId.Int64), filter.Month.Value(), filter.Year)
		if err != nil {
			return dto.UserOverviewResponse{}, err
		}

		for key := range kpiPerformances {
			kpiPerformances[key].ChickenKpiScore = chickenKpiPerWeek[key]
		}
	}

	placements := make([]string, 0)
	if slices.Contains(entity.CageLocationTypeList, user.Role.Name) {
		data, err := s.placementService.GetCagePlacementByUserId(user.Id)
		if err != nil {
			return dto.UserOverviewResponse{}, err
		}

		for _, e := range data {
			placements = append(placements, fmt.Sprintf("%s - %s", enum.LocationTypeCage.String(), e.Cage.Name))
		}
	} else if slices.Contains(entity.WarehouseLocationTypeList, user.Role.Name) {
		data, err := s.placementService.GetWarehousePlacementByUserId(user.Id)
		if err != nil {
			return dto.UserOverviewResponse{}, err
		}

		for _, e := range data {
			placements = append(placements, fmt.Sprintf("%s - %s", enum.LocationTypeWarehouse.String(), e.Warehouse.Name))
		}
	} else if slices.Contains(entity.StoreLocationTypeList, user.Role.Name) {
		data, err := s.placementService.GetStorePlacementByUserId(user.Id)
		if err != nil {
			return dto.UserOverviewResponse{}, err
		}

		for _, e := range data {
			placements = append(placements, fmt.Sprintf("%s - %s", enum.LocationTypeStore.String(), e.Store.Name))
		}
	}

	userCashAdvances, err := s.cashflowService.GetUserCashAdvanceByUserId(user.Id)
	if err != nil {
		return dto.UserOverviewResponse{}, err
	}

	overviewResponse := dto.UserOverviewResponse{
		User:                    mapper.UserToResponse(&user),
		UseCashAdvances:         userCashAdvances,
		UserInformation:         userInformation,
		Placements:              placements,
		KPIPerformances:         kpiPerformances,
		UserPresenceInformation: userPresenceInformation,
		UserSalaryInformation:   userSalaryInformation,
		UserWorkInformation:     userWorkInformation,
	}

	return overviewResponse, nil
}

func (s *UserService) GetUserOverviewList(filter dto.GetUserOverviewListFilter) (dto.UserListOverviewPaginationResponse, error) {
	s.repository.UseTx(false)

	ownerRole, err := s.roleService.GetRoleByName(constant.RoleOwner)
	if err != nil {
		return dto.UserListOverviewPaginationResponse{}, err
	}

	filter.ExcludeRoleIds = []uint64{ownerRole.Id}
	users, err := s.repository.GetUserOverviewLists(&filter)
	if err != nil {
		s.log.Error("failed to get user overview", zap.Error(err))
		return dto.UserListOverviewPaginationResponse{}, err
	}

	responses := make([]dto.UserListOverviewResponse, 0)
	for _, user := range users {
		response := mapper.UserOverviewToListResponse(&user)

		var (
			withDeleted = true
		)

		additionalWorkUsers, err := s.workService.GetAdditionalWorkUserByUserId(user.Id,
			dto.GetAdditionalWorkUserFilter{
				Month:       param.MonthParam(time.Now().Month()),
				Year:        uint64(time.Now().Year()),
				WithDeleted: &withDeleted,
			})
		if err != nil {
			return dto.UserListOverviewPaginationResponse{}, err
		}

		dailyWorkUsers, err := s.workService.GetDailyWorkUserByUserId(user.Id,
			dto.GetDailyWorkUserFilter{
				Month:       param.MonthParam(time.Now().Month()),
				Year:        uint64(time.Now().Year()),
				WithDeleted: &withDeleted,
			})
		if err != nil {
			return dto.UserListOverviewPaginationResponse{}, err
		}

		userPresences, err := s.presenceService.GetUserPresencesByUserId(user.Id,
			dto.GetPresenceFilter{
				Month: param.MonthParam(time.Now().Month()),
				Year:  uint64(time.Now().Year()),
			})
		if err != nil {
			return dto.UserListOverviewPaginationResponse{}, err
		}

		presenceScore, workScore, _ := util.CalculateKPIScoreUserInMonthViaDTO(additionalWorkUsers, dailyWorkUsers, userPresences)
		kpiPerformance := (0.6 * presenceScore) + (0.4 * workScore)
		if kpiPerformance >= constant.KPIScoreGood {
			response.KpiStatus = constant.KPIStatusGood
		} else if kpiPerformance >= constant.KPIScoreMid && kpiPerformance < constant.KPIScoreGood {
			response.KpiStatus = constant.KPIStatusMid
		} else {
			response.KpiStatus = constant.KPIStatusBad
		}

		responses = append(responses, response)
	}

	totalData, err := s.repository.CountTotalUserOverviewList(&filter)
	if err != nil {
		s.log.Error("failed count user overview", zap.Error(err))
		return dto.UserListOverviewPaginationResponse{}, err
	}

	resp := dto.UserListOverviewPaginationResponse{
		Users: responses,
	}

	if filter.Page > 0 {
		resp.TotalData = uint64(totalData)
		resp.TotalPage = uint64(math.Ceil(float64(totalData) / float64(constant.PaginationDefaultLimit)))
	}

	return resp, nil
}

func (s *UserService) GetUserPerformanceOverview(filter dto.GetUserPerformanceOverviewFilter) (dto.UserPerformanceOverviewResponse, error) {
	s.repository.UseTx(false)

	weeks := util.GetFourWeekRanges(int(filter.Year), time.Month(filter.Month))

	ownerRole, err := s.repository.GetRoleByName(constant.RoleOwner)
	if err != nil {
		s.log.Error("failed get role by name", zap.Error(err))
		return dto.UserPerformanceOverviewResponse{}, err
	}

	users, err := s.repository.GetUsers(&dto.GetUserListFilter{
		ExcluseRoleIds: []uint64{ownerRole.Id},
		LocationId:     filter.LocationId,
	})
	if err != nil {
		s.log.Error("failed get users", zap.Error(err))
		return dto.UserPerformanceOverviewResponse{}, err
	}

	kpiUser := float64(0)
	kpiUserPerWeek := make(map[int]float64)

	kpiChicken, err := s.chickenService.GetKPIScoreChickenInMonth(filter.LocationId, enum.Month(filter.Month.Value()), filter.Year)
	if err != nil {
		return dto.UserPerformanceOverviewResponse{}, err
	}

	kpiChickenPerWeek, err := s.chickenService.GetKPIScoreChickenPerWeek(filter.LocationId, enum.Month(filter.Month.Value()), filter.Year)
	if err != nil {
		return dto.UserPerformanceOverviewResponse{}, err
	}

	for _, user := range users {
		withDeleted := true
		additionalWorkUsers, err := s.workService.GetAdditionalWorkUserByUserId(user.Id,
			dto.GetAdditionalWorkUserFilter{
				Month:       filter.Month,
				Year:        filter.Year,
				WithDeleted: &withDeleted,
			})
		if err != nil {
			return dto.UserPerformanceOverviewResponse{}, nil
		}

		dailyWorkUsers, err := s.workService.GetDailyWorkUserByUserId(user.Id,
			dto.GetDailyWorkUserFilter{
				Month:       filter.Month,
				Year:        filter.Year,
				WithDeleted: &withDeleted,
			})
		if err != nil {
			return dto.UserPerformanceOverviewResponse{}, nil
		}

		userPresences, err := s.presenceService.GetUserPresencesByUserId(user.Id,
			dto.GetPresenceFilter{
				Month: filter.Month,
				Year:  filter.Year,
			})
		if err != nil {
			return dto.UserPerformanceOverviewResponse{}, nil
		}

		totalPresentWeek := make(map[int]uint64)
		totalWorkHourWeek := make(map[int]uint64)

		var totalPresent uint64 = 0
		var totalWorkHour float64 = 0
		for _, userPresence := range userPresences.Presences {
			week := util.FindWeek(userPresence.CreatedAt, weeks)
			if week == 0 {
				continue
			}

			if userPresence.Status == enum.PresenceStatusPresent.String() {
				totalPresent++
				totalPresentWeek[week]++

				if userPresence.EndTime != "" {
					startTime, err := time.Parse("15:04", userPresence.StartTime)
					if err != nil {
						s.log.Error("invalid start time", zap.Error(err))
						return dto.UserPerformanceOverviewResponse{}, err
					}

					endTime, err := time.Parse("15:04", userPresence.EndTime)
					if err != nil {
						s.log.Error("invalid end time", zap.Error(err))
						return dto.UserPerformanceOverviewResponse{}, err
					}

					diffHours := endTime.Sub(startTime).Hours()
					if diffHours > 8 {
						totalWorkHour += 8.0
						totalWorkHourWeek[week] += 8.0
					} else {
						totalWorkHour += diffHours
						totalWorkHourWeek[week] += uint64(diffHours)
					}
				} else {
					totalWorkHourWeek[week] += 8.0
				}
			}
		}

		totalWorkDoneWeek := make(map[int]uint64)
		totalDailyWorkInWeek := make(map[int]uint64)
		totalAdditionalWorkInWeek := make(map[int]uint64)

		var totalWorkDone uint64 = 0
		for _, dailyWorkUser := range dailyWorkUsers.DailyWorkUsers {
			week := util.FindWeek(dailyWorkUser.CreatedAt, weeks)
			if week == 0 {
				continue
			}

			totalDailyWorkInWeek[week]++
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
			totalAdditionalWorkInWeek[week]++
			if additionalWorkUser.IsDone {
				totalWorkDone++
				totalWorkDoneWeek[week]++
			}
		}

		var workScore float64
		if len(dailyWorkUsers.DailyWorkUsers)+len(additionalWorkUsers.AdditionalWorkUsers) == 0 {
			workScore = 0
		} else {
			workScore = float64(totalWorkDone) / float64(len(dailyWorkUsers.DailyWorkUsers)+len(additionalWorkUsers.AdditionalWorkUsers)) * 100
		}

		var presenceScore float64
		if len(userPresences.Presences) == 0 {
			presenceScore = 0
		} else {
			presenceScore = float64(totalWorkHour) / float64(len(userPresences.Presences)*8) * 100
		}

		if kpiUser == 0 {
			kpiUser = (presenceScore * 0.6) + (workScore * 0.4)
		} else {
			kpiUser = ((presenceScore * 0.6) + (workScore * 0.4) + kpiUser) / 2
		}

		for key, value := range weeks {
			presenceScore := float64(totalPresentWeek[key]) / float64(value.TotalDays) * 100
			workScore := float64(totalWorkHourWeek[key]) / float64(8*value.TotalDays) * 100

			if kpiUserPerWeek[key] == 0 {
				kpiUserPerWeek[key] = (presenceScore * 0.6) + (workScore * 0.4)
			} else {
				kpiUserPerWeek[key] = ((presenceScore * 0.6) + (workScore * 0.4) + kpiUserPerWeek[key]) / 2
			}
		}
	}

	performancesGraphsResponses := make([]dto.PerformanceGraphResponse, 0)
	for key := range weeks {
		performancesGraphsResponses = append(performancesGraphsResponses, dto.PerformanceGraphResponse{
			Key:                   fmt.Sprintf("Minggu %d", key),
			KPIChickenPerformance: kpiChickenPerWeek[key],
			KPIUserPerformance:    kpiUserPerWeek[key],
		})

	}

	return dto.UserPerformanceOverviewResponse{
		UserPerformanceDetail: dto.UserPerformanceSummaryResponse{
			TotalUser:  uint64(len(users)),
			KPIUser:    kpiUser,
			KPIChicken: kpiChicken,
			KPIAll:     (kpiChicken + kpiUser) / 2,
		},
		UserPerformanceGraphs: performancesGraphsResponses,
	}, nil
}
