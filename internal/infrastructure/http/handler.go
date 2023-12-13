package http

import (
	"WBL0/internal/infrastructure/cache"
	"WBL0/internal/infrastructure/postgres"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type Handler struct {
	cache.OrderCache
	postgres.OrderPostgres
}

func NewHandler(oc cache.OrderCache, op postgres.OrderPostgres) *Handler {
	return &Handler{oc, op}
}

func (h *Handler) RunServer() {
	r := gin.Default()
	r.LoadHTMLFiles("internal/template/index.html")
	r.GET("/order", func(c *gin.Context) {
		htmlPath := "index.html"
		c.HTML(http.StatusOK, htmlPath, nil)

	})
	r.POST("/order", func(c *gin.Context) {
		orderUID := c.PostForm("order_uid")
		order, err := h.OrderCache.GetOrder(orderUID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}
		c.JSON(200, order)
		log.Println(order)
	})
	r.Run(":8080")
}
