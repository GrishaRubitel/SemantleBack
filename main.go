package main

import (
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
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
		query := ctx.Query("query")
		getSimilarityPerc(ctx, secretWord, vecHolder, query)
	})
	router.GET("/best_perc")
	router.GET("/hint", func(c *gin.Context) {
		query := c.Query("best_word")
		getShowHint(c, vecHolder, secretWord, query)
	})
	router.GET("/show_finish")

	router.Run("localhost:8080")
}

func getRandWord(c *gin.Context, secretWord *string, wordKeeper *[]string, dicPath string) {

	file, err := ioutil.ReadFile(dicPath)
	if err != nil {
		log.Fatal(err)
	}

	*wordKeeper = strings.Split(string(file), "\n")
	*secretWord = strings.ReplaceAll((*wordKeeper)[rand.Intn(len(*wordKeeper))]+"_NOUN", "\r", "")

	c.IndentedJSON(http.StatusCreated, *secretWord)
}

func getSimilarityPerc(c *gin.Context, secretWord string, vecHolder map[string][]float64, query string) {

	start := time.Now()
	var res string

	for k := range vecHolder {
		if k == query+"_NOUN" {
			res = "Word in model"
			c.IndentedJSON(http.StatusFound, processQuery(vecHolder[k], vecHolder[secretWord]))
			return
		} else {
			res = "Word not in model"
		}
	}

	log.Println("Time since start: ", time.Since(start), res)
	c.IndentedJSON(http.StatusNoContent, query)
}

func processQuery(queryVec []float64, secretVec []float64) float64 {
	var res float64 = 0
	var scalProd float64 = 0
	var queryMod float64 = 0
	var secretMod float64 = 0

	for i1 := 0; i1 < len(queryVec); i1++ {
		scalProd += queryVec[i1] * secretVec[i1]
		queryMod += math.Pow(queryVec[i1], 2)
		secretMod += math.Pow(secretVec[i1], 2)

	}
	res = scalProd / (math.Sqrt(queryMod) * math.Sqrt(secretMod))

	return math.Floor(math.Abs(res * 100))
}

func getBestPerc(c *gin.Context) {

}

func getShowHint(c *gin.Context, vecHolder map[string][]float64, secretWord string, query string) {

	var hintHolder []string
	minPerc, _ := strconv.ParseFloat(query, 64)

	for k := range vecHolder {
		if strings.Contains(k, "_NOUN") {
			if kPerc := processQuery(vecHolder[k], vecHolder[secretWord]); kPerc > minPerc {
				hintHolder = append(hintHolder, strconv.FormatFloat(kPerc, 'f', -1, 64)+" "+strings.ReplaceAll(k, "_NOUN", ""))
			}
			if len(hintHolder) == 100 {
				break
			}
		}
	}

	c.IndentedJSON(http.StatusFound, hintHolder[rand.Intn(len(hintHolder))])
}

func getShowFinish(c *gin.Context) {

}
