FROM python:3.11.9-slim

ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1 \
    POETRY_NO_INTERACTION=1 \
    POETRY_VIRTUALENVS_CREATE=0

# Setup system
RUN apt update && \
    apt install --yes build-essential ca-certificates wget unzip && \
    rm -rf /var/lib/apt/lists/*

# Setup poetry package manager
RUN pip install poetry==1.8.3
RUN poetry config installer.parallel false

# Install dependencies
COPY pyproject.toml poetry.lock ./
RUN poetry install --no-cache

# Install resources
WORKDIR /app/resources
RUN wget https://huggingface.co/prithivida/flashrank/resolve/main/rank-T5-flan.zip -O rank-T5-flan.zip && \
    unzip rank-T5-flan.zip -d rank-T5-flan && \
    rm -rf rank-T5-flan.zip

WORKDIR /app

# Copy other resources
COPY artifacts ./artifacts

# Api
EXPOSE 2222

# Phoenix
EXPOSE 2223

CMD [ "uvicorn", "src.main:server", "--host", "0.0.0.0", "--port", "2222", "--timeout-graceful-shutdown", "30", "--reload" ]
