package controllers

import (
	"context"
	"database/sql"
	db "example/web-service-gin/db/sqlc"
	"example/web-service-gin/token"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthController struct {
	db         *db.Queries
	ctx        context.Context
	TokenMaker token.Maker
}

func NewAuthController(db *db.Queries, ctx context.Context, tokenMaker token.Maker) *AuthController {
	return &AuthController{db, ctx, tokenMaker}
}

// Web-Service-gin godoc
// @Summary Login
// @Schemes
// @Description use contact id to get access token
// @Tags auth
// @Accept json
// @Produce json
// @Param   contactId 	path   string  true  "Contact ID"
// @Success 200
// @Router /auth/login/{contactId} [get]
func (ac *AuthController) Login(ctx *gin.Context) {
	contactId := ctx.Param("contactId")
	contact, err := ac.db.GetContactById(ctx, uuid.MustParse(contactId))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "failed", "message": "Failed to retrieve contact with this ID"})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "Failed retrieving contact", "error": err.Error()})
		return
	}

	generatedToken, err := ac.TokenMaker.CreateToken(contact.ContactID.String(), time.Minute*30)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"staus": "failed", "message": fmt.Sprintf("Failed to get token maker, %s", err.Error())})
	}

	ctx.JSON(200, gin.H{"token": generatedToken})
}
