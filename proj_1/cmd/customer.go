package main

import (
	"log/slog"
	"net/http"
	"pizza-tracker-go/internal/models"

	"github.com/gin-gonic/gin"
)

type CustomerData struct {
	Title    string
	Order    models.Order
	Statuses []string
}

// Cấu trúc chứa dữ liệu truyền vào giao diện form đặt hàng (order.tmpl)
// bao gồm danh sách các loại và kích cỡ pizza có sẵn.
type OrderFormData struct {
	PizzaTypes []string
	PizzaSizes []string
}

// Đại diện cho dữ liệu form gửi lên từ người dùng
// Các thẻ form:"..." ánh xạ tên trường HTML vào struct, và thẻ binding:"..." khai báo các quy tắc kiểm tra (validation) như bắt buộc nhập (required), giới hạn độ dài ký tự, hoặc dùng các custom validator đã đăng ký (valid_pizza_size, valid_pizza_type).
type OrderRequest struct {
	Name         string   `form:"name" binding:"required,min=2,max=100"`
	Phone        string   `form:"phone" binding:"required,min=10,max=20"`
	Address      string   `form:"address" binding:"required,min=5,max=200"`
	Sizes        []string `form:"size" binding:"required,min=1,dive,valid_pizza_size"`
	PizzaTypes   []string `form:"pizza" binding:"required,min=1,dive,valid_pizza_type"`
	Instructions []string `form:"instructions" binding:"max=200"`
}

// Xử lý yêu cầu truy cập trang chủ hiển thị form.
func (h *Handler) ServeNewOrderForm(c *gin.Context) {
	//c.HTML(...): Trả về một trang HTML sử dụng template có tên "order.tmpl" với mã trạng thái 200 OK, đồng thời truyền danh sách loại và kích cỡ pizza vào giao diện để người dùng lựa chọn.
	c.HTML(http.StatusOK, "order.tmpl", OrderFormData{
		PizzaTypes: models.PizzaTypes,
		PizzaSizes: models.PizzaSizes,
	})
}

// Hàm xử lý khi người dùng gửi form đặt hàng
func (h *Handler) HandleNewOrderPost(c *gin.Context) {
	var form OrderRequest

	//c.ShouldBind(&form): Tự động đọc dữ liệu từ form HTML gửi lên, gắn vào biến form và chạy các bộ kiểm tra hợp lệ (validator). Nếu dữ liệu không hợp lệ, trả về lỗi 400 Bad Request kèm theo thông báo lỗi.
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Chuyển đổi dữ liệu: Duyệt qua danh sách các món pizza khách đặt để tạo thành một mảng các mục chi tiết models.OrderItem (gồm kích cỡ, loại pizza và ghi chú tương ứng).
	orderItems := make([]models.OrderItem, len(form.Sizes))
	for i := range orderItems {
		orderItems[i] = models.OrderItem{
			Size:        form.Sizes[i],
			Pizza:       form.PizzaTypes[i],
			Instruction: form.Instructions[i],
		}
	}

	//Khởi tạo đơn hàng: Gom nhóm thông tin khách hàng, gán trạng thái ban đầu là trạng thái đầu tiên trong danh sách ("Order placed") và đính kèm danh sách các món ăn (Items) vào cấu trúc models.Order.
	order := models.Order{
		CustomerName: form.Name,
		Phone:        form.Phone,
		Address:      form.Address,
		Status:       models.OrderStatuses[0],
		Items:        orderItems,
	}

	//Lưu vào database: Gọi hàm CreateOrder để lưu đơn hàng xuống cơ sở dữ liệu. Nếu có lỗi xảy ra, ghi lại log bằng slog.Error và trả về mã lỗi 500 Internal Server Error.
	if err := h.orders.CreateOrder(&order); err != nil {
		slog.Error("Failed to create order", "error", err)
		c.String(http.StatusInternalServerError, "Something went wrong")
		return
	}

	//Thành công và chuyển hướng: Ghi log thông báo tạo đơn thành công, sau đó chuyển hướng (c.Redirect) người dùng đến trang theo dõi đơn hàng riêng của họ dựa theo mã ID (/customer/:id).
	slog.Info("Order created", "orderId", order.ID, "customer", order.CustomerName)

	c.Redirect(http.StatusSeeOther, "/customer/"+order.ID)
}

func (h *Handler) serveCustomer(c *gin.Context) {
	//c.Param("id"): Lấy tham số id động từ đường dẫn URL (ví dụ: /customer/abc123xyz). Nếu không tìm thấy ID, trả về lỗi 400.
	orderID := c.Param("id")
	if orderID == "" {
		c.String(http.StatusBadRequest, "Order ID is required")
	}

	//Truy vấn đơn hàng: Gọi hàm GetOrder để tìm thông tin đơn hàng cùng danh sách các món (Preload("Items")) trong database dựa vào mã ID. Nếu không tìm thấy, trả về lỗi 404 Not Found.
	order, err := h.orders.GetOrder(orderID)
	if err != nil {
		c.String(http.StatusNotFound, "Order not found")
		return
	}

	//Hiển thị giao diện: Trả về file giao diện "customer.tmpl" với mã trạng thái 200 OK, kèm theo đối tượng dữ liệu đơn hàng vừa tìm được để hiển thị tiến trình làm pizza cho khách hàng.
	c.HTML(http.StatusOK, "customer.tmpl", CustomerData{
		Title:    "PizzaOrder Status" + orderID,
		Order:    *order,
		Statuses: models.OrderStatuses,
	})
}
