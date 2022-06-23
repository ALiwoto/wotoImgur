package wotoImgur

import (
	"net/http"
	"time"
)

// Client used to for go-imgur
type ImgurClient struct {
	HTTPClient    *http.Client
	ImgurClientID string
	RapidAPIKey   string
}

type ImgurError struct {
	// Message is the error message.
	Message string

	// Err is the inner error.
	Err error

	// Status is the status of the error.
	Status int
}

type albumInfoDataWrapper struct {
	Ai      *AlbumInfo `json:"data"`
	Success bool       `json:"success"`
	Status  int        `json:"status"`
}

// AlbumInfo contains all album information provided by imgur
type AlbumInfo struct {
	ID          string      `json:"id"`                   // The ID for the album
	Title       string      `json:"title"`                // The title of the album in the gallery
	Description string      `json:"description"`          // The description of the album in the gallery
	DateTime    int         `json:"datetime"`             // Time inserted into the gallery, epoch time
	Cover       string      `json:"cover"`                // The ID of the album cover image
	CoverWidth  int         `json:"cover_width"`          // The width, in pixels, of the album cover image
	CoverHeight int         `json:"cover_height"`         // The height, in pixels, of the album cover image
	AccountURL  string      `json:"account_url"`          // The account username or null if it's anonymous.
	AccountID   int         `json:"account_id"`           // The account ID or null if it's anonymous.
	Privacy     string      `json:"privacy"`              // The privacy level of the album, you can only view public if not logged in as album owner
	Layout      string      `json:"layout"`               // The view layout of the album.
	Views       int         `json:"views"`                // The number of album views
	Link        string      `json:"link"`                 // The URL link to the album
	Favorite    bool        `json:"favorite"`             // Indicates if the current user favorited the image. Defaults to false if not signed in.
	Nsfw        bool        `json:"nsfw"`                 // Indicates if the image has been marked as nsfw or not. Defaults to null if information is not available.
	Section     string      `json:"section"`              // If the image has been categorized by our backend then this will contain the section the image belongs in. (funny, cats, adviceanimals, wtf, etc)
	Order       int         `json:"order"`                // Order number of the album on the user's album page (defaults to 0 if their albums haven't been reordered)
	DeleteHash  string      `json:"deletehash,omitempty"` // OPTIONAL, the deletehash, if you're logged in as the album owner
	ImagesCount int         `json:"images_count"`         // The total number of images in the album
	Images      []ImageInfo `json:"images"`               // An array of all the images in the album (only available when requesting the direct album)
	InGallery   bool        `json:"in_gallery"`           // True if the image has been submitted to the gallery, false if otherwise.
	Limit       *RateLimit  // Current rate limit
}

// Comment is an imgur comment
type Comment struct {
	ID         int       `json:"id"`          // The ID for the comment
	ImageID    string    `json:"image_id"`    //The ID of the image that the comment is for
	Comment    string    `json:"comment"`     // The comment itself.
	Author     string    `json:"author"`      // Username of the author of the comment
	AuthorID   int       `json:"author_id"`   // The account ID for the author
	OnAlbum    bool      `json:"on_album"`    // If this comment was done to an album
	AlbumCover string    `json:"album_cover"` // The ID of the album cover image, this is what should be displayed for album comments
	Ups        int       `json:"ups"`         //	Number of upvotes for the comment
	Downs      int       `json:"downs"`       // The number of downvotes for the comment
	Points     float32   `json:"points"`      // the number of upvotes - downvotes
	Datetime   int       `json:"datetime"`    // Timestamp of creation, epoch time
	ParentID   int       `json:"parent_id"`   // If this is a reply, this will be the value of the comment_id for the caption this a reply for.
	Deleted    bool      `json:"deleted"`     // Marked true if this caption has been deleted
	Vote       string    `json:"vote"`        // The current user's vote on the comment. null if not signed in or if the user hasn't voted on it.
	Children   []Comment `json:"children"`    // All of the replies for this comment. If there are no replies to the comment then this is an empty set.
}

// GenericInfo is returned from functions for which the final result type is not known beforehand.
// Only one pointer is != nil
type GenericInfo struct {
	Image  *ImageInfo
	Album  *AlbumInfo
	GImage *GalleryImageInfo
	GAlbum *GalleryAlbumInfo
	Limit  *RateLimit
}

type galleryAlbumInfoDataWrapper struct {
	Ai      *GalleryAlbumInfo `json:"data"`
	Success bool              `json:"success"`
	Status  int               `json:"status"`
}

