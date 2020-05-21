# Image resizer service
Service resize to preview-format (thumbnail) jpg and png images
POST request should contains json object
string `json:"id"`      SHA1 of image
uint   `json:"size"`    convert to size
string `json:"body"`    encoded bytes of image to base64
string `json:"format"`  extension support jpeg, jpg, png





