package controllers

//
// import (
// 	"net/http"
//
// 	"github.com/go-chi/render"
// )

//
// func (c *Controllers) GetBrokerageAccount(w http.ResponseWriter, r *http.Request) {
// 	brokerageAccountID := uuid.Must(uuid.Parse(chi.URLParam(r, "brokerage_account_id")))
//
// 	brokerageAccount, err := c.store.GetBrokerageAccount(brokerageAccountID)
//
// 	if err != nil {
// 		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})
//
// 		return
// 	}
//
// 	render.JSON(w, r, brokerageAccount)
// }
//
// func (c *Controllers) ListBrokerageAccounts(w http.ResponseWriter, r *http.Request) {
// 	query := map[string]interface{}{}
// 	userID := r.URL.Query().Get("user_id")
//
// 	if userID != "" {
// 		query["user_id"] = userID
// 	}
//
// 	brokerageAccounts, err := c.store.ListBrokerageAccounts(query)
//
// 	if err != nil {
// 		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})
//
// 		return
// 	}
//
// 	render.JSON(w, r, brokerageAccounts)
// }

// func (c *Controllers) CreateBrokerageAccount(w http.ResponseWriter, r *http.Request) {
// 	c.broker.CreateAccount("test", "test")
//
// 	w.WriteHeader(http.StatusCreated)
// 	render.JSON(w, r, map[string]string{"response": "success"})
// }

//
// func (c *Controllers) ModifyBrokerageAccount(w http.ResponseWriter, r *http.Request) {
// 	var brokerageAccountPayload models.BrokerageAccount
//
// 	brokerageAccountID := uuid.Must(uuid.Parse(chi.URLParam(r, "brokerage_account_id")))
//
// 	errDecode := json.NewDecoder(r.Body).Decode(&brokerageAccountPayload)
// 	if errDecode != nil {
// 		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errDecode})
//
// 		return
// 	}
//
// 	brokerageAccount, err := c.store.ModifyBrokerageAccount(brokerageAccountID, brokerageAccountPayload)
//
// 	if err != nil {
// 		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})
//
// 		return
// 	}
//
// 	render.JSON(w, r, brokerageAccount)
// }
//
// func (c *Controllers) DeleteBrokerageAccount(w http.ResponseWriter, r *http.Request) {
// 	brokerageAccountID := uuid.Must(uuid.Parse(chi.URLParam(r, "brokerage_account_id")))
//
// 	err := c.store.DeleteBrokerageAccount(brokerageAccountID)
//
// 	if err != nil {
// 		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})
//
// 		return
// 	}
//
// 	w.WriteHeader(http.StatusNoContent)
// }
