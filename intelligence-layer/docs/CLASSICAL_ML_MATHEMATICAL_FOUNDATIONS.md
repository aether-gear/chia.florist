# 📐 Classical Machine Learning: The Plain-English & Mathematical Guide

> **How to read this guide**: Every concept is explained in **casual, everyday language first**, followed by the official name in parentheses `(Technical Term: ...)`, along with the exact math formula and a **complete list of every variable, constant, vector, and matrix** defined below it.


## 🗺️ The Evolution Tree (Why New Models Were Born)

```
1. Simple Straight-Line Guessing (Ordinary Least Squares / OLS)
   │ ❌ Breaks when columns repeat the same info or when columns outnumber rows.
   ▼
2. Straight-Line with Guardrails to Shrink or Delete Useless Columns (Ridge L2 & Lasso L1)
   │ ❌ Breaks because it can only draw flat, straight boundaries. Cannot answer Yes/No probabilities.
   ▼
3. Squeezing Line Outputs into 0% to 100% Chances (Logistic Regression & GLMs)
   │ ❌ Breaks because it cannot learn combinations (e.g. LeadTime × Velocity) without manual math tricks.
   ▼
4. Bending Space to Separate Mixed Points with Wide Safety Buffers (Support Vector Machines / Kernel SVM)
   │ ❌ Breaks because it takes massive compute O(N³) and crawls to a halt on datasets over 50,000 rows.
   ▼
5. Playing a Game of "20 Questions" with Step-by-Step Rules (Single Decision Tree / CART)
   │ ❌ Breaks because one small change in data flips the whole tree (wildly unstable, overfits easily).
   ├──────────────────────────────────────────────────────┐
   ▼ (Fixing instability by averaging 100 trees)          ▼ (Fixing errors by learning from leftover mistakes)
6. A Council of 100 Random Trees Voting (Random Forest) 7. A Chain of Apprentices Correcting Errors (AdaBoost & GBM)
   │ ❌ Cannot extrapolate future trends beyond seen max.  │ ❌ Slow, only looks at error direction, not curve.
   │                                                      ▼
   └──────────────────────────────────────────────────► 8. The Modern Heavyweight: Smart Residual Learning
                                                           (XGBoost / LightGBM / CatBoost)
                                                           Uses 2nd derivatives (curvature) + built-in penalties
                                                           to calculate exact tree weights in milliseconds.
```


## 1. The Core Big-Picture Ideas

### 1.1 Teaching with Answers vs. Finding Hidden Shapes
* **Teaching with an answer key** `(Technical Term: Supervised Learning)`: You give the model input features $\mathbf{x}$ (e.g., flower price, shop location) and the correct target label $y$ (e.g., actual units sold). The model learns a function $\hat{y} = f(\mathbf{x})$ to predict future answers.
* **Finding patterns on your own without labels** `(Technical Term: Unsupervised Learning)`: You only give the model input data $\mathbf{x}$ (e.g., customer payment habits) with no answer key. The model discovers natural groupings `(Clustering)` or spots weird outliers `(Anomaly Detection)`.

### 1.4 The Universal Tug-of-War: Rigidness vs. Jumpiness
The goal of all ML optimization is minimizing prediction errors plus complexity:

$$
\min_{f \in \mathcal{F}} \left[ \frac{1}{N} \sum_{i=1}^N \mathcal{L}(y_i, f(\mathbf{x}_i)) + \Omega(f) \right]
$$

#### 📋 Variable & Constant Breakdown:
| Symbol | Type | What it Represents |
|---|---|---|
| $N$ | Integer Constant | Total number of training samples (rows of data). |
| $i$ | Integer Index | Current sample row index ($i = 1, 2, \dots, N$). |
| $\mathbf{x}_i$ | Vector ($\mathbb{R}^d$) | The input features of sample $i$ (e.g., `[price, lead_time, view_count]`). |
| $y_i$ | Scalar ($\mathbb{R}$ or $\{0,1\}$) | The true ground-truth answer for sample $i$ (actual sales or risk flag). |
| $f$ | Function | The model function we want to learn ($\hat{y}_i = f(\mathbf{x}_i)$). |
| $\mathcal{F}$ | Set | The hypothesis space (the universe of all allowed formulas/shapes). |
| $\mathcal{L}(y_i, f(\mathbf{x}_i))$ | Function | The Loss Function: measures how wrong the prediction is on sample $i$. |
| $\Omega(f)$ | Scalar Function | The Regularization Penalty: fines the model for becoming too complicated. |


