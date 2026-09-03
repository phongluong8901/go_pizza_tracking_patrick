package models

import (
	"time"

	"github.com/teris-io/shortid"
	"gorm.io/gorm"
)

var (
	//OrderStatuses: Khai báo một mảng chứa các trạng thái tuần tự của một đơn hàng từ lúc đặt cho đến khi sẵn sàng.
	OrderStatuses = []string{"Order placed", "Preparing", "Baking", "Quality Check", "Ready"}

	//PizzaTypes: Mảng chứa danh sách các hương vị/loại pizza cửa hàng có phục vụ.
	PizzaTypes = []string{
		"Margherita",
		"Pepperoni",
		"Vegetarian",
		"Hawaiian",
		"Bbq Chicken",
		"Meat Lovers",
		"Buffalo Chicken",
		"Supreme",
		"Truffle Mushroom",
		"Four Cheese",
	}

	//PizzaSizes: Mảng chứa các kích cỡ pizza có sẵn.
	PizzaSizes = []string{
		"Small", "Medium", "Large", "X-Large",
	}
)

// OrderModel: Một struct đóng vai trò là "service" hoặc "repository" chứa kết nối cơ sở dữ liệu (*gorm.DB). Các hàm thao tác với database của đơn hàng sẽ được gắn vào struct này.
type OrderModel struct {
	DB *gorm.DB
}

// Đại diện cho bảng orders trong cơ sở dữ liệu.
type Order struct {
	ID           string      `gorm:"primaryKey; size:14" json:"id"` //Khóa chính của bảng, kiểu chuỗi (string), giới hạn độ dài 14 ký tự (size:14). Khi trả về dạng JSON, trường này có tên là id.
	Status       string      `gorm:"not null" json:"status"`        //bắt buộc phải có giá trị (not null).
	CustomerName string      `gorm:"not null" json:"customer_name"`
	Phone        string      `gorm:"not null" json:"phone"`
	Address      string      `gorm:"not null" json:"address"`
	Items        []OrderItem `gorm:"foreignKey:OrderID" json:"pizzas"` //Danh sách các món ăn trong đơn này (kiểu mảng []OrderItem). Thẻ gorm:"foreignKey:OrderID" báo cho GORM biết bảng OrderItem liên kết với bảng này thông qua khóa ngoại OrderID.
	CreatedAt    time.Time   `json:"createdAt"`
}

// Đại diện cho bảng order_items (chi tiết từng chiếc pizza trong một đơn hàng).
type OrderItem struct {
	ID          string `gorm:"primaryKey;size:14" json:"id"`
	OrderID     string `gorm:index;` //Khóa ngoại trỏ ngược về bảng Order. gorm:"index;" giúp đánh chỉ mục để tìm kiếm nhanh hơn.
	Size        string `gorm:"not null" json:"size"`
	Pizza       string `gorm:"not null" json:"pizza"`
	Instruction string `json:"instruction"`
}

// -Gói hook tự động tạo ID
// Đây là một hàm hook đặc biệt của GORM. Nó sẽ tự động chạy ngay trước khi câu lệnh INSERT được thực thi xuống cơ sở dữ liệu cho bảng Order
func (o *Order) BeforeCreate(tx *gorm.DB) error {
	//Nếu người dùng chưa gán ID (o.ID == ""), hàm sẽ gọi thư viện shortid.MustGenerate() để tự tạo một chuỗi ID ngẫu nhiên, ngắn gọn và gán vào.
	if o.ID == "" {
		o.ID = shortid.MustGenerate()
	}

	return nil
}

// Tương tự như trên, nhưng áp dụng cho từng mục sản phẩm (OrderItem) để tự động sinh mã ID ngắn nếu nó chưa có sẵn.
func (oi *OrderItem) BeforeCreate(tx *gorm.DB) error {
	if oi.ID == "" {
		oi.ID = shortid.MustGenerate()
	}

	return nil
}

// -Các hàm thao tác với Cơ sở dữ liệu
// CreateOrder: Nhận vào một con trỏ kiểu *Order và thực hiện câu lệnh INSERT vào database thông qua o.DB.Create(order)
func (o *OrderModel) CreateOrder(order *Order) error {
	return o.DB.Create(order).Error
}

// GetOrder: Dùng để tìm kiếm một đơn hàng dựa vào mã id
func (o *OrderModel) GetOrder(id string) (*Order, error) {
	var order Order
	//Câu lệnh này cực kỳ quan trọng, nó bảo GORM thực hiện lấy luôn cả bảng chi tiết sản phẩm (Items) liên kết với đơn hàng đó (tránh việc chỉ lấy thông tin chung mà bỏ quên danh sách pizza khách đặt).
	//First(&order, "id = ?", id): Lấy bản ghi đầu tiên khớp với điều kiện id = id truyền vào.
	err := o.DB.Preload("Items").First(&order, "id = ?", id).Error
	return &order, err
}
