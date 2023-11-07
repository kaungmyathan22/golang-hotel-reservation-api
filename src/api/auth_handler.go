package api

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	repository "github.com/kaungmyathan22/golang-hotel-reservation/src/repositories"
	"github.com/kaungmyathan22/golang-hotel-reservation/src/types"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userStore repository.UserStore
}

func NewAuthHandler(userStore repository.UserStore) *AuthHandler {
	return &AuthHandler{
		userStore: userStore,
	}
}

func (h *AuthHandler) HandleLogin(c *fiber.Ctx) error {
	var payload Loginparams
	if err := c.BodyParser(&payload); err != nil {
		return err
	}
	user, err := h.userStore.GetUserByEmail(c.Context(), payload.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("invalid credentials")
		}
		return err
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		return fmt.Errorf("invalid username / password.")
	}
	token := createTokenFromUser(user)
	return c.JSON(map[string]any{
		"user":  user,
		"token": token,
	})
}

func createTokenFromUser(user *types.User) string {
	now := time.Now()
	validTill := now.Add(time.Hour * 4).Unix()
	claims := jwt.MapClaims{
		"userID":  user.ID,
		"email":   user.Email,
		"expires": validTill,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Println(err, secret)
		fmt.Println("Failed to sign token with secret")
	}
	return tokenStr
}
