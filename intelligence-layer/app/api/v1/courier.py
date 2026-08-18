from fastapi import APIRouter, HTTPException, status
from app.schemas import ResponseWrapper, CourierSLARequest, CourierSLAResponse
from app.services.predictor import prediction_service

router = APIRouter(prefix="/predict", tags=["Courier SLA"])

@router.post("/courier-sla", response_model=ResponseWrapper[CourierSLAResponse])
def predict_courier_sla(payload: CourierSLARequest):
    """
    Estimate expected delivery duration (hours) and SLA reliability score for a given courier route.
    """
    try:
        res = prediction_service.predict_courier_sla(payload)
        return ResponseWrapper(success=True, data=res)
    except RuntimeError as re:
        raise HTTPException(status_code=status.HTTP_503_SERVICE_UNAVAILABLE, detail=str(re))
    except Exception as e:
        raise HTTPException(status_code=status.HTTP_500_INTERNAL_SERVER_ERROR, detail=f"Prediction error: {str(e)}")
