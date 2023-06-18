// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 18.06.2023

package seeds

import (
	"github.com/kristijorgji/goseeder"
)

func countriesSeeder(s goseeder.Seeder) {
	s.FromJson("countries")
}
