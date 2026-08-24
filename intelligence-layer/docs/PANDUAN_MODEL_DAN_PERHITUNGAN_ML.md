# 🌸 chia.florist: Panduan Model ML & Hitung-Hitungan Kasus Nyata (Step-by-Step)

> **Tujuan Dokumen**: Dokumen ini ngebahas konsep dasar `Tree-Based Model`, letak implementasi kodenya di folder `intelligence-layer`, hitung-hitungan step-by-step dari **1 sampel data**, serta alasan kenapa model pilihan di project ini lebih unggul dibanding alternatif `Machine Learning` klasik lainnya.

## 🌳 Cara Kerja Model `Tree-Based` (Logika Main Tebak-Tebakan "20 Pertanyaan")

Sebelum masuk ke masing-masing fitur, mari pahami dulu kenapa model bertipe `Tree-Based` jadi raja untuk data tabular (tabel database seperti pesanan, inventaris, dan pengiriman):

```
                       [Apakah Stok <= 5?]
                           /         \
                     (Ya) /           \ (Tidak)
                         ▼             ▼
       [Lead Time Supplier > 3 Hari?]   [Prediksi: AMAN]
                 /            \
           (Ya) /              \ (Tidak)
               ▼                ▼
     [Prediksi: BAHAYA]   [Prediksi: WASPADA]
```

1. **Memotong Data Pakai Pertanyaan Ya/Tidak (`Binary Splitting`)**:  
   Pohon keputusan (`Decision Tree`) bekerja persis seperti manusia saat mengambil keputusan. Dia mencari satu kolom dan satu angka batas (`Split Threshold`) yang paling rapi memisahkan data bagus dan data bermasalah.
2. **Tidak Perlu Rumus Garis Lurus (`Non-Linear Relationships`)**:  
   Model linear memaksa hubungan data harus lurus. Padahal di dunia nyata, penjualan bunga melonjak drastis saat Valentine dan mendatar di hari biasa. `Tree-Based Model` bisa menangkap pola tangga dan lonjakan ini tanpa perlu trik matematika yang rumit.
3. **Evolusi dari 1 Pohon ke `XGBoost`**:
   * **`Single Decision Tree`**: 1 pohon tunggal. Cepat, tapi gampang ngafalin data latih sampai bikin salah tebak di data baru (`High Variance / Overfitting`).
   * **`Random Forest` (`Bagging`)**: Bikin 100 pohon secara independen dari potongan data acak, lalu ambil suara terbanyak (`Voting / Averaging`). Menghilangkan sifat labil dari 1 pohon.
   * **`Gradient Boosting` (`Boosting`)**: Pohon dibuat berurutan (estafet). Pohon ke-2 belajar dari sisa kesalahan (`Residuals`) pohon ke-1, pohon ke-3 belajar dari kesalahan pohon ke-2, dan seterusnya.
   * **`XGBoost` (`Extreme Gradient Boosting`)**: Versi modern dari `Gradient Boosting`. Memakai turunan kedua (`Hessian / Curvature`) untuk menghitung bobot daun optimal dalam hitungan milidetik, plus dilengkapi rem otomatis (`Built-in L1/L2 Regularization`) agar pohon tidak tumbuh liar.

---

## ⚔️ Perbandingan Filosofis: `XGBoost` vs `Isolation Forest` (Dua Pohon dengan Tujuan Berbeda)

Walaupun sama-sama berbentuk struktur pohon bercabang (`Tree`), rumus dan proses di balik keduanya bekerja dengan filosofi yang 180 derajat bertolak belakang:

```
          XGBoost: "Penebangan Terpandu Error"        Isolation Forest: "Tebasan Samurai Acak"
          (Supervised: Cari Split Terbaik)            (Unsupervised: Cari Titik Paling Terpencil)

                 [Split Optimal: Gain Max]                     [Split Acak: Fitur & Nilai Acak]
                        /        \                                    /        \
                       /          \                                  /          \
                Leaf: Nilai w*   Leaf: Nilai w*               Normal Data       Anomali
               (Angka Koreksi)  (Angka Koreksi)             (Kedalaman > 10)  (Kedalaman = 1-2)
```

