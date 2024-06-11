// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 11.06.2024

package helpers

import (
	"github.com/athletistics/core-api/constants"
	res "github.com/athletistics/core-api/internal/responses"
	"gorm.io/gorm"
)

type Pagination struct {
	PageIndex *int `query:"page_index"`
	PageSize  *int `query:"page_size"`
}

func Paginate(r *Pagination, p *res.Page) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		p.PageIndex = 0
		var pageSize = constants.DefaultPageSize

		if r.PageIndex != nil && *r.PageIndex >= 0 {
			p.PageIndex = *r.PageIndex
		}
		if r.PageSize != nil && *r.PageSize > 0 {
			pageSize = *r.PageSize
		}

		return db.Offset(p.PageIndex * pageSize).Limit(pageSize)
	}
}
