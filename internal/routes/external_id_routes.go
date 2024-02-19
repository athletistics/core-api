// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 06.07.2023

package routes

import (
	"github.com/athletistics/core-api/internal/handlers/externalId"
	"github.com/gofiber/fiber/v2"
)

func SetupExternalIdRoutesTotpProtected(router fiber.Router) {
	node := router.Group("/externalId")

	node.Get("/user", externalId.SelectExternalUserIds)
	node.Post("/user", externalId.CreateExternalUserId)
	node.Get("/idTypes", externalId.GetTypes)
}