// GalleryAlbumInfo contains all information provided by imgur of a gallery album
type GalleryAlbumInfo struct {
	ID           string      `json:"id"`               // The ID for the album
	Title        string      `json:"title"`            // The title of the album in the gallery
	Description  string      `json:"description"`      // The description of the album in the gallery
	DateTime     int         `json:"datetime"`         // Time inserted into the gallery, epoch time
	Cover        string      `json:"cover"`            // The ID of the album cover image
	CoverWidth   int         `json:"cover_width"`      // The width, in pixels, of the album cover image
	CoverHeight  int         `json:"cover_height"`     // The height, in pixels, of the album cover image
	AccountURL   string      `json:"account_url"`      // The account username or null if it's anonymous.
	AccountID    int         `json:"account_id"`       // The account ID or null if it's anonymous.
	Privacy      string      `json:"privacy"`          // The privacy level of the album, you can only view public if not logged in as album owner
	Layout       string      `json:"layout"`           // The view layout of the album.
	Views        int         `json:"views"`            // The number of album views
	Link         string      `json:"link"`             // The URL link to the album
	Ups          int         `json:"ups"`              // Upvotes for the image
	Downs        int         `json:"downs"`            // Number of downvotes for the image
	Points       int         `json:"points"`           // Upvotes minus downvotes
	Score        int         `json:"score"`            // Imgur popularity score
	IsAlbum      bool        `json:"is_album"`         // if it's an album or not
	Vote         string      `json:"vote"`             // The current user's vote on the album. null if not signed in or if the user hasn't voted on it.
	Favorite     bool        `json:"favorite"`         // Indicates if the current user favorited the image. Defaults to false if not signed in.
	Nsfw         bool        `json:"nsfw"`             // Indicates if the image has been marked as nsfw or not. Defaults to null if information is not available.
	CommentCount int         `json:"comment_count"`    // Number of comments on the gallery album.
	Topic        string      `json:"topic"`            // Topic of the gallery album.
	TopicID      int         `json:"topic_id"`         // Topic ID of the gallery album.
	ImagesCount  int         `json:"images_count"`     // The total number of images in the album
	Images       []ImageInfo `json:"images,omitempty"` // An array of all the images in the album (only available when requesting the direct album)
	InMostViral  bool        `json:"in_most_viral"`    // Indicates if the album is in the most viral gallery or not.
	Limit        *RateLimit  // Current rate limit
}

type galleryImageInfoDataWrapper struct {
	Ii      *GalleryImageInfo `json:"data"`
	Success bool              `json:"success"`
	Status  int               `json:"status"`
}

// GalleryImageInfo contains all gallery image information provided by imgur
type GalleryImageInfo struct {
	ID           string     `json:"id"`                   // The ID for the image
	Title        string     `json:"title"`                // The title of the image.
	Description  string     `json:"description"`          // Description of the image.
	Datetime     int        `json:"datetime"`             // Time uploaded, epoch time
	MimeType     string     `json:"type"`                 // Image MIME type.
	Animated     bool       `json:"animated"`             // is the image animated
	Width        int        `json:"width"`                // The width of the image in pixels
	Height       int        `json:"height"`               // The height of the image in pixels
	Size         int        `json:"size"`                 // The size of the image in bytes
	Views        int        `json:"views"`                // The number of image views
	Bandwidth    int        `json:"bandwidth"`            // Bandwidth consumed by the image in bytes
	DeleteHash   string     `json:"deletehash,omitempty"` // OPTIONAL, the deletehash, if you're logged in as the image owner
	Link         string     `json:"link"`                 // The direct link to the the image. (Note: if fetching an animated GIF that was over 20MB in original size, a .gif thumbnail will be returned)
	Gifv         string     `json:"gifv,omitempty"`       // OPTIONAL, The .gifv link. Only available if the image is animated and type is 'image/gif'.
	Mp4          string     `json:"mp4,omitempty"`        // OPTIONAL, The direct link to the .mp4. Only available if the image is animated and type is 'image/gif'.
	Mp4Size      int        `json:"mp4_size,omitempty"`   // OPTIONAL, The Content-Length of the .mp4. Only available if the image is animated and type is 'image/gif'. Note that a zero value (0) is possible if the video has not yet been generated
	Looping      bool       `json:"looping,omitempty"`    // OPTIONAL, Whether the image has a looping animation. Only available if the image is animated and type is 'image/gif'.
	Vote         string     `json:"vote"`                 // The current user's vote on the album. null if not signed in or if the user hasn't voted on it.
	Favorite     bool       `json:"favorite"`             // Indicates if the current user favorited the image. Defaults to false if not signed in.
	Nsfw         bool       `json:"nsfw"`                 // Indicates if the image has been marked as nsfw or not. Defaults to null if information is not available.
	CommentCount int        `json:"comment_count"`        // Number of comments on the gallery album.
	Topic        string     `json:"topic"`                // Topic of the gallery album.
	TopicID      int        `json:"topic_id"`             // Topic ID of the gallery album.
	Section      string     `json:"section"`              // If the image has been categorized by our backend then this will contain the section the image belongs in. (funny, cats, adviceanimals, wtf, etc)
	AccountURL   string     `json:"account_url"`          // The username of the account that uploaded it, or null.
	AccountID    int        `json:"account_id"`           // The account ID of the account that uploaded it, or null.
	Ups          int        `json:"ups"`                  // Upvotes for the image
	Downs        int        `json:"downs"`                // Number of downvotes for the image
	Points       int        `json:"points"`               // Upvotes minus downvotes
	Score        int        `json:"score"`                // Imgur popularity score
	IsAlbum      bool       `json:"is_album"`             // if it's an album or not
	InMostViral  bool       `json:"in_most_viral"`        // Indicates if the album is in the most viral gallery or not.
	Limit        *RateLimit // Current rate limit
}

