package delivery

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/obh/chargebackapp/models"
	"github.com/obh/chargebackapp/utils"
	"github.com/obh/chargebackapp/workflows"

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
	e.POST("/notify", cb.notifyMerchant)
	e.POST("/response/:chargebackId", cb.handleMerchantResponse)
}

func (h *ChargebackHandler) addChargeback(ctx echo.Context) error {
	req := &models.ChargebackRequest{}
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}
	var payment models.Payment
	fmt.Println("finding payment ... :", req.PaymentId)
	if err := h.db.First(&payment, req.PaymentId).Error; err != nil {
		fmt.Println(err)
		return ctx.JSON(http.StatusNotFound, errors.New("payment not found"))
	}
	fmt.Println("payment found: ", payment.ID)
	chargeback := &models.Chargeback{ChargebackRequest: *req}
	h.db.Create(chargeback)
	fmt.Println("chargeback result: ", chargeback)

	input := workflows.ChargebackWFInput{
		Chargeback: *chargeback,
		Payment:    payment,
		Merchant:   utils.GetMerchant(),
	}
	_, err := h.temporalClient.ExecuteWorkflow(
		ctx.Request().Context(),
		client.StartWorkflowOptions{
			TaskQueue: workflows.TaskQueue,
			ID:        workflows.ChargebackWorkflowId(int(input.Chargeback.ID)),
		},
		workflows.ChargebackProcess,
		&input,
	)

	if err != nil {
		log.Printf("failed to start workflow: %v", err)
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusCreated, *chargeback)
}

func (h *ChargebackHandler) handleMerchantResponse(ctx echo.Context) error {
	cId, err := strconv.ParseUint(ctx.Param("chargebackId"), 10, 32)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errors.New("invalid chargeback id"))
	}
	response := &models.MerchantResponse{}
	if err := ctx.Bind(response); err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}
	err = h.signal(cId)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, "something went wrongt")
	}
	return ctx.JSON(http.StatusOK, "ok")
}

func (h *ChargebackHandler) signal(chargebackId uint64) error {
	signal := workflows.MerchantSubmissionSignal{
		MerchantResponded: true,
		RespondedAt:       time.Now(),
		Message:           "response from merchant",
		Proof:             "this is the proof",
	}
	workflowId := fmt.Sprintf("merchant_response:%d", chargebackId)
	return h.temporalClient.SignalWorkflow(context.Background(), workflowId, "",
		workflows.MerchantSubmissionSignalName, signal)

}

func (h *ChargebackHandler) notifyMerchant(ctx echo.Context) error {
	notifyReq := &models.ChargebackNotifyMerchant{}
	if err := ctx.Bind(notifyReq); err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}
	fmt.Println("Found request in /notify: ", notifyReq)
	var payment models.Payment
	var chargeback models.Chargeback
	if err := h.db.First(&payment, notifyReq.PaymentId).Error; err != nil {
		return ctx.JSON(http.StatusNotFound, errors.New("payment not found"))
	}
	if err := h.db.First(&chargeback, uint(notifyReq.ChargebackID)).Error; err != nil {
		return ctx.JSON(http.StatusNotFound, errors.New("chargeback not found"))
	}
	f, _ := os.LookupEnv("TEMPLATE_FOLDER")
	b, err := os.ReadFile(f + "/updateMerchant.html")
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, errors.New("email template not found"))
	}
	t, err := template.New("email").Parse(string(b))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, errors.New("email template not parsed"))
	}
	err = utils.SendMail("rohit@cashfree.com", "rohit@cashfree.com", "Chargeback received", t, t, chargeback)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, errors.New("email not sent"))
	}
	return ctx.JSON(http.StatusOK, "sent")
}
