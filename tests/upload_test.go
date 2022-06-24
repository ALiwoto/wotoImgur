package tests

import (
	"testing"

	"github.com/ALiwoto/wotoImgur/wotoImgur"
)

func TestUploadFile(t *testing.T) {
	client, err := wotoImgur.NewImgurClient("", nil)
	if err != nil {
		t.Error("when tried to get new client: ", err.Error())
		return
	}

	info, err := client.UploadImageFromFile("temp.png", "", "file title", "file description")
	if err != nil {
		t.Error("when tried to upload new photo: ", err.Error())
		return
	}

	if info.Link == "" {
		t.Error("link of the uploaded image is empty")
		return
	}
}
