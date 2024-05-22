package controllers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	db "example/web-service-gin/db/sqlc"
	"example/web-service-gin/schemas"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ContactController struct {
	db      *sql.DB
	queries *db.Queries
	ctx     context.Context
}

func NewContactController(db *sql.DB, queries *db.Queries, ctx context.Context) *ContactController {
	return &ContactController{db, queries, ctx}
}

// Create contact  handler
func (cc *ContactController) CreateContact(ctx *gin.Context) {
	var payload *schemas.CreateContact

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "Failed payload", "error": err.Error()})
		return
	}

	//check if fisrt name already exist
	tx, err := cc.db.Begin()
	if err != nil {
		panic(err)
	}
	defer tx.Rollback()

	qtx := cc.queries.WithTx(tx)

	checkEmptyContact, err := qtx.GetContactById(ctx, uuid.New())
	if err != nil {
		fmt.Println("expected logic")
	} else {
		fmt.Println("unexpected logic")
		fmt.Println(checkEmptyContact)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "Failed to create contact", "error": err.Error()})
		return
	}

	now := time.Now()
	args := &db.CreateContactParams{
		FirstName:   payload.FirstName,
		LastName:    payload.LastName,
		PhoneNumber: payload.PhoneNumber,
		Street:      payload.Street,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	contact, err := qtx.CreateContact(ctx, *args)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "Failed retrieving contact", "error": err.Error()})
		return
	}

	tx.Commit()
	ctx.JSON(http.StatusOK, gin.H{"status": "successfully created contact", "contact": contact})
}

// Update contact handler
func (cc *ContactController) UpdateContact(ctx *gin.Context) {
	var payload *schemas.UpdateContact
	contactId := ctx.Param("contactId")

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "Failed payload", "error": err.Error()})
		return
	}

	now := time.Now()
	args := &db.UpdateContactParams{
		ContactID:   uuid.MustParse(contactId),
		FirstName:   sql.NullString{String: payload.FirstName, Valid: payload.FirstName != ""},
		LastName:    sql.NullString{String: payload.LastName, Valid: payload.LastName != ""},
		PhoneNumber: sql.NullString{String: payload.PhoneNumber, Valid: payload.PhoneNumber != ""},
		Street:      sql.NullString{String: payload.Street, Valid: payload.Street != ""},
		UpdatedAt:   sql.NullTime{Time: now, Valid: true},
	}

	contact, err := cc.queries.UpdateContact(ctx, *args)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "failed", "message": "Failed to retrieve contact with this ID"})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "Failed retrieving contact", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "successfully updated contact", "contact": contact})
}

// Get a single handler
func (cc *ContactController) GetContactById(ctx *gin.Context) {
	contactId := ctx.Param("contactId")

	contact, err := cc.queries.GetContactById(ctx, uuid.MustParse(contactId))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "failed", "message": "Failed to retrieve contact with this ID"})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "Failed retrieving contact", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "Successfully retrived id", "contact": contact})
}

// Retrieve all records handlers
func (cc *ContactController) GetAllContacts(ctx *gin.Context) {
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "10")

	reqPageID, _ := strconv.Atoi(page)
	reqLimit, _ := strconv.Atoi(limit)
	offset := (reqPageID - 1) * reqLimit

	args := &db.ListContactsParams{
		Limit:  int32(reqLimit),
		Offset: int32(offset),
	}

	contacts, err := cc.queries.ListContacts(ctx, *args)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "Failed to retrieve contacts", "error": err.Error()})
		return
	}

	if contacts == nil {
		contacts = []db.Contact{}
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "Successfully retrieved all contacts", "size": len(contacts), "contacts": contacts})
}

// Deleting contact handlers
func (cc *ContactController) DeleteContactById(ctx *gin.Context) {
	contactId := ctx.Param("contactId")

	_, err := cc.queries.GetContactById(ctx, uuid.MustParse(contactId))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "failed", "message": "Failed to retrieve contact with this ID"})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "Failed retrieving contact", "error": err.Error()})
		return
	}

	err = cc.queries.DeleteContact(ctx, uuid.MustParse(contactId))
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "failed", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{"status": "successfuly deleted"})

}
