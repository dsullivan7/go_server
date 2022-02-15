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

  var brokerageAccountRes []map[string]interface{}

  for _, brokerageAccount := range brokerageAccounts {
    var brkrgAccntRes map[string]interface{}

    brkrgAccntJson, errDecode := json.Marshal(brokerageAccount)

    if (errDecode != nil) {
      c.utils.HandleError(w, r, errors.HTTPServerError{Err: errDecode})

      return
    }

    errEncode := json.Unmarshal(brkrgAccntJson, &brkrgAccntRes)

    if (errEncode != nil) {
      c.utils.HandleError(w, r, errors.HTTPServerError{Err: errEncode})

      return
    }

    brokerRes, brokerErr := c.broker.GetAccount(*brokerageAccount.AlpacaAccountID)

    if (brokerErr != nil) {
      c.utils.HandleError(w, r, errors.HTTPServerError{Err: brokerErr})

      return
    }

    brkrgAccntRes["cash"] = brokerRes.Cash

    brokerageAccountRes = append(brokerageAccountRes, brkrgAccntRes)
  }

	render.JSON(w, r, brokerageAccountRes)
}

func (c *Controllers) CreateBrokerageAccount(w http.ResponseWriter, r *http.Request) {
  var brokerageAccountReq map[string]interface{}

	errDecode := json.NewDecoder(r.Body).Decode(&brokerageAccountReq)
	if errDecode != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errDecode})

		return
	}

  alpacaAccountID, errBroker := c.broker.CreateAccount(
    brokerageAccountReq["first_name"].(string),
  	brokerageAccountReq["last_name"].(string),
  	brokerageAccountReq["date_of_birth"].(string),
  	brokerageAccountReq["tax_id"].(string),
  	brokerageAccountReq["email_address"].(string),
  	brokerageAccountReq["phone_number"].(string),
  	brokerageAccountReq["street_address"].(string),
  	brokerageAccountReq["city"].(string),
  	brokerageAccountReq["state"].(string),
  	brokerageAccountReq["postal_code"].(string),
  	brokerageAccountReq["funding_source"].(string),
  	r.RemoteAddr,
  )

  if errBroker != nil {
    c.utils.HandleError(w, r, errors.HTTPUserError{Err: errBroker})

    return
  }

  userID := uuid.Must(uuid.Parse(brokerageAccountReq["user_id"].(string)))

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
