package services

import (
	"bytes"
	"context"
	"core/internal/dto"
	"core/internal/models"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

type ExcelService struct{}

func NewExcelService() *ExcelService {
	return &ExcelService{}
}

func (s *ExcelService) GenerateGroupTemplate() ([]byte, error) {
	f := excelize.NewFile()
	defer func() { _ = f.Close() }()

	// Лист 1: Общая информация о группе
	sheetInfo := "Информация о группе"
	_, err := f.NewSheet(sheetInfo)
	if err != nil {
		return nil, fmt.Errorf("create sheet info: %w", err)
	}
	f.SetCellStr(sheetInfo, "A1", "name")
	f.SetCellStr(sheetInfo, "B1", "specialty")
	f.SetCellStr(sheetInfo, "C1", "admission_year")
	f.SetCellStr(sheetInfo, "A2", "Пример: ИС-51")
	f.SetCellStr(sheetInfo, "B2", "Информационные системы")
	f.SetCellStr(sheetInfo, "C2", "2025")
	f.SetColWidth(sheetInfo, "A", "A", 20)
	f.SetColWidth(sheetInfo, "B", "B", 30)
	f.SetColWidth(sheetInfo, "C", "C", 15)

	// Лист 2: Дисциплины
	sheetSubjects := "Дисциплины"
	_, err = f.NewSheet(sheetSubjects)
	if err != nil {
		return nil, fmt.Errorf("create sheet subjects: %w", err)
	}
	headers := []string{"title", "semester", "hours_load", "start_date", "end_date"}
	for i, h := range headers {
		col := string(rune('A' + i))
		f.SetCellStr(sheetSubjects, col+"1", h)
	}
	// Пример
	f.SetCellStr(sheetSubjects, "A2", "Математика")
	f.SetCellInt(sheetSubjects, "B2", 1)
	f.SetCellInt(sheetSubjects, "C2", 72)
	f.SetCellStr(sheetSubjects, "D2", "2025-09-01")
	f.SetCellStr(sheetSubjects, "E2", "2025-12-31")
	f.SetColWidth(sheetSubjects, "A", "E", 15)

	// Лист 3: Студенты
	sheetStudents := "Студенты"
	_, err = f.NewSheet(sheetStudents)
	if err != nil {
		return nil, fmt.Errorf("create sheet students: %w", err)
	}
	f.SetCellStr(sheetStudents, "A1", "full_name")
	f.SetCellStr(sheetStudents, "B1", "email")
	f.SetCellStr(sheetStudents, "A2", "Иванов Иван Иванович")
	f.SetCellStr(sheetStudents, "B2", "ivanov@example.com")
	f.SetColWidth(sheetStudents, "A", "B", 25)

	f.DeleteSheet("Sheet1")
	// Устанавливаем активный лист
	f.SetActiveSheet(0)

	// Сохраняем в буфер
	buf := new(bytes.Buffer)
	if err := f.Write(buf); err != nil {
		return nil, fmt.Errorf("write excel: %w", err)
	}
	return buf.Bytes(), nil
}

func (s *ExcelService) ParseGroupFromExcel(r io.Reader) (*dto.CreateGroupWithCurriculumRequest, error) {
	f, err := excelize.OpenReader(r)
	if err != nil {
		return nil, fmt.Errorf("open excel: %w", err)
	}
	defer func() { _ = f.Close() }()

	// 1. Информация о группе
	group, err := s.parseGroupInfo(f)
	if err != nil {
		return nil, fmt.Errorf("parse group info: %w", err)
	}

	// 2. Дисциплины
	subjects, err := s.parseSubjects(f)
	if err != nil {
		return nil, fmt.Errorf("parse subjects: %w", err)
	}
	if len(subjects) == 0 {
		return nil, fmt.Errorf("subjects sheet must contain at least one discipline")
	}

	// 3. Студенты
	students, err := s.parseStudents(f)
	if err != nil {
		return nil, fmt.Errorf("parse students: %w", err)
	}
	if len(students) == 0 {
		return nil, fmt.Errorf("students sheet must contain at least one student")
	}

	return &dto.CreateGroupWithCurriculumRequest{
		Group:    *group,
		Subjects: subjects,
		Students: students,
	}, nil
}

func (s *ExcelService) parseGroupInfo(f *excelize.File) (*dto.CreateGroupRequest, error) {
	sheet := "Информация о группе"
	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, fmt.Errorf("get rows from %s: %w", sheet, err)
	}
	if len(rows) < 2 {
		return nil, fmt.Errorf("sheet '%s' must have at least 2 rows (header + data)", sheet)
	}
	// Первая строка — заголовки, вторая — данные
	data := rows[1]
	if len(data) < 3 {
		return nil, fmt.Errorf("not enough columns in group info: need name, specialty, admission_year")
	}
	admissionYear, err := strconv.Atoi(strings.TrimSpace(data[2]))
	if err != nil {
		return nil, fmt.Errorf("admission_year must be integer: %w", err)
	}
	return &dto.CreateGroupRequest{
		Name:          strings.TrimSpace(data[0]),
		Specialty:     strings.TrimSpace(data[1]),
		AdmissionYear: admissionYear,
	}, nil
}

