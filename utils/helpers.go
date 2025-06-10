package utils
import(
	"errors"
	"strconv"
	"net/http"
	"github.com/youssef-abbih/go-todo-list/middleware"
)
func GetUserID(r *http.Request) (uint, error) {
	// 1. Extract user ID from context
	userIDVal := r.Context().Value(middleware.UserContextKey)
	if userIDVal == nil {
		return 0, errors.New("user not authorized")
	}

	userIDStr, ok := userIDVal.(string)
	if !ok {
		return 0, errors.New("invalid user ID type")
	}

	// 2. Convert string to uint
	userIDInt, err := strconv.Atoi(userIDStr)
	if err != nil {
		return 0, err
	}

	return uint(userIDInt), nil
}

// parseID extracts the task ID from the URL path
// func parseID(path string) (int, error) {
// 	// Example: /tasks/5 → "5"
// 	idStr := strings.TrimPrefix(path, "/tasks/")
// 	return strconv.Atoi(idStr)
// }

