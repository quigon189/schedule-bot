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

func (s *ExcelService) GenerateGroupTemplate(req *dto.CreateGroupWithCurriculumRequest) ([]byte, error) {
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
	f.SetColWidth(sheetInfo, "A", "A", 20)
	f.SetColWidth(sheetInfo, "B", "B", 30)
	f.SetColWidth(sheetInfo, "C", "C", 15)

	f.SetCellStr(sheetInfo, "A2", req.Group.Name)
	f.SetCellStr(sheetInfo, "B2", req.Group.Specialty)
	f.SetCellInt(sheetInfo, "C2", int64(req.Group.AdmissionYear))

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
	f.SetColWidth(sheetSubjects, "A", "E", 15)

	if len(req.Subjects) > 0 {
		for i, subj := range req.Subjects {
			f.SetCellStr(sheetSubjects, fmt.Sprintf("A%d", i+2), subj.Title)
			f.SetCellInt(sheetSubjects, fmt.Sprintf("B%d", i+2), int64(subj.Semester))
			f.SetCellInt(sheetSubjects, fmt.Sprintf("C%d", i+2), int64(subj.HoursLoad))
			f.SetCellStr(sheetSubjects, fmt.Sprintf("D%d", i+2), subj.StartDate)
			f.SetCellStr(sheetSubjects, fmt.Sprintf("E%d", i+2), subj.EndDate)
		}
	} else {
		// Пример
		f.SetCellStr(sheetSubjects, "A2", "Математика")
		f.SetCellInt(sheetSubjects, "B2", 1)
		f.SetCellInt(sheetSubjects, "C2", 72)
		f.SetCellStr(sheetSubjects, "D2", "2025-09-01")
		f.SetCellStr(sheetSubjects, "E2", "2025-12-31")
	}

	// Лист 3: Студенты
	sheetStudents := "Студенты"
	_, err = f.NewSheet(sheetStudents)
	if err != nil {
		return nil, fmt.Errorf("create sheet students: %w", err)
	}

	f.SetColWidth(sheetStudents, "A", "B", 25)
	f.SetCellStr(sheetStudents, "A1", "full_name")
	f.SetCellStr(sheetStudents, "B1", "email")

	if len(req.Students) > 0 {
		for i, stud := range req.Students {
			f.SetCellStr(sheetStudents, fmt.Sprintf("A%d", i+2), stud.FullName)
			f.SetCellStr(sheetStudents, fmt.Sprintf("B%d", i+2), stud.Email)
		}
	} else {
		f.SetCellStr(sheetStudents, "A2", "Иванов Иван Иванович")
		f.SetCellStr(sheetStudents, "B2", "ivanov@example.com")
	}

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

// ParsePlannerTemplate читает заполненный Excel и возвращает запрос для планировщика
func (s *ExcelService) ParsePlannerTemplate(r io.Reader, scheduleService *ScheduleService) (*dto.PlanScheduleRequest, error) {
	f, err := excelize.OpenReader(r)
	if err != nil {
		return nil, fmt.Errorf("open excel: %w", err)
	}
	defer f.Close()

	// Читаем параметры
	paramsSheet := "Параметры"
	academicPeriodID, err := s.getCellInt(f, paramsSheet, "A2")
	if err != nil {
		return nil, fmt.Errorf("academic_period_id: %w", err)
	}
	days, err := s.getCellInt(f, paramsSheet, "C2")
	if err != nil {
		days = 5
	}
	slots, err := s.getCellInt(f, paramsSheet, "D2")
	if err != nil {
		slots = 5
	}
	seed, err := s.getCellInt(f, paramsSheet, "E2")
	if err != nil {
		seed = 0
	}

	// Читаем данные
	dataSheet := "Данные"
	rows, err := f.GetRows(dataSheet)
	if err != nil {
		return nil, fmt.Errorf("get rows from data sheet: %w", err)
	}
	if len(rows) < 2 {
		return nil, fmt.Errorf("data sheet must have at least one row")
	}

	var subjectRequests []dto.SubjectRequest
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		if len(row) < 10 {
			return nil, fmt.Errorf("parse %d row: must be filled 10 columns", i+1)
		}
		priority, _ := strconv.Atoi(strings.TrimSpace(row[5]))
		lessonsCount, _ := strconv.Atoi(strings.TrimSpace(row[6]))
		subjectID, err := s.getCellInt(f, dataSheet, fmt.Sprintf("H%d", i+1))
		if err != nil {
			return nil, fmt.Errorf("parse subject id in row %d: %w", i+1, err)
		}
		if subjectID == 0 {
			break
		}
		teacherID, err := s.getCellInt(f, dataSheet, fmt.Sprintf("I%d", i+1))
		if err != nil {
			return nil, fmt.Errorf("parse teacher id in row %d: %w", i+1, err)
		}
		audienceID, err := s.getCellInt(f, dataSheet, fmt.Sprintf("J%d", i+1))
		if err != nil {
			return nil, fmt.Errorf("parse audience id in row %d: %w", i+1, err)
		}
		subjectRequests = append(subjectRequests, dto.SubjectRequest{
			SubjectID:    subjectID,
			TeacherID:    teacherID,
			AudienceID:   audienceID,
			Priority:     priority,
			LessonsCount: lessonsCount,
		})
	}

	return &dto.PlanScheduleRequest{
		AcademicPeriodID: academicPeriodID,
		Days:             days,
		Slots:            slots,
		Seed:             seed,
		SubjectList:      subjectRequests,
	}, nil
}

