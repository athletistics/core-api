// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 20.06.2023

package res

import "github.com/google/uuid"

type UuidCreated struct {
	Id uuid.UUID `json:"id"`
}

type IdCreated struct {
	Id string `json:"id"`
}

type IntIdCreated struct {
	Id int64 `json:"id"`
}
