package helper

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"strings"

	"gakuren-system.com/pkg/model"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func ReturnResponse(c fiber.Ctx, statusCode int, message string, data interface{}, err error) error {
	response := make(map[string]interface{})
	response["message"] = message
	response["data"] = nil
	response["error"] = nil

	if data != nil {
		response["data"] = data
	}
	if err != nil {
		log.Println(err)
		response["error"] = err.Error()
	}
	return c.Status(statusCode).JSON(response)
}

func ValidateRequest(c fiber.Ctx) (uuid.UUID, uuid.UUID, error) {
	tenantUUID, err := uuid.Parse(c.Get("tenant_uuid"))
	if err != nil {
		return uuid.Nil, uuid.Nil, errors.New("invalid or missing tenant UUID")
	}
	userUUID, err := GetUserUUIDByAccessToken(c.Get("Authorization"))
	if err != nil || userUUID == nil {
		if err == nil {
			err = errors.New("missing user UUID")
		}
		return uuid.Nil, uuid.Nil, err
	}
	return tenantUUID, *userUUID, nil
}

func GetUserUUIDByAccessToken(authHeader string) (*uuid.UUID, error) {
	// authHeader := c.Get("Authorization")

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, errors.New("Missing or invalid token")
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if strings.TrimSpace(tokenString) == "" {
		return nil, errors.New("token is required")
	}

	publicKey, err := ParsePublicKey(os.Getenv("RSA_PUBLIC_KEY"))
	if err != nil {
		return nil, fmt.Errorf("parse RSA public key: %w", err)
	}

	claims := new(model.AccessTokenClaims)
	_, err = jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	if claims.RegisteredClaims.Subject == "" {
		return nil, errors.New("missing subject in token claims")
	}

	subjectUUID, err := uuid.Parse(claims.RegisteredClaims.Subject)
	if err != nil {
		return nil, fmt.Errorf("parse subject as UUID: %w", err)
	}
	return &subjectUUID, nil
}

func DecodeB64Bytes(src string) ([]byte, error) {
	src = strings.TrimSpace(src)
	if src == "" {
		return nil, errors.New("empty base64 input")
	}

	b, err := base64.StdEncoding.DecodeString(src)
	if err == nil {
		return b, nil
	}

	b, err = base64.URLEncoding.DecodeString(src)
	if err == nil {
		return b, nil
	}

	return nil, err
}

func CalculateDataStatisticResult(count *model.CountResult, payload model.SearchPayload, totalData int) model.DataStatistics {
	var dataStat model.DataStatistics
	dataStat.TotalRow = count.Count
	if payload.RowPerPage != nil {
		dataStat.RowPerPage = *payload.RowPerPage
	} else {
		dataStat.RowPerPage = DEFAULT_ROW_PER_PAGES
	}
	if payload.Page != nil {
		dataStat.CurrentPage = *payload.Page
	} else {
		dataStat.CurrentPage = DEFAULT_PAGES
	}
	dataStat.MaxPage = int(math.Ceil(float64(count.Count) / float64(dataStat.RowPerPage)))
	dataStat.CurrentRow = totalData
	if (dataStat.CurrentPage * dataStat.RowPerPage) > dataStat.TotalRow {
		dataStat.StartRow = dataStat.TotalRow - dataStat.CurrentRow + 1
	} else {
		dataStat.StartRow = (dataStat.CurrentPage * dataStat.RowPerPage) - (dataStat.RowPerPage - 1)
	}
	if dataStat.CurrentPage*dataStat.RowPerPage < dataStat.TotalRow {
		dataStat.EndRow = dataStat.CurrentPage * dataStat.RowPerPage
	} else {
		dataStat.EndRow = dataStat.TotalRow
	}

	return dataStat
}
