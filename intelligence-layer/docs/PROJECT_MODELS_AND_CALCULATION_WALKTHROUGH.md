# 🌸 chia.florist: Real-World Model Guide & Step-by-Step Calculations

> **Guide Purpose**: This guide breaks down the machine learning models used in the `chia.florist` project alongside simpler classical alternatives. For every model, we take a real flower shop business scenario and calculate the prediction step-by-step on **one single sample of data** from its base formula to error minimization.

## 1. Demand & Sales Forecasting (Predicting 7-Day Units Sold)

### Scenario: Valentine's Week Rose Demand
Valentine's Day is coming up next week. The shop manager needs to know how many bouquets of "Red Passion Roses" (`PROD-001`) will sell over the next 7 days so they can order enough fresh stems from the Bandung flower farm without over-purchasing perishable flowers.

### Model Used in Project: XGBoost Regressor (Tree Leaf Weight Calculation)

XGBoost makes a baseline guess, looks at the leftover error on our sample, and uses both the slope and the curvature to calculate the optimal correction step for the decision tree leaf.

$$
w^* = -\frac{g_i}{h_i + \lambda}
$$

Given:
* $g_i = -6.0$ (The 1st gradient: our previous guess was 6 units too low)
* $h_i = 1.0$ (The 2nd gradient / Hessian for squared error loss)
* $\lambda = 1.0$ ($L_2$ regularization penalty to prevent wild over-adjustments)
* $\hat{y}^{(0)} = 20.0$ (Our initial baseline guess in units)
* $\eta = 0.1$ (Learning rate / shrinkage)

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

Thus, the updated forecast output is **20.3 units sold**.

### Classical ML Alternative: Ridge Regression (Linear Weighted Sum)

Instead of building trees, a regularized linear model simply multiplies each florist signal (such as past sales and catalog page views) by a learned importance multiplier and adds a baseline number.

$$
\hat{y} = w_1 x_1 + w_2 x_2 + b
$$

Given:
* $x_1 = 25.0$ (Units sold in the previous 7 days)
* $x_2 = 120.0$ (Product catalog page views this week)
* $w_1 = 0.65$ (Importance weight for past sales)
* $w_2 = 0.05$ (Importance weight for page views)
* $b = 4.0$ (Base organic sales intercept)

$$
\hat{y} = (0.65 \times 25.0) + (0.05 \times 120.0) + 4.0
$$

$$
\hat{y} = 16.25 + 6.00 + 4.00
$$

$$
\hat{y} = 26.25
$$

Thus, the predicted demand output is **26.25 units sold**.

## 2. Inventory Stockout Risk (Will We Run Out Before Supplier Arrives?)

### Scenario: Running Out of Lilies Before Supplier Delivery
A boutique store currently has 12 White Lily bouquets on the shelf. The store sells an average of 3 bouquets per day, and the flower farm takes 5 days to deliver a fresh batch. We need a risk probability score between 0% and 100% to decide if we must trigger an emergency purchase order.

### Model Used in Project: XGBoost Classifier (Probability via Sigmoid Squashing)

The classification tree outputs a raw unconstrained score based on inventory ratios, which is then squashed into a calibrated 0.0 to 1.0 probability of running out of stock.

$$
p = \frac{1}{1 + e^{-z}}
$$

