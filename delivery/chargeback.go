package delivery

import (
	"chargebackapp/models"
	"chargebackapp/workflows"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.temporal.io/sdk/client"
	"gorm.io/gorm"
)

type Tabler interface {
	TableName() string
}

type ChargebackHandler struct {
	db             *gorm.DB
	temporalClient client.Client
}

func AddChargebackHandler(e *echo.Echo, db *gorm.DB, c client.Client) {
	cb := &ChargebackHandler{db: db, temporalClient: c}
	e.PUT("/chargeback", cb.addChargeback)
}

func (h *ChargebackHandler) addChargeback(ctx echo.Context) error {
	chgBack := &models.ChargebackRequest{}
	if err := ctx.Bind(chgBack); err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}
	result := h.db.Create(chgBack)
	fmt.Println("chargeback result: ", result)
	//start the workflow for temporal
	//respond back
	input := &workflows.ChargebackInput{
		Chargeback: *chgBack,
		Customer:   models.Customer{Name: "rohit", Email: "rohit+customer@cgocashfree.com", Phone: "9000900001"},
		Payment:    models.Payment{Id: 10, Amount: 100, Reference: "reference"},
	}
	_, err := h.temporalClient.ExecuteWorkflow(
		ctx.Request().Context(),
		client.StartWorkflowOptions{
			TaskQueue: workflows.TaskQueue,
			ID:        workflows.ChargebackWorkflowId(int(input.Chargeback.ID)),
			// SearchAttributes: map[string]interface{}{
			// 	"merchantEmail": email,
			// },
		},
		workflows.ChargebackProcess,
		input,
	)

	if err != nil {
		log.Printf("failed to start workflow: %v", err)
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusCreated, chgBack)
}
