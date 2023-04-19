package rest

import (
	"astro/dto"
	os2 "astro/os"
	"astro/utils"
	"astro/vpn"
	setup "astro/vpn/install"
	"astro/vpn/types"
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
)

type (
	RestAPI struct {
		DB       *badger.DB
		ConfigDB *badger.DB
		Router   fiber.Router
	}
)

func NewRestAPI(db *badger.DB, config *badger.DB, router fiber.Router) *RestAPI {
	return &RestAPI{
		DB:       db,
		ConfigDB: config,
		Router:   router,
	}
}

func (r *RestAPI) Handlers() {
	api := r.Router.Group("/vpn")

	api.Post("/create", r.createUser)
	api.Get("/list", r.listUsers)
	api.Delete("/user", r.deleteUser)
	api.Post("/install", r.install)
	api.Get("/uninstall", r.uninstall)
	api.Get("/download", r.downloadFile)
	api.Get("/config", r.config)
	api.Get("/port", r.portAvailable)
}

func (r *RestAPI) install(c *fiber.Ctx) error {

	response := dto.Response{
		Status: fiber.StatusInternalServerError,
	}

	err := setup.InstallVPN()
	if err != nil {
		response.Error = err.Error()
		return c.Status(response.Status).JSON(response)
	}

	response.Status = fiber.StatusOK
	response.Data = true

	return c.JSON(response)
}

func (r *RestAPI) config(c *fiber.Ctx) error {

	response := dto.Response{
		Status: fiber.StatusInternalServerError,
	}
	config, err := setup.GetConfig()
	if err != nil {
		response.Error = err.Error()
		return c.Status(response.Status).JSON(response)
	}

	response.Status = fiber.StatusOK
	response.Data = config

	return c.JSON(response)
}

func (r *RestAPI) portAvailable(c *fiber.Ctx) error {
	port := c.Query("available")
	response := dto.Response{
		Status: fiber.StatusInternalServerError,
	}

	if port == "" {
		response.Error = "Port is required"
		return c.Status(response.Status).JSON(response)
	}
	response.Status = fiber.StatusOK
	response.Data = os2.IsPortAvailable(port)
	return c.JSON(response)
}

func (r *RestAPI) uninstall(c *fiber.Ctx) error {

	response := dto.Response{
		Status: fiber.StatusInternalServerError,
	}

	err := setup.UninstallVPN()
	if err != nil {
		response.Error = err.Error()
		return c.Status(response.Status).JSON(response)
	}

	response.Status = fiber.StatusOK
	response.Data = true

	return c.JSON(response)
}

func (r *RestAPI) createUser(c *fiber.Ctx) error {

	response := dto.Response{
		Status: fiber.StatusOK,
		Data:   false,
	}

	var body CreateClientRequest
	err := c.BodyParser(&body)
	if err != nil {
		response.Status = fiber.StatusInternalServerError
		response.Error = err.Error()
		return c.Status(response.Status).JSON(response)
	}

	errors := utils.ValidateStruct(&body)
	if errors != nil {
		response.Status = fiber.StatusBadRequest
		jsonErrors, err := json.Marshal(errors)
		if err != nil {
			response.Error = err.Error()
			return c.Status(response.Status).JSON(response)
		}
		response.Error = string(jsonErrors)
		return c.Status(response.Status).JSON(response)
	}

	user, err := vpn.CreateUser(r.DB, body.Username, body.Password)
	if err != nil {
		response.Status = fiber.StatusBadRequest
		response.Error = err.Error()
		return c.Status(response.Status).JSON(response)
	}
	response.Data = user
	return c.JSON(response)
}

func (r *RestAPI) listUsers(c *fiber.Ctx) error {
	clients, err := vpn.GetAllUsers(r.DB)

	response := dto.Response{
		Status: fiber.StatusOK,
		Data:   []types.User{},
	}

	if err != nil {
		response.Status = fiber.StatusBadRequest
		response.Error = err.Error()
		return c.Status(response.Status).JSON(response)
	}

	response.Data = clients

	return c.JSON(response)
}

func (r *RestAPI) deleteUser(c *fiber.Ctx) error {
	username := c.Query("username")

	response := dto.Response{
		Status: fiber.StatusOK,
		Data:   false,
	}

	if username == "" {
		response.Status = fiber.StatusBadRequest
		response.Error = "Username is required"
		return c.Status(response.Status).JSON(response)
	}

	err := vpn.DeleteUser(r.DB, username)
	if err != nil {
		response.Status = fiber.StatusBadRequest
		response.Error = err.Error()
		return c.Status(response.Status).JSON(response)
	}

	response.Data = true

	return c.JSON(response)
}

func (r *RestAPI) downloadFile(c *fiber.Ctx) error {
	username := c.Query("username")

	data, err := vpn.GeneratePathFileOvpn(username)
	if err != nil {
		response := dto.Response{
			Status: fiber.StatusInternalServerError,
			Error:  err.Error(),
		}
		return c.Status(response.Status).JSON(response)
	}

	c.Response().Header.Set("Content-Type", "application/zip")
	c.Response().Header.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.zip", username))

	return c.Send([]byte(data))
}
