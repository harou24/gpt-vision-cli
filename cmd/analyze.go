package cmd

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var imagePath string
var userPrompt string
var apiKey string

// analyzeCmd represents the analyze command
var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyze an image using GPT-4 Vision capabilities",
	Long:  `Provide an image and a user prompt to analyze the image content using GPT-4 Vision capabilities.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Load environment variables from .env file
		_ = godotenv.Load()

		if apiKey == "" {
			apiKey = os.Getenv("OPENAI_API_KEY")
		}
		if apiKey == "" {
			fmt.Println("API key is required. Provide it via the --apikey flag or set it in the .env file.")
			return
		}

		if imagePath == "" || userPrompt == "" {
			fmt.Println("Image path and user prompt are required")
			return
		}

		// Read and encode image
		base64Image, err := encodeImageToBase64(imagePath)
		if err != nil {
			fmt.Println("Error encoding image:", err)
			return
		}

		// Call GPT Vision API
		response, err := callGPTVision(base64Image, userPrompt, apiKey)
		if err != nil {
			fmt.Println("Error calling GPT Vision API:", err)
			return
		}

		// Print the response
		fmt.Println("GPT Vision response:", response)
	},
}

func encodeImageToBase64(imagePath string) (string, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	imageData, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(imageData), nil
}

func callGPTVision(base64Image, prompt, apiKey string) (string, error) {
	url := "https://api.openai.com/v1/chat/completions"

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", apiKey),
	}

	payload := map[string]interface{}{
		"model": "gpt-4-turbo",
		"messages": []map[string]interface{}{
			{
				"role": "user",
				"content": []map[string]interface{}{
					{"type": "text", "text": prompt},
					{"type": "image_url", "image_url": map[string]string{"url": fmt.Sprintf("data:image/jpeg;base64,%s", base64Image)}},
				},
			},
		},
		"max_tokens": 300,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "", err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		var errResponse map[string]interface{}
		json.Unmarshal(body, &errResponse)
		return "", fmt.Errorf("error from API: %v", errResponse)
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v", result), nil
}

func init() {
	rootCmd.AddCommand(analyzeCmd)
	analyzeCmd.Flags().StringVarP(&imagePath, "image", "i", "", "Path to the image file")
	analyzeCmd.Flags().StringVarP(&userPrompt, "prompt", "p", "", "User prompt message")
	analyzeCmd.Flags().StringVarP(&apiKey, "apikey", "k", "", "OpenAI API key")
}
