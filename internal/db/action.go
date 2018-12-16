package db

//Action represents the outcome and info of an action to be performed
type Action struct {
	tableName struct{} `sql:"actions"`
	ActionID  int      `sql:"action_id,type:serial,pk"`
	RoomID    int      `sql:"room_id"`
	Type      int      `sql:"type"`
}
