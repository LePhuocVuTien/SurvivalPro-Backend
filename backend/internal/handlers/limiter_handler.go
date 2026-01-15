package handlers

import "github.com/LePhuocVuTien/SurvivalPro-Backend/infrastructure/ratelimit"

type LimiterHandler struct {
	Limiter *ratelimit.RedisLimiter
}

func NewHandler(limiter *ratelimit.RedisLimiter) *LimiterHandler {
	return &LimiterHandler{Limiter: limiter}
}
