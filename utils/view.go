package utils

import (
	"fmt"
	"time"

	"github.com/klauspost/lctime"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func FormataMoedaBR(value float64) string {
	p := message.NewPrinter(language.BrazilianPortuguese)
	return p.Sprintf("R$ %.2f", value)
}

func FormataDataLocaleBR(date time.Time) string {
	lctime.SetLocale("pt_BR")
	return lctime.Strftime("%d de %b %Y", date)
}

func FormataTempoRelativo(date time.Time) string {
	now := time.Now().Add(-3 * time.Hour)
	duration := now.Sub(date)

	if years := int(duration.Hours() / 24 / 365); years > 0 {
		return pluralize(years, "ano", "anos")
	}
	if months := int(duration.Hours() / 24 / 30); months > 0 {
		return pluralize(months, "mês", "meses")
	}
	if days := int(duration.Hours() / 24); days > 0 {
		return pluralize(days, "dia", "dias")
	}
	if hours := int(duration.Hours()); hours > 0 {
		return pluralize(hours, "hora", "horas")
	}
	if minutes := int(duration.Minutes()); minutes > 0 {
		return pluralize(minutes, "minuto", "minutos")
	}
	seconds := int(duration.Seconds())
	return pluralize(seconds, "segundo", "segundos")
}

func pluralize(value int, singular, plural string) string {
	if value == 1 {
		return fmt.Sprintf("1 %s", singular)
	}
	return fmt.Sprintf("%d %s", value, plural)
}
