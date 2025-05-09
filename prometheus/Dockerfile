
FROM node:lts-alpine3.21

WORKDIR /app

COPY package*.json ./

RUN npm install

COPY . .

CMD ["node","--watch", "server.js"]