## 2. The Straight-Line Family (Linear Models & Logistic Regression)

### 2.1 Basic Line Fitting `(Technical Term: Ordinary Least Squares / OLS Regression)`

$$
\hat{y} = f(\mathbf{x}) = \mathbf{w}^T \mathbf{x} + b = \tilde{\mathbf{w}}^T \tilde{\mathbf{x}}
$$

#### 📋 Variable & Constant Breakdown:
| Symbol | Type | What it Represents |
|---|---|---|
| $\hat{y}$ | Scalar Float ($\mathbb{R}$) | The model's predicted output value. |
| $\mathbf{x}$ | Column Vector ($d \times 1$) | The $d$ input feature values for a single data point. |
| $\mathbf{w}$ | Column Vector ($d \times 1$) | The learned weight coefficients (importance multiplier of each feature). |
| $\mathbf{w}^T$ | Row Vector ($1 \times d$) | Transpose of $\mathbf{w}$ (allows dot product $\mathbf{w}^T\mathbf{x} = \sum w_j x_j$). |
| $b$ | Scalar Float ($\mathbb{R}$) | Bias / Intercept term (the predicted value when all inputs $\mathbf{x} = 0$). |
| $\tilde{\mathbf{w}}$ | Vector ($(d+1) \times 1$) | Augmented weight vector $[b, w_1, w_2, \dots, w_d]^T$ combining bias and weights. |
| $\tilde{\mathbf{x}}$ | Vector ($(d+1) \times 1$) | Augmented input vector $[1, x_1, x_2, \dots, x_d]^T$ with a dummy $1$ for bias. |

#### The OLS Loss Function (Mean Squared Error):

$$
\mathcal{L}(\mathbf{w}) = \frac{1}{2N} \|\mathbf{X}\mathbf{w} - \mathbf{y}\|_2^2
$$

#### 📋 Variable & Constant Breakdown:
| Symbol | Type | What it Represents |
|---|---|---|
| $\mathcal{L}(\mathbf{w})$ | Scalar Float | Total Mean Squared Error across the entire dataset. |
| $N$ | Integer | Total number of training rows. |
| $\mathbf{X}$ | Matrix ($N \times d$) | The full design matrix containing all feature values for all $N$ samples. |
| $\mathbf{y}$ | Column Vector ($N \times 1$) | The true target values column for all $N$ samples. |
| $\|\cdot\|_2^2$ | Norm Operator | Squared $L_2$ Euclidean Norm: sums the squared error of all elements ($\sum e_i^2$). |
| $\frac{1}{2}$ | Mathematical Constant | Scaling factor chosen so the derivative $\frac{d}{dx}\frac{1}{2}x^2 = x$ cancels out cleanly. |

#### The One-Step Exact Solution `(Technical Term: The Normal Equation)`:

$$
\mathbf{w}^* = (\mathbf{X}^T \mathbf{X})^{-1} \mathbf{X}^T \mathbf{y}
$$

#### 📋 Variable & Constant Breakdown:
| Symbol | Type | What it Represents |
|---|---|---|
| $\mathbf{w}^*$ | Vector ($d \times 1$) | The mathematically optimal weights that achieve the lowest possible MSE. |
| $\mathbf{X}^T$ | Matrix ($d \times N$) | Transpose of feature matrix $\mathbf{X}$ (swapping rows and columns). |
| $\mathbf{X}^T \mathbf{X}$ | Square Matrix ($d \times d$) | Feature covariance matrix measuring how features correlate with each other. |
| $(\cdot)^{-1}$ | Matrix Operator | Matrix Inverse (the matrix algebra equivalent of dividing by a number). |


### 2.2 Adding Guardrails to Keep Weights Sane `(Technical Term: Regularization)`

#### Ridge Regression ($L_2$):

