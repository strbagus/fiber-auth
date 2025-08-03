package handlers

import (
	"context"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	d "github.com/strbagus/fiber-auth/database"
	m "github.com/strbagus/fiber-auth/models"
	u "github.com/strbagus/fiber-auth/utils"
)

func SignIn(c *fiber.Ctx) error {
	var creds m.UserCred
	if err := c.BodyParser(&creds); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request payload",
			"error":   err.Error(),
		})
	}

	var user m.User
	err := d.DB.QueryRow(`SELECT uuid, username, password, fullname FROM users WHERE username = $1`, creds.Username).
		Scan(&user.UUID, &user.Username, &user.Password, &user.Fullname)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid username or password",
		})
	}

	if !u.CheckPassword(user.Password, creds.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid username or password",
		})
	}

	accessTokenDuration := time.Minute * 15
	refreshTokenDuration := time.Hour * 16

	accessToken, err := u.GenerateToken(map[string]interface{}{
		"uid":   user.UUID,
		"uname": user.Username,
		"fname": user.Fullname,
	}, accessTokenDuration)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to generate access token",
			"error":   err.Error(),
		})
	}

	// Refresh token
	refreshToken, err := u.GenerateToken(map[string]interface{}{
		"uuid": user.UUID,
	}, refreshTokenDuration)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to generate refresh token",
			"error":   err.Error(),
		})
	}

	// Set cookies
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		HTTPOnly: true,
		SameSite: "Lax",
		Expires:  time.Now().Add(accessTokenDuration),
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HTTPOnly: true,
		SameSite: "Lax",
		Expires:  time.Now().Add(refreshTokenDuration),
	})

	// Only return necessary user fields
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user": fiber.Map{
			"uuid":     user.UUID,
			"username": user.Username,
			"fullname": user.Fullname,
		},
	})
}

type ReqToken struct {
	AccessToken string `json:"access_token"`
}

type MyCustomClaims struct {
	Foo string `json:"foo"`
	jwt.RegisteredClaims
}

func CheckToken(c *fiber.Ctx) error {
	if !u.IsJWTUsable(c.Cookies("access_token")) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Token is not valid",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Token is valid",
	})
}

func InvalidateToken(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")
	if !u.IsJWTUsable(refreshToken) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Refresh Token is not Valid",
		})
	}
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "JWT Parse Error",
		})
	}
	if token.Valid {
		expUnix := int64(claims["exp"].(float64))
		expTime := time.Unix(expUnix, 0)
		now := time.Now()
		d.RD.Set(context.Background(), "blacklist:"+refreshToken, true, expTime.Sub(now))
	}

	c.ClearCookie("access_token", "refresh_token")

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Logged Out",
	})
}

func RefreshToken(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")
	if !u.IsJWTUsable(refreshToken) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid or missing refresh token",
		})
	}

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid refresh token",
		})
	}

	userUUID, ok := claims["uuid"].(string)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid token claims",
		})
	}

	var user m.User
	err = d.DB.QueryRow(`SELECT uuid, username, fullname FROM users WHERE uuid = $1`, userUUID).
		Scan(&user.UUID, &user.Username, &user.Fullname)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "User not found or DB error",
			"error":   err.Error(),
		})
	}

	expiration := time.Minute * 15
	accessToken, err := u.GenerateToken(map[string]interface{}{
		"uid":   user.UUID,
		"uname": user.Username,
		"fname": user.Fullname,
	}, expiration)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Token generation failed",
			"error":   err.Error(),
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		HTTPOnly: true,
		SameSite: "Lax",
		Expires:  time.Now().Add(expiration),
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Access token refreshed",
	})
}
