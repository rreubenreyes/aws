package participant

import (
	"bufio"
	"bytes"
	"image/png"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/nfnt/resize"
)

var staticBucket = os.Getenv("STATIC_BUCKET_ID")

func (p *Participant) PhotoPNG(svc *s3.S3) ([]byte, error) {
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

	return dlbuf.Bytes(), nil
}

func WidenPhotoPNG(b []byte, dx, dy float32) ([]byte, error) {
	photo, err := png.Decode(bytes.NewReader(b))
	if err != nil {
		log.Printf("could not decode png")
		log.Println(err)
		return nil, err
	}
	size := photo.Bounds().Size()
	log.Println("intake photo size", size.X, size.Y)
	rwidth := uint(float32(size.X) * dx)
	rheight := uint(float32(size.Y) * dy)
	resized := resize.Resize(rwidth, rheight, photo, resize.Lanczos3)

	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)
	err = png.Encode(w, resized)
	if err != nil {
		log.Println("could not encode photo")
		log.Println(err)
		return nil, err
	}

	return buf.Bytes(), nil
}
