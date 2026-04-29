package llm

import (
	"context"
	"core/internal/dto"
	"core/internal/services"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"time"
)

type Tool struct {
	Name        string                                                            `json:"name"`
	Description string                                                            `json:"description"`
	Parameters  any                                                               `json:"parameters"`
	Handler     func(ctx context.Context, params json.RawMessage) (string, error) `json:"-"`
}

func initTools(svc services.ScheduleService) []Tool {
	return []Tool{
		getScheduleTool(svc),
		getLessonsTool(svc),
	}
}

func getScheduleTool(svc services.ScheduleService) Tool {
	type getScheduleParams struct {
		GroupName      string `json:"group_name" jsonschema:"description=трехзначное число (например 501)"`
		TeacherName    string `json:"teacher_name" jsonschema:"description=только фамилия преподавателя (может содержать имя или первую букву имени)"`
		AudienceNumber string `json:"audience_number" jsonschema:"description=номер аудитории трехзначное число, может в конце содержать букву (например 409 или 502а)"`
		DayOfWeek      int    `json:"day_of_week" jsonschema:"description=порядковый номер дня недели (например для понедельника - 1)"`
		PeriodYear     string `json:"period_year" jsonschema:"description=учебный период указывающий года обучения в формате YYYY/YYYY (например 2025/2026); если не указан, то подставляется текущий активный период"`
		PeriodSemester int    `json:"period_semester" jsonschema:"description=указывает семестр учебного года, принимает одно из двух значений: 1 или 2"`
	}
	return Tool{
		Name:        "get_schedule",
		Description: "Получить основное расписание занятий заданное на конкретный учебный период, сформированно заранее на учебную неделю для групп, преподавателей и аудиторий",
		Parameters:  getScheduleParams{},
		Handler: func(ctx context.Context, raw json.RawMessage) (string, error) {
			var params getScheduleParams
			if err := json.Unmarshal(raw, &params); err != nil {
				return "", fmt.Errorf("unmarshal params: %w", err)
			}
			filter := dto.ScheduleFiltersRequest{}
			if params.GroupName != "" {
				group, err := svc.GetGroupByName(ctx, params.GroupName)
				if err != nil {
					return "", fmt.Errorf("get group with name %s: %w", params.GroupName, err)
				}
				if group != nil {
					filter.GroupID = &group.ID
				}
			}
			if params.TeacherName != "" {
				parts := strings.Fields(params.TeacherName)
				name := parts[0]
				if len(parts) >= 2 {
					name += string([]rune(parts[1])[0])
				}
				teacher, err := svc.GetTeacherByName(ctx, name)
				if err != nil {
					return "", fmt.Errorf("get teacher with name %s: %w", params.TeacherName, err)
				}
				if teacher != nil {
					filter.TeacherID = &teacher.User.ID
				}
			}
			if params.AudienceNumber != "" {
				audience, err := svc.GetAudienceByNumber(ctx, params.AudienceNumber)
				if err != nil {
					return "", fmt.Errorf("get audience with number %s: %w", params.AudienceNumber, err)
				}
				if audience != nil {
					filter.AudienceID = &audience.ID
				}
			}
			if params.DayOfWeek != 0 {
				filter.DayOfWeek = &params.DayOfWeek
			}
			if params.PeriodYear != "" && params.PeriodSemester != 0 {
				period, err := svc.GetAcademicPeriodByYear(ctx, params.PeriodYear, params.PeriodSemester)
				if err == nil {
					filter.AcademicPeriodID = &period.ID
				}
			}

			schedule, err := svc.GetAllScheduleTemplates(ctx, &filter)
			if err != nil {
				return "", fmt.Errorf("get schedule: %w", err)
			}

			var builder strings.Builder
			builder.WriteString("Расписание:\n\n")

			dayOfWeek := map[int]string{
				1: "Понедельник",
				2: "Вторник",
				3: "Среда",
				4: "Четверг",
				5: "Пятница",
				6: "Суббота",
				7: "Воскресенье",
			}

			if len(schedule) > 0 {
				for _, s := range schedule {
					fmt.Fprintf(&builder, "%s: ", dayOfWeek[s.DayOfWeek])
					fmt.Fprintf(&builder, "Пара №%d: ", s.Number)
					fmt.Fprintf(&builder, "Предмет: %s ", s.Subject.Title)
					fmt.Fprintf(&builder, "| Аудитория: %s ", s.Audience.Number)
					fmt.Fprintf(&builder, "| Преподаватель: %s ", s.Teacher.User.FullName)
					fmt.Fprintf(&builder, "| Группа: %s\n", s.Subject.Group.Name)
				}
			} else {
				fmt.Fprintf(&builder, "Расписание отсутствует")
			}

			return builder.String(), nil

		},
	}
}

