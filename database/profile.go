package database

import (
	"fmt"

	"gorm.io/gorm"
)

type EmailInfo struct {
	Email         string `gorm:"column:email;type:varchar(255);primaryKey"`
	UserProfileID uint   `gorm:"column:user_profile_id;index"`
}

type TwitterInfo struct {
	TwitterID       string `gorm:"column:twitter_id;type:varchar(255);primaryKey"`
	TwitterName     string `gorm:"column:twitter_name"`
	TwitterUserName string `gorm:"column:twitter_username"`
	UserProfileID   uint   `gorm:"column:user_profile_id;index"`
}

type TelegramInfo struct {
	TelegramID        int    `gorm:"column:telegram_id;type:varchar(255);primaryKey"`
	TelegramFirstName string `gorm:"column:telegram_first_name"`
	TelegramLastName  string `gorm:"column:telegram_last_name"`
	TelegramUserName  string `gorm:"column:telegram_username"`
	TelegramPhoto     string `gorm:"column:telegram_photo"`
	UserProfileID     uint   `gorm:"column:user_profile_id;index"`
}

type DiscordInfo struct {
	DiscordID       string `gorm:"column:discord_id;type:varchar(255);primaryKey"`
	DiscordName     string `gorm:"column:discord_name"`
	DiscordUserName string `gorm:"column:discord_username"`
	DiscordEmail    string `gorm:"column:discord_email"`
	UserProfileID   uint   `gorm:"column:user_profile_id;index"`
}

type UserProfile struct {
	gorm.Model
	Address string `gorm:"column:address;type:varchar(255);uniqueIndex"`

	EmailInfo    EmailInfo    `gorm:"foreignKey:UserProfileID;references:ID"`
	TwitterInfo  TwitterInfo  `gorm:"foreignKey:UserProfileID;references:ID"`
	TelegramInfo TelegramInfo `gorm:"foreignKey:UserProfileID;references:ID"`
	DiscordInfo  DiscordInfo  `gorm:"foreignKey:UserProfileID;references:ID"`
}

func BindEmailInfo(db *gorm.DB, userID uint, email string) error {
	var userEmail EmailInfo
	if err := db.Where("user_profile_id = ?", userID).First(&userEmail).Error; err == nil {
		return fmt.Errorf("您已绑定过邮箱")
	}

	emailInfo := EmailInfo{
		Email:         email,
		UserProfileID: userID,
	}

	if err := db.Create(&emailInfo).Error; err != nil {
		return fmt.Errorf("绑定邮箱失败: %v", err)
	}

	return nil
}

func BindTelegramInfo(db *gorm.DB, userID uint, telegramID int, telegramFirstName, telegramLastName, telegramUserName, telegramPhoto string) error {
	// 检查该Telegram账号是否已被绑定
	var existingTelegram TelegramInfo
	if err := db.Where("telegram_id = ?", telegramID).First(&existingTelegram).Error; err == nil {
		return fmt.Errorf("该Telegram账号已被绑定")
	}

	// 检查当前用户是否已有Telegram绑定
	var userTelegram TelegramInfo
	if err := db.Where("user_profile_id = ?", userID).First(&userTelegram).Error; err == nil {
		return fmt.Errorf("您已绑定过Telegram账号")
	}

	// 创建新的Telegram绑定
	telegramInfo := TelegramInfo{
		TelegramID:        telegramID,
		TelegramFirstName: telegramFirstName,
		TelegramLastName:  telegramLastName,
		TelegramUserName:  telegramUserName,
		TelegramPhoto:     telegramPhoto,
		UserProfileID:     userID,
	}

	if err := db.Create(&telegramInfo).Error; err != nil {
		return fmt.Errorf("绑定Telegram信息失败: %v", err)
	}

	return nil
}

func BindDiscordInfo(db *gorm.DB, userID uint, discordID, discordName, discordUserName, discordEmail string) error {
	// 如果discordID不为空，检查是否已被绑定
	if discordID != "" {
		var existingDiscord DiscordInfo
		if err := db.Where("discord_id = ?", discordID).First(&existingDiscord).Error; err == nil {
			return fmt.Errorf("该Discord账号已被绑定")
		}
	}

	// 检查当前用户是否已有Discord绑定
	var userDiscord DiscordInfo
	if err := db.Where("user_profile_id = ?", userID).First(&userDiscord).Error; err == nil {
		return fmt.Errorf("您已绑定过Discord账号")
	}

	// 创建新的Discord绑定
	discordInfo := DiscordInfo{
		DiscordID:       discordID,
		DiscordName:     discordName,
		DiscordUserName: discordUserName,
		DiscordEmail:    discordEmail,
		UserProfileID:   userID,
	}

	if err := db.Create(&discordInfo).Error; err != nil {
		return fmt.Errorf("绑定Discord信息失败: %v", err)
	}

	return nil
}

func BindTwitterInfo(db *gorm.DB, userID uint, twitterID, twitterName, twitterUserName string) error {
	// 检查该Twitter账号是否已被绑定（无论是否当前用户）
	var existingTwitter TwitterInfo
	if err := db.Where("twitter_id = ?", twitterID).First(&existingTwitter).Error; err == nil {
		return fmt.Errorf("该Twitter账号已被绑定")
	}

	// 检查当前用户是否已有Twitter绑定
	var userTwitter TwitterInfo
	if err := db.Where("user_profile_id = ?", userID).First(&userTwitter).Error; err == nil {
		return fmt.Errorf("您已绑定过Twitter账号")
	}

	// 创建新的Twitter绑定
	twitterInfo := TwitterInfo{
		TwitterID:       twitterID,
		TwitterName:     twitterName,
		TwitterUserName: twitterUserName,
		UserProfileID:   userID,
	}

	if err := db.Create(&twitterInfo).Error; err != nil {
		return fmt.Errorf("绑定Twitter信息失败: %v", err)
	}

	return nil
}

func GetUserFullProfile(db *gorm.DB, userID uint) (*UserProfile, error) {
	var user UserProfile
	if err := db.Preload("EmailInfo").
		Preload("TwitterInfo").
		Preload("TelegramInfo").
		Preload("DiscordInfo").
		First(&user, userID).Error; err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %v", err)
	}
	return &user, nil
}
