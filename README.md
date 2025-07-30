# TeleFilesBot

![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/Argh94/TeleFilesBot/main.yml?branch=main&style=flat-square)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/Argh94/TeleFilesBot?style=flat-square)
![GitHub License](https://img.shields.io/github/license/Argh94/TeleFilesBot?style=flat-square)
![GitHub Release](https://img.shields.io/github/v/release/Argh94/TeleFilesBot?style=flat-square)
![Telegram](https://img.shields.io/badge/Telegram-Bot-blue?style=flat-square&logo=telegram)

**TeleFilesBot** is a powerful and easy-to-use Telegram bot that allows users to upload, manage, and download files, images, and videos. The bot is developed in Go using the `go-telegram-bot-api` library, providing features such as uploading, listing, deleting, and downloading files.

---

## Important Note
Wherever you see `Argh94/TeleFilesBot`, replace it with your own GitHub username and repository name.

---

## Features
- 📤 **File Upload:** Send and save files, images, and videos.
- 📋 **List Files:** Get a list of uploaded files with download links.
- 🗑️ **Delete Files:** Delete files by replying to messages or using a command.
- 📥 **Download Files:** Download files via command or link.
- 🔐 **Secure Storage:** Files are stored in a private Telegram group.
- ⚙️ **Proxy Support:** Optional proxy support for Telegram API connectivity.
- 🚀 **Automatic Binary Build:** Build binaries for multiple OS/architectures with GitHub Actions.

---

## Prerequisites
- **Go** version 1.18 or higher ([Download Go](https://golang.org/dl/))
- **Telegram bot token** (from [BotFather](https://t.me/BotFather))
- **Private Telegram group** for file storage
- **GitHub Account** (optional, for using GitHub Actions)
- **Operating System:** Linux, Windows, or macOS

---

## Installation & Setup

### 1. Clone the Repository
```bash
git clone https://github.com/Argh94/TeleFilesBot.git
cd TeleFilesBot
```

### 2. Create the Configuration File

In the project root, create a file named `config.json` with the following content:

```json
{
  "botToken": "YOUR_BOT_TOKEN",
  "botUsername": "@YourBotUsername",
  "privateChatID": 123456789,
  "cacheFilePath": "cache.json"
}
```
- **botToken**: The token you received from BotFather.
- **botUsername**: Your bot's username (e.g., @MyTeleFilesBot).
- **privateChatID**: The numeric ID of your private Telegram group (to get it, add the bot to the group and use a bot like [@GetIDsBot](https://t.me/GetIDsBot)).
- **cacheFilePath**: Path to the cache file (default is `cache.json`).

### 3. Install Dependencies
```bash
go mod download
```

### 4. Build the Project
To build a binary for your desired OS and architecture (example: Linux 64-bit):

```bash
GOOS=linux GOARCH=amd64 go build -o output/linux_amd64/TeleFilesBot
```
Replace `linux` and `amd64` with your target OS (`linux`, `windows`, `darwin`) and architecture (`amd64`, `arm64`).

### 5. Run the Bot
For Linux:
```bash
./output/linux_amd64/TeleFilesBot
```
For Windows:
```bash
.\output\windows_amd64\TeleFilesBot.exe
```

### 6. Using GitHub Actions (Optional)
This project supports automatic binary builds via GitHub Actions. The workflow file is at `.github/workflows/main.yml`.
1. Push your project to GitHub.
2. Go to the **Actions** tab in your repository and enable the workflow.
3. After successful runs, download the binary artifacts from the **Actions** section.

---

## Bot Commands

| Command | Description |
|---------|-------------|
| /start | Show welcome message and custom keyboard |
| /list | List uploaded files with download links |
| /delete <file_name> | Delete a file by name or by replying to its message |
| /start download_<file_name> | Download a file by name |

---

## Usage Examples

- **File Upload:**  
  Send a file, image, or video to the bot. The bot will reply with "File saved successfully".

- **List Files:**  
  Send `/list` to receive the list of files with download links.

- **Delete File:**  
  Reply to the file message with `/delete`, or use `/delete file_name` directly.

- **Download File:**  
  Click the download link in the file list or send `/start download_file_name`.

---

## Project Structure

```
TeleFilesBot/
├── main.go              # Bot main code
├── config.json          # Configuration file
├── cache.json           # File info cache
├── output/              # Built binaries output
├── go.mod               # Go module
├── go.sum               # Dependencies
└── .github/workflows/   # GitHub Actions workflows
    └── build.yml
```

---

## Troubleshooting

- **Permission denied in GitHub Actions:**  
  Ensure cache paths are properly cleaned. The provided `build.yml` handles this.

- **BotAPI initialization failed:**  
  Check your bot token in `config.json`. If using a proxy, verify its configuration.

- **File exists in cache:**  
  Make sure the cache cleaning step in `build.yml` runs before restoring the cache.

---

## Contributing

- Fork the repository
- Make your changes and create a pull request
- Ensure your code matches Go standards and passes all tests

---

## License

This project is released under the [MIT License](LICENSE).

---

## Contact

For questions or support, use GitHub Issues or contact via Telegram.

---

Made with ❤️ by Argh94
