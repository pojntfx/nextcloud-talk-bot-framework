FROM node

WORKDIR /app

COPY package-lock.json package.json ./

RUN npm install

COPY main.js proc.js ./

CMD node main.js