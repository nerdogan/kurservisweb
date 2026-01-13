package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
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
	MarketProductId  int     `json:"marketProductId"`
	UpdatedAt        string  `json:"updatedAt"`
	CustomerBuysAt   float64 `json:"customerBuysAt"`
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

	r.GET("/", func(c *gin.Context) {
		c.Header("Content-Type", "text/html")
		c.String(200, htmlPage)
	})

	r.GET("/price", priceHandler)

	log.Println("Server :8080")
	r.Run(":8080")
}

/* ================= PRICE API ================= */

func priceHandler(c *gin.Context) {
	productID, _ := strconv.Atoi(c.Query("productId"))
	gram, _ := strconv.ParseFloat(c.Query("gram"), 64)
	factor, _ := strconv.ParseFloat(c.Query("factor"), 64)

	var sellPrice float64
	err := db.QueryRow(`
		SELECT customer_sells_at 
		FROM kur 
		WHERE market_product_id=$1
		ORDER BY updated_at DESC 
		LIMIT 1
	`, productID).Scan(&sellPrice)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "price not found"})
		return
	}

	price := gram * sellPrice * factor

	c.JSON(200, gin.H{
		"price": fmt.Sprintf("%.2f", price),
	})
}

/* ================= TABLE ================= */

func createKurTable() {
	query := `
	CREATE TABLE IF NOT EXISTS kur (
		id SERIAL PRIMARY KEY,
		market_product_id INT,
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

/* ================= FETCHER ================= */

func startPriceFetcher() {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		fetchAndSave()
		<-ticker.C
	}
}

func fetchAndSave() {
	resp, err := http.Get(apiURL)
	if err != nil {
		log.Println("API hata:", err)
		return
	}
	defer resp.Body.Close()

	var apiResp ApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		log.Println("JSON parse hata:", err)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		return
	}

	stmt, _ := tx.Prepare(`
		INSERT INTO kur 
		(market_product_id, updated_at, customer_buys_at, customer_sells_at)
		VALUES ($1,$2,$3,$4)
	`)
	defer stmt.Close()

	for _, item := range apiResp.Data {
		stmt.Exec(
			item.MarketProductId,
			item.UpdatedAt,
			item.CustomerBuysAt,
			item.CustomerSellsAt,
		)
	}

	tx.Commit()
	log.Println("Kur güncellendi:", time.Now().Format("15:04:05"))
}

/* ================= HTML ================= */

var htmlPage = `
<!DOCTYPE html>
<html lang="tr">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Altın Hesaplama</title>
<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.2/dist/css/bootstrap.min.css" rel="stylesheet">
<script src="https://unpkg.com/vue@3/dist/vue.global.js"></script>
<style>
.btn{min-height:48px;font-size:16px}
</style>
</head>
<body class="bg-light">
<div id="app" class="container py-3"></div>

<script>
const { createApp, ref, watch } = Vue;

createApp({
 setup(){
  const gram=ref(1)
  const productId=ref(0)
  const price=ref(null)
  const ayar=ref({label:'22K',factor:0.916})

  const products=[
   {id:0,name:'Gram'},
   {id:1,name:'Çeyrek'},
   {id:2,name:'Yarım'},
   {id:3,name:'Tam'}
  ]

  const ayarlar=[
   {label:'14K',factor:0.585},
   {label:'18K',factor:0.750},
   {label:'21K',factor:0.875},
   {label:'22K',factor:0.916}
  ]

  const hesapla=async()=>{
   const r=await fetch(
    \`/price?productId=\${productId.value}&gram=\${gram.value}&factor=\${ayar.value.factor}\`
   )
   const d=await r.json()
   price.value=d.price
  }

  watch([gram,productId,ayar],hesapla,{deep:true})

  return{gram,productId,price,ayar,products,ayarlar}
 }
}).mount("#app")
</script>
</body>
</html>
`
