package services

import (
	"context"
	"core/internal/dto"
	"core/internal/models"
	"core/internal/repository"
	"errors"
	"fmt"
	"log"
	"slices"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ScheduleService struct {
	pool               *pgxpool.Pool
	audienceRepo       *repository.AudienceRepo
	groupRepo          *repository.GroupRepo
	teacherRepo        *repository.TeacherRepo
	studentRepo        *repository.StudentRepo
	userRepo           *repository.UserRepo
	roleRepo           *repository.RoleRepo
	subjectRepo        *repository.SubjectRepo
	academicPeriodRepo *repository.AcademicPeriodRepo
	scheduleRepo       *repository.ScheduleRepo
	lessonLogRepo      *repository.LessonLogRepo
}

func NewScheduleService(pool *pgxpool.Pool) *ScheduleService {
	return &ScheduleService{
		pool:               pool,
		audienceRepo:       repository.NewAudienceRepo(pool),
		groupRepo:          repository.NewGroupRepo(pool),
		teacherRepo:        repository.NewTeacherRepo(pool),
		studentRepo:        repository.NewStudentRepo(pool),
		userRepo:           repository.NewUserRepo(pool),
		roleRepo:           repository.NewRoleRepo(pool),
		subjectRepo:        repository.NewSubjectRepo(pool),
		academicPeriodRepo: repository.NewAcademicPeriodRepo(pool),
		scheduleRepo:       repository.NewScheduleRepo(pool),
		lessonLogRepo:      repository.NewLessonLogRepo(pool),
	}
}

// ---------- Audience methods ----------
func (s *ScheduleService) CreateAudience(ctx context.Context, req *dto.CreateAudienceRequest) (*models.Audience, error) {
	audience := models.Audience{
		Name:   req.Name,
		Number: req.Number,
	}
	err := s.audienceRepo.Create(ctx, &audience)
	return &audience, err
}

func (s *ScheduleService) GetAudience(ctx context.Context, id int) (*models.Audience, error) {
	return s.audienceRepo.Get(ctx, id)
}

func (s *ScheduleService) GetAllAudience(ctx context.Context) ([]models.Audience, error) {
	return s.audienceRepo.GetAll(ctx)
}

func (s *ScheduleService) UpdateAudience(ctx context.Context, id int, req *dto.UpdateAudienceRequest) (*models.Audience, error) {
	audience, err := s.GetAudience(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get audience: %w", err)
	}

	if req.Name != "" {
		audience.Name = req.Name
	}
	if req.Number != "" {
		audience.Number = req.Number
	}

	err = s.audienceRepo.Update(ctx, audience)
	return audience, err
}

func (s *ScheduleService) DeleteAudience(ctx context.Context, id int) error {
	return s.audienceRepo.Delete(ctx, id)
}

// ---------- Group methods ----------
func (s *ScheduleService) CreateGroup(ctx context.Context, req *dto.CreateGroupRequest) (*models.Group, error) {
	group := models.Group{
		Name:          req.Name,
		Specialty:     req.Specialty,
		AdmissionYear: req.AdmissionYear,
	}
	err := s.groupRepo.Create(ctx, &group)
	return &group, err
}

func (s *ScheduleService) GetGroup(ctx context.Context, id int) (*models.Group, error) {
	return s.groupRepo.GetByID(ctx, id)
}

func (s *ScheduleService) GetAllGroups(ctx context.Context) ([]models.Group, error) {
	return s.groupRepo.GetAll(ctx)
}

func (s *ScheduleService) UpdateGroup(ctx context.Context, id int, req *dto.UpdateGroupRequest) (*models.Group, error) {
	group, err := s.GetGroup(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get group: %w", err)
	}

	if req.Name != "" {
		group.Name = req.Name
	}

	if req.Specialty != "" {
		group.Specialty = req.Specialty
	}

	if req.AdmissionYear != 0 {
		group.AdmissionYear = req.AdmissionYear
	}

	err = s.groupRepo.Update(ctx, group)
	return group, err
}

func (s *ScheduleService) DeleteGroup(ctx context.Context, id int) error {
	return s.groupRepo.Delete(ctx, id)
}

func (s *ScheduleService) CreateTeacher(ctx context.Context, req *dto.CreateUserRequest) (*models.Teacher, error) {
	passwordHash, err := hashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}
	user := models.User{
		Name:         req.Username,
		FullName:     req.FullName,
		Email:        req.Email,
		PasswordHash: passwordHash,
	}
	if err := s.userRepo.Create(ctx, &user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	if err := s.teacherRepo.CreateTeacher(ctx, &user); err != nil {
		return nil, err
	}
	return &models.Teacher{User: user}, nil
}

func (s *ScheduleService) GetTeacher(ctx context.Context, userID int) (*models.Teacher, error) {
	return s.teacherRepo.GetTeacherByUserID(ctx, userID)
}

func (s *ScheduleService) GetAllTeachers(ctx context.Context) ([]models.Teacher, error) {
	return s.teacherRepo.GetAllTeachers(ctx)
}

func (s *ScheduleService) DeleteTeacher(ctx context.Context, userID int) error {
	return s.teacherRepo.DeleteTeacher(ctx, userID)
}

func (s *ScheduleService) CreateStudent(ctx context.Context, req *dto.CreateStudentRequest) (*models.Student, error) {
	group, err := s.groupRepo.GetByID(ctx, req.GroupID)
	if err != nil {
		return nil, fmt.Errorf("get group: %w", err)
	}
	if group == nil {
		return nil, fmt.Errorf("group with id %d not found", req.GroupID)
	}

	passwordHash, err := hashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}
	user := models.User{
		Name:         req.Username,
		FullName:     req.FullName,
		Email:        req.Email,
		PasswordHash: passwordHash,
	}
	if err := s.userRepo.Create(ctx, &user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	if err := s.studentRepo.CreateStudent(ctx, &user, req.GroupID); err != nil {
		return nil, err
	}
	return &models.Student{User: user, Group: group}, nil
}

func (s *ScheduleService) GetStudent(ctx context.Context, userID int) (*models.Student, error) {
	return s.studentRepo.GetStudentByUserID(ctx, userID)
}

func (s *ScheduleService) GetAllStudents(ctx context.Context) ([]models.Student, error) {
	return s.studentRepo.GetAllStudents(ctx)
}

func (s *ScheduleService) UpdateStudentGroup(ctx context.Context, userID int, groupID int) error {
	group, err := s.groupRepo.GetByID(ctx, groupID)
	if err != nil {
		return err
	}
	if group == nil {
		return fmt.Errorf("group with id %d not found", groupID)
	}
	return s.studentRepo.UpdateStudentGroup(ctx, userID, groupID)
}

func (s *ScheduleService) DeleteStudent(ctx context.Context, userID int) error {
	return s.studentRepo.DeleteStudent(ctx, userID)
}

func (s *ScheduleService) AssignStudent(ctx context.Context, userID int, groupID int) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}
	if user == nil {
		return errors.New("user not found")
	}

	isTeacher, err := s.teacherRepo.IsTeacher(ctx, userID)
	if err != nil {
		return fmt.Errorf("check teacher: %w", err)
	}
	if isTeacher {
		return errors.New("user is already a teacher, cannot assign student")
	}

	group, err := s.groupRepo.GetByID(ctx, groupID)
	if err != nil {
		return fmt.Errorf("get group: %w", err)
	}
	if group == nil {
		return errors.New("group not found")
	}

	if err := s.studentRepo.CreateStudent(ctx, user, groupID); err != nil {
		return fmt.Errorf("create student profile: %w", err)
	}

	return nil
}

