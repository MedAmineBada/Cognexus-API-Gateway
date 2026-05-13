package config

import "strings"

type RouteFlag struct {
	Path     string
	Method   string
	FlagName string
}

var AuthRouteFlags = []RouteFlag{
	// Public routes
	{Path: "/api/v1/auth/student/register", Method: "POST", FlagName: "auth.student.register"},
	{Path: "/api/v1/auth/student/login", Method: "POST", FlagName: "auth.student.login"},
	{Path: "/api/v1/auth/teacher/register", Method: "POST", FlagName: "auth.teacher.register"},
	{Path: "/api/v1/auth/teacher/login", Method: "POST", FlagName: "auth.teacher.login"},
	{Path: "/api/v1/auth/admin/login", Method: "POST", FlagName: "auth.admin.login"},
	{Path: "/api/v1/auth/verify-otp", Method: "POST", FlagName: "auth.verify-otp"},
	{Path: "/api/v1/auth/resend-otp", Method: "POST", FlagName: "auth.resend-otp"},
	{Path: "/api/v1/auth/forgot-password", Method: "POST", FlagName: "auth.forgot-password"},
	{Path: "/api/v1/auth/reset-password", Method: "POST", FlagName: "auth.reset-password"},

	// Protected routes
	{Path: "/api/v1/auth/logout", Method: "POST", FlagName: "auth.logout"},
	{Path: "/api/v1/auth/refresh", Method: "POST", FlagName: "auth.refresh"},

	// Admin user management
	{Path: "/api/v1/admin/users", Method: "POST", FlagName: "auth.admin.users.create"},

	// User profile
	{Path: "/api/v1/users/me", Method: "GET", FlagName: "auth.users.me.get"},
	{Path: "/api/v1/users/me", Method: "PUT", FlagName: "auth.users.me.update"},
	{Path: "/api/v1/users/me/password", Method: "PUT", FlagName: "auth.users.me.password"},
	{Path: "/api/v1/users/me/image", Method: "PUT", FlagName: "auth.users.me.image.upload"},
	{Path: "/api/v1/users/me/image", Method: "DELETE", FlagName: "auth.users.me.image.delete"},

	// User queries
	{Path: "/api/v1/users", Method: "GET", FlagName: "auth.users.list"},
	{Path: "/api/v1/users/students", Method: "GET", FlagName: "auth.users.students"},
	{Path: "/api/v1/users/teachers", Method: "GET", FlagName: "auth.users.teachers"},
	{Path: "/api/v1/users/lead-teachers", Method: "GET", FlagName: "auth.users.lead-teachers"},
	{Path: "/api/v1/users/by-email", Method: "GET", FlagName: "auth.users.by-email"},
	{Path: "/api/v1/users/:user_id", Method: "GET", FlagName: "auth.users.by-id"},

	// User management
	{Path: "/api/v1/users/:user_id/role", Method: "PATCH", FlagName: "auth.users.role"},
	{Path: "/api/v1/users/:user_id", Method: "DELETE", FlagName: "auth.users.delete"},
	{Path: "/api/v1/users/:user_id/reactivate", Method: "PATCH", FlagName: "auth.users.reactivate"},

	// Internal (service-to-service)
	{Path: "/api/v1/auth/status/:user_id", Method: "GET", FlagName: "auth.status.get"},
	{Path: "/api/v1/auth/status/:user_id", Method: "PATCH", FlagName: "auth.status.update"},
}

var ExamRouteFlags = []RouteFlag{
	{Path: "/api/v1/exam/create", Method: "POST", FlagName: "exam_orchestrator.create_exam"},
	{Path: "/api/v1/exam/save", Method: "POST", FlagName: "exam_orchestrator.save_exam"},
	{Path: "/api/v1/exam/get", Method: "GET", FlagName: "exam_orchestrator.get_exam"},
	{Path: "/api/v1/exam/correction/create", Method: "POST", FlagName: "exam_orchestrator.create_correction"},
	{Path: "/api/v1/exam/correction/save", Method: "POST", FlagName: "exam_orchestrator.save_correction"},
	{Path: "/api/v1/exam/correction", Method: "GET", FlagName: "exam_orchestrator.get_correction"},
}

var ClassesRouteFlags = []RouteFlag{
	{Path: "/api/v1/classes", Method: "POST", FlagName: "classes.create"},
	{Path: "/api/v1/classes", Method: "GET", FlagName: "classes.list"},
	{Path: "/api/v1/classes/:class_id", Method: "GET", FlagName: "classes.get"},
	{Path: "/api/v1/classes/:class_id", Method: "PUT", FlagName: "classes.update"},
	{Path: "/api/v1/classes/:class_id/soft", Method: "DELETE", FlagName: "classes.soft-delete"},
	{Path: "/api/v1/classes/:class_id/permanent", Method: "DELETE", FlagName: "classes.permanent-delete"},

	{Path: "/api/v1/classes/:class_id/subjects", Method: "POST", FlagName: "classes.subjects.add"},
	{Path: "/api/v1/classes/:class_id/subjects", Method: "GET", FlagName: "classes.subjects.list"},
	{Path: "/api/v1/classes/:class_id/subjects/:subject_id", Method: "PUT", FlagName: "classes.subjects.update"},
	{Path: "/api/v1/classes/:class_id/subjects/:subject_id/soft", Method: "DELETE", FlagName: "classes.subjects.soft-delete"},
	{Path: "/api/v1/classes/:class_id/subjects/:subject_id/permanent", Method: "DELETE", FlagName: "classes.subjects.permanent-delete"},

	{Path: "/api/v1/classes/:class_id/enroll", Method: "POST", FlagName: "classes.enroll"},
	{Path: "/api/v1/classes/:class_id/students", Method: "GET", FlagName: "classes.students.list"},
	{Path: "/api/v1/classes/:class_id/students/:student_id", Method: "PUT", FlagName: "classes.students.update"},
	{Path: "/api/v1/classes/:class_id/students/:student_id/soft", Method: "DELETE", FlagName: "classes.students.soft-delete"},
	{Path: "/api/v1/classes/:class_id/students/:student_id/permanent", Method: "DELETE", FlagName: "classes.students.permanent-delete"},
}

func GetFlagForPath(path, method string) string {
	if strings.HasSuffix(path, "/submit") && strings.HasPrefix(path, "/api/v1/exam/") {
		return "exam_orchestrator.submit"
	}

	for _, route := range AuthRouteFlags {
		if route.Path == path && route.Method == method {
			return route.FlagName
		}
	}

	for _, route := range ExamRouteFlags {
		if route.Path == path && route.Method == method {
			return route.FlagName
		}
	}

	for _, route := range ClassesRouteFlags {
		if route.Path == path && route.Method == method {
			return route.FlagName
		}
	}

	return ""
}
