FROM node:14.15.0
MAINTAINER rich<rich@salt-mine.io>

# 도커안의 앱 디렉터리 생성
WORKDIR /app

COPY . /app

RUN npm install -g pm2 node-gyp
RUN npm install
RUN ln -sf /usr/share/zoneinfo/Asia/Seoul /etc/localtime

# 앱 소스 추가
COPY . .

EXPOSE 80 443 3000
CMD [ "pm2-runtime", "start", "ecosystem.config.js" ]
