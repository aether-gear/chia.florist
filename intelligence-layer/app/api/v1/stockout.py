from fastapi import APIRouter, HTTPException, status
from app.schemas import ResponseWrapper, StockoutRiskRequest, StockoutRiskResponse
from app.services.predictor import prediction_service

router = APIRouter(prefix="/predict", tags=["Stockout Risk"])

@router.post("/stockout-risk", response_model=ResponseWrapper[StockoutRiskResponse])
def predict_stockout_risk(payload: StockoutRiskRequest):
    """
    Classify stockout probability within supplier lead time and provide safety reorder alerts.
    """
    try:
        res = prediction_service.predict_stockout_risk(payload)
        return ResponseWrapper(success=True, data=res)
    except RuntimeError as re:
        raise HTTPException(status_code=status.HTTP_503_SERVICE_UNAVAILABLE, detail=str(re))
    except Exception as e:
        raise HTTPException(status_code=status.HTTP_500_INTERNAL_SERVER_ERROR, detail=f"Prediction error: {str(e)}")