Given:
* $z = 1.3863$ (Raw margin score computed by tree leaves for our Lily SKU)
* $e = 2.71828$ (Euler's mathematical constant)

$$
p = \frac{1}{1 + (2.71828)^{-1.3863}}
$$

$$
p = \frac{1}{1 + 0.2500}
$$

$$
p = \frac{1}{1.2500} = 0.8000
$$

Thus, the stockout risk output is **80.00% (CRITICAL RISK)**.

### Classical ML Alternative: Logistic Regression Loss Evaluation

When training on this sample, if the store actually ran out of stock (true label = 1), we evaluate how well the model performed using Binary Cross-Entropy loss.

$$
\mathcal{L} = -\left[ y \ln(p) + (1 - y) \ln(1 - p) \right]
$$

Given:
* $y = 1.0$ (The store actually ran out of stock)
* $p = 0.80$ (Model predicted an 80% chance)
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

Thus, the training loss penalty on this sample is **0.2231**.

## 3. Operational Anomaly Detection (Spotting Delayed Payments & Webhook Failures)

### Scenario: High-Value Order with 24-Hour Payment Lag
A customer places an order for a luxury wedding flower arrangement worth Rp 3.500.000 using Manual Bank Transfer. Normal transfers are completed in 15 minutes (900 seconds), but this transfer took 86.400 seconds (24 hours). We want to check if this is an operational anomaly.

### Model Used in Project: Isolation Forest (Random Partition Anomaly Score)

Isolation Forest isolates weird points using random cuts. Because this 24-hour delay is far away from normal data, it gets isolated very quickly near the top of the tree.

$$
s(x, n) = 2^{-\frac{h(x)}{c(n)}}
$$

Given:
* $h(x) = 2.0$ (Number of random cuts needed to isolate this transaction)
* $n = 256$ (Total samples in the tree construction)
* $c(256) = 2(\ln(255) + 0.5772) - \frac{2(255)}{256} = 10.23$ (Average normal cut depth)

$$
s = 2^{-\frac{2.0}{10.23}}
$$

$$
s = 2^{-0.1955}
$$

$$
s = 0.8732
$$

Thus, the calculated anomaly score is **0.8732 (HIGH ANOMALY)**.

### Classical ML Alternative: Statistical Z-Score (Distance from Normal Crowd)

A fast statistical alternative measures how many standard deviations away this transaction's latency is from the historical shop average.

$$
z = \frac{x - \mu}{\sigma}
$$

Given:
* $x = 86400.0$ (Transaction completion latency in seconds)
* $\mu = 1800.0$ (Historical average transfer latency in seconds)
* $\sigma = 7200.0$ (Historical standard deviation in seconds)

$$
z = \frac{86400.0 - 1800.0}{7200.0}
$$

$$
z = \frac{84600.0}{7200.0}
$$

$$
z = 11.75
$$

Thus, the latency deviation output is **11.75 standard deviations away (EXTREME OUTLIER)**.

## 4. Courier SLA & Delivery Duration (Estimating Delivery Hours)

### Scenario: Friday Afternoon Delivery via SiCepat
A customer orders a fresh flower bouquet to be shipped across Jakarta via SiCepat on Friday at 3:00 PM (shipping cost Rp 15.000). Store staff needs to know the estimated delivery hours to guarantee the flowers arrive before wilting.

### Model Used in Project: Gradient Boosting Regressor (Sequential Step Update)

Gradient boosting starts with the citywide average delivery duration, and each tree adds an incremental adjustment based on the selected courier and Friday afternoon rush hour traffic.

$$
F_1(x) = F_0 + \eta \cdot f_1(x)
$$

Given:
* $F_0 = 24.0$ (Citywide baseline delivery duration in hours)
* $\eta = 0.1$ (Model shrinkage / learning rate)
* $f_1(x) = -60.0$ (Tree leaf adjustment score for SiCepat motorized express dispatch)

$$
F_1(x) = 24.0 + (0.1 \times -60.0)
$$

$$
F_1(x) = 24.0 - 6.0
$$

$$
F_1(x) = 18.0
$$

Thus, the estimated delivery duration output is **18.0 hours (ON TRACK)**.

### Classical ML Alternative: Decision Tree Regressor (Group Average Reduction)

A single decision tree sorts the order into a leaf matching `courier == SICEPAT` and `is_weekend == 0` and averages the historical delivery records in that bucket.

$$
\bar{y} = \frac{1}{N} \sum_{i=1}^N y_i
$$

Given:
* $N = 3$ (Three historical orders matching SiCepat weekday afternoon delivery)
* $y_1 = 16.0$ (Sample 1 delivery hours)
* $y_2 = 18.0$ (Sample 2 delivery hours)
* $y_3 = 20.0$ (Sample 3 delivery hours)

$$
\bar{y} = \frac{16.0 + 18.0 + 20.0}{3}
$$

$$
\bar{y} = \frac{54.0}{3}
$$

$$
\bar{y} = 18.0
$$

Thus, the partition average prediction output is **18.0 hours**.

## 5. Model Evaluation (How We Measure Error on One Sample)

### Scenario: Checking How Far Off Our Forecast Was
We forecasted that a store would sell **25 bouquets** of Sunflower arrangements today. At the end of the day, the register recorded **29 bouquets sold**. We want to calculate the Squared Error and Absolute Error for this day.

### Mean Absolute Error (MAE) on One Sample

$$
e_{\text{abs}} = |y - \hat{y}|
$$

Given:
* $y = 29.0$ (Actual registered sales)
* $\hat{y} = 25.0$ (Model forecast)

$$
e_{\text{abs}} = |29.0 - 25.0|
$$

$$
e_{\text{abs}} = |4.0|
$$

$$
e_{\text{abs}} = 4.0
$$

Thus, the absolute error output is **4.0 units**.

### Mean Squared Error (MSE) on One Sample

$$
e_{\text{sq}} = (y - \hat{y})^2
$$

Given:
* $y = 29.0$ (Actual registered sales)
* $\hat{y} = 25.0$ (Model forecast)

$$
e_{\text{sq}} = (29.0 - 25.0)^2
$$

$$
e_{\text{sq}} = (4.0)^2
$$

$$
e_{\text{sq}} = 16.0
$$

Thus, the squared error output is **16.0 units squared**.
