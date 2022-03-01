package alpaca_test

import (
  "testing"
	goServerOso "go_server/internal/authorization/oso"
	"go_server/internal/models"

	"github.com/google/uuid"
  "github.com/osohq/go-oso"

	"github.com/stretchr/testify/assert"
)

func TestOso(tParent *testing.T) {
  tParent.Parallel()

  o, errOso := oso.NewOso()

  assert.Nil(tParent, errOso)

  osoAuthorization := goServerOso.NewAuthorization(o)

  errInit := osoAuthorization.Init()
  assert.Nil(tParent, errInit)

  tParent.Run("User", func(t *testing.T) {
    userID1 := uuid.New()
    userID2 := uuid.New()

    user1 := models.User{ UserID: userID1 }
    user2 := models.User{ UserID: userID2 }
    user3 := models.User{ UserID: userID1 }

    errValidRead := osoAuthorization.Authorize(user1, "read", user3)
    assert.Nil(t, errValidRead)

    // errValidModify := osoAuthorization.Authorize(user1, "modify", user3)
    // assert.Nil(t, errValidModify)
    //
    // errValidCreate := osoAuthorization.Authorize(user1, "create", user3)
    // assert.Nil(t, errValidCreate)
    //
    // errValidDelete := osoAuthorization.Authorize(user1, "delete", user3)
    // assert.Nil(t, errValidDelete)

    errInvalidRead := osoAuthorization.Authorize(user1, "read", user2)
    assert.NotNil(t, errInvalidRead)
  })
}