func (s *ScheduleService) AssignTeacher(ctx context.Context, userID int) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}
	if user == nil {
		return errors.New("user not found")
	}

	isStudent, err := s.studentRepo.IsStudent(ctx, userID)
	if err != nil {
		return fmt.Errorf("check student: %w", err)
	}
	if isStudent {
		return errors.New("user is already student, cannot assign teacher")
	}

	if err := s.teacherRepo.CreateTeacher(ctx, user); err != nil {
		return fmt.Errorf("create teacher: %w", err)
	}

	return nil
}

func (s *ScheduleService) AssignRole(ctx context.Context, userID int, roleName string) error {
	if roleName != "admin" && roleName != "manager" {
		return errors.New("role must be 'admin' or 'manager'")
	}

	role, err := s.roleRepo.GetRoleByName(ctx, roleName)
	if err != nil {
		return fmt.Errorf("get role: %w", err)
	}
	if role == nil {
		return errors.New("role not found")
	}

	return s.roleRepo.AssignRoleToUser(ctx, userID, roleName)
}

func (s *ScheduleService) RemoveRole(ctx context.Context, userID int, roleName string) error {
	roles, err := s.roleRepo.GetRoles(ctx)
	if err != nil {
		return err
	}

	if !slices.ContainsFunc(roles, func(r models.Role) bool {
		return r.Name == roleName
	}) {
		return errors.New("invalid role name")
	}

	return s.roleRepo.RemoveRoleFromUser(ctx, userID, roleName)
}

