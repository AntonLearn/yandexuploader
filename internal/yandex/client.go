// Package yandex implements an optimized REST client for bulk operations against the Yandex.Disk cloud service API.
package yandex

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/antonlearn/yandex-uploader/internal/progress"
	"github.com/antonlearn/yandex-uploader/internal/report"
)

const restEndpoint = "https://cloud-api.yandex.net/v1/disk/resources"

// Client acts as an orchestration entity maintaining authentication states and HTTP transport pools.
type Client struct {
	httpClient *http.Client
	token      string
}

// uploadLinkResponse mapped schema for target API JSON response deserialization blocks.
type uploadLinkResponse struct {
	Href string `json:"href"`
}

// NewClient sets up a Client instance tailored for stable, high-throughput network operations.
func NewClient(token string) *Client {
	return &Client{
		token: token,
		httpClient: &http.Client{
			Timeout: 0, // Uncapped transfer timeout to facilitate multi-gigabyte files processing.
			Transport: &http.Transport{
				DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
					dialer := &net.Dialer{Timeout: 30 * time.Second, KeepAlive: 30 * time.Second}
					return dialer.DialContext(ctx, network, addr)
				},
				DisableCompression: true, // Suppress CPU overhead for arbitrary payloads processing.
				TLSClientConfig:    &tls.Config{InsecureSkipVerify: false},
				WriteBufferSize:    1024 * 1024, // Optimized 1MB chunk pipeline buffer for uploading.
				ReadBufferSize:     32 * 1024,
			},
		},
	}
}

// UploadFolder recurses through the local directories, replicates structures remotely, and populates the operational log registry.
func (c *Client) UploadFolder(localPath string, reportItems *[]report.Item) error {
	baseDir := filepath.Base(localPath)
	remoteRootDir := "/" + baseDir

	if err := c.createRemoteDir(remoteRootDir); err != nil {
		return fmt.Errorf("directory configuration failure on remote storage %s: %w", remoteRootDir, err)
	}

	if rootPubLink, err := c.PublishAndGetLink(remoteRootDir); err == nil {
		*reportItems = append(*reportItems, report.Item{
			Left:  fmt.Sprintf("[Folder] %s (Entire directory container tree)", baseDir),
			Right: rootPubLink,
		})
	}

	return filepath.WalkDir(localPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, _ := filepath.Rel(localPath, path)
		if rel == "." {
			return nil
		}

		remotePath := "/" + baseDir + "/" + filepath.ToSlash(rel)

		if d.IsDir() {
			errDir := c.createRemoteDir(remotePath)
			if errDir == nil {
				if pubLink, errPub := c.PublishAndGetLink(remotePath); errPub == nil {
					*reportItems = append(*reportItems, report.Item{
						Left:  fmt.Sprintf("[Folder]   ├── %s", d.Name()),
						Right: pubLink,
					})
				}
			}
			return errDir
		}

		errFile := c.UploadSingleFile(path, remotePath)
		leftText := fmt.Sprintf("[File]   ├── %s", filepath.Base(path))

		if errFile == nil {
			pubLink, errPub := c.PublishAndGetLink(remotePath)
			if errPub != nil {
				*reportItems = append(*reportItems, report.Item{Left: leftText, Right: fmt.Sprintf("(Publish failed: %v)", errPub)})
			} else {
				*reportItems = append(*reportItems, report.Item{Left: leftText, Right: pubLink})
			}
		} else {
			*reportItems = append(*reportItems, report.Item{Left: leftText, Right: fmt.Sprintf("(Transfer failed: %v)", errFile)})
		}
		return nil
	})
}

// UploadSingleFile transmits a detached single binary entity using wrapped telemetry streaming buffers.
func (c *Client) UploadSingleFile(localPath, remotePath string) error {
	file, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer file.Close()

	info, _ := file.Stat()
	fileSize := info.Size()

	var sizeStr string
	if fileSize < 1024*1024 {
		sizeStr = fmt.Sprintf("%.1f KB", float64(fileSize)/1024)
	} else {
		sizeStr = fmt.Sprintf("%.1f MB", float64(fileSize)/1024/1024)
	}

	fmt.Printf("\n[Processing] %s (%s)\n", filepath.Base(localPath), sizeStr)

	apiURL := restEndpoint + "/upload?path=" + url.QueryEscape(remotePath) + "&overwrite=true"
	reqAPI, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return err
	}
	reqAPI.Header.Set("Authorization", "OAuth "+c.token)

	respAPI, err := c.httpClient.Do(reqAPI)
	if err != nil {
		return err
	}
	defer respAPI.Body.Close()

	if respAPI.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(respAPI.Body)
		return fmt.Errorf("yandex API returned status %d: %s", respAPI.StatusCode, string(body))
	}

	var linkResp uploadLinkResponse
	if err := json.NewDecoder(respAPI.Body).Decode(&linkResp); err != nil {
		return err
	}

	progressReader := progress.NewReader(file, fileSize)
	reqUpload, err := http.NewRequest("PUT", linkResp.Href, progressReader)
	if err != nil {
		return err
	}
	reqUpload.ContentLength = fileSize
	reqUpload.Header.Set("Content-Type", "application/octet-stream")

	respUpload, err := c.httpClient.Do(reqUpload)
	if err != nil {
		return err
	}
	defer respUpload.Body.Close()

	if respUpload.StatusCode != http.StatusCreated && respUpload.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(respUpload.Body)
		return fmt.Errorf("storage upload failed with status %d: %s", respUpload.StatusCode, string(body))
	}

	fmt.Printf("\r   Uploading: 100.00%% | [Transfer Successfully Completed]                  \n")
	return nil
}

// PublishAndGetLink executes cloud asset publishing pipelines and fetches public delivery URLs.
func (c *Client) PublishAndGetLink(remotePath string) (string, error) {
	time.Sleep(350 * time.Millisecond) // Mitigation delay sequence preventing API throttling limit encounters.

	publishURL := restEndpoint + "/publish?path=" + url.QueryEscape(remotePath)
	reqPub, err := http.NewRequest("PUT", publishURL, nil)
	if err != nil {
		return "", err
	}
	reqPub.Header.Set("Authorization", "OAuth "+c.token)
	respPub, err := c.httpClient.Do(reqPub)
	if err != nil {
		return "", err
	}
	respPub.Body.Close()

	metaURL := restEndpoint + "?path=" + url.QueryEscape(remotePath) + "&fields=public_url"
	reqMeta, err := http.NewRequest("GET", metaURL, nil)
	if err != nil {
		return "", err
	}
	reqMeta.Header.Set("Authorization", "OAuth "+c.token)
	respMeta, err := c.httpClient.Do(reqMeta)
	if err != nil {
		return "", err
	}
	defer respMeta.Body.Close()

	var metaResult struct {
		PublicURL string `json:"public_url"`
	}
	_ = json.NewDecoder(respMeta.Body).Decode(&metaResult)
	return metaResult.PublicURL, nil
}

// createRemoteDir sets up directory infrastructure remotely if it hasn't been explicitly declared yet.
func (c *Client) createRemoteDir(remotePath string) error {
	apiURL := restEndpoint + "?path=" + url.QueryEscape(remotePath)
	req, err := http.NewRequest("PUT", apiURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "OAuth "+c.token)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}
