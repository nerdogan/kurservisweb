# 💰 Altın Kur Servisi (Go + Gin + PostgreSQL + Vue 3)

Tek bir **Go binary** içinde çalışan,  
- 📡 Dış API'den **15 saniyede bir** kur verisi çeken  
- 🗄️ PostgreSQL'e otomatik kaydeden  
- 📱 **Mobil-uyumlu** modern web arayüzü sunan  
- 🧮 Satış fiyatına göre altın hesaplayan  

tam entegre bir uygulama.

---

## 🚀 Özellikler

- ✅ Backend: Tek `main.go` dosyası (Go + Gin)
- ✅ Frontend: Vue 3 + Vite + Vue Router
- ✅ Veritabanı: PostgreSQL
- ✅ UI: Bootstrap 5 + Responsive tasarım
- ✅ Mobil-uyumlu arayüz
- ✅ `.env` ile yapılandırma
- ✅ Dış API'den otomatik kur çekme (15 saniye aralığı)
- ✅ Tablo yoksa otomatik oluşturma
- ✅ Satış fiyatı (`customerSellsAt`) bazlı hesaplama
- ✅ Pinia ile state management
- ✅ Axios ile API iletişimi

---

## 🧱 Proje Yapısı

```
kurservisweb/
├─ main.go                    # Go backend
├─ frontend/                  # Vue 3 + Vite uygulaması
│  ├─ src/
│  │  ├─ components/          # Vue bileşenleri
│  │  ├─ views/              # Sayfa bileşenleri
│  │  │  ├─ HomeView.vue     # USD fiyatları
│  │  │  ├─ HomeVieweur.vue  # EUR fiyatları
│  │  │  ├─ HomeViewtl.vue   # TL fiyatları
│  │  │  └─ HomeViewhas.vue  # HAS fiyatları
│  │  ├─ router/             # Vue Router yapılandırması
│  │  ├─ assets/             # CSS ve statik dosyalar
│  │  ├─ App.vue             # Root bileşen
│  │  └─ main.js             # Entry point
│  ├─ package.json
│  ├─ vite.config.js
│  └─ dist/                  # Build çıktısı
├─ templates/                # Backend HTML şablonları
└─ .env                      # Ortam değişkenleri
```

### Backend Mimarisi (main.go)
- **Gin Web Server**: API ve frontend servisi
- **Kur Fetcher**: 15 saniye aralığıyla dış API'yi sorgulama
- **PostgreSQL**: Altın kur verilerini depolama
- **API Routes**:
  - `GET /` → Frontend (Vue uygulaması)
  - `GET /api/prices` → Güncel kur bilgisi

---

## 📦 Kurulum

### 1️⃣ Gereksinimler

- Go 1.20+
- Node.js 16+ (Frontend için)
- PostgreSQL
- Git

### 2️⃣ Projeyi Klonla

```bash
git clone <repo-url>
cd kurservisweb
```

### 3️⃣ .env Dosyası Oluştur

Backend kök dizinde `.env` dosyası oluşturun:

```env
API_URL=https://api.example.com/prices
DB_DSN=host=localhost user=postgres password=postgres dbname=gold sslmode=disable
PORT=8080
```

### 4️⃣ Backend Bağımlılıklarını Yükle

```bash
go mod download
```

### 5️⃣ Frontend Bağımlılıklarını Yükle

```bash
cd frontend
npm install
cd ..
```

### 6️⃣ Uygulamayı Çalıştır

**Seçenek A: Geliştirme Modu (Vite dev server)**

```bash
# Terminal 1: Backend
go run main.go

# Terminal 2: Frontend dev server
cd frontend
npm run dev
```

- Backend: `http://localhost:8080`
- Frontend: `http://localhost:5173` (Vite dev server)

**Seçenek B: Production Build**

```bash
# Frontend'i build et
cd frontend
npm run build
cd ..

# Backend'i çalıştır (frontend dist/ klasörü otomatik serve edilir)
go run main.go
```

Tarayıcıdan aç: `http://localhost:8080`

---

## 📊 PostgreSQL Tablosu

Uygulama otomatik olarak aşağıdaki tabloyu oluşturur:

```sql
CREATE TABLE kur (
    id SERIAL PRIMARY KEY,
    market_product_id INT,
    updated_at TIMESTAMP,
    customer_buys_at NUMERIC(18,5),
    customer_sells_at NUMERIC(18,5),
    created_at TIMESTAMP DEFAULT NOW()
);
```

---

## 🔁 Kur Çekme Mekanizması

- ⏱️ **Aralık**: Her 15 saniyede bir
- 🌐 **API**: `.env` içindeki `API_URL` adresine istek gönderilir
- 📥 **JSON Alanları**: 
  - `marketProductId`
  - `updatedAt`
  - `customerBuysAt`
  - `customerSellsAt`
- 💾 **Depolama**: PostgreSQL kur tablosuna kaydedilir

---

## 🖥️ Frontend Sayfaları

| Rota | Açıklama |
|------|----------|
| `/` | USD cinsinden altın fiyatları |
| `/eur` | EUR cinsinden altın fiyatları |
| `/tl` | TL cinsinden altın fiyatları |
| `/has` | HAS cinsinden altın fiyatları |

---

## 🛠️ Geliştirme

### Frontend Geliştirme Sunucusu

```bash
cd frontend
npm run dev
```

### Frontend Build

```bash
cd frontend
npm run build
```

Output: `frontend/dist/`

### Frontend Preview

```bash
cd frontend
npm run preview
```

---

## 📦 Kullanılan Kütüphaneler

### Backend (Go)
- **Gin**: Web framework
- **lib/pq**: PostgreSQL driver
- **godotenv**: .env dosyası yönetimi

### Frontend (Vue 3)
- **Vue 3**: Progressive JavaScript framework
- **Vue Router**: Client-side routing
- **Pinia**: State management
- **Bootstrap 5**: CSS framework
- **Axios**: HTTP client
- **Vite**: Build tool ve dev server

---

## 📱 Mobil Uyumluluk

- ✅ Responsive tasarım
- ✅ Meta viewport yapılandırma
- ✅ Bootstrap 5 grid sistemi
- ✅ Touch-friendly arayüz
- ✅ Dark mode desteği

---

## 👨‍💻 Geliştirici

Namık ERDOĞAN - 2024
