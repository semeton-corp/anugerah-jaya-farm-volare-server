package mapper

import (
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/constant"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
)

func DailyWorkDetailToResponse(dailyWork *entity.DailyWork) dto.DailyWorkDetailResponse {
	return dto.DailyWorkDetailResponse{
		Id:          dailyWork.Id,
		Description: dailyWork.Description,
		StartTime:   dailyWork.StartTime.Time.Format("15:04"),
		EndTime:     dailyWork.EndTime.Time.Format("15:04"),
	}
}

// Note : Without additonal work user information
func AdditionalWorkToResponse(additionalWork *entity.AdditionalWork) dto.AdditionalWorkResponse {
	response := dto.AdditionalWorkResponse{
		Id:           additionalWork.Id,
		Name:         additionalWork.Name,
		Description:  additionalWork.Description,
		Location:     LocationToResponse(&additionalWork.Location),
		LocationType: additionalWork.LocationType.String(),
		Date:         additionalWork.WorkDate.Format("02-01-2006"),
		Time:         additionalWork.WorkDate.Format("15:04"),
		Slot:         additionalWork.Slot,
		Salary:       additionalWork.Salary.String(),
	}

	switch additionalWork.LocationType {
	case enum.LocationWorkTypeCage:
		response.Place = additionalWork.Cage.Name
	case enum.LocationWorkTypeStore:
		response.Place = additionalWork.Store.Name
	case enum.LocationWorkTypeWarehouse:
		response.Place = additionalWork.Warehouse.Name
	}

	return response
}

func AdditionalWorkUserInformationToResponse(additionalWorkUser *entity.AdditionalWorkUser) dto.AdditionalWorkUserInformationResponse {
	return dto.AdditionalWorkUserInformationResponse{
		Id:       additionalWorkUser.Id,
		RoleName: additionalWorkUser.User.Role.Name,
		UserName: additionalWorkUser.User.Name,
		IsDone:   additionalWorkUser.IsDone,
	}
}

// Note : without status
func AdditionalWorkToListResponse(additionalWork *entity.AdditionalWork) dto.AdditionalWorkListResponse {
	response := dto.AdditionalWorkListResponse{
		Id:            additionalWork.Id,
		Date:          additionalWork.WorkDate.Format("02 Jan 2006"),
		Time:          additionalWork.WorkDate.Format("15:04"),
		Name:          additionalWork.Name,
		Location:      additionalWork.Location.Name,
		RemainingSlot: additionalWork.Slot - uint64(len(additionalWork.AdditionalWorkUsers)),
	}

	switch additionalWork.LocationType {
	case enum.LocationWorkTypeCage:
		response.Place = additionalWork.LocationType.String() + ", " + additionalWork.Cage.Name
	case enum.LocationWorkTypeStore:
		response.Place = additionalWork.LocationType.String() + ", " + additionalWork.Store.Name
	case enum.LocationWorkTypeWarehouse:
		response.Place = additionalWork.LocationType.String() + ", " + additionalWork.Warehouse.Name
	}

	if additionalWork.Slot == uint64(len(additionalWork.AdditionalWorkUsers)) {
		response.Status = constant.AdditionalWorkFullWorker
	} else {
		response.Status = constant.AdditionalWorkNeedWorker
	}

	return response
}

func DailyWorkUserToResponse(dailyWorkUser *entity.DailyWorkUser) dto.DailyWorkUserResponse {
	return dto.DailyWorkUserResponse{
		Id:     dailyWorkUser.Id,
		IsDone: dailyWorkUser.IsDone,
		Note:   dailyWorkUser.Note,
		DailyWork: dto.DailyWorkDetailResponse{
			Id:          dailyWorkUser.DailyWork.Id,
			Description: dailyWorkUser.DailyWork.Description,
			StartTime:   dailyWorkUser.DailyWork.StartTime.Time.Format("15:04"),
			EndTime:     dailyWorkUser.DailyWork.EndTime.Time.Format("15:04"),
		},
		CreatedAt: dailyWorkUser.CreatedAt,
	}
}

func AdditionalWorkUserToResponse(additionalWorkUser *entity.AdditionalWorkUser) dto.AdditionalWorkUserResponse {
	return dto.AdditionalWorkUserResponse{
		Id:     additionalWorkUser.Id,
		IsDone: additionalWorkUser.IsDone,
		Note:   additionalWorkUser.Note,
		AdditionalWork: dto.AdditionalWorkDetailResponse{
			Id:          additionalWorkUser.AdditionalWork.Id,
			Description: additionalWorkUser.AdditionalWork.Description,
			Date:        additionalWorkUser.AdditionalWork.CreatedAt.Format("2006-01-02"),
			Time:        additionalWorkUser.AdditionalWork.CreatedAt.Format("15:04"),
			Salary:      additionalWorkUser.AdditionalWork.Salary.String(),
		},
		CreatedAt: additionalWorkUser.CreatedAt,
	}
}
