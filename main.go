package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/dsoprea/go-exif/v3"
	jpegstructure "github.com/dsoprea/go-jpeg-image-structure/v2"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
)

const TATUM_API_KEY = "TODO"
const imageFile = "public/images/ape.jpg"

var iv = []byte{'\x0f', '\x0f', '\x0f', '\x0f', '\x0f', '\x0f', '\x0f', '\x0f', '\x0f', '\x0f', '\x0f', '\x0f', '\x0f', '\x0f', '\x0f', '\x0f'}
var key = []byte("qwertyuioplkjhgfdsazxcvbnmqwerty")

type IPFSUploadResponse struct {
	IPFSHash string `json:"ipfsHash"`
}

type MintNFTRequest struct {
	Chain string `json:"chain"`
	To    string `json:"to"`
	Url   string `json:"url"`
}

type MintNFTResponse struct {
	TxId   string `json:"txId"`
	Failed bool   `json:"failed"`
}

func main() {
	visibleImageBytes, err := ioutil.ReadFile(imageFile)
	if err != nil {
		log.Fatalf("%v", err)
	}

	realImageBytes, err := ioutil.ReadFile(imageFile)
	if err != nil {
		log.Fatalf("%v", err)
	}

	visibleImage, err := alterVisibleImage(visibleImageBytes)
	if err != nil {
		log.Fatalf("%v", err)
	}

	combinedImage, err := writeDataToDNGTag(visibleImage, realImageBytes)
	if err != nil {
		log.Fatalf("%v", err)
	}

	err = ioutil.WriteFile("public/images/encrypted-image.jpg", combinedImage, fs.ModePerm)
	if err != nil {
		log.Fatalf("%v", err)
	}

	ipfsUploadResponse, err := uploadToIPFS(combinedImage)
	if err != nil {
		log.Fatalf("%v", err)
	}

	println("IPFS:", ipfsUploadResponse.IPFSHash)
	//
	//mintNFTResponse, err := mintNFT(ipfsUploadResponse.IPFSHash)
	//if err != nil {
	//	log.Fatalf("%v", err)
	//}
	//
	//println("mintNFTResponse.TxId:", mintNFTResponse.TxId, "mintNFTResponse.failed", mintNFTResponse.Failed)
}

func alterVisibleImage(imageBytes []byte) ([]byte, error) {
	src, _, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		return nil, err
	}

	rgba, ok := src.(*image.RGBA)
	if !ok {
		b := src.Bounds()
		rgba = image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
		draw.Draw(rgba, rgba.Bounds(), src, b.Min, draw.Src)
	}

	min := rgba.Bounds().Min
	max := rgba.Bounds().Max

	for x := min.X; x < max.X; x++ {
		for y := min.Y; y < max.Y; y++ {
			currentColor := rgba.RGBA64At(x, y)
			if x%2 == 0 || y%2 == 0 {
				rgba.SetRGBA64(x, y, color.RGBA64{
					R: 0,
					G: 0,
					B: 0,
					A: currentColor.A,
				})
			}

		}
	}

	var buf bytes.Buffer
	err = jpeg.Encode(&buf, rgba, nil)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func writeDataToDNGTag(visibleImageBytes []byte, data []byte) ([]byte, error) {
	jmp := jpegstructure.NewJpegMediaParser()

	intfc, err := jmp.ParseBytes(visibleImageBytes)
	if err != nil {
		return nil, err
	}

	sl := intfc.(*jpegstructure.SegmentList)

	rootIb, err := sl.ConstructExifBuilder()
	if err != nil {
		return nil, err
	}

	ifdIb, err := exif.GetOrCreateIbFromRootIb(rootIb, "IFD0")
	if err != nil {
		return nil, err
	}

	s, err := encryptToBase64([]byte("message message message"))
	if err != nil {
		return nil, err
	}
	err = ifdIb.SetStandardWithName("DateTime", s)
	if err != nil {
		return nil, err
	}

	base64Image := base64.StdEncoding.EncodeToString(data)
	encryptedImage, err := encryptToBase64([]byte(base64Image))
	if err != nil {
		return nil, err
	}

	err = ifdIb.SetStandardWithName("ImageHistory", encryptedImage)
	if err != nil {
		return nil, err
	}

	// Update the exif segment.
	err = sl.SetExif(rootIb)
	if err != nil {
		return nil, err
	}

	b := new(bytes.Buffer)
	err = sl.Write(b)
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func encryptToBase64(toEncrypt []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, len(toEncrypt))
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertext, toEncrypt)

	// CTR mode is the same for both encryption and decryption, so we can
	// also decrypt that ciphertext with NewCTR.
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func mintNFT(ipfsID string) (MintNFTResponse, error) {
	mintNFTRequest := MintNFTRequest{
		Chain: "ETH",
		To:    "0x53e8577C4347C365E4e0DA5B57A589cB6f2AB848",
		Url:   "ipfs://" + ipfsID,
	}

	jsonBytes, err := json.Marshal(mintNFTRequest)
	if err != nil {
		return MintNFTResponse{}, err
	}

	req, err := http.NewRequest("POST", "https://api-eu1.tatum.io/v3/nft/mint", bytes.NewReader(jsonBytes))
	if err != nil {
		return MintNFTResponse{}, err
	}

	req.Header.Add("content-type", "application/json")
	req.Header.Add("x-testnet-type", "ethereum-ropsten")
	req.Header.Add("x-api-key", TATUM_API_KEY)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return MintNFTResponse{}, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return MintNFTResponse{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return MintNFTResponse{}, fmt.Errorf("Got error code while minting NFT: %d. Resp body: %s", resp.StatusCode, respBody)
	}

	var mintNFTResponse MintNFTResponse
	if err := json.Unmarshal(respBody, &mintNFTResponse); err != nil {
		return MintNFTResponse{}, err
	}

	return mintNFTResponse, nil
}

func uploadToIPFS(file []byte) (IPFSUploadResponse, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", "file.jpg")
	if err != nil {
		return IPFSUploadResponse{}, err
	}

	_, err = io.Copy(part, bytes.NewReader(file))
	if err != nil {
		return IPFSUploadResponse{}, err
	}

	err = writer.Close()
	if err != nil {
		return IPFSUploadResponse{}, err
	}

	req, err := http.NewRequest(http.MethodPost, "https://api-eu1.tatum.io/v3/ipfs", body)
	if err != nil {
		return IPFSUploadResponse{}, err
	}
	req.Header.Set("x-api-key", TATUM_API_KEY)
	req.Header.Add("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return IPFSUploadResponse{}, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return IPFSUploadResponse{}, err
	}
	if resp.StatusCode != http.StatusCreated {
		return IPFSUploadResponse{}, fmt.Errorf("Got error code while uploading to IPFS: %d. Body: %s", resp.StatusCode, respBody)
	}

	var ipfsResponse IPFSUploadResponse
	if err := json.Unmarshal(respBody, &ipfsResponse); err != nil {
		return IPFSUploadResponse{}, err
	}

	return ipfsResponse, nil
}
