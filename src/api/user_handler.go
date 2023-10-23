package api

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/kaungmyathan22/golang-hotel-reservation/src/db"
	"github.com/kaungmyathan22/golang-hotel-reservation/src/types"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
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

	user, err := h.userStore.CreateUser(c.Context(), parsedPayload)
	if err != nil {
		return err
	}
	return c.JSON(user)
}
