// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 18.06.2023

package seeds

import "github.com/baeriswyln/goseeder"

func init() {
	goseeder.SetDataPath("database/seeds/data")
	goseeder.Register(languageSeeder)
	goseeder.Register(genderSeeder)
	goseeder.Register(accessRoleSeeder)
	goseeder.Register(eventTypeSeeder)

	goseeder.Register(countriesSeeder)
	goseeder.Register(statesSeeder)
	goseeder.Register(placesSeeder)

	goseeder.Register(sportsSeeder)
	goseeder.Register(clubsSeeder)

	goseeder.Register(eventSeeder)
	goseeder.Register(disciplineSeeder)

	goseeder.Register(usersSeeder)
	goseeder.Register(authorisationSeeder)
}
