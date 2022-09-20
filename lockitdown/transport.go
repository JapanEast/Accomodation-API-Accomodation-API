
package lockitdown

import (
	"fmt"
	"strconv"
)

type (
	TransportRobot struct {
		Player        int  `json:"player"`
		Dir           Pair `json:"dir"`
		IsLocked      bool `json:"isLocked"`
		IsBeamEnabled bool `json:"isBeamEnabled"`
	}
	TransportRobots []interface{}
	TransportState  struct {
		GameDef          GameDef           `json:"gameDef"`