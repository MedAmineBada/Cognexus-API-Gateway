package config

import (
	"github.com/google/uuid"
	"regexp"
	"strconv"
	"strings"
)

var RouteRoles = map[string][]string{
	// ==================== AUTH ====================
	// Public routes - no roles needed (empty = no restriction)

	// Admin user management
	"POST:/api/v1/admin/users": {"admin"},

	// User profile - all authenticated users (no restriction, JWT is enough)
	// GET:/api/v1/users/me        - any authenticated
	// PUT:/api/v1/users/me        - any authenticated
	// PUT:/api/v1/users/me/password - any authenticated
	// PUT:/api/v1/users/me/image  - any authenticated
	// DELETE:/api/v1/users/me/image - any authenticated

	// User queries
	"GET:/api/v1/users":               {"admin", "lead_teacher"},
	"GET:/api/v1/users/students":      {"admin", "lead_teacher"},
	"GET:/api/v1/users/teachers":      {"admin", "lead_teacher"},
	"GET:/api/v1/users/lead-teachers": {"admin", "lead_teacher"},
	"GET:/api/v1/users/by-email":      {"admin", "lead_teacher", "teacher"},
	"GET:/api/v1/users/:user_id":      {"admin", "lead_teacher", "teacher"},

	// User management
	"PATCH:/api/v1/users/:user_id/role":       {"admin"},
	"DELETE:/api/v1/users/:user_id":           {"admin", "lead_teacher"},
	"PATCH:/api/v1/users/:user_id/reactivate": {"admin"},

	// ==================== CLASSES ====================
	// Class CRUD
	"POST:/api/v1/classes":                       {"admin", "lead_teacher"},
	"GET:/api/v1/classes":                        {"admin", "lead_teacher", "teacher", "student"},
	"GET:/api/v1/classes/:class_id":              {"admin", "lead_teacher", "teacher", "student"},
	"PUT:/api/v1/classes/:class_id":              {"admin", "lead_teacher"},
	"DELETE:/api/v1/classes/:class_id/soft":      {"admin", "lead_teacher"},
	"DELETE:/api/v1/classes/:class_id/permanent": {"admin"},

	// Subjects
	"POST:/api/v1/classes/:class_id/subjects":                         {"admin", "lead_teacher"},
	"GET:/api/v1/classes/:class_id/subjects":                          {"admin", "lead_teacher", "teacher", "student"},
	"PUT:/api/v1/classes/:class_id/subjects/:subject_id":              {"admin", "lead_teacher"},
	"DELETE:/api/v1/classes/:class_id/subjects/:subject_id/soft":      {"admin", "lead_teacher"},
	"DELETE:/api/v1/classes/:class_id/subjects/:subject_id/permanent": {"admin"},

	// Enrollments
	"POST:/api/v1/classes/:class_id/enroll":                           {"admin", "lead_teacher"},
	"GET:/api/v1/classes/:class_id/students":                          {"admin", "lead_teacher", "teacher"},
	"PUT:/api/v1/classes/:class_id/students/:student_id":              {"admin", "lead_teacher"},
	"DELETE:/api/v1/classes/:class_id/students/:student_id/soft":      {"admin", "lead_teacher"},
	"DELETE:/api/v1/classes/:class_id/students/:student_id/permanent": {"admin"},

	// ==================== EXAM ====================
	"POST:/api/v1/exam/create":            {"lead_teacher", "teacher"},
	"POST:/api/v1/exam/save":              {"lead_teacher", "teacher"},
	"GET:/api/v1/exam/get":                {"lead_teacher", "teacher"},
	"GET:/api/v1/exam/:exam_id":           {"lead_teacher", "teacher", "student"},
	"POST:/api/v1/exam/correction/create": {"lead_teacher", "teacher"},
	"POST:/api/v1/exam/correction/save":   {"lead_teacher", "teacher"},
	"GET:/api/v1/exam/correction":         {"lead_teacher", "teacher"},
	"POST:/api/v1/exam/:exam_id/submit":   {"student"},
	"GET:/api/v1/exam/:exam_id/report":    {"lead_teacher", "teacher"},
}

func NormalizePath(path string) string {
	segments := strings.Split(path, "/")
	for i, segment := range segments {
		if isID(segment) {
			// Check context to determine placeholder name
			if i > 0 {
				prevSegment := segments[i-1]
				switch prevSegment {
				case "users":
					segments[i] = ":user_id"
				case "classes":
					segments[i] = ":class_id"
				case "subjects":
					segments[i] = ":subject_id"
				case "students":
					segments[i] = ":student_id"
				case "exam":
					segments[i] = ":exam_id"
				case "status":
					segments[i] = ":user_id" // for /auth/status/{user_id}
				default:
					segments[i] = ":id"
				}
			} else {
				segments[i] = ":id"
			}
		}
	}
	return strings.Join(segments, "/")
}

func isID(segment string) bool {
	// Check if it's a number
	if _, err := strconv.Atoi(segment); err == nil {
		return true
	}

	// Check if it's a UUID
	if _, err := uuid.Parse(segment); err == nil {
		return true
	}

	// Check if it's a MongoDB ObjectID (24 hex chars)
	if matched, _ := regexp.MatchString("^[a-f0-9]{24}$", segment); matched {
		return true
	}

	// Check if it's a typical ID pattern (alphanumeric, may contain hyphens/underscores)
	if matched, _ := regexp.MatchString("^[a-zA-Z0-9_-]{8,}$", segment); matched {
		return true
	}

	return false
}
