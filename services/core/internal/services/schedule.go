package services

import (
	"context"
	"core/internal/dto"
	"core/internal/models"
	"core/internal/repository"
	"errors"
	"fmt"
	"io"
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

func (s *ScheduleService) GetAudienceByNumber(ctx context.Context, number string) (*models.Audience, error) {
	return s.audienceRepo.GetByNumber(ctx, number)
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

func (s *ScheduleService) GetGroupByName(ctx context.Context, name string) (*models.Group, error) {
	return s.groupRepo.GetByName(ctx, name)
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

func (s *ScheduleService) GetTeacherByName(ctx context.Context, name string) (*models.Teacher, error) {
	return s.teacherRepo.GetTeacherByName(ctx, name)
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

func (s *ScheduleService) GetStudentsByGroupID(ctx context.Context, id int) ([]models.Student, error) {
	return s.studentRepo.GetStudentsByGroupID(ctx, id)
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

func (s *ScheduleService) GetSubjectsByDate(ctx context.Context, startDate, endDate time.Time) ([]models.Subject, error) {
	return s.subjectRepo.GetSubjectsByDate(ctx, startDate, endDate)
}

func (s *ScheduleService) GetSubjectsByGroupID(ctx context.Context, groupID int) ([]models.Subject, error) {
	return s.subjectRepo.GetSubjectsByGroupID(ctx, groupID)
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

func (s *ScheduleService) GetActiveAcademicPeriod(ctx context.Context) (*models.AcademicPeriod, error) {
	return s.academicPeriodRepo.GetActive(ctx)
}

func (s *ScheduleService) GetAcademicPeriodByYear(ctx context.Context, year string, semester int) (*models.AcademicPeriod, error) {
	return s.academicPeriodRepo.GetByYear(ctx, year, semester)
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

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := s.checkScheduleTemplateConflicts(ctx, tx, &template); err != nil {
		return nil, err
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

	if repoFilters.AcademicPeriodID == nil {
		period, err := s.academicPeriodRepo.GetActive(ctx)
		if err != nil {
			return nil, fmt.Errorf("get active academic period: %w", err)
		}
		repoFilters.AcademicPeriodID = &period.ID
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

	currentSubject, err := s.subjectRepo.GetSubjectByID(ctx, template.SubjectID)
	if err != nil {
		return nil, fmt.Errorf("get subject: %w", err)
	}
	if currentSubject == nil {
		return nil, errors.New("subject not found")
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

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := s.scheduleRepo.DeleteWithTx(ctx, tx, template.ID); err != nil {
		return nil, fmt.Errorf("delete old template: %w", err)
	}

	if err := s.checkScheduleTemplateConflicts(ctx, tx, template); err != nil {
		return nil, err
	}

	err = s.scheduleRepo.CreateWithTx(ctx, tx, template)
	if err != nil {
		return nil, fmt.Errorf("create academic period: %w", err)
	}

	err = s.updatePlannedSubjectLessons(ctx, tx, currentSubject, template.AcademicPeriodID)
	if err != nil {
		return nil, fmt.Errorf("update planned lessons: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	return s.scheduleRepo.GetByID(ctx, template.ID)
}

func (s *ScheduleService) DeleteScheduleTemplate(ctx context.Context, id int) error {
	tmpl, err := s.GetScheduleTemplate(ctx, id)
	if err != nil {
		return err
	}
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	err = s.scheduleRepo.DeleteWithTx(ctx, tx, id)
	if err != nil {
		return err
	}
	err = s.updatePlannedSubjectLessons(ctx, tx, &tmpl.Subject, tmpl.AcademicPeriodID)
	if err != nil {
		return fmt.Errorf("update planned lessons: %w", err)
	}

	return tx.Commit(ctx)
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

		if err := s.checkScheduleTemplateConflicts(ctx, tx, &template); err != nil {
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

func (s *ScheduleService) checkScheduleTemplateConflicts(ctx context.Context, tx pgx.Tx, tmpl *models.ScheduleTemplate) error {
	// check teachers
	templates, err := s.scheduleRepo.GetAllWithTx(ctx, tx, repository.ScheduleFilters{
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
	templates, err = s.scheduleRepo.GetAllWithTx(ctx, tx, repository.ScheduleFilters{
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
	templates, err = s.scheduleRepo.GetAllWithTx(ctx, tx, repository.ScheduleFilters{
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
	templates, err := s.scheduleRepo.GetAllWithTx(ctx, tx, repository.ScheduleFilters{
		SubjectID: &subject.ID,
	})
	if err != nil {
		return fmt.Errorf("get schedule template for subject %d: %w", subject.ID, err)
	}
	if len(templates) == 0 {
		return s.deleteAllPlannedSubjectLessons(ctx, tx, subject.ID)
	}

	currentLogs, err := s.lessonLogRepo.GetAllWithTx(ctx, tx, repository.LessonLogFilters{
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

	toDelete, toUpdate, toCreate := s.diffPlannedLessons(plannedLessons, expectedLessons)

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

	log.Printf("Updated planned lessons logs for subject %d: deleted %d, updated %d, created %d",
		subject.ID, len(toDelete), len(toUpdate), len(toCreate))

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
	repoFilters := repository.LessonLogFilters{
		GroupID:          filters.GroupID,
		SubjectID:        filters.SubjectID,
		TeacherID:        filters.TeacherID,
		AudienceID:       filters.AudienceID,
		AcademicPeriodID: filters.AcademicPeriodID,
		DateFrom:         filters.DateFrom,
		DateTo:           filters.DateTo,
		Status:           filters.Status,
		Number:           filters.Number,
	}

	if repoFilters.AcademicPeriodID == nil {
		period, err := s.academicPeriodRepo.GetActive(ctx)
		if err != nil {
			return nil, fmt.Errorf("get active academic period: %w", err)
		}

		repoFilters.AcademicPeriodID = &period.ID
	}

	return s.lessonLogRepo.GetAll(ctx, repoFilters)
}

func (s *ScheduleService) CancelLesson(ctx context.Context, logID int, req *dto.CancelLessonRequest) error {
	lesson, err := s.lessonLogRepo.GetByID(ctx, logID)
	if err != nil {
		return fmt.Errorf("get lesson log: %w", err)
	}

	if lesson.Status == models.LessonStatusCanceled {
		return fmt.Errorf("lesson already is calceled")
	}

	lesson.Comment = req.Comment
	lesson.Status = models.LessonStatusCanceled

	subject, err := s.subjectRepo.GetSubjectByID(ctx, lesson.SubjectID)
	if err != nil {
		return fmt.Errorf("get lesson subject: %w", err)
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	if err := s.lessonLogRepo.UpdateWithTx(ctx, tx, lesson); err != nil {
		return fmt.Errorf("cancel lesson: %w", err)
	}
	if err := s.updatePlannedSubjectLessons(ctx, tx, subject, lesson.AcademicPeriodID); err != nil {
		return fmt.Errorf("update planned lessons: %w", err)
	}

	return tx.Commit(ctx)
}

func (s *ScheduleService) RescheduleLesson(ctx context.Context, req *dto.ReplaceScheduleRequest) error {

	// не видит конфликтов в lesson logs
	// надо добавить проверку lesson на уже имеющиеся

	period, err := s.academicPeriodRepo.GetActive(ctx)
	if err != nil {
		return fmt.Errorf("get active academic period: %w", err)
	}
	subject, err := s.subjectRepo.GetSubjectByID(ctx, req.SubjectID)
	if err != nil {
		return fmt.Errorf("get subject: %w", err)
	}

	lessonDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return fmt.Errorf("parse lesson date: %w", err)
	}

	lesson := models.LessonLog{
		Date:             lessonDate,
		Number:           req.Number,
		SubjectID:        req.SubjectID,
		TeacherID:        req.TeacherID,
		AudienceID:       req.AudienceID,
		Comment:          req.Comment,
		Status:           models.LessonStatusRescheduled,
		AcademicPeriodID: period.ID,
	}

	if err := s.checkLessonConflict(ctx, &lesson); err != nil {
		return err
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := s.lessonLogRepo.CreateWithTx(ctx, tx, &lesson); err != nil {
		return fmt.Errorf("create rescheduled lesson: %w", err)
	}

	if err := s.updatePlannedSubjectLessons(ctx, tx, subject, period.ID); err != nil {
		return fmt.Errorf("update planned lessons: %w", err)
	}

	return tx.Commit(ctx)
}

func (s *ScheduleService) checkLessonConflict(ctx context.Context, l *models.LessonLog) error {

	// Фильтр для поиска существующих занятий
	filters := repository.LessonLogFilters{
		DateFrom:         &l.Date,
		DateTo:           &l.Date,
		Number:           &l.Number,
		AcademicPeriodID: &l.AcademicPeriodID,
	}

	existing, err := s.lessonLogRepo.GetAll(ctx, filters)
	if err != nil {
		return fmt.Errorf("failed to check lesson uniqueness: %w", err)
	}

	if len(existing) > 0 {

		for _, lesson := range existing {
			if lesson.Status == models.LessonStatusCanceled {
				continue
			}

			if lesson.TeacherID == l.TeacherID {
				return fmt.Errorf("teacher %d already has a lesson at %s, lesson %d",
					l.TeacherID, l.Date.Format("2006-01-02"), l.Number)
			}

			if lesson.AudienceID == l.AudienceID {
				return fmt.Errorf("audience %d already has a lesson at %s, lesson %d",
					l.AudienceID, l.Date.Format("2006-01-02"), l.Number)
			}

			if lesson.Subject.GroupID == l.Subject.GroupID {
				return fmt.Errorf("group %d already has a lesson at %s, lesson %d",
					l.Subject.GroupID, l.Date.Format("2006-01-02"), l.Number)
			}
		}

	}

	return nil
}

func (s *ScheduleService) CompleteLessonFromDate(ctx context.Context, date time.Time, comment string) error {
	lessons, err := s.lessonLogRepo.GetAll(ctx, repository.LessonLogFilters{
		DateFrom: &date,
		DateTo:   &date,
	})
	if err != nil {
		return fmt.Errorf("get lesson from date %s", date.Format("2006-01-02"))
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	for _, lesson := range lessons {
		if lesson.Status == models.LessonStatusCompleted {
			continue
		}
		if lesson.Status == models.LessonStatusCanceled {
			continue
		}
		if err := s.lessonLogRepo.UpdateStatusWithTx(ctx, tx, lesson.ID, models.LessonStatusCompleted, comment); err != nil {
			return fmt.Errorf("update status for lesson %d: %w", lesson.ID, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func (s *ScheduleService) CreateGroupWithCurriculumAndStudents(ctx context.Context, req *dto.CreateGroupWithCurriculumRequest) (*dto.GroupWithCurriculumResponse, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	group := models.Group{
		Name:          req.Group.Name,
		Specialty:     req.Group.Specialty,
		AdmissionYear: req.Group.AdmissionYear,
	}
	if err := s.groupRepo.CreateWithTx(ctx, tx, &group); err != nil {
		return nil, fmt.Errorf("create group: %w", err)
	}

	var createdSubjects []models.Subject
	for i, subjReq := range req.Subjects {
		startDate, err := time.Parse("2006-01-02", subjReq.StartDate)
		if err != nil {
			return nil, fmt.Errorf("parse start_date for subject %d: %w", i, err)
		}
		endDate, err := time.Parse("2006-01-02", subjReq.EndDate)
		if err != nil {
			return nil, fmt.Errorf("parse end_date for subject %d: %w", i, err)
		}
		if startDate.After(endDate) {
			return nil, fmt.Errorf("start_date must be before end_date for subject %d: %w", i, err)
		}
		subject := models.Subject{
			Title:     subjReq.Title,
			Semester:  subjReq.Semester,
			HoursLoad: subjReq.HoursLoad,
			StartDate: startDate,
			EndDate:   endDate,
			GroupID:   group.ID,
		}

		if err := s.subjectRepo.CreateSubjectWithTx(ctx, tx, &subject); err != nil {
			return nil, fmt.Errorf("create subject %d: %w", i, err)
		}

		subject.Group = group
		createdSubjects = append(createdSubjects, subject)
	}

	var studentsResult []dto.StudentCreationResult
	for i, studReq := range req.Students {
		baseUsername := generateUsername(studReq.FullName, group.Name)
		username, err := ensureUniqueUsername(ctx, s.userRepo, baseUsername)
		if err != nil {
			return nil, fmt.Errorf("generate unique username: %w", err)
		}

		plainPassword, err := generateRandomPassword()
		if err != nil {
			return nil, fmt.Errorf("generate password: %w", err)
		}
		passwordHash, err := hashPassword(plainPassword)
		if err != nil {
			return nil, fmt.Errorf("hash password: %w", err)
		}

		user := models.User{
			Name:         username,
			FullName:     studReq.FullName,
			Email:        studReq.Email,
			PasswordHash: passwordHash,
		}
		if err := s.userRepo.CreateWithTx(ctx, tx, &user); err != nil {
			return nil, fmt.Errorf("create user %d: %w", i, err)
		}

		if err := s.studentRepo.CreateStudentWithTx(ctx, tx, &user, group.ID); err != nil {
			return nil, fmt.Errorf("create student %d profile: %w", i, err)
		}

		studentsResult = append(studentsResult, dto.StudentCreationResult{
			User:     user,
			Password: plainPassword,
			GroupID:  group.ID,
		})
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	return &dto.GroupWithCurriculumResponse{
		Group:    group,
		Subjects: createdSubjects,
		Students: studentsResult,
	}, nil
}

// CreateGroupFromExcel создаёт группу, дисциплины и студентов из Excel файла
func (s *ScheduleService) CreateGroupFromExcel(ctx context.Context, file io.Reader) (*dto.GroupWithCurriculumResponse, error) {
	excelSvc := NewExcelService()
	req, err := excelSvc.ParseGroupFromExcel(file)
	if err != nil {
		return nil, fmt.Errorf("parse excel: %w", err)
	}
	return s.CreateGroupWithCurriculumAndStudents(ctx, req)
}

// GetLessonStatistics возвращает статистику занятий с учётом фильтров
func (s *ScheduleService) GetLessonStatistics(ctx context.Context, req *dto.StatisticsRequest) (*dto.StatisticsResponse, error) {
	filters := repository.LessonLogFilters{
		GroupID:          req.GroupID,
		TeacherID:        req.TeacherID,
		SubjectID:        req.SubjectID,
		AcademicPeriodID: req.AcademicPeriodID,
	}
	if filters.AcademicPeriodID == nil {
		period, err := s.academicPeriodRepo.GetActive(ctx)
		if err != nil {
			return nil, fmt.Errorf("get active academic period: %w", err)
		}
		filters.AcademicPeriodID = &period.ID
	}
	logs, err := s.lessonLogRepo.GetAll(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("get lesson logs: %w", err)
	}

	// Группируем логи по предметам
	statsMap := make(map[int]*dto.SubjectStatistics)
	for _, log := range logs {
		stat, ok := statsMap[log.SubjectID]
		if !ok {
			subj, err := s.subjectRepo.GetSubjectByID(ctx, log.SubjectID)
			if err != nil || subj == nil {
				continue
			}
			totalLessons := (subj.HoursLoad + 1) / 2
			stat = &dto.SubjectStatistics{
				SubjectID:    subj.ID,
				SubjectTitle: subj.Title,
				GroupID:      subj.GroupID,
				GroupName:    subj.Group.Name,
				TotalHours:   subj.HoursLoad,
				TotalLessons: totalLessons,
			}
			statsMap[log.SubjectID] = stat
		}
		switch log.Status {
		case models.LessonStatusCompleted:
			stat.CompletedLessons++
		case models.LessonStatusRescheduled:
			stat.RescheduledLessons++
		case models.LessonStatusCanceled:
			stat.CancelledLessons++
		case models.LessonStatusPlanned:
			stat.PlannedLessons++
		}
	}

	// Рассчитываем RemainingLessons, CompletionPercent, IsOnSchedule
	totalCompletedAll := 0
	totalLessonsAll := 0
	totalRemainingAll := 0

	for _, stat := range statsMap {
		// Оставшиеся пары = всего - проведено
		stat.RemainingLessons = stat.TotalLessons - stat.CompletedLessons
		if stat.RemainingLessons < 0 {
			stat.RemainingLessons = 0
		}
		if stat.TotalLessons > 0 {
			stat.CompletionPercent = float64(stat.CompletedLessons) / float64(stat.TotalLessons) * 100
		}
		// Определяем, успеваем ли по графику: если к текущей дате проведено не меньше, чем ожидалось по равномерному распределению
		// Для этого нужно знать дату начала предмета и окончания
		if stat.PlannedLessons >= stat.RemainingLessons {
			stat.IsOnSchedule = true
		}
		totalCompletedAll += stat.CompletedLessons
		totalLessonsAll += stat.TotalLessons
		totalRemainingAll += stat.RemainingLessons
	}

	// Общая статистика
	overallProgress := 0.0
	if totalLessonsAll > 0 {
		overallProgress = float64(totalCompletedAll) / float64(totalLessonsAll) * 100
	}

	resp := &dto.StatisticsResponse{
		TotalSubjects:       len(statsMap),
		OverallProgress:     overallProgress,
		TotalLessonsAll:     totalLessonsAll,
		CompletedLessonsAll: totalCompletedAll,
		RemainingLessonsAll: totalRemainingAll,
		Subjects:            make([]dto.SubjectStatistics, 0, len(statsMap)),
	}
	for _, stat := range statsMap {
		resp.Subjects = append(resp.Subjects, *stat)
	}
	return resp, nil
}

// CreateStudentsFromExcel создаёт студентов из Excel-файла
func (s *ScheduleService) CreateStudentsFromExcel(ctx context.Context, r io.Reader) ([]dto.StudentCreationResult, error) {
	excelSvc := NewExcelService()
	studentsData, err := excelSvc.ParseStudentsFromExcel(r)
	if err != nil {
		return nil, fmt.Errorf("parse excel: %w", err)
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	var results []dto.StudentCreationResult
	for _, data := range studentsData {
		// Найти группу по имени
		group, err := s.groupRepo.GetByName(ctx, data.GroupName)
		if err != nil {
			return nil, fmt.Errorf("find group '%s': %w", data.GroupName, err)
		}
		if group == nil {
			return nil, fmt.Errorf("group '%s' not found", data.GroupName)
		}

		// Генерация username
		baseUsername := generateUsername(data.FullName, group.Name)
		username, err := ensureUniqueUsername(ctx, s.userRepo, baseUsername)
		if err != nil {
			return nil, fmt.Errorf("generate username for %s: %w", data.FullName, err)
		}

		// Генерация пароля
		plainPassword, err := generateRandomPassword()
		if err != nil {
			return nil, fmt.Errorf("generate password: %w", err)
		}
		passwordHash, err := hashPassword(plainPassword)
		if err != nil {
			return nil, fmt.Errorf("hash password: %w", err)
		}

		user := models.User{
			Name:         username,
			FullName:     data.FullName,
			Email:        data.Email,
			PasswordHash: passwordHash,
		}

		// Создаём в транзакции

		if err := s.userRepo.CreateWithTx(ctx, tx, &user); err != nil {
			return nil, fmt.Errorf("create user: %w", err)
		}
		if err := s.studentRepo.CreateStudentWithTx(ctx, tx, &user, group.ID); err != nil {
			return nil, fmt.Errorf("create student profile: %w", err)
		}

		results = append(results, dto.StudentCreationResult{
			User:     user,
			Password: plainPassword,
			GroupID:  group.ID,
		})
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit tx: %w", err)
	}

	return results, nil
}

// CreateTeachersFromExcel создаёт преподавателей из Excel-файла
func (s *ScheduleService) CreateTeachersFromExcel(ctx context.Context, r io.Reader) ([]dto.TeacherCreationResult, error) {
	excelSvc := NewExcelService()
	teachersData, err := excelSvc.ParseTeachersFromExcel(r)
	if err != nil {
		return nil, fmt.Errorf("parse excel: %w", err)
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	var results []dto.TeacherCreationResult
	for _, data := range teachersData {
		// Генерация username из ФИО (латиница)
		baseUsername := generateUsername(data.FullName, "") // без группы
		username, err := ensureUniqueUsername(ctx, s.userRepo, baseUsername)
		if err != nil {
			return nil, fmt.Errorf("generate username for %s: %w", data.FullName, err)
		}

		plainPassword, err := generateRandomPassword()
		if err != nil {
			return nil, fmt.Errorf("generate password: %w", err)
		}
		passwordHash, err := hashPassword(plainPassword)
		if err != nil {
			return nil, fmt.Errorf("hash password: %w", err)
		}

		user := models.User{
			Name:         username,
			FullName:     data.FullName,
			Email:        data.Email,
			PasswordHash: passwordHash,
		}

		if err := s.userRepo.CreateWithTx(ctx, tx, &user); err != nil {
			return nil, fmt.Errorf("create user: %w", err)
		}
		if err := s.teacherRepo.CreateTeacherWithTx(ctx, tx, &user); err != nil { 
			return nil, fmt.Errorf("create teacher profile: %w", err)
		}

		results = append(results, dto.TeacherCreationResult{
			User:     user,
			Password: plainPassword,
		})
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit tx: %w", err)
	}
	return results, nil
}

// CreateAudiencesFromExcel создаёт аудитории из Excel-файла
func (s *ScheduleService) CreateAudiencesFromExcel(ctx context.Context, r io.Reader) ([]models.Audience, error) {
    excelSvc := NewExcelService()
    audiencesData, err := excelSvc.ParseAudiencesFromExcel(r)
    if err != nil {
        return nil, fmt.Errorf("parse excel: %w", err)
    }

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("create transaction: %w", err)
	}
	defer tx.Rollback(ctx)
    
    var created []models.Audience
    for _, data := range audiencesData {
        audience := models.Audience{
            Name:   data.Name,
            Number: data.Number,
        }
        if err := s.audienceRepo.CreateWithTx(ctx, tx, &audience); err != nil {
            return nil, fmt.Errorf("create audience %s: %w", data.Name, err)
        }
        created = append(created, audience)
    }

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

    return created, nil
}
