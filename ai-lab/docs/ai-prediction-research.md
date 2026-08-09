# AI Prediction for E-Commerce Trend & Summary
## Research & Approach Selection — chia-florist B2C

> **Context**: Evaluating AI options for 2 microservices:
> - **Service 1 (General)**: Train on public e-commerce data → understand general seasonal/market trends
> - **Service 2 (Specific)**: Fine-tune/adapt on chia-florist B2C data → predict flower shop-specific events
>
> **Constraints**: ⏱ 3 Days | 💰 Free or ≤ Rp 700.000 (~$43 USD) / month for production
>
> *Researched: 2026-07-12 | Sources: Web research + subagent analysis*

---

## 📊 Quick Verdict

| Approach | 3-Day Feasibility | Cost (Production) | Recommendation |
|---|---|---|---|
| **ML (Prophet + XGBoost)** | ✅ Very Doable | ✅ Free (HF Spaces) | ⭐ Best for forecasting numbers |
| **Deep Learning (LSTM/Transformer)** | ⚠️ Tight | ✅ Free (HF Spaces) | ❌ Overkill for this scale |
| **LLM API Only** | ✅ Very Easy | ⚠️ Depends on volume | ⚠️ Good for summaries, not forecasting |
| **Hybrid (ML + LLM Summary)** | ✅ Doable | ✅ Free possible | ⭐⭐ **RECOMMENDED** |

---

## Option A — Traditional ML (Prophet + XGBoost)

### What it is
- **Prophet** (Meta/Facebook): Decompose time-series into trend + seasonality + holidays. Plug-and-play with minimal code.
- **XGBoost**: Tree-based model; learns from engineered features (day-of-week, month, lag sales, is_holiday, promo flags).

### Suitability for Flower Shop
Flower shops have **strong seasonal spikes**: Valentine's Day, Mother's Day, Eid, Christmas, graduation season. Both Prophet and XGBoost excel at this.

### 3-Day Implementation Plan

```
Day 1 — Data & EDA
  ├── Download public e-commerce dataset (Online Retail UCI)
  ├── Clean, resample to daily/weekly aggregates
  ├── Feature engineering: day_of_week, month, is_holiday, lag_7d, lag_30d
  └── EDA in Jupyter Notebook

Day 2 — Train Service 1 (General Model)
  ├── Prophet baseline on UCI Online Retail dataset
  ├── XGBoost on feature-engineered version
  └── Evaluate with MAE/RMSE, plot actual vs predicted

Day 3 — Train Service 2 + API
  ├── Adapt model using chia-florist transaction data
  ├── Wrap inference as FastAPI endpoint
  └── Deploy to Hugging Face Spaces (free)
```

### Cost Estimate (Production)
| Resource | Platform | Cost |
|---|---|---|
| ML Inference API | Hugging Face Spaces (CPU Basic) | **Free** — 2 vCPU, 16 GB RAM |
| Model storage | HF Spaces Git repo | **Free** |
| Alternative | Render Hobby | **Free** (cold start 30-60s) |
| Alternative | Modal.com Starter | **Free** ($30/mo credits reset monthly) |

### Libraries
```
prophet==1.1.x
xgboost==2.x
scikit-learn
pandas
fastapi + uvicorn
```

### Pros / Cons
| ✅ Pros | ❌ Cons |
|---|---|
| Fast to implement (2-3 days) | Needs feature engineering |
| Very interpretable (trend/seasonality breakdown) | Less accurate on very small datasets |
| Runs on CPU (no GPU cost) | Requires periodic retraining |
| Industry standard for retail forecasting | No natural language output (needs LLM for summaries) |

---

## Option B — Deep Learning (LSTM / Transformer)

### What it is
- **LSTM**: Recurrent neural network; learns sequential patterns in time-series
- **Temporal Fusion Transformer (TFT)** / **TiDE**: Transformer-based time-series models; state-of-the-art accuracy

### 3-Day Feasibility: ⚠️ TIGHT
- Setting up LSTM properly (hypertuning, handling vanishing gradients) = Day 1-2 alone
- TFT / TiDE with PyTorch Forecasting: steep learning curve
- **Verdict**: Possible only if you already have hands-on DL experience. Not recommended for 3-day sprint.

### Cost Estimate
| Resource | Cost |
|---|---|
| HF Spaces CPU | **Free** (slow inference ~2-5 sec) |
| HF Spaces T4 GPU | ~$0.60/hour on demand |
| Google Colab (training) | **Free** (T4 GPU, 12h session limit) |

