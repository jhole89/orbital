FROM node
RUN mkdir /code

COPY *.json /code/
COPY *.html /code/
COPY webpack.config.js /code/

COPY src /code/src/

WORKDIR /code

RUN npm install

ARG ORBITAL_SERVER_ADDRESS
ENV ORBITAL_SERVER_ADDRESS=${ORBITAL_SERVER_ADDRESS}

RUN npm run build
CMD ["npm", "run", "server"]
