FROM python:3.9-slim AS builder-image

# create and activate virtual environment
# using final folder name to avoid path issues with packages
RUN python -m venv /home/myuser/venv
ENV PATH="/home/myuser/venv/bin:$PATH"

# install requirements
COPY python-libs.txt .
RUN pip install --no-cache-dir wheel
RUN pip install --no-cache-dir -r python-libs.txt

FROM python:3.9-slim AS runner-image

RUN useradd --create-home myuser
COPY --from=builder-image /home/myuser/venv /home/myuser/venv

USER myuser
RUN mkdir /home/myuser/code
WORKDIR /home/myuser/code
COPY . .

EXPOSE 6000

# make sure all messages always reach console
ENV PYTHONUNBUFFERED=1

# activate virtual environment
ENV VIRTUAL_ENV=/home/myuser/venv
ENV PATH="/home/myuser/venv/bin:$PATH"
# In your Dockerfile
ENTRYPOINT ["sh", "-c"]
CMD ["python main.py --port=${PORT:-6000} --debug=${DEBUG:-true} --nats_url=${NATS_URL:-nats://nats-server:4222}"]