### Pros / Cons
| ✅ Pros | ❌ Cons |
|---|---|
| Higher accuracy on long sequences | Complex to implement in 3 days |
| Can capture complex non-linear patterns | Needs larger dataset to generalize |
| Portable model artifacts | Slow inference on CPU |
| | Risk of overfitting on small chia-florist data |

---

## Option C — LLM API Only

### What it is
Use an LLM (GPT-4o-mini, Gemini Flash, Groq Llama) to:
1. Accept structured sales data as context (CSV / JSON summary)
2. Generate natural-language trend analysis and predictions
3. Summarize upcoming events/patterns

### Pricing Comparison (July 2026)

| Provider | Model | Input ($/1M tokens) | Output ($/1M tokens) | Free Tier |
|---|---|---|---|---|
| **Google AI Studio** | Gemini 3.1 Flash-Lite | $0.25 (paid) | $1.50 (paid) | ✅ Yes (rate-limited, non-production) |
| **Google AI Studio** | Gemini 3.5 Flash | $1.50 (paid) | $9.00 (paid) | ✅ Yes (rate-limited) |
| **OpenAI** | GPT-4o-mini | $0.15 | $0.60 | ❌ No free tier |
| **Groq** | Llama 4 Scout 17B | Pay-as-you-go | Pay-as-you-go | ✅ 1,000 req/day (not prod-ready) |
| **Anthropic** | Claude Haiku 3.5 | $0.80 | $4.00 | ❌ No free tier |

### Budget Simulation (≤ Rp 700.000 / ~$43 USD / month)
Assume: 50 analyses/day × 30 days = 1,500 requests/month
Each request: ~2,000 input tokens + ~500 output tokens

| Model | Monthly Cost |
|---|---|
| GPT-4o-mini | (1500 × 2000 × $0.15/1M) + (1500 × 500 × $0.60/1M) = **~$0.90/mo** ✅ |
| Gemini 3.1 Flash-Lite (paid) | (1500 × 2000 × $0.25/1M) + (1500 × 500 × $1.50/1M) = **~$1.88/mo** ✅ |
| Gemini Free Tier | **$0** ✅ but rate-limited, non-production ToS |

> LLM API costs are negligible at low volume. At 50 req/day, even paid tiers are well under $5/month.

### Limitations of LLM-Only Approach
- ❌ LLMs do NOT reliably predict numbers — they hallucinate sales figures
- ❌ Cannot truly forecast future demand like a time-series model
- ❌ Stateless — every call is independent; no persistent model learning from your data
- ✅ Excellent for **summarizing** trend outputs from a real forecasting model
- ✅ Very easy to implement (few hours)

---

## Option D — Hybrid Architecture ⭐ RECOMMENDED

### Architecture Overview

```
┌─────────────────────────────────────────────────────────┐
│                   Service 1: General Trend               │
│  Public Dataset (UCI Online Retail, Flower Kaggle)       │
│  → Prophet/XGBoost → General Seasonality Model          │
│  → Exposed as: GET /predict/general?horizon=30           │
└─────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────┐
│                 Service 2: chia-florist Specific         │
│  chia-florist B2C order data                            │
│  → Prophet/XGBoost fine-tuned → Shop-Specific Model     │
│  → Exposed as: POST /predict/shop { data: [...] }        │
└─────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────┐
│                    LLM Summary Layer                     │
│  Numeric forecast output → structured prompt             │
│  → LLM API (GPT-4o-mini / Gemini Flash-Lite)            │
│  → Natural language insight: "Demand will spike 40%      │
│     next week due to Mother's Day. Prepare 150 bouquets" │
└─────────────────────────────────────────────────────────┘
```

### Why Hybrid Wins
1. **Reliable numbers**: ML models give deterministic, verifiable forecasts
2. **Human-readable output**: LLM wraps the numbers in actionable insight
3. **Cost efficient**: ML inference is free (HF Spaces), LLM cost is minimal at low volume
4. **Maintainable**: Two clear layers — forecasting & summarization
5. **Implementable in 3 days**: Each layer is independently simple

### 3-Day Implementation Timeline

