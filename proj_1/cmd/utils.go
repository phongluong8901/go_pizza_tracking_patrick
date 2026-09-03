package main

import (
	"encoding/json"
	"html/template"
	"os"

	"github.com/gin-contrib/sessions"
	gormsessions "github.com/gin-contrib/sessions/gorm"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// : Định nghĩa một cấu trúc chứa các thông số cấu hình cốt lõi của ứng dụng, gồm cổng chạy server (Port) và đường dẫn file cơ sở dữ liệu (DBPath).
type Config struct {
	Port   string
	DBPath string
}

// loadConfig: Hàm khởi tạo và trả về một đối tượng Config. Nó gọi hàm phụ trợ getEnv để lấy giá trị từ hệ thống, nếu không có sẽ lấy giá trị mặc định (cổng mặc định là 8080, cơ sở dữ liệu là pizza_tracker.db).
func loadConfig() Config {
	return Config{
		Port:   getEnv("PORT", "8080"),
		DBPath: getEnv("DB_PATH", "pizza_tracker.db"),
	}
}

// getEnv: Kiểm tra xem biến môi trường với tên key đã được thiết lập trên hệ thống hay chưa (os.Getenv(key)). Nếu có giá trị, nó trả về giá trị đó; ngược lại, nó trả về defaultValue.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Hàm nhận vào con trỏ quản lý router của Gin và trả về lỗi nếu có.
func loadTemplates(router *gin.Engine) error {
	//template.FuncMap: Cho phép đăng ký thêm các hàm tùy chỉnh để sử dụng trực tiếp bên trong các file giao diện HTML (.tmpl). Ở đây, hàm "add" được tạo ra để thực hiện phép cộng hai số nguyên (a + b) ngay trên giao diện.
	functions := template.FuncMap{
		"add": func(a, b int) int { return a + b },
		"json": func(v interface{}) template.JS {
			a, err := json.Marshal(v)
			if err != nil {
				return template.JS("")
			}
			return template.JS(string(a))
		},
	}

	tmpl, err := template.New("").Funcs(functions).ParseGlob("templates/*.tmpl")
	if err != nil {
		return err
	}

	router.SetHTMLTemplate(tmpl)
	return nil
}

func setupSessionStore(db *gorm.DB, secretKey []byte) sessions.Store {
	store := gormsessions.NewStore(db, true, secretKey)
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   86400,
		HttpOnly: true,
		Secure:   true,
		SameSite: 3,
	})

	return store
}

func SetSessionValue(c *gin.Context, key string, value interface{}) error {
	session := sessions.Default(c)
	session.Set(key, value)
	return session.Save()
}

func GetSessionString(c *gin.Context, key string) string {
	session := sessions.Default(c)
	val := session.Get(key)
	if val == nil {
		return ""
	}

	str, _ := val.(string)
	return str
}

func ClearSession(c *gin.Context) error {
	session := sessions.Default(c)
	session.Clear()
	return session.Save()
}
