package main

import (
	"pizza-tracker-go/internal/models"
	"slices"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Hàm đăng ký Custom Validator
// Hàm này dùng để cấu hình và đăng ký các luật kiểm tra dữ liệu mới vào hệ thống validator của Gin.
func RegisterCustomValidators() {
	//binding.Validator.Engine().(*validator.Validate): Lấy ra instance thực tế của thư viện validator đang được Gin quản lý. Câu lệnh (*validator.Validate) là cú pháp ép kiểu trong Go.
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// Đăng ký một quy tắc (rule) kiểm tra mới với tên (tag) tùy chỉnh:
		//"valid_pizza_type": Kiểm tra xem loại pizza khách chọn có nằm trong danh sách models.PizzaTypes hay không.
		v.RegisterValidation("valid_pizza_type", createSliceValidator(models.PizzaTypes))
		//"valid_pizza_size": Kiểm tra xem kích cỡ pizza có nằm trong danh sách models.PizzaSizes hay không.
		v.RegisterValidation("valid_pizza_size", createSliceValidator(models.PizzaSizes))
	}
}

// Hàm tạo Validator dựa trên mảng
// Đây là một Higher-order function (hàm nhận vào một hàm hoặc trả về một hàm). Nó nhận vào một mảng các giá trị được phép (allowedValues []string) và trả về một hàm kiểm tra kiểu validator.Func mà thư viện validator yêu cầu.
func createSliceValidator(allowedValues []string) validator.Func {
	return func(fl validator.FieldLevel) bool {
		// slices.Contains(...): Kiểm tra xem giá trị chuỗi đó có nằm trong mảng
		return slices.Contains(allowedValues, fl.Field().String())
	}
}
