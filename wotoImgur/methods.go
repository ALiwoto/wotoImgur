package wotoImgur

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/AnimeKaizoku/ssg/ssg"
)

// GetAlbumInfo queries imgur for information on a album
// returns album info, status code of the request, error
func (client *ImgurClient) GetAlbumInfo(id string) (*AlbumInfo, error) {
	body, rl, err := client.getURL("album/" + id)
	if err != nil {
		return nil, getErr(-1, "Problem getting URL for album info ID "+id+" - "+err.Error())
	}
	//client.Log.Debugf("%v\n", body)

	dec := json.NewDecoder(strings.NewReader(body))
	var alb albumInfoDataWrapper
	if err := dec.Decode(&alb); err != nil {
		return nil, getErr(-1, "Problem decoding json for albumID "+id+" - "+err.Error())
	}

	if !alb.Success {
		return nil, getErr(alb.Status, "Request to imgur failed for albumID "+id+" - "+strconv.Itoa(alb.Status))
	}

	alb.Ai.Limit = rl
	return alb.Ai, nil
}

// GetInfoFromURL tries to query imgur based on information identified in the URL.
// returns image/album info, status code of the request, error
func (client *ImgurClient) GetInfoFromURL(url string) (*GenericInfo, error) {
	url = strings.TrimSpace(url)

	// https://i.imgur.com/<id>.jpg -> image
	if strings.Contains(url, "://i.imgur.com/") {
		return client.directImageURL(url)
	}

	// https://imgur.com/a/<id> -> album
	if strings.Contains(url, "://imgur.com/a/") || strings.Contains(url, "://m.imgur.com/a/") {
		return client.albumURL(url)
	}

	// https://imgur.com/gallery/<id> -> gallery album
	if strings.Contains(url, "://imgur.com/gallery/") || strings.Contains(url, "://m.imgur.com/gallery/") {
		return client.galleryURL(url)
	}

	// https://imgur.com/<id> -> image
	if strings.Contains(url, "://imgur.com/") || strings.Contains(url, "://m.imgur.com/") {
		return client.imageURL(url)
	}

	return nil, getErr(-1, "URL pattern matching for URL "+url+" failed.")
}

func (client *ImgurClient) directImageURL(url string) (*GenericInfo, error) {
	var ret GenericInfo
	start := strings.LastIndex(url, "/") + 1
	end := strings.LastIndex(url, ".")
	if start+1 >= end {
		return nil, getErr(-1, "Could not find ID in URL "+url+". I was going down i.imgur.com path.")
	}
	id := url[start:end]
	// client.Log.Debugf("Detected imgur image ID %v. Was going down the i.imgur.com/ path.", id)
	gii, err := client.GetGalleryImageInfo(id)
	if err == nil {
		ret.GImage = gii
	} else {
		var ii *ImageInfo
		ii, err = client.GetImageInfo(id)
		ret.Image = ii
	}
	return &ret, err
}

func (client *ImgurClient) albumURL(url string) (*GenericInfo, error) {
	var ret GenericInfo

	start := strings.LastIndex(url, "/") + 1
	end := strings.LastIndex(url, "?")
	if end == -1 {
		end = len(url)
	}
	id := url[start:end]
	if id == "" {
		return nil, getErr(-1, "Could not find ID in URL "+url+". I was going down imgur.com/a/ path.")
	}
	// client.Log.Debugf("Detected imgur album ID %v. Was going down the imgur.com/a/ path.", id)
	ai, err := client.GetAlbumInfo(id)
	ret.Album = ai
	return &ret, err
}

func (client *ImgurClient) galleryURL(url string) (*GenericInfo, error) {
	var ret GenericInfo

	start := strings.LastIndex(url, "/") + 1
	end := strings.LastIndex(url, "?")
	if end == -1 {
		end = len(url)
	}
	id := url[start:end]
	if id == "" {
		return nil, getErr(-1, "Could not find ID in URL "+url+". I was going down imgur.com/gallery/ path.")
	}

	// client.Log.Debugf("Detected imgur gallery ID %v. Was going down the imgur.com/gallery/ path.", id)
	ai, err := client.GetGalleryAlbumInfo(id)
	if err == nil {
		ret.GAlbum = ai
		return &ret, err
	}
	// fallback to GetGalleryImageInfo
	// client.Log.Debugf("Failed to retrieve imgur gallery album. Attempting to retrieve imgur gallery image. err: %v status: %d", err, status)
	ii, err := client.GetGalleryImageInfo(id)
	ret.GImage = ii
	return &ret, err
}

