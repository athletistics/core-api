// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 18.06.2023

package seeds

import "github.com/baeriswyln/goseeder"

func placesSeeder(s goseeder.Seeder) {
	s.FromJson("places")
}
