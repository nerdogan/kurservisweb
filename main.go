package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

/* ================= GLOBAL ================= */

var (
	db     *sql.DB
	apiURL string
)

/* ================= API MODELLER ================= */

type ApiResponse struct {
	Data []ApiItem `json:"data"`
}

type ApiItem struct {
	MarketProductId int     `json:"marketProductId"`
	UpdatedAt       string  `json:"updatedAt"`
	CustomerBuysAt  float64 `json:"customerBuysAt"`
	CustomerSellsAt float64 `json:"customerSellsAt"`
}

/* ================= MAIN ================= */

func main() {
	var err error

	// .env yükle
	if err := godotenv.Load(); err != nil {
		log.Fatal(".env yüklenemedi")
	}

	apiURL = os.Getenv("API_URL")
	dsn := os.Getenv("DB_DSN")

	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	// tablo oluştur
	createKurTable()

	// 15 saniyede bir kur çek (goroutine)
	go startPriceFetcher()

	// GIN
	r := gin.Default()
	// CORS Middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	r.Static("/dist", "./frontend/dist")
	r.Static("/assets", "./frontend/dist/assets")

	r.GET("/", func(c *gin.Context) {
		c.File("./frontend/dist/index.html")
	})

	r.GET("/price", priceHandler)
	r.NoRoute(func(c *gin.Context) {
		c.File("./frontend/dist/index.html")
	})

	log.Println("Server :8080")
	r.Run(":8080")
}

/* ================= PRICE API ================= */

func priceHandler(c *gin.Context) {
	productID, _ := strconv.Atoi(c.Query("productId"))
	//gram, _ := strconv.ParseFloat(c.Query("gram"), 64)
	//factor, _ := strconv.ParseFloat(c.Query("factor"), 64)

	var sellPrice float64
	var updatedAt time.Time

	err := db.QueryRow(`
		SELECT customer_sells_at, updated_at
		FROM kur 
		WHERE market_product_id=$1
		ORDER BY updated_at DESC 
		LIMIT 1
	`, productID).Scan(&sellPrice, &updatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "price not found"})
		return
	}

	if productID == 3 || productID == 4 {
		sellPrice = sellPrice / 1000 // usd ve eur için gram fiyatı
	}

	//price := gram * sellPrice * factor
	price := math.Round(sellPrice*100) / 100 // 2 ondalık basamak

	c.JSON(200, []gin.H{
		{
			"tutar":  price,
			"tarih":  updatedAt.Format("2006-01-02 15:04:05"),
			"masano": strconv.Itoa(productID),
		},
	})
}

/* ================= TABLE ================= */

func createKurTable() {
	// market_product tablosu oluştur
	createMarketProductTable()

	query := `
	CREATE TABLE IF NOT EXISTS kur (
		id SERIAL PRIMARY KEY,
		market_product_id INT REFERENCES market_product(id),
		updated_at TIMESTAMP,
		customer_buys_at NUMERIC(18,5),
		customer_sells_at NUMERIC(18,5),
		created_at TIMESTAMP DEFAULT NOW()
	);`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("kur tablosu oluşturulamadı:", err)
	}
}

func createMarketProductTable() {
	query := `
	CREATE TABLE IF NOT EXISTS market_product (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT NOW()
	);`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("market_product tablosu oluşturulamadı:", err)
	}
}

/* ================= FETCHER ================= */

func startPriceFetcher() {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for {
		fetchAndSave(apiURL)
		<-ticker.C
	}
}
func fetchAndSave(apiURL string) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(apiURL)
	if err != nil {
		log.Println("API error:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("API status error:", resp.Status)
		return
	}

	var apiResp ApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		log.Println("JSON decode error:", err)
		return
	}

	limit := 8
	if len(apiResp.Data) < 8 {
		limit = len(apiResp.Data)
	}

	layout := "2006-01-02T15:04:05.999999"

	for i := 0; i < limit; i++ {
		item := apiResp.Data[i]

		updatedAt, err := time.Parse(layout, item.UpdatedAt)
		if err != nil {
			log.Println("Time parse error:", err)
			continue
		}

		_, err = db.Exec(`
			INSERT INTO kur 
			(market_product_id, updated_at, customer_buys_at, customer_sells_at)
			VALUES ($1, $2, $3, $4)
		`,
			item.MarketProductId,
			updatedAt,
			item.CustomerBuysAt,
			item.CustomerSellsAt,
		)

		if err != nil {
			log.Println("DB insert error:", err)
		}
	}
}
