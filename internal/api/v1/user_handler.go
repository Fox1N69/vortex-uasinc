package v1

import (
	"net/http"
	"strconv"
	"test-task/common/http/response"
	"test-task/internal/models"
	"test-task/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	GetAllUsers(c *gin.Context)
	UsersWithFiltersAndPagination(c *gin.Context)
	CreateUser(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
}

type userHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) UserHandler {
	return &userHandler{service: service}
}

// @Summary Get all users
// @Description Retrieves a list of all users.
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} models.User "users"
// @Failure 500 {object} models.Response "error"
// @Router /api/user [get]
func (h *userHandler) GetAllUsers(c *gin.Context) {
	users, err := h.service.Users()
	if err != nil {
		response.New(c).Error(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, users)
}

// @Summary Get users with filters and pagination
// @Description Retrieves users with filters and pagination parameters.
// @Tags users
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param pageSize query int false "Page size"
// @Success 200 {array} models.User "users"
// @Failure 500 {object} models.Response "error"
// @Router /api/user/filters [get]
func (h *userHandler) UsersWithFiltersAndPagination(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	page, _ := strconv.Atoi(pageStr)
	pageSizeStr := c.DefaultQuery("pageSize", "10")
	pageSize, _ := strconv.Atoi(pageSizeStr)

	users, err := h.service.UsersWithFiltersAndPagination(models.UserFilters{}, models.Pagination{Page: page, PageSize: pageSize})
	if err != nil {
		response.New(c).Error(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, users)
}

// @Summary Create a new user
// @Description Creates a new user based on data received from the request.
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.User true "User object to be created"
// @Success 201 {object} models.User "createdUser"
// @Failure 400 {object} models.Response "error"
// @Failure 500 {object} models.Response "error"
// @Router /api/user [post]
func (h *userHandler) CreateUser(c *gin.Context) {
	var input struct {
		PassportNumber string `json:"passportNumber"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	createdUser, err := h.service.CreateUser(input.PassportNumber)
	if err != nil {
		response.New(c).Error(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "create user success", "user": createdUser})
}

// @Summary Update a user
// @Description Updates user data based on data received from the request.
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body models.User true "Updated user object"
// @Success 200 {object} models.User "updatedUser"
// @Failure 400 {object} models.Response "error"
// @Failure 500 {object} models.Response "error"
// @Router /api/user/{id} [patch]
func (h *userHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	err := h.service.UpdateUser(uint(id), &user)
	if err != nil {
		response.New(c).Error(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// @Summary Delete a user
// @Description Deletes a user based on the ID received from the request.
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.SuccessResponse "message"
// @Failure 500 {object} models.Response "error"
// @Router /api/user/{id} [delete]
func (h *userHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	err := h.service.DeleteUser(uint(id))
	if err != nil {
		response.New(c).Error(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