func (s *ExcelService) parseSubjects(f *excelize.File) ([]dto.CreateGroupSubjectRequest, error) {
	sheet := "Дисциплины"
	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, fmt.Errorf("get rows from %s: %w", sheet, err)
	}
	if len(rows) < 2 {
		return nil, fmt.Errorf("sheet '%s' must have header and at least one discipline", sheet)
	}
	var subjects []dto.CreateGroupSubjectRequest
	// Пропускаем заголовок
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		if len(row) < 5 {
			continue // пропускаем неполные строки
		}
		semester, err := strconv.Atoi(strings.TrimSpace(row[1]))
		if err != nil {
			return nil, fmt.Errorf("row %d: semester must be integer: %w", i+1, err)
		}
		hoursLoad, err := strconv.Atoi(strings.TrimSpace(row[2]))
		if err != nil {
			return nil, fmt.Errorf("row %d: hours_load must be integer: %w", i+1, err)
		}
		startDate := strings.TrimSpace(row[3])
		endDate := strings.TrimSpace(row[4])

		// Валидация формата дат
		if _, err := time.Parse("2006-01-02", startDate); err != nil {
			return nil, fmt.Errorf("row %d: start_date must be in YYYY-MM-DD format", i+1)
		}
		if _, err := time.Parse("2006-01-02", endDate); err != nil {
			return nil, fmt.Errorf("row %d: end_date must be in YYYY-MM-DD format", i+1)
		}

		subjects = append(subjects, dto.CreateGroupSubjectRequest{
			Title:     strings.TrimSpace(row[0]),
			Semester:  semester,
			HoursLoad: hoursLoad,
			StartDate: startDate,
			EndDate:   endDate,
		})
	}
	return subjects, nil
}

func (s *ExcelService) parseStudents(f *excelize.File) ([]dto.CreateStudentWithCredentialsRequest, error) {
	sheet := "Студенты"
	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, fmt.Errorf("get rows from %s: %w", sheet, err)
	}
	if len(rows) < 2 {
		return nil, fmt.Errorf("sheet '%s' must have header and at least one student", sheet)
	}
	var students []dto.CreateStudentWithCredentialsRequest
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		if len(row) < 2 {
			continue
		}
		fullName := strings.TrimSpace(row[0])
		email := strings.TrimSpace(row[1])
		if fullName == "" || email == "" {
			continue
		}
		students = append(students, dto.CreateStudentWithCredentialsRequest{
			FullName: fullName,
			Email:    email,
		})
	}
	return students, nil
}

// GenerateStudentTemplate создаёт Excel-шаблон для студентов
func (s *ExcelService) GenerateStudentTemplate() ([]byte, error) {
	f := excelize.NewFile()
	defer f.Close()

	sheet := "Студенты"
	f.SetSheetName("Sheet1", sheet)

	// Заголовки
	f.SetCellStr(sheet, "A1", "full_name")
	f.SetCellStr(sheet, "B1", "email")
	f.SetCellStr(sheet, "C1", "group_name")

	// Пример
	f.SetCellStr(sheet, "A2", "Иванов Иван Иванович")
	f.SetCellStr(sheet, "B2", "ivanov@example.com")
	f.SetCellStr(sheet, "C2", "ИС-51")

	f.SetColWidth(sheet, "A", "C", 25)

	buf := new(bytes.Buffer)
	if err := f.Write(buf); err != nil {
		return nil, fmt.Errorf("write excel: %w", err)
	}
	return buf.Bytes(), nil
}

