package messages

import (
	"fmt"
	"github.com/Niromash/niromash-api/api"
	"github.com/Niromash/niromash-api/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"regexp"
	"strings"
	"time"
)

func ListMessageFromScopeRoute() *api.Route {
	return &api.Route{
		Path:   "/:scope",
		Method: api.MethodGet,
		Handler: func(c *gin.Context, user api.User, service api.MainService) {
			scope := c.Param("scope")
			key := c.Param("key")
			autoTranslateSubTranslation := c.Query("autoTranslateSubTranslation") == "true"

			if len(key) == 0 {
				entireScope, err := service.Messages().GetAllMessagesFromScope(scope)
				if err != nil {
					return
				}

				if autoTranslateSubTranslation {
					replaceAllTranslations(service, entireScope...)
				}

				c.JSON(200, entireScope)
				return
			}

			message, err := service.Messages().GetMessage(scope, key)
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					c.JSON(404, gin.H{
						"message": "Message not found",
					})
				}
				return
			}

			if autoTranslateSubTranslation {
				replaceAllTranslations(service, message)
			}

			c.JSON(200, message)
		},
	}
}

func replaceAllTranslations(service api.MainService, msgs ...*model.Message) {
	startTime := time.Now()
	regex, err := regexp.Compile(`%([^%]*)@(.*?)%`)
	if err == nil {
		for _, message := range msgs {
			for _, translation := range message.Translations {
				translationKey := translation.Locale
				translationValue := translation.Value
				matchs := regex.FindAllStringSubmatch(translationValue, -1)
				if len(matchs) > 0 { // Found sub translations, so translate them
					for _, match := range matchs {
						allGroup, keyFound, scopeFound := match[0], match[1], match[2]

						replaceOriginalTranslation := func(foundMessage *model.Message) {
							langToFetch := translationKey
							var foundTranslation *model.MessageTranslation
							var foundUsTranslation *model.MessageTranslation
							for _, messageTranslation := range foundMessage.Translations {
								if messageTranslation.Locale == langToFetch {
									foundTranslation = messageTranslation
									break
								}
								if messageTranslation.Locale == "en_US" {
									foundUsTranslation = messageTranslation
								}
							}

							if foundTranslation == nil {
								if foundUsTranslation == nil {
									if len(foundMessage.Translations) > 0 {
										foundTranslation = foundMessage.Translations[0]
									}
								} else {
									foundTranslation = foundUsTranslation
								}
							}

							if foundTranslation != nil {
								translation.Value = strings.ReplaceAll(translationValue, allGroup, foundTranslation.Value)
							}
						}

						if scopeFound == message.Scope.Name { // if the scope is the same, look through all messages in the current scope
							var foundMsg *model.Message
							for _, msg := range msgs {
								if msg.Key == keyFound {
									foundMsg = msg
									break
								}
							}

							if foundMsg != nil {
								replaceOriginalTranslation(foundMsg)
								continue
							}
						}

						// The scope is not the same, or the translation was not found, so fetch it from the database
						trans, err := service.Messages().GetMessage(scopeFound, keyFound)
						if err != nil || len(trans.Translations) == 0 {
							continue
						}
						replaceOriginalTranslation(trans)
					}
				}
			}
		}
	}
	fmt.Printf("Replaced all sub translations in %s\n", time.Since(startTime))
}