$$
\mathcal{L}_{\text{Ridge}}(\mathbf{w}) = \frac{1}{2N} \|\mathbf{X}\mathbf{w} - \mathbf{y}\|_2^2 + \frac{\lambda}{2} \|\mathbf{w}\|_2^2
$$

$$
\mathbf{w}^*_{\text{Ridge}} = (\mathbf{X}^T \mathbf{X} + N\lambda \mathbf{I})^{-1} \mathbf{X}^T \mathbf{y}
$$

#### 📋 Variable & Constant Breakdown:
| Symbol | Type | What it Represents |
|---|---|---|
| $\lambda$ | Hyperparameter ($\ge 0$) | Regularization strength constant: higher values shrink weights closer to zero. |
| $\|\mathbf{w}\|_2^2$ | Scalar Float | Sum of all squared weights $\sum_{j=1}^d w_j^2$ ($L_2$ Penalty). |
| $\mathbf{I}$ | Identity Matrix ($d \times d$) | Square matrix with $1$s on the diagonal and $0$s elsewhere. |
| $N\lambda\mathbf{I}$ | Matrix ($d \times d$) | Boosts the diagonal of $\mathbf{X}^T\mathbf{X}$ by $+N\lambda$, guaranteeing it is **always invertible**. |

#### Lasso Regression ($L_1$):

$$
\mathcal{L}_{\text{Lasso}}(\mathbf{w}) = \frac{1}{2N} \|\mathbf{X}\mathbf{w} - \mathbf{y}\|_2^2 + \lambda \sum_{j=1}^d |w_j|
$$

$$
w_j^* = S\left( \frac{1}{N}\mathbf{x}_j^T (\mathbf{y} - \mathbf{X}_{-j}\mathbf{w}_{-j}), \lambda \right)
$$

$$
S(z, \gamma) = \text{sign}(z) \max(0, |z| - \gamma)
$$

#### 📋 Variable & Constant Breakdown:
| Symbol | Type | What it Represents |
|---|---|---|
| $|w_j|$ | Scalar Float | Absolute value of weight $j$ ($L_1$ Penalty). |
| $w_j^*$ | Scalar Float | Optimal weight for feature $j$ updated via Coordinate Descent. |
| $\mathbf{x}_j$ | Column Vector ($N \times 1$) | The values of feature column $j$ across all $N$ rows. |
| $\mathbf{X}_{-j}\mathbf{w}_{-j}$ | Vector ($N \times 1$) | Prediction made by all other features excluding feature $j$. |
| $\mathbf{y} - \mathbf{X}_{-j}\mathbf{w}_{-j}$ | Vector ($N \times 1$) | The leftover residual that feature $j$ is responsible for explaining. |
| $S(z, \gamma)$ | Operator | Soft-Thresholding function: clamps values below threshold $\gamma$ to exact $0.0$. |
| $\text{sign}(z)$ | Function | Returns $+1$ if $z > 0$, $-1$ if $z < 0$, and $0$ if $z = 0$. |


### 2.3 Squeezing Line Numbers into 0% to 100% Probabilities `(Technical Term: Logistic Regression)`

$$
p = \sigma(z) = \frac{1}{1 + e^{-z}} \quad \text{where } z = \mathbf{w}^T \mathbf{x} + b
$$

#### 📋 Variable & Constant Breakdown:
| Symbol | Type | What it Represents |
|---|---|---|
| $p$ | Scalar Probability ($[0, 1]$) | Predicted chance of class 1 (e.g. $p = 0.85$ means 85% chance of stockout). |
| $\sigma(\cdot)$ | Function | Logistic Sigmoid activation function that squashes real numbers into $(0, 1)$. |
| $z$ | Scalar Float ($(-\infty, +\infty)$) | Raw linear score (logit) computed from input features. |
| $e$ | Mathematical Constant | Euler's number ($e \approx 2.71828\dots$). |

#### Binary Cross-Entropy Loss:

$$
\mathcal{L}(\mathbf{w}) = -\frac{1}{N} \sum_{i=1}^N \left[ y_i \ln \hat{y}_i + (1 - y_i) \ln(1 - \hat{y}_i) \right]
$$

