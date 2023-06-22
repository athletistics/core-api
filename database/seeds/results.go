// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 22.06.2023

package seeds

import "github.com/baeriswyln/goseeder"

func resultSeeder(s goseeder.Seeder) {
	s.FromJson("results")
}
