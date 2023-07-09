// Brief   : ID generator using Soyflake (the same library as Zitadel)
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 09.07.2023

package helpers

import (
	"github.com/sony/sonyflake"
	"time"
)

var IdGen *sonyflake.Sonyflake

// InitIdGen initialises the sonyflake ID generator, a distributed unique ID generator.
// This initialisation is currently not checking the uniqueness of the machine ID, and
// the lower 16 bits of the IP address are used as a machine ID. This should be updated
// if the app is running distributed.
func InitIdGen() {
	IdGen = sonyflake.NewSonyflake(sonyflake.Settings{
		StartTime:      time.Now(),
		MachineID:      nil,
		CheckMachineID: nil,
	})
}
