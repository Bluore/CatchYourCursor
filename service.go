package main

import (
	"encoding/json"
	"fmt"
)

func AddCursor(hub *Hub, req Operation) error {
	originCursorList := hub.cursorList

	var cursor Cursor
	err := json.Unmarshal(req.Data, &cursor)
	if err != nil {
		return err
	}

	cursor.ID = req.ID
	hub.cursorList = append(hub.cursorList, cursor)

	err = responseMassageOne(hub, Listed, originCursorList, req.ID)
	if err != nil {
		return err
	}

	err = responseMassageAll(hub, Added, cursor, req.ID)
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
	case CursorCheck:
		err = GetCursorCheck(hub, operation)
	case Catch:
		err = CatchCursorService(hub, operation)
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

	return responseMassageAll(hub, Moved, cursor, cursor.ID)
}

func responseMassageAll(hub *Hub, status operStatus, data interface{}, exp string) error {
	response := &Response{
		Status: status,
		Data:   data,
	}
	return sendMassageAll(hub, response, exp)
}

func sendMassageAll(hub *Hub, msg interface{}, exp string) error {
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
			hub.unregister <- client
		}
	}

	return err
}
func responseMassageOne(hub *Hub, status operStatus, data interface{}, aim string) error {
	response := &Response{
		Status: status,
		Data:   data,
	}
	return sendMassageOne(hub, response, aim)
}

func sendMassageOne(hub *Hub, msg interface{}, aim string) error {
	sendMsg, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	for client := range hub.clients {
		if client.id != aim {
			continue
		}
		select {
		case client.send <- sendMsg:
		default:
			hub.unregister <- client
		}
		break
	}

	return err
}

func SendCursorCheck(hub *Hub) error {
	originCursorList := hub.cursorList
	hub.cursorList = hub.cursorList[:0]
	return responseMassageAll(hub, CursorCheck, originCursorList, "")
}

func GetCursorCheck(hub *Hub, req Operation) error {
	var cursor Cursor
	err := json.Unmarshal(req.Data, &cursor)
	if err != nil {
		return err
	}

	hub.cursorList = append(hub.cursorList, cursor)
	return nil
}

func CatchCursorService(hub *Hub, req Operation) error {
	var cursor Cursor
	err := json.Unmarshal(req.Data, &cursor)
	if err != nil {
		return err
	}

	return responseMassageAll(hub, Catch, cursor, cursor.ID)
}
