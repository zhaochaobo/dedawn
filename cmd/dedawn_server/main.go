package main

import (
	"dedawn/card"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func main() {
	addr := ":8080"
	flag.StringVar(&addr, "addr", ":8080", "listen address")
	flag.Parse()
	fmt.Println(addr)

	cards := make(map[string]card.Card)
	adminCard := card.Card{
		No:       "dedawn",
		Secret:   "dedawn@2024",
		Amount:   10000 * time.Hour,
		CreateAt: time.Now(),
	}
	cards[adminCard.No] = adminCard

	r := gin.Default()

	r.Use(func(ctx *gin.Context) {
		fmt.Println("in middleware")
		ctx.Next()
		if err := ctx.Err(); err != nil {
			fmt.Printf("the handle error, %v\n", err)
		}
		fmt.Printf("the handle error, %v\n", ctx.Errors)
		if len(ctx.Errors) > 0 {
			ctx.JSON(ctx.Writer.Status(), ctx.Errors)
		}
	})

	r.GET("/check", func(c *gin.Context) {
		q := struct {
			No     string `form:"no"`
			Secret string `form:"secret"`
		}{}
		if err := c.BindQuery(&q); err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		cd, ok := cards[q.No]
		if !ok {
			_ = c.AbortWithError(http.StatusNotFound, fmt.Errorf("the card no %s not exists", q.No))
			return
		}

		if cd.Secret != strings.TrimSpace(q.Secret) {
			_ = c.AbortWithError(http.StatusForbidden, fmt.Errorf("the card no or secret is incurrect"))
			return
		}
		cd.ConsumeAt = time.Now()

		c.JSON(http.StatusOK, cd)
	})

	// grant play duration
	r.GET("/grant", func(c *gin.Context) {

		q := struct {
			Amount time.Duration `form:"amount"`
		}{}
		if err := c.BindQuery(&q); err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		// clean outdated cards
		for n, c := range cards {
			if time.Now().Sub(c.CreateAt) > 24*time.Hour {
				delete(cards, n)
			}
		}

		no := rand.Intn(10000) + 10000
		secret := randSeq(5)
		cd := card.Card{
			No:       strconv.Itoa(no),
			Secret:   secret,
			Amount:   q.Amount,
			CreateAt: time.Now(),
		}
		cards[cd.No] = cd

		c.JSON(http.StatusOK, cd)
	})

	r.GET("/deduct", func(c *gin.Context) {
		q := struct {
			No     string        `form:"no"`
			Amount time.Duration `form:"amount"`
		}{}
		if err := c.BindQuery(&q); err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		cd, ok := cards[q.No]
		if !ok {
			_ = c.AbortWithError(http.StatusNotFound, fmt.Errorf("the card [%s] not found", q.No))
			return
		}

		fmt.Printf("card %s deduct %d", q.No, q.Amount)

		cd.Amount = cd.Amount - q.Amount
		if cd.Amount <= 0 {
			_ = c.AbortWithError(http.StatusGone, fmt.Errorf("the card has depleted"))
			delete(cards, q.No)
			return
		}
		cards[q.No] = cd

		c.JSON(http.StatusOK, cd)

	})
	_ = r.Run(addr) // list
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
