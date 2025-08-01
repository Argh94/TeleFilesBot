package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net"
    "net/http"
    "net/url"
    "path/filepath"
    "strings"
    "time"

    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Config struct {
    BotToken      string `json:"botToken"`
    BotUsername   string `json:"botUsername"`
    PrivateChatID int64  `json:"privateChatID"`
    CacheFilePath string `json:"cacheFilePath"`
}

var (
    config    Config
    fileStore = make(map[string]string)
    useProxy  = true
)

func main() {
    log.SetFlags(log.LstdFlags | log.Lshortfile)
    log.Println("Start the robot...")

    if err := loadConfig("config.json"); err != nil {
        log.Fatalf("Error loading configuration: %v", err)
    }

    loadCache()

    var httpClient *http.Client
    if useProxy {
        proxyURL, err := url.Parse("http://127.0.0.1:7890")
        if err != nil {
            log.Fatalf("Proxy URL parsing failed: %v", err)
        }
        httpClient = &http.Client{
            Transport: &http.Transport{
                Proxy: http.ProxyURL(proxyURL),
                DialContext: (&net.Dialer{
                    Timeout:   5 * time.Second,
                    KeepAlive: 5 * time.Second,
                    DualStack: false,
                }).DialContext,
            },
        }
        log.Println("Proxy Enabled")
    } else {
        httpClient = http.DefaultClient
        log.Println("Proxy Disabled")
    }

    bot, err := tgbotapi.NewBotAPIWithClient(config.BotToken, tgbotapi.APIEndpoint, httpClient)
    if err != nil {
        log.Fatalf("BotAPI initialization failed: %v", err)
    }

    bot.Debug = false
    log.Printf("Successfully logged into the robot: %s", bot.Self.UserName)

    u := tgbotapi.NewUpdate(0)
    u.Timeout = 60

    updates := bot.GetUpdatesChan(u)

    for update := range updates {
        if update.Message == nil {
            continue
        }

        if update.Message.Command() == "start" && strings.HasPrefix(update.Message.CommandArguments(), "download_") {
            fileName := strings.ReplaceAll(strings.TrimPrefix(update.Message.CommandArguments(), "download_"), "_", " ")
            downloadFile(bot, update.Message.Chat.ID, fileName)
            continue
        }

        command := strings.TrimSpace(update.Message.Text)
        switch command {
        case "/start", "Help":
            sendCustomKeyboard(bot, update.Message.Chat.ID)
            msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Welcome to ZSNET Bot!\nSend me a file to upload\nUse /list to view the file list\nUse /delete to delete files")
            sendMessageWithLog(bot, msg, "Welcome message sent successfully")

        case "/list", "My Files":
            if len(fileStore) == 0 {
                msg := tgbotapi.NewMessage(update.Message.Chat.ID, "No files found")
                sendMessageWithLog(bot, msg, "Send empty file list message")
                continue
            }

            fileList := "File list:\n"
            i := 1
            for fileName := range fileStore {
                nameWithoutExt := strings.TrimSuffix(fileName, filepath.Ext(fileName))
                escapedFileName := escapeMarkdownV2(nameWithoutExt)
                downloadLink := fmt.Sprintf("https://t.me/%s?start=download_%s", config.BotUsername, strings.ReplaceAll(escapedFileName, " ", "_"))
                fileList += fmt.Sprintf("%d [%s](%s)\n", i, escapedFileName, escapeMarkdownV2(downloadLink))
                i++
            }
            msg := tgbotapi.NewMessage(update.Message.Chat.ID, fileList)
            msg.ParseMode = "MarkdownV2"
            msg.DisableWebPagePreview = true
            sendMessageWithLog(bot, msg, "File list sent successfully")

        case "/delete", "Delete File":
            if update.Message.ReplyToMessage != nil && update.Message.ReplyToMessage.Document != nil {
                fileName := update.Message.ReplyToMessage.Document.FileName
                deleteFile(bot, update.Message.Chat.ID, fileName)
            } else {
                args := update.Message.CommandArguments()
                if args == "" {
                    msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please enter the file name to delete, or reply to a message containing a file and use /delete")
                    sendMessageWithLog(bot, msg, "Prompt the user to enter a file name or reply to a file message")
                    continue
                }
                deleteFile(bot, update.Message.Chat.ID, args)
            }

        default:
            if update.Message.Document != nil {
                fileID := update.Message.Document.FileID
                fileName := update.Message.Document.FileName
                fileStore[fileName] = fileID
                saveCache()
                forward := tgbotapi.NewForward(config.PrivateChatID, update.Message.Chat.ID, update.Message.MessageID)
                _, err := bot.Send(forward)
                if err != nil {
                    msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Failed to save file")
                    sendMessageWithLog(bot, msg, fmt.Sprintf("File save failed: %s", fileName))
                    continue
                }
                msg := tgbotapi.NewMessage(update.Message.Chat.ID, "File saved successfully")
                sendMessageWithLog(bot, msg, fmt.Sprintf("The file has been saved and forwarded to the group: %s", fileName))
            } else if update.Message.Photo != nil {
                photo := update.Message.Photo[len(update.Message.Photo)-1]
                fileID := photo.FileID
                fileName := fmt.Sprintf("photo_%d.jpg", time.Now().Unix())
                fileStore[fileName] = fileID
                saveCache()
                msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Image saved")
                sendMessageWithLog(bot, msg, fmt.Sprintf("Image saved: %s", fileName))
            } else if update.Message.Video != nil {
                fileID := update.Message.Video.FileID
                fileName := fmt.Sprintf("video_%d.mp4", time.Now().Unix())
                fileStore[fileName] = fileID
                saveCache()
                msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Video saved")
                sendMessageWithLog(bot, msg, fmt.Sprintf("Video saved: %s", fileName))
            } else {
                msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Invalid command")
                sendMessageWithLog(bot, msg, "User sent an invalid command")
            }
        }
    }
}

