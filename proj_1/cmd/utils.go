package main

import (
	"html/template"
	"os"

	"github.com/gin-gonic/gin"
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
	}

	tmpl, err := template.New("").Funcs(functions).ParseGlob("templates/*.tmpl")
	if err != nil {
		return err
	}

	router.SetHTMLTemplate(tmpl)
	return nil
}