```
Day 1: Data & General Model (Service 1)
  AM  → Download UCI Online Retail II dataset from Kaggle
  AM  → EDA, cleaning, feature engineering in notebook
  PM  → Train Prophet model → evaluate → serialize to .pkl
  PM  → Quick FastAPI wrapper for /predict/general endpoint

Day 2: chia-florist Model + LLM Summary (Service 2 + Layer 3)
  AM  → Adapt feature engineering for chia-florist order data
  AM  → Train XGBoost on chia-florist data (Prophet if small dataset)
  PM  → Integrate LLM API (Gemini free or GPT-4o-mini)
  PM  → Design prompt template: "Given this forecast: {json}, summarize..."
  PM  → Test end-to-end: data → forecast → natural language summary

Day 3: API + Deployment
  AM  → Merge both services into clean FastAPI app
  AM  → Add /predict/general and /predict/shop endpoints
  PM  → Deploy to Hugging Face Spaces (Gradio or Docker)
  PM  → Write API documentation + smoke test with Postman/curl
```

---

## 🏗️ Hosting & Infrastructure

### Recommended Free Stack

| Layer | Tool | Platform | Monthly Cost |
|---|---|---|---|
| ML Inference API | FastAPI + Prophet/XGBoost | Hugging Face Spaces (CPU Basic) | **Free** |
| LLM Summary | Gemini API (Google AI Studio) | Google Cloud | **Free** (dev) / ~$1-2 (prod, low vol) |
| Model Storage | .pkl / .joblib files | HF Spaces repo | **Free** |
| Vector DB (RAG) | ChromaDB (local) | Embedded in app | **Free** |

### Alternative Hosting Options

| Platform | Cost | Notes |
|---|---|---|
| **Hugging Face Spaces** | Free (2 vCPU, 16 GB) | Best for ML. Supports Gradio, Streamlit, Docker |
| **Render Hobby** | Free | Cold starts (30-60s delay after 15min idle) |
| **Modal.com Starter** | Free ($30/mo reset) | Serverless, scale-to-zero |
| **Google Cloud Run** | Free tier (2M req/mo) | Requires Docker image |
| **Streamlit Community Cloud** | Free | Good for dashboards, not pure API |

---

## 📦 Free Training Datasets

| Dataset | Source | Size | Best For |
|---|---|---|---|
| **Online Retail II (UCI)** | Kaggle / UCI | 1M+ rows | General e-commerce seasonality |
| **Retail Sales w/ Seasonality** | Kaggle (search: retail seasonal sales) | Synthetic | Holiday + promo effect modeling |
| **Flower Shop Sales** | Kaggle (search: flower shop sales) | Small | Direct florist analog |
| **Store Sales - Time Series** | Kaggle Competition | Large | Multi-store seasonal forecasting |
| **E-Commerce Behavior Data** | Kaggle (mkechinov) | 70M events | Purchase funnel behavior |

---

## ⚙️ Final Stack Decision

### ✅ Go With: **Hybrid (ML + LLM)**

```yaml
service_1_general:
  framework: Prophet (primary) + XGBoost (secondary)
  dataset: UCI Online Retail II + Kaggle Flower Shop dataset
  output: numeric forecast (daily/weekly demand)
  hosting: Hugging Face Spaces CPU (free)

service_2_chia_florist:
  framework: Prophet or XGBoost (depends on data volume)
  dataset: chia-florist B2C order history
  output: product-level demand forecast
  hosting: same HF Space or separate endpoint

llm_summary_layer:
  provider: Google Gemini API (Gemini free tier for dev)
  fallback: GPT-4o-mini (~$0.90/mo at 50 req/day)
  purpose: numeric forecast → natural language business insight

total_estimated_cost:
  free_tier: $0/month (HF Spaces + Gemini free)
  low_prod:  ~$2-5/month (Gemini paid or GPT-4o-mini)
  budget_check: ✅ Well within Rp 700.000 ($43 USD) limit
```

### ❌ Why NOT Deep Learning?
- 3 days is too short to properly tune LSTM/TFT
- Small chia-florist dataset → DL overfits easily
- Prophet + XGBoost achieve ~85-92% of DL accuracy on tabular time-series
- No GPU cost = stays within budget

### ❌ Why NOT LLM-Only?
- LLMs fabricate sales numbers — unreliable for forecasting
- Cannot maintain learned patterns across sessions without expensive RAG
- Stateless by nature; can't do real prediction

---

## 🌸 Flower Shop–Specific Insights (Indonesia Context)

Prophet natively supports Indonesian holidays via `add_country_holidays('ID')`.
Key seasonal peaks to encode:

