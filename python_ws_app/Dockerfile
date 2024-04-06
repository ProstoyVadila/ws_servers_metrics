FROM python:3.11-slim as poetry-stage

ENV POETRY_VERSION=1.6.1 \
    POETRY_HOME="/opt/poetry" \
    POETRY_NO_INTERACTION=1

ENV PATH="$POETRY_HOME/bin:$PATH"

RUN apt-get update && \
    apt-get install --no-install-recommends -y curl && \
    curl -sSL https://install.python-poetry.org | python3 -

WORKDIR /poetry_deps

COPY ./poetry.lock /poetry_deps
COPY ./pyproject.toml /poetry_deps
RUN poetry export --without-hashes -f requirements.txt --output requirements.txt

FROM python:3.11-slim as main-stage

WORKDIR /app
# Install any needed packages specified in requirements.txt
COPY --from=poetry-stage /poetry_deps/requirements.txt /app
RUN pip install --prefer-binary -r /app/requirements.txt

# Copy the current directory contents into the container at /app
COPY ./app/ /app

# store revision
ARG REVISION
RUN echo $REVISION > /app/revision.info


ENV LC_ALL ru_RU.UTF-8
ENV LANG ru_RU.UTF-8
ENV LANGUAGE ru_RU.UTF-8

# Run app.py when the container launches
CMD [ "python", "-u", "app.py" ]
