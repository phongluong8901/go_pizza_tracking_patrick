package models

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Tạo một struct tổng hợp để chứa tất cả các mô hình dữ liệu
type DBModel struct {
	Order OrderModel
}

func InitDB(dataSourceName string) (*DBModel, error) {
	//gorm.Open(...): Mở kết nối đến cơ sở dữ liệu SQLite thông qua driver
	//sqlite.Open(dataSourceName). Biến &gorm.Config{} dùng để cấu hình nâng cao cho GORM (ở đây để trống nghĩa là dùng mặc định).
	db, err := gorm.Open(sqlite.Open(dataSourceName), &gorm.Config{})

	if err != nil {
		return nil, fmt.Errorf("Failed to migrate dtabse: %v", err)
	}

	//db.AutoMigrate(...): Tính năng tự động tạo bảng (Migration) của GORM. Nó sẽ tự đọc cấu trúc của struct Order và OrderItem để tự động tạo ra các bảng orders và order_items trong file SQLite tương ứng.
	err = db.AutoMigrate(&Order{}, &OrderItem{})
	if err != nil {
		return nil, fmt.Errorf("Failed to migrate database %v", err)
	}

	//Khởi tạo và trả về: Tạo một thể hiện (instance) của DBModel, trong đó gắn kết nối db vừa tạo vào bên trong OrderModel (thông qua trường DB: db)
	dbModel := &DBModel{
		Order: OrderModel{DB: db},
	}

	return dbModel, nil
}
