package app

import (
	"context"
	"fmt"
	"github.com/nozgurozturk/marvin/notifier/entity"
	"github.com/nozgurozturk/marvin/notifier/internal/config"
	"github.com/nozgurozturk/marvin/notifier/internal/service"
	"github.com/nozgurozturk/marvin/pkg/utils"
	"time"
)

/*
	1. GetAllSubscribers
	2. ChecksubscriberTimePast
	3. GetRepositoryOutedatedPackeges
	4. sendmail
*/

func checkExactTime(sub *entity.SubscriberDTO) bool {

	now := time.Now().UTC()
	if sub.Notify.Minute == now.Minute() {
		if sub.Notify.Frequency == entity.Day {
			if sub.Notify.Hour == now.Hour() {
				return true
			} else {
				return false
			}
		}

		if sub.Notify.Frequency == entity.Week {
			if sub.Notify.Weekday == now.Weekday() && sub.Notify.Hour == now.Hour() {
				return true
			} else {
				return false
			}
		}
		return true
	}
	return false
}

func checkAvailableSubscriber(sub *entity.SubscriberDTO) bool {
	now := time.Now().UTC()
	if sub.Notify.Frequency == entity.Week {
		if sub.Notify.Weekday < now.Weekday() {
			return false
		}
		if sub.Notify.Hour < now.Hour() {
			return false
		}
	}
	if sub.Notify.Frequency == entity.Day {
		if sub.Notify.Hour < now.Hour() {
			return false
		}
		if sub.Notify.Minute < now.Minute() {
			return false
		}
	}

	if sub.Notify.Frequency == entity.Hour {
		if sub.Notify.Minute < now.Minute() {
			return false
		}
	}

	return true
}

func getAllAvailableSubscribers(s service.SubscriberService) []*entity.SubscriberDTO {
	allSubs, err := s.GetAll()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	var availableSubs []*entity.SubscriberDTO
	for _, sub := range allSubs {
		if checkAvailableSubscriber(sub) {
			availableSubs = append(availableSubs, sub)
		}
	}
	return availableSubs

}

func findOutdatedPackage(repoDTO *entity.RepoDTO) []*entity.Package {
	var outdatedPackages []*entity.Package

	for _, pkg := range repoDTO.PackageList {
		if pkg.IsOutdated {
			outdatedPackages = append(outdatedPackages, pkg)
		}
	}
	return outdatedPackages
}

func SendNotification(s service.SubscriberService, r service.RepoService) {
	ctx := context.Background()

	var availableSubs []*entity.SubscriberDTO

	checkAllSubscribers := func(c context.Context) {
		availableSubs = getAllAvailableSubscribers(s)
	}

	checkNotifyTimeAndSend := func(c context.Context) {
		for _, sub := range availableSubs {
			fmt.Println(sub.ID, sub.Email, sub.Notify.Hour, sub.Notify.Minute)
			go func(s *entity.SubscriberDTO) {

				if checkExactTime(s) {
					repo, err := r.FindById(s.RepoID)
					if err != nil {
						fmt.Println(err)
						return
					}

					outdatedPackages := findOutdatedPackage(repo)
					if outdatedPackages != nil {

						token, err := CreateToken(s)
						if err != nil {
							fmt.Println(err)
							return
						}

						cnf := config.Get().HTTP

						notifyTime := fmt.Sprintf("%02d:%02d", sub.Notify.Hour, sub.Notify.Minute)
						notifyFrequency := fmt.Sprintf("every %s",s.Notify.Frequency)

						if s.Notify.Frequency == entity.Hour {
							notifyTime = fmt.Sprintf("%02d", s.Notify.Minute)
						}
						if s.Notify.Frequency == entity.Week {
							notifyFrequency = fmt.Sprintf("every %s",s.Notify.Weekday.String())
						}
						templateData := struct {
							RepoName             string
							RepoLink             string
							NotifyUpdateLink     string
							NotifyFrequency      string
							NotifyTime           string
							OutdatedPackageCount int
							PackageList          []*entity.Package
							UnsubscribeLink      string
						}{
							OutdatedPackageCount: len(outdatedPackages),
							PackageList:          outdatedPackages,
							NotifyFrequency:      notifyFrequency,
							NotifyTime:           notifyTime,
							RepoName:             repo.Name,
							RepoLink:             repo.Path,
							NotifyUpdateLink:     "http://" + cnf.MainHost + cnf.MainPort + "/subscriber?t=" + token.Token,
							UnsubscribeLink:      "http://" + cnf.MainHost + cnf.MainPort + "/subscriber/unsubscribe?t=" + token.Token,
						}

						emailBody, err := utils.ParseHTMLTemplate("./web/email-sub-notify.html", templateData)
						if err != nil {
							fmt.Println(err)
							return
						}
						subject := "Outdated Packages" + " 📦 " + repo.Name
						err = SendEmail(s.Email, subject, emailBody)
						if err != nil {
							fmt.Println(err)
						}
					}
				}
			}(sub)
		}
	}

	worker := NewScheduler()
	// Get all subscriber every UTC hour
	worker.Add(ctx, checkAllSubscribers, time.Hour)
	// Check subscriber time every UTC minute
	worker.Add(ctx, checkNotifyTimeAndSend, time.Minute)

}
