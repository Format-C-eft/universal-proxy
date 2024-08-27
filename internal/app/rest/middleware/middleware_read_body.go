package middleware

import (
	"bytes"
	"io"

	"github.com/Format-C-eft/universal-proxy/internal/config"
	"github.com/Format-C-eft/utils/logger"
	"github.com/gin-gonic/gin"
)

func (m *StoreImpl) ReadRequestBody() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		rawData, err := ginContext.GetRawData()
		if err != nil {
			logger.ErrorKV(
				ginContext,
				"ReadRequestBody.GetRawData",
				"err", err,
				"method", ginContext.Request.Method,
				"url", ginContext.Request.URL.Path,
			)

			ginContext.Set(config.ContextRequestBodyKey, []byte{})

			return
		}

		ginContext.Set(config.ContextRequestBodyKey, rawData)

		ginContext.Request.GetBody = func() (io.ReadCloser, error) {
			ginContext.Request.Body = io.NopCloser(bytes.NewBuffer(rawData))
			buffer := bytes.NewBuffer(rawData)
			closer := io.NopCloser(buffer)
			return closer, nil
		}

		body, _ := ginContext.Request.GetBody()
		ginContext.Request.Body = body
	}
}