func (client *ImgurClient) imageURL(url string) (*GenericInfo, error) {
	var ret GenericInfo

	start := strings.LastIndex(url, "/") + 1
	end := strings.LastIndex(url, "?")
	if end == -1 {
		end = len(url)
	}
	id := url[start:end]
	if id == "" {
		return nil, getErr(-1, "Could not find ID in URL "+url+". I was going down imgur.com/ path.")
	}
	// client.Log.Debugf("Detected imgur image ID %v. Was going down the imgur.com/ path.", id)
	ii, err := client.GetGalleryImageInfo(id)
	if err == nil {
		ret.GImage = ii
		return &ret, nil
	}

	i, err := client.GetImageInfo(id)
	ret.Image = i
	return &ret, err
}

// GetGalleryAlbumInfo queries imgur for information on a gallery album
// returns album info, status code of the request, error
func (client *ImgurClient) GetGalleryAlbumInfo(id string) (*GalleryAlbumInfo, error) {
	body, rl, err := client.getURL("gallery/album/" + id)
	if err != nil {
		return nil, getErr(-1, "Problem getting URL for gallery album info ID "+id+" - "+err.Error())
	}
	// client.Log.Debugf("%v\n", body)

	dec := json.NewDecoder(strings.NewReader(body))
	var alb galleryAlbumInfoDataWrapper
	if err := dec.Decode(&alb); err != nil {
		return nil, getErr(-1, "Problem decoding json for gallery albumID "+id+" - "+err.Error())
	}
	alb.Ai.Limit = rl

	if !alb.Success {
		return nil, getErr(alb.Status, "Request to imgur failed for gallery albumID "+id+" - "+strconv.Itoa(alb.Status))
	}
	return alb.Ai, nil
}

// GetGalleryImageInfo queries imgur for information on a image
// returns image info, status code of the request, error
func (client *ImgurClient) GetGalleryImageInfo(id string) (*GalleryImageInfo, error) {
	body, rl, err := client.getURL("gallery/image/" + id)
	if err != nil {
		return nil, getErr(-1, "Problem getting URL for gallery image info ID "+id+" - "+err.Error())
	}
	// client.Log.Debugf("%v\n", body)

	dec := json.NewDecoder(strings.NewReader(body))
	var img galleryImageInfoDataWrapper
	if err := dec.Decode(&img); err != nil {
		return nil, getErr(-1, "Problem decoding json for gallery imageID "+id+" - "+err.Error())
	}
	img.Ii.Limit = rl

	if !img.Success {
		return nil, getErr(img.Status, "Request to imgur failed for gallery imageID "+id+" - "+strconv.Itoa(img.Status))
	}
	return img.Ii, nil
}

func (client *ImgurClient) createAPIURL(u string) string {
	if client.RapidAPIKey == "" {
		return apiEndpoint + u
	}
	return apiEndpointRapidAPI + u
}

// getURL returns
// - body as string
// - RateLimit with current limits
// - error in case something broke
func (client *ImgurClient) getURL(URL string) (string, *RateLimit, error) {
	URL = client.createAPIURL(URL)
	// client.Log.Infof("Requesting URL %v\n", URL)
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return "", nil, errors.New("Could not create request for " + URL + " - " + err.Error())
	}

	req.Header.Add("Authorization", "Client-ID "+client.ImgurClientID)
	if client.RapidAPIKey != "" {
		req.Header.Add("x-rapidapi-host", "imgur-apiv3.p.rapidapi.com")
		req.Header.Add("x-rapidapi-key", client.RapidAPIKey)
	}

	// Make a request to the sourceURL
	res, err := client.HTTPClient.Do(req)
	if err != nil {
		return "", nil, errors.New("Could not get " + URL + " - " + err.Error())
	}
	defer res.Body.Close()

	if !(res.StatusCode >= 200 && res.StatusCode <= 300) {
		return "", nil, errors.New("HTTP status indicates an error for " + URL + " - " + res.Status)
	}

	// Read the whole body
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", nil, errors.New("Problem reading the body for " + URL + " - " + err.Error())
	}

	// Get RateLimit headers
	rl, err := extractRateLimits(res.Header)
	if err != nil {
		// client.Log.Infof("Problem with extracting rate limits: %v", err)
	}

	return string(body[:]), rl, nil
}

// GetImageInfo queries imgur for information on a image
// returns image info, status code of the request, error
func (client *ImgurClient) GetImageInfo(id string) (*ImageInfo, error) {
	body, rl, err := client.getURL("image/" + id)
	if err != nil {
		return nil, getErr(-1, "Problem getting URL for image info ID "+id+" - "+err.Error())
	}

	dec := json.NewDecoder(strings.NewReader(body))
	var img imageInfoDataWrapper
	if err := dec.Decode(&img); err != nil {
		return nil, getErr(-1, "Problem decoding json for imageID "+id+" - "+err.Error())
	}
	img.Ii.Limit = rl

	if !img.Success {
		return nil, getErr(img.Status, "Request to imgur failed for imageID "+id+" - "+strconv.Itoa(img.Status))
	}
	return img.Ii, nil
}

