package main

import (
	"auth_service/internal/config"
	"auth_service/internal/models"
	"auth_service/internal/repository"
	"fmt"
	"log"
	"slices"
)

func initAdmins(cfg *config.Config, userRepo repository.UserRepository) {
	admins, err := userRepo.GetByRole(models.AdminRole)
	if err != nil {
		admins = []models.User{}
	}
	adminStr := "Admins: "
	for _, admin := range admins {
		adminStr += fmt.Sprintf(" %d", admin.TelegramID)
	}
	log.Print(adminStr)

	for _, admin_id := range cfg.Admins {
		if !slices.ContainsFunc(admins, func(u models.User) bool {
			return u.TelegramID == admin_id
		}) {
			admin, err := userRepo.Get(admin_id)
			if err != nil {
				admin, err = userRepo.Create(&models.User{
					TelegramID: admin_id,
				})
				if err != nil {
					continue
				}
			}
			roles := []string{models.AdminRole}
			userRepo.UpdateUserRoles(admin.TelegramID, roles)
			admin, _ = userRepo.Get(admin.TelegramID)
		}
	}
}
