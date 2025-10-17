import uvicorn
from fastapi import FastAPI
from .api.handler import router as sso_router

app = FastAPI()

ROUTERS = [
    sso_router,
]

for router in ROUTERS:
    app.include_router(router, prefix="/api")

if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=4001, reload=true)