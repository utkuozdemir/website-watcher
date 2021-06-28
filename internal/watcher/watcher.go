package watcher

import (
	"fmt"
	"github.com/mxschmitt/playwright-go"
	log "github.com/sirupsen/logrus"
	"github.com/utkuozdemir/website-watcher/internal/notification"
	"github.com/utkuozdemir/website-watcher/internal/website"
	"sync"
	"time"
)

type Watcher interface {
	Start() error
}

func New(websites []website.Website) Watcher {
	return &watcher{websites: websites}
}

type watcher struct {
	websites []website.Website
}

type probeResult struct {
	conditionMet bool
	pageSource   *string
	status       int
}

func (w *watcher) Start() error {
	pw, err := playwright.Run(&playwright.RunOptions{Browsers: []string{"chromium"}})
	if err != nil {
		return err
	}
	defer stopPlaywright(pw)

	browser, err := pw.Chromium.Launch()
	if err != nil {
		return err
	}
	defer closeBrowser(browser)

	var wg sync.WaitGroup
	for _, wb := range w.websites {
		// copy website to avoid data race: https://stackoverflow.com/a/46874398/1005102
		ws := wb
		wg.Add(1)
		go w.watch(&wg, browser, ws)
	}

	wg.Wait()
	return nil
}

func (w *watcher) watch(wg *sync.WaitGroup, browser playwright.Browser, website website.Website) {
	defer wg.Done()
	websiteLogger := log.WithField("url", website.URL())
	for {
		pr, err := w.probe(browser, website)
		if err != nil {
			website.ResetConsecutiveSuccessCount()
			websiteLogger.Warn("Failed to probe website")
		} else if pr.conditionMet {
			website.IncrementConsecutiveSuccessCount()
			successCount := website.ConsecutiveSuccessCount()
			successful := successCount >= website.SuccessThreshold()
			didNotNotifyRecently := time.Since(website.LastNotificationTime()) > website.NotificationDelay()
			if successful && didNotNotifyRecently {
				pi := notification.NewPageInfo(website.Name(), website.URL(), pr.status)
				h := website.NotificationHandler()
				h.Handle(pi)
			} else if successful && !didNotNotifyRecently {
				websiteLogger.Info("Successful but notified recently, skipping notification")
			} else if !successful {
				websiteLogger.WithField("success", successCount).
					WithField("threshold", website.SuccessThreshold()).
					Info("On the pathway to success")
			}
		} else {
			website.ResetConsecutiveSuccessCount()
			websiteLogger.Info("Condition was not met")
		}

		delay := website.Delay()
		websiteLogger.
			WithField("delay", fmt.Sprintf("%s", delay)).
			Info("Sleeping...")
		time.Sleep(delay)
	}
}

func (w *watcher) probe(browser playwright.Browser, website website.Website) (*probeResult, error) {
	pageSource, status, err := getPageSource(browser, website.URL())
	if err != nil {
		return nil, err
	}

	conditionMet := website.Condition()(pageSource, status)
	log.WithField("url", website.URL()).
		WithField("page_source", pageSource).
		Debug()

	r := probeResult{
		conditionMet: conditionMet,
		pageSource:   pageSource,
		status:       status,
	}

	return &r, nil
}

func getPageSource(browser playwright.Browser, url string) (*string, int, error) {
	p, err := browser.NewPage()
	if err != nil {
		return nil, -1, err
	}
	defer closePage(p)

	resp, err := p.Goto(url)
	if err != nil {
		return nil, -1, err
	}

	body, err := resp.Body()
	if err != nil {
		return nil, -1, err
	}

	if resp.Status() != 200 {
		log.WithField("url", url).WithField("status", resp.Status()).
			Warn("Non-200 response status")
	}

	s := string(body)
	return &s, resp.Status(), nil
}

func closePage(p playwright.Page) {
	err := p.Close()
	if err != nil {
		log.WithError(err).Error("Couldn't close page")
	}
}

func closeBrowser(browser playwright.Browser) {
	err := browser.Close()
	if err != nil {
		log.WithError(err).Error("Couldn't close browser")
	}
}

func stopPlaywright(pw *playwright.Playwright) {
	err := pw.Stop()
	if err != nil {
		log.WithError(err).Error("Couldn't stop playwright")
	}
}
