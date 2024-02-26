package main

import (
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	var secretWord string
	var dicPath string = "dataset.txt"
	wordKeeper := []string{}
	//vecHolder := VecBaseInicialisation()

	router := gin.Default()

	router.GET("/random_word", func(c *gin.Context) {
		getRandWord(c, &secretWord, &wordKeeper, dicPath /*, vecHolder*/)
	})
	router.GET("/similarity")
	router.GET("/best_perc")
	router.GET("/hint")
	router.GET("/show_finish")

	router.Run("localhost:8080")
}

func getRandWord(c *gin.Context, secretWord *string, wordKeeper *[]string, dicPath string /*, vecHolder map[string][]float64*/) {

	file, err := ioutil.ReadFile(dicPath)
	if err != nil {
		log.Fatal(err)
	}

	*wordKeeper = strings.Split(string(file), "\n")
	*secretWord = (*wordKeeper)[rand.Intn(len(*wordKeeper))]

	c.IndentedJSON(http.StatusOK, secretWord)
}

func getSimilarityPerc(c *gin.Context) {

}

func getBestPerc(c *gin.Context) {

}

func getShowHint(c *gin.Context) {

}

func getShowFinish(c *gin.Context) {

}
