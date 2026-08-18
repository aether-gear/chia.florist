from fastapi import APIRouter
from app.config import settings
from app.schemas import ResponseWrapper, HealthResponse
from app.services.model_loader import model_registry

router = APIRouter(tags=["Health & Status"])

@router.get("/health", response_model=ResponseWrapper[HealthResponse])
def get_health_status():
    data = HealthResponse(
        status="healthy",
        version=settings.VERSION,
        loaded_models=model_registry.get_loaded_names(),
        model_details=model_registry.get_status_summary()
    )
    return ResponseWrapper(success=True, data=data)
