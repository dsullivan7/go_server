package controllers


import (
  "encoding/json"
	"go_server/internal/errors"
	"go_server/internal/models"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)


func (c *Controllers) GetBrokerageAccount(w http.ResponseWriter, r *http.Request) {
	brokerageAccountID := uuid.Must(uuid.Parse(chi.URLParam(r, "brokerage_account_id")))

	brokerageAccount, err := c.store.GetBrokerageAccount(brokerageAccountID)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.JSON(w, r, brokerageAccount)
}

func (c *Controllers) ListBrokerageAccounts(w http.ResponseWriter, r *http.Request) {
	query := map[string]interface{}{}
	userID := r.URL.Query().Get("user_id")

	if userID != "" {
		query["user_id"] = userID
	}

	brokerageAccounts, err := c.store.ListBrokerageAccounts(query)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.JSON(w, r, brokerageAccounts)
}

func (c *Controllers) CreateBrokerageAccount(w http.ResponseWriter, r *http.Request) {
  var brokerageAccountReq map[string]interface{}

	errDecode := json.NewDecoder(r.Body).Decode(&brokerageAccountReq)
	if errDecode != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errDecode})

		return
	}

  alpacaAccountID, errBroker := c.broker.CreateAccount(
    brokerageAccountReq["givenName"].(string),
  	brokerageAccountReq["familyName"].(string),
  	brokerageAccountReq["dateOfBirth"].(string),
  	brokerageAccountReq["taxID"].(string),
  	brokerageAccountReq["emailAddress"].(string),
  	brokerageAccountReq["phoneNumber"].(string),
  	brokerageAccountReq["streetAddress"].(string),
  	brokerageAccountReq["city"].(string),
  	brokerageAccountReq["state"].(string),
  	brokerageAccountReq["postalCode"].(string),
  	brokerageAccountReq["fundingSource"].(string),
  	brokerageAccountReq["ipAddress"].(string),
  )

  if errBroker != nil {
    c.utils.HandleError(w, r, errors.HTTPUserError{Err: errBroker})

    return
  }

  userID := uuid.Must(uuid.Parse(brokerageAccountReq["userID"].(string)))

  brokerageAccountPayload := models.BrokerageAccount{
    UserID: &userID,
    AlpacaAccountID: &alpacaAccountID,
  }

	brokerageAccount, err := c.store.CreateBrokerageAccount(brokerageAccountPayload)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, brokerageAccount)
}


func (c *Controllers) ModifyBrokerageAccount(w http.ResponseWriter, r *http.Request) {
	var brokerageAccountPayload models.BrokerageAccount

	brokerageAccountID := uuid.Must(uuid.Parse(chi.URLParam(r, "brokerage_account_id")))

	errDecode := json.NewDecoder(r.Body).Decode(&brokerageAccountPayload)
	if errDecode != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errDecode})

		return
	}

	brokerageAccount, err := c.store.ModifyBrokerageAccount(brokerageAccountID, brokerageAccountPayload)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.JSON(w, r, brokerageAccount)
}

func (c *Controllers) DeleteBrokerageAccount(w http.ResponseWriter, r *http.Request) {
	brokerageAccountID := uuid.Must(uuid.Parse(chi.URLParam(r, "brokerage_account_id")))

	err := c.store.DeleteBrokerageAccount(brokerageAccountID)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

  render.NoContent(w, r)
}
