from fastapi import APIRouter
from app.api.v1 import health, demand, stockout, courier, anomaly

api_v1_router = APIRouter()

api_v1_router.include_router(health.router)
api_v1_router.include_router(demand.router)
api_v1_router.include_router(stockout.router)
api_v1_router.include_router(courier.router)
api_v1_router.include_router(anomaly.router)
