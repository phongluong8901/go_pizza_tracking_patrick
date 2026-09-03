package main

import (
	"log/slog"
	"os"
	"pizza-tracker-go/internal/models"

	"github.com/gin-gonic/gin"
)

func main() {
	//Gọi hàm nạp cấu hình từ biến môi trường (cổng chạy server và đường dẫn file database).
	cfg := loadConfig()

	//Cấu hình Logger: Tạo một trình ghi log dạng văn bản (TextHandler) xuất ra màn hình console (os.Stdout) và đặt nó làm bộ ghi log mặc định toàn cục cho ứng dụng (slog.SetDefault).
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	//Khởi tạo Database: Gọi hàm mở kết nối và chạy migration cho SQLite dựa trên đường dẫn cấu hình (cfg.DBPath). Nếu gặp lỗi, ghi log lỗi rồi dừng ứng dụng ngay lập tức với mã os.Exit(1). Nếu thành công, ghi log thông báo.
	dbModel, err := models.InitDB(cfg.DBPath)
	if err != nil {
		slog.Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}

	slog.Info("Database initialized successfully")

	//Đăng ký Validator & Khởi tạo Handler: Kích hoạt các quy tắc kiểm tra dữ liệu tùy chỉnh, tạo đối tượng điều khiển (Handler) chứa kết nối cơ sở dữ liệu, và khởi tạo router Gin mặc định (đã tích hợp sẵn middleware logger và recovery).
	RegisterCustomValidators()

	h := NewHandler(dbModel)

	router := gin.Default()

	//Nạp Template và Định tuyến: Tải toàn bộ các file giao diện HTML (.tmpl) vào router, sau đó gắn kết các đường dẫn URL với các hàm xử lý tương ứng thông qua setupRoutes.
	if err := loadTemplates(router); err != nil {
		slog.Error("Failed to load templates", "error", err)
		os.Exit(1)
	}

	setupRoutes(router, h)

	//Khởi động Server: In log thông báo địa chỉ truy cập ứng dụng, sau đó gọi router.Run để lắng nghe các yêu cầu HTTP từ client trên cổng đã cấu hình (ví dụ: :8080).
	slog.Info("Server starting", "url", "http://localhost:"+cfg.Port)

	router.Run(":" + cfg.Port)

}