// GenerateTeacherTemplate создаёт Excel-шаблон для преподавателей
func (s *ExcelService) GenerateTeacherTemplate() ([]byte, error) {
	f := excelize.NewFile()
	defer f.Close()

	sheet := "Преподаватели"
	f.SetSheetName("Sheet1", sheet)

	f.SetCellStr(sheet, "A1", "full_name")
	f.SetCellStr(sheet, "B1", "email")

	f.SetCellStr(sheet, "A2", "Петров Петр Петрович")
	f.SetCellStr(sheet, "B2", "petrov@example.com")

	f.SetColWidth(sheet, "A", "B", 30)

	buf := new(bytes.Buffer)
	if err := f.Write(buf); err != nil {
		return nil, fmt.Errorf("write excel: %w", err)
	}
	return buf.Bytes(), nil
}

// ParseStudentsFromExcel читает файл и возвращает список данных студентов с group_name
func (s *ExcelService) ParseStudentsFromExcel(r io.Reader) ([]dto.StudentImport, error) {
	f, err := excelize.OpenReader(r)
	if err != nil {
		return nil, fmt.Errorf("open excel: %w", err)
	}
	defer f.Close()

	sheet := "Студенты"
	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, fmt.Errorf("get rows: %w", err)
	}
	if len(rows) < 2 {
		return nil, fmt.Errorf("sheet must have header and at least one student")
	}

	var students []dto.StudentImport
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		if len(row) < 3 {
			continue
		}
		fullName := strings.TrimSpace(row[0])
		email := strings.TrimSpace(row[1])
		groupName := strings.TrimSpace(row[2])
		if fullName == "" || email == "" || groupName == "" {
			continue
		}
		students = append(students, dto.StudentImport{
			FullName:  fullName,
			Email:     email,
			GroupName: groupName,
		})
	}
	if len(students) == 0 {
		return nil, fmt.Errorf("no valid student rows found")
	}
	return students, nil
}

// ParseTeachersFromExcel читает файл и возвращает список данных преподавателей
func (s *ExcelService) ParseTeachersFromExcel(r io.Reader) ([]dto.TeacherImport, error) {
	f, err := excelize.OpenReader(r)
	if err != nil {
		return nil, fmt.Errorf("open excel: %w", err)
	}
	defer f.Close()

	sheet := "Преподаватели"
	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, fmt.Errorf("get rows: %w", err)
	}
	if len(rows) < 2 {
		return nil, fmt.Errorf("sheet must have header and at least one teacher")
	}

	var teachers []dto.TeacherImport
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		if len(row) < 2 {
			continue
		}
		fullName := strings.TrimSpace(row[0])
		email := strings.TrimSpace(row[1])
		if fullName == "" || email == "" {
			continue
		}
		teachers = append(teachers, dto.TeacherImport{
			FullName: fullName,
			Email:    email,
		})
	}
	if len(teachers) == 0 {
		return nil, fmt.Errorf("no valid teacher rows found")
	}
	return teachers, nil
}

// GenerateAudienceTemplate создаёт Excel-шаблон для аудиторий
func (s *ExcelService) GenerateAudienceTemplate() ([]byte, error) {
	f := excelize.NewFile()
	defer f.Close()

	sheet := "Аудитории"
	f.SetSheetName("Sheet1", sheet)

	f.SetCellStr(sheet, "A1", "name")
	f.SetCellStr(sheet, "B1", "number")

	f.SetCellStr(sheet, "A2", "Аудитория 101")
	f.SetCellStr(sheet, "B2", "101")

	f.SetColWidth(sheet, "A", "B", 25)

	buf := new(bytes.Buffer)
	if err := f.Write(buf); err != nil {
		return nil, fmt.Errorf("write excel: %w", err)
	}
	return buf.Bytes(), nil
}

