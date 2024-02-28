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

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var siphKey string

func main() {
	var secretWord string
	var dicPath string = "dataset.txt"
	wordKeeper := []string{}
	vecHolder := VecBaseInicialisation()

	siphKey = createKey()

	router := gin.Default()
	router.Use(cors.Default())

	router.GET("/random_word", func(c *gin.Context) {
		getRandWord(c, &secretWord, &wordKeeper, dicPath)
	})

	router.GET("/similarity", func(c *gin.Context) {
		secretWord := decrypt(strings.ToLower(c.Query("secret_word")))
		query := c.Query("query")
		getSimilarityPerc(c, secretWord, vecHolder, query)
	})

	router.GET("/check", func(c *gin.Context) {
		query := c.Query("query")
		secretWord := decrypt(c.Query("secret_word"))
		getCheck(c, secretWord, query)
	})

	router.GET("/hint", func(c *gin.Context) {
		query := c.Query("query")
		secretWord := decrypt(c.Query("secret_word"))
		getShowHint(c, vecHolder, secretWord, query)
	})

	router.GET("/show_finish", func(c *gin.Context) {
		secretWord := decrypt(c.Query("secret_word"))
		getShowFinish(c, secretWord)
	})

	router.Run("localhost:8080")
}

func getRandWord(c *gin.Context, secretWord *string, wordKeeper *[]string, dicPath string) {

	file, err := ioutil.ReadFile(dicPath)
	if err != nil {
		log.Fatal(err)
	}

	*wordKeeper = strings.Split(string(file), "\n")
	*secretWord = strings.ReplaceAll((*wordKeeper)[rand.Intn(len(*wordKeeper))]+"_NOUN", "\r", "")
	log.Println(*secretWord)
	word := encrypt(*secretWord)
	log.Println(word)

	c.IndentedJSON(http.StatusCreated, word)
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

func getCheck(c *gin.Context, secretWord string, query string) {
	c.IndentedJSON(http.StatusFound, secretWord == query+"_NOUN")
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

func getShowFinish(c *gin.Context, secretWord string) {
	c.IndentedJSON(http.StatusOK, secretWord)
}
