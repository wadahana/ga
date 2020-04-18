package display

import (
	"encoding/binary"
	"math"
	"github.com/wadahana/memu"
)

const (
	MouseDown  int = memu.MouseDown
	MouseUp    int = memu.MouseUp
	MouseWheel int = memu.MouseWheel
	MouseMove  int = memu.MouseMove

	LeftButton   int = memu.LeftButton
	MiddleButton int = memu.MiddleButton
	RightButton  int = memu.RightButton
)

type MouseEvent struct {
	eventType   int
	mouseType 	int
	buttonType  int
	x,y 		float32
}

func newMouseEvent(msg []byte) Event {
	eventType := int(msg[0]);
	if eventType != 1 {
		return nil;
	}
	e := &MouseEvent{eventType: 1};
	e.mouseType   = int(msg[1]);
	e.buttonType  = byteToInt16(msg[2:4])
	e.x           = byteToFloat32(msg[4:8])
	e.y 		  = byteToFloat32(msg[8:12])
	return e;
}

func (e *MouseEvent) getEventType() int {
	return e.eventType;
}
func (e *MouseEvent) getMouseType() int {
	return e.mouseType;
}
func (e *MouseEvent) getButtonType() int {
	return e.buttonType;
}
func (e *MouseEvent) getWheelDelta() int {
	return e.buttonType;
}
func (e *MouseEvent) getPos() (float32, float32) {
	return e.x, e.y
}


type KeyboardEvent struct {
	eventType   int
	press       bool
	keyCode     int32
}

func newKeyboardEvent(msg []byte) Event {
	eventType := int(msg[0]);
	if eventType != 2 {
		return nil;
	}
	e := &KeyboardEvent{eventType:2};

	e.press   = msg[1] != 0;
	e.keyCode = int32(binary.BigEndian.Uint32(msg[4:8]))
	
	return e;
}
func (e *KeyboardEvent) getEventType() int {
	return e.eventType;
}

func (e *KeyboardEvent) getPress() bool {
	return e.press;
}

func (e *KeyboardEvent) getKeyCode() int32 {
	return e.keyCode;
}

func byteToInt16(bytes []byte) int {
	v := int16(binary.BigEndian.Uint16(bytes))
	return int(v);
}

func byteToFloat32(bytes []byte) float32 {
    bits := binary.BigEndian.Uint32(bytes)
    return math.Float32frombits(bits)
}