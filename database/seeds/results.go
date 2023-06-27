// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 22.06.2023

package seeds

import "github.com/baeriswyln/goseeder"

func resultSeeder(s goseeder.Seeder) {
	s.FromJson("results")
}

func resultStaticSeeder(s goseeder.Seeder) {
	s.FromJson("result_statics")
}

func resultStaticTypeSeeder(s goseeder.Seeder) {
	s.FromJson("result_static_types")
}

func categorySeeder(s goseeder.Seeder) {
	s.FromJson("categories")
}