| Event | Timing | Impact |
|---|---|---|
| Valentine's Day | Feb 14 | Single-day spike, red roses dominant |
| Mother's Day | 2nd Sunday May | Week-long, highest annual revenue |
| Eid al-Fitr (Idul Fitri) | Varies (lunar) | 🇮🇩 Indonesian gift-giving peak |
| Eid al-Adha | Varies (lunar) | Secondary peak |
| Graduation Season | May–June | Significant for Indonesian florists |
| Wedding Season | June–September | Dry season correlation |
| Christmas / New Year | Dec 25 – Jan 1 | Moderate peak |

### Recommended Prophet Config for Indonesia
```python
from prophet import Prophet

model = Prophet(
    yearly_seasonality=True,
    weekly_seasonality=True,
    daily_seasonality=False,
    changepoint_prior_scale=0.05  # lower = more stable trend line
)
# Add all Indonesian public holidays automatically
model.add_country_holidays(country_name='ID')
```

---

## ⚡ Bonus: TimeGPT (Zero-Shot Option)

[Nixtla TimeGPT](https://docs.nixtla.io/) is a pre-trained foundation model for time-series — no training required.
- **Free tier**: Available for low-volume use
- **Usage**: Pass your historical data → get forecast via API (no model training)
- **Verdict for chia-florist**: Interesting as a fallback or quick validation, but Prophet gives more control and interpretability for this use case.

```python
# TimeGPT zero-shot example
from nixtla import NixtlaClient
client = NixtlaClient(api_key="your_key")
forecast = client.forecast(df=your_df, h=30, freq="D")  # 30-day forecast, daily
```

---

## 📋 Decision Matrix

| Criteria | ML Only | Deep Learning | LLM Only | Hybrid (⭐) |
|---|---|---|---|---|
| 3-day feasibility | ✅ Easy | ❌ Hard | ✅ Easiest | ✅ Feasible |
| Numerical accuracy | ✅ High | ✅✅ Highest | ❌ Poor | ✅ High |
| Natural language output | ❌ None | ❌ None | ✅✅ Excellent | ✅✅ Excellent |
| Monthly cost | $0 | $0–$30 | $0–$1 | $0–$5 |
| Infrastructure complexity | Low | Med-High | Low | Low-Med |
| Handles flower seasonality | ✅ Excellent | ✅ Good | ⚠️ Via context | ✅ Excellent |
| Indonesian holiday support | ✅ Prophet built-in | Manual | Via prompt | ✅ Prophet + prompt |
| Works with small dataset | ✅ | ⚠️ | ✅ | ✅ |

---

## 🚀 Immediate Next Steps

1. **Setup environment**:
   ```bash
   pip install prophet xgboost fastapi uvicorn pandas scikit-learn joblib google-generativeai
   ```
2. **Download dataset**: UCI Online Retail II from Kaggle → save to `data/raw/`
3. **Create notebook**: `notebooks/01_eda_general_trend.ipynb`
4. **Get API key**: [Google AI Studio](https://aistudio.google.com) → free Gemini API key (Flash-Lite, no billing needed)
5. **Alternative free LLM**: [Groq Console](https://console.groq.com) → free API key (Llama 4, no CC needed)
6. **Extend this repo**: Use existing `src/` structure, add `src/forecaster.py` and `src/summarizer.py`

### Recommended `requirements.txt` additions
```
prophet>=1.1.5
xgboost>=2.0.0
scikit-learn>=1.3.0
fastapi>=0.110.0
uvicorn>=0.27.0
pandas>=2.0.0
google-generativeai>=0.5.0   # OR groq>=0.5.0 (free tier)
# chromadb>=0.4.0             # optional: RAG vector store
```

---

## 📚 References

- [Prophet Quick Start](https://facebook.github.io/prophet/docs/quick_start.html)
- [XGBoost Time Series (Kaggle Tutorial)](https://www.kaggle.com/code/robikscube/time-series-forecasting-with-machine-learning)
- [UCI Online Retail Dataset](https://archive.ics.uci.edu/ml/datasets/online+retail)
- [Kaggle Store Sales Competition](https://www.kaggle.com/competitions/store-sales-time-series-forecasting)
- [Hugging Face Spaces Overview](https://huggingface.co/docs/hub/spaces-overview)
- [Google AI Studio Free Tier](https://aistudio.google.com/app/apikey)
- [Modal.com Pricing](https://modal.com/pricing)

---

*Generated: 2026-07-12 | chia-florist ai-lab research*
