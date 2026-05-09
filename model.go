package main

import "encoding/json"

type operStatus int

const (
	Moved       operStatus = 1
	Deleted     operStatus = 2
	Listed      operStatus = 3
	Catch       operStatus = 4
	Added       operStatus = 5
	GetID       operStatus = 6
	CursorCheck operStatus = 7
)

type Cursor struct {
	ID string  `json:"id"`
	X  float64 `json:"x"`
	Y  float64 `json:"y"`
}

type Operation struct {
	Status operStatus      `json:"status"`
	ID     string          `json:"id"`
	Data   json.RawMessage `json:"data"`
}

type Response struct {
	Status operStatus  `json:"status"`
	Data   interface{} `json:"data"`
}

type AddCursorResp struct {
	ID string `json:"id"`
}

type GetClientIDResp struct {
	ID string `json:"id"`
}
