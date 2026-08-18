import logging
from contextlib import asynccontextmanager
from fastapi import FastAPI, Request
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import JSONResponse

from app.config import settings
from app.api.v1.router import api_v1_router
from app.services.model_loader import model_registry
from app.schemas.base import ResponseWrapper

# Configure Logging
logging.basicConfig(
    level=settings.LOG_LEVEL,
    format="%(asctime)s - %(name)s - %(levelname)s - %(message)s"
)
logger = logging.getLogger("bootloader")

@asynccontextmanager
async def lifespan(app: FastAPI):
    """
    Application lifespan context manager:
    Executes model checkpoint loading into system memory on startup.
    """
    logger.info("Initializing intelligence layer...")
    model_registry.load_all_models()
    yield
    logger.info("Shutting down AI Lab App Server...")

app = FastAPI(
    title=settings.APP_NAME,
    version=settings.VERSION,
    description="Real-time ML Model Inference & Anomaly Detection Server for chia.florist staff backend (service-core).",
    lifespan=lifespan,
    docs_url="/docs",
    redoc_url="/redoc"
)

# CORS Middleware Setup
if settings.ALLOWED_ORIGINS:
    app.add_middleware(
        CORSMiddleware,
        allow_origins=settings.ALLOWED_ORIGINS,
        allow_credentials=True,
        allow_methods=["*"],
        allow_headers=["*"],
    )

# Include API v1 Routes
app.include_router(api_v1_router, prefix=settings.API_V1_PREFIX)

@app.exception_handler(Exception)
async def global_exception_handler(request: Request, exc: Exception):
    logger.error(f"Unhandled exception on path {request.url.path}: {exc}", exc_info=True)
    return JSONResponse(
        status_code=500,
        content=ResponseWrapper(
            success=False,
            error=f"Internal Server Error: {str(exc)}"
        ).model_dump()
    )

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(
        "app.main:app",
        host=settings.HOST,
        port=settings.PORT,
        reload=True
    )
