package controllers

import (
	"net/http"

	"github.com/ZhuoyangM/ConfigLeak/internal/store"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserService *store.UserService
}

func (c *UserController) Register(ctx *gin.Context) {
	var req store.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := c.UserService.CreateUser(&req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func (c *UserController) Login(ctx *gin.Context) {
	var req store.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	token, err := c.UserService.AuthenticateUser(&req)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Login successfully", "token": token})
}

func (c *UserController) GetUserInfo(ctx *gin.Context) {
	userId, ok := ctx.Get("userId")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	user, err := c.UserService.GetUserByID(userId.(uint))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	resp := store.ToGetUserResponse(user)
	ctx.JSON(http.StatusOK, gin.H{"message": "User info retrieved successfully", "user": resp})
}
