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

	now := time.Now()
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
	now := time.Now()
	if sub.Notify.Frequency == entity.Week {
		if sub.Notify.Weekday < now.Weekday() {
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
		if sub.Notify.Hour < now.Hour() {
			return false
		}
	}
	return true
}

func getAllAvailableSubscribers(s service.SubscriberService) []*entity.SubscriberDTO {
	allSubs, err := s.GetAll()
	if err != nil {
		fmt.Println(err)
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
			fmt.Println(sub.Email, sub.Notify.Minute, checkExactTime(sub))
			if checkExactTime(sub) {
				repo, err := r.FindById(sub.RepoID)
				if err != nil {
					fmt.Println(err)
					continue
				}

				outdatedPackages := findOutdatedPackage(repo)
				if outdatedPackages != nil {

					token, err := CreateToken(sub)
					if err != nil {
						fmt.Println(err)
						continue
					}

					cnf := config.Get().HTTP
					notifyTime := fmt.Sprintf("%02d:%02d", sub.Notify.Hour, sub.Notify.Minute)
					templateData := struct {
						RepoName             string
						RepoLink             string
						NotifyUpdateLink     string
						NotifyFrequency      entity.Frequency
						NotifyTime           string
						OutdatedPackageCount int
						PackageList          []*entity.Package
						UnsubscribeLink      string
					}{
						OutdatedPackageCount: len(outdatedPackages),
						PackageList:          outdatedPackages,
						NotifyFrequency:      sub.Notify.Frequency,
						NotifyTime:           notifyTime,
						RepoName:             repo.Name,
						RepoLink:             repo.Path,
						NotifyUpdateLink:     "http://" + cnf.MainHost + cnf.MainPort + "/subscriber?t=" + token.Token,
						UnsubscribeLink:      "http://" + cnf.MainHost + cnf.MainPort + "/subscriber/unsubscribe?t=" + token.Token,
					}

					emailBody, err := utils.ParseHTMLTemplate("./web/email-sub-notify.html", templateData)
					if err != nil {
						fmt.Println(err)
						continue
					}
					// TODO: Make it concurrent
					subject := repo.Name + " 📦 " + "Outdated Package List"
					err = SendEmail(sub.Email, subject, emailBody)
					if err != nil {
						fmt.Println(err)
						continue
					}

				}
			}
		}
	}

	worker := NewScheduler()
	worker.Add(ctx, checkAllSubscribers, time.Hour*1)
	worker.Add(ctx, checkNotifyTimeAndSend, time.Minute*1)

}
