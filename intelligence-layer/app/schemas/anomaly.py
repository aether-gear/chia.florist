from typing import List
from pydantic import BaseModel, Field

class AnomalyCheckRequest(BaseModel):
    amount: float = Field(350000.0, ge=0.0, description="Total order/payment transaction amount in IDR")
    time_to_pay_sec: float = Field(120.0, ge=0.0, description="Elapsed time to complete payment in seconds")
    is_manual_transfer: float = Field(0.0, ge=0.0, le=1.0, description="1.0 if manual bank transfer, 0.0 for instant gateway")
    is_failed_status: float = Field(0.0, ge=0.0, le=1.0, description="1.0 if transaction previously failed, 0.0 otherwise")

class AnomalyCheckResponse(BaseModel):
    is_anomaly: bool
    anomaly_score: float
    severity: str
    reasons: List[str]
