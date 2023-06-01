CREATE SCHEMA IF NOT EXISTS ibooking;

CREATE TABLE ibooking.users (
                                User_ID        BIGSERIAL PRIMARY KEY,
                                User_Num       TEXT NOT NULL,
                                Password_Hash  TEXT NOT NULL,
                                User_Name      TEXT NOT NULL,
                                Credit         INT NOT NULL DEFAULT 0,
                                UNIQUE (User_Num)
);
-- SQLc converts snake_case to CamelCase

CREATE TABLE ibooking.orders (
                                 Order_ID   BIGSERIAL PRIMARY KEY,
                                 User_ID    BIGINT NOT NULL,
                                 Room_ID    BIGINT NOT NULL,
                                 Seat_ID    BIGINT NOT NULL,
                                 Start_Time TEXT NOT NULL,
                                 End_Time   TEXT NOT NULL,
                                 Date       TEXT NOT NULL,
                                 Status     INT NOT NULL DEFAULT 0
);


CREATE TABLE ibooking.rooms (
                                    Room_ID     BIGINT PRIMARY KEY,
                                    Room_Name   TEXT NOT NULL,
                                    Description TEXT NOT NULL
);

CREATE TABLE ibooking.seats (
                               Seat_ID          BIGSERIAL PRIMARY KEY,
                               Room_ID          BIGINT NOT NULL,
                               Seat_Num         TEXT NOT NULL,
                               Description      TEXT NOT NULL
);