#### 📋 Variable & Constant Breakdown:
| Symbol | Type | What it Represents |
|---|---|---|
| $y_i$ | Binary Integer ($\{0, 1\}$) | True label of row $i$ ($1$ = event occurred, $0$ = did not occur). |
| $\hat{y}_i$ | Scalar Float ($[0, 1]$) | Model's predicted probability for row $i$. |
| $\ln(\hat{y}_i)$ | Scalar Float | Logarithm of predicted probability (penalizes confident wrong guesses severely). |


## 3. Bending Space with Wide Buffers (Support Vector Machines)

$$
\min_{\mathbf{w}, b, \boldsymbol{\xi}} \frac{1}{2} \|\mathbf{w}\|_2^2 + C \sum_{i=1}^N \xi_i \quad \text{subject to} \quad y_i(\mathbf{w}^T \mathbf{x}_i + b) \ge 1 - \xi_i, \quad \xi_i \ge 0
$$

$$
K(\mathbf{x}, \mathbf{z}) = \exp\left(-\gamma \|\mathbf{x} - \mathbf{z}\|_2^2\right)
$$

#### 📋 Variable & Constant Breakdown:
| Symbol | Type | What it Represents |
|---|---|---|
| $\mathbf{w}$ | Vector ($d \times 1$) | Weight vector defining the orientation of the separating hyperplane. |
| $\|\mathbf{w}\|_2$ | Scalar Float | Euclidean length of vector $\mathbf{w}$ (Buffer width = $\frac{2}{\|\mathbf{w}\|_2}$). |
| $C$ | Hyperparameter ($\ge 0$) | Tradeoff penalty constant: balances wide buffer vs. number of misclassified dots. |
| $\xi_i$ | Scalar Float ($\ge 0$) | Slack variable: distance by which sample $i$ violates the safe margin. |
| $K(\mathbf{x}, \mathbf{z})$ | Function | Kernel function: calculates high-dimensional similarity between two data points. |
| $\gamma$ | Hyperparameter ($\ge 0$) | RBF kernel spread constant: higher $\gamma$ creates tighter, more complex boundaries. |


## 4. Decision Trees & Ensembles

### 4.1 Decision Trees (CART)

$$
\text{Gini}(R_m) = 1 - \sum_{k=1}^K p_{mk}^2
$$

#### 📋 Variable & Constant Breakdown:
| Symbol | Type | What it Represents |
|---|---|---|
| $R_m$ | Subset of Space | The $m$-th bounding box region defined by tree question splits. |
| $K$ | Integer | Total number of target classes. |
| $p_{mk}$ | Probability ($[0, 1]$) | Fraction of samples in region $R_m$ that belong to class $k$. |
| $\text{Gini}(R_m)$ | Scalar Float | Impurity score: $0.0$ = 100% pure clean box, $0.5$ = completely mixed 50/50. |

### 4.2 Random Forest

$$
\text{Total Forest Variance} = \rho \sigma^2 + \frac{1 - \rho}{B} \sigma^2
$$

#### 📋 Variable & Constant Breakdown:
| Symbol | Type | What it Represents |
|---|---|---|
| $B$ | Integer Constant | Total number of decision trees trained in parallel. |
| $\sigma^2$ | Scalar Float | Variance (jumpiness) of an individual single decision tree. |
| $\rho$ | Correlation ($[0, 1]$) | Average similarity between trees (reduced by random feature subspacing). |

### 4.3 Gradient Boosting (GBM Foundation)

$$
F_{\text{Final}}(\mathbf{x}) = F_0 + \sum_{m=1}^M \eta \cdot f_m(\mathbf{x})
$$

$$
r_{im} = -\left[ \frac{\partial L(y_i, F(\mathbf{x}_i))}{\partial F(\mathbf{x}_i)} \right]_{F = F_{m-1}}
$$

