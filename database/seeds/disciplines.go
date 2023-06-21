// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 21.06.2023

package seeds

import "github.com/baeriswyln/goseeder"

func disciplineSeeder(s goseeder.Seeder) {
	s.FromJson("disciplines")
}
