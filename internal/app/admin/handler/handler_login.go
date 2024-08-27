package handler

import (
	"net/http"

	"github.com/Format-C-eft/utils/headers"
	"github.com/Format-C-eft/utils/logger"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *StoreImpl) Login(ginContext *gin.Context) {
	session := sessions.Default(ginContext)
	sessionID, _ := uuid.NewUUID()
	session.Set("sessionID", sessionID.String())
	if err := session.Save(); err != nil {
		logger.ErrorKV(ginContext, "cant generate admin sessionID", "err", err)
		ginContext.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ginContext.Writer.Header().Set(headers.Location, "/")
	ginContext.AbortWithStatus(http.StatusFound)
}
