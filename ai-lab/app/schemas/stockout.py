from typing import Optional
from pydantic import BaseModel, Field

class StockoutRiskRequest(BaseModel):
    product_id: str = Field(..., description="Unique product SKU identifier")
    stock: float = Field(10.0, ge=0.0, description="Current available unreserved stock")
    reserved_stock: float = Field(2.0, ge=0.0, description="Stock reserved for active orders")
    stock_burn_rate_7d: float = Field(2.5, ge=0.0, description="Average daily stock depletion rate over past 7 days")
    supplier_lead_time_days: float = Field(7.0, ge=0.0, description="Supplier lead time in days")
    estimated_days_to_stockout: Optional[float] = Field(None, ge=0.0, description="Estimated days remaining until stock reaches zero")
    reorder_urgency_ratio: Optional[float] = Field(None, ge=0.0, description="Ratio of lead time demand to current stock")

class StockoutRiskResponse(BaseModel):
    product_id: str
    stockout_probability: float
    will_stockout: bool
    risk_level: str
