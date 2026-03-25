package widgets

import (
	"fmt"

	"go.a8l.eu/ccstatusline/internal/config"
	"go.a8l.eu/ccstatusline/internal/input"
)

func init() {
	Register("rate-limit-five-hour", rateLimitFiveHourWidget{})
	Register("rate-limit-seven-day", rateLimitSevenDayWidget{})
}

type rateLimitFiveHourWidget struct{}

func (rateLimitFiveHourWidget) Render(si *input.StatusInput, _ config.WidgetItem) string {
	if si == nil || si.RateLimits == nil || si.RateLimits.FiveHour == nil {
		return ""
	}
	p := si.RateLimits.FiveHour.UsedPercentage
	if p == nil {
		return ""
	}
	return fmt.Sprintf("%.0f%%", *p)
}

type rateLimitSevenDayWidget struct{}

func (rateLimitSevenDayWidget) Render(si *input.StatusInput, _ config.WidgetItem) string {
	if si == nil || si.RateLimits == nil || si.RateLimits.SevenDay == nil {
		return ""
	}
	p := si.RateLimits.SevenDay.UsedPercentage
	if p == nil {
		return ""
	}
	return fmt.Sprintf("%.0f%%", *p)
}
