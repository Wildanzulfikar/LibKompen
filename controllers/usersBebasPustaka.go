package controllers

import (
	"LibKompen/database"
	"LibKompen/models"
	"LibKompen/models/request"

	"github.com/gofiber/fiber/v2"
)

func UsersBebasPustakaGetAll(c *fiber.Ctx) error {
	var users []models.UsersBebasPustaka
	database.DB.Find(&users)
	return c.JSON(users)
}

func CreateUsersBebasPustaka(c *fiber.Ctx) error {
	createUsers := new(request.UsersRequestBebasPustaka)

	if err := c.BodyParser(createUsers); err != nil {
		return err
	}

	newUsers := models.UsersBebasPustaka{
		Name:     createUsers.Name,
		Email:    createUsers.Email,
		Username: createUsers.Username,
		Password: createUsers.Password,
		Role:     createUsers.Role,
		Status:   createUsers.Status,
	}

	errCreateUsers := database.DB.Create(&newUsers)

	if errCreateUsers != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "failed to store data",
		})
	}

	return c.JSON(fiber.Map{
		"message": "success",
		"data":    newUsers,
	})
}

func DeleteUsersBebasPustaka(c *fiber.Ctx) error {
	userId := c.Params("id_users")
	var user models.UsersBebasPustaka

	err := database.DB.Debug().First(&user, "id_users=?", userId).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "user not found",
		})
	}

	errDelete := database.DB.Debug().Delete(&user).Error
	if errDelete != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "internal server error",
		})
	}

	return c.JSON(fiber.Map{
		"message": "user was deleted",
	})
}