func getLessonsTool(svc services.ScheduleService) Tool {
	type getLessonsParams struct {
		DateFrom       string `json:"date_from" jsonschema:"description=начальная дата в формате ГГГГ-ММ-ДД"`
		DateTo         string `json:"date_to" jsonschema:"description=конечная дата в формате ГГГГ-ММ-ДД"`
		Status         string `json:"status" jsonschema:"description=одно из 4 значений planned - запланировано на основе расписания, canceled - отменено, rescheduled - подставлено, completed - проведено"`
		Number         int    `json:"number" jsonschema:"description=номер занятия (пары)"`
		GroupName      string `json:"group_name" jsonschema:"description=трехзначное число (например 501)"`
		TeacherName    string `json:"teacher_name" jsonschema:"description=только фамилия преподавателя (может содержать имя или первую букву имени)"`
		AudienceNumber string `json:"audience_number" jsonschema:"description=номер аудитории трехзначное число, может в конце содержать букву (например 409 или 502а)"`
	}

	return Tool{
		Name:        "get_lessons",
		Description: "Получить спискок занятий по заданным фильтрам, используй когда пользователь хочет узнать занятия или расписание на определенную дату",
		Parameters:  getLessonsParams{},
		Handler: func(ctx context.Context, raw json.RawMessage) (string, error) {
			var params getLessonsParams
			if err := json.Unmarshal(raw, &params); err != nil {
				return "", fmt.Errorf("unmarshal params: %w", err)
			}
			filter := dto.LessonLogFiltersRequest{}
			if params.DateFrom != "" {
				dateFrom, err := time.Parse("2006-01-02", params.DateFrom)
				if err != nil {
					return "", fmt.Errorf("parse date_from: %w", err)
				}
				filter.DateFrom = &dateFrom
			}
			if params.DateTo != "" {
				dateTo, err := time.Parse("2006-01-02", params.DateTo)
				if err != nil {
					return "", fmt.Errorf("parse date_to: %w", err)
				}
				filter.DateTo = &dateTo
			} else {
				if params.DateFrom != "" {
					dateFrom, err := time.Parse("2006-01-02", params.DateFrom)
					if err != nil {
						return "", fmt.Errorf("parse date_from: %w", err)
					}
					filter.DateTo = &dateFrom
				}
			}

			if params.Status != "" {
				statuses := []string{"planned", "canceled", "rescheduled", "completed"}
				if slices.Contains(statuses, params.Status) {
					filter.Status = &params.Status
				}
			}
			if params.Number != 0 {
				filter.Number = &params.Number
			}
			if params.GroupName != "" {
				group, err := svc.GetGroupByName(ctx, params.GroupName)
				if err != nil {
					return "", fmt.Errorf("get group with name %s: %w", params.GroupName, err)
				}
				if group != nil {
					filter.GroupID = &group.ID
				}
			}
			if params.TeacherName != "" {
				parts := strings.Fields(params.TeacherName)
				name := parts[0]
				teacher, err := svc.GetTeacherByName(ctx, name)
				if err != nil {
					return "", fmt.Errorf("get teacher with name %s: %w", params.TeacherName, err)
				}
				if teacher != nil {
					filter.TeacherID = &teacher.User.ID
				}
			}
			if params.AudienceNumber != "" {
				audience, err := svc.GetAudienceByNumber(ctx, params.AudienceNumber)
				if err != nil {
					return "", fmt.Errorf("get audience with number %s: %w", params.AudienceNumber, err)
				}
				if audience != nil {
					filter.AudienceID = &audience.ID
				}
			}

			lessons, err := svc.GetLessonLogs(ctx, &filter)
			if err != nil {
				return "", fmt.Errorf("get lessons: %w", err)
			}

			var builder strings.Builder
			dayOfWeek := map[int]string{
				1: "Понедельник",
				2: "Вторник",
				3: "Среда",
				4: "Четверг",
				5: "Пятница",
				6: "Суббота",
				7: "Воскресенье",
			}

			lessonStatus := map[string]string{
				"planned":     "По расписанию",
				"canceled":    "Отменена",
				"rescheduled": "Подставлена",
				"completed":   "Проведена",
			}

			fmt.Fprintf(&builder, "Расписание занятий\n")

			if len(lessons) > 0 {
				for _, l := range lessons {
					fmt.Fprintf(&builder, "Дата: %s | %s | ", l.Date.Format("02.01.2006"), dayOfWeek[int(l.Date.Weekday())])
					fmt.Fprintf(&builder, "Пара №%d: ", l.Number)
					fmt.Fprintf(&builder, "Предмет: %s ", l.Subject.Title)
					fmt.Fprintf(&builder, "| Аудитория: %s ", l.Audience.Number)
					fmt.Fprintf(&builder, "| Преподаватель: %s ", l.Teacher.User.FullName)
					fmt.Fprintf(&builder, "| Группа: %s ", l.Subject.Group.Name)
					fmt.Fprintf(&builder, "| Стутус занятия: %s\n", lessonStatus[l.Status])
				}
				builder.WriteString("\n")
			} else {
				fmt.Fprintf(&builder, "Занятия отсутствуют")
			}

			return builder.String(), nil

		},
	}

}
