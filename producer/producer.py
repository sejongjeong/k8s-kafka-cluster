import json
import time
import random
from kafka import KafkaProducer

# Kafka Producer 생성
producer = KafkaProducer(
    bootstrap_servers='study-cluster-kafka-bootstrap.kafka.svc.cluster.local:9092',
    value_serializer=lambda v: json.dumps(v).encode('utf-8')  # JSON 데이터를 직렬화
)

try:
    for i in range(1, 3601):  # 1부터 3600까지의 숫자
        # Dummy JSON 데이터 생성
        message = {
            "id": i,
            "data": f"kafka test data",
            "random_value": random.randint(1, 100)
        }

        # Kafka 토픽으로 데이터 전송
        producer.send('log', value=message)
        print(f"Sent: {message}")

        # 1초 대기
        time.sleep(1)
except KeyboardInterrupt:
    print("\nStopping producer...")
finally:
    producer.close()