func (s *ScheduleService) RemoveStudent(ctx context.Context, userID int) error {
	isStudent, err := s.studentRepo.IsStudent(ctx, userID)
	if err != nil {
		return fmt.Errorf("check student: %w", err)
	}
	if !isStudent {
		return errors.New("user is not a student")
	}

	if err := s.studentRepo.DeleteStudent(ctx, userID); err != nil {
		return fmt.Errorf("delete student profile: %w", err)
	}

	if err := s.roleRepo.RemoveRoleFromUser(ctx, userID, "student"); err != nil {
		return fmt.Errorf("remove student role: %w", err)
	}

	return nil
}

func (s *ScheduleService) RemoveTeacher(ctx context.Context, userID int) error {
	isTeacher, err := s.teacherRepo.IsTeacher(ctx, userID)
	if err != nil {
		return fmt.Errorf("check teacher: %w", err)
	}
	if !isTeacher {
		return errors.New("user is not a teacher")
	}

	if err := s.teacherRepo.DeleteTeacher(ctx, userID); err != nil {
		return fmt.Errorf("delete teacher profile: %w", err)
	}

	if err := s.roleRepo.RemoveRoleFromUser(ctx, userID, "teacher"); err != nil {
		return fmt.Errorf("remove student role: %w", err)
	}

	return nil
}

func (s *ScheduleService) GetUserRoles(ctx context.Context, userID int) ([]models.Role, error) {
	return s.roleRepo.GetUserRoles(ctx, userID)
}

// --------Subject Methods
func (s *ScheduleService) CreateSubject(ctx context.Context, req *dto.CreateSubjectRequest) (*models.Subject, error) {
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("parse start_date: %w", err)
	}
	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, fmt.Errorf("parse end_date: %w", err)
	}
	if startDate.After(endDate) {
		return nil, errors.New("start_date must be before end_date")
	}
	subject := models.Subject{
		Title:     req.Title,
		Semester:  req.Semester,
		HoursLoad: req.HoursLoad,
		StartDate: startDate,
		EndDate:   endDate,
		GroupID:   req.GroupID,
	}
	err = s.subjectRepo.CreateSubject(ctx, &subject)
	return &subject, err
}

func (s *ScheduleService) GetSubjectByID(ctx context.Context, id int) (*models.Subject, error) {
	return s.subjectRepo.GetSubjectByID(ctx, id)
}

func (s *ScheduleService) GetAllSubjects(ctx context.Context, req *dto.PaginatedSubjectsRequest) (*dto.PaginatedSubjectsResponse, error) {
	return s.subjectRepo.GetAllSubjects(ctx, req)
}

func (s *ScheduleService) GetSubjectsByGroupID(ctx context.Context, groupID int, req *dto.PaginatedSubjectsRequest) (*dto.PaginatedSubjectsResponse, error) {
	return s.subjectRepo.GetSubjectsByGroupID(ctx, groupID, req)
}

func (s *ScheduleService) UpdateSubject(ctx context.Context, id int, req *dto.UpdateSubjectRequest) (*models.Subject, error) {
	subject, err := s.GetSubjectByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get subject: %w", err)
	}

	if req.Title != "" {
		subject.Title = req.Title
	}

	if req.Semester != 0 {
		subject.Semester = req.Semester
	}

	if req.HoursLoad != 0 {
		subject.HoursLoad = req.HoursLoad
	}

	if req.StartDate != "" {
		startDate, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			return nil, fmt.Errorf("parse start_date: %w", err)
		}
		subject.StartDate = startDate
	}

	if req.EndDate != "" {
		endDate, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			return nil, fmt.Errorf("parse end_date: %w", err)
		}
		subject.EndDate = endDate
	}

	if subject.StartDate.After(subject.EndDate) {
		return nil, errors.New("start_date must be before end_date")
	}

	err = s.subjectRepo.UpdateSubject(ctx, subject)
	return subject, err
}

func (s *ScheduleService) DeleteSubject(ctx context.Context, id int) error {
	return s.subjectRepo.DeleteSubject(ctx, id)
}

