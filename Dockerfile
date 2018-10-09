FROM python:3.6-alpine

ENV LANG=en_US.utf8

# Pipenv
RUN pip3 install pipenv

# Add code
COPY . /code

# Install dependencies
RUN cd /code && pipenv install

# Service
WORKDIR /code/validator
ENV HOST=0.0.0.0 PORT=80
EXPOSE 80
CMD ["/usr/local/bin/pipenv", "run", "python", "server.py"]
