package yandex_cloud

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"log"
	"os"
	"strings"
)

func CreateWebsite() {
	if err := configureBucketACL(bucket); err != nil {
		exitErrorf("Ошибка при настройке публичного доступа к бакету: %v", err)
	}

	// Настройка хостинга статического сайта
	indexPage := "index.html"
	errorPage := "error.html"
	if err := configureBucketWebsite(bucket, indexPage, errorPage); err != nil {
		exitErrorf("Ошибка при настройке хостинга статического сайта: %v", err)
	}

	albums, _ := GetAllRootFolders()

	for i, album := range albums {
		html := generateAlbumPage(album) // Здесь вы должны реализовать свою функцию генерации HTML страницы альбома
		pageKey := fmt.Sprintf("album%d.html", i)
		if err := saveObjectToBucket(bucket, pageKey, html); err != nil {
			exitErrorf("Ошибка при сохранении HTML страницы альбома: %v", err)
		}
	}

	// Генерация и сохранение HTML документа для индексной страницы
	indexHTML := generateIndexPage() // Здесь вы должны реализовать свою функцию генерации HTML индексной страницы
	if err := saveObjectToBucket(bucket, indexPage, indexHTML); err != nil {
		exitErrorf("Ошибка при сохранении HTML индексной страницы: %v", err)
	}

	// Генерация и сохранение HTML документа для страницы ошибки
	errorHTML := generateErrorPage() // Здесь вы должны реализовать свою функцию генерации HTML страницы ошибки
	if err := saveObjectToBucket(bucket, errorPage, errorHTML); err != nil {
		exitErrorf("Ошибка при сохранении HTML страницы ошибки: %v", err)
	}

	// Вывод ссылки на сайт
	fmt.Printf("Ссылка на сайт: https://%s.website.yandexcloud.net/\n", bucket)
}

func generateIndexPage() string {
	folderNames, _ := GetAllRootFolders()

	links := ""
	for i, name := range folderNames {
		links = links + fmt.Sprintf("<li><a href=\"album%d.html\">%s</a></li>\n", i, name)
	}

	return fmt.Sprintf(`<!doctype html>
<html>
    <head>
        <title>Фотоархив</title>
 		<meta content="text/html; charset=UTF-8" http-equiv="Content-Type">
    </head>
<body>
    <h1>Фотоархив</h1>
    <ul>
        %s
    </ul>
</body`, links)
}

func generateErrorPage() string {
	return fmt.Sprintf(`<!doctype html>
<html>
    <head>
		<meta content="text/html; charset=UTF-8" http-equiv="Content-Type">
        <title>Фотоархив</title>
    </head>
<body>
    <h1>Ошибка</h1>
    <p>Ошибка при доступе к фотоархиву. Вернитесь на <a href="index.html">главную страницу</a> фотоархива.</p>
</body>
</html>`)
}

func generateAlbumPage(album string) string {
	prefix := album + "/"
	listInput := &s3.ListObjectsV2Input{
		Bucket: &bucket,
		Prefix: &prefix,
	}
	listOutput, err := client.ListObjectsV2(context.TODO(), listInput)

	if err != nil {
		log.Fatalf("failed to list objects: %v", err)
	}

	var files []string
	links := ""

	for _, object := range listOutput.Contents {
		sTemp := *object.Key
		if sTemp[len(sTemp)-5:] == ".jpeg" || sTemp[len(sTemp)-4:] == ".jpg" {
			files = append(files, *object.Key)
			links = links + fmt.Sprintf("<img src=\"https://%s.website.yandexcloud.net/%s\" data-title=\"%s\">", bucket, sTemp, sTemp)
		}
	}

	return fmt.Sprintf(`<!doctype html>
<html>
    <head>
        <link rel="stylesheet" type="text/css" href="https://cdnjs.cloudflare.com/ajax/libs/galleria/1.6.1/themes/classic/galleria.classic.min.css" />
        <style>
            .galleria{ width: 960px; height: 540px; background: #000 }
        </style>
        <script src="https://ajax.googleapis.com/ajax/libs/jquery/3.6.0/jquery.min.js"></script>
        <script src="https://cdnjs.cloudflare.com/ajax/libs/galleria/1.6.1/galleria.min.js"></script>
        <script src="https://cdnjs.cloudflare.com/ajax/libs/galleria/1.6.1/themes/classic/galleria.classic.min.js"></script>
		<meta content="text/html; charset=UTF-8" http-equiv="Content-Type">
    </head>
    <body>
        <div class="galleria">
            %s
        </div>
        <p>Вернуться на <a href="index.html">главную страницу</a> фотоархива</p>
        <script>
            (function() {
                Galleria.run('.galleria');
            }());
        </script>
    </body>
</html>`, links)

	//// Перебираем каждый объект в списке
	//for _, obj := range resp.Contents {
	//	// Формируем публичный URL для каждого объекта
	//	url := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucket, *obj.Key)
	//
	//	fmt.Printf("File: %s, URL: %s\n", *obj.Key, url)
	//}
}

// Функция для настройки публичного доступа на чтение объектов бакета
func configureBucketACL(bucketName string) error {
	_, err := client.PutBucketAcl(context.TODO(), &s3.PutBucketAclInput{
		ACL:    types.BucketCannedACLPublicRead,
		Bucket: aws.String(bucketName),
	})
	return err
}

// Функция для настройки хостинга статического сайта
func configureBucketWebsite(bucketName, indexPage, errorPage string) error {
	_, err := client.PutBucketWebsite(context.TODO(), &s3.PutBucketWebsiteInput{
		Bucket: aws.String(bucketName),
		WebsiteConfiguration: &types.WebsiteConfiguration{
			IndexDocument: &types.IndexDocument{
				Suffix: aws.String(indexPage),
			},
			ErrorDocument: &types.ErrorDocument{
				Key: aws.String(errorPage),
			},
		},
	})
	return err
}

// Функция для сохранения объекта (HTML страницы) в бакете
func saveObjectToBucket(bucketName, key, html string) error {
	_, err := client.PutObject(context.TODO(), &s3.PutObjectInput{
		Body:        strings.NewReader(html),
		Bucket:      aws.String(bucketName),
		Key:         aws.String(key),
		ContentType: aws.String("text/html"),
	})
	return err
}

// Функция для обработки ошибок и вывода сообщения перед выходом
func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
