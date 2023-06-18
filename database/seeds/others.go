// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 19.06.2023

package seeds

import "github.com/baeriswyln/goseeder"

func sportsSeeder(s goseeder.Seeder) {
	s.FromJson("sports")
}

func clubsSeeder(s goseeder.Seeder) {
	s.FromJson("clubs")
}