func (s *ScheduleService) CreateAcademicPeriod(ctx context.Context, req *dto.CreateAcademicPeriodRequest) (*models.AcademicPeriod, error) {
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("parse start_date: %w", err)
	}
	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, fmt.Errorf("parse end_date: %w", err)
	}
	if startDate.After(endDate) {
		return nil, errors.New("start_date must be before end_date")
	}
	acdemicPeriod := models.AcademicPeriod{
		Year:      req.Year,
		Semester:  req.Semester,
		StartDate: startDate,
		EndDate:   endDate,
	}

	err = s.academicPeriodRepo.Create(ctx, &acdemicPeriod)
	return &acdemicPeriod, err
}

func (s *ScheduleService) GetAcademicPeriodByID(ctx context.Context, id int) (*models.AcademicPeriod, error) {
	return s.academicPeriodRepo.GetByID(ctx, id)
}

func (s *ScheduleService) GetAllAcademicPeriods(ctx context.Context) ([]models.AcademicPeriod, error) {
	return s.academicPeriodRepo.GetAll(ctx)
}

func (s *ScheduleService) UpdateAcademicPeriod(ctx context.Context, id int, req *dto.UpdateAcademicPeriodRequest) (*models.AcademicPeriod, error) {
	academicPeriod, err := s.GetAcademicPeriodByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get academic period: %w", err)
	}

	if req.Year != "" {
		academicPeriod.Year = req.Year
	}

	if req.Semester != 0 {
		academicPeriod.Semester = req.Semester
	}

	if req.StartDate != "" {
		startDate, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			return nil, fmt.Errorf("parse start_date: %w", err)
		}
		academicPeriod.StartDate = startDate
	}

	if req.EndDate != "" {
		endDate, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			return nil, fmt.Errorf("parse end_date: %w", err)
		}
		academicPeriod.EndDate = endDate
	}

	if academicPeriod.StartDate.After(academicPeriod.StartDate) {
		return nil, errors.New("start_date must be before end_date")
	}

	err = s.academicPeriodRepo.Update(ctx, academicPeriod)
	return academicPeriod, err
}

func (s *ScheduleService) DeleteAcademicPeriod(ctx context.Context, id int) error {
	return s.academicPeriodRepo.Delete(ctx, id)
}