type imageInfoDataWrapper struct {
	Ii      *ImageInfo `json:"data"`
	Success bool       `json:"success"`
	Status  int        `json:"status"`
}

// ImageInfo contains all image information provided by imgur
type ImageInfo struct {
	ID          string     `json:"id"`                   // The ID for the image
	Title       string     `json:"title"`                // The title of the image.
	Description string     `json:"description"`          // Description of the image.
	Datetime    int        `json:"datetime"`             // Time uploaded, epoch time
	MimeType    string     `json:"type"`                 // Image MIME type.
	Animated    bool       `json:"animated"`             // is the image animated
	Width       int        `json:"width"`                // The width of the image in pixels
	Height      int        `json:"height"`               // The height of the image in pixels
	Size        int        `json:"size"`                 // The size of the image in bytes
	Views       int        `json:"views"`                // The number of image views
	Bandwidth   int        `json:"bandwidth"`            // Bandwidth consumed by the image in bytes
	DeleteHash  string     `json:"deletehash,omitempty"` // OPTIONAL, the deletehash, if you're logged in as the image owner
	Name        string     `json:"name,omitempty"`       // OPTIONAL, the original filename, if you're logged in as the image owner
	Section     string     `json:"section"`              // If the image has been categorized by our backend then this will contain the section the image belongs in. (funny, cats, adviceanimals, wtf, etc)
	Link        string     `json:"link"`                 // The direct link to the the image. (Note: if fetching an animated GIF that was over 20MB in original size, a .gif thumbnail will be returned)
	Gifv        string     `json:"gifv,omitempty"`       // OPTIONAL, The .gifv link. Only available if the image is animated and type is 'image/gif'.
	Mp4         string     `json:"mp4,omitempty"`        // OPTIONAL, The direct link to the .mp4. Only available if the image is animated and type is 'image/gif'.
	Mp4Size     int        `json:"mp4_size,omitempty"`   // OPTIONAL, The Content-Length of the .mp4. Only available if the image is animated and type is 'image/gif'. Note that a zero value (0) is possible if the video has not yet been generated
	Looping     bool       `json:"looping,omitempty"`    // OPTIONAL, Whether the image has a looping animation. Only available if the image is animated and type is 'image/gif'.
	Favorite    bool       `json:"favorite"`             // Indicates if the current user favorited the image. Defaults to false if not signed in.
	Nsfw        bool       `json:"nsfw"`                 // Indicates if the image has been marked as nsfw or not. Defaults to null if information is not available.
	Vote        string     `json:"vote"`                 // The current user's vote on the album. null if not signed in, if the user hasn't voted on it, or if not submitted to the gallery.
	InGallery   bool       `json:"in_gallery"`           // True if the image has been submitted to the gallery, false if otherwise.
	Limit       *RateLimit // Current rate limit
}

type rateLimitDataWrapper struct {
	Rl      *rateLimitInternal `json:"data"`
	Success bool               `json:"success"`
	Status  int                `json:"status"`
}

// internal representation used for the json parser
type rateLimitInternal struct {
	UserLimit       int64
	UserRemaining   int64
	UserReset       int64
	ClientLimit     int64
	ClientRemaining int64
}

// RateLimit details can be found here: https://api.imgur.com/#limits
type RateLimit struct {
	// Total credits that can be allocated.
	UserLimit int64
	// Total credits available.
	UserRemaining int64
	// Timestamp for when the credits will be reset.
	UserReset time.Time
	// Total credits that can be allocated for the application in a day.
	ClientLimit int64
	// Total credits remaining for the application in a day.
	ClientRemaining int64
}