func loadConfig(configFile string) error {
    if _, err := ioutil.ReadFile(configFile); err != nil {
        log.Println("Config file not found, generating a new one...")

        defaultConfig := Config{
            BotToken:      "your-bot-token-here",
            BotUsername:   "your-bot-username-here",
            PrivateChatID: 123456789,
            CacheFilePath: "cache.json",
        }

        data, err := json.MarshalIndent(defaultConfig, "", "  ")
        if err != nil {
            return fmt.Errorf("failed to marshal default config: %v", err)
        }

        if err := ioutil.WriteFile(configFile, data, 0644); err != nil {
            return fmt.Errorf("failed to write default config file: %v", err)
        }

        log.Println("Default config file created successfully.")

        config = defaultConfig
        return nil
    }

    data, err := ioutil.ReadFile(configFile)
    if err != nil {
        return fmt.Errorf("failed to read config file: %v", err)
    }

    if err := json.Unmarshal(data, &config); err != nil {
        return fmt.Errorf("failed to parse config file: %v", err)
    }

    return nil
}

func sendCustomKeyboard(bot *tgbotapi.BotAPI, chatID int64) {
    keyboard := tgbotapi.NewReplyKeyboard(
        tgbotapi.NewKeyboardButtonRow(
            tgbotapi.NewKeyboardButton("Help"),
            tgbotapi.NewKeyboardButton("My Files"),
        ),
        tgbotapi.NewKeyboardButtonRow(
            tgbotapi.NewKeyboardButton("Delete File"),
            tgbotapi.NewKeyboardButton("Download File"),
        ),
    )

    msg := tgbotapi.NewMessage(chatID, "Please select an action:")
    msg.ReplyMarkup = keyboard
    if _, err := bot.Send(msg); err != nil {
        log.Printf("Failed to send custom keyboard: %v", err)
    }
}

func deleteFile(bot *tgbotapi.BotAPI, chatID int64, fileName string) {
    if _, exists := fileStore[fileName]; exists {
        delete(fileStore, fileName)
        saveCache()
        msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("File deleted: %s", fileName))
        sendMessageWithLog(bot, msg, fmt.Sprintf("File deleted successfully: %s", fileName))
    } else {
        msg := tgbotapi.NewMessage(chatID, "Specified file not found")
        sendMessageWithLog(bot, msg, fmt.Sprintf("file not found: %s", fileName))
    }
}

func downloadFile(bot *tgbotapi.BotAPI, chatID int64, fileName string) {
    var fileID string
    var exists bool
    var displayName string
    for key, value := range fileStore {
        if strings.TrimSuffix(key, filepath.Ext(key)) == fileName {
            fileID = value
            displayName = key
            exists = true
            break
        }
    }
    if exists {
        file, err := bot.GetFile(tgbotapi.FileConfig{FileID: fileID})
        if err != nil {
            log.Printf("Failed to get file: %v", err)
            return
        }
        filePath := file.FilePath
        msg := tgbotapi.NewDocument(chatID, tgbotapi.FileURL("https://api.telegram.org/file/bot" + config.BotToken + "/" + filePath))
        msg.Caption = displayName
        if _, err := bot.Send(msg); err != nil {
            log.Printf("Failed to send file: %v", err)
        }
    } else {
        msg := tgbotapi.NewMessage(chatID, "File does not exist")
        sendMessageWithLog(bot, msg, "File not found")
    }
}

func saveCache() {
    data, err := json.Marshal(fileStore)
    if err != nil {
        log.Printf("Failed to save cache: %v", err)
        return
    }

    if err := ioutil.WriteFile(config.CacheFilePath, data, 0644); err != nil {
        log.Printf("Failed to write cache file: %v", err)
    }
}

func loadCache() {
    data, err := ioutil.ReadFile(config.CacheFilePath)
    if err != nil {
        log.Println("No cache file found, continuing with an empty store.")
        return
    }

    if err := json.Unmarshal(data, &fileStore); err != nil {
        log.Printf("Failed to load cache: %v", err)
    }
}

func sendMessageWithLog(bot *tgbotapi.BotAPI, msg tgbotapi.MessageConfig, logMessage string) {
    _, err := bot.Send(msg)
    if err != nil {
        log.Printf("Error sending message: %v", err)
    }
    log.Println(logMessage)
}

func escapeMarkdownV2(input string) string {

    return strings.ReplaceAll(input, "_", "\\_")
}
