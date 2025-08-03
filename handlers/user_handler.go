package handlers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	d "github.com/strbagus/fiber-auth/database"
	m "github.com/strbagus/fiber-auth/models"
)

func ListUsers(c *fiber.Ctx) error {
	p := new(m.Pages)
	if err := c.QueryParser(p); err != nil {
		return err
	}
	var total int
	offset := p.Offset
	limit := p.Limit
	if limit == 0 {
		limit = 10
	}
	err := d.DB.QueryRow("SELECT count(uuid) FROM users").Scan(&total)
	rows, err := d.DB.Query("SELECT uuid, username, fullname FROM users LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		log.Fatalf("DB Query Error: %s", err)
	}
	defer rows.Close()
	users := []m.User{}
	for rows.Next() {
		var user m.User
		err = rows.Scan(&user.UUID, &user.Username, &user.Fullname)
		if err != nil {
			log.Fatalf("DB Scan Error: %s", err)
		}
		users = append(users, user)
	}
	response := m.UserListResponse{
		Data: users,
		Pagination: m.Pages{
			Total:  total,
			Limit:  limit,
			Offset: offset,
		},
	}
	return c.Status(fiber.StatusOK).JSON(response)
}
