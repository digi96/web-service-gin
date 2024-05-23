package controllers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	db "example/web-service-gin/db/sqlc"
	"example/web-service-gin/rabbitmqconnect"
	"example/web-service-gin/schemas"
	"example/web-service-gin/util"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/vmihailenco/msgpack/v5"
)

type ContactController struct {
	db  *db.Queries
	ctx context.Context
}

var (
	mycache = cache.New(&cache.Options{
		Redis: redis.NewRing(&redis.RingOptions{
			Addrs: map[string]string{
				"server": ":6379",
			},
		}),
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})
)

func NewContactController(db *db.Queries, ctx context.Context) *ContactController {
	return &ContactController{db, ctx}
}

// Web-Service-gin godoc
// @Summary Create a new contact
// @Schemes
// @Description Create a new contact in DB.
// @Tags contacts
// @Accept json
// @Produce json
// @Param   contact 	body   schemas.CreateContact  true  "Contact JSON"
// @Success 200 {object} db.Contact
// @Router /contacts [post]
func (cc *ContactController) CreateContact(ctx *gin.Context) {
	var payload *schemas.CreateContact

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "Failed payload", "error": err.Error()})
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

	contact, err := cc.db.CreateContact(ctx, *args)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "Failed retrieving contact", "error": err.Error()})
		return
	}

	contactJson, err := json.Marshal(contact)
	if err != nil {
		util.FailOnError(err, "Failed to Marshal contact struct")
	}

	rabbit := rabbitmqconnect.RabbitMQ{Body: string(contactJson), QueueName: "defaultqueue"}
	rabbit.Puplish()

	ctx.JSON(http.StatusOK, gin.H{"status": "successfully created contact", "contact": contact})
}

// Web-Service-gin godoc
// @Summary Update a contact
// @Schemes
// @Description Update a contact in DB.
// @Tags contacts
// @Accept json
// @Produce json
// @Param   contactId 	path   string  true  "Contact ID"
// @Param   contact 	body   schemas.UpdateContact  true  "Contact JSON"
// @Success 200 {object} db.Contact
// @Router /contacts/{contactId} [patch]
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

	contact, err := cc.db.UpdateContact(ctx, *args)

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

// Web-Service-gin godoc
// @Summary Show a contact
// @Schemes
// @Description get contact by contact_id
// @Tags contacts
// @Accept json
// @Produce json
// @Param   contactId 	path   string  true  "Contact ID"
// @Success 200 {object} db.Contact
// @Router /contacts/{contactId} [get]
func (cc *ContactController) GetContactById(ctx *gin.Context) {
	contactId := ctx.Param("contactId")

	var wantedCheck []byte
	if err := mycache.Get(ctx, contactId, &wantedCheck); err == nil {
		fmt.Println(wantedCheck)
		var cachedContact db.Contact
		err = msgpack.Unmarshal(wantedCheck, &cachedContact)
		if err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"status": "Failed retrieving contact", "error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"status": "Successfully retrived id(cached)", "contact": cachedContact})
		return
	}

	contact, err := cc.db.GetContactById(ctx, uuid.MustParse(contactId))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "failed", "message": "Failed to retrieve contact with this ID"})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "Failed retrieving contact", "error": err.Error()})
		return
	}

	//ctx := context.TODO()
	//key :=

	b, err := msgpack.Marshal(contact)

	if err != nil {
		fmt.Println("failed to marshal contact object before storing to redis", err.Error())
	}

	if err := mycache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   contactId,
		Value: b,
		TTL:   time.Hour,
	}); err != nil {
		fmt.Println("failed to store object to redis", err.Error())
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "Successfully retrived id", "contact": contact})
}

// Web-Service-gin godoc
// @Summary Get all contacts
// @Schemes
// @Description Get all contacts
// @Tags contacts
// @Produce json
// @Success 200 {object} []db.Contact
// @Router /contacts [get]
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

	contacts, err := cc.db.ListContacts(ctx, *args)
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

	_, err := cc.db.GetContactById(ctx, uuid.MustParse(contactId))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "failed", "message": "Failed to retrieve contact with this ID"})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "Failed retrieving contact", "error": err.Error()})
		return
	}

	err = cc.db.DeleteContact(ctx, uuid.MustParse(contactId))
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "failed", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{"status": "successfuly deleted"})

}
