package main

import "pizza-tracker-go/internal/models"

// Định nghĩa cấu trúc Handler đóng vai trò là tầng điều khiển (Controller), nơi chứa các hàm xử lý logic HTTP request từ người dùng.
type Handler struct {
	orders              *models.OrderModel
	users               *models.UserModel
	notificationManager *NotificationManager
}

// Hàm khởi tạo (constructor) nhận vào một con trỏ *models.DBModel (chứa toàn bộ kết nối database đã được thiết lập sẵn) và trả về một con trỏ *Handler đã được cấu hình hoàn chỉnh.
func NewHandler(dbModel *models.DBModel) *Handler {
	return &Handler{
		orders:              &dbModel.Order, //Trích xuất phần quản lý đơn hàng từ DBModel tổng rồi gắn vào trường orders của Handler
		users:               &dbModel.User,
		notificationManager: NewNotificationManager(),
	}
}