// ------- Schedule Templates Methods --------
func (s *ScheduleService) CreateScheduleTemplate(ctx context.Context, req *dto.CreateScheduleTemplateRequest) (*models.ScheduleTemplate, error) {
	subject, err := s.GetSubjectByID(ctx, req.SubjectID)
	if err != nil {
		return nil, fmt.Errorf("get subject with id %d: %w", req.SubjectID, err)
	}
	if subject == nil {
		return nil, errors.New("subject not found")
	}

	teacher, err := s.GetTeacher(ctx, req.TeacherID)
	if err != nil {
		return nil, fmt.Errorf("get teacher with id %d: %w", req.TeacherID, err)
	}
	if teacher == nil {
		return nil, errors.New("teacher not found")
	}

	audience, err := s.GetAudience(ctx, req.AudienceID)
	if err != nil {
		return nil, fmt.Errorf("get audience with id %d: %w", req.AudienceID, err)
	}
	if audience == nil {
		return nil, errors.New("audience not found")
	}

	period, err := s.GetAcademicPeriodByID(ctx, req.AcademicPeriodID)
	if err != nil {
		return nil, fmt.Errorf("get academic period with id %d: %w", req.AcademicPeriodID, err)
	}
	if period == nil {
		return nil, errors.New("academic period not found")
	}

	template := models.ScheduleTemplate{
		DayOfWeek:        req.DayOfWeek,
		WeekType:         req.WeekType,
		Number:           req.Number,
		SubjectID:        req.SubjectID,
		TeacherID:        req.TeacherID,
		AudienceID:       req.AudienceID,
		AcademicPeriodID: req.AcademicPeriodID,
		Subject:          *subject,
	}

	if err := s.checkScheduleTemplateConflicts(ctx, &template); err != nil {
		return nil, err
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	err = s.scheduleRepo.CreateWithTx(ctx, tx, &template)
	if err != nil {
		return nil, fmt.Errorf("create schedule template: %w", err)
	}
	err = s.updatePlannedSubjectLessons(ctx, tx, &template.Subject, template.AcademicPeriodID)
	if err != nil {
		return nil, fmt.Errorf("update planned lessons for subject %d: %w", subject.ID, err)
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	return s.scheduleRepo.GetByID(ctx, template.ID)
}

func (s *ScheduleService) GetScheduleTemplate(ctx context.Context, id int) (*models.ScheduleTemplate, error) {
	return s.scheduleRepo.GetByID(ctx, id)
}

func (s *ScheduleService) GetAllScheduleTemplates(ctx context.Context, filters *dto.ScheduleFiltersRequest) ([]models.ScheduleTemplate, error) {
	repoFilters := repository.ScheduleFilters{
		GroupID:          filters.GroupID,
		SubjectID:        filters.SubjectID,
		TeacherID:        filters.TeacherID,
		AudienceID:       filters.AudienceID,
		DayOfWeek:        filters.DayOfWeek,
		WeekType:         filters.WeekType,
		AcademicPeriodID: filters.AcademicPeriodID,
	}

	return s.scheduleRepo.GetAll(ctx, repoFilters)
}

func (s *ScheduleService) UpdateScheduleTemplate(ctx context.Context, id int, req *dto.UpdateScheduleTemplateRequest) (*models.ScheduleTemplate, error) {
	template, err := s.GetScheduleTemplate(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get schedule template: %w", err)
	}
	if template == nil {
		return nil, errors.New("schedule template no found")
	}

	if req.DayOfWeek != nil {
		template.DayOfWeek = *req.DayOfWeek
	}
	if req.Number != nil {
		template.Number = *req.Number
	}
	if req.WeekType != nil {
		template.WeekType = *req.WeekType
	}
	if req.SubjectID != nil {
		subject, err := s.subjectRepo.GetSubjectByID(ctx, *req.SubjectID)
		if err != nil {
			return nil, fmt.Errorf("get subject: %w", err)
		}
		if subject == nil {
			return nil, errors.New("subject not found")
		}
		template.SubjectID = *req.SubjectID
	}
	if req.TeacherID != nil {
		teacher, err := s.teacherRepo.GetTeacherByUserID(ctx, *req.TeacherID)
		if err != nil {
			return nil, fmt.Errorf("get teacher: %w", err)
		}
		if teacher == nil {
			return nil, errors.New("teacher not found")
		}
		template.TeacherID = *req.TeacherID
	}
	if req.AudienceID != nil {
		audience, err := s.audienceRepo.Get(ctx, *req.AudienceID)
		if err != nil {
			return nil, fmt.Errorf("get audience: %w", err)
		}
		if audience == nil {
			return nil, errors.New("audience not found")
		}
		template.AudienceID = *req.AudienceID
	}
	if req.AcademicPeriodID != nil {
		period, err := s.academicPeriodRepo.GetByID(ctx, *req.AcademicPeriodID)
		if err != nil {
			return nil, fmt.Errorf("get academic period: %w", err)
		}
		if period == nil {
			return nil, errors.New("academic period not found")
		}
		template.AcademicPeriodID = *req.AcademicPeriodID
	}

	if err := s.checkScheduleTemplateConflicts(ctx, template); err != nil {
		return nil, err
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	err = s.scheduleRepo.UpdateWithTx(ctx, tx, template)
	if err != nil {
		return nil, fmt.Errorf("create academic period: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	return s.scheduleRepo.GetByID(ctx, template.ID)
}

func (s *ScheduleService) DeleteScheduleTemplate(ctx context.Context, id int) error {
	return s.scheduleRepo.Delete(ctx, id)
}

func (s *ScheduleService) GetGroupSchedule(ctx context.Context, groupID int, periodID *int) ([]models.ScheduleTemplate, error) {
	if periodID == nil {
		period, err := s.academicPeriodRepo.GetActive(ctx)
		if err != nil {
			return nil, fmt.Errorf("get active academic period: %w", err)
		}

		periodID = &period.ID
	}

	filters := repository.ScheduleFilters{
		GroupID:          &groupID,
		AcademicPeriodID: periodID,
	}

	return s.scheduleRepo.GetAll(ctx, filters)
}

func (s *ScheduleService) GetTeacherSchedule(ctx context.Context, teacherID int, periodID *int) ([]models.ScheduleTemplate, error) {
	if periodID == nil {
		period, err := s.academicPeriodRepo.GetActive(ctx)
		if err != nil {
			return nil, fmt.Errorf("get active academic period: %w", err)
		}
		periodID = &period.ID
	}

	filters := repository.ScheduleFilters{
		TeacherID:        &teacherID,
		AcademicPeriodID: periodID,
	}

	return s.scheduleRepo.GetAll(ctx, filters)
}

func (s *ScheduleService) GetAudienceSchedule(ctx context.Context, audienceID int, periodID *int) ([]models.ScheduleTemplate, error) {
	if periodID == nil {
		period, err := s.academicPeriodRepo.GetActive(ctx)
		if err != nil {
			return nil, fmt.Errorf("get active academic period: %w", err)
		}
		periodID = &period.ID
	}

	filters := repository.ScheduleFilters{
		AudienceID:       &audienceID,
		AcademicPeriodID: periodID,
	}

	return s.scheduleRepo.GetAll(ctx, filters)
}

func (s *ScheduleService) CreateSemesterSchedule(ctx context.Context, req *dto.CreateSemesterScheduleRequest) error {
	period, err := s.GetAcademicPeriodByID(ctx, req.AcademicPeriodID)
	if err != nil {
		return fmt.Errorf("get academic period: %w", err)
	}
	if period == nil {
		return errors.New("academic period not found")
	}

	group, err := s.GetGroup(ctx, req.GroupID)
	if err != nil {
		return fmt.Errorf("get group: %w", err)
	}
	if group == nil {
		return errors.New("group not found")
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Проверяем, нет ли уже расписания у группы на заданный период
	existingsFilters := repository.ScheduleFilters{
		GroupID:          &req.GroupID,
		AcademicPeriodID: &req.AcademicPeriodID,
	}

	existingsTemplstes, err := s.scheduleRepo.GetAll(ctx, existingsFilters)
	if err != nil {
		return fmt.Errorf("check existings schedule: %w", err)
	}
	if len(existingsTemplstes) > 0 {
		return errors.New("schedule already exists for this group and academic period")
	}

	var subjects []models.Subject

	for _, entry := range req.ScheduleEntries {

		subject, err := s.GetSubjectByID(ctx, entry.SubjectID)
		if err != nil {
			return fmt.Errorf("get subject %d: %w", entry.SubjectID, err)
		}
		if subject == nil {
			return fmt.Errorf("subject with id %d not found", entry.SubjectID)
		}

		if !slices.ContainsFunc(subjects, func(s models.Subject) bool {
			return s.ID == entry.SubjectID
		}) {
			subjects = append(subjects, *subject)
		}

		if subject.GroupID != req.GroupID {
			return fmt.Errorf("subject %d does not belong to group %d", entry.SubjectID, req.GroupID)
		}

		template := models.ScheduleTemplate{
			DayOfWeek:        entry.DayOfWeek,
			WeekType:         entry.WeekType,
			Number:           entry.Number,
			SubjectID:        entry.SubjectID,
			TeacherID:        entry.TeacherID,
			AudienceID:       entry.AudienceID,
			AcademicPeriodID: req.AcademicPeriodID,
			Subject:          *subject,
		}

		if err := s.checkScheduleTemplateConflicts(ctx, &template); err != nil {
			return err
		}

		err = s.scheduleRepo.CreateWithTx(ctx, tx, &template)
		if err != nil {
			return fmt.Errorf("create schedule template: %w", err)
		}
	}

	for _, subject := range subjects {
		err = s.updatePlannedSubjectLessons(ctx, tx, &subject, req.AcademicPeriodID)
		if err != nil {
			return fmt.Errorf("update planned: %w", err)
		}
	}

	return tx.Commit(ctx)
}

func (s *ScheduleService) checkScheduleTemplateConflicts(ctx context.Context, tmpl *models.ScheduleTemplate) error {
	// check teachers
	templates, err := s.scheduleRepo.GetAll(ctx, repository.ScheduleFilters{
		TeacherID:        &tmpl.TeacherID,
		AcademicPeriodID: &tmpl.AcademicPeriodID,
		DayOfWeek:        &tmpl.DayOfWeek,
	})
	if err != nil {
		return fmt.Errorf("check teacher conflicts: %w", err)
	}
	for _, t := range templates {
		if t.Number == tmpl.Number {
			if t.WeekType == 0 || t.WeekType == tmpl.WeekType {
				return fmt.Errorf("teacher %d already has a lesson at day %d, lesson %d, week type %d",
					tmpl.TeacherID, tmpl.DayOfWeek, tmpl.Number, tmpl.WeekType)
			}
		}
	}

	// check audiences
	templates, err = s.scheduleRepo.GetAll(ctx, repository.ScheduleFilters{
		AudienceID:       &tmpl.AudienceID,
		AcademicPeriodID: &tmpl.AcademicPeriodID,
		DayOfWeek:        &tmpl.DayOfWeek,
	})
	if err != nil {
		return fmt.Errorf("check audience conflicts: %w", err)
	}
	for _, t := range templates {
		if t.Number == tmpl.Number {
			if t.WeekType == 0 || t.WeekType == tmpl.WeekType {
				return fmt.Errorf("audience %d already has a lesson at day %d, lesson %d week type %d",
					tmpl.TeacherID, tmpl.DayOfWeek, tmpl.Number, tmpl.WeekType)
			}
		}
	}

	// check group
	templates, err = s.scheduleRepo.GetAll(ctx, repository.ScheduleFilters{
		GroupID:          &tmpl.Subject.GroupID,
		AcademicPeriodID: &tmpl.AcademicPeriodID,
		DayOfWeek:        &tmpl.DayOfWeek,
	})
	if err != nil {
		return fmt.Errorf("check group conflicts: %w", err)
	}
	for _, t := range templates {
		if t.Number == tmpl.Number {
			if t.WeekType == 0 || t.WeekType == tmpl.WeekType {
				return fmt.Errorf("group %d already has a lesson at day %d, lesson %d, week type %d",
					tmpl.TeacherID, tmpl.DayOfWeek, tmpl.Number, tmpl.WeekType)
			}
		}
	}

	return nil
}

func (s *ScheduleService) updatePlannedSubjectLessons(ctx context.Context, tx pgx.Tx, subject *models.Subject, periodID int) error {
	templates, err := s.scheduleRepo.GetAll(ctx, repository.ScheduleFilters{
		SubjectID: &subject.ID,
	})
	if err != nil {
		return fmt.Errorf("get schedule template for subject %d: %w", subject.ID, err)
	}
	if len(templates) == 0 {
		return s.deleteAllPlannedSubjectLessons(ctx, tx, subject.ID)
	}

	currentLogs, err := s.lessonLogRepo.GetAll(ctx, repository.LessonLogFilters{
		SubjectID: &subject.ID,
	})
	if err != nil {
		return fmt.Errorf("get current lesson logs with subject %d: %w", subject.ID, err)
	}

	expectedLessons := s.generateExpectedSubjectLessons(templates, currentLogs, subject, periodID)
	var plannedLessons []models.LessonLog
	for _, curr := range currentLogs {
		if curr.Status == models.LessonStatusPlanned {
			plannedLessons = append(plannedLessons, curr)
		}
	}

	log.Printf("expectedLessons %+v", expectedLessons)
	log.Printf("currentLessons %+v", currentLogs)

	toDelete, toUpdate, toCreate := s.diffPlannedLessons(plannedLessons, expectedLessons)

	log.Printf("toDelete %+v", toDelete)
	log.Printf("toUpdate %+v", toUpdate)
	log.Printf("toCreate %+v", toCreate)

	for _, lesson := range toDelete {
		if err := s.lessonLogRepo.DeleteWithTx(ctx, tx, lesson.ID); err != nil {
			return fmt.Errorf("delete lesson %d: %w", lesson.ID, err)
		}
	}

	for _, lesson := range toUpdate {
		if err := s.lessonLogRepo.UpdateWithTx(ctx, tx, &lesson); err != nil {
			return fmt.Errorf("update lesson %d: %w", lesson.ID, err)
		}
	}

	for _, lesson := range toCreate {
		if err := s.lessonLogRepo.CreateWithTx(ctx, tx, &lesson); err != nil {
			return fmt.Errorf("create lesson %d: %w", lesson.ID, err)
		}
	}

	return nil
}

func (s *ScheduleService) deleteAllPlannedSubjectLessons(ctx context.Context, tx pgx.Tx, subjectID int) error {
	plannedStatus := models.LessonStatusPlanned
	lessons, err := s.lessonLogRepo.GetAll(ctx, repository.LessonLogFilters{
		SubjectID: &subjectID,
		Status:    &plannedStatus,
	})
	if err != nil {
		return fmt.Errorf("get planned lessons for deletion: %w", err)
	}

	for _, lesson := range lessons {
		if err := s.lessonLogRepo.DeleteWithTx(ctx, tx, lesson.ID); err != nil {
			return fmt.Errorf("delete lesson %d: %w", lesson.ID, err)
		}
	}
	log.Printf("Deleted %d planned lessons for subject %d (no templates)", len(lessons), subjectID)
	return nil
}

func (s *ScheduleService) generateExpectedSubjectLessons(templates []models.ScheduleTemplate, current []models.LessonLog, subject *models.Subject, periodID int) []models.LessonLog {
	var expected []models.LessonLog

	lessonCount := subject.HoursLoad / 2
	if subject.HoursLoad%2 == 1 {
		lessonCount++
	}

	scheduledLessons := 0
	for _, curr := range current {
		if curr.Status != models.LessonStatusPlanned && curr.Status != models.LessonStatusCanceled {
			scheduledLessons++
		}
	}

	if scheduledLessons >= lessonCount {
		return expected
	}

	for currentDate := subject.StartDate; !currentDate.After(subject.EndDate); currentDate = currentDate.AddDate(0, 0, 1) {
		dayOfWeek := int(currentDate.Weekday())
		if dayOfWeek == 0 {
			dayOfWeek = 7
		}

		weekType := 1
		_, week := currentDate.ISOWeek()
		if week%2 != 1 {
			weekType = 2
		}

		for _, tmpl := range templates {
			if tmpl.DayOfWeek != dayOfWeek {
				continue
			}

			if tmpl.WeekType != 0 {
				if tmpl.WeekType != weekType {
					continue
				}
			}

			if slices.ContainsFunc(current, func(l models.LessonLog) bool {
				return l.Date.Equal(currentDate) &&
					l.Number == tmpl.Number &&
					l.Status != models.LessonStatusPlanned
			}) {
				continue
			}

			lesson := models.LessonLog{
				Date:             currentDate,
				Number:           tmpl.Number,
				Status:           models.LessonStatusPlanned,
				Comment:          "Сгенерировано автоматически",
				SubjectID:        tmpl.SubjectID,
				TeacherID:        tmpl.TeacherID,
				AudienceID:       tmpl.AudienceID,
				AcademicPeriodID: periodID,
				IsFromTemplate:   true,
			}
			expected = append(expected, lesson)
			scheduledLessons++
			if scheduledLessons >= lessonCount {
				return expected
			}
		}
	}
	return expected
}

func (s *ScheduleService) diffPlannedLessons(current, expected []models.LessonLog) (toDelete, toUpdate, toCreate []models.LessonLog) {
	expectedMap := make(map[string]models.LessonLog)
	for _, exp := range expected {
		key := fmt.Sprintf("%s_%d", exp.Date.Format("2006-01-02"), exp.Number)
		expectedMap[key] = exp
	}

	for _, cur := range current {
		key := fmt.Sprintf("%s_%d", cur.Date.Format("2006-01-02"), cur.Number)
		if exp, ok := expectedMap[key]; ok {
			if cur.SubjectID != exp.SubjectID ||
				cur.TeacherID != exp.TeacherID ||
				cur.AudienceID != exp.AudienceID {
				cur.SubjectID = exp.SubjectID
				cur.TeacherID = exp.TeacherID
				cur.AudienceID = exp.AudienceID
				cur.Comment = exp.Comment
				toUpdate = append(toUpdate, cur)
			}

			delete(expectedMap, key)
		} else {
			toDelete = append(toDelete, cur)
		}
	}

	for _, exp := range expectedMap {
		toCreate = append(toCreate, exp)
	}

	return
}

func (s *ScheduleService) GetLessonLogs(ctx context.Context, filters *dto.LessonLogFiltersRequest) ([]models.LessonLog, error) {
	repoFilters := repository.LessonLogFilters{}

	if filters.GroupID != nil {
		repoFilters.GroupID = filters.GroupID
	}
	if filters.SubjectID != nil {
		repoFilters.SubjectID = filters.SubjectID
	}
	if filters.TeacherID != nil {
		repoFilters.TeacherID = filters.TeacherID
	}
	if filters.AudienceID != nil {
		repoFilters.AudienceID = filters.AudienceID
	}
	if filters.AcademicPeriodID != nil {
		repoFilters.AcademicPeriodID = filters.AcademicPeriodID
	}
	if filters.DateFrom != nil {
		repoFilters.DateFrom = filters.DateFrom
	}
	if filters.DateTo != nil {
		repoFilters.DateTo = filters.DateTo
	}
	if filters.Status != nil {
		repoFilters.Status = filters.Status
	}
	if filters.Number != nil {
		repoFilters.Number = filters.Number
	}

	return s.lessonLogRepo.GetAll(ctx, repoFilters)
}
