package domain

import "math"

type Status struct {
	MessageBytes  float64
	ConnectionNum float64
}

func (stat *Status) Avg(windowSize float64) *Status {
	return &Status{
		MessageBytes:  stat.MessageBytes / windowSize,
		ConnectionNum: stat.ConnectionNum / windowSize,
	}
}

func (sum *Status) Sub(stat *Status) {
	if stat != nil {
		sum.ConnectionNum -= stat.ConnectionNum
		sum.MessageBytes -= stat.MessageBytes
	}
}

func (sum *Status) Add(stat *Status) {
	if stat != nil {
		sum.ConnectionNum += stat.ConnectionNum
		sum.MessageBytes += stat.MessageBytes
	}
}

func (stat *Status) CalculateActiveSorce() float64 {
	return getGB(stat.MessageBytes)
}

func (stat *Status) CalculateStaticSorce() float64 {
	return stat.ConnectionNum
}

func getGB(m float64) float64 {
	return decimal(m / (1 << 30))
}

func decimal(value float64) float64 {
	return math.Trunc(value*1e2+0.5) * 1e-2
}
