FROM --platform=linux/amd64 python:3.11.9-slim AS app

ENV PYTHONUNBUFFERED=1 \
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
RUN poetry install --no-cache --no-dev

# Install resources
WORKDIR /app/resources
RUN wget https://huggingface.co/prithivida/flashrank/resolve/main/rank-T5-flan.zip -O rank-T5-flan.zip && \
    unzip rank-T5-flan.zip -d rank-T5-flan && \
    rm -rf rank-T5-flan.zip

WORKDIR /app

# Copy source files
COPY src ./src

# Copy other resources
COPY artifacts ./artifacts

# Run container as non-root
RUN groupadd -g 1000 app && \
    useradd -u 1000 -g app -d /app -s /bin/sh -m app && \
    chown -R app:app /app

# Api
EXPOSE 2222

USER app

CMD [ "gunicorn", "src.main:server", "--bind", "0.0.0.0:2222", "--workers", "1", "--worker-class", "uvicorn.workers.UvicornWorker", "--graceful-timeout", "30" ]
