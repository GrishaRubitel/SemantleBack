package main

import (
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	var secretWord string
	var dicPath string = "dataset.txt"
	wordKeeper := []string{}
	vecHolder := VecBaseInicialisation()

	router := gin.Default()

	router.GET("/random_word", func(c *gin.Context) {
		getRandWord(c, &secretWord, &wordKeeper, dicPath)
	})
	router.GET("/similarity", func(ctx *gin.Context) {
		query := ctx.Query("secretWord")
		getSimilarityPerc(ctx, &secretWord, &vecHolder, query)
	})
	router.GET("/best_perc")
	router.GET("/hint")
	router.GET("/show_finish")

	router.Run("localhost:8080")
}

func getRandWord(c *gin.Context, secretWord *string, wordKeeper *[]string, dicPath string) {

	file, err := ioutil.ReadFile(dicPath)
	if err != nil {
		log.Fatal(err)
	}

	*wordKeeper = strings.Split(string(file), "\n")
	*secretWord = (*wordKeeper)[rand.Intn(len(*wordKeeper))]

	c.IndentedJSON(http.StatusCreated, strings.ReplaceAll(*secretWord, "\r", ""))
}

func getSimilarityPerc(c *gin.Context, secretWord *string, vecHolder *map[string][]float64, query string) {

	start := time.Now()
	var res string

	for k := range *vecHolder {
		if k == query+"_NOUN" {
			res = "Word in model"
			c.IndentedJSON(http.StatusFound, k)
			return
		} else {
			res = "Word not in model"
		}
	}

	log.Println("Time since start: ", time.Since(start), res)
	c.IndentedJSON(http.StatusNoContent, query)
}

func getBestPerc(c *gin.Context) {

}

func getShowHint(c *gin.Context) {

}

func getShowFinish(c *gin.Context) {

}
