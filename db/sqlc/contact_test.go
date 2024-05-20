package db

import (
	"context"
	"database/sql"
	"example/web-service-gin/util"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func CreateRandomContact(t *testing.T) Contact {
	now := time.Now()
	arg := CreateContactParams{
		FirstName:   util.RandomString(5),
		LastName:    util.RandomString(5),
		PhoneNumber: strconv.FormatInt(util.RandomInt(9999999, 99999999), 10),
		Street:      util.RandomAddress(),
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	contact, err := testQueries.CreateContact(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, contact)

	require.Equal(t, arg.PhoneNumber, contact.PhoneNumber)
	require.Equal(t, arg.Street, contact.Street)

	require.NotZero(t, contact.ContactID)
	require.NotZero(t, contact.CreatedAt)

	return contact
}

func TestCreateContact(t *testing.T) {
	CreateRandomContact(t)
}

func TestUpdateStore(t *testing.T) {
	contact1 := CreateRandomContact(t)

	arg := UpdateContactParams{
		ContactID:   contact1.ContactID,
		FirstName:   sql.NullString{String: contact1.FirstName, Valid: contact1.FirstName != ""},
		LastName:    sql.NullString{String: contact1.LastName, Valid: contact1.LastName != ""},
		PhoneNumber: sql.NullString{String: contact1.PhoneNumber, Valid: contact1.PhoneNumber != ""},
		Street:      sql.NullString{String: "no physical address"},
	}

	contact2, err := testQueries.UpdateContact(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, contact2)

	require.Equal(t, contact1.ContactID, contact2.ContactID)
	require.Equal(t, contact1.FirstName, contact2.FirstName)
	require.Equal(t, contact1.LastName, contact2.LastName)
	require.Equal(t, contact1.Street, contact2.Street)
	require.WithinDuration(t, contact1.CreatedAt, contact2.CreatedAt, time.Second)

}

func TestDeleteStore(t *testing.T) {
	contact1 := CreateRandomContact(t)

	err := testQueries.DeleteContact(context.Background(), contact1.ContactID)
	require.NoError(t, err)

	contact2, err := testQueries.GetContactById(context.Background(), contact1.ContactID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, contact2)

}

func TestListContact(t *testing.T) {
	arg := ListContactsParams{
		Limit:  5,
		Offset: 5,
	}

	createdContacts := []Contact{}

	for i := 0; i < 10; i++ {
		tempContact := CreateRandomContact(t)
		createdContacts = append(createdContacts, tempContact)
	}

	contacts, err := testQueries.ListContacts(context.Background(), arg)

	if assert.Nil(t, err) {
		require.NotEmpty(t, contacts)
	}

	require.NoError(t, err)
	require.NotEmpty(t, contacts)
	require.Len(t, contacts, 5)

	for _, contact := range contacts {
		require.NotEmpty(t, contact)
	}

	for _, contact := range createdContacts {
		err := testQueries.DeleteContact(context.Background(), contact.ContactID)
		require.NoError(t, err)
	}
}
