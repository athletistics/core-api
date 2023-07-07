// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 07.07.2023

package seeds

import "github.com/baeriswyln/goseeder"

func externalIdTypeSeeder(s goseeder.Seeder) {
	s.FromJson("external_id_types")
}
