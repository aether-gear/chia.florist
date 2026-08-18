from typing import Optional
from pydantic import BaseModel, Field

class CourierSLARequest(BaseModel):
    courier_code: str = Field("jne", description="Courier partner code (jne, jnt, sicepat)")
    shipping_cost: float = Field(25000.0, ge=0.0, description="Shipping fee amount in IDR")
    dispatch_day_of_week: int = Field(1, ge=0, le=6, description="Day of week package is dispatched (0=Mon, 6=Sun)")
    dispatch_hour: int = Field(14, ge=0, le=23, description="Hour of dispatch (0-23)")
    dispatch_is_weekend: float = Field(0.0, ge=0.0, le=1.0, description="1.0 if weekend dispatch, 0.0 otherwise")

class CourierSLAResponse(BaseModel):
    courier_code: str
    estimated_duration_hours: float
    sla_confidence_score: float
    delivery_status: str
