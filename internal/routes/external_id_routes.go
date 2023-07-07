// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 06.07.2023

package routes

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/internal/handlers/externalId"
)

func SetupExternalIdRoutesTotpProtected(router fiber.Router) {
	node := router.Group("/externalId")

	node.Get("/user", externalId.SelectExternalUserIds)
	node.Post("/user", externalId.CreateExternalUserId)
	node.Get("/idTypes", externalId.GetTypes)
}
