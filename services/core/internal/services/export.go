package services

type ScheduleExportService struct {
	scheduleService *ScheduleService
}

func NewScheduleExportService(scheduleService *ScheduleService) *ScheduleExportService {
	return &ScheduleExportService{scheduleService: scheduleService}
}

type GroupScheduleDate struct {
	GroupName string
}
