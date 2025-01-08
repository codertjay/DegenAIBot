package helper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/dghubble/oauth1"
	"io"
	"log"
	"net/http"
)

func (h *Helper) twitterHelper(ctx context.Context, payload []byte, url string, requestMethod string) ([]byte, int, error) {

	config := oauth1.NewConfig(h.cfg.TwitterAPIKey, h.cfg.TwitterAPISecretKey)
	token := &oauth1.Token{Token: h.cfg.TwitterAccessToken, TokenSecret: h.cfg.TwitterAccessTokenSecret}
	httpClient := config.Client(oauth1.NoContext, token)

	req, _ := http.NewRequest(requestMethod, url, bytes.NewBuffer(payload))

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", h.cfg.TwitterBearerToken))

	res, err := httpClient.Do(req)
	if err != nil {
		return nil, 0, err
	}

	log.Println("Response headers:", res.Header)

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, res.StatusCode, err
	}

	log.Println(string(body))
	return body, res.StatusCode, nil
}

func (h *Helper) SendTweet(tweet string) (err error) {

	if len(tweet) > 250 {
		tweet = tweet[:250]
	}

	//Create a map to represent the payload
	payloadData := map[string]string{
		"text": tweet,
	}

	log.Println(tweet)

	// Properly marshal the payload into JSON
	payload, err := json.Marshal(payloadData)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	_, _, err = h.twitterHelper(context.Background(), payload, h.cfg.TwitterSendTweetRoute, "POST")
	if err != nil {
		return err
	}
	return nil
}
