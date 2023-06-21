// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 21.06.2023

package seeds

import "github.com/baeriswyln/goseeder"

func eventTypeSeeder(s goseeder.Seeder) {
	s.FromJson("event_types")
}

func eventSeeder(s goseeder.Seeder) {
	s.FromJson("events")
}
