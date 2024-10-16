from typing import Optional

from pydantic import BaseModel

from src.config import Config


class HealthEndpoints:
    def __init__(self, config: Config) -> None:
        self.config = config

    class GetHealthResponse(BaseModel):
        class Item(BaseModel):
            error: Optional[str]
            latency: int

        system: Item

    async def get_health(self) -> GetHealthResponse:
        return self.GetHealthResponse(
            system=self.GetHealthResponse.Item(
                error=None,
                latency=0,
            ),
        )