// ParseAudiencesFromExcel читает файл и возвращает список аудиторий
func (s *ExcelService) ParseAudiencesFromExcel(r io.Reader) ([]dto.AudienceImport, error) {
	f, err := excelize.OpenReader(r)
	if err != nil {
		return nil, fmt.Errorf("open excel: %w", err)
	}
	defer f.Close()

	sheet := "Аудитории"
	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, fmt.Errorf("get rows: %w", err)
	}
	if len(rows) < 2 {
		return nil, fmt.Errorf("sheet must have header and at least one audience")
	}

	var audiences []dto.AudienceImport
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		if len(row) < 2 {
			continue
		}
		name := strings.TrimSpace(row[0])
		number := strings.TrimSpace(row[1])
		if name == "" || number == "" {
			continue
		}
		audiences = append(audiences, dto.AudienceImport{
			Name:   name,
			Number: number,
		})
	}
	if len(audiences) == 0 {
		return nil, fmt.Errorf("no valid audience rows found")
	}
	return audiences, nil
}

// GeneratePlannerTemplate создаёт Excel-шаблон для генерации расписания
func (s *ExcelService) GeneratePlannerTemplate(ctx context.Context, scheduleService *ScheduleService, period *models.AcademicPeriod) ([]byte, error) {
	f := excelize.NewFile()
	defer f.Close()

	// 1. Лист "Параметры"
	paramsSheet := "Параметры"
	_, err := f.NewSheet(paramsSheet)
	if err != nil {
		return nil, fmt.Errorf("create params sheet: %w", err)
	}
	f.SetCellStr(paramsSheet, "A1", "academic_period_id")
	f.SetCellStr(paramsSheet, "B1", "academic_period_name")
	f.SetCellStr(paramsSheet, "C1", "days")
	f.SetCellStr(paramsSheet, "D1", "slots")
	f.SetCellStr(paramsSheet, "E1", "seed")

	f.SetColWidth(paramsSheet, "A", "A", 20)
	f.SetColWidth(paramsSheet, "B", "B", 30)

	periodName := fmt.Sprintf("%s (%d семестр)", period.Year, period.Semester)
	// Пример заполнения (можно оставить пустым, пользователь выберет)
	f.SetCellInt(paramsSheet, "A2", int64(period.ID))
	f.SetCellStr(paramsSheet, "B2", periodName)
	f.SetCellInt(paramsSheet, "C2", 5)
	f.SetCellInt(paramsSheet, "D2", 5)
	f.SetCellStr(paramsSheet, "E2", "")
	f.SetColVisible(paramsSheet, "A", false)

	// listsSheet := "Справочник_учебные_периоды"
	// _, err = f.NewSheet(listsSheet)
	// if err != nil {
	// 	return nil, fmt.Errorf("create lists sheet: %w", err)
	// }
	// for i, p := range periods {
	// 	f.SetCellStr(listsSheet, fmt.Sprintf("A%d", i+1), periodNames[i])
	// 	f.SetCellInt(listsSheet, fmt.Sprintf("B%d", i+1), int64(p.ID))
	// }
	//
	// dv := excelize.NewDataValidation(true)
	// dv.Sqref = "B2"
	// dv.SetSqrefDropList(fmt.Sprintf("%s!A1:A%d", listsSheet, len(periods)))
	// f.AddDataValidation(paramsSheet, dv)
	//
	// formula := fmt.Sprintf("=VLOOKUP(B2, %s!A1:B%d, 2, FALSE)", listsSheet, len(periods))
	// f.SetCellFormula(paramsSheet, "A2", formula)
	// f.SetColVisible(paramsSheet, "A", false)
	// f.SetSheetVisible(listsSheet, false)

	// 2. Лист "Данные" (основной)
	dataSheet := "Данные"
	_, err = f.NewSheet(dataSheet)
	if err != nil {
		return nil, fmt.Errorf("create data sheet: %w", err)
	}
	headers := []struct {
		Text  string
		Width float64
	}{
		{Text: "Группа", Width: 15},
		{Text: "Дисциплина", Width: 40},
		{Text: "Семестр", Width: 10},
		{Text: "Преподаватель", Width: 40},
		{Text: "Аудитория", Width: 12},
		{Text: "Приоритет", Width: 12},
		{Text: "Количество_пар", Width: 20},
		{Text: "Subject_ID", Width: 10},
		{Text: "Teacher_ID", Width: 10},
		{Text: "Audience_ID", Width: 10},
	}
	for i, h := range headers {
		col := string(rune('A' + i))
		f.SetCellStr(dataSheet, col+"1", h.Text)
		f.SetColWidth(dataSheet, col, col, h.Width)
	}

	subjects, err := scheduleService.GetSubjectsByDate(ctx, period.StartDate, period.EndDate)
	if err != nil {
		return nil, fmt.Errorf("get subject for academic period: %w", err)
	}

	for i, s := range subjects {
		f.SetCellStr(dataSheet, fmt.Sprintf("A%d", i+2), s.Group.Name)
		f.SetCellStr(dataSheet, fmt.Sprintf("B%d", i+2), s.Title)
		f.SetCellInt(dataSheet, fmt.Sprintf("C%d", i+2), int64(s.Semester))

		diff := s.EndDate.Sub(s.StartDate)
		weeks := math.Ceil(diff.Hours() / 24 / 7)
		lessonCount := math.Ceil(float64(s.HoursLoad) / weeks)

		f.SetCellInt(dataSheet, fmt.Sprintf("G%d", i+2), int64(lessonCount))
		f.SetCellInt(dataSheet, fmt.Sprintf("H%d", i+2), int64(s.ID))
	}
	f.SetColVisible(dataSheet, "H", false)
	f.SetColVisible(dataSheet, "I", false)
	f.SetColVisible(dataSheet, "J", false)

	// Преподаватели
	teachersSheet := "Справочник_преподавателей"
	_, err = f.NewSheet(teachersSheet)
	if err != nil {
		return nil, fmt.Errorf("create teachers sheet: %w", err)
	}
	teachers, err := scheduleService.GetAllTeachers(ctx)
	if err != nil {
		return nil, fmt.Errorf("get teachers: %w", err)
	}
	for i, t := range teachers {
		f.SetCellStr(teachersSheet, fmt.Sprintf("A%d", i+1), t.User.FullName)
		f.SetCellInt(teachersSheet, fmt.Sprintf("B%d", i+1), int64(t.User.ID))
	}

	dvT := excelize.NewDataValidation(true)
	dvT.Sqref = "D2:D10000"
	dvT.SetSqrefDropList(fmt.Sprintf("%s!A$1:A%d", teachersSheet, len(teachers)))
	f.AddDataValidation(dataSheet, dvT)

	for i := 2; i < 10000; i++ {
		formula := fmt.Sprintf("=VLOOKUP(D%d, %s!A1:B%d, 2, FALSE)", i, teachersSheet, len(teachers))
		f.SetCellFormula(dataSheet, fmt.Sprintf("I%d", i), formula)
	}
	f.SetSheetVisible(teachersSheet, false)

	// Аудитории
	audiencesSheet := "Справочник_аудиторий"
	_, err = f.NewSheet(audiencesSheet)
	if err != nil {
		return nil, fmt.Errorf("create audiences sheet: %w", err)
	}
	audiences, err := scheduleService.GetAllAudience(ctx)
	if err != nil {
		return nil, fmt.Errorf("get audiences: %w", err)
	}
	for i, a := range audiences {
		f.SetCellStr(audiencesSheet, fmt.Sprintf("A%d", i+1), a.Number)
		f.SetCellInt(audiencesSheet, fmt.Sprintf("B%d", i+1), int64(a.ID))
	}

	symb := "@"
	customStyle, err := f.NewStyle(&excelize.Style{
		CustomNumFmt: &symb,
	})
	if err != nil {
		return nil, fmt.Errorf("create style: %w", err)
	}
	dvA := excelize.NewDataValidation(true)
	dvA.Sqref = "E2:E10000"
	dvA.SetSqrefDropList(fmt.Sprintf("%s!A$1:A%d", audiencesSheet, len(audiences)))
	f.AddDataValidation(dataSheet, dvA)
	f.SetColStyle(dataSheet, "E", customStyle)

	for i := 2; i < 10000; i++ {
		formula := fmt.Sprintf("=VLOOKUP(E%d, %s!A1:B%d, 2, FALSE)", i, audiencesSheet, len(audiences))
		f.SetCellFormula(dataSheet, fmt.Sprintf("J%d", i), formula)
	}
	f.SetSheetVisible(audiencesSheet, false)

	// Удаляем дефолтный лист
	f.DeleteSheet("Sheet1")

	// Устанавливаем активный лист
	f.SetActiveSheet(0) // данные

	buf := new(bytes.Buffer)
	if err := f.Write(buf); err != nil {
		return nil, fmt.Errorf("write excel: %w", err)
	}
	return buf.Bytes(), nil
}