#### 📋 Variable & Constant Breakdown:
| Symbol | Type | What it Represents |
|---|---|---|
| $F_{\text{Final}}(\mathbf{x})$ | Scalar Float | Final aggregated prediction after $M$ sequential boosting steps. |
| $F_0$ | Scalar Float | Base starting guess (usually the dataset target mean $\bar{y}$). |
| $M$ | Integer Constant | Total number of boosting rounds/trees in the chain. |
| $m$ | Integer Index | Current boosting stage index ($m = 1, 2, \dots, M$). |
| $\eta$ | Hyperparameter ($[0.01, 0.1]$) | Learning rate / Shrinkage: dials down each tree's contribution to prevent overfitting. |
| $f_m(\mathbf{x})$ | Function | Tree $m$, trained to predict the leftover errors of stage $m-1$. |
| $r_{im}$ | Scalar Float | The pseudo-residual (negative gradient of loss) for row $i$ at iteration $m$. |

---

### 4.4 Formula Inti Terpadu XGBoost (`The Unified XGBoost Mathematical Engine`)

Semua varian XGBoost (baik Regresi maupun Klasifikasi) digerakkan oleh **satu mesin matematika yang sama** berbasis Deret Taylor Orde ke-2 dan pengelompokan daun:

#### 1. Aproksimasi Taylor Orde ke-2 & Pengelompokan Daun ($j = 1 \dots T$):

$$
\mathcal{L}^{(t)} \approx \sum_{j=1}^T \left[ G_j w_j + \frac{1}{2} (H_j + \lambda) w_j^2 \right] + \gamma T
$$

Di mana $G_j = \sum_{i \in I_j} g_i$ adalah total gradien daun $j$, dan $H_j = \sum_{i \in I_j} h_i$ adalah total hessian daun $j$.

#### 2. Nilai Bobot Daun Optimal ($w_j^*$):

$$
w_j^* = -\frac{G_j}{H_j + \lambda}
$$

#### 3. Skor Kualitas Pohon (Hasil Substitusi $w_j^*$ ke $\mathcal{L}^{(t)}$):

$$
\tilde{\mathcal{L}}^{(t)} = -\frac{1}{2} \sum_{j=1}^T \frac{G_j^2}{H_j + \lambda} + \gamma T
$$

#### 4. Kriteria Pemotongan Cabang (`Split Gain Criterion`):

$$
\text{Split Gain} = \frac{1}{2} \left[ \frac{G_L^2}{H_L + \lambda} + \frac{G_R^2}{H_R + \lambda} - \frac{G_P^2}{H_P + \lambda} \right] - \gamma
$$

```
         Parent Node [Gp, Hp]
               /      \
              /        \
    Left Node           Right Node
    [GL, HL]            [GR, HR]
```

#### 📋 Variable & Constant Breakdown:
| Symbol | Type | What it Represents |
|---|---|---|
| $G_L, H_L$ | Scalar Floats | Sum of gradients and Hessians for samples in the proposed **left child** ($I_L$). |
| $G_R, H_R$ | Scalar Floats | Sum of gradients and Hessians for samples in the proposed **right child** ($I_R$). |
| $G_P, H_P$ | Scalar Floats | Sum of gradients and Hessians in the **parent node** ($I_P$) before splitting. |
| $w_j^*$ | Scalar Float | Optimal prediction score output by leaf node $j$. |
| $\lambda$ | Hyperparameter ($\ge 0$) | $L_2$ leaf regularization penalty preventing extreme leaf weights. |
| $\gamma$ | Hyperparameter ($\ge 0$) | Minimum loss reduction threshold (pruning cost) required to add an extra leaf. |

---

### 4.5 XGBoost di Aplikasi `chia.florist`: Skenario & Perbandingan Ruang Himpunan Nilai

Di aplikasi `chia.florist`, mesin XGBoost yang sama ini diterapkan pada 2 use case berbeda. Perbedaan fundamental keduanya terletak pada **Ruang Himpunan Nilai Angka (Mathematical Output Space)** yang dihasilkan:

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

