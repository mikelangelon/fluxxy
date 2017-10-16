docker build -t fluxxy .
docker tag fluxxy:latest 118472457945.dkr.ecr.us-east-2.amazonaws.com/fluxxy:latest

#docker run -p 8000:8080 fluxxy