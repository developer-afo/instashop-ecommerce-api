package handler

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/developer-afo/instashop-ecommerce-api/payload/response"
	"github.com/developer-afo/instashop-ecommerce-api/repository"
)

func GetUserId(c *fiber.Ctx) uuid.UUID {
	userId := c.Locals("userId").(uuid.UUID)

	return userId
}

func Index(c *fiber.Ctx) error {

	var resp response.Response

	var about struct {
		Name    string `json:"name"`
		Version string `json:"version"`
		Author  string `json:"author"`
	}

	about.Name = "Instashop API"
	about.Version = "0.0.1"
	about.Author = "Afolabi Salawu"

	resp.Status = http.StatusOK
	resp.Message = http.StatusText(http.StatusOK)
	resp.Data = map[string]interface{}{"about": about}

	return c.JSON(resp)
}

func NotFound(c *fiber.Ctx) error {
	var resp response.Response

	resp.Status = http.StatusNotFound
	resp.Message = "Route not found"

	return c.Status(http.StatusNotFound).JSON(resp)
}

func GeneratePageable(context *fiber.Ctx) (pageable repository.Pageable) {

	pageable.Page = 1
	pageable.Size = 20
	pageable.SortBy = "created_at"
	pageable.SortDirection = "desc"
	pageable.Search = ""

	size, err := strconv.Atoi(context.Query("size", "0"))
	if (size > 0) && err == nil {
		pageable.Size = size
	}

	page, err := strconv.Atoi(context.Query("page", "1"))
	if (page > 0) && err == nil {
		pageable.Page = page
	}

	orderBy := context.Query("sort_by", "")
	if orderBy != "" {
		pageable.SortBy = orderBy
	}

	sortDir := context.Query("sort_dir", "")
	if sortDir != "" {
		pageable.SortBy = sortDir
	}

	search := context.Query("search", "")
	if search != "" {
		pageable.Search = search
	}

	return pageable
}
