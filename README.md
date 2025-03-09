# yankyt

_yankyt_ is a lightweight, Go-based script that downloads videos from a YouTube playlist using [yt-dlp](https://github.com/yt-dlp/yt-dlp) and [aria2c](https://aria2.github.io/). It ensures that videos are not re-downloaded if they already exist in the specified output directory by checking file hashes. Once all downloads are complete, it sends a desktop notification (on Linux using `notify-send`).

## Features

- **Playlist Downloads:** Download entire YouTube playlists with ease.
- **Accelerated Downloads:** Uses `aria2c` for faster and more reliable downloads.
- **Avoid Re-downloads:** Checks for existing files by comparing hashes, so only new videos are downloaded.
- **Desktop Notifications:** Notifies you upon completion of the download process.

## Requirements

- **Go:** Version 1.16 or later.
- **yt-dlp:** Download from [yt-dlp GitHub](https://github.com/yt-dlp/yt-dlp) or install via your package manager.
- **aria2c:** Download from [aria2 GitHub](https://github.com/aria2/aria2) or install via your package manager.
- **notify-send:** For desktop notifications on Linux.

### Installation on Arch-based Systems (me btw)

```sh
yay -S aria2 notify-send
pipx install yt-dlp
```

### Installation on Debian-based Systems

```sh
sudo apt update
sudo apt install yt-dlp aria2 notify-send
```

## Installation

### Using `go install`

You can install _yankyt_ directly using the `go install` command:

```sh
go install github.com/ahmydyasser/yankyt
```

If you want to install a specific version, tag your release (e.g., `v1.0.0`) and then install it like so:

```sh
go install github.com/ahmydyasser/yankyt@v1.0.0
```

### From Source

Alternatively, you can clone the repository and run the script:

```sh
git clone https://github.com/ahmydyasser/yankyt.git
cd yankyt
go run main.go -u "https://www.youtube.com/playlist?list=YOUR_PLAYLIST_ID" -o "/path/to/output/directory"
```

## Usage

After installation, you can use the tool with the following command-line flags:

- `-u` or `--url`: URL of the YouTube playlist.
- `-o` or `--output`: Path to the output directory.

Example:

```sh
yankyt -u "https://www.youtube.com/playlist?list=YOUR_PLAYLIST_ID" -o "/path/to/output/directory"
```

### Help Command

You can check the available options using:

```sh
yankyt -h
```

Output:

```
Usage of yankyt:
  -o string
        Output directory (default: current directory)
  -u string
        URL of the YouTube playlist
```

## How It Works

1. **Parsing Flags:** The script accepts command-line flags to specify the YouTube playlist URL and the output directory.
2. **Fetching Playlist Info:** Uses `yt-dlp` to fetch playlist metadata in JSON format.
3. **Directory Management:** Creates a directory based on the playlist title and channel.
4. **Hash Checking:** Scans the directory for existing files and computes their hashes to prevent duplicate downloads.
5. **Downloading Videos:** Downloads missing videos using `yt-dlp` for metadata extraction and `aria2c` for actual download tasks.
6. **Notification:** Sends a desktop notification once the download process is complete.

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository.
2. Create a new branch: `git checkout -b feature/your-feature`.
3. Commit your changes: `git commit -am 'Add new feature'`.
4. Push to the branch: `git push origin feature/your-feature`.
5. Open a pull request.

## License

This project is licensed under the [MIT License](LICENSE).

## Contact

For any questions or suggestions, please open an issue on the [GitHub repository](https://github.com/ahmydyasser/yankyt).
