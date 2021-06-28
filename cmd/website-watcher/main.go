package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/utkuozdemir/website-watcher/internal/notification"
	"github.com/utkuozdemir/website-watcher/internal/watcher"
	"github.com/utkuozdemir/website-watcher/internal/website"
	"os"
	"strings"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		PadLevelText:  true,
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {
	pushoverToken := os.Getenv("PUSHOVER_TOKEN")
	pushoverRecipient := os.Getenv("PUSHOVER_USER")
	if pushoverToken == "" || pushoverRecipient == "" {
		log.Error("Missing pushover env variables PUSHOVER_TOKEN or PUSHOVER_USER")
		os.Exit(1)
	}

	//handler := notification.NewPushoverHandler(pushoverToken, pushoverRecipient)
	handler := notification.NewLoggingHandler()

	websites := []website.Website{
		website.MustNewWithDefaults("otto.de",
			"https://www.otto.de/technik/gaming/playstation/ps5/",
			doesNotContainAnd200Condition("leider ausverkauft"), handler),
		website.MustNewWithDefaults(
			"mediamarkt.de",
			"https://www.mediamarkt.de/de/product/_sony-playstation%C2%AE5-2661938.html",
			doesNotContainAnd200Condition("aktuell nicht verfügbar"), handler,
		),
		website.MustNewWithDefaults(
			"saturn.de",
			"https://www.saturn.de/de/product/_sony-playstation%C2%AE5-2661938.html",
			doesNotContainAnd200Condition("aktuell nicht verfügbar"),
			handler,
		),
		website.MustNewWithDefaults(
			"euronics.de disc",
			"https://www.euronics.de/spiele-und-konsolen-film-und-musik/spiele-und-konsolen/playstation-5/spielekonsole/playstation-5-konsole-4061856837826",
			doesNotContainAnd200Condition("Onlineshop nicht verfügbar"),
			handler,
		),
		website.MustNewWithDefaults(
			"euronics.de digital",
			"https://www.euronics.de/spiele-und-konsolen-film-und-musik/spiele-und-konsolen/playstation-5/spielekonsole/playstation-5-digital-edition-konsole-4061856837833",
			doesNotContainAnd200Condition("Onlineshop nicht verfügbar"),
			handler,
		),
		website.MustNewWithDefaults(
			"amazon.de disc",
			"https://www.amazon.de/dp/B08H93ZRK9/?coliid=ICXS62JGGZSY5&colid=1MS6PECO7SBJ8&psc=0&ref_=lv_ov_lig_dp_it&tag=mmo-deals-21",
			doesNotContainAnd200Condition("Derzeit nicht verfügbar"),
			handler,
		),
		website.MustNewWithDefaults(
			"amazon.de digital",
			"https://www.amazon.de/dp/B08H98GVK8/?coliid=I2DTD5BHDDHCBL&colid=1MS6PECO7SBJ8&psc=0&ref_=lv_ov_lig_dp_it_im&tag=mmo-deals-21",
			doesNotContainAnd200Condition("Derzeit nicht verfügbar"),
			handler,
		),
		website.MustNewWithDefaults(
			"mueller.de disc",
			"https://www.mueller.de/p/playstation-5-konsole-2675305/",
			doesNotContainAnd200Condition("Vielen Dank an alle"),
			handler,
		),
		website.MustNewWithDefaults(
			"mueller.de digital",
			"https://www.mueller.de/p/playstation-5-konsole-digital-edition-2675308/",
			doesNotContainAnd200Condition("Vielen Dank an alle"),
			handler,
		),
		website.MustNewWithDefaults(
			"medimax.de digital",
			"https://www.medimax.de/p/1315337/playstation-5-digital-edition-825gb-ssd",
			doesNotContainAnd200Condition("derzeit keine"),
			handler,
		),
		website.MustNewWithDefaults(
			"medimax.de disc",
			"https://www.medimax.de/p/1315336/play-station-5-825gb-ssd",
			doesNotContainAnd200Condition("derzeit keine"),
			handler,
		),
		website.MustNewWithDefaults(
			"mytoys.de disc",
			"https://www.mytoys.de/sony-ps5-playstation-5-konsole-inkl-laufwerk-17235694.html",
			doesNotContainAnd200Condition("Leider nicht verfügbar"),
			handler,
		),
		website.MustNewWithDefaults(
			"mytoys.de digital",
			"https://www.mytoys.de/sony-ps5-playstation-5-konsole-digital-version-17235665.html",
			doesNotContainAnd200Condition("Leider nicht verfügbar"),
			handler,
		),

	}

	w := watcher.New(websites)
	err := w.Start()
	if err != nil {
		log.WithError(err).Error("Failed to start watcher")
		os.Exit(1)
	}
}

func doesNotContainAnd200Condition(text string) func(*string, int) bool {
	return func(pageSource *string, status int) bool {
		return status == 200 && !strings.Contains(*pageSource, text)
	}
}
