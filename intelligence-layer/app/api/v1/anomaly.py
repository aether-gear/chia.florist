from fastapi import APIRouter, HTTPException, status
from app.schemas import ResponseWrapper, AnomalyCheckRequest, AnomalyCheckResponse
from app.services.predictor import prediction_service

router = APIRouter(prefix="/predict", tags=["Anomaly Detection"])

@router.post("/anomaly", response_model=ResponseWrapper[AnomalyCheckResponse])
def detect_anomaly(payload: AnomalyCheckRequest):
    """
    Perform operational anomaly detection on payment latency, fulfillment delay, and stock discrepancy features.
    """
    try:
        res = prediction_service.detect_anomaly(payload)
        return ResponseWrapper(success=True, data=res)
    except RuntimeError as re:
        raise HTTPException(status_code=status.HTTP_503_SERVICE_UNAVAILABLE, detail=str(re))
    except Exception as e:
        raise HTTPException(status_code=status.HTTP_500_INTERNAL_SERVER_ERROR, detail=f"Prediction error: {str(e)}")
