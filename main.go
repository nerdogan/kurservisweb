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
	resp, err := http.Get(apiURL)
	if err != nil {
		fmt.Println("API error:", err)
		return
	}
	defer resp.Body.Close()

	var apiResp ApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		fmt.Println("JSON decode error:", err)
		return
	}

	limit := 8
	if len(apiResp.Data) < 8 {
		limit = len(apiResp.Data)
	}

	for i := 0; i < limit; i++ {
		item := apiResp.Data[i]
		layout := "2006-01-02T15:04:05.999999"

		updatedAt, _ := time.Parse(layout, item.UpdatedAt)

		_, err := db.Exec(`
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
			fmt.Println("DB insert error:", err)
		}
	}
}

/* ================= HTML ================= */

var htmlPage = `
<!DOCTYPE html>
<html lang="tr">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Altın & Gümüş Hesaplama</title>

<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.2/dist/css/bootstrap.min.css" rel="stylesheet">
<script src="https://unpkg.com/vue@3/dist/vue.global.js"></script>

<style>
body { background:#f8f9fa; }
.btn { min-height:48px; font-size:16px; }
</style>
</head>

<body>
<div id="app" class="container py-3">

  <div class="card shadow-sm">
    <div class="card-body">

      <h5 class="text-center mb-3">Altın Fiyat Hesaplama</h5>

      <div class="mb-3">
        <div class="btn-group w-100">
          <button v-for="p in goldProducts"
            class="btn"
            :class="goldProductId === p.id ? 'btn-primary' : 'btn-outline-primary'"
            @click="goldProductId = p.id">
            {{ p.name }}
          </button>
        </div>
      </div>

      <div class="mb-3">
        <input type="number" class="form-control form-control-lg"
          v-model.number="goldGram" placeholder="Gram">
      </div>

      <div class="mb-3">
        <div class="btn-group w-100">
          <button v-for="a in goldAyarlar"
            class="btn"
            :class="goldAyar.label === a.label ? 'btn-success' : 'btn-outline-success'"
            @click="goldAyar = a">
            {{ a.label }}
          </button>
        </div>
      </div>

      <div v-if="goldPrice" class="alert alert-success text-center fs-5">
        {{ goldPrice }}
      </div>

      <hr>

      <h5 class="text-center mb-3">Gümüş Fiyat Hesaplama</h5>

      <div class="mb-3">
        <div class="btn-group w-100">
          <button v-for="p in silverProducts"
            class="btn"
            :class="silverProductId === p.id ? 'btn-primary' : 'btn-outline-primary'"
            @click="silverProductId = p.id">
            {{ p.name }}
          </button>
        </div>
      </div>

      <div class="mb-3">
        <input type="number" class="form-control form-control-lg"
          v-model.number="silverGram" placeholder="Gram">
      </div>

      <div class="mb-3">
        <div class="btn-group w-100">
          <button v-for="a in silverAyarlar"
            class="btn"
            :class="silverAyar.label === a.label ? 'btn-success' : 'btn-outline-success'"
            @click="silverAyar = a">
            {{ a.label }}
          </button>
        </div>
      </div>

      <div v-if="silverPrice" class="alert alert-success text-center fs-5">
        {{ silverPrice }}
      </div>

    </div>
  </div>
</div>

<script>
const { createApp, ref, watch } = Vue;

createApp({
  setup() {

    const goldGram = ref(1);
    const goldProductId = ref(1);
    const goldAyar = ref({ label: '14K', factor: 0.585 });
    const goldPrice = ref(null);

    const silverGram = ref(1);
    const silverProductId = ref(6);
    const silverAyar = ref({ label: '935', factor: 0.935 });
    const silverPrice = ref(null);

    const goldProducts = [
      { id: 1, name: 'Altın TL' },
      { id: 3, name: 'Altın USD' },
      { id: 4, name: 'Altın EUR' }
    ];

    const silverProducts = [
      { id: 6, name: 'Gümüş TL' },
      { id: 7, name: 'Gümüş USD' },
      { id: 8, name: 'Gümüş EUR' }
    ];

    const goldAyarlar = [
      { label: '14K', factor: 0.585 },
      { label: '18K', factor: 0.750 },
      { label: '21K', factor: 0.875 },
      { label: '22K', factor: 0.916 }
    ];

    const silverAyarlar = [
      { label: '935', factor: 0.935 },
      { label: '1000', factor: 1.0 }
    ];

    const hesapla = async function(productId, gram, factor, target) {
      if (gram <= 0) return;
      const r = await fetch(
        "/price?productId=" + productId +
        "&gram=" + gram +
        "&factor=" + factor
      );
      const d = await r.json();
      target.value = d.price;
    };

    watch([goldGram, goldProductId, goldAyar], function() {
      hesapla(goldProductId.value, goldGram.value, goldAyar.value.factor, goldPrice);
    }, { deep: true });

    watch([silverGram, silverProductId, silverAyar], function() {
      hesapla(silverProductId.value, silverGram.value, silverAyar.value.factor, silverPrice);
    }, { deep: true });

    return {
      goldGram, goldProductId, goldAyar, goldPrice,
      silverGram, silverProductId, silverAyar, silverPrice,
      goldProducts, goldAyarlar,
      silverProducts, silverAyarlar
    };
  }
}).mount("#app");
</script>

</body>
</html>
`
