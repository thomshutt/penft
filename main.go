package main

func main() {
	//	visibleImageBytes, err := ioutil.ReadFile(imageFile)
	//	if err != nil {
	//		log.Fatalf("%v", err)
	//	}
	//
	//	realImageBytes, err := ioutil.ReadFile(imageFile)
	//	if err != nil {
	//		log.Fatalf("%v", err)
	//	}
	//
	//	visibleImage, err := alterVisibleImage(visibleImageBytes)
	//	if err != nil {
	//		log.Fatalf("%v", err)
	//	}
	//
	//	combinedImage, err := writeDataToDNGTag(visibleImage, realImageBytes)
	//	if err != nil {
	//		log.Fatalf("%v", err)
	//	}
	//
	//	err = ioutil.WriteFile("public/images/encrypted-image.jpg", combinedImage, fs.ModePerm)
	//	if err != nil {
	//		log.Fatalf("%v", err)
	//	}
	//
	//	ipfsUploadResponse, err := uploadToIPFS(combinedImage)
	//	if err != nil {
	//		log.Fatalf("%v", err)
	//	}
	//
	//	println("IPFS:", ipfsUploadResponse.IPFSHash)
	//	//
	//	//mintNFTResponse, err := mintNFT(ipfsUploadResponse.IPFSHash)
	//	//if err != nil {
	//	//	log.Fatalf("%v", err)
	//	//}
	//	//
	//	//println("mintNFTResponse.TxId:", mintNFTResponse.TxId, "mintNFTResponse.failed", mintNFTResponse.Failed)
}