### 1. `XGBoost`: Dipandu oleh Gradien Error
* **Tujuan**: Memprediksi angka target $y$ seakurat mungkin.
* **Proses Splitting**: Sangat selektif. Menggunakan rumus `Split Gain` untuk mencari titik potong yang menghasilkan penurunan error terbesar.
* **Hasil Akhir Daun**: Mengeluarkan **angka koreksi bobot $w^*$** ($-\frac{\sum g_i}{\sum h_i + \lambda}$) yang ditambahkan ke tebakan model sebelumnya.

### 2. `Isolation Forest`: Tebasan Acak untuk Mengukur Kerapatan Ruang
* **Tujuan**: Memisahkan satu titik dari populasi tanpa perlu tahu target jawaban $y$.
* **Proses Splitting**: Murni acak (`Uniform Random`). Memilih kolom dan angka potong tanpa mempedulikan gradien atau fungsi loss sama sekali.
* **Hasil Akhir Daun**: **Bukan angka prediksi**, melainkan hanya menghitung **kedalaman potongan ($h(x)$)**:
  * Titik normal yang berkerumun padat butuh banyak tebasan acak ($h(x)$ dalam $\to$ Skor anomali $\approx 0.0$).
  * Titik anomali yang berdiri sendirian di ujung ekstrem langsung terisolasi dalam 1-2 tebasan ($h(x)$ dangkal $\to$ Skor anomali $\approx 1.0$).

| Aspek | `XGBoost` | `Isolation Forest` |
|---|---|---|
| **Tipe Pembelajaran** | `Supervised Learning` (Wajib ada label $y$) | `Unsupervised Learning` (Hanya butuh input $\mathbf{x}$) |
| **Penentu Titik Potong** | Optimasi matematika turunan ke-1 (`Gradient`) & ke-2 (`Hessian`) | Pemilihan acak murni (`Uniform Random Selection`) |
| **Output Daun Pohon** | Nilai bobot numerik penambah prediksi ($w^*$) | Jumlah langkah kedalaman dari akar ($h(x)$) |
| **Fungsi Matematika Akhir** | Penjumlahan aditif berantai $F_0 + \sum \eta f_m(x)$ | Fungsi peluruhan eksponensial $2^{-\frac{\mathbb{E}(h(x))}{c(n)}}$ |

---

## 1. Prediksi Permintaan & Penjualan (Demand Forecasting)


