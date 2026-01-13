# 💰 Gold Price Service (Go + Gin + PostgreSQL)

Tek bir **Go binary** içinde çalışan,  
- 📡 Dış API’den **15 saniyede bir** kur verisi çeken  
- 🗄️ PostgreSQL’e otomatik kaydeden  
- 📱 **Mobil-first** web arayüzü sunan  
- 🧮 Satış fiyatına göre altın hesaplayan  

tam entegre bir uygulama.

---

## 🚀 Özellikler

- ✅ Tek `main.go` dosyası
- ✅ Go + Gin
- ✅ PostgreSQL
- ✅ Vue 3 (CDN) + Bootstrap 5
- ✅ Mobil-first UI
- ✅ `.env` ile yapılandırma
- ✅ Dış API’den otomatik kur çekme
- ✅ Tablo yoksa otomatik oluşturma
- ✅ Satış fiyatı (`customerSellsAt`) bazlı hesaplama

---

## 🧱 Mimari

main.go
├─ Gin Web Server
│ ├─ / → Mobil UI (Vue + Bootstrap)
│ └─ /price → Fiyat API
│
├─ Kur Fetcher (15 sn)
│ └─ External API
│
└─ PostgreSQL
└─ kur tablosu


---

## 📦 Kurulum

### 1️⃣ Gereksinimler

- Go 1.20+
- PostgreSQL
- Git

---

### 2️⃣ Projeyi Klonla

```bash
git clone https://github.com/kullanici/gold-price-service.git
cd gold-price-service

3️⃣ .env Dosyası Oluştur
API_URL=https://api.ornek.com/prices
DB_DSN=host=localhost user=postgres password=postgres dbname=gold sslmode=disable

4️⃣ Bağımlılıkları Yükle
go mod init gold-price-service
go get github.com/gin-gonic/gin
go get github.com/lib/pq
go get github.com/joho/godotenv

5️⃣ Çalıştır
go run main.go


Tarayıcıdan aç:

http://localhost:8080

📊 PostgreSQL Tablosu

Uygulama otomatik olarak aşağıdaki tabloyu oluşturur:

CREATE TABLE kur (
    id SERIAL PRIMARY KEY,
    market_product_id INT,
    updated_at TIMESTAMP,
    customer_buys_at NUMERIC(18,5),
    customer_sells_at NUMERIC(18,5),
    created_at TIMESTAMP DEFAULT NOW()
);

🔁 Kur Çekme Mekanizması

⏱️ Her 15 saniyede bir

🌐 .env içindeki API_URL adresine istek atar

📥 JSON içinden şu alanları alır:

marketProductId

updatedAt

customerBuysAt

customerSellsAt

💾 PostgreSQL kur tablosuna kaydeder

🧮 Fiyat Hesaplama Mantığı
Fiyat = Gram × customerSellsAt × Ayar Katsayısı

Ayar Katsayıları
Ayar	Katsayı
14K	0.585
18K	0.750
21K	0.875
22K	0.916
📱 Mobil UI Özellikleri

Büyük dokunmatik butonlar

Ürün seçimi (Gram, Çeyrek, Yarım, Tam)

Ayar seçimi (14K – 22K)

Otomatik hesaplama

Tek kolon, mobil-first tasarım

🛠️ API Endpoint
GET /price

Query Params

Param	Açıklama
productId	Ürün ID
gram	Gram
factor	Ayar katsayısı

Response

{
  "price": "12345.67"
}

🔒 Hata Yönetimi

Dış API down olsa bile server çalışmaya devam eder

DB hataları loglanır

UI çökmeyecek şekilde tasarlanmıştır
