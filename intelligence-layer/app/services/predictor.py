import logging
import pandas as pd
import numpy as np
from typing import Dict, Any

from app.services.model_loader import model_registry
from app.schemas import (
    DemandForecastRequest, DemandForecastResponse,
    StockoutRiskRequest, StockoutRiskResponse,
    CourierSLARequest, CourierSLAResponse,
    AnomalyCheckRequest, AnomalyCheckResponse,
)

logger = logging.getLogger("predictor")

class PredictionService:
    """
    Inference service that converts validated Pydantic DTOs into DataFrame feature matrices
    matching trained model feature schemas.
    """
    
    @staticmethod
    def predict_demand(payload: DemandForecastRequest) -> DemandForecastResponse:
        model = model_registry.get_model("demand")
        
        data_dict = {
            "gross_margin_pct": [payload.gross_margin_pct],
            "view_count": [payload.view_count],
            "day_of_week": [payload.day_of_week],
            "month": [payload.month],
            "is_weekend": [payload.is_weekend],
            "units_sold_lag_1": [payload.units_sold_lag_1],
            "units_sold_lag_7": [payload.units_sold_lag_7],
            "units_sold_lag_14": [payload.units_sold_lag_14],
            "units_sold_lag_30": [payload.units_sold_lag_30],
            "rolling_7d_mean": [payload.rolling_7d_mean],
            "rolling_7d_std": [payload.rolling_7d_std],
            "rolling_30d_mean": [payload.rolling_30d_mean],
            "rolling_30d_std": [payload.rolling_30d_std],
        }
        df_features = pd.DataFrame(data_dict)
        
        preds = model.predict(df_features)
        predicted_units = float(np.maximum(0.0, preds[0]))
        
        return DemandForecastResponse(
            product_id=payload.product_id,
            predicted_units_sold_7d=round(predicted_units, 2),
            confidence_tier="high" if predicted_units > 0 else "low"
        )

    @staticmethod
    def predict_stockout_risk(payload: StockoutRiskRequest) -> StockoutRiskResponse:
        model = model_registry.get_model("stockout")
        
        burn_rate = payload.stock_burn_rate_7d
        est_days = payload.estimated_days_to_stockout
        if est_days is None:
            est_days = payload.stock / max(0.1, burn_rate)

        urgency_ratio = payload.reorder_urgency_ratio
        if urgency_ratio is None:
            lead_demand = burn_rate * payload.supplier_lead_time_days
            urgency_ratio = lead_demand / max(1.0, payload.stock)

        # Exact feature columns expected by stockout_risk.json
        data_dict = {
            "stock": [payload.stock],
            "reserved_stock": [payload.reserved_stock],
            "stock_burn_rate_7d": [burn_rate],
            "supplier_lead_time_days": [payload.supplier_lead_time_days],
            "estimated_days_to_stockout": [est_days],
            "reorder_urgency_ratio": [urgency_ratio],
        }
        df_features = pd.DataFrame(data_dict)
        
        if hasattr(model, "predict_proba"):
            probs = model.predict_proba(df_features)[0]
            prob_stockout = float(probs[1]) if len(probs) > 1 else float(probs[0])
        else:
            raw_pred = model.predict(df_features)[0]
            prob_stockout = float(raw_pred)

        will_stockout = bool(prob_stockout >= 0.5)
        
        if prob_stockout >= 0.75:
            risk_level = "CRITICAL"
        elif prob_stockout >= 0.4:
            risk_level = "WARNING"
        else:
            risk_level = "NORMAL"

        return StockoutRiskResponse(
            product_id=payload.product_id,
            stockout_probability=round(prob_stockout, 4),
            will_stockout=will_stockout,
            risk_level=risk_level
        )

    @staticmethod
    def predict_courier_sla(payload: CourierSLARequest) -> CourierSLAResponse:
        model = model_registry.get_model("courier")
        
        code = payload.courier_code.lower()
        courier_jne = 1.0 if "jne" in code else 0.0
        courier_jnt = 1.0 if "jnt" in code or "j&t" in code else 0.0
        courier_sicepat = 1.0 if "sicepat" in code else 0.0

        # Exact feature columns expected by courier_sla.json
        data_dict = {
            "courier_jne": [courier_jne],
            "courier_jnt": [courier_jnt],
            "courier_sicepat": [courier_sicepat],
            "shipping_cost": [payload.shipping_cost],
            "dispatch_day_of_week": [float(payload.dispatch_day_of_week)],
            "dispatch_hour": [float(payload.dispatch_hour)],
            "dispatch_is_weekend": [float(payload.dispatch_is_weekend)],
        }
        df_features = pd.DataFrame(data_dict)
        
        preds = model.predict(df_features)
        est_hours = float(np.maximum(1.0, preds[0]))

        confidence = 0.95 if (courier_jne or courier_jnt or courier_sicepat) else 0.80

        return CourierSLAResponse(
            courier_code=payload.courier_code,
            estimated_duration_hours=round(est_hours, 1),
            sla_confidence_score=confidence,
            delivery_status="ON_TRACK" if est_hours <= 48.0 else "DELAY_RISK"
        )

    @staticmethod
    def detect_anomaly(payload: AnomalyCheckRequest) -> AnomalyCheckResponse:
        model = model_registry.get_model("anomaly")
        
        # Exact feature columns and order expected by isolation_forest.pkl
        data_dict = {
            "amount": [payload.amount],
            "time_to_pay_sec": [payload.time_to_pay_sec],
            "is_failed_status": [payload.is_failed_status],
            "is_manual_transfer": [payload.is_manual_transfer],
        }
        df_features = pd.DataFrame(data_dict)

        raw_pred = model.predict(df_features)[0]
        score = float(model.decision_function(df_features)[0]) if hasattr(model, "decision_function") else 0.0
        
        is_anomaly = bool(raw_pred == -1)

        reasons = []
        if is_anomaly:
            if payload.time_to_pay_sec > 1800:
                reasons.append(f"Excessive payment completion delay ({payload.time_to_pay_sec:.0f}s)")
            if payload.is_failed_status == 1.0:
                reasons.append("Prior transaction failure status flag set")
            if payload.amount > 2000000.0:
                reasons.append(f"High transaction value amount anomaly (IDR {payload.amount:,.0f})")
            if not reasons:
                reasons.append("Out-of-distribution operational transaction pattern")

        severity = "HIGH" if (is_anomaly and score < -0.1) else ("MEDIUM" if is_anomaly else "NORMAL")

        return AnomalyCheckResponse(
            is_anomaly=is_anomaly,
            anomaly_score=round(score, 4),
            severity=severity,
            reasons=reasons
        )


prediction_service = PredictionService()
