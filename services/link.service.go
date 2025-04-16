package services

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"parse_photo_go/domains"
	"parse_photo_go/models"

	"github.com/PuerkitoBio/goquery"
)

type LinkService struct {
	linkDbService     domains.LinksDbService
	mediafilesService MediafilesService
}

func NewLinkService(linkDbService domains.LinksDbService, mediafilesService MediafilesService) *LinkService {
	return &LinkService{linkDbService, mediafilesService}
}

func (s *LinkService) Create(link models.CreateLinkDto) error {
	url, err := checkURL(link.Path)

	if err != nil {

		return fmt.Errorf("path is not a URL")
	} else {
		return s.linkDbService.CreateLink(link.Path, url[1])
	}

}

func (s *LinkService) GetAll(isReachable, showDuplicate bool) ([]models.LinkWithDuplicatePath, error) {
	return s.linkDbService.GetAll(isReachable, showDuplicate)
}

func (s *LinkService) Remove(id int64) error {
	return s.linkDbService.Remove(id)
}

func (s *LinkService) DownloadFiles(id int64) error {
	link, err := s.linkDbService.GetOne(id)
	if err != nil {
		return fmt.Errorf("failed to get link: %s", err.Error())
	}
	page, err := getPage(link.Path)
	if err != nil {
		return fmt.Errorf("failed to fetch page: %s", err.Error())
	}
	isOsosedki := isOsosedkiDomain(link.Path)
	var isTelegraph string
	if strings.Contains(link.Path, "telegra.ph") {
		isTelegraph = "telegra.ph"
	}
	urls := getMediaUrls(page, isOsosedki, isTelegraph)
	dirPath, err := createDirectory(link.Name)
	if err != nil {
		return fmt.Errorf("failed to create directory: %s", err.Error())
	}
	var downloadedCount, totalFiles int

	for _, url := range urls {
		cleanUrl := strings.Split(url, "?")[0]
		fileName := filepath.Base(cleanUrl)
		ext := strings.TrimPrefix(filepath.Ext(fileName), ".")
		if _, ok := EXTENSIONS[ext]; !ok {
			continue
		}
		totalFiles++
		filePath := filepath.Join(dirPath, fileName)
		fullUrl := url
		if !strings.HasPrefix(url, "http") {
			fullUrl = link.Path + url
		}

		downloadedMediafiles := make([]models.CreateMediafileDto, 0, len(urls))

		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			fullUrl = getHighResUrl(fullUrl)
			res, err := s.mediafilesService.DownloadFile(fullUrl, filePath, link.ID)
			if err != nil {
				log.Fatalf("failed to download file: %s", err.Error())
			}
			downloadedMediafiles = append(downloadedMediafiles, res)
			downloadedCount++

		} else {
			downloadedCount++
			continue
		}

	}

	return nil
}

func (s *LinkService) ScanFilesForLink(id int64) error {
	// implementation for scanning files for a link by id
	return nil
}

func (s *LinkService) CheckDownloaded(id int64) (bool, error) {
	// implementation for checking if files are downloaded for a link by id
	return false, nil
}

func (s *LinkService) TagUnreachable(id int64, reachable bool) error {
	return s.linkDbService.TagUnreachable(id, reachable)
}

func checkURL(url string) ([]string, error) {
	re := regexp.MustCompile(`(http[s]?:\/\/[^\/\s]+\/)(.*)`)
	urlMatches := re.FindStringSubmatch(url)

	if len(urlMatches) == 3 {
		return urlMatches[1:], nil
	}
	return nil, fmt.Errorf("invalid URL format")
}

func isOsosedkiDomain(url string) bool {
	return strings.Contains(url, "ososedki.com")
}

func getMediaUrls(page *goquery.Document, absoluteOnly bool, domain string) []string {
	urlSet := make(map[string]struct{})

	page.Find("img, video").Each(func(i int, media *goquery.Selection) {
		src, exists := media.Attr("src")
		if !exists || src == "" {
			return
		}

		isAbsolute := strings.HasPrefix(src, "http://") || strings.HasPrefix(src, "https://")
		if !absoluteOnly || isAbsolute {
			if domain != "" && !isAbsolute {
				src = domain + src
			}
			urlSet[src] = struct{}{}
		}
	})
	urls := make([]string, 0, len(urlSet))
	for url := range urlSet {
		urls = append(urls, url)
	}

	return urls
}

func getPage(url string) (*goquery.Document, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch page: %s", resp.Status)
	}

	return goquery.NewDocumentFromReader(resp.Body)
}

func createDirectory(name string) (string, error) {
	dirPath := filepath.Join("result", name)
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return "", fmt.Errorf("%s is not exists and could not create directory: %w", dirPath, err)
		}
	}
	return dirPath, nil
}

func getHighResUrl(url string) string {
	highResUrl := strings.Replace(url, "/a/604/", "/a/1280/", 1)
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Head(highResUrl)
	if err != nil {
		return url
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		return highResUrl
	}
	return url

}

var EXTENSIONS = map[string]struct{}{"jpeg": {}, "jpg": {}, "mp4": {}, "png": {}, "gif": {}, "webp": {}}
