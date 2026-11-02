package domain

import "time"

func timePtr(value time.Time) *time.Time { return &value }
