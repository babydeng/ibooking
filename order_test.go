package main

import (
	"Group4/ibooking-back/internal"
	"Group4/ibooking-back/store"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateNewOrder(t *testing.T) {
	// 创建一个模拟的数据库连接
	// 这里可以使用测试数据库、内存数据库或者模拟的数据库对象
	db := createMockDB()

	// 创建一个 HTTP 请求的记录器
	recorder := httptest.NewRecorder()

	// 构建一个请求体
	order := store.CreateUserOrderParams{
		Date:      "2023-05-18",
		UserID:    1,
		RoomID:    1,
		SeatID:    1,
		StartTime: "15",
		EndTime:   "18",
		// 其他订单信息字段...
	}
	orderJSON, _ := json.Marshal(order)
	req, _ := http.NewRequest("POST", "/order", bytes.NewBuffer(orderJSON))

	// 调用处理函数

	handler := handleCreateNewOrder(db)
	handler.ServeHTTP(recorder, req)

	// 检查响应状态码
	if recorder.Code == http.StatusAccepted {
		t.Log("The seat already booked by users")
	}

	// 解析响应体
	var response store.CreateUserOrderParams
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to parse response body: %v", err)
	}

	// 根据预期结果进行断言
	expectedOrder := store.CreateUserOrderParams{
		Date:      "2023-05-18",
		UserID:    1,
		RoomID:    1,
		SeatID:    1,
		StartTime: "15",
		EndTime:   "18",
		// 其他订单信息字段...
	}
	if response != expectedOrder {
		t.Errorf("Expected order %+v, but got %+v", expectedOrder, response)
	} else {
		t.Log("Create order successfully")
		t.Logf("%+v", response)
	}
}

// 创建一个模拟的数据库连接的辅助函数
func createMockDB() *sql.DB {
	// 创建并返回一个模拟的数据库连接
	dbURI := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		internal.GetAsString("DB_USER", "pi"),
		internal.GetAsString("DB_PASSWORD", "123456"),
		internal.GetAsString("DB_HOST", "10.177.29.226"),
		internal.GetAsInt("DB_PORT", 5432),
		internal.GetAsString("DB_NAME", "restapi"),
	)
	// Open the database
	db, err := sql.Open("postgres", dbURI)
	if err != nil {
		log.Fatalln("Error opening database:", err)
	}

	// Connectivity check
	if err := db.Ping(); err != nil {
		log.Fatalln("Error from database ping:", err)
	}
	return db
}

func TestHandleListOrders(t *testing.T) {
	// 创建一个临时数据库用于测试
	db := createMockDB()
	// 创建一个带有路由的路由器
	router := mux.NewRouter()
	router.HandleFunc("/order/{date}", handleListOrders(db)).Methods("GET")

	t.Run("Valid Request", func(t *testing.T) {
		// 构造请求
		req, err := http.NewRequest("GET", "/order/2023-05-19?roomID=1", nil)
		assert.NoError(t, err, "Failed to create request")

		// 创建响应记录器
		rr := httptest.NewRecorder()
		// 将请求发送到路由器进行处理
		router.ServeHTTP(rr, req)
		// 验证响应状态码
		assert.Equal(t, http.StatusOK, rr.Code, "Unexpected status code")

		// 解码响应的 JSON 数据
		var response []store.IbookingOrder
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err, "Failed to decode response JSON")
		t.Logf("The response %+v", response)

	})

	t.Run("Invalid Date Format", func(t *testing.T) {
		// 构造请求
		req, err := http.NewRequest("GET", "/order/2023-5-18?roomID=1", nil)
		assert.NoError(t, err, "Failed to create request")

		// 创建响应记录器
		rr := httptest.NewRecorder()
		// 将请求发送到路由器进行处理
		router.ServeHTTP(rr, req)
		// 验证响应状态码
		assert.Equal(t, http.StatusBadRequest, rr.Code, "Unexpected status code")

		// 解码响应的 JSON 数据
		var response []store.IbookingOrder
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err, "Failed to decode response JSON")
		t.Logf("The response %+v", response)
	})

}

func TestHandleDeleteOrder(t *testing.T) {
	// 创建一个临时数据库用于测试
	db := createMockDB()
	// 创建一个带有路由的路由器
	router := mux.NewRouter()
	router.HandleFunc("/order/{order_id}", handleDeleteOrder(db)).Methods("DELETE")

	// 创建一个虚拟的 HTTP 请求和响应对象
	req, err := http.NewRequest("DELETE", "/order/13", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	// 检查响应状态码
	if rr.Code != http.StatusOK {
		t.Errorf("Delete fail %+v, but got status %d", rr.Body, rr.Code)
	} else {
		t.Logf("Delete order successfully")
	}

}
