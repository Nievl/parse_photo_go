package services

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"parse_photo_go/domains"
	"parse_photo_go/models"
	"parse_photo_go/utils"

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
		isTelegraph = "https://telegra.ph"
	}
	urls := getMediaUrls(page, isOsosedki, isTelegraph)
	dirPath, err := createDirectory(link.Name)
	if err != nil {
		return fmt.Errorf("failed to create directory: %s", err.Error())
	}

	var (
		downloadedCount      int64
		downloadedMediafiles = make([]models.CreateMediafileDto, 0, len(urls))
		downloadedMu         sync.Mutex
		wg                   sync.WaitGroup
		sem                  = make(chan struct{}, 5) // лимит 5 одновременных загрузок
	)

	for _, url := range urls {
		url := url
		wg.Add(1)
		sem <- struct{}{} // блокируем семафор
		go func() {
			defer wg.Done()
			defer func() { <-sem }() // освобождаем семафор

			cleanUrl := strings.Split(url, "?")[0]
			fileName := filepath.Base(cleanUrl)
			ext := strings.TrimPrefix(filepath.Ext(fileName), ".")
			if _, ok := EXTENSIONS[ext]; !ok {
				return
			}

			filePath := filepath.Join(dirPath, fileName)
			fullUrl := url
			if !strings.HasPrefix(url, "http") {
				fullUrl = link.Path + url
			}

			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				fullUrl = getHighResUrl(fullUrl)
				res, err := s.mediafilesService.DownloadFile(fullUrl, filePath, link.ID)
				if err != nil {
					fmt.Printf("failed to download file: %s", err.Error())
				} else {
					downloadedMu.Lock()
					downloadedMediafiles = append(downloadedMediafiles, res)
					downloadedCount++
					downloadedMu.Unlock()
				}

			} else {
				atomic.AddInt64(&downloadedCount, 1)

			}
		}()
	}

	wg.Wait()
	totalFiles := len(urls)

	for _, mediafile := range downloadedMediafiles {
		err := s.mediafilesService.Create(mediafile)
		if err != nil {
			fmt.Printf("failed to create mediafile: %s\n", err.Error())
		}
	}

	linkDto := models.UpdateLinkDto{
		IsDownloaded:         downloadedCount == int64(totalFiles),
		Progress:             int((float64(downloadedCount) / float64(totalFiles)) * 100),
		Mediafiles:           totalFiles,
		DownloadedMediafiles: int(downloadedCount),
	}

	s.linkDbService.UpdateFilesNumber(id, linkDto)

	for _, mediafile := range downloadedMediafiles {
		err := s.mediafilesService.Create(mediafile)
		if err != nil {
			fmt.Printf("failed to create mediafile: %s\n", err.Error())
		}
	}

	return nil
}

func (s *LinkService) ScanFilesForLink(id int64) (string, error) {
	link, err := s.linkDbService.GetOne(id)
	if err != nil {
		return "", fmt.Errorf("failed to get link: %s", err.Error())
	}
	dirPath := filepath.Join("result", link.Name)
	_, err = os.Stat(dirPath)
	if os.IsNotExist(err) {
		return "", fmt.Errorf("directory %s does not exist", dirPath)
	}
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return "", fmt.Errorf("failed to read directory: %s", err.Error())
	}
	if len(files) > 0 {
		existedMediafiles, err := s.mediafilesService.GetAllByLinkId(link.ID)
		if err != nil {
			return "", fmt.Errorf("failed to get mediafiles: %s", err.Error())
		}
		existedMediaFilesSet := make(map[string]struct{}, len(existedMediafiles))
		for _, mediafile := range existedMediafiles {
			existedMediaFilesSet[mediafile.Name] = struct{}{}
		}
		for _, file := range files {
			fileName := file.Name()
			if _, ok := existedMediaFilesSet[fileName]; !ok {
				filePath := filepath.Join(dirPath, fileName)
				hash, err := utils.GetHashByPath(filePath)
				if err != nil {
					fmt.Printf("failed to get hash: %s", err.Error())
				}
				info, _ := os.Stat(filePath)
				mediaFile := models.CreateMediafileDto{
					Name:   fileName,
					Path:   filePath,
					Hash:   hash,
					Size:   info.Size(),
					LinkID: link.ID,
				}
				err = s.mediafilesService.Create(mediaFile)
				if err != nil {
					fmt.Printf("failed to create mediafile: %s", err.Error())
				}
			}

		}

	} else {
		return "", fmt.Errorf("directory %s is empty", dirPath)
	}

	return fmt.Sprintf("files link %s scanned"), nil
}

func (s *LinkService) CheckDownloaded(id int64) (string, error) {
	link, err := s.linkDbService.GetOne(id)
	if err != nil {
		return "", fmt.Errorf("failed to get link: %s", err.Error())
	}

	page, _ := getPage(link.Path)
	dirPath := filepath.Join("result", link.Name)
	dir, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		dir = nil
	}

	if dir != nil && page == nil {
		files, err := os.ReadDir(dirPath)
		if err != nil {
			return "", fmt.Errorf("failed to read directory: %s", err.Error())
		}
		if len(files) > 0 {
			linkDto := models.UpdateLinkDto{
				IsDownloaded:         true,
				Progress:             100,
				Mediafiles:           len(files),
				DownloadedMediafiles: len(files),
			}
			s.linkDbService.UpdateFilesNumber(id, linkDto)
			return fmt.Sprintf("%d files in %s and \n page not found", len(files), link.Path), nil
		} else {
			return fmt.Sprintf("0 files in %s and \n page not found", link.Path), nil
		}
	}

	isOsosedki := isOsosedkiDomain(link.Path)
	var isTelegraph string
	if strings.Contains(link.Path, "telegra.ph") {
		isTelegraph = "https://telegra.ph"
	}

	if dir != nil && page != nil {
		urls := getMediaUrls(page, isOsosedki, isTelegraph)
		files, err := os.ReadDir(dirPath)
		if err != nil {
			return "", fmt.Errorf("failed to read directory: %s", err.Error())
		}
		linkDto := models.UpdateLinkDto{
			IsDownloaded:         len(files) == len(urls),
			Progress:             int((float64(len(files)) / float64(len(urls))) * 100),
			Mediafiles:           len(urls),
			DownloadedMediafiles: len(files),
		}
		s.linkDbService.UpdateFilesNumber(id, linkDto)
		return fmt.Sprintf("%d files in %s and \n page exists", len(files), link.Path), nil
	}
	if dir == nil && page != nil {
		urls := getMediaUrls(page, isOsosedki, isTelegraph)
		linkDto := models.UpdateLinkDto{
			IsDownloaded:         false,
			Progress:             0,
			Mediafiles:           len(urls),
			DownloadedMediafiles: 0,
		}
		s.linkDbService.UpdateFilesNumber(id, linkDto)
		return fmt.Sprintf("0 files in %s and \n page exists", link.Path), nil

	}

	return fmt.Sprintf("%s  is not exists and \n page not found", link.Path), nil

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
