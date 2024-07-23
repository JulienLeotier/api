package utils

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func TransactionMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tx := db.Begin()
		if tx.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
			c.Abort()
			return
		}

		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
				log.Println("Transaction rolled back due to panic")
				panic(r)
			}
		}()

		c.Set("tx", tx)
		c.Next()

		if len(c.Errors) > 0 {
			tx.Rollback()
			log.Println("Transaction rolled back due to errors")
			return
		}

		if err := tx.Commit().Error; err != nil {
			log.Println("Transaction commit failed:", err)
			c.Abort()
			tx.Rollback()
		}
	}
}
