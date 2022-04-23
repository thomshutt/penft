package main

func createResponse(ipfsHash string) string {
	return `
	<!DOCTYPE html>
<html lang="en" data-theme="light">
<head>

    <!-- Basic Page Needs
    –––––––––––––––––––––––––––––––––––––––––––––––––– -->
    <meta charset="utf-8">
    <title>peNFT</title>
    <meta name="description" content="">
    <meta name="author" content="">

    <!-- CSS
    –––––––––––––––––––––––––––––––––––––––––––––––––– -->
    <link rel="stylesheet" href="/css/pico.min.css">

</head>
<body>

<header class="container" style="padding-bottom: 0px;">
    <hgroup style="text-align: center;">
        <img src="/images/banner.png" style="max-width:75%;"/>
    </hgroup>
    <hgroup style="text-align:center;">
        <h2>partially-encryptedNFT: Proof-of-concept of partial encryption of image NFTs, to display and prove ownership
            while avoiding the artwork being stolen</h2>
    </hgroup>

</header>

<main class="container" style="padding-top: 0px;">
Success! Click <a href="/display.html?ipfs=` + ipfsHash + `">here</a> to see your image
</main>
</body>
</html>
`
}
