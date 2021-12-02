package docker

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/ish-xyz/regtester/pkg/config"
	"github.com/ish-xyz/regtester/pkg/controller"
	"github.com/ish-xyz/regtester/pkg/logger"
)

var client = http.Client{}

func NewDockerClient(c config.Config) *DockerClient {
	return &DockerClient{
		Username:        c.Connection.BasicAuth.Username,
		Password:        c.Connection.BasicAuth.Password,
		ExtraHeaders:    c.Connection.ExtraHeaders,
		ManifestVersion: c.Connection.ManifestVersion,
		CheckIntegrity:  c.Workload.CheckIntegrity,
		MaxConcLayers:   c.Workload.MaxConcLayers,
	}
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func sumDownloadTimeArray(times []float64) float64 {
	sumItems := 0.0
	for _, value := range times {
		sumItems += value
	}
	return sumItems
}

// Create http request
func createRequest(method, url string, headers http.Header) (*http.Request, error) {

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header = headers

	return req, nil
}

// Download a Docker Layer
func (d *DockerClient) GetLayer(url string, headers http.Header) (float64, error) {

	//TODO implement retry
	req, err := createRequest("GET", url, headers)
	if err != nil {
		return 0.0, err
	}

	start := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		return (float64(time.Since(start)) / 1000000000), err

	}
	defer resp.Body.Close()

	secondsElapsed := float64(time.Since(start)) / 1000000000

	if d.CheckIntegrity {
		logger.DebugLogger.Println("Running integrity Check")
		// Get downloaded layer sha256
		layerData, _ := ioutil.ReadAll(resp.Body)
		layerSHA25632Bytes := sha256.Sum256(layerData)
		layerSHA256 := hex.EncodeToString(layerSHA25632Bytes[:])

		// Get desired sha256
		urlSplit := strings.Split(url, "/")
		desiredSHA256 := strings.Split(urlSplit[len(urlSplit)-1:][0], ":")[1]
		fmt.Println(desiredSHA256)

		if desiredSHA256 != layerSHA256 {
			logger.WarningLogger.Println("Integrity check failed!")
			return secondsElapsed, fmt.Errorf(
				"checksum failed for layer sha256:%s",
				desiredSHA256,
			)
		}
	}

	return secondsElapsed, nil
}

// Get the manifest of a Docker image from a given registry
func (d *DockerClient) GetManifest(url string, headers http.Header) (ImageManifestV2, error) {

	var manifest ImageManifestV2

	req, err := createRequest("GET", url, headers)
	if err != nil {
		return ImageManifestV2{}, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return ImageManifestV2{}, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &manifest)

	return manifest, err
}

// Pull a Docker image
func (d *DockerClient) Pull(registry, image string) (controller.Report, error) {

	var wg sync.WaitGroup

	maxConcLayers := make(chan int, d.MaxConcLayers)
	report := controller.Report{
		Image:            image,
		Registry:         registry,
		FailedLayers:     0,
		DownloadedLayers: 0,
		AvgDownloadTime:  0.0,
		Success:          false,
		ManifestDownload: true,
	}
	avgTime := []float64{}
	headers := make(http.Header)
	imageParts := strings.Split(image, ":")
	imageName, tag := imageParts[0], imageParts[1]
	manifestUrl := fmt.Sprintf("%s/v2/%s/manifests/%s", registry, imageName, tag)

	headers["Accept"] = []string{"application/vnd.docker.distribution.manifest.v2+json"}
	if d.Username != "" && d.Password != "" {
		headers["Authorization"] = []string{"Basic " + basicAuth(d.Username, d.Password)}
	}

	manifest, err := d.GetManifest(manifestUrl, headers)
	if err != nil {
		report.ManifestDownload = true
		return report, err
	}

	// Download layers in parallel with a limit
	for index, layer := range manifest.Layers {

		maxConcLayers <- index
		wg.Add(1)

		go func(wg *sync.WaitGroup, layer Layer) {

			defer wg.Done()

			layerUrl := fmt.Sprintf("%s/v2/%s/blobs/%s", registry, imageName, layer.Digest)
			headers["Accept"] = []string{}

			logger.DebugLogger.Printf("Downloading layer: %s\n", layerUrl)
			_, err := d.GetLayer(layerUrl, headers)
			if err != nil {
				report.FailedLayers += 1
			} else {
				report.DownloadedLayers += 1
			}

			<-maxConcLayers
		}(&wg, layer)
	}

	wg.Wait()
	/*
		for recording := range recordExecutionTimes {
			avgTime = append(avgTime, recording)
		}
	*/
	report.AvgDownloadTime = sumDownloadTimeArray(avgTime) / float64(len(avgTime))

	if report.FailedLayers == 0 {
		report.Success = true
	}
	//create report

	return report, nil
}
