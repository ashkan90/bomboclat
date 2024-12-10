package main

import (
	"context"
	"log"
)

//go:generate go run ../generator/main.go

func main() {
	ctx := context.Background()

	userClient := &UserClientImpl{}
	userInfoClient := &UserInfoClientImpl{}
	basketClient := &BasketClientImpl{}

	// Type-safe spawn kullanımı
	result := gummy.Process(ctx,
		gummy.UserSpawn(userClient.GetUserById, ctx, 13),
		gummy.UserInfoSpawn(userInfoClient.GetUserInfo, ctx, 13),
		gummy.BasketSpawn(basketClient.GetBasket, ctx, 13),
	)

	if result.Err != nil {
		log.Fatal(result.Err)
	}

	// Type assertion with generated helper
	user := result.Get[*UserModel](0)
	userInfo := result.Get[*UserInfoModel](1)
	basket := result.Get[*BasketModel](2)
}
