FROM python:3.9-slim
WORKDIR /app
COPY producer.py /app/
RUN pip install kafka-python
CMD ["python", "producer.py"]