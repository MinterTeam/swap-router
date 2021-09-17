package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"time"
)

type requestLogger struct {
	log *log.Entry
}

func newRequestLogger(c *gin.Context) *requestLogger {
	rid, _ := uuid.NewUUID()
	rlog := log.WithFields(log.Fields{"req": c.Request.RequestURI, "rid": rid, "ip": c.ClientIP()})
	rlog = rlog.WithContext(context.WithValue(context.Background(), "time", time.Now()))
	return &requestLogger{rlog}
}

func (r *requestLogger) Debug(text string) {
	r.log.WithField("t", time.Since(r.log.Context.Value("time").(time.Time))).Debug(text)
}
