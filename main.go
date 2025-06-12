package main

import (
	"fmt"
	"os"

	"test-db/database"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db_user := os.Getenv("DB_USER")
	db_password := os.Getenv("DB_PASSWORD")
	db_host := os.Getenv("DB_HOST")
	db_port := os.Getenv("DB_PORT")
	db_name := "test"

	dsn := db_user + ":" + db_password + "@tcp(" + db_host + ":" + db_port + ")/" + db_name + "?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// 自动迁移
	if err := AutoMigrate(db); err != nil {
		panic(err)
	}

	// 创建用户
	user := database.UserProfile{Address: "0x12345..."}
	db.Create(&user)

	// 绑定Twitter
	err = database.BindTwitterInfo(db, user.ID, "12345", "Twitter User", "twitteruser")
	if err != nil {
		fmt.Println(err)
	}

	// 绑定Discord (此时DiscordID可以为空)
	err = database.BindDiscordInfo(db, user.ID, "12345", "Discord User", "discorduser", "")
	if err != nil {
		fmt.Println(err)
	}

	err = database.BindTelegramInfo(db, user.ID+1, 12345, "Telegram User1", "Telegram User2", "Telegramuser", "")
	if err != nil {
		fmt.Println(err)
	}

	// 获取用户完整信息
	fullProfile, err := database.GetUserFullProfile(db, user.ID)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%+v\n", fullProfile)
	}

	var infos []database.TelegramInfo
	err = db.Model(&database.TelegramInfo{}).Find(&infos).Error
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%+v\n", infos)
	}
}

func AutoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&database.UserProfile{},
		&database.EmailInfo{},
		&database.TwitterInfo{},
		&database.TelegramInfo{},
		&database.DiscordInfo{},
	)
	if err != nil {
		return fmt.Errorf("failed to auto migrate: %v", err)
	}
	return nil
}