### 📍 Letak File & Arsitektur di `intelligence-layer`:
* **Trainer Script**: [`src/forecasting_trainer.py`](file:///d:/__Projects/kage/chia.florist/intelligence-layer/src/forecasting_trainer.py) (Class: `DemandForecasterTrainer`, model: `XGBRegressor`)
* **Feature Engineering**: [`src/feature_engineering.py`](file:///d:/__Projects/kage/chia.florist/intelligence-layer/src/feature_engineering.py) (Class: `TimeSeriesFeatureBuilder`)
* **Model Checkpoint**: `models/demand_forecasting_model.json`
* **Inference Endpoint**: `POST /api/v1/forecast/demand` di [`app/api/v1/forecast.py`](file:///d:/__Projects/kage/chia.florist/intelligence-layer/app/api/v1/forecast.py)
* **Prediction Service**: `PredictionService.predict_demand()` di [`app/services/predictor.py`](file:///d:/__Projects/kage/chia.florist/intelligence-layer/app/services/predictor.py)

### Skenario: Stok Mawar untuk Hari Valentine
Minggu depan sudah masuk Hari Valentine. Manajer toko bunga butuh estimasi berapa buket "Mawar Merah Passion" (`PROD-001`) yang bakal laku dalam 7 hari ke depan. Tujuannya agar bisa pesan tangkai mawar segar ke kebun bunga di Bandung secara pas, tidak kekurangan stok dan tidak kebanyakan sampai bunganya busuk di gudang.

### Model yang Dipakai di Project: `XGBoost Regressor` (Perhitungan Bobot Daun)

`XGBoost` membuat tebakan baseline awal, melihat sisa error pada sampel kita, lalu memanfaatkan arah kemiringan (`Gradient`) dan kelengkungan (`Hessian`) untuk menghitung nilai koreksi terbaik pada daun pohon keputusan.

$$
w^* = -\frac{g_i}{h_i + \lambda}
$$

Diketahui:
* $g_i = -6.0$ (Gradien ke-1: tebakan awal model masih kurang 6 buket dari kenyataan)
* $h_i = 1.0$ (Gradien ke-2 / Hessian untuk Squared Error Loss)
* $\lambda = 1.0$ (Penalti $L_2$ Regularization agar koreksi tidak melompat terlalu ekstrem)
* $\hat{y}^{(0)} = 20.0$ (Tebakan baseline awal dalam jumlah buket)
* $\eta = 0.1$ (Learning Rate / Shrinkage)

$$
w^* = -\frac{-6.0}{1.0 + 1.0}
$$

$$
w^* = \frac{6.0}{2.0} = 3.0
$$

$$
\hat{y}^{(1)} = 20.0 + (0.1 \times 3.0)
$$

$$
\hat{y}^{(1)} = 20.0 + 0.3 = 20.3
$$

Dengan demikian, hasil tebakan penjualan mawar yang sudah dikoreksi adalah **20.3 buket terjual**.

### Alternatif ML Klasik: `Ridge Regression` (Linear Weighted Sum)

Model linear biasa mengalikan setiap fitur data toko dengan bobot kepentingannya masing-masing, lalu ditambahkan nilai intercept dasar.

$$
\hat{y} = w_1 x_1 + w_2 x_2 + b
$$

Diketahui:
* $x_1 = 25.0$ (Jumlah buket terjual dalam 7 hari kemarin)
* $x_2 = 120.0$ (Jumlah view katalog web minggu ini)
* $w_1 = 0.65$ (Bobot pengaruh riwayat penjualan)
* $w_2 = 0.05$ (Bobot pengaruh pengunjung web)
* $b = 4.0$ (Penjualan dasar toko secara organik)

$$
\hat{y} = (0.65 \times 25.0) + (0.05 \times 120.0) + 4.0
$$

$$
\hat{y} = 16.25 + 6.00 + 4.00
$$

$$
\hat{y} = 26.25
$$

Dengan demikian, hasil estimasi permintaan dari model linear adalah **26.25 buket terjual**.

### 💡 Mengapa Memilih `XGBoost Regressor` Dibanding `Ridge Regression`?
1. **Mampu Menangkap Pola Non-Linear**: `Ridge Regression` memaksakan pola garis lurus. Jika pengunjung web naik 10x lipat saat Valentine, model linear mengira penjualan naik 10x lipat secara datar, padahal kapasitas produksi bunga toko ada batas maksimalnya (`Cap`). `XGBoost` dengan mudah membuat batas cabang atas ini.
2. **Kombinasi Fitur Otomatis (`Feature Interactions`)**: `XGBoost` bisa otomatis mengenali pola kombinasi seperti: *"HANYA JIKA hari = Sabtu/Minggu DAN event = Valentine, MAKA penjualan meledak"*. Pada `Ridge Regression`, kombinasi perkalian ini harus dibuat manual satu per satu lewat matematika tambahan.

## 2. Risiko Kehabisan Stok (Inventory Stockout Risk)

### 📍 Letak File & Arsitektur di `intelligence-layer`:
* **Trainer Script**: [`src/stockout_trainer.py`](file:///d:/__Projects/kage/chia.florist/intelligence-layer/src/stockout_trainer.py) (Class: `StockoutRiskTrainer`, model: `XGBClassifier`)
* **Feature Engineering**: [`src/feature_engineering.py`](file:///d:/__Projects/kage/chia.florist/intelligence-layer/src/feature_engineering.py) (Class: `InventoryStockoutFeatureBuilder`)
* **Model Checkpoint**: `models/stockout_risk_model.json`
* **Inference Endpoint**: `POST /api/v1/inventory/stockout-risk` di [`app/api/v1/inventory.py`](file:///d:/__Projects/kage/chia.florist/intelligence-layer/app/api/v1/inventory.py)
* **Prediction Service**: `PredictionService.predict_stockout_risk()` di [`app/services/predictor.py`](file:///d:/__Projects/kage/chia.florist/intelligence-layer/app/services/predictor.py)

### Skenario: Stok Lili Putih Menipis Menjelang Akhir Pekan
Toko sedang memegang sisa 12 buket Bunga Lili Putih di etalase. Rata-rata penjualan harian adalah 3 buket/hari, sedangkan supplier kebun butuh waktu 5 hari untuk mengirimkan stok baru. Kita butuh skor probabilitas 0.0 sampai 1.0 untuk menentukan apakah sistem inventaris harus memicu peringatan darurat pembelian ulang (`Emergency Purchase Order`).

### Model yang Dipakai di Project: `XGBoost Classifier` (Konversi Skor ke Probabilitas via Sigmoid)

Pohon keputusan klasifikasi menghitung skor margin mentah berdasarkan rasio stok dan kecepatan jual, lalu dikonversi menggunakan fungsi `Sigmoid` menjadi probabilitas valid antara 0% hingga 100%.

$$
p = \frac{1}{1 + e^{-z}}
$$

Diketahui:
* $z = 1.3863$ (Skor logit mentah dari daun pohon keputusan untuk produk Lili ini)
* $e = 2.71828$ (Konstanta matematika Euler)

$$
p = \frac{1}{1 + (2.71828)^{-1.3863}}
$$

$$
p = \frac{1}{1 + 0.2500}
$$

$$
p = \frac{1}{1.2500} = 0.8000
$$

Dengan demikian, peluang risiko kehabisan stok adalah **80.00% (CRITICAL RISK / BUTUH RESTOCK)**.

### Alternatif ML Klasik: Evaluasi Loss `Logistic Regression` (Log-Loss)

Ketika model sedang dievaluasi pada data latih di mana stok toko ternyata benar-benar habis (label aktual $y = 1$), kualitas prediksi dihitung menggunakan `Binary Cross-Entropy Loss`.

$$
\mathcal{L} = -\left[ y \ln(p) + (1 - y) \ln(1 - p) \right]
$$

Diketahui:
* $y = 1.0$ (Kondisi riil: toko benar-benar kehabisan stok)
* $p = 0.80$ (Model memprediksi probabilitas 80%)
* $\ln(0.80) = -0.2231$

$$
\mathcal{L} = -\left[ (1.0 \times \ln(0.80)) + (0.0 \times \ln(0.20)) \right]
$$

$$
\mathcal{L} = -(-0.2231)
$$

$$
\mathcal{L} = 0.2231
$$

Dengan demikian, penalti error latihan pada data ini adalah **0.2231**.

### 💡 Mengapa Memilih `XGBoost Classifier` Dibanding `Logistic Regression`?
1. **Risiko Inventaris Bekerja Berdasarkan Batas Ambang Kritis (`Threshold-Based`)**: Jika sisa stok masih 10 buket, risiko kehabisan stok mendekati 0%. Tapi begitu stok turun melewati angka 3 buket, risikonya melonjak drastis ke 90%. `XGBoost` sangat ahli memotong batas angka kritis (`Split at Stock <= 3`), sedangkan `Logistic Regression` mencoba meratakannya dalam kurva probabilitas yang landai.
2. **Kekebalan terhadap Skala Angka Berbeda (`Scale Invariance`)**: `XGBoost` tidak peduli apakah fitur harga bernilai jutaan rupiah sedangkan sisa stok bernilai satuan kecil (3 pcs). `Logistic Regression` wajib melakukan normalisasi data (`StandardScaler`) terlebih dahulu agar perhitungannya tidak rusak.

## 3. Deteksi Anomali Operasional (Operational Anomaly Detection)

### 📍 Letak File & Arsitektur di `intelligence-layer`:
* **Trainer Script**: [`src/anomaly_trainer.py`](file:///d:/__Projects/kage/chia.florist/intelligence-layer/src/anomaly_trainer.py) (Class: `OperationalAnomalyTrainer`, model: `IsolationForest`)
* **Feature Engineering**: [`src/feature_engineering.py`](file:///d:/__Projects/kage/chia.florist/intelligence-layer/src/feature_engineering.py) (Class: `OperationalAnomalyFeatureBuilder`)
* **Model Checkpoint**: `models/anomaly_detector/isolation_forest.pkl` & `scaler.pkl`
* **Inference Endpoint**: `POST /api/v1/anomalies/detect` di [`app/api/v1/anomalies.py`](file:///d:/__Projects/kage/chia.florist/intelligence-layer/app/api/v1/anomalies.py)
* **Prediction Service**: `PredictionService.detect_anomaly()` di [`app/services/predictor.py`](file:///d:/__Projects/kage/chia.florist/intelligence-layer/app/services/predictor.py)

### Skenario: Orderan Nilai Besar dengan Keterlambatan Pembayaran 24 Jam
Ada pesanan rangkaian bunga dekorasi pernikahan mewah senilai Rp 3.500.000 via Transfer Bank Manual. Rata-rata pelanggan lain mengunggah bukti transfer dalam 15 menit (900 detik), namun orderan ini baru masuk buktinya setelah 86.400 detik (24 jam). Sistem perlu mendeteksi apakah ini anomali operasional / kegagalan webhook payment gateway.

### Model yang Dipakai di Project: `Isolation Forest` (Skor Kedalaman Partisi Acak)

`Isolation Forest` bekerja dengan membuat partisi acak pada fitur data. Titik anomali yang lokasinya terpencil di ruang data akan terisolasi sangat cepat pada kedalaman pohon yang sangat dangkal ($h(x)$ kecil).

$$
s(x, n) = 2^{-\frac{h(x)}{c(n)}}
$$

Diketahui:
* $h(x) = 2.0$ (Jumlah potongan cabang acak sampai transaksi ini terisolasi sendirian)
* $n = 256$ (Jumlah sampel data yang diambil saat membangun pohon)
* $c(256) = 10.23$ (Rata-rata kedalaman normal pada struktur Binary Search Tree)

$$
s = 2^{-\frac{2.0}{10.23}}
$$

$$
s = 2^{-0.1955}
$$

$$
s = 0.8732
$$

Dengan demikian, skor anomali transaksi ini adalah **0.8732 (HIGH ANOMALY / PERLU AUDIT)**.

### Alternatif ML Klasik: `Z-Score` Statistik (Jarak Deviasi Standar)

Metode statistik sederhana mengukur berapa standar deviasi jarak nilai data tertentu dari rata-rata populasi normal.

$$
z = \frac{x - \mu}{\sigma}
$$

Diketahui:
* $x = 86400.0$ (Durasi penyelesaian transaksi dalam detik)
* $\mu = 1800.0$ (Rata-rata durasi normal transaksi dalam detik)
* $\sigma = 7200.0$ (Standar deviasi riwayat transaksi dalam detik)

$$
z = \frac{86400.0 - 1800.0}{7200.0}
$$

$$
z = \frac{84600.0}{7200.0}
$$

$$
z = 11.75
$$

Dengan demikian, jarak penyimpangannya adalah **11.75 standar deviasi (EXTREME OUTLIER)**.

### 💡 Mengapa Memilih `Isolation Forest` Dibanding `Z-Score`?
1. **Deteksi Multidimensi (`Multi-Dimensional Isolation`)**: `Z-Score` hanya bisa memeriksa 1 kolom saja dalam satu waktu (misal hanya durasi atau hanya harga). `Isolation Forest` sanggup mendeteksi kombinasi aneh sekaligus: misal order Rp 50.000 tapi durasinya 24 jam, atau order Rp 50.000.000 tapi selesai dalam 2 detik.
2. **Bebas Asumsi Distribusi Data**: `Z-Score` mengasumsikan data terdistribusi normal simetris (`Bell Curve / Gaussian`). Kenyataannya, data waktu e-commerce sangat menceng (`Skewed`). `Isolation Forest` tidak membutuhkan asumsi distribusi kurva sama sekali (`Non-Parametric`).

## 4. Estimasi Durasi Pengiriman Kurir & SLA (Courier SLA Estimator)

### 📍 Letak File & Arsitektur di `intelligence-layer`:
* **Trainer Script**: [`src/courier_sla_trainer.py`](file:///d:/__Projects/kage/chia.florist/intelligence-layer/src/courier_sla_trainer.py) (Class: `CourierSLATrainer`, model: `GradientBoostingRegressor`)
* **Feature Engineering**: [`src/feature_engineering.py`](file:///d:/__Projects/kage/chia.florist/intelligence-layer/src/feature_engineering.py) (Class: `CourierSLAFeatureBuilder`)
* **Model Checkpoint**: `models/courier_sla_model.json`
* **Inference Endpoint**: `POST /api/v1/courier/sla-prediction` di [`app/api/v1/courier.py`](file:///d:/__Projects/kage/chia.florist/intelligence-layer/app/api/v1/courier.py)
* **Prediction Service**: `PredictionService.predict_courier_sla()` di [`app/services/predictor.py`](file:///d:/__Projects/kage/chia.florist/intelligence-layer/app/services/predictor.py)

### Skenario: Pengiriman Hari Jumat Sore Menggunakan SiCepat
Pelanggan membeli buket bunga segar untuk dikirim antarkota Jakarta menggunakan kurir SiCepat pada hari Jumat jam 15:00 WIB (ongkir Rp 15.000). Toko bunga perlu mengestimasi durasi perjalanan kurir agar bunga tidak sampai layu sebelum diterima pelanggan.

### Model yang Dipakai di Project: `Gradient Boosting Regressor` (Pembaruan Residu Bertahap)

`Gradient Boosting` memulai dari rata-rata durasi dasar pengiriman kota ($F_0$), lalu setiap pohon secara berurutan menambahkan koreksi waktu spesifik berdasarkan kurir dan jam sibuk hari kerja.

$$
F_1(x) = F_0 + \eta \cdot f_1(x)
$$

Diketahui:
* $F_0 = 24.0$ (Durasi dasar rata-rata pengiriman sekota dalam jam)
* $\eta = 0.1$ (Learning Rate / Shrinkage)
* $f_1(x) = -60.0$ (Nilai koreksi pohon untuk kurir ekspres motor SiCepat)

$$
F_1(x) = 24.0 + (0.1 \times -60.0)
$$

$$
F_1(x) = 24.0 - 6.0
$$

$$
F_1(x) = 18.0
$$

Dengan demikian, estimasi durasi pengiriman bunga adalah **18.0 jam (ON TRACK / SEGAR)**.

### Alternatif ML Klasik: `Decision Tree Regressor` (Rata-rata Kelompok Tunggal)

Satu pohon keputusan tunggal mengelompokkan pesanan ini ke daun kriteria `courier == SICEPAT` dan `is_weekend == 0`, lalu menghitung nilai rata-rata sampel masa lalu yang jatuh pada kelompok tersebut.

$$
\bar{y} = \frac{1}{N} \sum_{i=1}^N y_i
$$

Diketahui:
* $N = 3$ (Ada 3 data pengiriman SiCepat di jam sibuk hari kerja)
* $y_1 = 16.0$ (Data historis 1: 16 jam)
* $y_2 = 18.0$ (Data historis 2: 18 jam)
* $y_3 = 20.0$ (Data historis 3: 20 jam)

$$
\bar{y} = \frac{16.0 + 18.0 + 20.0}{3}
$$

$$
\bar{y} = \frac{54.0}{3}
$$

$$
\bar{y} = 18.0
$$

Dengan demikian, hasil estimasi durasi rata-rata kelompok adalah **18.0 jam**.

### 💡 Mengapa Memilih `Gradient Boosting Regressor` Dibanding `Single Decision Tree`?
1. **Hasil Estimasi Lebih Halus dan Akurat**: `Single Decision Tree` menghasilkan prediksi berupa tangga kasar (semua order SiCepat sore hari akan ditebak sama persis 18.0 jam). `Gradient Boosting` menggabungkan puluhan pohon koreksi kecil sehingga mampu memperhalus prediksi sesuai variasi ongkir, jarak km, dan cuaca.
2. **Mencegah Overfitting pada Data Baru**: Pohon tunggal sangat mudah terpengaruh oleh 1 atau 2 data aneh di masa lalu. `Gradient Boosting` dengan parameter `learning_rate = 0.1` membatasi agar tidak ada satu pohon pun yang mendominasi keputusan akhir.

## 5. Cara Menghitung Error Tebakan Model pada 1 Sampel

### Skenario: Evaluasi Akurasi Prediksi Penjualan Toko
Model memprediksi toko akan menjual **25 buket** Bunga Matahari hari ini. Di penghujung hari, data kasir mencatat penjualan riil sebanyak **29 buket**. Kita ingin mengukur tingkat kesalahan tebakan pada hari tersebut.

### `Mean Absolute Error` (MAE) pada 1 Data

$$
e_{\text{abs}} = |y - \hat{y}|
$$

Diketahui:
* $y = 29.0$ (Penjualan aktual di kasir)
* $\hat{y} = 25.0$ (Prediksi model)

$$
e_{\text{abs}} = |29.0 - 25.0|
$$

$$
e_{\text{abs}} = |4.0|
$$

$$
e_{\text{abs}} = 4.0
$$

Dengan demikian, selisih error absolutnya adalah **4.0 buket**.

### `Mean Squared Error` (MSE) pada 1 Data

$$
e_{\text{sq}} = (y - \hat{y})^2
$$

Diketahui:
* $y = 29.0$ (Penjualan aktual di kasir)
* $\hat{y} = 25.0$ (Prediksi model)

$$
e_{\text{sq}} = (29.0 - 25.0)^2
$$

$$
e_{\text{sq}} = (4.0)^2
$$

$$
e_{\text{sq}} = 16.0
$$

Dengan demikian, nilai error kuadratnya adalah **16.0**.
