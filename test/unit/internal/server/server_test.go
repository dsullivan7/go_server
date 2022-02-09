package server_test

import (
	"testing"

	"go_server/test/utils"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

func TestServer(tParent *testing.T) {
	tParent.Parallel()

	testServer, err := utils.NewTestServer()
	assert.Nil(tParent, err)

	tParent.Run("Test Init", func(t *testing.T) {
		t.Parallel()

		tctx := chi.NewRouteContext()
		assert.True(t, testServer.Router.Match(tctx, "GET", "/api/users"))
	})
}
