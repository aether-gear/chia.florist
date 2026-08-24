# 🧠 Filosofi & Arsitektur Matematika Lengkap: Model Machine Learning `chia.florist`

> **Tujuan Dokumen**: Dokumen ini menjelaskan **filosofi visual (mental model)** dari setiap model `Machine Learning` yang digunakan di `chia.florist`, serta menyajikan **Formula Besar Pemersatu (`Master Unifying Formula`)** yang menggabungkan seluruh rumus internal (seperti *Objective Function*, *Gradient*, *Hessian*, *Leaf Weight*, *Split Gain*, dan *Additive Updates*) ke dalam satu siklus matematika yang utuh dan tidak terputus.

## 🏛️ Fondasi Induk: Formula Terbesar dalam Seluruh Dunia Machine Learning

Sebelum melihat rumus masing-masing model, sadari bahwa hampir semua model `Supervised Learning` di dunia (mulai dari Regresi Linear, SVM, sampai XGBoost) tunduk pada satu **Formula Induk Terbesar** yang disebut **Empirical Risk Minimization with Regularization**:

$$
\min_{\Theta} \mathcal{J}(\Theta) = \underbrace{\frac{1}{N} \sum_{i=1}^N \mathcal{L}(y_i, f(\mathbf{x}_i; \Theta))}_{\text{Kualitas Prediksi (Seberapa Cocok dengan Data)}} + \underbrace{\Omega(\Theta)}_{\text{Rem / Regularisasi (Mencegah Overfitting)}}
$$

* Kalau $\mathcal{L}$ memakai *Squared Error* dan $\Omega$ memakai kuadrat bobot $\lambda \|\mathbf{w}\|_2^2 \to$ Lahirlah **`Ridge Regression`**.
* Kalau $\mathcal{L}$ memakai *Log-Loss* dan fungsi $f$ diperas oleh `Sigmoid` $\to$ Lahirlah **`Logistic Regression`**.
* Kalau fungsi $f$ dibangun dari kumpulan pohon bertingkat yang dioptimasi memakai turunan kedua Taylor $\to$ Lahirlah **`XGBoost`**.

## 1. Mesin Inti Terpadu XGBoost (`The Unified XGBoost Engine`)

