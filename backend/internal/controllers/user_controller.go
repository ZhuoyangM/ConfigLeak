package controllers

import (
	"github.com/ZhuoyangM/ConfigLeak/internal/store"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserService store.UserService
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (c *UserController) Register(ctx *gin.Context) {
	var req RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	user := store.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	if err := c.UserService.CreateUser(&user); err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to create user"})
		return
	}

	ctx.JSON(201, gin.H{"message": "User created successfully"})
}

func (c *UserController) Login(ctx *gin.Context) {}

// GET /api/user (for testing purposes)
func (c *UserController) GetUserInfo(ctx *gin.Context) {}
