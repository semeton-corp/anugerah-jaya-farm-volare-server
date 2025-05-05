package mapper

import (
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
)

func DailyWorkDetailToResponse(dailyWork *entity.DailyWork) dto.DailyWorkDetailResponse {
	return dto.DailyWorkDetailResponse{
		Id:          dailyWork.Id,
		Description: dailyWork.Description,
		StartTime:   dailyWork.StartTime.Time.Format("15:04"),
		EndTime:     dailyWork.EndTime.Time.Format("15:04"),
	}
}

// Note : Without additonal work staff
func AdditionalWorkToResponse(additionalWork *entity.AdditionalWork) dto.AdditionalWorkResponse {
	return dto.AdditionalWorkResponse{
		Id:          additionalWork.Id,
		Description: additionalWork.Description,
		Location:    additionalWork.Location.String(),
		Slot:        additionalWork.Slot,
	}
}

func AdditionalWorkStaffInformationToResponse(additionalWorkStaff *entity.AdditionalWorkStaff) dto.AdditionalWorkStaffInformationResponse {
	return dto.AdditionalWorkStaffInformationResponse{
		Id:        additionalWorkStaff.Id,
		Date:      additionalWorkStaff.CreatedAt.Format("2006-01-02"),
		Time:      additionalWorkStaff.CreatedAt.Format("15:04"),
		StaffName: additionalWorkStaff.Staff.Name,
		IsDone:    additionalWorkStaff.IsDone,
	}
}

// Note : without status
func AdditionalWorkToListResponse(additionalWork *entity.AdditionalWork) dto.AdditionalWorkListResponse {
	return dto.AdditionalWorkListResponse{
		Id:            additionalWork.Id,
		Date:          additionalWork.CreatedAt.Format("02 Jan 2006"),
		Description:   additionalWork.Description,
		Location:      additionalWork.Location.String(),
		RemainingSlot: additionalWork.Slot - uint64(len(additionalWork.AdditionalWorkStaff)),
	}
}

func DailyWorkStaffToResponse(dailyWorkStaff *entity.DailyWorkStaff) dto.DailyWorkStaffResponse {
	return dto.DailyWorkStaffResponse{
		Id:     dailyWorkStaff.Id,
		IsDone: dailyWorkStaff.IsDone,
		DailyWork: dto.DailyWorkDetailResponse{
			Id:          dailyWorkStaff.DailyWork.Id,
			Description: dailyWorkStaff.DailyWork.Description,
			StartTime:   dailyWorkStaff.DailyWork.StartTime.Time.Format("15:04"),
			EndTime:     dailyWorkStaff.DailyWork.EndTime.Time.Format("15:04"),
		},
	}
}

func AdditionalWorkStaffToResponse(dailyWorkStaff *entity.AdditionalWorkStaff) dto.AdditionalWorkStaffResponse {
	return dto.AdditionalWorkStaffResponse{
		Id:     dailyWorkStaff.Id,
		IsDone: dailyWorkStaff.IsDone,
		AdditionalWork: dto.AdditionalWorkDetailResponse{
			Id:          dailyWorkStaff.AdditionalWork.Id,
			Description: dailyWorkStaff.AdditionalWork.Description,
			Date:        dailyWorkStaff.AdditionalWork.CreatedAt.Format("2006-01-02"),
			Time:        dailyWorkStaff.AdditionalWork.CreatedAt.Format("15:04"),
		},
	}
}
