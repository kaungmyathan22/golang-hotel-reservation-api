package middlewares

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	customeErrors "github.com/kaungmyathan22/golang-hotel-reservation/src/errors"
	repository "github.com/kaungmyathan22/golang-hotel-reservation/src/repositories"
)

func JWTAuthentication(userStore repository.UserStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth_headers, ok := c.GetReqHeaders()["Authorization"]
		if !ok {
			return fmt.Errorf("unauthorized")
		}
		splitted_token := strings.Split(auth_headers[0], "Bearer ")

		token := splitted_token[1]
		claims, err := parseJWTToken(token)
		if err != nil {
			return err
		}

		expiresFloat := claims["expires"].(float64)
		expires := int64(expiresFloat)
		if time.Now().Unix() > expires {
			return fmt.Errorf("token expired")
		}
		userID := claims["userID"].(string)
		user, err := userStore.GetUserByID(c.Context(), userID)
		if err != nil {
			return fmt.Errorf("unauthorized")
		}
		c.Context().SetUserValue("user", user)
		return c.Next()
	}
}

func parseJWTToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("invalid signing method", token.Header["alg"])
			return nil, customeErrors.ErrUnAuthorized()
		}
		secret := os.Getenv("JWT_SECRET")
		return []byte(secret), nil
	})
	if err != nil {
		fmt.Println("failed to parse JWT token:", err)
		return nil, customeErrors.ErrUnAuthorized()
	}
	if !token.Valid {
		fmt.Println("invalid token")
		return nil, customeErrors.ErrUnAuthorized()
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, customeErrors.ErrUnAuthorized()
	}
	return claims, nil
}
