// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2

package store

import ()

type IbookingOrder struct {
	OrderID   int64  `json:"order_id"`
	UserID    int64  `json:"user_id"`
	RoomID    int64  `json:"room_id"`
	SeatID    int64  `json:"seat_id"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Date      string `json:"date"`
	Status    int32  `json:"status"`
}

type IbookingRoom struct {
	RoomID      int64  `json:"room_id"`
	RoomName    string `json:"room_name"`
	Description string `json:"description"`
}

type IbookingSeat struct {
	SeatID      int64  `json:"seat_id"`
	RoomID      int64  `json:"room_id"`
	SeatNum     string `json:"seat_num"`
	Description string `json:"description"`
}

type IbookingUser struct {
	UserID       int64  `json:"user_id"`
	UserNum      string `json:"user_num"`
	PasswordHash string `json:"password_hash"`
	UserName     string `json:"user_name"`
	Credit       int32  `json:"credit"`
}