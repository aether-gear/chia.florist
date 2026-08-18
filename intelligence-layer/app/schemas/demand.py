from typing import Optional
from pydantic import BaseModel, Field

class DemandForecastRequest(BaseModel):
    product_id: str = Field(..., description="Unique product SKU identifier")
    gross_margin_pct: float = Field(0.35, ge=0.0, le=1.0, description="Gross profit margin percentage (0.0 to 1.0)")
    view_count: int = Field(500, ge=0, description="Product page view count in last 30 days")
    day_of_week: int = Field(1, ge=0, le=6, description="Day of week (0=Mon, 6=Sun)")
    month: int = Field(8, ge=1, le=12, description="Month (1-12)")
    is_weekend: float = Field(0.0, ge=0.0, le=1.0, description="Weekend flag (1.0 for Sat/Sun, 0.0 otherwise)")
    units_sold_lag_1: float = Field(5.0, ge=0.0, description="Units sold 1 day ago")
    units_sold_lag_7: float = Field(6.0, ge=0.0, description="Units sold 7 days ago")
    units_sold_lag_14: float = Field(4.0, ge=0.0, description="Units sold 14 days ago")
    units_sold_lag_30: float = Field(5.0, ge=0.0, description="Units sold 30 days ago")
    rolling_7d_mean: float = Field(5.2, ge=0.0, description="7-day rolling average daily sales")
    rolling_7d_std: float = Field(1.1, ge=0.0, description="7-day rolling standard deviation")
    rolling_30d_mean: float = Field(4.8, ge=0.0, description="30-day rolling average daily sales")
    rolling_30d_std: float = Field(1.4, ge=0.0, description="30-day rolling standard deviation")

class DemandForecastResponse(BaseModel):
    product_id: str
    predicted_units_sold_7d: float
    confidence_tier: str = "high"
