FROM golang:1.14-buster

CMD /bin/init.sh /bin/rss-bot

COPY . /code

WORKDIR /code

RUN go get ./...

RUN make gh_linux_amd64
RUN chmod -R g+rwx /code && cp build/telegram-rss-bot-linux-amd64 /bin/rss-bot && cp init.sh /bin