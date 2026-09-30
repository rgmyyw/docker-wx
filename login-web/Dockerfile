FROM node:22-alpine
WORKDIR /app
ENV TZ=Asia/Shanghai
COPY package.json ./
RUN npm install --omit=dev --no-audit --no-fund
COPY server.js ./
EXPOSE 8058
CMD ["node", "server.js"]