#### Skenario 1: `XGBoost Regressor` (Demand Forecasting)
* **Letak File**: [`src/forecasting_trainer.py`](file:///d:/__Projects/kage/chia.florist/intelligence-layer/src/forecasting_trainer.py) $\to$ [`app/api/v1/forecast.py`](file:///d:/__Projects/kage/chia.florist/intelligence-layer/app/api/v1/forecast.py)
* **Skenario Bisnis**: Memprediksi berapa buket "Mawar Merah Passion" (`PROD-001`) yang bakal terjual dalam 7 hari ke depan untuk persiapan stok Hari Valentine.
* **Ruang Himpunan Input & Output**:
  * Gradien & Hessian: $g_i = \hat{y}_i - y_i$, $h_i = 1.0$.
  * Output Daun $w_j^*$: Bernilai langsung di dalam **Himpunan Bilangan Riil Kontinu** ($\mathbb{R}$, satuan fisik buket).
  * Prediksi Final: $\hat{y} = \hat{y}^{(0)} + \sum \eta w_{q(\mathbf{x})}^* \in [0, +\infty)$.
* **Hasil Nyata**: $\hat{y} = 20.3$ buket mawar.

#### Skenario 2: `XGBoost Classifier` (Inventory Stockout Risk)
* **Letak File**: [`src/stockout_trainer.py`](file:///d:/__Projects/kage/chia.florist/intelligence-layer/src/stockout_trainer.py) $\to$ [`app/api/v1/inventory.py`](file:///d:/__Projects/kage/chia.florist/intelligence-layer/app/api/v1/inventory.py)
* **Skenario Bisnis**: Memprediksi apakah stok 12 buket Lili Putih akan habis sebelum supplier kebun tiba dalam 5 hari pengiriman.
* **Ruang Himpunan Input & Output (3 Tahapan Ruang)**:
  1. **Ruang Log-Odds Internal ($\mathbb{R}$)**: Nilai daun $w_j^* = -\frac{\sum (p_i - y_i)}{\sum p_i(1-p_i) + \lambda}$ beroperasi di ruang logit $z \in (-\infty, +\infty)$.
  2. **Ruang Probabilitas Kontinu ($[0, 1]$)**: Skor total $z$ dipetakan melalui fungsi `Sigmoid` $p = \frac{1}{1 + e^{-z}} \in [0.0, 1.0]$ (e.g. $p = 0.8000$).
  3. **Ruang Biner Diskrit ($\{0, 1\}$)**: Untuk keperluan operasional gudang, probabilitas kontinu diputuskan menjadi label diskrit tegas:
     $$\hat{y} = \begin{cases} 1 \quad (\text{STATUS: BUTUH RESTOCK DARURAT}) & \text{jika } p \ge 0.50 \\ 0 \quad (\text{STATUS: STOK AMAN}) & \text{jika } p < 0.50 \end{cases}$$

---

#### 📊 Tabel Matriks Perbandingan Ruang Himpunan Nilai (Output Space Comparison)

| Dimensi Perbandingan | `XGBoost Regressor` (Demand) | `XGBoost Classifier` (Stockout) |
|---|---|---|
| **Fungsi Loss Dasar** | Mean Squared Error ($\frac{1}{2}(y - \hat{y})^2$) | Binary Cross-Entropy (Log-Loss) |
| **Ruang Gradien ($g_i$)** | $\mathbb{R}$ (Selisih unit riil: $\hat{y}_i - y_i$) | $[-1, +1]$ (Selisih probabilitas: $p_i - y_i$) |
| **Ruang Hessian ($h_i$)** | $\{1.0\}$ (Konstanta skalar tetap) | $[0, 0.25]$ (Variansi kurva binomial $p(1-p)$) |
| **Ruang Nilai Daun ($w^*$)** | $\mathbb{R}$ (Unit besaran fisik target) | $\mathbb{R}$ (Unit delta log-odds / logit) |
| **Fungsi Aktivasi Penghubung** | Identitas ($f(z) = z$) | **`Sigmoid`** ($\sigma(z) = \frac{1}{1 + e^{-z}}$) |
| **Ruang Output Final** | **Kontinu** $[0, +\infty) \subset \mathbb{R}$ | **Probabilitas** $[0, 1] \to$ **Diskrit Biner** $\{0, 1\}$ |



## 5. Unsupervised Learning & Anomaly Detection

### 5.1 $k$-Means Clustering

$$
J = \sum_{k=1}^K \sum_{\mathbf{x}_i \in C_k} \|\mathbf{x}_i - \boldsymbol{\mu}_k\|_2^2
$$

#### 📋 Variable & Constant Breakdown:
| Symbol | Type | What it Represents |
|---|---|---|
| $J$ | Scalar Float | Within-Cluster Sum of Squares (Inertia: measures how compact clusters are). |
| $K$ | Integer Constant | The number of clusters chosen beforehand (e.g. $K=3$ customer tiers). |
| $k$ | Integer Index | Current cluster index ($k = 1, 2, \dots, K$). |
| $C_k$ | Set of Samples | The group of sample points currently assigned to cluster $k$. |
| $\mathbf{x}_i$ | Vector ($d \times 1$) | Coordinate vector of sample data point $i$. |
| $\boldsymbol{\mu}_k$ | Vector ($d \times 1$) | Centroid vector: the center pin (mean coordinate) of cluster $k$. |


### 5.2 Principal Component Analysis (PCA)

$$
\mathbf{\Sigma} \mathbf{w} = \lambda \mathbf{w}
$$

#### 📋 Variable & Constant Breakdown:
| Symbol | Type | What it Represents |
|---|---|---|
| $\mathbf{\Sigma}$ | Square Matrix ($d \times d$) | Sample covariance matrix $\mathbf{\Sigma} = \frac{1}{N}\mathbf{X}^T\mathbf{X}$ (measures how all columns correlate). |
| $\mathbf{w}$ | Unit Vector ($d \times 1$) | Eigenvector: an orthogonal direction of variance (Principal Component axis). |
| $\lambda$ | Scalar Float ($\ge 0$) | Eigenvalue: amount of data variance captured along direction $\mathbf{w}$. |


### 5.3 Isolation Forest Anomaly Detection

```
       Normal Point (Deep inside tree: Path length h(x) is long)
                 o
             o  o  o  o  o
            o  o (x) o  o
             o  o  o  o

       Anomaly (Alone in open space: Path length h(x) is super short!)
                                           o
                                          / \
                                        (x)  o
```

$$
s(\mathbf{x}, n) = 2^{-\frac{\mathbb{E}(h(\mathbf{x}))}{c(n)}}
$$

$$
c(n) = 2 \left( \ln(n - 1) + 0.5772156649 \right) - \frac{2(n - 1)}{n}
$$

#### 📋 Variable & Constant Breakdown:
| Symbol | Type | What it Represents |
|---|---|---|
| $s(\mathbf{x}, n)$ | Scalar Float ($[0, 1]$) | Normalized Anomaly Score ($\to 1.0$ = Definite Anomaly, $< 0.5$ = Normal). |
| $\mathbf{x}$ | Vector ($d \times 1$) | The single transaction sample being scored. |
| $n$ | Integer Constant | Subsampling size used to construct each isolation tree (typically $n = 256$). |
| $h(\mathbf{x})$ | Integer | Path length: number of random splits required to isolate point $\mathbf{x}$. |
| $\mathbb{E}(h(\mathbf{x}))$ | Scalar Float | Average path length of point $\mathbf{x}$ across all trees in the forest. |
| $c(n)$ | Scalar Float | Benchmark average path length of an unsuccessful search in a Binary Search Tree of size $n$. |
| $0.5772156649$ | Mathematical Constant | Euler-Mascheroni constant ($\gamma_{\text{Euler}}$). |


## 6. Time-Series Forecasting: Additive Decomposition (Prophet)

```
1. Long-term Growth Trend g(t)    2. Repeating Cycles s(t)          3. Holiday Spikes h(t)
   (Straight line with bends)        (Weekly/Yearly waves)             (Valentine's / Mother's Day)
   ▲                                 ▲                                 ▲
   │    /                            │   /\    /\    /\                │       |     |
   │   /                             │  /  \  /  \  /  \               │       |     |
   └──/────────►                     └─/────\/────\/────►              └───────┴─────┴──►
```

$$
y(t) = g(t) + s(t) + h(t) + \epsilon_t
$$

$$
s(t) = \sum_{n=1}^N \left[ a_n \cos\left( \frac{2\pi n t}{P} \right) + b_n \sin\left( \frac{2\pi n t}{P} \right) \right]
$$

#### 📋 Variable & Constant Breakdown:
| Symbol | Type | What it Represents |
|---|---|---|
| $y(t)$ | Scalar Float | Total forecasted sales volume at future timestamp $t$. |
| $g(t)$ | Function | Piecewise linear growth trend (captures overall company growth or decline). |
| $s(t)$ | Function | Periodic seasonality waves (weekly and yearly consumer patterns). |
| $h(t)$ | Function | Holiday impact function (scheduled shock multiplier on Valentine's / Mother's Day). |
| $\epsilon_t$ | Scalar Float | Gaussian noise / unpredictable day-to-day fluctuation at time $t$. |
| $P$ | Scalar Constant | Period length of the cycle ($P=7$ for 7-day week cycle, $P=365.25$ for yearly cycle). |
| $N$ | Integer | Number of Fourier harmonic terms (higher $N$ models sharper, more complex waves). |
| $a_n, b_n$ | Parameters ($\mathbb{R}$) | Learned amplitude coefficients determining the height and phase of seasonal waves. |
| $\pi$ | Mathematical Constant | Archimedes' constant ($\pi \approx 3.14159\dots$). |


## 7. How to Read Evaluation Numbers (The Math of Evaluation)

### 7.1 Coefficient of Determination ($R^2$ Score)

$$
R^2 = 1 - \frac{\text{SS}_{\text{res}}}{\text{SS}_{\text{tot}}} = 1 - \frac{\sum_{i=1}^N (y_i - \hat{y}_i)^2}{\sum_{i=1}^N (y_i - \bar{y})^2}
$$

#### 📋 Variable & Constant Breakdown:
| Symbol | Type | What it Represents |
|---|---|---|
| $R^2$ | Scalar Float ($(-\infty, 1.0]$) | Coefficient of determination: proportion of variance explained by model. |
| $y_i$ | Scalar Float | True ground-truth value of sample $i$. |
| $\hat{y}_i$ | Scalar Float | Model's predicted value for sample $i$. |
| $\bar{y}$ | Scalar Float | Simple historical average of all true $y$ values ($\bar{y} = \frac{1}{N}\sum y_i$). |
| $\text{SS}_{\text{res}}$ | Scalar Float | Residual Sum of Squares: total squared error made by your model. |
| $\text{SS}_{\text{tot}}$ | Scalar Float | Total Sum of Squares: total squared error if you just blindly guessed the average $\bar{y}$. |

* $R^2 = 1.0$ (100%): Flawless prediction ($\text{SS}_{\text{res}} = 0$).
* $R^2 = 0.0$ (0%): The model is doing the exact same thing as guessing the average $\bar{y}$.
* **Why your Demand Forecasting had $R^2 = 0.0537$**: Synthetic data was pure uniform random noise; the model found no real signal, so $\text{SS}_{\text{res}} \approx \text{SS}_{\text{tot}}$, yielding $R^2 \approx 5\%$.


### 7.2 Classification Alarm Metrics

$$
\text{Precision} = \frac{TP}{TP + FP}, \quad \text{Recall} = \frac{TP}{TP + FN}, \quad F_1 = \frac{2TP}{2TP + FP + FN}
$$

#### 📋 Variable & Constant Breakdown:
| Symbol | Type | What it Represents |
|---|---|---|
| $TP$ | Integer Count | **True Positives**: Real emergencies that the alarm correctly triggered on. |
| $FP$ | Integer Count | **False Positives (Type I Error)**: False alarms (alarm went off, but nothing happened). |
| $FN$ | Integer Count | **False Negatives (Type II Error)**: Missed emergencies (crisis happened, alarm slept). |
| $TN$ | Integer Count | **True Negatives**: Normal quiet days that the alarm correctly ignored. |
| $\text{Precision}$ | Ratio ($[0, 1]$) | When the alarm rings, probability that it is a **real fire** ($\frac{TP}{TP+FP}$). |
| $\text{Recall}$ | Ratio ($[0, 1]$) | Out of all real fires, percentage that the alarm **successfully caught** ($\frac{TP}{TP+FN}$). |
| $F_1$ | Harmonic Mean ($[0, 1]$) | Balanced score between not crying wolf (Precision) and not sleeping through fires (Recall). |
| $\text{ROC-AUC}$ | Area ($[0.5, 1.0]$) | Ranking ability: probability that a random true anomaly gets a higher score than a normal row ($0.5$ = Coin flip, $1.0$ = Perfect ranking). |

