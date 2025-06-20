package db

import "LiveLive/model"

func CreateCourse(course *model.Course) (err error) {
	return Mysql.Create(course).Error
}

func FindCourseByClassnameAndTeacherId(classname string, TeacherId int64) (*model.Course, error) {
	var course model.Course
	err := Mysql.Where("teacher_id = ? AND classname = ?", TeacherId, classname).First(&course).Error
	return &course, err

}

func FindCourseByTeacherId(teacherId int) (*model.Course, error) {
	var course model.Course
	err := Mysql.Where("teacher_id = ?", teacherId).First(&course).Error
	return &course, err
}

func FindCourseByCourseName(courseName string) (*model.Course, error) {
	var course model.Course
	err := Mysql.Where("classname = ?", courseName).First(&course).Error
	return &course, err
}

func AddStudentCourse(course *model.CourseMember) error {
	return Mysql.Create(course).Error
}

func FindCourseInviteByCode(code string) (*model.CourseInvite, error) {
	var courseInvite model.CourseInvite
	err := Mysql.Where("code = ?", code).First(&courseInvite).Error
	return &courseInvite, err

}

func AddCourseInvite(course *model.CourseInvite) error {
	return Mysql.Create(course).Error
}

func FindCourseInviteByCourseId(courseId uint) (*model.CourseInvite, error) {
	var courseInvite model.CourseInvite
	err := Mysql.Where("course_id = ?", courseId).First(&courseInvite).Error
	return &courseInvite, err
}

func FindCourseMemberByCourseIdAndStudentId(courseId int64, studentId int64) (*model.CourseMember, error) {
	var courseMember model.CourseMember
	err := Mysql.Where("course_id = ? AND student_id = ?", courseId, studentId).First(&courseMember).Error
	return &courseMember, err
}

func FindStudentByTeacherNameAndClassName(teacherName string, ClassName string) (*[]model.CourseMember, error) {
	var students []model.CourseMember
	err := Mysql.Where("classname = ? AND teacher_name = ?", ClassName, teacherName).Find(&students).Error
	return &students, err
}
