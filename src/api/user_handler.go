package api

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	db "github.com/kaungmyathan22/golang-hotel-reservation/src/repositories"
	"github.com/kaungmyathan22/golang-hotel-reservation/src/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}

type Loginparams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := h.userStore.GetUsers(c.Context())
	if err != nil {
		log.Fatal(err)
	}
	return c.JSON(users)
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := h.userStore.GetUserByID(c.Context(), id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("User with given id not found.")
		}
		return err
	}
	return c.JSON(user)
}

func (h *UserHandler) HandleCreateUsers(c *fiber.Ctx) error {
	var payload types.CreateUserPayload
	if err := c.BodyParser(&payload); err != nil {
		return err
	}
	if errors := payload.Validate(); len(errors) > 0 {
		return c.JSON(errors)
	}
	parsedPayload, err := types.NewUserFromParams(payload)
	if err != nil {
		return err
	}
	if alreadyExistedUser, _ := h.userStore.GetUserByEmail(c.Context(), payload.Email); alreadyExistedUser != nil {
		return fmt.Errorf("user with that email already existed.")
	}
	user, err := h.userStore.CreateUser(c.Context(), parsedPayload)
	if err != nil {
		return err
	}
	return c.JSON(user)
}

func (h *UserHandler) HandleDeleteUsers(c *fiber.Ctx) error {
	userID := c.Params("id")
	err := h.userStore.DeleteUserByID(c.Context(), userID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("User with given id not found.")
		}
		return err
	}
	return c.JSON(map[string]string{"message": "Successfully deleted user with " + userID})
}

func (h *UserHandler) HandleUpdateUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	// var payload bson.M
	var payload types.UpdateUserPayload
	if err := c.BodyParser(&payload); err != nil {
		return err
	}
	objectId, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": objectId}
	err = h.userStore.UpdateUser(c.Context(), filter, payload)
	if err != nil {
		return err
	}
	return c.JSON(map[string]string{"message": "Successfully updated user. "})
}
