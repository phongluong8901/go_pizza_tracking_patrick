package main

import "github.com/gin-gonic/gin"

// Hàm nhận vào con trỏ quản lý router tổng của Gin (*gin.Engine) và một con trỏ tầng xử lý *Handler (h) để gán các hàm logic tương ứng với từng đường dẫn URL.
func setupRoutes(router *gin.Engine, h *Handler) {
	router.GET("/", h.ServeNewOrderForm)            //ServeNewOrderForm để hiển thị giao diện form đặt món cho người dùng.
	router.POST("/new-order", h.HandleNewOrderPost) //HandleNewOrderPost để tiếp nhận, kiểm tra dữ liệu và lưu đơn hàng mới từ form.
	router.GET("/customer/:id", h.serveCustomer)    //serveCustomer để tra cứu và hiển thị trang theo dõi tiến trình đơn hàng của khách hàng cụ thể.

	//Cấu hình phục vụ các tệp tài nguyên tĩnh (như file CSS, hình ảnh, JavaScript)
	router.Static("/static", "/templates/static")
}
