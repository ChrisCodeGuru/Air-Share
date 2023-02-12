package api

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/IsaacKoh88/infosecurity_project/back-end/utils/cryptography"
	"github.com/IsaacKoh88/infosecurity_project/back-end/utils/detection"
	"github.com/IsaacKoh88/infosecurity_project/back-end/utils/logging"
	"github.com/IsaacKoh88/infosecurity_project/back-end/utils/sensitive"
	"github.com/IsaacKoh88/infosecurity_project/back-end/utils/token"
	"github.com/gorilla/mux"
)

func UploadFile(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	// Get route variables
	directory := mux.Vars(r)["directory"]
	// Validate parent id is acceptable
	if directory != "root" {
		matched, _ := regexp.MatchString(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`, directory)
		if !matched {
			// throw 400 bad request
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("request parameters are incorrect"))
			return
		}
	}

	// Extract JWT from cookie
	cookie_token, err := token.ExtractToken(r)
	if err != nil {
		// throw 401 error if no token
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("you are not authorised to upload files"))
		return
	}

	// Get JWT claims
	JWTClaims, err := token.Claims(cookie_token)
	if err != nil {
		fmt.Println(err)
		return
	}

	var csrfreturns string
	t := time.Now().UTC()

	tokenerr := db.QueryRow("SELECT csrf_check($1)", JWTClaims.UUID).Scan(&csrfreturns)
	if tokenerr != nil {
		logging.ErrorLogger(tokenerr, r)
		fmt.Println(tokenerr)
	}

	csrfToken := strings.Split(csrfreturns[1:len(csrfreturns)-1], ",")[0]
	yyyy, _ := strconv.Atoi(strings.Split(csrfreturns[1:len(csrfreturns)-1], ",")[1])
	mmmm, _ := strconv.Atoi(strings.Split(csrfreturns[1:len(csrfreturns)-1], ",")[2])
	dd, _ := strconv.Atoi(strings.Split(csrfreturns[1:len(csrfreturns)-1], ",")[3])
	hhhh, _ := strconv.Atoi(strings.Split(csrfreturns[1:len(csrfreturns)-1], ",")[4])
	mm, _ := strconv.Atoi(strings.Split(csrfreturns[1:len(csrfreturns)-1], ",")[5])
	ss, _ := strconv.Atoi(strings.Split(csrfreturns[1:len(csrfreturns)-1], ",")[6])
	fullmmmm := time.Month(mmmm)

	if len(r.Header["Origin"]) == 0 || r.Header["Origin"][0] != "https://localhost" || len(r.Header["Referer"]) == 0 || r.Header["Referer"][0] != "https://localhost/fileshare" || len(r.Header["X-Csrf-Token"]) == 0 {
		w.WriteHeader(401)
		fmt.Println("CSRF")
		return
	} else if t.After(time.Date(yyyy, fullmmmm, dd, hhhh, mm, ss, 0, time.UTC)) || csrfToken == "" || r.Header["X-Csrf-Token"][0] != csrfToken {
		w.WriteHeader(401)
		fmt.Println("CSRF")
		return
	}

	// Check if user has permissions to directory
	var permisisonLevel int

	if directory != "root" {
		// Check if directory exists
		var count int

		existenceerr := db.QueryRow("SELECT COUNT(*) FROM folder WHERE id=$1", directory).Scan(&count)
		if existenceerr != nil {
			fmt.Println(existenceerr)
		}

		if count == 0 {
			// throw 404 not found
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("folder not found"))
			return
		}

		permissionerr := db.QueryRow("SELECT permission FROM account_folder WHERE account_id=$1 AND folder_id=$2", JWTClaims.UUID, directory).Scan(&permisisonLevel)
		if permissionerr != nil {
			// throw 500 internal server error
			w.WriteHeader(http.StatusInternalServerError)
			logging.ErrorLogger(permissionerr, r)
			fmt.Println(permissionerr)
			return
		}

		if permisisonLevel < 2 {
			// throw 401 unauthorised
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("you are not authorised to upload file to this folder"))
			return
		}
	}

	// Maximum upload of 10 MB file
	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		// throw 413 request too large
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		w.Write([]byte("object is too large"))
		return
	}

	// Get handler for filename, size and headers
	file, handler, err := r.FormFile("selectedFile")
	if err != nil {
		// throw 400 bad request
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// Check file types
	fileTypeList := []string{
		"application/pdf",
		"video/mp4",
		"audio/mpeg",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"application/msword",
		"image/png",
		"image/jpeg",
		"text/csv",
		"text/plain",
	}

	// Verify file type
	fileType, exists := handler.Header["Content-Type"]

	if exists {
		allowed := false
		for _, b := range fileTypeList {
			if b == fileType[0] {
				allowed = true
			}
		}

		if !allowed {
			// throw 400 bad request
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("file type not allowed"))
			return
		}
	} else {
		// throw 400 bad request
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("request parameters are incorrect"))
		return
	}

	// Extract file bytes
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		// throw 500 internal server error
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	//Hash testing
	sha256filehash := sha256.Sum256(fileBytes)

	// Generate 32 byte AES key
	aesKey := make([]byte, 32)
	_, generatekeyerr := rand.Read(aesKey)
	if generatekeyerr != nil {
		// throw 500 internal server error
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	// Create AES ciphers and encrypt file
	c, err := aes.NewCipher(aesKey)
	if err != nil {
		fmt.Println(err)
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		fmt.Println(err)
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		fmt.Println(err)
	}

	// Encrypt AES key
	encryptedAESKey, err := cryptography.Encrypt(aesKey)
	if err != nil {
		// throw 500 internal server error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Add file to database
	var created_id string

	if (fileType[0] == "image/png" || fileType[0] == "image/jpeg") && handler.Size < 5242880 {
		// detect text in images using aws rekognition if object file is an image
		imageText, err := detection.ImageOCR(fileBytes)
		if err != nil {
			// throw 500 internal server error
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println(err)
			return
		}

		// extract sensitive data from image text
		sensitiveData := detection.ImageSensitiveData(imageText.TextDetections)

		if len(sensitiveData) > 0 && directory != "root" && permisisonLevel != 4 {
			// throw 405 unauthorised
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Println("uploading sensitive data to another user's directory not allowed")
			return
		}

		if directory == "root" {
			filewriteerr := db.QueryRow("SELECT create_file($1, $2, $3, $4, $5)", JWTClaims.UUID, handler.Filename, encryptedAESKey.CiphertextBlob, encryptedAESKey.KeyId, fmt.Sprintf("%x", sha256filehash)).Scan(&created_id)
			if filewriteerr != nil {
				// throw 500 internal server error
				w.WriteHeader(http.StatusInternalServerError)
				logging.ErrorLogger(filewriteerr, r)
				fmt.Println(filewriteerr)
				return
			}
		} else {
			filewriteerr := db.QueryRow("SELECT create_file($1, $2, $3, $4, $5, $6)", JWTClaims.UUID, handler.Filename, encryptedAESKey.CiphertextBlob, encryptedAESKey.KeyId, fmt.Sprintf("%x", sha256filehash), directory).Scan(&created_id)
			if filewriteerr != nil {
				// throw 500 internal server error
				w.WriteHeader(http.StatusInternalServerError)
				logging.ErrorLogger(filewriteerr, r)
				fmt.Println(filewriteerr)
				return
			}
		}

		if len(sensitiveData) > 0 {
			// set image sensitivity in database
			var fileName string
			filewriteerr := db.QueryRow("SELECT sensitive_update($1)", created_id).Scan(&fileName)
			if filewriteerr != nil {
				// throw 500 internal server error
				w.WriteHeader(http.StatusInternalServerError)
				logging.ErrorLogger(filewriteerr, r)
				fmt.Println(filewriteerr)
				return
			}

			// get image config
			_, err = file.Seek(0, 0)
			if err != nil {
				// throw 500 internal server error
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Println("file seek error")
				return
			}
			imgConfig, _, err := image.DecodeConfig(file)
			if err != nil {
				// throw 500 internal server error
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Println("get image config err:", err)
				return
			}

			// decode image to interact with it
			_, err = file.Seek(0, 0)
			if err != nil {
				// throw 500 internal server error
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Println("file seek error")
				return
			}
			img, _, err := image.Decode(file)
			if err != nil {
				// throw 500 internal server error
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Println("image decoding err:", err)
				return
			}

			censoredImage := sensitive.NewUserImg(img)

			// iterate through all sensitive data
			for _, item := range sensitiveData {
				xstart := int(float64(imgConfig.Width)*float64(*item.Geometry.BoundingBox.Left) + 0.5)                                 // calculate horizontal start
				xend := int(float64(imgConfig.Width)*float64(*item.Geometry.BoundingBox.Left+*item.Geometry.BoundingBox.Width) + 0.5)  // calculate horizontal end
				ystart := int(float64(imgConfig.Height)*float64(*item.Geometry.BoundingBox.Top) + 0.5)                                 // calculate vertical start
				yend := int(float64(imgConfig.Height)*float64(*item.Geometry.BoundingBox.Top+*item.Geometry.BoundingBox.Height) + 0.5) // calculate vertical end

				// iterate through all pixels on the x axis
				for xpixel := xstart; xpixel <= xend; xpixel++ {
					//iterate through all pixels on the y axis
					for ypixel := ystart; ypixel <= yend; ypixel++ {
						// set pixel to black
						censoredImage.Set(xpixel, ypixel, color.RGBA{0, 0, 0, 255})
					}
				}
			}

			// create censored image file
			censoredfile, err := os.Create(fmt.Sprintf("./objects/%s-censored", created_id))
			if err != nil {
				// throw 500 internal server eror
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Println(err)
				return
			}
			defer censoredfile.Close()

			if fileType[0] == "image/jpeg" {
				var b bytes.Buffer

				// encode censored image into buffer
				encodingerr := jpeg.Encode(&b, censoredImage, nil)
				if encodingerr != nil {
					//throw 500 internal server error
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Println("writing to buffer:", encodingerr)
					return
				}

				// encrypt censored image
				encryptedCensored := gcm.Seal(nonce, nonce, b.Bytes(), nil)

				censoredfile.Write(encryptedCensored)

			} else if fileType[0] == "image/png" {
				var b bytes.Buffer

				// encode censored image into buffer
				encodingerr := png.Encode(&b, censoredImage)
				if encodingerr != nil {
					//throw 500 internal server error
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Println("writing to buffer:", encodingerr)
					return
				}

				// encrypt censored image
				encryptedCensored := gcm.Seal(nonce, nonce, b.Bytes(), nil)

				censoredfile.Write(encryptedCensored)
			}
		}
	} else if fileType[0] == "text/plain" {
		// Detect sensitive data in txt file
		fileData := fileBytes

		if detection.SensitiveOccurrences(string(fileData)) > 0 && directory != "root" && permisisonLevel != 4 {
			// throw 405 unauthorised
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Println("uploading sensitive data to another user's directory not allowed")
			return
		}

		if directory == "root" {
			filewriteerr := db.QueryRow("SELECT create_file($1, $2, $3, $4, $5)", JWTClaims.UUID, handler.Filename, encryptedAESKey.CiphertextBlob, encryptedAESKey.KeyId, fmt.Sprintf("%x", sha256filehash)).Scan(&created_id)
			if filewriteerr != nil {
				// throw 500 internal server error
				w.WriteHeader(http.StatusInternalServerError)
				logging.ErrorLogger(filewriteerr, r)
				fmt.Println(filewriteerr)
				return
			}
		} else {
			filewriteerr := db.QueryRow("SELECT create_file($1, $2, $3, $4, $5, $6)", JWTClaims.UUID, handler.Filename, encryptedAESKey.CiphertextBlob, encryptedAESKey.KeyId, fmt.Sprintf("%x", sha256filehash), directory).Scan(&created_id)
			if filewriteerr != nil {
				// throw 500 internal server error
				w.WriteHeader(http.StatusInternalServerError)
				logging.ErrorLogger(filewriteerr, r)
				fmt.Println(filewriteerr)
				return
			}
		}

		if detection.SensitiveOccurrences(string(fileData)) > 0 {
			// set image sensitivity in database
			var fileName string
			filewriteerr := db.QueryRow("SELECT sensitive_update($1)", created_id).Scan(&fileName)
			if filewriteerr != nil {
				// throw 500 internal server error
				w.WriteHeader(http.StatusInternalServerError)
				logging.ErrorLogger(filewriteerr, r)
				fmt.Println(filewriteerr)
				return
			}

			// redact sensitive data
			fileData = detection.Sensitive.USLicensePlate.ReplaceAll(fileData, []byte("***redacted***"))
			fileData = detection.Sensitive.SocialSecurityNumber.ReplaceAll(fileData, []byte("***redacted***"))
			fileData = detection.Sensitive.PaymentCard.ReplaceAll(fileData, []byte("***redacted***"))
			fileData = detection.Sensitive.NRIC.ReplaceAll(fileData, []byte("***redacted***"))

			// create censored text file
			censoredfile, err := os.Create(fmt.Sprintf("./objects/%s-censored", created_id))
			if err != nil {
				// throw 500 internal server eror
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Println(err)
				return
			}
			defer censoredfile.Close()

			// encrypt censored image
			encryptedCensored := gcm.Seal(nonce, nonce, fileData, nil)

			censoredfile.Write(encryptedCensored)
		}
	} else {
		if directory == "root" {
			filewriteerr := db.QueryRow("SELECT create_file($1, $2, $3, $4, $5)", JWTClaims.UUID, handler.Filename, encryptedAESKey.CiphertextBlob, encryptedAESKey.KeyId, fmt.Sprintf("%x", sha256filehash)).Scan(&created_id)
			if filewriteerr != nil {
				// throw 500 internal server error
				w.WriteHeader(http.StatusInternalServerError)
				logging.ErrorLogger(filewriteerr, r)
				fmt.Println(filewriteerr)
				return
			}
		} else {
			filewriteerr := db.QueryRow("SELECT create_file($1, $2, $3, $4, $5, $6)", JWTClaims.UUID, handler.Filename, encryptedAESKey.CiphertextBlob, encryptedAESKey.KeyId, fmt.Sprintf("%x", sha256filehash), directory).Scan(&created_id)
			if filewriteerr != nil {
				// throw 500 internal server error
				w.WriteHeader(http.StatusInternalServerError)
				logging.ErrorLogger(filewriteerr, r)
				fmt.Println(filewriteerr)
				return
			}
		}
	}

	// File Storing
	myfile, err := os.Create("./objects/" + created_id)
	if err != nil {
		//throw 500 internal server error
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	defer myfile.Close()

	encryptedBytes := gcm.Seal(nonce, nonce, fileBytes, nil)

	// write this byte array to our temporary file
	myfile.Write(encryptedBytes)
	fmt.Printf("%x\n", sha256filehash)

	var email string
	accounterr := db.QueryRow("SELECT get_user_email($1)", JWTClaims.UUID).Scan(&email)
	if accounterr != nil {
		logging.ErrorLogger(accounterr, r)
		fmt.Println(accounterr)
	}

	folderDetails := map[string]interface{}{
		"Email":     email,
		"File Name": handler.Filename,
		"File Size": handler.Size,
		"MIME type": handler.Header["Content-Type"],
	}
	folderDetailsJSON, err := json.Marshal(folderDetails)
	if err != nil {
		logging.ErrorLogger(err, r)
		panic(err)
	}
	logging.Logger(string(folderDetailsJSON), 201, r)

	w.WriteHeader(201)
	fmt.Fprintf(w, "Successfully Uploaded File\n")
}
