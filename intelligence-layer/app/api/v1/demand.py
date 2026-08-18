from fastapi import APIRouter, HTTPException, status
from app.schemas import ResponseWrapper, DemandForecastRequest, DemandForecastResponse
from app.services.predictor import prediction_service

router = APIRouter(prefix="/predict", tags=["Demand Forecasting"])

@router.post("/demand", response_model=ResponseWrapper[DemandForecastResponse])
def predict_demand(payload: DemandForecastRequest):
    """
    Predict 7-day SKU sales & unit demand volume based on time-series lags and product metrics.
    """
    try:
        res = prediction_service.predict_demand(payload)
        return ResponseWrapper(success=True, data=res)
    except RuntimeError as re:
        raise HTTPException(status_code=status.HTTP_503_SERVICE_UNAVAILABLE, detail=str(re))
    except Exception as e:
        raise HTTPException(status_code=status.HTTP_500_INTERNAL_SERVER_ERROR, detail=f"Prediction error: {str(e)}")
