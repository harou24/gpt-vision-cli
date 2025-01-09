# GPT Vision CLI

GPT Vision CLI is a command-line tool that leverages OpenAI's GPT-4 Vision capabilities to analyze images and provide detailed insights based on user prompts. This tool allows you to upload images and get descriptive responses about the content of the images directly from the command line.

## Features

- Analyze images using GPT-4 Vision capabilities.
- Provide user prompts to get specific information about the image content.
- Supports API key configuration via CLI flags or `.env` file for flexibility.

## Prerequisites

- Go (version 1.16+)
- OpenAI API key

## Quick run

Clone the repository, install the dependencies, and run:
   ```bash
   git clone https://github.com/yourusername/gpt-vision-cli.git
   cd gpt-vision-cli
   go mod tidy
   go run main.go analyze --image ./example.jpg --prompt "What objects are visible in this image?" --apikey your_openai_api_key
   ```

## Configuration

Copy **.env.example** to **.env** and set your OpenAI API key in the .env file.

### Build and run

```bash
go build -o gpt-vision-cli
./gpt-vision-cli analyze --image ./example.jpg --prompt "What objects are visible in this image?"
```
