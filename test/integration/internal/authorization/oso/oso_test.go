package alpaca_test

import (
  "testing"
	goServerOso "go_server/internal/authorization/oso"
	"go_server/internal/models"

	"github.com/google/uuid"
  "github.com/osohq/go-oso"

	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
  t.Parallel()

  o, errOso := oso.NewOso()

  assert.Nil(t, errOso)

  osoAuthorization := goServerOso.NewAuthorization(o)

  errInit := osoAuthorization.Init()
  assert.Nil(t, errInit)

  userID1 := uuid.New()
  userID2 := uuid.New()

  user1 := models.User{ UserID: userID1 }
  user2 := models.User{ UserID: userID2 }
  user3 := models.User{ UserID: userID1 }

  errValidRead := osoAuthorization.Authorize(user1, "read", user3)
  assert.Nil(t, errValidRead)

  errValidModify := osoAuthorization.Authorize(user1, "modify", user3)
  assert.Nil(t, errValidModify)

  errValidDelete := osoAuthorization.Authorize(user1, "delete", user3)
  assert.Nil(t, errValidDelete)

  errInvalidRead := osoAuthorization.Authorize(user1, "read", user2)
  assert.NotNil(t, errInvalidRead)
}
