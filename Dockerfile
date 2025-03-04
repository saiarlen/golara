FROM golang:1.23

# All Shell Binaries
RUN apt-get update && apt-get install -y ghostscript imagemagick ffmpeg poppler-utils fonts-open-sans

#mods
RUN sed -i 's/<policy domain="coder" rights="none" pattern="PDF"/<policy domain="coder" rights="read|write" pattern="PDF"/g' /etc/ImageMagick-6/policy.xml


WORKDIR /ekycdir

# Copy the go mod and sum files
COPY go.mod go.sum ./

RUN go mod download
RUN go mod tidy

RUN go install github.com/air-verse/air@latest

COPY . .

#Pull Binaries
RUN cd xbin && wget https://scripts.appxcube.in/gobin/wkhtmltopdf-amd64 && chmod +x wkhtmltopdf-amd64


EXPOSE 9000
CMD ["air"]


