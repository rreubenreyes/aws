package participant

import (
	"bufio"
	"bytes"
	"image/png"
	"io"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/nfnt/resize"
)

const (
	PHOTO_X_SCALE = 1.2
	PHOTO_Y_SCALE = 0.8
)

var staticBucket = os.Getenv("STATIC_BUCKET_ID")

func (p *Participant) PhotoPNG(svc *s3.S3) (io.Reader, error) {
	log.Printf("getting participant photo for %s\n", p.Name)

  // get photo
  dlbuf := &aws.WriteAtBuffer{}
	downloader := s3manager.NewDownloaderWithClient(svc)
	_, err := downloader.Download(dlbuf, &s3.GetObjectInput{
		Bucket: aws.String(staticBucket),
		Key:    aws.String(p.PhotoS3Key),
	})
	if err != nil {
		log.Println("could not get photo")
    log.Println(err)
		return nil, err
	}
  data := dlbuf.Bytes()

  // resize photo
  rdData := bytes.NewReader(data)
  // rdData.Seek(0,0)
	photo, err := png.Decode(rdData)
	if err != nil {
		log.Printf("could not decode png")
    log.Println(err)
		return nil, err
	}
	size := photo.Bounds().Size()
  log.Println("intake photo size", size.X, size.Y)
	rwidth := uint(float32(size.X) * PHOTO_X_SCALE)
	rheight := uint(float32(size.Y) * PHOTO_Y_SCALE)
	resized := resize.Resize(rwidth, rheight, photo, resize.Lanczos3)

  // upload resized photo to s3
	var buf bytes.Buffer
	wrResized := bufio.NewWriter(&buf)
	err = png.Encode(wrResized, resized)
	if err != nil {
		log.Println("could not encode photo")
    log.Println(err)
		return nil, err
	}

	var encoded []byte
	_, err = wrResized.Write(encoded)
	if err != nil {
		log.Println("could not write photo to in-memory buffer")
    log.Println(err)
		return nil, err
	}

	rdEncoded := bytes.NewReader(encoded)
	log.Printf("uploading resized participant photo for %s\n", p.Name)
	uploader := s3manager.NewUploaderWithClient(svc)
  _, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(staticBucket),
		Key:    aws.String(p.PhotoS3Key),
		Body:   rdEncoded,
  })
	if err != nil {
		log.Println("could not write resized photo")
    log.Println(err)
		return nil, err
	}

	return rdEncoded, nil
}
