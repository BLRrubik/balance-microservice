package server

import (
	"balance-microservice/internal/app/model"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

func (s *server) configureRouter() {
	//health check endpoint
	s.router.GET("/health", s.healthCheck)

	//users endpoints
	s.router.GET("/users/:userId", s.getBalanceOfUser)
	s.router.POST("/users/deposit", s.replenishmentAccount)

	//service endpoint
	s.router.GET("/services/:serviceId", s.getServiceById)

	//bills endpoint
	s.router.POST("/bills", s.reserveFunds)
	s.router.PATCH("/bill/:billId/approve", s.approveReservation)
	s.router.PATCH("/bill/:billId/reject", s.rejectReservation)

	//accounting
	s.router.GET("/accounting", s.getAccounting)
	s.router.GET("/accounting/csv", s.exportCSV)

	//transactions
	s.router.GET("/transactions/:userId", s.getTransactionsOfUser)
}

// health
func (s *server) healthCheck(c *gin.Context) {
	type response struct {
		Status string `json:"status"`
		Time   string `json:"time"`
	}

	res := &response{
		Time:   time.Now().Format("02-01-2006 15:04"),
		Status: "UP",
	}

	c.JSON(200, res)
}

// users methods
func (s *server) replenishmentAccount(c *gin.Context) {
	var req model.UserDepositRequest

	if err := c.ShouldBind(&req); err != nil {
		sendError(c, 400, err.Error())
		return
	}

	u, e := s.store.UserRepository().ReplenishmentAccount(&req)
	if e != nil {
		sendError(c, e.StatusCode, e.Message)
		return
	}

	c.JSON(200, u.ToDto())
}

func (s *server) getBalanceOfUser(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		sendError(c, 400, err.Error())
		return
	}

	u, e := s.store.UserRepository().GetBalanceOfUser(userId)
	if e != nil {
		sendError(c, e.StatusCode, e.Message)
		return
	}

	c.JSON(200, u.ToDto())
}

// accounting methods
func (s *server) getAccounting(c *gin.Context) {
	records, err := s.store.AccountingRepository().FindAll()
	if err != nil {
		sendError(c, err.StatusCode, err.Message)
	}

	c.JSON(200, records)
}

func (s *server) exportCSV(c *gin.Context) {
	date := c.Query("date")

	filepath, err := s.store.AccountingRepository().ExportCSV(date)
	if err != nil {
		sendError(c, err.StatusCode, err.Message)
		return
	}
	c.Header("Content-Type", "application/octet-stream")

	c.File(*filepath)
}

// service methods
func (s *server) getServiceById(c *gin.Context) {
	serviceId, err := strconv.Atoi(c.Param("serviceId"))
	if err != nil {
		sendError(c, 400, err.Error())
		return
	}

	service, e := s.store.ServiceRepository().FindById(serviceId)
	if e != nil {
		sendError(c, e.StatusCode, e.Message)
		return
	}

	c.JSON(200, service)
}

// bill methods
func (s *server) reserveFunds(c *gin.Context) {
	var req model.BillCreateRequest

	if err := c.ShouldBind(&req); err != nil {
		sendError(c, 400, err.Error())
		return
	}

	bill, e := s.store.BillRepository().ReservedFunds(&req)
	if e != nil {
		sendError(c, e.StatusCode, e.Message)
		return
	}

	c.JSON(200, bill.ToDto())
}

func (s *server) approveReservation(c *gin.Context) {
	billId, err := strconv.Atoi(c.Param("billId"))
	if err != nil {
		sendError(c, 400, err.Error())
		return
	}

	bill, e := s.store.BillRepository().ApproveReservation(billId)
	if e != nil {
		sendError(c, e.StatusCode, e.Message)
		return
	}

	c.JSON(200, bill.ToDto())
}

func (s *server) rejectReservation(c *gin.Context) {
	billId, err := strconv.Atoi(c.Param("billId"))
	if err != nil {
		sendError(c, 400, err.Error())
		return
	}

	bill, e := s.store.BillRepository().RejectReservation(billId)
	if e != nil {
		sendError(c, e.StatusCode, e.Message)
		return
	}

	c.JSON(200, bill.ToDto())
}

// methods for transaction
func (s *server) getTransactionsOfUser(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		sendError(c, 400, err.Error())
		return
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "0"))
	if err != nil {
		sendError(c, 400, err.Error())
		return
	}

	size, err := strconv.Atoi(c.DefaultQuery("size", "10"))
	if err != nil {
		sendError(c, 400, err.Error())
		return
	}

	sort := c.DefaultQuery("sort", "date")

	dir := c.DefaultQuery("dir", "ASC")

	transactions, e := s.store.TransactionRepository().GetTransactionsOfUser(userId, page, size, sort, dir)
	if e != nil {
		sendError(c, e.StatusCode, e.Message)
		return
	}

	c.JSON(200, transactions)
}

// method for send error
func sendError(c *gin.Context, code int, err string) {
	c.JSON(code, ErrorType{
		Message: err,
	})
}
