package docker

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
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
	}
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
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

// Download Layer
func (d *DockerClient) GetLayer(url string, headers http.Header) error {

	req, err := createRequest("GET", url, headers)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// check if the body is right
	return nil
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

	return manifest, nil
}

// Pull a docker image
func (d *DockerClient) Pull(registry, image string) (controller.Report, error) {
	// Get manifest of the image

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

	for _, layer := range manifest.Layers {

		layerUrl := fmt.Sprintf("%s/v2/%s/blobs/%s", registry, imageName, layer.Digest)
		// Remove Accept header as it's not needed for layers
		headers["Accept"] = []string{}

		logger.DebugLogger.Printf("Downloading layer: %s\n", layerUrl)
		start := time.Now()
		err := d.GetLayer(layerUrl, headers)
		if err != nil {
			report.FailedLayers += 1
		} else {
			report.DownloadedLayers += 1
		}
		secondsElapsed := float64(time.Now().Sub(start)) / 1000000000
		avgTime = append(avgTime, secondsElapsed)
	}

	sumArray := func(times []float64) float64 {
		sumItems := 0.0
		for _, value := range times {
			sumItems += value
		}
		return sumItems
	}

	report.AvgDownloadTime = sumArray(avgTime) / float64(len(avgTime))

	if report.FailedLayers == 0 {
		report.Success = true
	}
	//create report

	return report, nil
}