// Вспомогательные методы для парсинга Excel
func (s *ExcelService) getCellInt(f *excelize.File, sheet, cell string) (int, error) {
	val, err := f.GetCellValue(sheet, cell)
	if err != nil {
		return 0, err
	}
	if val == "" {
		return 0, nil
	}
	return strconv.Atoi(strings.TrimSpace(val))
}

func (s *ExcelService) GenerateGroupSchedule(sch *dto.WeeklySchedule) ([]byte, error) {
	f := excelize.NewFile()
	defer f.Close()

	weekDays := map[int]string{
		1: "Понедельник",
		2: "Вторник",
		3: "Среда",
		4: "Четверг",
		5: "Пятница",
		6: "Суббота",
	}

	for _, group := range sch.Groups {
		groupSheet := group.Name
		_, err := f.NewSheet(groupSheet)
		if err != nil {
			return nil, fmt.Errorf("create params sheet: %w", err)
		}

		f.SetCellStr(groupSheet, "A1", "День недели")
		f.SetCellStr(groupSheet, "B1", "Номер пары")
		f.SetCellStr(groupSheet, "C1", "Предмет")
		f.SetCellStr(groupSheet, "D1", "Аудитория")
		f.SetCellStr(groupSheet, "E1", "Преподаватель")

		f.SetColWidth(groupSheet, "A", "E", 20)

		type SchRow struct {
			WeekDay  string
			Slot     int
			Subject  string
			Audience string
			Teacher  string
		}
		schRows := []SchRow{}
		for day, daySchedule := range sch.Grid {
			for slot, data := range daySchedule {
				for _, lesson := range data {
					if lesson.Subject.GroupID == group.ID {
						var subjectTitle string
						if lesson.WeekType == 1 {
							subjectTitle += "нечетная "
						}
						if lesson.WeekType == 2 {
							subjectTitle += "четная"
						}
						schRows = append(schRows, SchRow{
							WeekDay:  weekDays[day],
							Slot:     slot,
							Subject:  subjectTitle + lesson.Subject.Title,
							Audience: lesson.Audience.Number,
							Teacher:  lesson.Teacher.User.GetShortName(),
						})
					}
				}
			}
		}

		for i, row := range schRows {
			f.SetCellStr(groupSheet, fmt.Sprintf("A%d", i+2), row.WeekDay)
			f.SetCellInt(groupSheet, fmt.Sprintf("B%d", i+2), int64(row.Slot))
			f.SetCellStr(groupSheet, fmt.Sprintf("C%d", i+2), row.Subject)
			f.SetCellStr(groupSheet, fmt.Sprintf("D%d", i+2), row.Audience)
			f.SetCellStr(groupSheet, fmt.Sprintf("E%d", i+2), row.Teacher)
		}
	}

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