### 📍 Letak File di Codebase:
* **Demand Forecasting (Regressor)**: [`src/forecasting_trainer.py`](file:///d:/__Projects/kage/chia.florist/intelligence-layer/src/forecasting_trainer.py) $\to$ [`app/api/v1/forecast.py`](file:///d:/__Projects/kage/chia.florist/intelligence-layer/app/api/v1/forecast.py)
* **Stockout Risk (Classifier)**: [`src/stockout_trainer.py`](file:///d:/__Projects/kage/chia.florist/intelligence-layer/src/stockout_trainer.py) $\to$ [`app/api/v1/inventory.py`](file:///d:/__Projects/kage/chia.florist/intelligence-layer/app/api/v1/inventory.py)
* **Feature Engineering**: [`src/feature_engineering.py`](file:///d:/__Projects/kage/chia.florist/intelligence-layer/src/feature_engineering.py)
* **Inference Service**: [`app/services/predictor.py`](file:///d:/__Projects/kage/chia.florist/intelligence-layer/app/services/predictor.py)

---

### 🎨 Filosofi Visual: *"Pahat Patung Bertingkat Mengikuti Kontur Lembah Error"*
Bayangkan kamu sedang memahat balok kayu besar untuk membuat patung yang mirip bentuk grafik penjualan atau risiko kehabisan stok bunga:
* **Pohon ke-1** memotong balok secara kasar menjadi 2 tingkat.
* **Pohon ke-2** melihat bagian mana yang masih ketebalan atau kekurusan (sisa *Error / Residuals*), lalu memahat bagian tersebut lebih halus.
* **Pohon ke-100** tinggal menghaluskan bagian lekukan terkecil.

Model linear hanya bisa menarik garis lurus miring dengan gradien konstan ($y = wx + b$). `XGBoost` menyusun tumpukan balok lego non-linear bertingkat yang sanggup mengikuti lonjakan penjualan hari Valentine atau batas kritis stok bunga setinggi apa pun.

---

### 🔄 Siklus Rumus Besar yang Mengayomi (`Master Mathematical Loop` XGBoost)

Berikut adalah alur utuh bagaimana *Objective Function* induk diturunkan langkah demi langkah sampai menghasilkan rumus bobot daun ($w^*$) dan rumus *Split Gain*:

```
[1. Master Objective]  ──►  [2. Taylor 2nd-Order]  ──►  [3. Leaf Grouping]
         │                                                      │
         ▼                                                      ▼
[6. Prediction Update] ◄──  [5. Split Gain Formula] ◄── [4. Optimal Leaf w*]
```

#### Langkah 1: Master Objective Function pada Iterasi ke-$t$
Pada pohon ke-$t$, kita ingin mencari fungsi pohon baru $f_t(\mathbf{x})$ yang jika ditambahkan ke tebakan lama $\hat{y}_i^{(t-1)}$ akan menghasilkan error terkecil:

$$
\mathcal{L}^{(t)} = \sum_{i=1}^N \mathcal{L}\left(y_i, \hat{y}_i^{(t-1)} + f_t(\mathbf{x}_i)\right) + \gamma T + \frac{1}{2}\lambda \sum_{j=1}^T w_j^2
$$

---

#### Langkah 2: Aproksimasi Kurva Menggunakan Deret Taylor Orde ke-2

> 🔍 **Apa yang disubstitusikan?**  
> Kita mengganti fungsi error umum $\mathcal{L}(y_i, \hat{y}_i^{(t-1)} + f_t(\mathbf{x}_i))$ dengan rumus aproksimasi deret Taylor di sekitar tebakan lama $\hat{y}_i^{(t-1)}$:  
>
> $f(a + \Delta x) \approx f(a) + f'(a)\Delta x + \frac{1}{2}f''(a)(\Delta x)^2$
> Di sini $a = \hat{y}_i^{(t-1)}$ dan $\Delta x = f_t(\mathbf{x}_i)$.

Bentuk setelah substitusi:

$$
\mathcal{L}^{(t)} \approx \sum_{i=1}^N \left[ \mathcal{L}(y_i, \hat{y}_i^{(t-1)}) + g_i f_t(\mathbf{x}_i) + \frac{1}{2} h_i f_t^2(\mathbf{x}_i) \right] + \gamma T + \frac{1}{2}\lambda \sum_{j=1}^T w_j^2
$$

Di mana turunan ke-1 (`Gradient`) dan turunan ke-2 (`Hessian`) adalah:
$$
g_i = \frac{\partial \mathcal{L}(y_i, \hat{y}_i^{(t-1)})}{\partial \hat{y}_i^{(t-1)}}, \quad h_i = \frac{\partial^2 \mathcal{L}(y_i, \hat{y}_i^{(t-1)})}{\partial (\hat{y}_i^{(t-1)})^2}
$$

Karena suku pertama $\mathcal{L}(y_i, \hat{y}_i^{(t-1)})$ adalah konstanta dari masa lalu (tidak terpengaruh pohon baru), suku ini bisa diabaikan saat optimasi.

---

#### Langkah 3: Mengelompokkan Sampel Data Berdasarkan Daun ($j = 1 \dots T$)

> 🔍 **Apa yang disubstitusikan?**  
> Semua sampel data $i$ yang jatuh ke dalam daun yang sama ($i \in I_j$) akan mendapatkan nilai output yang sama persis, yaitu $f_t(\mathbf{x}_i) = w_j$.  
> Kita ganti variabel $f_t(\mathbf{x}_i)$ menjadi $w_j$ dan kelompokkan penjumlahan $\sum_{i=1}^N$ menjadi penjumlahan per daun $\sum_{j=1}^T \sum_{i \in I_j}$.

Bentuk setelah pengelompokan:

$$
\tilde{\mathcal{L}}^{(t)} = \sum_{j=1}^T \left[ \left(\sum_{i \in I_j} g_i\right) w_j + \frac{1}{2} \left(\sum_{i \in I_j} h_i\right) w_j^2 \right] + \gamma T + \frac{1}{2}\lambda \sum_{j=1}^T w_j^2
$$

Gabungkan suku $w_j^2$ di dalam satu kurung:

$$
\tilde{\mathcal{L}}^{(t)} = \sum_{j=1}^T \left[ \left(\sum_{i \in I_j} g_i\right) w_j + \frac{1}{2} \left(\sum_{i \in I_j} h_i + \lambda\right) w_j^2 \right] + \gamma T
$$

Biar ringkas, kita substitusikan simbol $G_j = \sum_{i \in I_j} g_i$ dan $H_j = \sum_{i \in I_j} h_i$:

$$
\tilde{\mathcal{L}}^{(t)} = \sum_{j=1}^T \left[ G_j w_j + \frac{1}{2} (H_j + \lambda) w_j^2 \right] + \gamma T
$$

---

#### Langkah 4: Menurunkan Rumus Bobot Daun Optimal ($w_j^*$)

> 🔍 **Apa yang disubstitusikan?**  
> Persamaan di atas adalah fungsi kuadrat terhadap $w_j$ berbentuk $A w_j + \frac{1}{2} B w_j^2$.  
> Titik minimumnya dicari dengan menurunkan $\tilde{\mathcal{L}}^{(t)}$ terhadap $w_j$ dan menyamakan hasilnya ke nol ($\frac{\partial \tilde{\mathcal{L}}}{\partial w_j} = 0$).

Proses turunan:

$$
\frac{\partial \tilde{\mathcal{L}}^{(t)}}{\partial w_j} = G_j + (H_j + \lambda) w_j = 0
$$

Pindahkan $G_j$ ke ruas kanan dan bagi dengan $(H_j + \lambda)$:

$$
(H_j + \lambda) w_j = -G_j \implies \mathbf{w_j^* = -\frac{G_j}{H_j + \lambda}}
$$

---

#### Langkah 5: Substitusi $w_j^*$ Kembali untuk Menghasilkan Rumus `Split Gain`

> 🔍 **Apa yang disubstitusikan?**  
> Kita mengganti setiap variabel $w_j$ di dalam persamaan $\tilde{\mathcal{L}}^{(t)}$ pada Langkah 3 dengan rumus optimalnya: $w_j^* = -\frac{G_j}{H_j + \lambda}$.

Tulis bentuk mentah substitusi untuk suku di dalam daun $j$:

$$
\text{Suku Daun } j = G_j \left( -\frac{G_j}{H_j + \lambda} \right) + \frac{1}{2} (H_j + \lambda) \left( -\frac{G_j}{H_j + \lambda} \right)^2
$$

Penyederhanaan aljabar langkah demi langkah:
1. Kuadratkan pecahan kedua: $\left(-\frac{G_j}{H_j + \lambda}\right)^2 = \frac{G_j^2}{(H_j + \lambda)^2}$
2. Masukkan ke persamaan:
   $$= -\frac{G_j^2}{H_j + \lambda} + \frac{1}{2} (H_j + \lambda) \cdot \frac{G_j^2}{(H_j + \lambda)^2}$$
3. Coret satu $(H_j + \lambda)$ pada pembilang dan penyebut di suku kanan:
   $$= -\frac{G_j^2}{H_j + \lambda} + \frac{1}{2} \cdot \frac{G_j^2}{H_j + \lambda}$$
4. Samakan penyebut ($-1 + \frac{1}{2} = -\frac{1}{2}$):
   $$= \left( -1 + \frac{1}{2} \right) \frac{G_j^2}{H_j + \lambda} = -\frac{1}{2} \frac{G_j^2}{H_j + \lambda}$$

Dengan demikian, nilai skor kualitas optimal dari seluruh pohon adalah:

$$
\tilde{\mathcal{L}}^{(t)}(q) = -\frac{1}{2} \sum_{j=1}^T \frac{G_j^2}{H_j + \lambda} + \gamma T
$$

---

#### Menghitung `Split Gain` dari Pemotongan Cabang

> 🔍 **Apa yang disubstitusikan?**  
> Keuntungan pemotongan (`Split Gain`) adalah selisih antara skor error cabang Induk (Parent $P$) sebelum dipotong dikurangi total skor error 2 cabang anak (Left $L$ dan Right $R$) setelah dipotong:  
> $$\text{Split Gain} = \text{Score}_{\text{Parent}} - (\text{Score}_{\text{Left}} + \text{Score}_{\text{Right}})$$

Bentuk mentah substitusi:

$$
\text{Gain} = \left( -\frac{1}{2} \frac{G_P^2}{H_P + \lambda} + \gamma \right) - \left( -\frac{1}{2} \left[ \frac{G_L^2}{H_L + \lambda} + \frac{G_R^2}{H_R + \lambda} \right] + 2\gamma \right)
$$

Buka kurung dan kelompokkan $-\frac{1}{2}$ dengan $+2\gamma - \gamma = +\gamma$:

$$
\mathbf{\text{Split Gain} = \frac{1}{2} \left[ \frac{G_L^2}{H_L + \lambda} + \frac{G_R^2}{H_R + \lambda} - \frac{G_P^2}{H_P + \lambda} \right] - \gamma}
$$

* Jika $\text{Gain} > 0 \implies$ Pemotongan cabang menguntungkan (error turun lebih besar dari biaya daun $\gamma$).
* Jika $\text{Gain} \le 0 \implies$ Pemotongan cabang dibatalkan (`Pruning`).

---

#### Langkah 6: Pembaruan Prediksi Akhir (Pohon Baru Ditambahkan)
$$
\mathbf{\hat{y}_i^{(t)} = \hat{y}_i^{(t-1)} + \eta \cdot w_{q(\mathbf{x}_i)}^*}
$$

---

### 1.1 Implementasi Fitur di Aplikasi `chia.florist` & Perbandingan Ruang Himpunan Nilai

Di aplikasi `chia.florist`, mesin XGBoost tunggal ini diaplikasikan pada 2 use case bisnis yang berbeda. Perbedaan fundamental keduanya terletak pada **Ruang Himpunan Nilai Angka (Mathematical Output Space)** yang dihasilkan:

```
                            MESIN MATEMATIKA XGBOOST
                               (w* = -G / (H + λ))
                                        │
             ┌──────────────────────────┴──────────────────────────┐
             ▼                                                     ▼
   [XGBoost Regressor]                                   [XGBoost Classifier]
  (Demand Forecasting)                                  (Stockout Risk Alarm)
             │                                                     │
             ▼                                                     ▼
Ruang Bilangan Riil Kontinu                           Ruang Log-Odds / Logit
     y_hat ∈ [0, +∞)                                      z ∈ (-∞, +∞)
(e.g. 20.3 buket mawar terjual)                                    │
                                                                   ▼ [Pemetaan Sigmoid]
                                                              Ruang Probabilitas
                                                                 p ∈ [0.0, 1.0]
                                                            (e.g. 80.0% Risiko Habis)
                                                                   │
                                                                   ▼ [Threshold Keputusan]
                                                              Ruang Biner Diskrit
                                                                 y_hat ∈ {0, 1}
                                                            (1 = Restock, 0 = Aman)
```

#### A. Skenario 1: `XGBoost Regressor` (Demand Forecasting)
* **Use Case**: Memprediksi berapa buket "Mawar Merah Passion" (`PROD-001`) yang bakal terjual dalam 7 hari ke depan untuk persiapan stok Hari Valentine.
* **Ruang Himpunan**:
  * Gradien & Hessian: $g_i = \hat{y}_i - y_i \in \mathbb{R}$, $h_i = 1.0 \in \{1.0\}$.
  * Output Daun $w_j^*$: Bernilai langsung di dalam **Himpunan Bilangan Riil Kontinu** ($\mathbb{R}$, satuan fisik buket).
  * Prediksi Final: $\hat{y} \in [0, +\infty) \subset \mathbb{R}$ (Hasil: **20.3 buket mawar**).

#### B. Skenario 2: `XGBoost Classifier` (Inventory Stockout Risk)
* **Use Case**: Memprediksi apakah stok 12 buket Lili Putih akan habis sebelum supplier kebun tiba dalam 5 hari pengiriman.
* **Ruang Himpunan (3 Tingkat Pemetaan)**:
  1. **Ruang Log-Odds Internal ($\mathbb{R}$)**: Nilai daun $w_j^* = -\frac{\sum (p_i - y_i)}{\sum p_i(1-p_i) + \lambda}$ beroperasi di ruang logit $z \in (-\infty, +\infty)$.
  2. **Ruang Probabilitas Kontinu ($[0, 1]$)**: Skor total $z$ dipetakan melalui fungsi `Sigmoid` $p = \frac{1}{1 + e^{-z}} \in [0.0, 1.0]$ (Hasil: **$80.0\%$ Peluang Habis**).
  3. **Ruang Biner Diskrit ($\{0, 1\}$)**: Untuk keperluan operasional gudang, probabilitas kontinu diputuskan menjadi label diskrit tegas:
     $$\hat{y} = \begin{cases} 1 \quad (\text{STATUS: BUTUH RESTOCK DARURAT}) & \text{jika } p \ge 0.50 \\ 0 \quad (\text{STATUS: STOK AMAN}) & \text{jika } p < 0.50 \end{cases}$$

---

#### 📊 Tabel Matriks Perbandingan Ruang Himpunan Nilai

| Dimensi Perbandingan | `XGBoost Regressor` (Demand) | `XGBoost Classifier` (Stockout) |
|---|---|---|
| **Fungsi Loss Dasar** | Mean Squared Error ($\frac{1}{2}(y - \hat{y})^2$) | Binary Cross-Entropy (Log-Loss) |
| **Ruang Gradien ($g_i$)** | $\mathbb{R}$ (Selisih unit riil: $\hat{y}_i - y_i$) | $[-1, +1]$ (Selisih probabilitas: $p_i - y_i$) |
| **Ruang Hessian ($h_i$)** | $\{1.0\}$ (Konstanta skalar tetap) | $[0, 0.25]$ (Variansi kurva binomial $p(1-p)$) |
| **Ruang Nilai Daun ($w^*$)** | $\mathbb{R}$ (Unit besaran fisik target) | $\mathbb{R}$ (Unit delta log-odds / logit) |
| **Fungsi Aktivasi Penghubung** | Identitas ($f(z) = z$) | **`Sigmoid`** ($\sigma(z) = \frac{1}{1 + e^{-z}}$) |
| **Ruang Output Final** | **Kontinu** $[0, +\infty) \subset \mathbb{R}$ | **Probabilitas** $[0, 1] \to$ **Diskrit Biner** $\{0, 1\}$ |


---

## 2. `Isolation Forest` (Fitur: Operational Anomaly Detection)

### 📍 Letak File di Codebase:
* **Trainer**: [`src/anomaly_trainer.py`](file:///d:/__Projects/kage/chia.florist/intelligence-layer/src/anomaly_trainer.py) (`OperationalAnomalyTrainer`)
* **Fitur**: [`src/feature_engineering.py`](file:///d:/__Projects/kage/chia.florist/intelligence-layer/src/feature_engineering.py) (`OperationalAnomalyFeatureBuilder`)
* **Service**: [`app/services/predictor.py`](file:///d:/__Projects/kage/chia.florist/intelligence-layer/app/services/predictor.py) (`detect_anomaly`)

---

### 🎨 Filosofi Visual: *"Tebasan Pedang Acak di Ruang Kosong"*
Bayangkan kamu punya ruangan berisi kerumunan 1000 orang di tengah ruangan, dan ada 1 orang aneh yang berdiri sendirian di pojok ruangan yang sangat jauh:
* Kalau kamu menebas ruangan secara acak dengan sekat dinding, orang yang berdiri sendirian di pojok akan langsung terkurung sendirian hanya dalam **1 atau 2 sekat** (`Short Path Length`).
* Orang yang berada di tengah kerumunan padat butuh **puluhan sekat** sampai dia bisa terpisah sendirian (`Deep Path Length`).

---

### 🔄 Siklus Rumus Besar yang Mengayomi (`Isolation & Anomaly Scoring Loop`)

#### Langkah 1: Kedalaman Teoretis Rata-Rata Binary Search Tree ($c(n)$)
Untuk mengetahui apakah kedalaman suatu pohon itu termasuk "dangkal" atau "dalam", kita butuh angka pembanding standar teoretis dari struktur *Binary Search Tree* (BST) dengan $n$ data:

$$
c(n) = 2 \left( \ln(n - 1) + \gamma_{\text{Euler}} \right) - \frac{2(n - 1)}{n}
$$

*(Di mana $\gamma_{\text{Euler}} \approx 0.5772156649$ adalah konstanta Euler-Mascheroni).*

#### Langkah 2: Menghitung Rata-Rata Kedalaman Lintasan dari $T$ Pohon Acak
Data transaksi $\mathbf{x}$ dilewatkan ke seluruh $T$ pohon independen (biasanya $T = 100$), lalu dihitung rata-rata kedalamannya ($\mathbb{E}(h(\mathbf{x}))$):

$$
\mathbb{E}(h(\mathbf{x})) = \frac{1}{T} \sum_{t=1}^T h_t(\mathbf{x})
$$

#### Langkah 3: Master Exponential Decay Scoring Function

> 🔍 **Apa yang disubstitusikan?**  
> Rasio kedalaman relatif $\frac{\mathbb{E}(h(\mathbf{x}))}{c(n)}$ disubstitusikan ke dalam eksponen basis 2:

$$
\mathbf{s(\mathbf{x}, n) = 2^{-\frac{\mathbb{E}(h(\mathbf{x}))}{c(n)}}}
$$

* Jika $\mathbb{E}(h(\mathbf{x})) \to 0$ (sangat cepat terisolasi): $s \to 2^0 = \mathbf{1.0}$ (**Pasti Anomali**).
* Jika $\mathbb{E}(h(\mathbf{x})) \to c(n)$ (kedalaman rata-rata normal): $s \to 2^{-1} = \mathbf{0.5}$ (**Kondisi Normal / Wajar**).
* Jika $\mathbb{E}(h(\mathbf{x})) \to n-1$ (sangat susah terisolasi): $s \to 2^{-\infty} = \mathbf{0.0}$ (**Sangat Padat / Normal Murni**).

#### Langkah 4: Keputusan Ambang Batas Alert Operasional
$$
\text{Status}(\mathbf{x}) = \begin{cases} \text{ANOMALY (Picu Investigasi)} & \text{jika } s(\mathbf{x}, n) \ge 0.60 \\ \text{NORMAL (Transaksi Sah)} & \text{jika } s(\mathbf{x}, n) < 0.60 \end{cases}
$$

---

## 3. `Gradient Boosting Regressor` (Fitur: Courier SLA Estimator)

### 📍 Letak File di Codebase:
* **Trainer**: [`src/courier_sla_trainer.py`](file:///d:/__Projects/kage/chia.florist/intelligence-layer/src/courier_sla_trainer.py) (`CourierSLATrainer`)
* **Fitur**: [`src/feature_engineering.py`](file:///d:/__Projects/kage/chia.florist/intelligence-layer/src/feature_engineering.py) (`CourierSLAFeatureBuilder`)
* **Service**: [`app/services/predictor.py`](file:///d:/__Projects/kage/chia.florist/intelligence-layer/app/services/predictor.py) (`predict_courier_sla`)

---

### 🎨 Filosofi Visual: *"Lari Estafet Guru dan Barisan Murid"*
* **Pelari #0 (Guru)** membuat tebakan rata-rata durasi kota: *"Semua pengiriman butuh 24 jam"*.
* **Murid #1** datang dan melihat sisa selisihnya: *"Untuk SiCepat motor, tebakan guru kelebihan 6 jam"*. Murid #1 mengoreksi $-6$ jam.
* **Murid #2** datang melihat sisa kesalahan murid #1: *"Kalau jam sibuk hari Jumat sore, ada macet tambahan $+2$ jam"*. Murid #2 menambahkan $+2$ jam.
* Hasil akhir adalah penjumlahan kerja sama seluruh rantai murid.

---

### 🔄 Siklus Rumus Besar yang Mengayomi (`Functional Gradient Descent Loop`)

#### Langkah 1: Inisialisasi Tebakan Dasar Konstan ($F_0$)
Model mencari angka awal konstan yang meminimalkan total loss pada seluruh populasi data (untuk *Squared Error*, ini adalah rata-rata target $\bar{y}$):

$$
F_0(\mathbf{x}) = \arg\min_c \sum_{i=1}^N \mathcal{L}(y_i, c) = \bar{y}
$$

---

#### Langkah 2: Menghitung Sisa Error / Pseudo-Residuals ($r_{im}$)

> 🔍 **Apa yang disubstitusikan?**  
> Untuk fungsi loss kuadrat $\mathcal{L}(y_i, F) = \frac{1}{2}(y_i - F)^2$, kita turunkan terhadap $F$:  
> $$\frac{\partial \mathcal{L}}{\partial F} = \frac{1}{2} \cdot 2(y_i - F) \cdot (-1) = -(y_i - F)$$  
> Substitusikan hasil ini ke rumus pseudo-residual $r_{im} = -\left[ \frac{\partial \mathcal{L}}{\partial F} \right]$:

$$
r_{im} = -\left( -(y_i - F_{m-1}(\mathbf{x}_i)) \right) = \mathbf{y_i - F_{m-1}(\mathbf{x}_i)}
$$

*Minus bertemu minus menjadi plus, sehingga residu persis menjadi selisih tebakan dengan kenyataan riil!*

---

#### Langkah 3: Melatih Pohon Keputusan Baru ($h_m$) pada Data Residu
Pohon keputusan baru dilatih bukan untuk menebak target asli $y$, melainkan khusus untuk mencocokkan nilai sisa error $r_{im}$:

$$
\min_{h_m} \sum_{i=1}^N \left( r_{im} - h_m(\mathbf{x}_i) \right)^2
$$

#### Langkah 4: Pembaruan Model Aditif Berantai
Gabungkan pohon baru ke dalam model utama dengan pengali *Learning Rate* $\eta \in (0, 1]$:

$$
F_m(\mathbf{x}) = F_{m-1}(\mathbf{x}) + \eta \cdot h_m(\mathbf{x})
$$

#### Langkah 5: Master Final Prediction Equation
Setelah seluruh $M$ pohon selesai dilatih, fungsi prediksi akhir pengiriman adalah:

$$
\mathbf{F_{\text{Final}}(\mathbf{x}) = F_0(\mathbf{x}) + \sum_{m=1}^M \eta \cdot h_m(\mathbf{x})}
$$

---

## 🧭 Rangkuman Perbandingan Seluruh Formula Induk

| Fitur di `chia.florist` | Model ML | Mental Model Visual | Formula Besar Pemersatu (`Master Loop`) |
|---|---|---|---|
| **Demand Forecast** | `XGBoost Regressor` | Balok lego bertingkat mengikuti kontur lembah error | $\mathcal{L}^{(t)} \approx \sum [g_i f_t + \frac{1}{2}h_i f_t^2] + \Omega \implies w_j^* = -\frac{G_j}{H_j + \lambda} \implies \text{Gain}$ |
| **Stockout Risk** | `XGBoost Classifier` | Batas kritis dipotong lalu diperas jadi probabilitas | $p = \sigma(z) \implies g_i = p_i - y_i, \ h_i = p_i(1-p_i) \implies w_j^* = -\frac{\sum g_i}{\sum h_i + \lambda}$ |
| **Operational Anomaly** | `Isolation Forest` | Tebasan pedang acak mengukur keterpencilan ruang | $c(n) = 2\ln(n-1) + 2\gamma - \frac{2(n-1)}{n} \implies s(\mathbf{x}, n) = 2^{-\frac{\mathbb{E}(h(\mathbf{x}))}{c(n)}}$ |
| **Courier SLA** | `Gradient Boosting` | Lari estafet berantai mengoreksi sisa selisih guru | $F_0 = \bar{y} \implies r_{im} = y_i - F_{m-1} \implies F_{\text{Final}}(\mathbf{x}) = F_0 + \sum_{m=1}^M \eta h_m(\mathbf{x})$ |