// GetRateLimit returns the current rate limit without doing anything else
func (client *ImgurClient) GetRateLimit() (*RateLimit, error) {
	// We are requesting any URL and parse the returned HTTP headers
	body, rl, err := client.getURL("account/kaffeeshare")

	if err != nil {
		return nil, errors.New("Problem getting URL for rate - " + err.Error())
	}
	//client.Log.Debugf("%v\n", body)

	dec := json.NewDecoder(strings.NewReader(body))

	var bodyDecoded rateLimitDataWrapper
	if err := dec.Decode(&bodyDecoded); err != nil {
		return nil, errors.New("Problem decoding json for ratelimit - " + err.Error())
	}

	if !bodyDecoded.Success {
		return nil, errors.New("Request to imgur failed for ratelimit - " + strconv.Itoa(bodyDecoded.Status))
	}

	var ret RateLimit
	ret.ClientLimit = rl.ClientLimit
	ret.ClientRemaining = rl.ClientRemaining
	ret.UserLimit = rl.UserLimit
	ret.UserRemaining = rl.UserRemaining
	ret.UserReset = rl.UserReset

	return &ret, nil
}

// UploadImage uploads the image to imgur
// image                Can be a binary file, base64 data, or a URL for an image. (up to 10MB)
// album       optional The id of the album you want to add the image to.
//                      For anonymous albums, album should be the deleteHash that is returned at creation.
// dType                The type of the file that's being sent; file, base64 or URL
// title       optional The title of the image.
// description optional The description of the image.
// returns image info, status code of the upload, error
func (client *ImgurClient) UploadImage(image []byte, album, dType, title, description string) (*ImageInfo, error) {
	if image == nil {
		return nil, getErr(-1, "Invalid image")
	}
	if dType != "file" && dType != "base64" && dType != "URL" {
		return nil, getErr(-1, "Passed invalid dType: "+dType+". Please use file/base64/URL.")
	}

	form := createUploadForm(image, album, dType, title, description)

	URL := client.createAPIURL("image")
	req, err := http.NewRequest("POST", URL, bytes.NewBufferString(form.Encode()))
	// client.Log.Debugf("Posting to URL %v\n", URL)
	if err != nil {
		return nil, getErr(-1, "Could create request for "+URL+" - "+err.Error())
	}

	req.Header.Add("Authorization", "Client-ID "+client.ImgurClientID)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if client.RapidAPIKey != "" {
		req.Header.Add("X-RapidAPI-Key", client.RapidAPIKey)
	}

	res, err := client.HTTPClient.Do(req)
	if err != nil {
		return nil, getErr(-1, "Could not post "+URL+" - "+err.Error())
	}
	defer res.Body.Close()

	// Read the whole body
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, getErr(-1, "Problem reading the body of "+URL+" - "+err.Error())
	}

	// client.Log.Debugf("%v\n", string(body[:]))

	dec := json.NewDecoder(bytes.NewReader(body))
	var img imageInfoDataWrapper
	if err = dec.Decode(&img); err != nil {
		return nil, getErr(-1, "Problem decoding json result from image upload - "+err.Error()+". JSON(?): "+string(body))
	}

	if !img.Success {
		return nil, getErr(img.Status, "Upload to imgur failed with status: "+strconv.Itoa(img.Status))
	}

	img.Ii.Limit, _ = extractRateLimits(res.Header)

	return img.Ii, nil
}

// UploadImageFromFile uploads a file given by the filename string to imgur.
func (client *ImgurClient) UploadImageFromFile(filename, album, title, description string) (*ImageInfo, error) {
	// client.Log.Infof("*** IMAGE UPLOAD ***\n")
	f, err := os.Open(filename)
	if err != nil {
		return nil, getErrF(500, "Could not open file %v - Error: %v", filename, err)
	}
	defer f.Close()
	fileInfo, err := f.Stat()
	if err != nil {
		return nil, getErrF(500, "Could not stat file %v - Error: %v", filename, err)
	}
	size := fileInfo.Size()
	b := make([]byte, size)
	n, err := f.Read(b)
	if err != nil || int64(n) != size {
		return nil, getErrF(500, "Could not read file %v - Error: %v", filename, err)
	}

	return client.UploadImage(b, album, "file", title, description)
}

// --------------------------------------------------------

func (e *ImgurError) Error() string {
	myStr := ""
	if e.Status != 0 {
		myStr += "[" + ssg.ToBase10(int64(e.Status)) + "] : "
	}

	if e.Message != "" {
		myStr += e.Message + " : "
	}

	if e.Err != nil {
		myStr += e.Err.Error()
	}

	myStr = strings.Trim(strings.TrimSpace(myStr), ":")

	return myStr
}
