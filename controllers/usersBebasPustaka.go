package controllers

import (
	"LibKompen/database"
	"LibKompen/models"
	"LibKompen/models/request"
	"log"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func UsersBebasPustakaGetAll(c *fiber.Ctx) error {
	var users []models.UsersBebasPustaka
	database.DB.Find(&users)
	return c.JSON(users)
}

func CreateUsersBebasPustaka(c *fiber.Ctx) error {
	createUsers := new(request.UsersRequestBebasPustaka)

	// Step 1: Parse request
	if err := c.BodyParser(createUsers); err != nil {
		log.Println("BODY PARSER ERROR:", err)
		return c.Status(400).JSON(fiber.Map{
			"message": "invalid request",
			"error":   err.Error(),
		})
	}
	log.Println("REQUEST BODY:", createUsers)

	// Step 2: Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(createUsers.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("HASH PASSWORD ERROR:", err)
		return c.Status(500).JSON(fiber.Map{
			"message": "failed to hash password",
			"error":   err.Error(),
		})
	}
	log.Println("HASHED PASSWORD:", string(hashedPassword))
	log.Println("LENGTH HASHED PASSWORD:", len(hashedPassword))

	// Step 3: Prepare new user struct
	newUsers := models.UsersBebasPustaka{
		Name:     createUsers.Name,
		Email:    createUsers.Email,
		Username: createUsers.Username,
		Password: string(hashedPassword),
		Role:     createUsers.Role,
		Status:   createUsers.Status,
	}

	// Step 4: Debug raw SQL
	db := database.DB.Debug()
	errCreateUsers := db.Create(&newUsers).Error
	if errCreateUsers != nil {
		log.Println("ERROR CREATE USER:", errCreateUsers)
		return c.Status(500).JSON(fiber.Map{
			"message": "failed to store data",
			"error":   errCreateUsers.Error(),
		})
	}

	// Step 5: Success
	log.Println("USER CREATED SUCCESSFULLY:", newUsers)
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
