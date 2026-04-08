package llm

import (
	"context"
	"core/internal/dto"
	"core/internal/services"
	"fmt"
	"slices"
	"strings"
	"time"
)

type Tool struct {
	Name        string                                                           `json:"name"`
	Description string                                                           `json:"description"`
	Parameters  map[string]string                                                `json:"parameters"`
	Handler     func(ctx context.Context, params map[string]any) (string, error) `json:"-"`
}

func initTools(svc services.ScheduleService) []Tool {
	return []Tool{
		{
			Name:        "get_schedule",
			Description: "Получить основное расписание занятий заданное на конкретный учебный период, сформированно заранее на чебную неделю для групп, преподавателей и аудиторий",
			Parameters: map[string]string{
				"group_name":      "string, опциональный, трехзначное число (например 501)",
				"teacher_name":    "string, опциональный, только фамилия преподавателя (может содержать имя или первую букву имени)",
				"audience_number": "string, опциональный, номер аудитории трехзначное число, может в конце содержать букву (например 409 или 502а)",
				"day_of_week":     "int, опциональный, порядковый номер дня недели (например для понедельника - 1)",
				"period_year":     "string, опциональный, учебный период указывающий года обучения в формате YYYY/YYYY (например 2025/2026); если не указан, то подставляется текущий активный период",
				"period_semester": "int, опциональный (должен быть задан, если задан period_year), указывает семестр учебного года, принимает одно из двух значений: 1 или 2",
			},
			Handler: func(ctx context.Context, params map[string]any) (string, error) {
				filter := dto.ScheduleFiltersRequest{}
				if groupName, ok := params["group_name"].(string); ok && groupName != "" {
					group, err := svc.GetGroupByName(ctx, groupName)
					if err != nil {
						return "", fmt.Errorf("get group with name %s: %w", groupName, err)
					}
					if group != nil {
						filter.GroupID = &group.ID
					}
				}
				if teacherName, ok := params["teacher_name"].(string); ok && teacherName != "" {
					parts := strings.Fields(teacherName)
					name := parts[0]
					if len(parts) >= 2 {
						name += string([]rune(parts[1])[0])
					}
					teacher, err := svc.GetTeacherByName(ctx, name)
					if err != nil {
						return "", fmt.Errorf("get teacher with name %s: %w", teacherName, err)
					}
					if teacher != nil {
						filter.TeacherID = &teacher.User.ID
					}
				}
				if audienceNumber, ok := params["audience_number"].(string); ok && audienceNumber != "" {
					audience, err := svc.GetAudienceByNumber(ctx, audienceNumber)
					if err != nil {
						return "", fmt.Errorf("get audience with number %s: %w", audienceNumber, err)
					}
					if audience != nil {
						filter.AudienceID = &audience.ID
					}
				}
				if dayOfWeek, ok := params["day_of_week"].(float64); ok && int(dayOfWeek) != 0 {
					intDay := int(dayOfWeek)
					filter.DayOfWeek = &intDay
				}
				periodYear, _ := params["period_year"].(string)
				periodSemester, _ := params["period_semester"].(float64)
				if periodYear != "" && int(periodSemester) != 0 {
					period, err := svc.GetAcademicPeriodByYear(ctx, periodYear, int(periodSemester))
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
		},
		{
			Name:        "get_lessons",
			Description: "Получить спискок занятий по заданным фильтрам, используй когда пользователь хочет узнать занятия или расписание на определенную дату",
			Parameters: map[string]string{
				"date_from":       "string, опциональный, начальная дата в формате ГГГГ-ММ-ДД",
				"date_to":         "string, опциональный, конечная дата в формате ГГГГ-ММ-ДД",
				"status":          "string, опциональный, одно из 4 значений planned - запланировано на основе расписания, canceled - отменено, rescheduled - подставлено, completed - проведено",
				"number":          "int, опциональный, номер занятия (пары)",
				"group_name":      "string, опциональный, трехзначное число (например 501)",
				"teacher_name":    "string, опциональный, только фамилия преподавателя (может содержать имя или первую букву имени)",
				"audience_number": "string, опциональный, номер аудитории трехзначное число, может в конце содержать букву (например 409 или 502а)",
			},
			Handler: func(ctx context.Context, params map[string]any) (string, error) {
				filter := dto.LessonLogFiltersRequest{}
				if dateF, ok := params["date_from"].(string); ok && dateF != "" {
					dateFrom, err := time.Parse("2006-01-02", dateF)
					if err != nil {
						return "", fmt.Errorf("parse date_from: %w", err)
					}
					filter.DateFrom = &dateFrom
				}
				if dateT, ok := params["date_to"].(string); ok && dateT != "" {
					dateTo, err := time.Parse("2006-01-02", dateT)
					if err != nil {
						return "", fmt.Errorf("parse date_to: %w", err)
					}
					filter.DateTo = &dateTo
				} else {
					if dateF, ok := params["date_from"].(string); ok && dateF != "" {
						dateFrom, err := time.Parse("2006-01-02", dateF)
						if err != nil {
							return "", fmt.Errorf("parse date_from: %w", err)
						}
						filter.DateTo = &dateFrom
					}
				}

				if status, ok := params["status"].(string); ok && status != "" {
					statuses := []string{"planned", "canceled", "rescheduled", "completed"}
					if slices.Contains(statuses, status) {
						filter.Status = &status
					}
				}
				if number, ok := params["number"].(float64); ok && number != 0 {
					intNumber := int(number)
					filter.Number = &intNumber
				}
				if groupName, ok := params["group_name"].(string); ok && groupName != "" {
					group, err := svc.GetGroupByName(ctx, groupName)
					if err != nil {
						return "", fmt.Errorf("get group with name %s: %w", groupName, err)
					}
					if group != nil {
						filter.GroupID = &group.ID
					}
				}
				if teacherName, ok := params["teacher_name"].(string); ok && teacherName != "" {
					parts := strings.Fields(teacherName)
					name := parts[0]
					// if len(parts) >= 2 {
					// 	name += string([]rune(parts[1])[0])
					// }
					teacher, err := svc.GetTeacherByName(ctx, name)
					if err != nil {
						return "", fmt.Errorf("get teacher with name %s: %w", teacherName, err)
					}
					if teacher != nil {
						filter.TeacherID = &teacher.User.ID
					}
				}
				if audienceNumber, ok := params["audience_number"].(string); ok && audienceNumber != "" {
					audience, err := svc.GetAudienceByNumber(ctx, audienceNumber)
					if err != nil {
						return "", fmt.Errorf("get audience with number %s: %w", audienceNumber, err)
					}
					if audience != nil {
						filter.AudienceID = &audience.ID
					}
				}

				periodYear, _ := params["period_year"].(string)
				periodSemester, _ := params["period_semester"].(float64)
				if periodYear != "" && int(periodSemester) != 0 {
					period, err := svc.GetAcademicPeriodByYear(ctx, periodYear, int(periodSemester))
					if err == nil {
						filter.AcademicPeriodID = &period.ID
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
		},
	}
}
