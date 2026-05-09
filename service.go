package main

import (
	"encoding/json"
	"fmt"
)

func AddCursor(hub *Hub, req Operation) error {
	//msg, ok := req.Data.([]byte)
	//if !ok {
	//	return fmt.Errorf("AddCursor: invalid data type %T", req.Data)
	//}

	var cursor Cursor
	err := json.Unmarshal(req.Data, &cursor)
	if err != nil {
		return err
	}

	cursor.ID = req.ID
	hub.cursorList = append(hub.cursorList, cursor)

	err = responseMassage(hub, Listed, hub.cursorList, cursor.ID)
	if err != nil {
		return err
	}

	return nil
}

func DealMassage(hub *Hub, msg []byte) error {
	var operation Operation
	err := json.Unmarshal(msg, &operation)
	if err != nil {
		fmt.Println(fmt.Sprintf("operation unmarshal error: %v", err))
		return err
	}

	switch operation.Status {
	case Added:
		err = AddCursor(hub, operation)
	case Moved:
		err = MoveCursor(hub, operation)
		// todo other oper

	}
	if err != nil {
		fmt.Println(fmt.Sprintf("operation deal cursor error: %v", err))
		return err
	}

	return nil
}

func GetClientID(c *Client) error {
	msg := &Response{
		Status: GetID,
		Data:   GetClientIDResp{ID: c.id},
	}
	sendMsg, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	select {
	case c.send <- sendMsg:
	default:
	}
	return err
}

func MoveCursor(hub *Hub, req Operation) error {
	var cursor Cursor
	err := json.Unmarshal(req.Data, &cursor)
	if err != nil {
		return err
	}

	return responseMassage(hub, Moved, cursor, cursor.ID)
}

func responseMassage(hub *Hub, status operStatus, data interface{}, exp string) error {
	response := &Response{
		Status: status,
		Data:   data,
	}
	return sendMassage(hub, response, exp)
}

func sendMassage(hub *Hub, msg interface{}, exp string) error {
	sendMsg, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	for client := range hub.clients {
		if client.id == exp {
			continue
		}
		select {
		case client.send <- sendMsg:
		default:
		}
	}

	return err
}
