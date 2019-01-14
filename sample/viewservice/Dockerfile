FROM python:3-slim

RUN pip3 install -q Flask==0.11.1 requests==2.18.4
COPY viewservice.py /opt/microservices/
EXPOSE 9080
WORKDIR /opt/microservices
CMD python3 viewservice.py
