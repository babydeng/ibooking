package main

import (
	"Group4/ibooking-back/internal"
	"Group4/ibooking-back/internal/api"
	"Group4/ibooking-back/store"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
	"strconv"
)

var (
	cookieStore = sessions.NewCookieStore([]byte("forDemo"))
)

func init() {
	cookieStore.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   60 * 15,
		HttpOnly: true,
	}
}

func handleLogin(db *sql.DB) http.HandlerFunc {
	return func(wr http.ResponseWriter, req *http.Request) {

		// Thanks to our middleware, we know we have JSON
		// we'll decode it into our request type and see if it's valid
		type loginRequest struct {
			Username string `json:"username,omitempty"`
			Password string `json:"password,omitempty"`
		}

		payload := loginRequest{}
		if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
			log.Println("Error decoding the body", err)
			api.JSONError(wr, http.StatusBadRequest, "Error decoding JSON")
			return
		}

		querier := store.New(db)
		user, err := querier.GetUserByName(req.Context(), payload.Username)
		if errors.Is(err, sql.ErrNoRows) || !internal.CheckPasswordHash(payload.Password, user.PasswordHash) {
			api.JSONError(wr, http.StatusForbidden, "Bad Credentials")
			return
		}
		if err != nil {
			log.Println("Received error looking up user", err)
			api.JSONError(wr, http.StatusInternalServerError, "Couldn't log you in due to a server error")
			return
		}

		// We're valid. Let's tell the user and set a cookie
		// 如果获取名为session-name的会话，如果会话不存在，则返回一个新会话
		session, err := cookieStore.Get(req, "session-name")
		if err != nil {
			log.Println("Cookie store failed with", err)
			api.JSONError(wr, http.StatusInternalServerError, "Session Error")
		}
		// 将以下键值对存储在session中
		session.Values["userAuthenticated"] = true
		session.Values["userID"] = user.UserID
		session.Save(req, wr)
	}
}

func checkSecret(db *sql.DB) http.HandlerFunc {
	return func(wr http.ResponseWriter, req *http.Request) {
		userDetails, _ := userFromSession(req)

		querier := store.New(db)
		user, err := querier.GetUser(req.Context(), userDetails.UserID)
		if errors.Is(err, sql.ErrNoRows) {
			api.JSONError(wr, http.StatusForbidden, "User not found")
			return
		}

		api.JSONMessage(wr, http.StatusOK, fmt.Sprintf("Hello there %s", user.UserName))
	}
}

func handleCreateNewOrder(db *sql.DB) http.HandlerFunc {
	return func(wr http.ResponseWriter, req *http.Request) {
		// 从请求中解析订单信息
		var orderParam store.CreateUserOrderParams
		err := json.NewDecoder(req.Body).Decode(&orderParam)
		if err != nil {
			api.JSONError(wr, http.StatusBadRequest, "Invalid request payload", err.Error())
			return
		}

		var timeParam store.GetOrderBySeatParams
		timeParam.Date = orderParam.Date
		timeParam.RoomID = orderParam.RoomID
		timeParam.SeatID = orderParam.SeatID
		querier := store.New(db)
		// query orders by time
		orders, err := querier.GetOrderBySeat(req.Context(), timeParam)
		if err != nil {
			api.JSONError(wr, http.StatusInternalServerError, err.Error())
			return
		}

		// convert start and end time to float64
		startTime, _ := strconv.ParseFloat(orderParam.StartTime, 64)
		endTime, _ := strconv.ParseFloat(orderParam.EndTime, 64)

		// check if the seat is already booked
		if len(orders) != 0 {
			for _, order := range orders {
				// convert order start and end time to float64
				orderStartTime, _ := strconv.ParseFloat(order.StartTime, 64)
				orderEndTime, _ := strconv.ParseFloat(order.EndTime, 64)
				if endTime > orderStartTime && orderEndTime > startTime {
					api.JSONMessage(wr, http.StatusAccepted, "The seat already booked by users", order.StartTime, order.EndTime)
					return
				}
			}
		}

		res, err := querier.CreateUserOrder(req.Context(), orderParam)
		if err != nil {
			api.JSONError(wr, http.StatusInternalServerError, err.Error())
			return
		}

		json.NewEncoder(wr).Encode(&res)
	}
}

func handleListOrders(db *sql.DB) http.HandlerFunc {
	return func(wr http.ResponseWriter, req *http.Request) {
		date := mux.Vars(req)["date"]
		fmt.Printf("The date: %v\n", date)
		// fmt.Printf("The result is: %v\n", internal.CheckDateFormat(date))

		if !internal.CheckDateFormat(date) {
			api.JSONError(wr, http.StatusBadRequest, "Bad date format")
			return
		}

		roomIDStr := req.URL.Query().Get("roomID")
		roomID, err := strconv.Atoi(roomIDStr)
		if err != nil {
			api.JSONError(wr, http.StatusBadRequest, "Bad roomID format")
			return
		}

		var oderByRoomParams store.GetOrderByRoomParams
		oderByRoomParams.Date = date
		oderByRoomParams.RoomID = int64(roomID)

		querier := store.New(db)
		orders, err := querier.GetOrderByRoom(req.Context(), oderByRoomParams)
		// fmt.Println("get orders: ", orders)
		if err != nil {
			api.JSONError(wr, http.StatusInternalServerError, err.Error())
			return
		}
		json.NewEncoder(wr).Encode(&orders)
	}
}

func handleDeleteOrder(db *sql.DB) http.HandlerFunc {
	// 判断时间冲突，在时间内不运行取消预约
	return func(wr http.ResponseWriter, req *http.Request) {
		// fmt.Println("??")
		// fmt.Println(mux.Vars(req)["order_id"])
		orderID, err := strconv.Atoi(mux.Vars(req)["order_id"])
		if err != nil {
			api.JSONError(wr, http.StatusBadRequest, "Bad order_id")
			return
		}

		querier := store.New(db)
		order, err := querier.GetOrder(req.Context(), int64(orderID))
		if err != nil {
			api.JSONError(wr, http.StatusInternalServerError, err.Error())
			return
		}

		// if order time < now, can't delete
		if !internal.CheckTime(order) {
			api.JSONError(wr, http.StatusBadRequest, "Can't delete order")
			return
		}

		err = store.New(db).DeleteOrder(req.Context(), int64(orderID))
		if err != nil {
			api.JSONError(wr, http.StatusBadRequest, "Bad order_id", err.Error())
			return
		}

		api.JSONMessage(wr, http.StatusOK, fmt.Sprintf("Order %+v is deleted", order))
	}
}
