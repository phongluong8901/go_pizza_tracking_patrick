1. Gin
2. Gorm
3. Ginsessions
4. Gorm store
5. Bcrypt
6. validtor
7. shortID

# --- lib
"time": Thư viện chuẩn dùng để xử lý thời gian (ví dụ: ngày giờ tạo đơn).

(https://github.com/teris-io/shortid)": Thư viện bên ngoài dùng để tạo mã định danh (ID) ngắn, độc nhất và an toàn cho URL.

"gorm.io/gorm": Thư viện ORM (Object-Relational Mapping) phổ biến trong Go giúp tương tác với cơ sở dữ liệu quan hệ (PostgreSQL, MySQL, SQLite,...) thông qua ngôn ngữ lập trình thay vì viết câu lệnh SQL thủ công.

"gorm.io/driver/sqlite": Driver kết nối cơ sở dữ liệu SQLite dành riêng cho GORM, giúp Go có thể đọc ghi file cơ sở dữ liệu SQLite.

"fmt": Thư viện chuẩn dùng để định dạng chuỗi, ở đây được dùng để tạo thông báo lỗi tùy chỉnh kèm theo lỗi gốc (fmt.Errorf).

"slices": Thư viện chuẩn của Go (từ phiên bản 1.21 trở lên) cung cấp các hàm tiện ích xử lý mảng/slice, ví dụ như kiểm tra phần tử có tồn tại trong mảng không.

(https://github.com/gin-gonic/gin/binding)": Thư viện ràng buộc dữ liệu của framework Gin, cho phép truy cập vào bộ kiểm tra dữ liệu (validator) mặc định của Gin.

(https://github.com/go-playground/validator/v10)": Thư viện validator ngầm định mà Gin sử dụng để kiểm tra tính hợp lệ của struct (ví dụ: validate:"required").

"log/slog" và "os": Thư viện dùng để ghi log chuẩn và thao tác với hệ điều hành (như kết thúc chương trình bằng mã lỗi).

# --- stack