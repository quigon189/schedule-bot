package llm

import (
	"context"
	"core/internal/dto"
	"core/internal/services"
	"fmt"
	"log"
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
			Name:        "get_weekly_schedule_group",
			Description: "Получить основное расписание группы (без изменений) на неделю для указаного учебного периода",
			Parameters: map[string]string{
				"group_name":      "string, обязательный, трехзначиное число (например 501)",
				"period_year":     "string, опциональный, учебный период указывающий года обучения в формате YYYY/YYYY (например 2025/2026); если не указан, то подставляется текущий активный период",
				"period_semester": "int, опциональный (должен быть задан, если задан period_year), указывает семестр учебного года, принимает одно из двух значений: 1 или 2",
			},
			Handler: func(ctx context.Context, params map[string]any) (string, error) {
				log.Printf("Параметры: %+v", params)
				groupName, ok := params["group_name"].(string)
				if !ok {
					return "", fmt.Errorf("group_name обязательный параметр")
				}

				var periodID *int
				periodYear, ok := params["period_year"].(string)
				if ok && periodYear != "" {
					semester, okk := params["period_semester"].(float64)
					if !okk {
						return "", fmt.Errorf("period_semester обязательный, если указан period_year")
					}

					period, err := svc.GetAcademicPeriodByYear(ctx, periodYear, int(semester))
					if err != nil {
						return "", fmt.Errorf("get academic period: %w", err)
					}

					periodID = &period.ID
				}

				group, err := svc.GetGroupByName(ctx, groupName)
				if err != nil {
					return "", fmt.Errorf("get group by name: %w", err)
				}

				schedule, err := svc.GetAllScheduleTemplates(ctx, &dto.ScheduleFiltersRequest{
					GroupID:          &group.ID,
					AcademicPeriodID: periodID,
				})

				var builder strings.Builder
				builder.WriteString("Основное распиание занятий для группы " + group.Name + ":\n\n")

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
					for day := range 7 {
						fmt.Fprintf(&builder, "%s: ", dayOfWeek[day])
						for _, s := range schedule {
							if s.DayOfWeek == day {
								fmt.Fprintf(&builder, "Пара №%d: ", s.Number)
								fmt.Fprintf(&builder, "%s ", s.Subject.Title)
								fmt.Fprintf(&builder, "| %s ", s.Audience.Number)
								fmt.Fprintf(&builder, "| %s\n", s.Teacher.User.FullName)

							}
						}
						builder.WriteString("\n")
					}
				} else {
					fmt.Fprintf(&builder, "Расписание отсутствует")
				}

				return builder.String(), nil
			},
		},
		{
			Name:        "get_day_schedule_group",
			Description: "Получить основное расписание группы (без изменений) на заданный день недели для указанного учебного периода",
			Parameters: map[string]string{
				"group_name":      "string, обязательный, название группы трехзначное число (например 501)",
				"day_of_week":     "int, обязательный, порядковый номер дня недели (например для понедельника - 1)",
				"period_year":     "string, опциональный, учебный период указывающий года обучения в формате YYYY/YYYY (например 2025/2026); если не указан, то подставляется текущий активный период",
				"period_semester": "int, опциональный (должен быть задан, если задан period_year), указывает семестр учебного года, принимает одно из двух значений: 1 или 2",
			},
			Handler: func(ctx context.Context, params map[string]any) (string, error) {
				groupName, ok := params["group_name"].(string)
				if !ok {
					return "", fmt.Errorf("group_name обязательный параметр")
				}

				group, err := svc.GetGroupByName(ctx, groupName)
				if err != nil {
					return "", fmt.Errorf("get group by name %s: %w", groupName, err)
				}

				var periodID *int
				periodYear, ok := params["period_year"].(string)
				if ok && periodYear != "" {
					periodSemester, ok := params["period_semester"].(float64)
					if ok {
						period, err := svc.GetAcademicPeriodByYear(ctx, periodYear, int(periodSemester))
						if err != nil {
							return "", fmt.Errorf("get academic period: %w", err)
						}

						periodID = &period.ID
					}
				}

				day, ok := params["day_of_week"].(float64)
				if !ok {
					return "", fmt.Errorf("day_of_week required")
				}

				intDay := int(day)

				schedule, err := svc.GetAllScheduleTemplates(ctx, &dto.ScheduleFiltersRequest{
					GroupID:          &group.ID,
					AcademicPeriodID: periodID,
					DayOfWeek:        &intDay,
				})

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

				fmt.Fprintf(&builder, "Распиание занятий для группы %s на %s:\n\n", group.Name, dayOfWeek[intDay])

				if len(schedule) > 0 {
					for _, s := range schedule {
						fmt.Fprintf(&builder, "Пара №%d: ", s.Number)
						fmt.Fprintf(&builder, "%s ", s.Subject.Title)
						fmt.Fprintf(&builder, "| %s ", s.Audience.Number)
						fmt.Fprintf(&builder, "| %s\n", s.Teacher.User.FullName)
					}
				} else {
					fmt.Fprintf(&builder, "Расписание отсутствует")
				}

				return builder.String(), nil

			},
		},
		{
			Name:        "get_weekly_schedule_teacher",
			Description: "Получить основное расписание преподавателя (без изменений) на неделю для указаного учебного периода",
			Parameters: map[string]string{
				"teacher_lastname":  "string, обязательный, фамилия преподавателя",
				"teacher_firstname": "string, опциональный, имя преподавателя или первая буква имени",
				"period_year":       "string, опциональный, учебный период указывающий года обучения в формате YYYY/YYYY (например 2025/2026); если не указан, то подставляется текущий активный период",
				"period_semester":   "int, опциональный (должен быть задан, если задан period_year), указывает семестр учебного года, принимает одно из двух значений: 1 или 2",
			},
			Handler: func(ctx context.Context, params map[string]any) (string, error) {
				lastName, ok := params["teacher_lastname"].(string)
				if !ok {
					return "", fmt.Errorf("teacher_lastname обязательный параметр")
				}

				firstName, _ := params["teacher_firstname"].(string)

				teacherName := lastName + firstName

				var periodID *int
				periodYear, ok := params["period_year"].(string)
				if ok && periodYear != "" {
					semester, okk := params["period_semester"].(float64)
					if !okk {
						return "", fmt.Errorf("period_semester обязательный, если указан period_year")
					}

					period, err := svc.GetAcademicPeriodByYear(ctx, periodYear, int(semester))
					if err != nil {
						return "", fmt.Errorf("get academic period: %w", err)
					}

					periodID = &period.ID
				}

				teacher, err := svc.GetTeacherByName(ctx, teacherName)
				if err != nil || teacher == nil {
					return "", fmt.Errorf("get teacher by name %s: %w", teacherName, err)
				}

				schedule, err := svc.GetAllScheduleTemplates(ctx, &dto.ScheduleFiltersRequest{
					TeacherID:        &teacher.User.ID,
					AcademicPeriodID: periodID,
				})

				var builder strings.Builder
				builder.WriteString("Основное распиание занятий для преподавателя " + teacher.User.FullName + ":\n\n")

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
					for day := range 7 {
						fmt.Fprintf(&builder, "%s: ", dayOfWeek[day])
						for _, s := range schedule {
							if s.DayOfWeek == day {
								fmt.Fprintf(&builder, "Пара №%d: ", s.Number)
								fmt.Fprintf(&builder, "%s ", s.Subject.Title)
								fmt.Fprintf(&builder, "| %s ", s.Audience.Number)
								fmt.Fprintf(&builder, "| %s\n", s.Subject.Group.Name)

							}
						}
						builder.WriteString("\n")
					}
				} else {
					fmt.Fprintf(&builder, "Расписание отсутствует")
				}

				return builder.String(), nil
			},
		},
		{
			Name:        "get_day_schedule_teacher",
			Description: "Получить основное расписание преподавателя (без изменений) на заданый день недели для указаного учебного периода",
			Parameters: map[string]string{
				"teacher_lastname":  "string, обязательный, фамилия преподавателя",
				"teacher_firstname": "string, опциональный, имя преподавателя или первая буква имени",
				"day_of_week":       "int, обязательный, порядковый номер дня недели (например для понедельника - 1)",
				"period_year":       "string, опциональный, учебный период указывающий года обучения в формате YYYY/YYYY (например 2025/2026); если не указан, то подставляется текущий активный период",
				"period_semester":   "int, опциональный (должен быть задан, если задан period_year), указывает семестр учебного года, принимает одно из двух значений: 1 или 2",
			},
			Handler: func(ctx context.Context, params map[string]any) (string, error) {
				lastName, ok := params["teacher_lastname"].(string)
				if !ok {
					return "", fmt.Errorf("teacher_lastname обязательный параметр")
				}

				firstName, _ := params["teacher_firstname"].(string)

				teacherName := lastName + firstName

				var periodID *int
				periodYear, ok := params["period_year"].(string)
				if ok && periodYear != "" {
					semester, okk := params["period_semester"].(float64)
					if !okk {
						return "", fmt.Errorf("period_semester обязательный, если указан period_year")
					}

					period, err := svc.GetAcademicPeriodByYear(ctx, periodYear, int(semester))
					if err != nil {
						return "", fmt.Errorf("get academic period: %w", err)
					}

					periodID = &period.ID
				}

				day, ok := params["day_of_week"].(float64)
				if !ok {
					return "", fmt.Errorf("day_of_week обязаятельный")
				}
				intDay := int(day)

				teacher, err := svc.GetTeacherByName(ctx, teacherName)
				if err != nil || teacher == nil {
					return "", fmt.Errorf("get teacher by name %s: %w", teacherName, err)
				}

				schedule, err := svc.GetAllScheduleTemplates(ctx, &dto.ScheduleFiltersRequest{
					DayOfWeek:        &intDay,
					TeacherID:        &teacher.User.ID,
					AcademicPeriodID: periodID,
				})

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

				fmt.Fprintf(&builder, "Основное расписание занятий для преподавателя %s\n%s:\n", teacher.User.FullName, dayOfWeek[intDay])

				if len(schedule) > 0 {
					for _, s := range schedule {
						fmt.Fprintf(&builder, "Пара №%d: ", s.Number)
						fmt.Fprintf(&builder, "%s ", s.Subject.Title)
						fmt.Fprintf(&builder, "| %s ", s.Audience.Number)
						fmt.Fprintf(&builder, "| %s\n", s.Subject.Group.Name)
					}
					builder.WriteString("\n")
				} else {
					fmt.Fprintf(&builder, "Расписание отсутствует")
				}

				return builder.String(), nil
			},
		},
		{
			Name:        "get_day_lessons_group",
			Description: "Получить спискок занятий на заданый день для учебной группы",
			Parameters: map[string]string{
				"group_name": "string, обязательный, трехзначное число (например 501)",
				"date":       "string, обязательный, дата в формате ГГГГ-ММ-ДД",
			},
			Handler: func(ctx context.Context, params map[string]any) (string, error) {
				groupName, ok := params["group_name"].(string)
				if !ok {
					return "", fmt.Errorf("group_name required")
				}
				group, err := svc.GetGroupByName(ctx, groupName)
				if err != nil {
					return "", fmt.Errorf("get group by name %s: %w", groupName, err)
				}

				dateStr, ok := params["date"].(string)
				if !ok {
					return "", fmt.Errorf("date required")
				}
				date, err := time.Parse("2006-01-02", dateStr)
				if err != nil {
					return "", fmt.Errorf("parse date: %w", err)
				}

				lessons, err := svc.GetLessonLogs(ctx, &dto.LessonLogFiltersRequest{
					GroupID:  &group.ID,
					DateFrom: &date,
					DateTo:   &date,
				})
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

				fmt.Fprintf(&builder, "Расписание занятий для группы %s\n Дата:%s:\n", group.Name, date.Format("02.01.2006"))
				fmt.Fprintf(&builder, "%s:\n\n", dayOfWeek[int(date.Weekday())])

				if len(lessons) > 0 {
					for _, l := range lessons {
						fmt.Fprintf(&builder, "Пара №%d: ", l.Number)
						fmt.Fprintf(&builder, "%s ", l.Subject.Title)
						fmt.Fprintf(&builder, "| %s ", l.Audience.Number)
						fmt.Fprintf(&builder, "| %s ", l.Teacher.User.FullName)
						fmt.Fprintf(&builder, "| %s\n", lessonStatus[l.Status])
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
