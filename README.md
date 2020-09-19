## LSB Encoding

CLI application that hides(embeds) texts inside images

Inspired by: [CryptoTheLlama Youtube Video](https://www.youtube.com/watch?v=o_uXEete0K4&ab_channel=CryptoTheLlama)

## Installation:
### If you have go environment:
> go get github.com/DavudSafarli/LSB-encoding

If your added GOPATH/bin to PATH env variable, you should be able to use `LSB-encoding` command

### If you don't:
Just download the executable from inside repository files, and place it wherever you want



## How it works:
### embedding a secret inside image:

> LSB-encoding.exe embed -text "this is a secret text" -path "./image-path-to-embed.png"

creates a new image with <image-path-to-embed_embedded.png> in the same directory


### extracting secret from image:
> LSB-encoding.exe extract -path "./image-path-to-embed_embedded.png"

prints the secret hidden inside the image

### Sidenote
embedding only works with PNG images, because JPG images are compressed and pixel values are modified by decoding algorithms





