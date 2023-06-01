-- name: CreateUserOrder :one
INSERT INTO ibooking.orders (
    User_ID,
    Room_ID,
    Seat_ID,
    Start_Time,
    Date,
    End_Time
) VALUES (
             $1,
             $2,
             $3,
             $4,
             $5,
             $6
         )  RETURNING *;

-- name: ListUserOrders :many
SELECT *
FROM ibooking.orders
WHERE User_ID = $1;

-- name: GetOrderBySeat :many
SELECT *
FROM ibooking.orders
WHERE Date = $1 AND Seat_ID = $2 AND Room_ID = $3;

-- name: GetOrderByRoom :many
SELECT *
FROM ibooking.orders
WHERE Date = $1 AND Room_ID = $2;

-- name: GetOrder :one
SELECT *
FROM ibooking.orders
WHERE Order_ID = $1;

-- name: GetOrderByUserID :many
SELECT *
FROM ibooking.orders
WHERE User_ID = $1;

-- name: DeleteOrder :exec
DELETE FROM ibooking.orders
WHERE Order_ID = $1;

-- name: CreateRoom :one
INSERT INTO ibooking.rooms (
    Room_Name,
    Description
) VALUES (
             $1,
             $2
         ) RETURNING *;

-- name: CreateSeat :one
INSERT INTO ibooking.seats (
    Seat_Num,
    Room_ID,
    Description
) VALUES (
             $1,
             $2,
             $3
         ) RETURNING *;


-- name: UpdateOrderStatus :one
UPDATE ibooking.orders SET
        Status = $1
WHERE Order_ID = $2 RETURNING *;
