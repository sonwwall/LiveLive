package db

import "LiveLive/model"

func CreateCourse(course *model.Course) (err error) {
	return Mysql.Create(course).Error
}

func FindCourseByClassname(classname string) (*model.Course, error) {
	var course model.Course
	err := Mysql.Where("classname = ?", classname).First(&course).Error
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
