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
	now := time.Now()
	duration := now.Sub(date)

	years := int(duration.Hours() / 24 / 365)
	months := int(duration.Hours() / 24 / 30)
	days := int(duration.Hours() / 24)
	hours := int(duration.Hours())
	minutes := int(duration.Minutes())
	seconds := int(duration.Seconds())

	if years > 0 {
		if years == 1 {
			return "1 ano"
		}
		return fmt.Sprintf("%d anos", years)
	} else if months > 0 {
		if months == 1 {
			return "1 mês"
		}
		return fmt.Sprintf("%d meses", months)
	} else if days > 0 {
		if days == 1 {
			return "1 dia"
		}
		return fmt.Sprintf("%d dias", days)
	} else if hours > 0 {
		if hours == 1 {
			return "1 hora"
		}
		return fmt.Sprintf("%d horas", hours)
	} else if minutes > 0 {
		if minutes == 1 {
			return "1 minuto"
		}
		return fmt.Sprintf("%d minutos", minutes)
	} else {
		if seconds == 1 {
			return "1 segundo"
		}
		return fmt.Sprintf("%d segundos", seconds)
	}
}
